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

type AddressHandler struct {
	addressService service.AddressService
	i18n           *i18n.I18n
}

func NewAddressHandler(
	addressService service.AddressService,
	i18n *i18n.I18n,
) *AddressHandler {
	return &AddressHandler{
		addressService: addressService,
		i18n:           i18n,
	}
}

// CreateAddress creates a new address
func (h *AddressHandler) CreateAddress(c *gin.Context) {
	var req dto.CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	address, err := h.addressService.CreateAddress(c.Request.Context(), req)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "address.create.failed"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "address.create.success"),
		address,
	)
	c.JSON(http.StatusCreated, response)
}

// GetAddress retrieves an address by ID
func (h *AddressHandler) GetAddress(c *gin.Context) {
	addressIDStr := c.Param("id")
	addressID, err := strconv.ParseInt(addressIDStr, 10, 32)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			"Invalid address ID",
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	address, err := h.addressService.GetAddressByID(c.Request.Context(), int(addressID))
	if err != nil {
		var message string
		switch err.Error() {
		case "address not found":
			message = middleware.Translate(c, "address.not_found")
		default:
			message = middleware.Translate(c, "common.internal_server_error")
		}

		response := dto.NewErrorResponse(message, nil)
		c.JSON(http.StatusNotFound, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "address.get.success"),
		address,
	)
	c.JSON(http.StatusOK, response)
}

// SearchAddresses searches addresses by location
func (h *AddressHandler) SearchAddresses(c *gin.Context) {
	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var filters dto.AddressFilterRequest
	if err := c.ShouldBindQuery(&filters); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	addresses, paginationResp, err := h.addressService.GetAddressesWithFilters(c.Request.Context(), filters, pagination)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "address.search.failed"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "address.search.success"),
		dto.ListResponse{
			Items:      addresses,
			Pagination: *paginationResp,
		},
	)
	c.JSON(http.StatusOK, response)
}

// GetAddressesByCity retrieves addresses by city
func (h *AddressHandler) GetAddressesByCity(c *gin.Context) {
	city := c.Param("city")

	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	addresses, paginationResp, err := h.addressService.GetAddressesByCity(c.Request.Context(), city, pagination)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "address.city.search.failed"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "address.city.search.success"),
		dto.ListResponse{
			Items:      addresses,
			Pagination: *paginationResp,
		},
	)
	c.JSON(http.StatusOK, response)
}
