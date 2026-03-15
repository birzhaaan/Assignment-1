package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"practice5/internal/models"
	"practice5/internal/repository"
)

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) GetUsers(c *gin.Context) {
	var filter models.UserFilter

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "5"))
	if err != nil {
		pageSize = 5
	}

	filter.Page = page
	filter.PageSize = pageSize
	filter.OrderBy = c.DefaultQuery("order_by", "id")

	if idStr := c.Query("id"); idStr != "" {
		if id, err := strconv.Atoi(idStr); err == nil {
			filter.ID = &id
		}
	}

	if name := c.Query("name"); name != "" {
		filter.Name = &name
	}

	if email := c.Query("email"); email != "" {
		filter.Email = &email
	}

	if gender := c.Query("gender"); gender != "" {
		filter.Gender = &gender
	}

	if birthDate := c.Query("birth_date"); birthDate != "" {
		filter.BirthDate = &birthDate
	}

	result, err := h.repo.GetPaginatedUsers(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetCommonFriends(c *gin.Context) {
	user1, err1 := strconv.Atoi(c.Query("user1"))
	user2, err2 := strconv.Atoi(c.Query("user2"))

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user ids",
		})
		return
	}

	friends, err := h.repo.GetCommonFriends(user1, user2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, friends)
}
