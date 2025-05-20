package dto

type RegisterRequest struct {
    Name      string `json:"name" binding:"required"`
    Email     string `json:"email" binding:"required,email"`
    Password  string `json:"password" binding:"required,min=8"`
    BirthYear int    `json:"birth_year" binding:"required"`
    Phone     string `json:"phone" binding:"required"`
}

type VerifyOTPRequest struct {
    Email    string `json:"email" binding:"required,email"`
    OTP      string `json:"otp" binding:"required,len=6"`
    Name     string `json:"name" binding:"required"`
    Password string `json:"password" binding:"required,min=8"`
    Phone    string `json:"phone" binding:"required"`
    BirthYear int    `json:"birth_year" binding:"required"`
}

type RegisterResponse struct {
    Message string `json:"message"`
    Email   string `json:"email"`
}

type VerifyOTPResponse struct {
    Message string `json:"message"`
    Token   string `json:"token,omitempty"`
}
