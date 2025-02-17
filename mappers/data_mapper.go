package mappers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"os"
	"rpl-service/constants"
	"rpl-service/models"
	"strings"
)

func GetCourseRequest(r *http.Request) (models.Course, error) {
	var body models.Course
	if decodeErr := json.NewDecoder(r.Body).Decode(&body); decodeErr != nil {
		return models.Course{}, errors.New("invalid request body")
	}
	return body, nil
}

func GetEnrolledRequest(r *http.Request) (EnrollmentRequest, error) {
	var result EnrollmentRequest
	if json.NewDecoder(r.Body).Decode(&result) != nil {
		return result, errors.New("invalid request body")
	}
	return result, nil
}

func GetDeleteRequest(r *http.Request) (DeleteRequest, error) {
	var result DeleteRequest
	if json.NewDecoder(r.Body).Decode(&result) != nil {
		return result, errors.New("invalid request body")
	}
	return result, nil
}

// GetUserID TODO: Check this when possible.
func GetUserID(r *http.Request) (string, error) {
	// Auth0 configuration
	issuerURL := os.Getenv("AUTH0_DOMAIN")
	audience := []string{os.Getenv("AUTH0_AUDIENCE")}

	// Parse the issuer URL
	issuer, err := url.Parse(issuerURL)
	if err != nil {
		return constants.EmptyString, fmt.Errorf("failed to parse issuer URL: %w", err)
	}

	// Set up token validator
	provider := jwks.NewCachingProvider(issuer, constants.ProviderDuration)
	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL,
		audience,
	)

	if err != nil {
		return constants.EmptyString, fmt.Errorf("failed to setup validator: %w", err)
	}

	// Validate token
	token, err := jwtValidator.ValidateToken(r.Context(), r.Header.Get("Authorization"))
	if err != nil {
		return constants.EmptyString, fmt.Errorf("invalid token: %w", err)
	}

	// Extract claims
	claims, ok := token.(*validator.ValidatedClaims)
	if !ok {
		return constants.EmptyString, errors.New("invalid token claims")
	}

	// Get user ID from subject
	sub := claims.RegisteredClaims.Subject
	parts := strings.Split(sub, "|")
	if len(parts) != constants.IDPartsAmount {
		return constants.EmptyString, errors.New("invalid subject format")
	}

	return parts[1], nil
}

// AUX

type EnrollmentRequest struct {
	UserID   uuid.UUID `json:"UserID"`
	CourseID uuid.UUID `json:"CourseID"`
}

type DeleteRequest struct {
	UserID    uuid.UUID `json:"UserID"`
	CourseID  uuid.UUID `json:"CourseID"`
	StudentID uuid.UUID `json:"StudentID"`
}
