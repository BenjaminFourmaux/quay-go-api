package Repositories

import (
	"quay-go-api/Database"
	"quay-go-api/Entities/Models"
)

func GetAccessTokenFromName(tokenName string) (Models.OauthAccessToken, error) {
	var accessToken Models.OauthAccessToken

	err := Database.DB.Where("token_name = ?", tokenName).First(&accessToken).Error

	return accessToken, err
}
