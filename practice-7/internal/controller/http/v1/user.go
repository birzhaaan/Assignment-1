package v1

import (
	"net/http"
	"practice-7/internal/entity"
	"practice-7/internal/usecase"
	"practice-7/utils"
	"github.com/gin-gonic/gin"
)

type userRoutes struct {
	t usecase.UserInterface
}

func NewUserRoutes(handler *gin.RouterGroup, t usecase.UserInterface) {
	r := &userRoutes{t}
	h := handler.Group("/users")
	{
		h.POST("/register", r.Register)
		h.POST("/login", r.Login)
		
		protected := h.Group("/")
		protected.Use(utils.JWTAuthMiddleware()) 
		{
			protected.GET("/me", r.GetMe)
			
			protected.PATCH("/promote/:id", utils.RoleMiddleware("admin"), r.Promote)
		}
	}
}

func (r *userRoutes) Register(c *gin.Context) {
	var dto entity.CreateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashed, _ := utils.HashPassword(dto.Password)
	user := &entity.User{
		Username: dto.Username,
		Email:    dto.Email,
		Password: hashed,
		Role:     "user", 
	}
	created, token, err := r.t.RegisterUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": created, "token": token})
}

func (r *userRoutes) Login(c *gin.Context) {
	var dto entity.LoginUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := r.t.LoginUser(&dto)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (r *userRoutes) GetMe(c *gin.Context) {
	userID := c.GetString("userID")
	user, err := r.t.GetMe(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (r *userRoutes) Promote(c *gin.Context) {
	id := c.Param("id")
	if err := r.t.PromoteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Promotion failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User promoted to admin successfully"})
}