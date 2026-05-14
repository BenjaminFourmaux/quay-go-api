package Avatar

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Entities/Models"
	"strconv"
	"strings"
)

var AvatarColors = []string{"#969696", "#aec7e8", "#ff7f0e", "#ffbb78", "#2ca02c", "#98df8a", "#d62728", "#ff9896", "#9467bd", "#c5b0d5", "#8c564b", "#c49c94",
	"#e377c2", "#f7b6d2", "#7f7f7f", "#c7c7c7", "#bcbd22", "#1f77b4", "#17becf", "#9edae5", "#393b79", "#5254a3", "#6b6ecf", "#9c9ede", "#9ecae1", "#31a354",
	"#b5cf6b", "#a1d99b", "#8c6d31", "#ad494a", "#e7ba52", "#a55194"}

/*
Computes and returns the full data block for the avatar
*/
func getAvatar(name string, emailOrId string, kind string) Dto.Avatar {

	// NOTE: email_or_id can be null if gotten from external auth email is disabled, so use the username in that case
	var usernameEmailOrId string
	if emailOrId == "" {
		usernameEmailOrId = name
	} else {
		usernameEmailOrId = emailOrId
	}

	hashValue := computeHash(usernameEmailOrId)
	hashColor := getHashColor(hashValue, AvatarColors)

	return Dto.Avatar{
		Name:  name,
		Hash:  hashValue,
		Color: hashColor,
		Kind:  kind,
	}
}

func GetAvatarForUser(user Models.User) Dto.Avatar {
	var kind string
	if user.Robot {
		kind = "robot"
	} else {
		kind = "user"
	}
	return getAvatar(user.Username, user.Email, kind)
}

func GetAvatarForOrg(org Models.User) Dto.Avatar {
	return getAvatar(org.Username, org.Email, "org")
}

func GetAvatarForTeam(team Models.Team) Dto.Avatar {
	return getAvatar(team.Name, strconv.Itoa(team.ID), "team")
}

func computeHash(usernameEmailOrId string) string {
	// Trim whitespace, convert to lowercase, and encode to UTF-8
	processedInput := strings.TrimSpace(strings.ToLower(usernameEmailOrId))

	// Compute SHA256 hash
	hash := sha256.Sum256([]byte(processedInput))

	// Convert hash to hexadecimal string
	return hex.EncodeToString(hash[:])
}

func getHashColor(hashValue string, colors []string) string {
	// Calculate the number of bytes needed
	byteCount := int(math.Ceil(math.Log(float64(len(colors))) / math.Log(16)))

	// Extract the required number of bytes from the hash
	byteData, _ := hex.DecodeString(hashValue[:byteCount*2])

	// Convert the byte data to an integer
	var intValue int
	for _, b := range byteData {
		intValue = (intValue << 8) + int(b)
	}

	// Get the corresponding color
	hashColor := colors[intValue%len(colors)]

	return hashColor
}
