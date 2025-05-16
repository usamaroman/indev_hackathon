package main

import (
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/app"
)

// @title CRM Backend API
// @version 1.0
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	app.Run()
}
