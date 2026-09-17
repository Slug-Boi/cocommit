package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/google/go-github/v91/github"
)

type GithubProfile struct {
	Login string `json:"login"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func checkGHCLI() bool {
	// Check if the gh command line tool is installed
	cmd := exec.Command("gh", "auth", "status")
	out, err := cmd.CombinedOutput()
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
		panic(fmt.Sprint("Error fetching github profile", err))
	}
	return out

}

func FetchGithubProfile(ctx context.Context, username string) (User, error) {
	var profile GithubProfile
	var err error

	if ctx == nil {
		fmt.Println("No context provided, using gh-cli to fetch github profile")

		if !checkGHCLI() {
			return User{}, fmt.Errorf("fetching github profile: no context provided and gh-cli is not installed")
		}

		data := useGHCLI(username)
		if err := json.Unmarshal(data, &profile); err != nil {
			return User{}, fmt.Errorf("parsing gh-cli output: %w", err)
		}
	} else {
		profile, err = fetchGithubProfileViaAPI(ctx, username)
		if err != nil {
			fmt.Println("go-github fetch failed, falling back to gh-cli:", err)

			if !checkGHCLI() {
				return User{}, fmt.Errorf("fetching github profile: go-github failed (%w) and gh-cli is not installed", err)
			}

			fmt.Println("Using gh-cli to fetch github profile")
			data := useGHCLI(username)
			if err := json.Unmarshal(data, &profile); err != nil {
				return User{}, fmt.Errorf("parsing gh-cli output: %w", err)
			}
		} else {
			fmt.Println("Using go-github to fetch github profile")
		}
	}

	if profile.Name == "" {
		return User{}, fmt.Errorf("no name found in github profile for %q — something went wrong fetching the profile", username)
	}

	shortname := profile.Name
	if len(shortname) > 2 {
		shortname = shortname[:2]
	}

	return User{
		Shortname: strings.ToLower(shortname),
		Longname:  profile.Name,
		Username:  profile.Login,
		Email:     profile.Email,
		Ex:        false,
		Groups:    []string{},
		Platform:  "github",
	}, nil
}

func fetchGithubProfileViaAPI(ctx context.Context, username string) (GithubProfile, error) {
	tok, err := Login(ctx)
	if err != nil {
		return GithubProfile{}, fmt.Errorf("authenticating: %w", err)
	}

	client, err := github.NewClient(github.WithAuthToken(tok.AccessToken))
	if err != nil {
		return GithubProfile{}, err
	}

	authedUser, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return GithubProfile{}, fmt.Errorf("fetching authenticated user: %w", err)
	}

	if username != "" && !strings.EqualFold(authedUser.GetLogin(), username) {
		return GithubProfile{}, fmt.Errorf("authenticated as %q but requested profile for %q — private email is only available for your own account", authedUser.GetLogin(), username)
	}

	return GithubProfile{
		Name:  authedUser.GetName(),
		Login: authedUser.GetLogin(),
		Email: authedUser.GetEmail(),
	}, nil
}
