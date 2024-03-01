package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"forum/internal/models"
	"forum/internal/web/handlers/helpers"
	"net/http"
	"net/url"
	"strings"
)

func (h *Handler) GoogleAuthHandler(w http.ResponseWriter, r *http.Request) {
	scope := url.QueryEscape("email profile https://www.googleapis.com/auth/drive.file")
	url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&prompt=select_account", models.GoogleAuthURL, models.GoogleClientID, models.GoogleRedirectURL, scope)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code") // temporary token given by Google

	// fmt.Println("Code: ", code) // 4/0AfJohXnLtvZf6XVZjrQRMyaH_CdPg6yB7XoUnrdBqv0wd6RkJDo361ff3yc6qzaMdS6oKQ

	if code == "" {
		helpers.ErrorHandler(w, http.StatusUnauthorized, errors.New("Temporary Token is invalid"))
		return
	}

	tokenRes, err := getGoogleOauthToken(code)
	if err != nil {
		helpers.ErrorHandler(w, http.StatusBadGateway, errors.New("Token wasn't received properly!"))
		return
	}

	googleUser, err := getGoogleUser(tokenRes.AccessToken, tokenRes.TokenID)
	if err != nil {
		helpers.ErrorHandler(w, http.StatusBadGateway, errors.New("The information received from Google"))
		return
	}

	// creating the struct type of User after Google Auth
	googleData := models.GoogleLoginUserData{
		Name:       googleUser.Name,
		Email:      googleUser.Email,
		FirstName:  googleUser.Given_name,
		SecondName: googleUser.Family_name,
		Provider:   "Google",
		Password:   googleUser.Password,
	}

	session, err := h.service.GoogleAuthorization(&googleData)
	if err != nil {
		helpers.ErrorHandler(w, http.StatusBadRequest, err)
		return
	} else {
		helpers.SessionCookieSet(w, session.Token, session.ExpTime)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

func getGoogleOauthToken(code string) (*models.GoogleResponseToken, error) {
	// The URL to which the POST request is sent. This is the token endpoint of the OAuth 2.0 provider (in this case, Google).
	const rootURL = "https://oauth2.googleapis.com/token" // URL for getting the access token using code

	// map is convenient for encoding form values in the request body
	values := url.Values{}
	values.Set("code", code)
	values.Set("client_id", models.GoogleClientID)
	values.Set("client_secret", models.GoogleClientSecret)
	values.Set("redirect_uri", models.GoogleRedirectURL)
	values.Set("grant_type", "authorization_code")

	// Make a POST request to the Google token endpoint
	response, err := http.Post(rootURL, "application/x-www-form-urlencoded", strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var tokenRes models.GoogleResponseToken
	err = json.NewDecoder(response.Body).Decode(&tokenRes)
	if err != nil {
		return nil, err
	}

	return &tokenRes, nil
}

func getGoogleUser(AccessToken string, TokenID string) (*models.GoogleUserResult, error) {
	rootURL := fmt.Sprintf("https://www.googleapis.com/oauth2/v1/userinfo?alt=json&access_token=%s", AccessToken)

	// Prepare the request to the Google People API
	req, err := http.NewRequest("GET", rootURL, nil)
	if err != nil {
		return nil, err
	}

	// Set the Authorization header with the access token
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", TokenID))

	// Make the request to the Google People API
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var userRes models.GoogleUserResult
	err = json.NewDecoder(res.Body).Decode(&userRes)
	if err != nil {
		return nil, err
	}

	return &userRes, nil
}
