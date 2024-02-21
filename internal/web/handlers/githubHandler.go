package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"forum/internal/models"
	"forum/internal/web/handlers/helpers"
	"io/ioutil"
	"net/http"
)

func (h *Handler) GithubAuthHandler(w http.ResponseWriter, r *http.Request) {
	// prompt=consent    ---- this is needed for prompting the user to confirm authorisation
	url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=offline&prompt=consent", models.GitHubAuthURL, models.GitHubClientID, models.GitHubRedirectURL)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) GithubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code") // temporary token given by Github
	if code == "" {
		helpers.ErrorHandler(w, http.StatusUnauthorized, errors.New("Temporary token is invalid"))
	}

	tokenRes, err := getGithubOauthToken(code)
	if err != nil {
		helpers.ErrorHandler(w, http.StatusBadGateway, errors.New("The information received from Github"))
		return
	}

	githubData, err := getGithubData(tokenRes.AccessToken)
	if err != nil {
		helpers.ErrorHandler(w, http.StatusBadGateway, errors.New("The information received from GitHub"))
		return
	}

	userData, err := getUserData(githubData)
	if err != nil {
		helpers.ErrorHandler(w, http.StatusBadGateway, errors.New("The information received from GitHub"))
		return
	}

	session, err := h.service.GitHubAuthorization(&userData)
	if err != nil {
		helpers.ErrorHandler(w, http.StatusBadRequest, err)
		return
	} else {
		helpers.SessionCookieSet(w, session.Token, session.ExpTime)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

func getUserData(data string) (models.GitHubLoginUserData, error) {
	userData := models.GitHubLoginUserData{}
	if err := json.Unmarshal([]byte(data), &userData); err != nil {
		return models.GitHubLoginUserData{}, err
	}

	return userData, nil
}

func getGithubOauthToken(code string) (*models.GitHubResponseToken, error) {
	requestBodyMap := map[string]string{
		"client_id":     models.GitHubClientID,
		"client_secret": models.GitHubClientSecret,
		"code":          code,
	}
	requestJSON, err := json.Marshal(requestBodyMap)
	if err != nil {
		return nil, err
	}

	req, reqerr := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer(requestJSON),
	)
	if reqerr != nil {
		return nil, reqerr
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		return nil, resperr
	}

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ghresp models.GitHubResponseToken
	if err := json.Unmarshal(respbody, &ghresp); err != nil {
		return nil, err
	}

	return &ghresp, nil
}

func getGithubData(accessToken string) (string, error) {
	req, reqerr := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if reqerr != nil {
		return "", reqerr
	}

	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		return "", resperr
	}

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respbody), nil
}
