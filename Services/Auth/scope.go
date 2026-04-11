package Auth

type Scope struct {
	ID          string // ID of the scope, e.g., "repo:read"
	Name        string // Name of the scope, e.g., "ReadRepo"
	Title       string // Title the human-readable name of the scope, e.g., "Read access to repositories"
	Description string // Description of the scope, e.g., "Allows read access to repositories and their metadata."
	Icon        string // Icon font-awesome class for the scope, e.g., "fa-solid fa-lock"
	Dangerous   bool   // Dangerous indicates if the scope is considered dangerous, e.g., true for "user:admin
}

var (
	ReadRepo = Scope{
		ID:          "repo:read",
		Name:        "ReadRepo",
		Title:       "View all visible repositories",
		Description: "This application will be able to view and pull all repositories visible to the granting user",
		Icon:        "fa-hdd-o",
		Dangerous:   false,
	}

	WriteRepo = Scope{
		ID:          "repo:write",
		Name:        "WriteRepo",
		Title:       "Read/Write to any accessible repositories",
		Description: "This application will be able to view, push and pull to all repositories to which the granting user has write access",
		Icon:        "fa-hdd-o",
		Dangerous:   false,
	}

	AdminRepo = Scope{
		ID:          "repo:admin",
		Name:        "AdminRepo",
		Title:       "Administer repositories",
		Description: "This application will have administrator access to all repositories to which the granting user has access",
		Icon:        "fa-hdd-o",
		Dangerous:   false,
	}

	CreateRepo = Scope{
		ID:          "repo:create",
		Name:        "CreateRepo",
		Title:       "Create repositories",
		Description: "This application will be able to create repositories in all namespaces that the granting user is allowed to create repositories",
		Icon:        "fa-plus",
		Dangerous:   false,
	}

	ReadUser = Scope{
		ID:          "user:read",
		Name:        "ReadUser",
		Title:       "Read user information",
		Description: "This application will be abele to read user information such as username and email address",
		Icon:        "fa-user",
		Dangerous:   false,
	}

	AdminUser = Scope{
		ID:          "user:admin",
		Name:        "AdminUser",
		Title:       "Administer user",
		Description: "This application will be able to administer your account including creating robots and granting them permissions to your repositories. You should have absolute trust in the requesting application before granting this permission",
		Icon:        "fa-gear",
		Dangerous:   true,
	}

	OrgAdmin = Scope{
		ID:          "org:admin",
		Name:        "OrgAdmin",
		Title:       "Administer organization",
		Description: "This application will be able to administer your organizations including creating robots, creating teams, adjusting team membership, and changing billing settings. You should have absolute trust in the requesting application before granting this permission",
		Icon:        "fa-gear",
		Dangerous:   true,
	}

	DirectLogin = Scope{
		ID:          "direct_user_login",
		Name:        "DirectLogin",
		Title:       "Full Access",
		Description: "This scope should not available to OAuth applications. Never approve a request for this scope!",
		Icon:        "fa-exclamation-triangle",
		Dangerous:   true,
	}

	SuperUser = Scope{
		ID:          "super:user",
		Name:        "SuperUser",
		Title:       "Super user access",
		Description: "This application will be able to administer your installation including managing users, managing organizations and other features found in the superuser panel. You should have absolute trust in the requesting application before granting this permission",
		Icon:        "fa-street-view",
		Dangerous:   true,
	}
)

func GetScopeFromID(scopeId string) Scope {
	switch scopeId {
	case "repo:read":
		return ReadRepo
	case "repo:write":
		return WriteRepo
	case "repo:admin":
		return AdminRepo
	case "repo:create":
		return CreateRepo
	case "user:read":
		return ReadUser
	case "user:admin":
		return AdminUser
	case "org:admin":
		return OrgAdmin
	case "direct_user_login":
		return DirectLogin
	case "super:user":
		return SuperUser
	default:
		return Scope{}
	}
}
