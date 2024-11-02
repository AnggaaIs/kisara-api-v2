package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOAuthConfig *oauth2.Config

func init() {
	GoogleOAuthConfig = &oauth2.Config{
		ClientID:     AppConfig.GoogleClientID,
		ClientSecret: AppConfig.GoogleClientSecret,
		RedirectURL:  AppConfig.GoogleRedirectURL,
		Endpoint:     google.Endpoint,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	}
}
