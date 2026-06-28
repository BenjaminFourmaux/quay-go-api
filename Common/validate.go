package Common

import (
	"net/mail"
	"quay-go-api/Common/Errors"
	"quay-go-api/Entities/Dto"
	"regexp"
	"strings"
)

/* DEV NOTE
- function to check if a field is valid must be named like: IsValide<FieldName> and return a bool
- function to validate a Dto must be named like: Validate<DtoName> and return an error
*/

/*
ValidateMessageSeverity checks if the severity is valid (e.g., "info", "warning", "error") and return true if the severity given is a valid Message severity, otherwise false
*/
func ValidateMessageSeverity(severity string) bool {
	return severity == "info" || severity == "warning" || severity == "error"
}

/*
ValidateCreateOrganization organization metadata for creating a new organization.
*/
func ValidateCreateOrganization(organizationMetadata Dto.CreateOrganization) error {
	// Validate org name
	if !IsValidOrganizationOrUserName(organizationMetadata.Name) {
		return Errors.OrganizationNameInvalid()
	}

	return nil
}

/*
IsValidOrganizationOrUserName checks if the organization or user name is valid (e.g., "my-org", "user_name", "username123")
Rules:
1. Not be empty
2. Be between 2 and 255 characters long
3. Contain only alphanumeric characters, dashes, or underscores
*/
func IsValidOrganizationOrUserName(name string) bool {
	// 1. empty value
	if name == "" {
		return false
	}

	// 2. value length
	if len(name) < 2 || len(name) > 255 {
		return false
	}

	// 3. valid characters (alphanumeric, dash and underscore)
	for _, char := range name {
		if !(char >= 'a' && char <= 'z') && !(char >= '0' && char <= '9') && char != '-' && char != '_' {
			return false
		}
	}

	return true
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
ValidateTeamRole cheks if the role is valid (e.g., "owners", "admin", "member")
*/
func ValidateTeamRole(role string) bool {
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
	if team.Role != nil && !ValidateTeamRole(*team.Role) {
		return Errors.TeamRoleInvalid()
	}

	return nil
}

func ValidateCreateRepository(repositoryMetadata Dto.CreateRepository) error {
	// Validate repository name
	if !IsValidRepositoryName(repositoryMetadata.Name) {
		return Errors.RepositoryNameInvalid()
	}

	// Validate Namespace if present
	if repositoryMetadata.Namespace != nil && !IsValidOrganizationOrUserName(*repositoryMetadata.Namespace) {
		return Errors.RepositoryNamespaceInvalid()
	}

	// Validate kind
	if !IsValidRepositoryKind(repositoryMetadata.Kind) {
		return Errors.RepositoryKindInvalid()
	}

	return nil
}

/*
IsValidRepositoryKind cheks if the kind is valid (e.g., "image", or "application")
*/
func IsValidRepositoryKind(kind string) bool {
	return kind == "image" || kind == "application"
}

/*
IsValidRepositoryName checks if the repository name is valid (e.g., "my-repo", "user_name/my-repo", "username123/my-repo")
*/
func IsValidRepositoryName(repositoryName string) bool {
	repositoryName = strings.TrimSpace(repositoryName)
	if repositoryName == "" {
		return false
	}

	var reRepositoryName = regexp.MustCompile(`^[a-z0-9]+(?:(?:[._]|__|[-]+)[a-z0-9]+)*(?:/[a-z0-9]+(?:(?:[._]|__|[-]+)[a-z0-9]+)*)*$`)
	if len(repositoryName) > 255 || !reRepositoryName.MatchString(repositoryName) {
		return false
	}
	return true
}
