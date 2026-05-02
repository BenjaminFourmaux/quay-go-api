package Auth

/*
AuthenticatedUser struct for storing information (scopes and id) of the authenticated user
*/
type AuthenticatedUser struct {
	ID     int
	Scopes []Scope
}
