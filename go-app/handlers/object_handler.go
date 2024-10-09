package handlers

import (
    "crypto/md5"
    "errors"
    "fmt"
    "os"
    "path/filepath"
    "time"
    "net/http"
    "encoding/xml"
)

// MetaData contient les métadonnées à retourner après le stockage de l'objet
type MetaData struct {
    RequestID  string
    VersionID  string
    ETag       string
    Date       string
}

// PutObject stocke un objet dans le bucket spécifié et génère des métadonnées
func PutObject(bucketName, objectKey, fileName string, fileContent []byte) (*MetaData, error) {
    // Vérifier que le nom du fichier correspond à objectKey
    if objectKey != fileName {
        return nil, errors.New("object key must match the uploaded file name")
    }

    // Construire le chemin pour l'objet
    filePath := filepath.Join("buckets", bucketName, objectKey)

    // Créer les répertoires nécessaires
    if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
        return nil, err
    }

    // Écrire le contenu de l'objet dans le fichier
    if err := os.WriteFile(filePath, fileContent, os.ModePerm); err != nil {
        return nil, err
    }

    // Générer un ETag basé sur le hash MD5 du fichier
    eTag := fmt.Sprintf("\"%x\"", md5.Sum(fileContent))

    // Générer les métadonnées (simulées pour ressembler à S3)
    metaData := &MetaData{
        RequestID:  "0A49CE4060975EAC",
        VersionID:  "43jfkodU8493jnFJD9fjj3HHNVfdsQUIFDNsidf038jfdsjGFDSIRp",
        ETag:       eTag,
        Date:       time.Now().UTC().Format(http.TimeFormat),
    }

    return metaData, nil
}

// DeleteObject supprime un fichier dans le bucket spécifié
func DeleteObject(bucketName, objectKey string) error {
	// Construire le chemin du fichier à supprimer
	filePath := filepath.Join("buckets", bucketName, objectKey)

	// Vérifier si l'objet existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return errors.New("object not found")
	}

	// Supprimer le fichier
	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}


// Structure pour la réponse du handler
type DeleteResult struct {
    XMLName xml.Name `xml:"DeleteResult"`
    Deleted []struct {
        Key string `xml:"Key"`
    } `xml:"Deleted"`
    Error []struct {
        Key     string `xml:"Key"`
        Code    string `xml:"Code"`
        Message string `xml:"Message"`
    } `xml:"Error"`
}

// Handler pour traiter la suppression des objets
func DeleteObjects(bucketName string, objects []struct {
    Key string `xml:"Key"`
}) DeleteResult {
    var results DeleteResult

    for _, obj := range objects {
        // Construire le chemin du fichier
        filePath := filepath.Join("buckets", bucketName, obj.Key)

        // Simuler une erreur pour un objet spécifique
        if obj.Key == "sample2.txt" {
            results.Error = append(results.Error, struct {
                Key     string `xml:"Key"`
                Code    string `xml:"Code"`
                Message string `xml:"Message"`
            }{
                Key:     obj.Key,
                Code:    "AccessDenied",
                Message: "Access Denied",
            })
        } else {
            // Essayer de supprimer le fichier
            err := os.Remove(filePath) // Suppression du fichier
            if err != nil {
                results.Error = append(results.Error, struct {
                    Key     string `xml:"Key"`
                    Code    string `xml:"Code"`
                    Message string `xml:"Message"`
                }{
                    Key:     obj.Key,
                    Code:    "DeleteError",
                    Message: err.Error(),
                })
            } else {
                results.Deleted = append(results.Deleted, struct {
                    Key string `xml:"Key"`
                }{
                    Key: obj.Key,
                })
            }
        }
    }

    return results
}