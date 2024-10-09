package struc

import (
	"encoding/xml"
)

// ListObjectsResponse structure pour la réponse XML
type ListObjectsResponse struct {
	XMLName         xml.Name             `xml:"ListBucketResult"`
	Name            string               `xml:"Name"`
	Prefix          string               `xml:"Prefix"`
	Marker          string               `xml:"Marker"`
	MaxKeys         int                  `xml:"MaxKeys"`
	IsTruncated     bool                 `xml:"IsTruncated"`
	Contents        []Object             `xml:"Contents"`
	CommonPrefixes  []CommonPrefix       `xml:"CommonPrefixes"`
}

// Object structure pour chaque objet
type Object struct {
	Key          string `xml:"Key"`
	LastModified string `xml:"LastModified"`
	ETag         string `xml:"ETag"`
	Size         int64  `xml:"Size"`
	StorageClass string `xml:"StorageClass"`
	Owner        Owner  `xml:"Owner"`
}

// Owner structure pour le propriétaire de l'objet
type Owner struct {
	ID          string `xml:"ID"`
	DisplayName string `xml:"DisplayName"`
}

// CommonPrefix structure pour les préfixes communs
type CommonPrefix struct {
	Prefix string `xml:"Prefix"`
}

type DeleteObjectRequest struct {
    XMLName xml.Name        `xml:"Delete"`
    Objects []struct {
        Key string `xml:"Key"`
    } `xml:"Object"`
    Quiet bool `xml:"Quiet"`
}