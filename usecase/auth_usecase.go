package usecase

import (
	"database/sql"
	"errors"
	"log"
	"login-jwt-otp/model"
	"login-jwt-otp/repository"
	"login-jwt-otp/utils/service"
)

type AuthUsecase struct {
	userRepo   repository.UserRepoInterface
	jwtService service.JwtService
}

type AuthUsecaseInterface interface {
	Login(email string, password string) (model.AuthResponse, error)
	Register(user model.Users, password string) (model.AuthResponse, error)
	IsEmailExists(email string) (bool, error)
	GenerateToken(user model.Users) (string, error)
}

func (u *AuthUsecase) IsEmailExists(email string) (bool, error) {
	return u.userRepo.IsEmailExists(email)
}

func NewAuthUsecase(userRepo repository.UserRepoInterface, jwtService service.JwtService) *AuthUsecase {
	return &AuthUsecase{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (u *AuthUsecase) GenerateToken(user model.Users) (string, error) {
	return u.jwtService.CreateToken(user)
}

func (u *AuthUsecase) Login(email string, password string) (model.AuthResponse, error) {
	user, err := u.userRepo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.AuthResponse{}, errors.New("email not found")
		}
		return model.AuthResponse{}, errors.New("failed to retrieve user")
	}

	if !service.CheckPasswordHash(password, user.PasswordHash) {
		log.Println("Password mismatch for email:", email)
		return model.AuthResponse{}, errors.New("incorrect password")
	}

	token, err := u.GenerateToken(user)
	if err != nil {
		return model.AuthResponse{}, err
	}

	return model.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (u *AuthUsecase) Register(user model.Users, plainPassword string) (model.AuthResponse, error) {
	// Check if email is already registered
	exists, err := u.userRepo.IsEmailExists(user.Email)
	if err != nil {
		return model.AuthResponse{}, err
	}
	if exists {
		return model.AuthResponse{}, errors.New("email already registered")
	}

	// Hash password and create user
	hashedPassword, err := service.HashPassword(plainPassword)
	if err != nil {
		return model.AuthResponse{}, err
	}
	user.PasswordHash = hashedPassword
	user.Role = "USER"

	createdUser, err := u.userRepo.CreateUser(user)
	if err != nil {
		return model.AuthResponse{}, err
	}

	log.Printf("User successfully registered!")

	return model.AuthResponse{
		User: createdUser,
	}, nil
}
