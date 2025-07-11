package routes

import (
	"ticket-service/api/controllers"
	"ticket-service/pkg/websocket"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine,
	managerTicketController *controllers.ManagerTicketController,
	ticketController *controllers.TicketController,
	tokenTestController *controllers.TokenTestController,
	checkinController *controllers.CheckinController,
	wsManager *websocket.Manager,
) {
	ticketGroup := r.Group("/api/v1")
	{
		ticketGroup.POST("/initiate-booking", ticketController.InitiateBookingHandler)

		// Handler cho route WebSocket được định nghĩa ngay tại đây
		// để có thể truyền wsManager đã được khởi tạo từ main vào
		ticketGroup.GET("/ws/track/:bookingId", wsManager.HandleConnection)

		//User get our ticket
		ticketGroup.GET("/tickets/:id", ticketController.GetTicketHandler)
		ticketGroup.GET("/tickets", ticketController.GetAllTicketHandler)
		ticketGroup.POST("/tickets", ticketController.CreateTicketHandler)

		// New Public Route to get a ticket by ID
		ticketGroup.GET("/public/ticket/:id", ticketController.GetPublicTicketInfoByIDHandler)

		//User get our ticket
		ticketGroup.GET("/tickets/all", ticketController.GetAllTicketsPaginatedHandler) // New route for all tickets (paginated)

		ticketGroup.POST("/ticket-by-phone", ticketController.GetInfoTicketHandler)

		// Staff create ticket (cash payment)
		ticketGroup.POST("/staff/tickets", ticketController.CreateTicketByStaffHandler) // Requires staff authentication

		//Get ticket have order
		ticketGroup.GET("/tickets-available/:id", ticketController.GetAvailableHandler)
		ticketGroup.POST("/trips-available-seats", ticketController.GetAvailableMultiTripsHandler)

		//Payment
		ticketGroup.POST("/payments", managerTicketController.UpdateManagerTicketHandler)

		// Nhóm route mới cho check-in
		checkinGroup := ticketGroup.Group("/checkin")
		{
			// POST /api/v1/checkin
			// Route này xử lý việc check-in một vé
			checkinGroup.POST("/", checkinController.CheckinHandler)

			// GET /api/v1/checkin/trip/{tripID}
			// Route này lấy tất cả các lượt check-in của một chuyến đi
			checkinGroup.GET("/trip/:tripID", checkinController.GetTripCheckinsHandler)
		}

		// New token test route
		ticketGroup.GET("/token-test", tokenTestController.ValidateTokenHandler)
	}
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}
