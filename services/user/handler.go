package user

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wafi04/backend/pkg/middleware"
	httpresponse "github.com/wafi04/backend/pkg/response"
	"github.com/wafi04/backend/pkg/types"
	request "github.com/wafi04/backend/pkg/types/req"
)

type UserHandler struct {
	userrepo UserRepository
}

func NewUserHandler(userepo UserRepository) *UserHandler {
	return &UserHandler{
		userrepo: userepo,
	}
}

func (h *UserHandler) HandleCreateUserDetails(c *gin.Context) {
	user, err := middleware.GetUserFromGinContext(c)
	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	var req struct {
		PlaceBirth  string            `json:"place_birth"`
		DateBirth   *time.Time        `json:"date_birth,omitempty"`
		Gender      string            `json:"gender"`
		PhoneNumber string            `json:"phone_number"`
		Bio         string            `json:"bio"`
		Preferences types.Preferences `json:"preferences"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "error : %s", err.Error())
		return
	}

	create, err := h.userrepo.CreateUserDetails(c, &request.ReqCreateUserDetails{
		UserID:      user.UserID,
		PlaceBirth:  req.PlaceBirth,
		DateBirth:   req.DateBirth,
		Gender:      req.Gender,
		PhoneNumber: req.PhoneNumber,
		Bio:         req.Bio,
		Preferences: req.Preferences,
	})
	if err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed to create user details ", err.Error())
		return
	}

	httpresponse.SendSuccessResponse(c, http.StatusCreated, "Created User Details Successfully", create)
}

func (h *UserHandler) GetUserDetails(c *gin.Context) {
	user, err := middleware.GetUserFromGinContext(c)

	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userdet, err := h.userrepo.GetUserDetails(c, user.UserID)
	if err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed to get Profiles Details ", err.Error())
		return
	}

	httpresponse.SendSuccessResponse(c, http.StatusOK, "Get User Details successfully", userdet)

}

func (h *UserHandler) HandleUpdateProfiles(c *gin.Context) {
	user, err := middleware.GetUserFromGinContext(c)

	if err != nil {
		httpresponse.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	var req struct {
		PlaceBirth  *string            `json:"place_birth,omitempty"`
		DateBirth   *time.Time         `json:"date_birth,omitempty"`
		Gender      *string            `json:"gender,omitempty"`
		PhoneNumber *string            `json:"phone_number,omitempty"`
		Bio         *string            `json:"bio,omitempty"`
		Preferences *types.Preferences `json:"preferences,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "error : %s", err.Error())
		return
	}

	update, err := h.userrepo.UpdateUserDetails(c, &request.ReqUpdateUserDetails{
		UserID:      user.UserID,
		PlaceBirth:  req.PlaceBirth,
		DateBirth:   req.DateBirth,
		Gender:      req.Gender,
		PhoneNumber: req.PhoneNumber,
		Bio:         req.Bio,
		Preferences: req.Preferences,
	})

	if err != nil {
		httpresponse.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed to update user details", err.Error())
		return
	}

	httpresponse.SendSuccessResponse(c, http.StatusOK, "Update User details successfuly", update)

}
