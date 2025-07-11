package services

import (
	// Replace 'api_gateway' with your actual Go module name

	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// Service defines the properties of a backend service.
type Service struct {
	Name     string
	URL      string
	Prefix   string
	Priority int // Higher priority means it's checked first for a matching prefix
}

// ServiceRegistry holds a list of registered services.
type ServiceRegistry struct {
	Services []Service
}

// NewServiceRegistry creates a new instance of ServiceRegistry.
func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{Services: []Service{}}
}

// RegisterService adds a new service to the registry.
func (sr *ServiceRegistry) RegisterService(name, serviceURL, prefix string, priority int) {
	sr.Services = append(sr.Services, Service{Name: name, URL: serviceURL, Prefix: prefix, Priority: priority})
	log.Printf("Service Registered: Name=%s, URL=%s, Prefix=%s, Priority=%d", name, serviceURL, prefix, priority)
}

// GetServiceForPath finds the best matching service for a given request path
// based on prefix length and priority.
func (sr *ServiceRegistry) GetServiceForPath(path string) (*Service, bool) {
	var bestMatch *Service
	highestPriority := -1  // Initialize with a value lower than any possible priority
	longestPrefixLen := -1 // Initialize to track the longest prefix for a given highest priority

	for i := range sr.Services {
		service := &sr.Services[i]
		if strings.HasPrefix(path, service.Prefix) {
			currentPrefixLen := len(service.Prefix)

			// If current service has higher priority, it's the new best match
			if service.Priority > highestPriority {
				highestPriority = service.Priority
				longestPrefixLen = currentPrefixLen
				bestMatch = service
			} else if service.Priority == highestPriority {
				// If priorities are the same, prefer the one with a longer (more specific) prefix
				if currentPrefixLen > longestPrefixLen {
					longestPrefixLen = currentPrefixLen
					bestMatch = service
				}
			}
		}
	}
	if bestMatch != nil {
		log.Printf("Path '%s' matched service '%s' (Prefix: '%s', Priority: %d)", path, bestMatch.Name, bestMatch.Prefix, bestMatch.Priority)
	} else {
		log.Printf("No service found for path: %s", path)
	}
	return bestMatch, bestMatch != nil
}

// ProxyHandler forwards requests to the appropriate backend service.
func (sr *ServiceRegistry) ProxyHandler(c *gin.Context) {
	requestPath := c.Request.URL.Path
	service, found := sr.GetServiceForPath(requestPath)
	if !found {
		log.Printf("ProxyHandler: Service not found for path %s", requestPath)
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "Service not found for this path"})
		return
	}

	targetURL, err := url.Parse(service.URL)
	if err != nil {
		log.Printf("ProxyHandler: Error parsing target URL '%s' for service '%s': %v", service.URL, service.Name, err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "Could not parse target service URL"})
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Modify the request before sending it to the target service
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req) // Sets req.URL.Scheme, req.URL.Host

		// Ensure the path sent to the backend service is the full original path
		// that matched the prefix, not just the remainder.
		// Or, if services expect paths relative to their prefix, adjust accordingly.
		// Based on original code, it seems the full path is intended.
		req.URL.Path = requestPath

		// Set standard proxy headers
		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Real-IP", c.ClientIP())

		// Handle X-Forwarded-For
		if prior, ok := req.Header["X-Forwarded-For"]; ok {
			req.Header.Set("X-Forwarded-For", strings.Join(prior, ", ")+", "+c.ClientIP())
		} else {
			req.Header.Set("X-Forwarded-For", c.ClientIP())
		}

		// Handle X-Forwarded-Proto
		if proto := c.Request.Header.Get("X-Forwarded-Proto"); proto != "" {
			req.Header.Set("X-Forwarded-Proto", proto)
		} else {
			if c.Request.TLS != nil {
				req.Header.Set("X-Forwarded-Proto", "https")
			} else {
				req.Header.Set("X-Forwarded-Proto", "http")
			}
		}

		// Custom headers for downstream services
		req.Header.Set("X-Gateway-Target-Service", service.Name)
		req.Header.Set("X-Original-Gateway-Path", requestPath) // The original path received by gateway

		// Pass user information if available from middleware
		if userIDVal, exists := c.Get("userID"); exists {
			if userID, ok := userIDVal.(int); ok {
				req.Header.Set("X-User-ID", fmt.Sprintf("%d", userID))
			}
		}
		if userRoleVal, exists := c.Get("userRole"); exists {
			if userRole, ok := userRoleVal.(string); ok {
				req.Header.Set("X-User-Role", userRole)
			}
		}
		// Clear default User-Agent if you don't want to pass Gin's default
		// req.Header.Del("User-Agent")
	}

	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		log.Printf("Proxy error to service %s (%s) for path %s: %v", service.Name, service.URL, req.URL.Path, err)
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(rw).Encode(gin.H{
			"code":    http.StatusBadGateway,
			"message": "Upstream service error",
			"details": err.Error(),
		})
	}

	log.Printf("Proxying request for '%s' to service '%s' at '%s'. UserID: %v, Role: %v",
		requestPath, service.Name, service.URL, c.Value("userID"), c.Value("userRole"))
	proxy.ServeHTTP(c.Writer, c.Request)
}
