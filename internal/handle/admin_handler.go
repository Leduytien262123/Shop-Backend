package handle

import (
	"backend/internal/consts"
	"backend/internal/helpers"
	"backend/internal/model"
	"backend/internal/repo"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminHandler struct {
	userRepo *repo.UserRepository
}

func NewAdminHandler(userRepo *repo.UserRepository) *AdminHandler {
	return &AdminHandler{userRepo: userRepo}
}

func (h *AdminHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userRepo.GetAllUsers()
	if err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, consts.MSG_INTERNAL_ERROR, err)
		return
	}

	var response []model.UserResponse
	for _, user := range users {
		response = append(response, model.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			Role:     user.Role,
			IsActive: user.IsActive,
		})
	}

	helpers.SuccessResponse(c, consts.MSG_SUCCESS, response)
}

func (h *AdminHandler) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		helpers.ValidationErrorResponse(c, "Invalid user ID")
		return
	}

	user, err := h.userRepo.GetUserByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ErrorResponse(c, http.StatusNotFound, consts.MSG_USER_NOT_FOUND, nil)
			return
		}
		helpers.ErrorResponse(c, http.StatusInternalServerError, consts.MSG_INTERNAL_ERROR, err)
		return
	}

	response := model.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.Role,
		IsActive: user.IsActive,
	}

	helpers.SuccessResponse(c, consts.MSG_SUCCESS, response)
}

func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		helpers.ValidationErrorResponse(c, "Invalid user ID")
		return
	}

	var input struct {
		Role string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.ValidationErrorResponse(c, consts.MSG_VALIDATION_ERROR)
		return
	}

	// Validate role
	if input.Role != consts.ROLE_ADMIN && input.Role != consts.ROLE_USER {
		helpers.ValidationErrorResponse(c, "Invalid role")
		return
	}

	user, err := h.userRepo.GetUserByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ErrorResponse(c, http.StatusNotFound, consts.MSG_USER_NOT_FOUND, nil)
			return
		}
		helpers.ErrorResponse(c, http.StatusInternalServerError, consts.MSG_INTERNAL_ERROR, err)
		return
	}

	user.Role = input.Role
	if err := h.userRepo.UpdateUser(user); err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, consts.MSG_INTERNAL_ERROR, err)
		return
	}

	response := model.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.Role,
		IsActive: user.IsActive,
	}

	helpers.SuccessResponse(c, consts.MSG_SUCCESS, response)
}

func (h *AdminHandler) ToggleUserStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		helpers.ValidationErrorResponse(c, "Invalid user ID")
		return
	}

	user, err := h.userRepo.GetUserByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ErrorResponse(c, http.StatusNotFound, consts.MSG_USER_NOT_FOUND, nil)
			return
		}
		helpers.ErrorResponse(c, http.StatusInternalServerError, consts.MSG_INTERNAL_ERROR, err)
		return
	}

	user.IsActive = !user.IsActive
	if err := h.userRepo.UpdateUser(user); err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, consts.MSG_INTERNAL_ERROR, err)
		return
	}

	response := model.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.Role,
		IsActive: user.IsActive,
	}

	helpers.SuccessResponse(c, consts.MSG_SUCCESS, response)
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		helpers.ValidationErrorResponse(c, "Invalid user ID")
		return
	}

	// Check if user exists
	_, err = h.userRepo.GetUserByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ErrorResponse(c, http.StatusNotFound, consts.MSG_USER_NOT_FOUND, nil)
			return
		}
		helpers.ErrorResponse(c, http.StatusInternalServerError, consts.MSG_INTERNAL_ERROR, err)
		return
	}

	if err := h.userRepo.DeleteUser(uint(id)); err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, consts.MSG_INTERNAL_ERROR, err)
		return
	}

	helpers.SuccessResponse(c, "User deleted successfully", nil)
}
