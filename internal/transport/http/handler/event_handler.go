package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/louco-event/internal/dto"
	"github.com/louco-event/internal/i18n"
	"github.com/louco-event/internal/middleware"
	"github.com/louco-event/internal/service"
)

type EventHandler struct {
	eventService      service.EventService
	addressService    service.AddressService
	ticketService     service.TicketService
	invitationService service.InvitationService
	creatorService    service.CreatorService
	i18n              *i18n.I18n
}

func NewEventHandler(
	eventService service.EventService,
	addressService service.AddressService,
	ticketService service.TicketService,
	invitationService service.InvitationService,
	creatorService service.CreatorService,
	i18n *i18n.I18n,
) *EventHandler {
	return &EventHandler{
		eventService:      eventService,
		addressService:    addressService,
		ticketService:     ticketService,
		invitationService: invitationService,
		creatorService:    creatorService,
		i18n:              i18n,
	}
}

// CreateEvent creates a new event
func (h *EventHandler) CreateEvent(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	var req dto.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	event, err := h.eventService.CreateEvent(c.Request.Context(), int(userID), req)
	if err != nil {
		var message string
		switch err.Error() {
		case "user not found":
			message = middleware.Translate(c, "user.not_found")
		case "creator profile not found":
			message = middleware.Translate(c, "creator.not_found")
		case "event.image_not_found":
			message = middleware.Translate(c, "event.image_not_found")
		case "event.video_not_found":
			message = middleware.Translate(c, "event.video_not_found")
		case "event.address_not_found":
			message = middleware.Translate(c, "event.address_not_found")
		case "event.category_not_found":
			message = middleware.Translate(c, "event.category_not_found")
		default:
			message = middleware.Translate(c, "event.create.failed")
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "event.create.success"),
		event,
	)
	c.JSON(http.StatusCreated, response)
}

// GetEvent retrieves an event by ID
func (h *EventHandler) GetEvent(c *gin.Context) {
	eventIDStr := c.Param("id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 32)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			"Invalid event ID",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var userID *int
	if uid, exists := middleware.GetCurrentUserID(c); exists {
		uidInt := int(uid)
		userID = &uidInt
	}

	event, err := h.eventService.GetEventByID(c.Request.Context(), int(eventID), userID)
	if err != nil {
		var message string
		switch err.Error() {
		case "event not found":
			message = middleware.Translate(c, "event.not_found")
		case "access denied":
			message = middleware.Translate(c, "event.access_denied")
		default:
			message = middleware.Translate(c, "common.internal_server_error")
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(http.StatusNotFound, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "event.get.success"),
		event,
	)
	c.JSON(http.StatusOK, response)
}

// UpdateEvent updates an existing event
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	eventIDStr := c.Param("id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 32)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			"Invalid event ID",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var req dto.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Get creator by user ID first
	// Note: UpdateEvent service method should be updated to accept userID instead of creatorID
	event, err := h.eventService.UpdateEvent(c.Request.Context(), int(eventID), int(userID), req)
	if err != nil {
		var message string
		switch err.Error() {
		case "event not found":
			message = middleware.Translate(c, "event.not_found")
		case "access denied":
			message = middleware.Translate(c, "event.access_denied")
		default:
			message = middleware.Translate(c, "event.update.failed")
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "event.update.success"),
		event,
	)
	c.JSON(http.StatusOK, response)
}

// DeleteEvent deletes an event
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	eventIDStr := c.Param("id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 32)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			"Invalid event ID",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Note: DeleteEvent service method should be updated to accept userID instead of creatorID
	err = h.eventService.DeleteEvent(c.Request.Context(), int(eventID), int(userID))
	if err != nil {
		var message string
		switch err.Error() {
		case "event not found":
			message = middleware.Translate(c, "event.not_found")
		case "access denied":
			message = middleware.Translate(c, "event.access_denied")
		default:
			message = middleware.Translate(c, "event.delete.failed")
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "event.delete.success"),
		nil,
	)
	c.JSON(http.StatusOK, response)
}

// GetMyEvents retrieves events created by the authenticated user
func (h *EventHandler) GetMyEvents(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Note: GetCreatorEvents service method should be updated to accept userID instead of creatorID
	events, paginationResp, err := h.eventService.GetCreatorEvents(c.Request.Context(), int(userID), pagination)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "event.list.failed"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "event.list.success"),
		dto.ListResponse{
			Items:      events,
			Pagination: *paginationResp,
		},
	)
	c.JSON(http.StatusOK, response)
}

// GetEvents retrieves events with filtering
func (h *EventHandler) GetEvents(c *gin.Context) {
	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Use GetPublicEvents method for now
	events, paginationResp, err := h.eventService.GetPublicEvents(c.Request.Context(), pagination)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "event.list.failed"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "event.list.success"),
		dto.ListResponse{
			Items:      events,
			Pagination: *paginationResp,
		},
	)
	c.JSON(http.StatusOK, response)
}

// UpdateEventStatus updates event status
func (h *EventHandler) UpdateEventStatus(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	eventIDStr := c.Param("id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 32)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			"Invalid event ID",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var req dto.UpdateEventStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Note: UpdateEventStatus service method should be updated to accept userID instead of creatorID
	event, err := h.eventService.UpdateEventStatus(c.Request.Context(), int(eventID), int(userID), req)
	if err != nil {
		var message string
		switch err.Error() {
		case "event not found":
			message = middleware.Translate(c, "event.not_found")
		case "access denied":
			message = middleware.Translate(c, "event.access_denied")
		case "invalid status transition":
			message = middleware.Translate(c, "event.invalid_status_transition")
		default:
			message = middleware.Translate(c, "event.status.update.failed")
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "event.status.update.success"),
		event,
	)
	c.JSON(http.StatusOK, response)
}

// GetEventStatistics retrieves event statistics
func (h *EventHandler) GetEventStatistics(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// Note: GetEventStats service method should be updated to accept userID instead of creatorID
	stats, err := h.eventService.GetEventStats(c.Request.Context(), int(userID))
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "event.statistics.failed"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "event.statistics.success"),
		stats,
	)
	c.JSON(http.StatusOK, response)
}

// SearchEventsByLocation searches events by location
func (h *EventHandler) SearchEventsByLocation(c *gin.Context) {
	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	city := c.Query("city")
	country := c.Query("country")

	events, paginationResp, err := h.eventService.GetPublicEventsByLocation(c.Request.Context(), city, country, pagination)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "event.location.search.failed"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "event.location.search.success"),
		dto.ListResponse{
			Items:      events,
			Pagination: *paginationResp,
		},
	)
	c.JSON(http.StatusOK, response)
}

// CreateTicket creates a new ticket for an event
func (h *EventHandler) CreateTicket(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	eventIDStr := c.Param("event_id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 32)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			"Invalid event ID format: "+eventIDStr,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var req dto.CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Get creator ID from user ID using creator service
	creator, err := h.creatorService.GetCreatorByUserID(c.Request.Context(), int(userID))
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "creator.not_found"),
			nil,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	ticket, err := h.ticketService.CreateTicket(c.Request.Context(), int(eventID), creator.ID, req)
	if err != nil {
		var message string
		switch err.Error() {
		case "event not found":
			message = middleware.Translate(c, "event.not_found")
		case "access denied":
			message = middleware.Translate(c, "event.access_denied")
		default:
			message = middleware.Translate(c, "ticket.create.failed")
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "ticket.create.success"),
		ticket,
	)
	c.JSON(http.StatusCreated, response)
}

// GetEventTickets retrieves tickets for an event
func (h *EventHandler) GetEventTickets(c *gin.Context) {
	eventIDStr := c.Param("event_id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 32)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			"Invalid event ID",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	tickets, err := h.ticketService.GetTicketsByEventID(c.Request.Context(), int(eventID))
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "ticket.list.failed"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "ticket.list.success"),
		tickets,
	)
	c.JSON(http.StatusOK, response)
}

// CreateInvitation creates a new invitation for an event
func (h *EventHandler) CreateInvitation(c *gin.Context) {
	_, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	eventIDStr := c.Param("event_id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 32)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			"Invalid event ID",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var req dto.CreateInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	invitation, err := h.invitationService.CreateInvitation(c.Request.Context(), int(eventID), req)
	if err != nil {
		var message string
		switch err.Error() {
		case "event not found":
			message = middleware.Translate(c, "event.not_found")
		case "access denied":
			message = middleware.Translate(c, "event.access_denied")
		case "invitation already exists":
			message = middleware.Translate(c, "invitation.already_exists")
		default:
			message = middleware.Translate(c, "invitation.create.failed")
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "invitation.create.success"),
		invitation,
	)
	c.JSON(http.StatusCreated, response)
}

// GetEventInvitations retrieves invitations for an event
func (h *EventHandler) GetEventInvitations(c *gin.Context) {
	_, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	eventIDStr := c.Param("event_id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 32)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			"Invalid event ID",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	invitations, paginationResp, err := h.invitationService.GetInvitationsByEventID(c.Request.Context(), int(eventID), pagination)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "invitation.list.failed"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "invitation.list.success"),
		dto.ListResponse{
			Items:      invitations,
			Pagination: *paginationResp,
		},
	)
	c.JSON(http.StatusOK, response)
}

// RespondToInvitation responds to an event invitation
func (h *EventHandler) RespondToInvitation(c *gin.Context) {
	_, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	invitationIDStr := c.Param("invitation_id")
	invitationID, err := strconv.ParseInt(invitationIDStr, 10, 32)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			"Invalid invitation ID",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var req dto.UpdateInvitationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	invitation, err := h.invitationService.UpdateInvitationStatus(c.Request.Context(), int(invitationID), req)
	if err != nil {
		var message string
		switch err.Error() {
		case "invitation not found":
			message = middleware.Translate(c, "invitation.not_found")
		case "access denied":
			message = middleware.Translate(c, "invitation.access_denied")
		case "invitation expired":
			message = middleware.Translate(c, "invitation.expired")
		default:
			message = middleware.Translate(c, "invitation.respond.failed")
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "invitation.respond.success"),
		invitation,
	)
	c.JSON(http.StatusOK, response)
}

// GetMyDraftEvents retrieves draft events created by the authenticated user
func (h *EventHandler) GetMyDraftEvents(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Note: GetCreatorDraftEvents service method should be updated to accept userID instead of creatorID
	events, paginationResp, err := h.eventService.GetCreatorDraftEvents(c.Request.Context(), int(userID), pagination)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "event.list.failed"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "event.list.success"),
		dto.ListResponse{
			Items:      events,
			Pagination: *paginationResp,
		},
	)
	c.JSON(http.StatusOK, response)
}

// GetMyPublishedEvents retrieves published events created by the authenticated user
func (h *EventHandler) GetMyPublishedEvents(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Note: GetCreatorPublishedEvents service method should be updated to accept userID instead of creatorID
	events, paginationResp, err := h.eventService.GetCreatorPublishedEvents(c.Request.Context(), int(userID), pagination)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "event.list.failed"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "event.list.success"),
		dto.ListResponse{
			Items:      events,
			Pagination: *paginationResp,
		},
	)
	c.JSON(http.StatusOK, response)
}

// SubmitForReview submits an event for review
func (h *EventHandler) SubmitForReview(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	eventIDStr := c.Param("id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 32)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			"Invalid event ID",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Note: SubmitEventForReview service method should be updated to accept userID instead of creatorID
	event, err := h.eventService.SubmitEventForReview(c.Request.Context(), int(eventID), int(userID))
	if err != nil {
		var message string
		switch err.Error() {
		case "event not found":
			message = middleware.Translate(c, "event.not_found")
		case "access denied":
			message = middleware.Translate(c, "event.access_denied")
		case "invalid status transition":
			message = middleware.Translate(c, "event.invalid_status_transition")
		default:
			message = middleware.Translate(c, "event.submit.failed")
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "event.submit.success"),
		event,
	)
	c.JSON(http.StatusOK, response)
}

// PublishEvent publishes an event
func (h *EventHandler) PublishEvent(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	eventIDStr := c.Param("id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 32)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			"Invalid event ID",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Note: PublishEvent service method should be updated to accept userID instead of creatorID
	event, err := h.eventService.PublishEvent(c.Request.Context(), int(eventID), int(userID))
	if err != nil {
		var message string
		switch err.Error() {
		case "event not found":
			message = middleware.Translate(c, "event.not_found")
		case "access denied":
			message = middleware.Translate(c, "event.access_denied")
		case "invalid status transition":
			message = middleware.Translate(c, "event.invalid_status_transition")
		default:
			message = middleware.Translate(c, "event.publish.failed")
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "event.publish.success"),
		event,
	)
	c.JSON(http.StatusOK, response)
}

// CancelEvent cancels an event
func (h *EventHandler) CancelEvent(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	eventIDStr := c.Param("id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 32)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			"Invalid event ID",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Note: CancelEvent service method should be updated to accept userID instead of creatorID
	event, err := h.eventService.CancelEvent(c.Request.Context(), int(eventID), int(userID))
	if err != nil {
		var message string
		switch err.Error() {
		case "event not found":
			message = middleware.Translate(c, "event.not_found")
		case "access denied":
			message = middleware.Translate(c, "event.access_denied")
		case "invalid status transition":
			message = middleware.Translate(c, "event.invalid_status_transition")
		default:
			message = middleware.Translate(c, "event.cancel.failed")
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "event.cancel.success"),
		event,
	)
	c.JSON(http.StatusOK, response)
}

// GetEventStats retrieves event statistics for the authenticated creator
func (h *EventHandler) GetEventStats(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// Note: GetEventStats service method should be updated to accept userID instead of creatorID
	stats, err := h.eventService.GetEventStats(c.Request.Context(), int(userID))
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "event.statistics.failed"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "event.statistics.success"),
		stats,
	)
	c.JSON(http.StatusOK, response)
}

// GetPublicEvents retrieves public events
func (h *EventHandler) GetPublicEvents(c *gin.Context) {
	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	events, paginationResp, err := h.eventService.GetPublicEvents(c.Request.Context(), pagination)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "event.list.failed"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "event.list.success"),
		dto.ListResponse{
			Items:      events,
			Pagination: *paginationResp,
		},
	)
	c.JSON(http.StatusOK, response)
}

// SearchEvents searches public events
func (h *EventHandler) SearchEvents(c *gin.Context) {
	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var searchReq dto.EventSearchRequest
	if err := c.ShouldBindQuery(&searchReq); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	events, paginationResp, err := h.eventService.SearchPublicEvents(c.Request.Context(), searchReq, pagination)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "event.search.failed"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "event.search.success"),
		dto.ListResponse{
			Items:      events,
			Pagination: *paginationResp,
		},
	)
	c.JSON(http.StatusOK, response)
}

// GetEventsByLocation retrieves events by location
func (h *EventHandler) GetEventsByLocation(c *gin.Context) {
	city := c.Param("city")
	country := c.Query("country")

	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	events, paginationResp, err := h.eventService.GetPublicEventsByLocation(c.Request.Context(), city, country, pagination)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "event.location.search.failed"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "event.location.search.success"),
		dto.ListResponse{
			Items:      events,
			Pagination: *paginationResp,
		},
	)
	c.JSON(http.StatusOK, response)
}

// GetUpcomingEvents retrieves upcoming events
func (h *EventHandler) GetUpcomingEvents(c *gin.Context) {
	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	events, paginationResp, err := h.eventService.GetUpcomingEvents(c.Request.Context(), pagination)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "event.upcoming.failed"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "event.upcoming.success"),
		dto.ListResponse{
			Items:      events,
			Pagination: *paginationResp,
		},
	)
	c.JSON(http.StatusOK, response)
}

// GetEventsByCategory retrieves events by category
func (h *EventHandler) GetEventsByCategory(c *gin.Context) {
	categoryIDStr := c.Param("category_id")
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 32)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			"Invalid category ID",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	events, paginationResp, err := h.eventService.GetPublicEventsByCategory(c.Request.Context(), int(categoryID), pagination)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "event.category.search.failed"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "event.category.search.success"),
		dto.ListResponse{
			Items:      events,
			Pagination: *paginationResp,
		},
	)
	c.JSON(http.StatusOK, response)
}
