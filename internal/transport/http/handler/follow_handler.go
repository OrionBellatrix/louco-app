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

// FollowHandler handles follow-related HTTP requests
type FollowHandler struct {
	followService *service.FollowService
	i18n          *i18n.I18n
}

// NewFollowHandler creates a new follow handler instance
func NewFollowHandler(followService *service.FollowService, i18n *i18n.I18n) *FollowHandler {
	return &FollowHandler{
		followService: followService,
		i18n:          i18n,
	}
}

// Helper function to safely get string value from pointer
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Follow handles POST /follows - follow a user
func (h *FollowHandler) Follow(c *gin.Context) {
	var req dto.FollowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Get current user from auth middleware
	currentUserID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// Follow the user
	if err := h.followService.Follow(c.Request.Context(), currentUserID, req.UserID); err != nil {
		// Check if it's a domain error (business logic error)
		if domainErr, ok := err.(interface{ Error() string }); ok {
			response := dto.NewErrorResponse(
				middleware.Translate(c, domainErr.Error()),
				nil,
			)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.internal_server_error"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "follow.success"),
		nil,
	)
	c.JSON(http.StatusOK, response)
}

// Unfollow handles DELETE /follows/:user_id - unfollow a user
func (h *FollowHandler) Unfollow(c *gin.Context) {
	userIDStr := c.Param("user_id")
	targetUserID, err := strconv.Atoi(userIDStr)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			nil,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Get current user from auth middleware
	currentUserID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// Unfollow the user
	if err := h.followService.Unfollow(c.Request.Context(), currentUserID, targetUserID); err != nil {
		// Check if it's a domain error (business logic error)
		if domainErr, ok := err.(interface{ Error() string }); ok {
			response := dto.NewErrorResponse(
				middleware.Translate(c, domainErr.Error()),
				nil,
			)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.internal_server_error"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "follow.unfollow_success"),
		nil,
	)
	c.JSON(http.StatusOK, response)
}

// GetFollowStatus handles GET /follows/status/:user_id - check follow status
func (h *FollowHandler) GetFollowStatus(c *gin.Context) {
	userIDStr := c.Param("user_id")
	targetUserID, err := strconv.Atoi(userIDStr)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			nil,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Get current user from auth middleware
	currentUserID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// Check if current user follows target user
	isFollowing, err := h.followService.IsFollowing(c.Request.Context(), currentUserID, targetUserID)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.internal_server_error"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Check if target user follows current user
	isFollower, err := h.followService.IsFollowing(c.Request.Context(), targetUserID, currentUserID)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.internal_server_error"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	responseData := dto.FollowStatusResponse{
		IsFollowing: isFollowing,
		IsFollower:  isFollower,
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "follow.status_retrieved"),
		responseData,
	)
	c.JSON(http.StatusOK, response)
}

// GetFollowers handles GET /follows/followers - get user's followers
func (h *FollowHandler) GetFollowers(c *gin.Context) {
	var pagination dto.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Get current user from auth middleware
	currentUserID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	page, limit := pagination.GetDefaultPagination()
	offset := pagination.GetOffset()

	// Get followers
	followers, total, err := h.followService.GetFollowers(c.Request.Context(), currentUserID, limit, offset)
	if err != nil {
		// Check if it's a domain error (business logic error)
		if domainErr, ok := err.(interface{ Error() string }); ok {
			response := dto.NewErrorResponse(
				middleware.Translate(c, domainErr.Error()),
				nil,
			)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.internal_server_error"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Convert to response format
	followResponses := make([]*dto.FollowResponse, len(followers))
	for i, follow := range followers {
		followResponses[i] = &dto.FollowResponse{
			ID:         follow.ID,
			FollowedAt: follow.CreatedAt,
			User: &dto.UserFollowInfo{
				ID:             follow.Follower.ID,
				FirstName:      follow.Follower.FullName,
				LastName:       "",
				Username:       getStringValue(follow.Follower.Username),
				FollowersCount: follow.Follower.FollowersCount,
				FollowingCount: follow.Follower.FollowingCount,
			},
		}

		// Add profile picture if exists
		if follow.Follower.ProfilePicture != nil {
			followResponses[i].User.ProfilePicture = &dto.MediaResponse{
				ID:           follow.Follower.ProfilePicture.ID,
				UserID:       follow.Follower.ProfilePicture.UserID,
				OriginalName: follow.Follower.ProfilePicture.OriginalName,
				FileName:     follow.Follower.ProfilePicture.FileName,
				FileURL:      follow.Follower.ProfilePicture.FileURL,
				MediaType:    string(follow.Follower.ProfilePicture.MediaType),
				MimeType:     follow.Follower.ProfilePicture.MimeType,
				FileSize:     follow.Follower.ProfilePicture.FileSize,
				Width:        follow.Follower.ProfilePicture.Width,
				Height:       follow.Follower.ProfilePicture.Height,
				Duration:     follow.Follower.ProfilePicture.Duration,
				IsConverted:  follow.Follower.ProfilePicture.IsConverted,
				CreatedAt:    follow.Follower.ProfilePicture.CreatedAt,
				UpdatedAt:    follow.Follower.ProfilePicture.UpdatedAt,
			}
		}
	}

	responseData := dto.FollowersResponse{
		Followers:  followResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: dto.CalculateTotalPages(total, limit),
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "follow.followers_retrieved"),
		responseData,
	)
	c.JSON(http.StatusOK, response)
}

// GetFollowing handles GET /follows/following - get user's following
func (h *FollowHandler) GetFollowing(c *gin.Context) {
	var pagination dto.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Get current user from auth middleware
	currentUserID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	page, limit := pagination.GetDefaultPagination()
	offset := pagination.GetOffset()

	// Get following
	following, total, err := h.followService.GetFollowing(c.Request.Context(), currentUserID, limit, offset)
	if err != nil {
		// Check if it's a domain error (business logic error)
		if domainErr, ok := err.(interface{ Error() string }); ok {
			response := dto.NewErrorResponse(
				middleware.Translate(c, domainErr.Error()),
				nil,
			)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.internal_server_error"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Convert to response format
	followResponses := make([]*dto.FollowResponse, len(following))
	for i, follow := range following {
		followResponses[i] = &dto.FollowResponse{
			ID:         follow.ID,
			FollowedAt: follow.CreatedAt,
			User: &dto.UserFollowInfo{
				ID:             follow.Following.ID,
				FirstName:      follow.Following.FullName,
				LastName:       "",
				Username:       getStringValue(follow.Following.Username),
				FollowersCount: follow.Following.FollowersCount,
				FollowingCount: follow.Following.FollowingCount,
			},
		}

		// Add profile picture if exists
		if follow.Following.ProfilePicture != nil {
			followResponses[i].User.ProfilePicture = &dto.MediaResponse{
				ID:           follow.Following.ProfilePicture.ID,
				UserID:       follow.Following.ProfilePicture.UserID,
				OriginalName: follow.Following.ProfilePicture.OriginalName,
				FileName:     follow.Following.ProfilePicture.FileName,
				FileURL:      follow.Following.ProfilePicture.FileURL,
				MediaType:    string(follow.Following.ProfilePicture.MediaType),
				MimeType:     follow.Following.ProfilePicture.MimeType,
				FileSize:     follow.Following.ProfilePicture.FileSize,
				Width:        follow.Following.ProfilePicture.Width,
				Height:       follow.Following.ProfilePicture.Height,
				Duration:     follow.Following.ProfilePicture.Duration,
				IsConverted:  follow.Following.ProfilePicture.IsConverted,
				CreatedAt:    follow.Following.ProfilePicture.CreatedAt,
				UpdatedAt:    follow.Following.ProfilePicture.UpdatedAt,
			}
		}
	}

	responseData := dto.FollowingResponse{
		Following:  followResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: dto.CalculateTotalPages(total, limit),
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "follow.following_retrieved"),
		responseData,
	)
	c.JSON(http.StatusOK, response)
}

// GetMutualFollows handles GET /follows/mutual/:user_id - get mutual follows
func (h *FollowHandler) GetMutualFollows(c *gin.Context) {
	userIDStr := c.Param("user_id")
	targetUserID, err := strconv.Atoi(userIDStr)
	if err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			nil,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var pagination dto.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.validation_failed"),
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Get current user from auth middleware
	currentUserID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.unauthorized"),
			nil,
		)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	page, limit := pagination.GetDefaultPagination()
	offset := pagination.GetOffset()

	// Get mutual follows
	mutualFollows, total, err := h.followService.GetMutualFollows(c.Request.Context(), currentUserID, targetUserID, limit, offset)
	if err != nil {
		// Check if it's a domain error (business logic error)
		if domainErr, ok := err.(interface{ Error() string }); ok {
			response := dto.NewErrorResponse(
				middleware.Translate(c, domainErr.Error()),
				nil,
			)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		response := dto.NewErrorResponse(
			middleware.Translate(c, "common.internal_server_error"),
			nil,
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Convert to response format
	followResponses := make([]*dto.FollowResponse, len(mutualFollows))
	for i, follow := range mutualFollows {
		followResponses[i] = &dto.FollowResponse{
			ID:         follow.ID,
			FollowedAt: follow.CreatedAt,
			User: &dto.UserFollowInfo{
				ID:             follow.Following.ID,
				FirstName:      follow.Following.FullName,
				LastName:       "",
				Username:       getStringValue(follow.Following.Username),
				FollowersCount: follow.Following.FollowersCount,
				FollowingCount: follow.Following.FollowingCount,
			},
		}

		// Add profile picture if exists
		if follow.Following.ProfilePicture != nil {
			followResponses[i].User.ProfilePicture = &dto.MediaResponse{
				ID:           follow.Following.ProfilePicture.ID,
				UserID:       follow.Following.ProfilePicture.UserID,
				OriginalName: follow.Following.ProfilePicture.OriginalName,
				FileName:     follow.Following.ProfilePicture.FileName,
				FileURL:      follow.Following.ProfilePicture.FileURL,
				MediaType:    string(follow.Following.ProfilePicture.MediaType),
				MimeType:     follow.Following.ProfilePicture.MimeType,
				FileSize:     follow.Following.ProfilePicture.FileSize,
				Width:        follow.Following.ProfilePicture.Width,
				Height:       follow.Following.ProfilePicture.Height,
				Duration:     follow.Following.ProfilePicture.Duration,
				IsConverted:  follow.Following.ProfilePicture.IsConverted,
				CreatedAt:    follow.Following.ProfilePicture.CreatedAt,
				UpdatedAt:    follow.Following.ProfilePicture.UpdatedAt,
			}
		}
	}

	responseData := dto.MutualFollowsResponse{
		MutualFollows: followResponses,
		Total:         total,
		Page:          page,
		Limit:         limit,
		TotalPages:    dto.CalculateTotalPages(total, limit),
	}

	response := dto.NewSuccessResponse(
		middleware.Translate(c, "follow.mutual_follows_retrieved"),
		responseData,
	)
	c.JSON(http.StatusOK, response)
}
