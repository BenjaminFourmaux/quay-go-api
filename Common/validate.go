package Common

import (
	"net/mail"
	"quay-go-api/Common/Errors"
	"quay-go-api/Entities/Dto"
	"regexp"
	"strings"
)

/*
ValidateCreateOrganization organization metadata for creating a new organization. Rules:
The organization name must:
1. Not be empty
2. Be between 2 and 255 characters long
3. Contain only alphanumeric characters, dashes, or underscores
*/
func ValidateCreateOrganization(organizationMetadata Dto.CreateOrganization) error {
	// Validate org name
	// 1. empty value
	if organizationMetadata.Name == "" {
		return Errors.OrganizationNameInvalid()
	}

	// 2. value length
	if len(organizationMetadata.Name) < 2 || len(organizationMetadata.Name) > 255 {
		return Errors.OrganizationNameInvalid()
	}

	// 3. valid characters (alphanumeric, dash and underscore)
	for _, char := range organizationMetadata.Name {
		if !(char >= 'a' && char <= 'z') && !(char >= '0' && char <= '9') && char != '-' && char != '_' {
			return Errors.OrganizationNameInvalid()
		}
	}

	return nil
}

func ValidateUpdateOrganization(organizationMetadata Dto.UpdateOrganization) error {
	if organizationMetadata.Email != nil {
		email := strings.TrimSpace(*organizationMetadata.Email)
		if email == "" {
			return Errors.OrganizationEmailInvalid()
		}

		if _, err := mail.ParseAddress(email); err != nil {
			return Errors.OrganizationEmailInvalid()
		}
	}

	if organizationMetadata.InvoiceEmailAddress != nil {
		email := strings.TrimSpace(*organizationMetadata.InvoiceEmailAddress)
		if email == "" {
			return Errors.OrganizationEmailInvalid()
		}

		if _, err := mail.ParseAddress(email); err != nil {
			return Errors.OrganizationEmailInvalid()
		}
	}

	if organizationMetadata.TagExpirationS != nil && *organizationMetadata.TagExpirationS < 0 {
		return Errors.OrganizationTagExpirationInvalid()
	}

	return nil
}

/*
ValidateRole cheks if the role is valid (e.g., "owners", "admin", "member")
*/
func ValidateRole(role string) bool {
	return role == "admin" || role == "creator" || role == "member"
}

func ValidateTeam(team Dto.CreateTeam) error {
	// Validate team name (required)
	if team.Name == nil ||
		*team.Name == "" ||
		strings.TrimSpace(*team.Name) == "" {
		return Errors.TeamNameRequired()
	}

	var reName = regexp.MustCompile(`^[a-z][a-z0-9]+$`)
	if (len(*team.Name) < 2 && len(*team.Name) > 255) ||
		!reName.MatchString(*team.Name) {
		return Errors.TeamNameInvalid()
	}

	// Validate Role (optional)
	if team.Role != nil && !ValidateRole(*team.Role) {
		return Errors.TeamRoleInvalid()
	}

	return nil
}
