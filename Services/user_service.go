package Services

import (
	"quay-go-api/Entities/Dto"
	"quay-go-api/Repositories"
	"quay-go-api/Services/Auth"
)

func GetMeInfo(userId int, userScopes []Auth.Scope) (Dto.UserMeResponse, error) {
	userModel := Repositories.GetUserByIdWithUserInformation(userId)

	// Convert models to dto
	userLogins := []Dto.UserLogin{}
	for _, federatedLogin := range userModel.FederatedLogins {
		userLogins = append(userLogins, Dto.UserLogin{
			Service:           federatedLogin.Service.Name,
			ServiceIdentifier: federatedLogin.ServiceIdent,
			Metadata:          federatedLogin.MetadataJson,
		})
	}

	userPrompts := []string{}
	for _, prompt := range userModel.Prompts {
		userPrompts = append(userPrompts, prompt.Kind.Name)
	}

	userOrgs := []Dto.Organization{}
	if Auth.Can(Auth.ReadUser, userScopes) {
		orgsModel, _ := Repositories.GetUserOrganizations(userModel.Username)

		for _, org := range orgsModel {
			userOrgs = append(userOrgs, Dto.Organization{
				Name:               org.Username,
				Avatar:             GetAvatarForOrg(org),
				CanCreateRepo:      Auth.Can(Auth.CreateRepo, userScopes),
				Public:             false, // TODO: check if the org name not in list of public Namespaces
				IsOrgAdmin:         Auth.Can(Auth.OrgAdmin, userScopes),
				PreferredNamespace: !(!userModel.StripeId.Valid || userModel.StripeId.String == ""),
			})
		}
	}

	userDto := Dto.UserMeResponse{
		Anonymous:           false,
		Username:            userModel.Username,
		Avatar:              GetAvatarForUser(userModel),
		CanCreateRepo:       true,
		IsMe:                true, // get me
		Verified:            userModel.Verified,
		Email:               userModel.Email,
		Logins:              userLogins,
		InvoiceEmail:        userModel.InvoiceEmail,
		InvoiceEmailAddress: Dto.NullString(userModel.InvoiceEmailAddress),
		PreferredNamespace:  !(!userModel.StripeId.Valid || userModel.StripeId.String == ""),
		TagExpirationS:      userModel.RemovedTagExpirationS,
		Prompts:             userPrompts,
		Company:             Dto.NullString(userModel.Company),
		FamilyName:          Dto.NullString(userModel.FamilyName),
		GivenName:           Dto.NullString(userModel.GivenName),
		Location:            Dto.NullString(userModel.Location),
		IsFreeAccount:       !userModel.StripeId.Valid || userModel.StripeId.String == "", // if stripe id is empty, it's a free account
		HasPasswordSet:      userModel.PasswordHash.Valid && userModel.PasswordHash.String != "",
		Organizations:       userOrgs,
		SuperUser:           true, // where is this info ?
	}

	return userDto, nil
}
