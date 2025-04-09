package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

)

type GithubProfile struct {
	Login string `json:"login"`
	Name string `json:"name"`
	Email string `json:"email"`
}

func checkGHCLI() bool {
	// Check if the gh command line tool is installed
	cmd := exec.Command("gh", "auth", "status")
	out ,err := cmd.CombinedOutput()
	if err == nil {
		if strings.Contains(string(out), "Logged in to") {
			return true
		}
	} else {
		return false
	}

	return false
	
}

func useGHCLI(username string) []byte {
	cmd := exec.Command("gh", "api", fmt.Sprintf("/users/%s", username))
	
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprint("Error fetching github profile",err))
	}
	return out

}

func FetchGithubProfile(username string) User {
	// Fetch the github profile and create a user with everything except the email

	var profile GithubProfile

	check := checkGHCLI()

	if check {
		// If the gh command line tool is installed, use it to fetch the github profile
		fmt.Println("Using gh-cli to fetch github profile")
		data := useGHCLI(username)

		err := json.Unmarshal(data, &profile)
		if err != nil {
			panic(fmt.Sprint("Error parsing github profile: ", err))
		}
		if profile.Name == "" {
			panic(fmt.Sprint("Error: No name found in github profile something went wrong whilst fetching using gh-cli \n(output from cmd:)", string(data)))
		}
	} else {

		fmt.Println("Using http request to fetch github profile")
		// If the gh command line tool is not installed, use the http request
		url := fmt.Sprintf("https://api.github.com/users/%s", username)

		resp, err := http.Get(url)
		if err != nil {
			panic(fmt.Sprint("Error fetching github profile: ", err))
		}
		defer resp.Body.Close()

		// Parse the response and create a user
		err = json.NewDecoder(resp.Body).Decode(&profile)
		if err != nil {
			panic(fmt.Sprint("Error parsing github profile: ", err))
		}

		if profile.Name == "" {
			panic(fmt.Sprintf("Error: No name found in github profile something went wrong whilst fetching \n(http response): %s", resp.Status))
		}
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

