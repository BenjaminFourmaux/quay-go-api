package Auth

/*
AuthenticatedUser struct for storing information (scopes and id) of the authenticated user
*/
type AuthenticatedUser struct {
	ID       int     // User ID in the Database
	Username string  // User username in the Database (username or robot account name)
	Scopes   []Scope // List of scopes the user has
	Token    string  // Token used for authentication
}
