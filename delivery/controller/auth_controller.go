package controller

import (
	"login-jwt-otp/model"
	"login-jwt-otp/model/dto"
	"login-jwt-otp/usecase"
	"login-jwt-otp/utils/service"
	"net/http"

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
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	// Check if email exists
	exists, err := a.AuthUsecase.IsEmailExists(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to process request",
			Error:   "internal server error",
		})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, dto.ErrorResponse{
			Message: "Registration failed",
			Error:   "email already registered",
		})
		return
	}

	// Generate and send OTP
	_, err = a.UserUsecase.GenerateOTP(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to process request",
			Error:   "failed to generate OTP",
		})
		return
	}

	// Return success response without sensitive data
	c.JSON(http.StatusOK, dto.RegisterResponse{
		Message: "OTP sent successfully. Please check your email.",
		Email:   req.Email,
	})
}

// VerifyOtp handles OTP verification and completes registration
func (a *AuthController) VerifyOtp(c *gin.Context) {
	var req dto.VerifyOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	// Verify OTP
	_, err := a.UserUsecase.VerifyOTP(req.Email, req.OTP)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Verification failed",
			Error:   err.Error(),
		})
		return
	}

	// Create user with verified email
	newUser := model.Users{
		Name:      req.Name,
		Email:     req.Email,
		BirthYear: req.BirthYear,
		Phone:     req.Phone,
		Role:      "USER",
	}

	// Register the user
	authResp, err := a.AuthUsecase.Register(newUser, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Registration failed",
			Error:   "failed to create user account: " + err.Error(),
		})
		return
	}

	// Return success response with token
	c.JSON(http.StatusOK, dto.Response{
		Message: "Registration successful",
		Data:    authResp,
	})
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

	c.JSON(http.StatusOK, dto.Response{
		Message: "Login successful",
		Data:    authResp,
	})
}
