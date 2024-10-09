package main

import (
	"github.com/gofiber/fiber/v2"

	router "bucket-s3-app/router"
)

func main() {
	app := fiber.New()

	app.Static("/", "./buckets") 

	router.InitRouter(app)

	app.Listen(":3000")
}