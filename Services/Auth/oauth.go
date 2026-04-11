package Auth

import (
	"quay-go-api/Entities/Models"
	"quay-go-api/Repositories"
	"quay-go-api/Services/Logger"
	"strings"
)

const AccessTokenPrefixLength = 20

/*
ValidateBearerToken Validate an OAuth token found inside the Authorization header and indicate whether it's a valid OAuth token
*/
func ValidateBearerToken(authHeader string) (bool, Models.OauthAccessToken) {
	// We assume that the token is not empty (validate by the first condition of the authorized middleware)

	normalized := strings.Split(authHeader, " ")
	if len(normalized) != 2 || strings.ToLower(normalized[0]) != "bearer" {
		return false, Models.OauthAccessToken{}
	}
	oauthToken := normalized[1]
	return validateOauthToken(oauthToken)
}

func validateOauthToken(token string) (bool, Models.OauthAccessToken) {
	if isJwt(token) {
		return validateSsoOauthToken(token)
	} else {
		return validateAppOauthToken(token)
	}
}

/*
validateAppOauthToken Validates the specified OAuth token, returning whether it points to a valid OAuth token
*/
func validateAppOauthToken(token string) (bool, Models.OauthAccessToken) {
	return validateAccessToken(token)
}

func validateSsoOauthToken(token string) (bool, Models.OauthAccessToken) {
	// TODO: to implement
	return false, Models.OauthAccessToken{}
}

func validateAccessToken(token string) (bool, Models.OauthAccessToken) {
	if len(token) <= AccessTokenPrefixLength {
		return false, Models.OauthAccessToken{}
	}

	tokenName := token[:AccessTokenPrefixLength]
	if tokenName == "" {
		return false, Models.OauthAccessToken{}
	}

	tokenCode := token[AccessTokenPrefixLength:]
	if tokenCode == "" {
		return false, Models.OauthAccessToken{}
	}

	// Search if the token exists in the database
	found, err := Repositories.GetAccessTokenFromName(tokenName)
	if err != nil {
		Logger.Error("Error while searching for token in database: " + err.Error())
		return false, Models.OauthAccessToken{}
	}

	if found.ID == 0 {
		Logger.Error("Token not found in database, authorization failed")
		return false, Models.OauthAccessToken{}
	}

	return true, found
}

func isJwt(token string) bool {
	// TODO: maybe find a better way to validate that is a jwt token
	// A JWT token consists of three parts separated by dots
	parts := strings.Split(token, ".")
	return len(parts) == 3
}
