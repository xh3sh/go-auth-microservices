package handler

import (
	"net/http"
	"time"

	"github.com/xh3sh/go-auth-microservices/internal/constants"
	"github.com/xh3sh/go-auth-microservices/internal/models"
	"github.com/xh3sh/go-auth-microservices/internal/repository"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	repo      repository.Repository
	eventRepo repository.EventRepository
}

func NewUserHandler(repo repository.Repository, eventRepo repository.EventRepository) *UserHandler {
	return &UserHandler{
		repo:      repo,
		eventRepo: eventRepo,
	}
}

func (h *UserHandler) GetEventRepo() repository.EventRepository {
	return h.eventRepo
}

// GetUsers РІРѕР·РІСЂР°С‰Р°РµС‚ СЃРїРёСЃРѕРє РІСЃРµС… РїРѕР»СЊР·РѕРІР°С‚РµР»РµР№
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.repo.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	if h.eventRepo != nil {
		userID := c.GetHeader("X-User-ID")

		go h.eventRepo.PublishUserActionEvent(models.UserActionEvent{
			UserID:    userID,
			Action:    constants.ActionUserListViewed,
			Resource:  constants.ResourceUsers,
			Timestamp: time.Now(),
			Status:    constants.StatusSuccess,
		})
	}

	c.JSON(http.StatusOK, users)
}

// GetUserByID РІРѕР·РІСЂР°С‰Р°РµС‚ РёРЅС„РѕСЂРјР°С†РёСЋ Рѕ РїРѕР»СЊР·РѕРІР°С‚РµР»Рµ РїРѕ РµРіРѕ ID
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := h.repo.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if h.eventRepo != nil {
		viewerID := c.GetHeader("X-User-ID")

		go h.eventRepo.PublishUserActionEvent(models.UserActionEvent{
			UserID:     viewerID,
			Action:     constants.ActionUserProfileViewed,
			Resource:   constants.ResourceUser,
			ResourceID: id,
			Timestamp:  time.Now(),
			Status:     constants.StatusSuccess,
			Metadata: map[string]interface{}{
				"target_user_id": id,
			},
		})
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser СѓРґР°Р»СЏРµС‚ РїРѕР»СЊР·РѕРІР°С‚РµР»СЏ РїРѕ РµРіРѕ ID
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.repo.DeleteUser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	if h.eventRepo != nil {
		userID := c.GetHeader("X-User-ID")

		go h.eventRepo.PublishUserActionEvent(models.UserActionEvent{
			UserID:     userID,
			Action:     constants.ActionUserDeleted,
			Resource:   constants.ResourceUser,
			ResourceID: id,
			Timestamp:  time.Now(),
			Status:     constants.StatusSuccess,
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
