package struc

import (
	"encoding/xml"
)

// Define the XML structure
type CreateBucketConfiguration struct {
	XMLName              xml.Name `xml:"CreateBucketConfiguration"`
	LocationConstraint   string   `xml:"LocationConstraint"`
	Location             Location `xml:"Location"`
	Bucket               Bucket   `xml:"Bucket"`
}

type Location struct {
	Name string `xml:"Name"`
	Type string `xml:"Type"`
}

type Bucket struct {
	DataRedundancy string `xml:"DataRedundancy"`
	Type           string `xml:"Type"`
	Name 		   string `xml:"Name"`
}

// Structure pour représenter un bucket
type BucketResponse struct {
	XMLName      xml.Name `xml:"Bucket"`
	CreationDate string   `xml:"CreationDate"`
	Name         string   `xml:"Name"`
}

// Structure pour représenter les buckets dans la réponse
type Buckets struct {
	XMLName xml.Name `xml:"Buckets"`
	Bucket  []BucketResponse `xml:"Bucket"`
}

// Structure pour représenter le propriétaire dans la réponse
type BucketOwner struct {
	XMLName    xml.Name `xml:"Owner"`
	DisplayName string   `xml:"DisplayName"`
	ID         string   `xml:"ID"`
}

// Structure pour représenter la réponse complète
type ListAllMyBucketsResult struct {
	XMLName             xml.Name `xml:"ListAllMyBucketsResult"`
	Buckets             Buckets  `xml:"Buckets"`
	Owner               BucketOwner    `xml:"Owner"`
	ContinuationToken   string   `xml:"ContinuationToken"`
}