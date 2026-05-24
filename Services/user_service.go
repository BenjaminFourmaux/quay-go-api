package Services

import (
	"quay-go-api/Common"
	"quay-go-api/Common/Errors"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Repositories"
	"quay-go-api/Services/Auth"
	"quay-go-api/Services/Avatar"
	logger "quay-go-api/Services/Logger"
)

func GetMeInfo(currentUser Auth.AuthenticatedUser) (Dto.UserMeResponse, error) {
	logger.Info("[User Service] Get Me Info")
	logger.Debug("Authenticated user ID: %d", currentUser.ID)

	logger.Info("Retrieving current user with information from database")
	userModel, err := Repositories.GetUserByIdWithUserInformation(currentUser.ID)
	if err != nil {
		logger.Error("Error retrieving current user with information from database: %s", err.Error())
		return Dto.UserMeResponse{}, err
	}
	if userModel.ID == 0 {
		logger.Warning("Current user not found in database: %d", currentUser.ID)
		customErr := Errors.CurrentUserNotFound()
		return Dto.UserMeResponse{}, customErr
	}

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

	userOrgs := []Dto.UserOrganization{}
	if Auth.Can(Auth.ReadUser, currentUser.Scopes) {
		logger.Info("Retrieving user organizations from database")
		orgsModel, orgsErr := Repositories.GetUserOrganizations(userModel.ID)
		if orgsErr != nil {
			logger.Error("Error retrieving user organizations from database: %s", orgsErr.Error())
		} else {
			logger.Debug("User organizations found: %d", len(orgsModel))
			userOrgs = Common.ConvertUserModelsToDto(orgsModel, userModel, currentUser.Scopes)
		}

	} else {
		logger.Debug("ReadUser scope missing, skipping organizations retrieval")
	}

	userDto := Dto.UserMeResponse{
		Anonymous:           false,
		Username:            userModel.Username,
		Avatar:              Avatar.GetAvatarForUser(userModel),
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
