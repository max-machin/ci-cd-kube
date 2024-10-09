package handlers

import (
	"fmt"
	"os"
	"log"
	"errors"
	"path/filepath"
	"time"
	"bucket-s3-app/struc"
)

// Handler pour lister les dossiers/buckets
func ListBuckets() (struc.ListAllMyBucketsResult, error) {
	var result struc.ListAllMyBucketsResult

	// Définir le chemin du répertoire contenant les buckets
	bucketsDir := "buckets"

	// Lire le contenu du répertoire
	files, err := os.ReadDir(bucketsDir)
	if err != nil {
		return result, fmt.Errorf("error reading buckets directory: %w", err)
	}

	// Préparer la liste des buckets
	for _, file := range files {
		if file.IsDir() {
			// Ici, on utilise une date de création fictive, à remplacer par la vraie date si disponible
			creationDate := time.Now().Format(time.RFC3339)
			result.Buckets.Bucket = append(result.Buckets.Bucket, struc.BucketResponse{
				CreationDate: creationDate,
				Name:         file.Name(),
			})
		}
	}

	// Ajouter des informations fictives pour le propriétaire et le token de continuation
	result.Owner = struc.BucketOwner{
		DisplayName: "exampleDisplayName",
		ID:          "exampleID",
	}
	result.ContinuationToken = "exampleContinuationToken"

	return result, nil
}

// CreateBucket prend le nom du bucket et crée le dossier correspondant
func CreateBucket(bucketName string) error {
	basePath := "buckets"
	bucketPath := fmt.Sprintf("%s/%s", basePath, bucketName)

	// Vérifiez si le répertoire de base existe
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		// Créez le répertoire de base s'il n'existe pas
		err := os.MkdirAll(basePath, os.ModePerm)
		if err != nil {
			log.Println("Error creating base directory:", err)
			return err
		}
	}

	// Créez le dossier du bucket
	err := os.Mkdir(bucketPath, os.ModePerm)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("Bucket already exists")
		}
		log.Println("Error creating bucket folder:", err)
		return err
	}

	return nil
}



// DeleteBucket supprime le dossier correspondant au bucket
func DeleteBucket(bucketName string) error {
	bucketPath := filepath.Join("buckets", bucketName)

	// Vérifiez si le bucket existe
	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		return fmt.Errorf("Bucket does not exists")
	}

	// Supprimez le bucket (le dossier) et tout son contenu
	err := os.RemoveAll(bucketPath)
	if err != nil {
		log.Println("Error deleting bucket:", err)
		return err
	}

	return nil
}



// ListObjects récupère la liste des objets dans un bucket
func ListObjects(bucketName string, prefix string, delimiter string, maxKeys int) (*struc.ListObjectsResponse, error) {
	bucketPath := filepath.Join("buckets", bucketName)

	// Vérifiez si le bucket existe
	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		return nil, errors.New("bucket does not exist")
	}

	// Liste tous les fichiers dans le bucket
	files, err := os.ReadDir(bucketPath)
	if err != nil {
		return nil, err
	}

	// Préparez la réponse
	response := &struc.ListObjectsResponse{
		Name:        bucketName,
		Prefix:      prefix,
		Marker:      "",
		MaxKeys:     maxKeys,
		IsTruncated: false,
	}

	// Ajoutez les objets à la réponse
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Obtenez les informations du fichier
		filePath := filepath.Join(bucketPath, file.Name())
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return nil, err
		}

		object := struc.Object{
			Key:          file.Name(),
			LastModified: fileInfo.ModTime().Format(time.RFC3339),
			ETag:         `"etag-placeholder"`, // Remplacez ceci par une vraie valeur ETag si nécessaire
			Size:         fileInfo.Size(),
			StorageClass: "STANDARD",
			Owner: struc.Owner{
				ID:          "owner-id-placeholder", // Remplacez ceci par une vraie valeur Owner ID si nécessaire
				DisplayName: "owner-display-name-placeholder", // Remplacez ceci par une vraie valeur Display Name si nécessaire
			},
		}
		response.Contents = append(response.Contents, object)
	}

	return response, nil
}