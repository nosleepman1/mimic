package domain

import "errors"

var (
	ErrEndpointNotFound    = errors.New("endpoint introuvable dans le contrat d'API")
	ErrMethodNotAllowed    = errors.New("methode HTTP non autorisee pour cette route")
	ErrInvalidContractFile = errors.New("fichier de contrat d'API invalide ou introuvable")
	ErrUnsupportedFormat   = errors.New("format de fichier non supporte (OpenAPI, Swagger, Postman ou YAML attendu)")
	ErrResourceNotFound    = errors.New("ressource introuvable en memoire")
	ErrValidationFailed    = errors.New("echec de validation du contrat d'API")
)
