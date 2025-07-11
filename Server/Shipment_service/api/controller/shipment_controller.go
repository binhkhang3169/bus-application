package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"shipment-service/internal/db"
	"shipment-service/internal/model"
	"shipment-service/internal/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ShipmentController struct {
	shipmentService service.ShipmentService
}

func NewShipmentController(shipmentService service.ShipmentService) *ShipmentController {
	return &ShipmentController{shipmentService: shipmentService}
}

// CreateShipmentHandler handles POST /trips/:trip_id/shipments
func (sc *ShipmentController) CreateShipmentHandler(c *gin.Context) {
	tripIDStr := c.Param("trip_id")
	tripID, err := strconv.ParseInt(tripIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trip_id format"})
		return
	}

	var req model.CreateShipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	if req.TripID != int32(tripID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trip ID in body (%d) does not match URL (%d)", req.TripID, tripID)})
		return
	}

	createdShipment, err := sc.shipmentService.CreateShipment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create shipment: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdShipment)
}

// GetShipmentByIDHandler handles GET /shipments/:id
func (sc *ShipmentController) GetShipmentByIDHandler(c *gin.Context) {
	shipmentIDStr := c.Param("id")
	shipmentID, err := strconv.ParseInt(shipmentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shipment ID format"})
		return
	}

	shipment, err := sc.shipmentService.GetShipmentByID(c.Request.Context(), int32(shipmentID))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve shipment"})
		return
	}

	c.JSON(http.StatusOK, shipment)
}

// ListShipmentsHandler handles GET /shipments
func (sc *ShipmentController) ListShipmentsHandler(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit format"})
		return
	}
	offset, err := strconv.ParseInt(offsetStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset format"})
		return
	}

	params := db.ListShipmentsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	shipments, err := sc.shipmentService.ListShipments(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve shipments"})
		return
	}
	if shipments == nil {
		shipments = []db.Shipment{}
	}
	c.JSON(http.StatusOK, shipments)
}

// GetShipmentsByTripIDHandler handles GET /trips/:trip_id/shipments
func (sc *ShipmentController) GetShipmentsByTripIDHandler(c *gin.Context) {
	tripIDStr := c.Param("trip_id")
	tripID, err := strconv.ParseInt(tripIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trip_id format"})
		return
	}

	shipments, err := sc.shipmentService.GetShipmentsByTripID(c.Request.Context(), int32(tripID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve shipments"})
		return
	}
	if shipments == nil {
		shipments = []db.Shipment{}
	}
	c.JSON(http.StatusOK, shipments)
}

// GetInvoiceByIDHandler handles GET /invoices/:id
func (sc *ShipmentController) GetInvoiceByIDHandler(c *gin.Context) {
	invoiceIDStr := c.Param("id")
	invoiceID, err := strconv.ParseInt(invoiceIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID format"})
		return
	}

	invoice, err := sc.shipmentService.GetInvoiceByID(c.Request.Context(), int32(invoiceID))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve invoice"})
		return
	}
	c.JSON(http.StatusOK, invoice)
}

// GetInvoiceForShipmentHandler handles GET /shipments/:id/invoice
func (sc *ShipmentController) GetInvoiceForShipmentHandler(c *gin.Context) {
	shipmentIDStr := c.Param("id")
	shipmentID, err := strconv.ParseInt(shipmentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shipment_id format"})
		return
	}

	invoice, err := sc.shipmentService.GetInvoiceForShipment(c.Request.Context(), int32(shipmentID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve invoice"})
		return
	}
	c.JSON(http.StatusOK, invoice)
}

// ListInvoicesHandler handles GET /invoices
func (sc *ShipmentController) ListInvoicesHandler(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit format"})
		return
	}
	offset, err := strconv.ParseInt(offsetStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset format"})
		return
	}

	params := db.ListInvoicesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	invoices, err := sc.shipmentService.ListInvoices(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve invoices"})
		return
	}
	if invoices == nil {
		invoices = []db.Invoice{}
	}
	c.JSON(http.StatusOK, invoices)
}

// ListInvoicesByTripIDHandler handles GET /trips/:trip_id/invoices
func (sc *ShipmentController) ListInvoicesByTripIDHandler(c *gin.Context) {
	tripIDStr := c.Param("trip_id")
	tripID, err := strconv.ParseInt(tripIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trip_id format"})
		return
	}

	invoices, err := sc.shipmentService.ListInvoicesByTripID(c.Request.Context(), int32(tripID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve invoices"})
		return
	}

	if invoices == nil {
		invoices = []db.Invoice{}
	}
	c.JSON(http.StatusOK, invoices)
}
