package controller

import (
	"net/http"
	"login-jwt-otp/model/dto"
	"login-jwt-otp/usecase"
	"login-jwt-otp/utils/service"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	rg          *gin.RouterGroup
	jwtService  service.JwtService
	AuthUsecase usecase.AuthUsecase
	UserUsecase usecase.UserUsecase
}

func NewAuthController(rg *gin.RouterGroup, jwt service.JwtService, authUC usecase.AuthUsecase, userUC usecase.UserUsecase) *AuthController {
	return &AuthController{
		rg:          rg,
		jwtService:  jwt,
		AuthUsecase: authUC,
		UserUsecase: userUC,
	}
}

func (ac *AuthController) Route() {
	ac.rg.POST("/register", ac.Register)
	ac.rg.POST("/verify-otp", ac.VerifyOtp)
	ac.rg.POST("/login", ac.Login)
}

// Register handles user registration and sends OTP
func (a *AuthController) Register(c *gin.Context) {
	var input struct {
		Name      string `json:"name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		BirthYear int    `json:"birth_year"`
		Phone     string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Bad Request",
			Error:   "invalid input",
		})
		return
	}

	// Generate OTP
	if _, err := a.UserUsecase.GenerateOTP(input.Email); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to generate OTP",
			Error:   err.Error(),
		})
		return
	}

	// Return OTP information
	c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent successfully",
		"email":   input.Email,
	})
}

// VerifyOtp handles OTP verification and completes registration
func (a *AuthController) VerifyOtp(c *gin.Context) {
	var input struct {
		Email string `json:"email"`
		Otp   string `json:"otp"`
		Name  string `json:"name"`
		Password  string `json:"password"`
		BirthYear int    `json:"birth_year"`
		Phone     string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Bad Request",
			Error:   "invalid input",
		})
		return
	}

	// Verify OTP and create user
	user, err := a.UserUsecase.VerifyOTP(input.Email, input.Otp)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Invalid OTP",
			Error:   err.Error(),
		})
		return
	}

	// Update user details
	user.Name = input.Name
	user.Phone = input.Phone
	user.BirthYear = input.BirthYear
	user.Role = "user"

	// Update password
	hashedPassword, err := service.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to encrypt password",
			Error:   "gagal mengenkripsi password",
		})
		return
	}
	user.PasswordHash = hashedPassword

	// Save updated user
	updatedUser, err := a.AuthUsecase.Register(user, input.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Registration failed",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, updatedUser)
}

func (a *AuthController) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	authResp, err := a.AuthUsecase.Login(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Invalid credentials",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, authResp)
}
