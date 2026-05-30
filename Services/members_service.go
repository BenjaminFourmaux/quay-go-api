package Services

import (
	"fmt"
	"quay-go-api/Common/Errors"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Entities/Models"
	"quay-go-api/Repositories"
	"quay-go-api/Services/Auth"
	"quay-go-api/Services/Avatar"
	logger "quay-go-api/Services/Logger"
)

func ListTeamMembers(orgName string, teamName string, filters map[string]string, currentUser Auth.AuthenticatedUser) ([]Dto.TeamMember, error) {
	logger.Info("[Members Service] List Team Members")
	logger.Debug("Organization name: %s", orgName)
	logger.Debug("Team name: %s", teamName)
	logger.Debug("With filters: %+v", filters)

	// Validating filters
	var filterIsRobot bool = false
	var filterIsInvited bool = false
	if isRobot, ok := filters["is_robot"]; ok {
		filterIsRobot = isRobot == "true"
	}
	if isInvited, ok := filters["is_invited"]; ok {
		filterIsInvited = isInvited == "true"
	}

	// Retrieve organization and check if exists
	logger.Info("Retrieving organization details from database")
	organizationModel, err := Repositories.GetOrganizationDetailsByName(orgName)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("Organization not found: %s", orgName)
			return nil, Errors.OrganizationNotFound(orgName)
		default:
			logger.Error("Error retrieving organization details from database: %s", err.Error())
			return nil, err
		}
	}

	// Browse teams and check if organization's team exist
	logger.Info("Checking if team exists in the organization")
	for _, team := range organizationModel.Teams {
		if team.Name == teamName {
			// Retrieve team with details (members, member.user, invite)
			logger.Info("Retrieving team details from database")
			teamDetailsModel, err := Repositories.GetTeamDetailsById(team.ID)
			if err != nil {
				switch err.Error() {
				case "record not found": // Not possible but "on sait jamais"
					logger.Warning("Team not found: %s", teamName)
					return nil, Errors.TeamNotFound(teamName)
				default:
					logger.Error("Error retrieving team details from database: %s", err.Error())
					return nil, err
				}
			}

			members := []Dto.TeamMember{}

			// Browse team members
			for _, member := range teamDetailsModel.Members {
				// Apply filters
				if filterIsRobot && !member.User.Robot { // keep only robot user
					continue
				}
				if filterIsInvited && !isInvitedMember(member.UserId, team.ID, team.TeamMemberInvites) {
					continue
				}

				members = append(members, Dto.TeamMember{
					Name:    member.User.Username,
					Kind:    "user", // can be something else ?
					IsRobot: member.User.Robot,
					Avatar:  Avatar.GetAvatarForUser(member.User),
					Invited: isInvitedMember(member.UserId, team.ID, team.TeamMemberInvites),
				})
			}

			return members, nil
		}
	}

	logger.Warning("Team not found: %s", teamName)
	return nil, Errors.TeamNotFound(teamName)
}

func AddMemberToTeam(orgName string, teamName string, memberName string, currentUser Auth.AuthenticatedUser) (Dto.TeamMember, error) {
	logger.Info("[Members Service] Add Team Member")
	logger.Debug("Organization name: %s", orgName)
	logger.Debug("Team name: %s", teamName)
	logger.Debug("Username: %s", memberName)

	// 1. Get organization by name and check if exists
	logger.Info("Get organization and check if exists")
	organizationModel, err := Repositories.GetOrganizationByName(orgName)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("Organization not found: %s", orgName)
			return Dto.TeamMember{}, Errors.OrganizationNotFound(orgName)
		default:
			logger.Error("Error retrieving organization details from database: %s", err.Error())
			return Dto.TeamMember{}, err
		}
	}
	logger.Debug("Organization found: %s (ID: %d)", organizationModel.Username, organizationModel.ID)

	// 2. Get organization teams with members
	logger.Info("Get organization teams")
	teamsModel, err := Repositories.GetOrganizationTeamsByOrgId(organizationModel.ID)
	if err != nil {
		switch err.Error() {
		case "record not found": // not possible but "on sait jamais"
			logger.Warning("Teams not found in organization: %s", orgName)
			return Dto.TeamMember{}, Errors.TeamNotFound(orgName)
		default:
			logger.Error("Error retrieving organization details from database: %s", err.Error())
			return Dto.TeamMember{}, err
		}
	}
	logger.Debug("Teams found in organization: %d", len(teamsModel))

	organizationModel.Teams = teamsModel

	// 3. Check if team exist
	logger.Info("Checking if team exists in the organization")
	for _, team := range organizationModel.Teams {
		if team.Name == teamName {
			// Team matching, continue the treatment
			logger.Info("Check if the current user can perform this action")
			if !Auth.Can(Auth.OrgAdmin, currentUser.Scopes) || !isUserIsOrgOwner(currentUser.ID, organizationModel) {
				logger.Warning("Current user cannot perform this action (not Scope OrgAdmin or not organization owner)")
				return Dto.TeamMember{}, Errors.UnauthorizedInsufficientRole()
			}

			// 4. Check if the user is already in this team
			logger.Info("Check if the current user is in the team")
			for _, member := range team.Members {
				if member.User.Username == memberName {
					logger.Warning("User is already in this team")
					return Dto.TeamMember{}, Errors.MemberAlreadyInTeam()
				}
			}

			// 5. Check if the user exists
			logger.Info("Check if the user to add exists")
			userModel, err := Repositories.GetUserByName(memberName)
			if err != nil {
				switch err.Error() {
				case "record not found":
					logger.Warning("User not found: %s", memberName)
					return Dto.TeamMember{}, Errors.UserNotExists(memberName)
				default:
					logger.Error("Error retrieving user details from database: %s", err.Error())
					return Dto.TeamMember{}, err
				}
			}

			// 6. All goods, add the user to the team
			logger.Info("Adding user to the team")
			teamMemberToAdd := Models.TeamMember{
				UserId: userModel.ID,
				TeamId: team.ID,
			}
			teamMemberCreated, err := Repositories.CreateTeamMember(teamMemberToAdd)
			if err != nil {
				return Dto.TeamMember{}, err
			}
			logger.Debug("Team member created with ID: %d", teamMemberCreated.ID)

			// 7. Make a TeamMember Dto
			teamMember := Dto.TeamMember{
				Name:    userModel.Username,
				Kind:    "user",
				IsRobot: false,
				Avatar:  Avatar.GetAvatarForUser(userModel),
				Invited: false,
			}
			return teamMember, nil
		}
	}

	return Dto.TeamMember{}, Errors.TeamNotFound(teamName)
}

func RemoveMemberToTeam(orgName string, teamName string, memberName string, currentUser Auth.AuthenticatedUser) error {
	logger.Info("[Members Service] Remove Team Member")
	logger.Debug("Organization name: %s", orgName)
	logger.Debug("Team name: %s", teamName)
	logger.Debug("Username: %s", memberName)

	// 1. Get organization by name and check if exists
	logger.Info("Get organization and check if exists")
	organizationModel, err := Repositories.GetOrganizationByName(orgName)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("Organization not found: %s", orgName)
			return Errors.OrganizationNotFound(orgName)
		default:
			logger.Error("Error retrieving organization details from database: %s", err.Error())
			return err
		}
	}
	logger.Debug("Organization found: %s (ID: %d)", organizationModel.Username, organizationModel.ID)

	// 2. Get organization teams with members
	logger.Info("Get organization teams")
	teamsModel, err := Repositories.GetOrganizationTeamsByOrgId(organizationModel.ID)
	if err != nil {
		switch err.Error() {
		case "record not found": // not possible but "on sait jamais"
			logger.Warning("Teams not found in organization: %s", orgName)
			return Errors.TeamNotFound(orgName)
		default:
			logger.Error("Error retrieving organization details from database: %s", err.Error())
			return err
		}
	}
	logger.Debug("Teams found in organization: %d", len(teamsModel))

	organizationModel.Teams = teamsModel

	// 3. Check if team exist
	logger.Info("Checking if team exists in the organization")
	for _, team := range organizationModel.Teams {
		if team.Name == teamName {
			// Team matching, continue the treatment
			logger.Info("Check if the current user can perform this action")
			if !Auth.Can(Auth.OrgAdmin, currentUser.Scopes) || !isUserIsOrgOwner(currentUser.ID, organizationModel) {
				logger.Warning("Current user cannot perform this action (not Scope OrgAdmin or not organization owner)")
				return Errors.UnauthorizedInsufficientRole()
			}

			// 4. Check if the user is on the team
			logger.Info("Check if the current user is in the team")
			for _, member := range team.Members {
				if member.User.Username == memberName {
					// 5. Remove it
					logger.Info(fmt.Sprintf("Removing user %d from the team", member.UserId))
					err = Repositories.DeleteTeamMember(member)
					if err != nil {
						logger.Error("Error removing user from the team: %s", err.Error())
						return err
					}
					return nil
				} else {
					logger.Warning("User is not on the team")
					return Errors.MemberNotInTeam()
				}
			}
		}
	}

	return Errors.TeamNotFound(teamName)
}

// <editor-fold desc="Private Methods">

func isInvitedMember(userId int, teamId int, invites []Models.TeamMemberInvite) bool {
	for _, invite := range invites {
		if invite.TeamId == teamId {
			if invite.UserId == userId {
				return true
			}
		}
	}
	return false
}

// </editor-fold>
