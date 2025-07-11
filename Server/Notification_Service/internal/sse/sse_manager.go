package sse

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"notification-service/internal/db" // Assuming db.Notification is your notification model
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Client represents a single SSE client.
type Client struct {
	UserID   string // To send user-specific notifications
	SendChan chan []byte
}

// SSEManager manages all active SSE client connections.
type SSEManager struct {
	mu           sync.RWMutex
	clients      map[*Client]struct{}            // Set of all active clients
	userClients  map[string]map[*Client]struct{} // Clients per user_id for targeted messages
	Broadcast    chan db.Notification            // Channel to receive new broadcast notifications
	UserSpecific chan db.Notification            // Channel to receive new user-specific notifications
}

// NewSSEManager creates and starts a new SSEManager.
func NewSSEManager() *SSEManager {
	m := &SSEManager{
		clients:      make(map[*Client]struct{}),
		userClients:  make(map[string]map[*Client]struct{}),
		Broadcast:    make(chan db.Notification, 100), // Buffered channel
		UserSpecific: make(chan db.Notification, 100), // Buffered channel
	}
	go m.dispatchNotifications() // Start a goroutine to dispatch messages
	return m
}

// RegisterClient adds a new client to the manager.
func (m *SSEManager) RegisterClient(userID string, sendChan chan []byte) *Client {
	m.mu.Lock()
	defer m.mu.Unlock()

	client := &Client{UserID: userID, SendChan: sendChan}
	m.clients[client] = struct{}{}

	if userID != "" {
		if _, ok := m.userClients[userID]; !ok {
			m.userClients[userID] = make(map[*Client]struct{})
		}
		m.userClients[userID][client] = struct{}{}
		log.Printf("SSE Manager: Registered client for user %s", userID)
	} else {
		// This case might be for clients that want ALL notifications, including all user-specific ones.
		// Or, it could be an error if a userID was expected. Adjust logic as needed.
		log.Printf("SSE Manager: Registered a client with no specific userID (will receive broadcasts).")
	}
	return client
}

// UnregisterClient removes a client from the manager.
func (m *SSEManager) UnregisterClient(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.clients, client)
	if client.UserID != "" {
		if userMap, ok := m.userClients[client.UserID]; ok {
			delete(userMap, client)
			if len(userMap) == 0 {
				delete(m.userClients, client.UserID)
			}
		}
	}
	close(client.SendChan) // Important to signal the handler to stop
	log.Printf("SSE Manager: Unregistered client for user %s", client.UserID)
}

// dispatchNotifications listens on the notification channels and sends them to appropriate clients.
func (m *SSEManager) dispatchNotifications() {
	for {
		select {
		case notification := <-m.Broadcast:
			jsonData, err := json.Marshal(notification)
			if err != nil {
				log.Printf("SSE Manager: Error marshalling broadcast notification: %v", err)
				continue
			}
			m.mu.RLock()
			for client := range m.clients { // Broadcast to ALL connected clients
				select {
				case client.SendChan <- jsonData:
				default:
					log.Printf("SSE Manager: Client %s send channel full for broadcast. Skipping.", client.UserID)
				}
			}
			m.mu.RUnlock()
			log.Printf("SSE Manager: Dispatched broadcast notification ID %s", notification.ID.String())

		case notification := <-m.UserSpecific:
			if !notification.UserID.Valid || notification.UserID.String == "" {
				log.Printf("SSE Manager: Received user-specific notification without a valid UserID. Skipping.")
				continue
			}
			userID := notification.UserID.String
			jsonData, err := json.Marshal(notification)
			if err != nil {
				log.Printf("SSE Manager: Error marshalling user-specific notification: %v", err)
				continue
			}

			m.mu.RLock()
			// Send to specific user
			if userClients, ok := m.userClients[userID]; ok {
				for client := range userClients {
					select {
					case client.SendChan <- jsonData:
					default:
						log.Printf("SSE Manager: Client %s send channel full for user-specific. Skipping.", client.UserID)
					}
				}
				log.Printf("SSE Manager: Dispatched user-specific notification ID %s to user %s", notification.ID.String(), userID)
			}
			// Also send user-specific notifications to clients registered without a UserID (if any and if desired)
			// This part depends on how you want clients without a specific UserID to behave.
			// For now, let's assume they only get explicit broadcasts.
			m.mu.RUnlock()
		}
	}
}

// StreamNotificationsHandler is the Gin handler for SSE connections.
// It's designed to stream notifications for a specific user.
func (m *SSEManager) StreamNotificationsHandler(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required for SSE stream"})
		return
	}

	// Set headers for SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Adjust for your CORS policy

	messageChan := make(chan []byte, 10) // Buffered channel for this client
	client := m.RegisterClient(userID, messageChan)
	defer m.UnregisterClient(client)

	// Notify client of connection (optional)
	_, err := c.Writer.WriteString("event: connected\ndata: Connection established for user " + userID + "\n\n")
	if err != nil {
		log.Printf("SSE Handler: Error writing connection confirmation to client %s: %v", userID, err)
		return
	}
	c.Writer.Flush()

	ctx := c.Request.Context()
	keepAliveTicker := time.NewTicker(20 * time.Second) // Send a keep-alive comment every 20s
	defer keepAliveTicker.Stop()

	log.Printf("SSE Handler: Client connected for user %s. Listening for notifications...", userID)

	for {
		select {
		case <-ctx.Done(): // Client disconnected
			log.Printf("SSE Handler: Client %s disconnected (context done).", userID)
			return
		case <-keepAliveTicker.C:
			_, err := c.Writer.WriteString(": keep-alive\n\n") // SSE comment for keep-alive
			if err != nil {
				log.Printf("SSE Handler: Error writing keep-alive to client %s: %v", userID, err)
				return // Client likely disconnected
			}
			c.Writer.Flush()
		case message, ok := <-messageChan:
			if !ok { // Channel closed, client unregistered
				log.Printf("SSE Handler: Message channel closed for client %s.", userID)
				return
			}
			_, err := c.Writer.WriteString(fmt.Sprintf("event: notification\ndata: %s\n\n", string(message)))
			if err != nil {
				log.Printf("SSE Handler: Error writing message to client %s: %v", userID, err)
				return // Client likely disconnected
			}
			c.Writer.Flush()
			log.Printf("SSE Handler: Sent notification to user %s", userID)
		}
	}
}
