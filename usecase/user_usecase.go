package usecase

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"login-jwt-otp/model"
	"login-jwt-otp/repository"
	"login-jwt-otp/utils/service"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	CreateUserUsecase(user model.Users) (model.Users, error)
	GenerateOTP(email string) (string, error)
	VerifyOTP(email string, otp string) (model.Users, error)
	GetUserByEmail(email string) (model.Users, error)
	GetAllUsersUsecase() ([]model.Users, error)
}

type userUsecase struct {
	UserRepo repository.UserRepoInterface
}

func NewUserUsecase(repo repository.UserRepoInterface) *userUsecase {
	return &userUsecase{
		UserRepo: repo,
	}
}

func (u *userUsecase) GenerateOTP(email string) (string, error) {
	// Validate email
	if !service.IsValidEmail(email) {
		return "", errors.New("invalid email format")
	}

	// Check if email exists
	exists, err := u.UserRepo.IsEmailExists(email)
	if err != nil {
		return "", err
	}
	if exists {
		return "", errors.New("email already registered")
	}

	// Generate 6 digit OTP

	optCode := fmt.Sprintf("%06d", rand.Intn(999999))

	// Create OTP record
	opt := model.OTP{
		Email:     email,
		Code:      optCode,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		Attempts:  0,
	}

	// Save OTP to database
	if err := u.UserRepo.SaveOTP(&opt); err != nil {
		return "", err
	}

	// Send OTP via email (implement email service)
	// For now, we'll just log it
	log.Printf("OTP sent to %s: %s", email, optCode)

	return optCode, nil
}

func (u *userUsecase) VerifyOTP(email string, otp string) (model.Users, error) {
	// Get OTP from database
	storedOTP, err := u.UserRepo.GetOTPByCode(otp)
	if err != nil {
		return model.Users{}, err
	}

	// Check if OTP expired
	if time.Now().After(storedOTP.ExpiresAt) {
		return model.Users{}, errors.New("OTP has expired")
	}

	// Check if too many attempts
	if storedOTP.Attempts >= 3 {
		return model.Users{}, errors.New("too many attempts")
	}

	// Create user
	user := model.Users{
		Email:        email,
		PasswordHash: storedOTP.Code, // Changed to use PasswordHash
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return model.Users{}, err
	}
	user.PasswordHash = string(hashedPassword)

	// Delete used OTP
	if err := u.UserRepo.DeleteOTP(otp); err != nil {
		log.Printf("Failed to delete OTP: %v", err)
	}

	return user, nil
}

func (u *userUsecase) CreateUserUsecase(user model.Users) (model.Users, error) {
	// Validate user input
	if !service.IsValidEmail(user.Email) {
		return model.Users{}, errors.New("invalid email format")
	}

	if err := service.IsValidPassword(user.PasswordHash); err != nil {
		return model.Users{}, err
	}

	exists, err := u.UserRepo.IsEmailExists(user.Email)
	if err != nil {
		return model.Users{}, err
	}
	if exists {
		return model.Users{}, errors.New("email already in use")
	}

	return u.UserRepo.CreateUser(user)
}

func (u *userUsecase) GetUserByIDUsecase(id int) (model.Users, error) {
	return u.UserRepo.GetUserByID(id)
}

// GetAllUsersUsecase retrieves all users from the repository
func (u *userUsecase) GetAllUsersUsecase() ([]model.Users, error) {
	users, err := u.UserRepo.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %v", err)
	}
	return users, nil
}

// GetUserByEmail retrieves a user by email from the repository
func (u *userUsecase) GetUserByEmail(email string) (model.Users, error) {
	if !service.IsValidEmail(email) {
		return model.Users{}, errors.New("invalid email format")
	}

	user, err := u.UserRepo.GetUserByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Users{}, fmt.Errorf("user with email %s not found", email)
		}
		return model.Users{}, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}
