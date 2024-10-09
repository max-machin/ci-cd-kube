package router

import (
	"github.com/gofiber/fiber/v2"

	Controller "bucket-s3-app/controllers"
)

func InitRouter(app *fiber.App) {
	
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	/****** Bucket routes ******/
	// Get list of all buckets
	app.Get("/buckets", Controller.ListBuckets)
	// Create new bucket
	app.Put("/buckets/:bucketName", Controller.CreateBucket)
	// Delete bucket
	app.Delete("/buckets/:bucketName", Controller.DeleteBucket)

	/****** Objects routes ******/
	// Get list of all object for one bucket
	app.Get("/buckets/:bucketName/objects", Controller.ListObjects)
	// Get one bucket object
	app.Get("buckets/:bucketName/objects/:objectName", Controller.GetObject)
	// Add one bucket object
	app.Put("/buckets/:bucketName/:objectName", Controller.PutObject)
	// Delete one bucket object
	app.Delete("/buckets/:bucketName/:objectName", Controller.DeleteObject)
	// Delete multiple bucket objects
	app.Delete("/buckets/:bucketName/objects/delete", Controller.DeleteObjects)


}