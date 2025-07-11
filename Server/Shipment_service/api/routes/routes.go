package routes

import (
	"shipment-service/api/controller"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, shipmentCtrl *controller.ShipmentController) {
	api := router.Group("/api/v1")

	// --- Trip specific routes ---
	trips := api.Group("/tripsshipments")
	{
		tripWithID := trips.Group("/:trip_id")
		// POST /api/v1/trips/:trip_id/shipments - Create shipment for a trip
		tripWithID.POST("/shipments", shipmentCtrl.CreateShipmentHandler)
		// GET /api/v1/trips/:trip_id/shipments - List all shipments for a trip
		tripWithID.GET("/shipments", shipmentCtrl.GetShipmentsByTripIDHandler)
		// GET /api/v1/trips/:trip_id/invoices - List all invoices for a trip
		tripWithID.GET("/invoices", shipmentCtrl.ListInvoicesByTripIDHandler)
	}

	// --- Shipment specific routes ---
	shipments := api.Group("/shipments")
	{
		// GET /api/v1/shipments - List all shipments with pagination
		shipments.GET("", shipmentCtrl.ListShipmentsHandler)

		shipmentWithID := shipments.Group("/:id")
		// GET /api/v1/shipments/:id - Get a single shipment by its ID
		shipmentWithID.GET("", shipmentCtrl.GetShipmentByIDHandler)
		// GET /api/v1/shipments/:id/invoice - Get the invoice for a specific shipment
		shipmentWithID.GET("/invoice", shipmentCtrl.GetInvoiceForShipmentHandler)
	}

	// --- Invoice specific routes ---
	invoices := api.Group("/shipinvoices")
	{
		// GET /api/v1/invoices - List all invoices with pagination
		invoices.GET("", shipmentCtrl.ListInvoicesHandler)
		// GET /api/v1/invoices/:id - Get a single invoice by its ID
		invoices.GET("/:id", shipmentCtrl.GetInvoiceByIDHandler)
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})
}
