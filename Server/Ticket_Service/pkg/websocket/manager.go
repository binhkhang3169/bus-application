// file: pkg/websocket/manager.go
package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"ticket-service/domain/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

const pendingWsConnectionsKey = "pending_ws_connections"

// Message định nghĩa cấu trúc tin nhắn chuẩn qua WebSocket
type Message struct {
	Type    string      `json:"type"` // "result", "error", "ack"
	Payload interface{} `json:"payload"`
}

// Client đại diện cho một kết nối WebSocket đang hoạt động
type Client struct {
	conn      *websocket.Conn
	bookingID string
	manager   *Manager
	send      chan []byte
	ack       chan bool
	ctx       context.Context    // Context để quản lý vòng đời của goroutine
	cancel    context.CancelFunc // Hàm để hủy context
}

// Manager quản lý các kết nối WebSocket và Redis client
type Manager struct {
	clients      map[string]*Client
	mu           sync.RWMutex
	redisClient  *redis.Client
	onAckTimeout func(bookingID string, ticketID string)
}

// NewManager khởi tạo một WebSocket Manager mới
func NewManager(redisClient *redis.Client, onAckTimeout func(bookingID string, ticketID string)) *Manager {
	return &Manager{
		clients:      make(map[string]*Client),
		redisClient:  redisClient,
		onAckTimeout: onAckTimeout,
	}
}

// getRedisChannel trả về tên kênh Redis được chuẩn hóa
func getRedisChannel(bookingID string) string {
	return fmt.Sprintf("booking-result:%s", bookingID)
}

// register đăng ký một client mới và khởi chạy các goroutine cần thiết
func (m *Manager) register(client *Client) {
	m.mu.Lock()
	m.clients[client.bookingID] = client
	m.mu.Unlock()

	log.Printf("[WebSocket] Client registered locally for bookingId: %s", client.bookingID)

	// Khởi chạy goroutine để ghi và lắng nghe Redis
	// readPump sẽ được chạy trên goroutine chính của handler.
	go client.writePump()
	go client.redisListenPump()
}

// unregister xóa một client khỏi manager
func (m *Manager) unregister(client *Client) {
	m.mu.Lock()
	if _, ok := m.clients[client.bookingID]; ok {
		delete(m.clients, client.bookingID)
		client.cancel() // Hủy context để dừng các goroutine liên quan
		close(client.send)
		log.Printf("[WebSocket] Client unregistered locally for bookingId: %s", client.bookingID)
	}
	m.mu.Unlock()
}

// redisListenPump: Goroutine lắng nghe tin nhắn từ kênh Redis và chuyển tiếp đến client
func (c *Client) redisListenPump() {
	channel := getRedisChannel(c.bookingID)
	pubsub := c.manager.redisClient.Subscribe(c.ctx, channel)
	defer pubsub.Close()

	log.Printf("[Redis] Subscribed to channel: %s for bookingId: %s", channel, c.bookingID)

	// --- LOGIC MỚI: KIỂM TRA LẠI TRẠNG THÁI NGAY SAU KHI SUBSCRIBE ---
	// Điều này để xử lý race condition khi consumer đã xử lý xong trước khi client kịp subscribe.
	redisStateKey := fmt.Sprintf("booking:state:%s", c.bookingID)
	stateData, err := c.manager.redisClient.HGetAll(c.ctx, redisStateKey).Result()
	if err == nil && (stateData["status"] == "COMPLETED" || stateData["status"] == "FAILED") {
		log.Printf("[Redis] Race condition detected for %s. State is already '%s'. Sending historical result.", c.bookingID, stateData["status"])

		var messageToSend Message
		if stateData["status"] == "COMPLETED" {
			var payload models.TicketReturn // Hoặc bất kỳ struct nào bạn lưu trong result_payload
			// Bạn có thể cần bỏ comment dòng dưới nếu dùng thư viện `models` trong package này
			if json.Unmarshal([]byte(stateData["result_payload"]), &payload) == nil {
				messageToSend = Message{Type: "result", Payload: payload}
			} else {
				// Fallback nếu unmarshal lỗi
				messageToSend = Message{Type: "error", Payload: map[string]string{"error": "Failed to parse historical result."}}
			}
		} else { // FAILED
			messageToSend = Message{Type: "error", Payload: map[string]string{"error": stateData["error_message"]}}
		}

		payloadBytes, _ := json.Marshal(messageToSend)
		select {
		case c.send <- payloadBytes:
			log.Printf("[Redis] Sent historical result for %s and closing listener.", c.bookingID)
		case <-c.ctx.Done():
			log.Printf("[Redis] Client disconnected before historical result could be sent for %s.", c.bookingID)
		}
		return // Kết thúc goroutine lắng nghe vì đã có kết quả cuối cùng.
	}
	// --- KẾT THÚC LOGIC MỚI ---

	// Vòng lặp lắng nghe như cũ
	ch := pubsub.Channel()

	for {
		select {
		case msg := <-ch:
			log.Printf("[Redis] Message received on channel %s, forwarding to client.", channel)
			select {
			case c.send <- []byte(msg.Payload):
			case <-c.ctx.Done():
				return // Client ngắt kết nối trong khi chờ gửi
			}
			// Một booking chỉ có một kết quả cuối cùng, sau khi nhận xong thì có thể dừng lắng nghe.
			return
		case <-c.ctx.Done():
			log.Printf("[Redis] Listener stopped for channel %s due to client disconnect.", channel)
			return
		}
	}
}

// readPump: Goroutine đọc tin nhắn từ client (chủ yếu là tin nhắn 'ack')
// **HÀM NÀY SẼ CHẶN (BLOCK) ĐỂ GIỮ KẾT NỐI**
func (c *Client) readPump() {
	defer func() {
		c.manager.unregister(c)
		c.conn.Close()
	}()
	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil })

	for {
		// conn.ReadMessage() là một lệnh blocking.
		// Nó sẽ đợi cho đến khi có tin nhắn hoặc có lỗi.
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[WebSocket] Error reading message for %s: %v", c.bookingID, err)
			}
			break // Thoát vòng lặp khi có lỗi hoặc client đóng kết nối
		}
		var msg Message
		if err := json.Unmarshal(message, &msg); err == nil && msg.Type == "ack" {
			select {
			case c.ack <- true:
			default:
			}
		}
	}
}

// writePump: Goroutine ghi tin nhắn đến client và quản lý ACK timeout
func (c *Client) writePump() {
	ackTimeoutDuration := 15 * time.Second
	pingTicker := time.NewTicker(45 * time.Second)
	var ackTimer *time.Timer

	defer func() {
		if ackTimer != nil {
			ackTimer.Stop()
		}
		pingTicker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("[WebSocket] Error writing message: %v", err)
				return
			}

			// Khởi động bộ đếm thời gian chờ ACK sau khi đã gửi tin nhắn kết quả
			if ackTimer != nil {
				ackTimer.Stop()
			}
			ackTimer = time.NewTimer(ackTimeoutDuration)

			go func(msgBytes []byte) {
				select {
				case <-ackTimer.C:
					log.Printf("[ACK] Timeout for bookingId: %s. Initiating ticket cancellation.", c.bookingID)
					var msgData Message
					if json.Unmarshal(msgBytes, &msgData) == nil && msgData.Type == "result" {
						// Giả sử payload là một map có chứa "ticket_id"
						if payloadMap, ok := msgData.Payload.(map[string]interface{}); ok {
							if ticketID, ok := payloadMap["ticket_id"].(string); ok {
								c.manager.onAckTimeout(c.bookingID, ticketID)
							}
						}
					}
					c.conn.Close()
				case <-c.ack:
					ackTimer.Stop()
					log.Printf("[ACK] Received for bookingId: %s. Closing connection.", c.bookingID)
					c.conn.Close()
				case <-c.ctx.Done():
					ackTimer.Stop()
					return
				}
			}(message)

		case <-pingTicker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}

// HandleConnection là handler cho Gin để xử lý kết nối WebSocket
// Hàm này sẽ thay thế cho helper ServeWs cũ.
func (m *Manager) HandleConnection(c *gin.Context) {
	bookingID := c.Param("bookingId")
	if bookingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bookingId is required"})
		return
	}

	removedCount, err := m.redisClient.ZRem(c.Request.Context(), pendingWsConnectionsKey, bookingID).Result()
	if err != nil {
		log.Printf("[WebSocket] Redis error defusing timeout for %s: %v", bookingID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing connection."})
		return
	}
	if removedCount == 0 {
		// ZREM trả về 0 có nghĩa là bookingId không có trong set.
		// Lý do: nó đã hết hạn và đã bị TimeoutWorker xử lý.
		log.Printf("[WebSocket] Connection rejected for %s: timeout already processed.", bookingID)
		c.JSON(http.StatusRequestTimeout, gin.H{"error": "Your booking session has expired because you did not connect in time."})
		return
	}

	log.Printf("[WebSocket] Connection for %s successful within timeout window.", bookingID)

	// --- KIỂM TRA TRẠNG THÁI BOOKING TRONG REDIS ---
	redisStateKey := fmt.Sprintf("booking:state:%s", bookingID)
	stateData, err := m.redisClient.HGetAll(c.Request.Context(), redisStateKey).Result()
	if err != nil {
		log.Printf("[WebSocket] Redis error checking state for %s: %v", bookingID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not check booking status."})
		return
	}
	if len(stateData) == 0 {
		log.Printf("[WebSocket] Booking session not found or expired for %s", bookingID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking session not found or expired."})
		return
	}

	status := stateData["status"]

	// --- NÂNG CẤP LÊN WEBSOCKET ---
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[WebSocket] Failed to upgrade connection for %s: %v", bookingID, err)
		return
	}

	// --- XỬ LÝ DỰA TRÊN TRẠNG THÁI ---
	if status == "COMPLETED" || status == "FAILED" {
		log.Printf("[WebSocket] Booking %s already processed (Status: %s). Sending immediate result.", bookingID, status)
		var messageToSend Message
		if status == "COMPLETED" {
			var payload models.TicketReturn
			json.Unmarshal([]byte(stateData["result_payload"]), &payload)
			messageToSend = Message{Type: "result", Payload: payload}
		} else {
			payload := map[string]string{"error": stateData["error_message"]}
			messageToSend = Message{Type: "error", Payload: payload}
		}

		msgBytes, _ := json.Marshal(messageToSend)
		conn.WriteMessage(websocket.TextMessage, msgBytes)
		conn.Close()
		return
	}

	// --- Trạng thái là QUEUED hoặc PROCESSING -> Tiến hành kết nối và lắng nghe ---
	ctx, cancel := context.WithCancel(c.Request.Context())
	client := &Client{
		conn: conn, bookingID: bookingID, manager: m,
		send: make(chan []byte, 256), ack: make(chan bool, 1),
		ctx: ctx, cancel: cancel,
	}
	m.register(client)

	client.readPump() // Chặn handler để giữ kết nối
}
