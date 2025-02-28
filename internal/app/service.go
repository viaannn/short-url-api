package app

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"short_url/config"
	"strings"
)

// Service represents the service layer
type Service struct {
	repository IRepository
}

// InitService initializes the service component
func InitService(repository IRepository) *Service {
	return &Service{repository}
}

// CreateMasterUrl creates a short link based on the given request
func (service *Service) CreateMasterUrl(request CreateShortLinkRequest) (*CreateShortLinkResponse, error) {
	// Hashing the target URL to create a unique key as a short link
	hash := sha256.Sum256([]byte(request.TargetUrl))      // Generate SHA-256 hash
	encoded := base64.URLEncoding.EncodeToString(hash[:]) // Base64 encode
	key := strings.TrimRight(encoded[:8], "=")            // Trim, to use the first 8 characters
	fmt.Println("Generated short link key: ", key)

	// CreateMasterUrl method from the repository layer to Insert new record in the database
	createdData, dbError := service.repository.CreateMasterUrl(key, request.TargetUrl)
	// Handle and return DB Error, if any
	if dbError != nil {
		return nil, dbError
	}

	// Get domain host from the environment variable, if not found use default value
	domainHost := config.GetEnv(config.EnvDomainName, "http://localhost:8080/")

	// Return the response
	return &CreateShortLinkResponse{
		ShortLink: domainHost + createdData.Key,
		TargetUrl: createdData.TargetUrl,
	}, nil
}

// GetTargetUrl returns the target URL based on the short link key
func (service *Service) GetTargetUrl(key string) (string, error) {
	// FindByKey method from the repository layer to get the target URL based on the short link key
	entity, dbError := service.repository.FindByKey(key)
	// Handle and return DB Error, if any
	if dbError != nil {
		return "", dbError
	}

	return entity.TargetUrl, nil
}
