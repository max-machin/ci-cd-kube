package controller

import (
    "encoding/xml"
    "fmt"
    "github.com/gofiber/fiber/v2"
    "log"
	"strings"
	"strconv"
	"io"
	"time"
	"os"
	"net/http"
    "bucket-s3-app/handlers" // Assurez-vous que ce chemin est correct pour vos gestionnaires
    "bucket-s3-app/struc"    // Importez le package struc
	"bucket-s3-app/middlewares"
)

// Fonction pour gérer les requêtes HTTP et retourner la liste des buckets
func ListBuckets(c *fiber.Ctx) error {
	// Appeler le handler pour obtenir la liste des buckets
	result, err := handlers.ListBuckets()
	if err != nil {
		log.Printf("Error: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to read buckets directory")
	}

	// Convertir la réponse en XML
	xmlResponse, err := xml.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Printf("Error generating XML response: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to generate XML response")
	}

	// Ajouter la déclaration XML
	xmlResponse = []byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n" + string(xmlResponse))

	// Définir les en-têtes de réponse
	c.Set("Content-Type", "application/xml")
	c.Set("Content-Length", fmt.Sprintf("%d", len(xmlResponse)))
	c.Set("Server", "AmazonS3")

	// Retourner la réponse avec un statut 200 OK
	return c.SendString(string(xmlResponse))
}


// Handler for creating a bucket
func CreateBucket(c *fiber.Ctx) error {
	// Vérifiez si la méthode est PUT
	if c.Method() != fiber.MethodPut {
		return c.Status(fiber.StatusMethodNotAllowed).SendString("Method not allowed")
	}

	// Extrayez le nom du bucket depuis les paramètres de l'URL
	bucketName := c.Params("bucketName")
	if bucketName == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Bucket name is missing in URL")
	}

	// Optionnellement, vous pouvez lire le corps de la requête
	xmlData := c.Body()
	log.Println("Received XML:")
	log.Println(string(xmlData))

	// Parse the XML data
	var config struc.CreateBucketConfiguration
	err := xml.Unmarshal(xmlData, &config)
	if err != nil {
		log.Println("Error parsing XML:", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid XML")
	}

	validator := middlewares.NewBucketNameValidator()
	result := validator.Validate(bucketName)

	
	if len(result) > 0 {
		errorMessage := strings.Join(result, "; ")
		return c.Status(fiber.StatusBadRequest).SendString(errorMessage)
	}

	// Créez le bucket en appelant la fonction du handler
	err = handlers.CreateBucket(bucketName)
	if err != nil {
		log.Println("Error creating bucket:", err)
		if err.Error() == "Bucket already exists" {
			return c.Status(fiber.StatusConflict).SendString("Bucket already exists")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to create bucket folder")
	}

	// Ajouter les en-têtes spécifiques à Amazon S3
	c.Set("Location", fmt.Sprintf("/%s", bucketName))
	c.Set("x-amz-request-id", "236A8905248E5A01") // ID fictif pour l'exemple
	c.Set("x-amz-id-2", "YgIPIfBiKa2bj0KMg95r/0zo3emzU4dzsD4rcKCHQUAdQkf3ShJTOOpXUueF6QKo") // ID fictif pour l'exemple
	c.Set("Date", "Wed, 01 Mar 2006 12:00:00 GMT") // Date fictive pour l'exemple
	c.Set("Content-Length", "0")
	c.Set("Connection", "close")
	c.Set("Server", "AmazonS3")

	return c.Status(fiber.StatusOK).SendString(fmt.Sprintf("Location : /%s", bucketName))
}


// DeleteBucket supprime un bucket en utilisant le nom fourni dans l'URL
func DeleteBucket(c *fiber.Ctx) error {
	// Vérifiez si la méthode est DELETE
	if c.Method() != fiber.MethodDelete {
		return c.Status(fiber.StatusMethodNotAllowed).SendString("Method not allowed")
	}

	// Extrayez le nom du bucket depuis les paramètres de l'URL
	bucketName := c.Params("bucketName")
	if bucketName == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Bucket name is missing in URL")
	}

	// Supprimez le bucket en appelant la fonction du handler
	err := handlers.DeleteBucket(bucketName)
	if err != nil {
		log.Println("Error deleting bucket:", err)
		
			return c.Status(fiber.StatusNotFound).SendString("Bucket not found")

	}

	c.Set("x-amz-request-id", "236A8905248E5A01") // ID fictif pour l'exemple
	c.Set("x-amz-id-2", "YgIPIfBiKa2bj0KMg95r/0zo3emzU4dzsD4rcKCHQUAdQkf3ShJTOOpXUueF6QKo") // ID fictif pour l'exemple
	c.Set("Date", "Wed, 01 Mar 2006 12:00:00 GMT") // Date fictive pour l'exemple
	c.Set("Connection", "close")
	c.Set("Server", "AmazonS3")

	// Répondre avec un code 204 No Content
	return c.SendStatus(fiber.StatusNoContent)
}


// ListObjectsHandler récupère la liste des objets dans un bucket
func ListObjects(c *fiber.Ctx) error {
	bucketName := c.Params("bucketName")
	if bucketName == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Bucket name is missing in URL")
	}

	// Paramètres de requête
	prefix := c.Query("prefix")
	delimiter := c.Query("delimiter")
	maxKeys := 1000
	if c.Query("max-keys") != "" {
		var err error
		maxKeys, err = strconv.Atoi(c.Query("max-keys"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid max-keys parameter")
		}
	}

	objects, err := handlers.ListObjects(bucketName, prefix, delimiter, maxKeys)
	if err != nil {
		log.Println("Error listing objects:", err)
		if err.Error() == "bucket does not exist" {
			return c.Status(fiber.StatusNotFound).SendString("Bucket not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to list objects")
	}

	c.Set("x-amz-request-id", "236A8905248E5A01") // ID fictif pour l'exemple
	c.Set("x-amz-id-2", "YgIPIfBiKa2bj0KMg95r/0zo3emzU4dzsD4rcKCHQUAdQkf3ShJTOOpXUueF6QKo") // ID fictif pour l'exemple
	c.Set("Date", "Wed, 01 Mar 2006 12:00:00 GMT") // Date fictive pour l'exemple
	c.Set("Connection", "close")
	c.Set("Server", "AmazonS3")

	// Répondre avec la réponse XML
	c.Type("xml")
	return c.SendString(serializeXML(objects))
}

func PutObject(c *fiber.Ctx) error {
    // Récupérer les paramètres URI
    bucketName := c.Params("bucketName")
    objectKey := c.Params("objectName")

    // Valider les paramètres
    if bucketName == "" || objectKey == "" {
        return c.Status(fiber.StatusBadRequest).SendString("Bucket name or object key is missing in URI")
    }

    // Lire le fichier depuis la requête
    fileHeader, err := c.FormFile("file")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).SendString("File upload failed")
    }

    // Lire le contenu du fichier
    fileContent, err := fileHeader.Open()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("Unable to open file")
    }
    defer fileContent.Close()

    // Convertir le fichier en byte array
    fileBytes, err := io.ReadAll(fileContent)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("Unable to read file content")
    }

    // Appeler le handler pour stocker l'objet et générer les métadonnées
    metaData, err := handlers.PutObject(bucketName, objectKey, fileHeader.Filename, fileBytes)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to store object: " + err.Error())
    }

    // Répondre avec les métadonnées générées
    c.Set("x-amz-id-2", metaData.RequestID)
    c.Set("x-amz-request-id", metaData.RequestID)
    c.Set("x-amz-version-id", metaData.VersionID)
    c.Set("ETag", metaData.ETag)
    c.Set("Date", metaData.Date)
    c.Set("Content-Length", "0")
    c.Set("Connection", "close")
    c.Set("Server", "AmazonS3")

    return c.SendStatus(fiber.StatusOK)
}

func DeleteObject(c *fiber.Ctx) error {

	bucketName := c.Params("bucketName")
	objectKey := c.Params("objectName")

	if bucketName == "" || objectKey == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Bucket name or object key is missing in URI")
	}
	
	err := handlers.DeleteObject(bucketName, objectKey)

	if err != nil {
		log.Println("Error deleting object:", err)
		if err.Error() == "object not found" {
			return c.Status(fiber.StatusNotFound).SendString("object not found")
		}

		return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete object")
	}


	 // En-têtes supplémentaires pour la suppression
	 c.Set("x-amz-delete-marker", "true")
	 c.Set("x-amz-request-charged", "requester")


	return c.SendStatus(fiber.StatusNoContent)
}



// Controller pour gérer la suppression des objets
func DeleteObjects(c *fiber.Ctx) error {
    log.Printf("Received DELETE request for batch deletion: %s", c.Request().URI().String())

    // Lire le corps de la requête
    bodyBytes := c.Body()
    log.Printf("Received body bytes: %v", bodyBytes)

    // Définir la structure de la requête
    var deleteReq struc.DeleteObjectRequest
    err := xml.Unmarshal(bodyBytes, &deleteReq)
    if err != nil {
        log.Printf("Error parsing XML: %v", err)
        return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
    }

    log.Printf("Received delete request: %+v", deleteReq)

    if len(deleteReq.Objects) == 0 {
        return c.Status(fiber.StatusBadRequest).SendString("No objects specified for deletion")
    }

    // Extraire le nom du bucket depuis les paramètres de l'URL ou de la requête
    bucketName := c.Params("bucketName") // ou utilisez c.Query("bucketName") si passé en tant que paramètre de requête

    // Appeler le handler pour traiter la suppression
    results := handlers.DeleteObjects(bucketName, deleteReq.Objects)

    // Convertir la réponse en XML
    xmlResponse, err := xml.MarshalIndent(results, "", "  ")
    if err != nil {
        log.Printf("Error generating XML response: %v", err)
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to generate XML response")
    }

    // Ajouter les en-têtes de réponse
    c.Set("x-amz-id-2", "5h4FxSNCUS7wP5z92eGCWDshNpMnRuXvETa4HH3LvvH6VAIr0jU7tH9kM7X+njXx")
    c.Set("x-amz-request-id", "A437B3B641629AEE")
    c.Set("Date", time.Now().UTC().Format(time.RFC1123))
    c.Set("Content-Type", "application/xml")
    c.Set("Content-Length", fmt.Sprintf("%d", len(xmlResponse)))
    c.Set("Server", "AmazonS3")

    // Retourner la réponse avec un statut 200 OK
    return c.SendString(string(xmlResponse))
}

// serializeXML convertit un ListObjectsResponse en XML
func serializeXML(data *struc.ListObjectsResponse) string {
	xmlData, err := xml.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println("Error marshalling XML:", err)
		return ""
	}
	return string(xmlData)
}

func GetObject(c *fiber.Ctx) error {
	fmt.Println("GetObject function called")

	// Extraire le nom du bucket et la clé de l'objet à partir de l'URL
	bucketName := c.Params("bucketName")
	key := c.Params("objectName")

	// Log les paramètres extraits
	fmt.Printf("Bucket Name: %s, Object Key: %s\n", bucketName, key)

	// Déterminer le chemin complet du fichier
	filePath := fmt.Sprintf("buckets/%s/%s", bucketName, key)
	fmt.Printf("File Path: %s\n", filePath)

	// Ouvrir le fichier
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File not found")
			return c.Status(fiber.StatusNotFound).SendString("File not found")
		}
		fmt.Printf("Error opening file: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error opening file")
	}
	defer file.Close()

	// Obtenir des informations sur le fichier
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("Error getting file info: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error getting file info")
	}

	// Définir les en-têtes de réponse
	c.Set("Last-Modified", fileInfo.ModTime().UTC().Format(http.TimeFormat))
	c.Set("x-amz-id-2", "sample-id")
	c.Set("x-amz-request-id", "sample-request-id")
	c.Set("Date", time.Now().UTC().Format(http.TimeFormat))
	c.Set("ETag", fmt.Sprintf("\"%s\"", fileInfo.ModTime().Unix()))
	c.Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// Log les en-têtes
	fmt.Printf("Response Headers: %+v\n", c.Response().Header)

	// Lire et retourner le contenu du fichier
	return c.SendFile(filePath)
}
