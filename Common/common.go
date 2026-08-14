package Common

import (
	"crypto/rand"
	"math/big"
	"quay-go-api/Entities/Models"
	"quay-go-api/Services/Auth"
	"strings"
)

func HasScope(scopes []Auth.Scope, scope Auth.Scope) bool {
	for _, scp := range scopes {
		if scp.ID == scope.ID {
			return true
		}
	}
	return false
}

/*
CanViewTeams checks if the user can view the team
A user can view a team if:
1. They are a member of that team (any role)
2. They are the scope org:admin
*/
func CanViewTeams(userId int, team Models.Team, userScopes []Auth.Scope) bool {
	if Auth.Can(Auth.OrgAdmin, userScopes) {
		return true
	}

	if team.Members == nil {
		panic("team members should be preloaded")
	}
	for _, teamMember := range team.Members {
		if teamMember.UserId == userId {
			return true
		}
	}
	return false
}

func GetTeamRoleIdFromRoleName(roleName string) int {
	switch roleName {
	case "admin":
		return 1
	case "creator":
		return 2
	case "member":
		return 3
	default:
		return 0
	}
}

func MapRepositoryStateStr(stateId int) string {
	switch stateId {
	case 0:
		return "NORMAL" // Regular repo where all actions are possible
	case 1:
		return "READ_ONLY" // Only read actions, such as pull, are allowed regardless of specific user permissions
	case 2:
		return "MIRROR" // Equivalent to READ_ONLY except that mirror robot has write permission
	case 3:
		return "MARKED_FOR_DELETION" // Indicates the repository has been marked for deletion and should be hidden and unusable.
	case 4:
		return "ORG_MIRROR" // Equivalent to MIRROR but for repositories created via organization-level mirroring
	default:
		return "UNKNOWN"
	}
}

func GetRoleIdFromRoleName(roleName string) int {
	switch roleName {
	case "admin":
		return 1
	case "write":
		return 2
	case "read":
		return 3
	default:
		return 0
	}
}

/*
MapMediaTypeName maps a media type ID to its corresponding MediaType Name.
Hardcoded values in Database
*/
func MapMediaTypeName(mediaTypeId int) string {
	switch mediaTypeId {
	case 1:
		return "text/plain"
	case 2:
		return "application/json"
	case 3:
		return "text/markdown"
	case 4:
		return "application/vnd.cnr.blob.v0.tar+gzip"
	case 5:
		return "application/vnd.cnr.package-manifest.helm.v0.json"
	case 6:
		return "application/vnd.cnr.package-manifest.kpm.v0.json"
	case 7:
		return "application/vnd.cnr.package-manifest.docker-compose.v0.json"
	case 8:
		return "application/vnd.cnr.package.kpm.v0.tar+gzip"
	case 9:
		return "application/vnd.cnr.package.helm.v0.tar+gzip"
	case 10:
		return "application/vnd.cnr.package.docker-compose.v0.tar+gzip"
	case 11:
		return "application/vnd.cnr.manifests.v0.json"
	case 12:
		return "application/vnd.cnr.manifest.list.v0.json"
	case 13:
		return "application/vnd.docker.distribution.manifest.v1+json"
	case 14:
		return "application/vnd.docker.distribution.manifest.v1+prettyjws"
	case 15:
		return "application/vnd.docker.distribution.manifest.v2+json"
	case 16:
		return "application/vnd.docker.distribution.manifest.list.v2+json"
	case 17:
		return "application/vnd.oci.image.index.v1+json"
	case 18:
		return "application/vnd.oci.image.manifest.v1+json"
	default:
		return "unknown"
	}
}

func MapLabelSourceType(sourceTypeId int) Models.LabelSourceType {
	switch sourceTypeId {
	case 1:
		return Models.LabelSourceType{ID: 1, Name: "manifest", Mutable: false}
	case 2:
		return Models.LabelSourceType{ID: 2, Name: "api", Mutable: true}
	case 3:
		return Models.LabelSourceType{ID: 3, Name: "internal", Mutable: false}
	default:
		return Models.LabelSourceType{ID: 0, Name: "unknown", Mutable: false}
	}
}

func MapLoginServiceName(loginServiceId int) string {
	switch loginServiceId {
	case 1:
		return "github"
	case 2:
		return "quayrobot"
	case 3:
		return "ldap"
	case 4:
		return "google"
	case 5:
		return "keystone"
	case 6:
		return "dex"
	case 7:
		return "jwtauthn"
	default:
		return "unknown"
	}
}

func MapLoginServiceId(loginServiceName string) int {
	switch strings.ToLower(loginServiceName) {
	case "github":
		return 1
	case "quayrobot":
		return 2
	case "ldap":
		return 3
	case "google":
		return 4
	case "keystone":
		return 5
	case "dex":
		return 6
	case "jwtauthn":
		return 7
	default:
		return 0
	}
}

func FormatRobotUsername(username string, robotName string) string {
	return username + "+" + robotName
}

func RandomStringGenerator(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" // ascii to upper + digits
	result := make([]byte, length)
	for i := range result {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[n.Int64()]
	}
	return string(result)
}
