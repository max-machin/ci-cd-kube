package middlewares

import (
	"regexp" // Regular expression search : https://pkg.go.dev/regexp
)

// Taille min (3) - max (63) du nom du Bucket.
const (
	minBucketNameLength = 3
	maxBucketNameLength = 63
)

// Validator struct encapsule les règles de validation pour les noms de bucket.
type Validator struct {
	prefixPattern   *regexp.Regexp
	suffixPattern   *regexp.Regexp
	ipAddressPattern *regexp.Regexp
	dotsPattern     *regexp.Regexp
	namePattern     *regexp.Regexp
}


// NewBucketNameValidator crée une nouvelle instance de Validator avec des règles prédéfinies.
func NewBucketNameValidator() *Validator {
	return &Validator{
		// Initialisation des patterns de regex pour la validation des noms de bucket.
		prefixPattern:   regexp.MustCompile(`^(xn--|sthree-|sthree-configurator-|amzn-s3-demo-)`),
		suffixPattern:   regexp.MustCompile(`(-s3alias|--ol-s3|\.mrap|--x-s3)$`),
		ipAddressPattern: regexp.MustCompile(`^\d{1,3}(\.\d{1,3}){3}$`),
		dotsPattern:     regexp.MustCompile(`\.\.`),
		namePattern:     regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9.-]*[a-zA-Z0-9])?$`),
	}
}


func (v *Validator) Validate(name string) []string {
    var errors []string

    // Vérification de la taille du nom
    if len(name) < minBucketNameLength || len(name) > maxBucketNameLength {
        errors = append(errors, "Nom du bucket doit être entre 3 et 63 caractères.")
    }

    // Vérification des préfixes invalides
    if v.prefixPattern.MatchString(name) {
        errors = append(errors, "Nom du bucket ne peut pas commencer par un préfixe invalide.")
    }

    // Vérification des suffixes invalides
    if v.suffixPattern.MatchString(name) {
        errors = append(errors, "Nom du bucket ne peut pas se terminer par un suffixe invalide.")
    }
    
    // Vérification pour s'assurer que le nom n'est pas une adresse IP
    if v.ipAddressPattern.MatchString(name) {
        errors = append(errors, "Nom du bucket ne peut pas être une adresse IP.")
    }

    // Vérification pour s'assurer qu'il n'y a pas de points consécutifs
    if v.dotsPattern.MatchString(name) {
        errors = append(errors, "Nom du bucket ne peut pas contenir des points consécutifs.")
    }

    // Vérification du pattern général du nom de bucket
    if !v.namePattern.MatchString(name) {
        errors = append(errors, "Nom du bucket doit commencer et se terminer par une lettre ou un chiffre.")
    }

    if len(errors) == 0 {
        // Le nom est valide, retour d'un tableau vide pour les messages
        return []string{}
    }

    return errors
}

