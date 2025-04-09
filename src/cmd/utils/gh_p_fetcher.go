package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type GithubProfile struct {
	Login string `json:"login"`
	Name string `json:"name"`
}

func FetchGithubProfile(username string) User {
	// Fetch the github profile and create a user with everything except the email

	url := fmt.Sprintf("https://api.github.com/users/%s", username)

	resp, err := http.Get(url)
	if err != nil {
		panic(fmt.Sprint("Error fetching github profile: ", err))
	}
	defer resp.Body.Close()

	// Parse the response and create a user
	var profile GithubProfile
	err = json.NewDecoder(resp.Body).Decode(&profile)
	if err != nil {
		panic(fmt.Sprint("Error parsing github profile: ", err))
	}

	// Create a user with the github profile
	return User{
		Shortname: strings.ToLower(profile.Name[:2]),	
		Longname: profile.Name,
		Username: profile.Login,
		Email: "",
		Ex: false,
		Groups: []string{},
	}

}

