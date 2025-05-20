package main

import (
	"login-jwt-otp/delivery"
)

func main() {
	server := delivery.NewServer()
	server.Run()
}

//ini adalah contoh perubahanf
