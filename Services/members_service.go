package Services

import (
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
