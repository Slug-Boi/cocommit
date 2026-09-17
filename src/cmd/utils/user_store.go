package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/google/go-github/v91/github"
	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

const (
	keyringService = "cocommit-github"
	keyringUser    = "token"
)

type IndexEntry struct {
	UUID     string `json:"uuid"`
	Login    string `json:"login"`
	Username string `json:"username"`
	Platform string `json:"platform"`
	Path     string `json:"path"`
}

func SaveToken(tok *oauth2.Token) error {
	data, err := json.Marshal(tok)
	if err != nil {
		return fmt.Errorf("marshaling token: %w", err)
	}
	return keyring.Set(keyringService, keyringUser, string(data))
}

func LoadToken() (*oauth2.Token, error) {
	data, err := keyring.Get(keyringService, keyringUser)
	if err != nil {
		if err == keyring.ErrNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("reading token from keyring: %w", err)
	}

	var tok oauth2.Token
	if err := json.Unmarshal([]byte(data), &tok); err != nil {
		return nil, fmt.Errorf("unmarshaling token: %w", err)
	}
	return &tok, nil
}

func SubmitAuthor(ctx context.Context, userToken, owner, repo, uuid, serializedBlob string) (*github.Issue, error) {
	client, err := github.NewClient(github.WithAuthToken(userToken))
	if err != nil {
		return nil, err
	}

	body := fmt.Sprintf("uuid: %s\n\n```\n%s\n```", uuid, serializedBlob)

	issue, _, err := client.Issues.Create(ctx, owner, repo, github.CreateIssueRequest{
		Title:  fmt.Sprintf("%s", "publish-author: "+uuid),
		Body:   new(body),
		Labels: []string{"author-submission"},
	})
	return issue, err
}

var githubOAuthConfig = &oauth2.Config{
	ClientID: "Ov23lieXDKh7KvTnEjGr",
	Endpoint: oauth2.Endpoint{
		AuthURL:       "https://github.com/login/oauth/authorize",
		DeviceAuthURL: "https://github.com/login/device/code",
		TokenURL:      "https://github.com/login/oauth/access_token",
	},
	Scopes: []string{"public_repo"}, // or "repo" if you need private repo access too
}

func Login(ctx context.Context) (*oauth2.Token, error) {
	// Step 1: request a device code
	resp, err := githubOAuthConfig.DeviceAuth(ctx)
	if err != nil {
		return nil, fmt.Errorf("requesting device code: %w", err)
	}

	// Step 2: show the user what to do
	fmt.Printf("First, copy your one-time code: %s\n", resp.UserCode)
	fmt.Printf("Then open: %s\n", resp.VerificationURI)
	// resp.VerificationURIComplete (if present) is a URL that pre-fills the code —
	// nice to print too, or to open automatically in a browser.

	// Step 3: poll until the user approves (this blocks until success/expiry/denial)
	token, err := githubOAuthConfig.DeviceAccessToken(ctx, resp)
	if err != nil {
		return nil, fmt.Errorf("waiting for authorization: %w", err)
	}

	return token, nil
}

func DeleteToken() error {
	return keyring.Delete(keyringService, keyringUser)
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

func indexURL(userStoreOwner, userStoreRepo, userStoreRef string) string {
	return fmt.Sprintf("https://cdn.jsdelivr.net/gh/%s/%s@%s/index.json", userStoreOwner, userStoreRepo, userStoreRef)
}

func blobURL(path, userStoreOwner, userStoreRepo, userStoreRef string) string {
	return fmt.Sprintf("https://cdn.jsdelivr.net/gh/%s/%s@%s/%s", userStoreOwner, userStoreRepo, userStoreRef, path)
}

func fetchIndex(userStoreOwner, userStoreRepo, userStoreRef string) ([]IndexEntry, error) {
	resp, err := httpClient.Get(indexURL(userStoreOwner, userStoreRepo, userStoreRef))
	if err != nil {
		return nil, fmt.Errorf("fetching index: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetching index: unexpected status %d", resp.StatusCode)
	}

	var index []IndexEntry
	if err := json.NewDecoder(resp.Body).Decode(&index); err != nil {
		return nil, fmt.Errorf("parsing index: %w", err)
	}
	return index, nil
}

func fetchBlob(path, userStoreOwner, userStoreRepo, userStoreRef string) (string, error) {

	resp, err := httpClient.Get(blobURL(path, userStoreOwner, userStoreRepo, userStoreRef))
	if err != nil {
		return "", fmt.Errorf("fetching blob: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("fetching blob: unexpected status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading blob: %w", err)
	}
	return string(data), nil
}

// FetchUserStoreUser looks up a username in the public author store and
// returns the matching decoded User(s). platform is optional — pass ""
// to match on username alone (may return multiple entries, since
// usernames aren't unique across platforms); pass a platform to narrow
// to a single match.
func FetchUserStoreUser(username, platform, userStoreOwner, userStoreRepo, userStoreRef string) string {
	index, err := fetchIndex(userStoreOwner, userStoreRepo, userStoreRef)
	if err != nil {
		fmt.Println("Error fetching user store index:", err)
		return ""
	}

	var matches []IndexEntry
	for _, e := range index {
		if e.Username != username {
			continue
		}
		if platform != "" && e.Platform != platform {
			continue
		}
		matches = append(matches, e)
	}

	if len(matches) == 0 {
		fmt.Printf("No user found for username %q\n", username)
		return ""
	}

	var result string
	for _, m := range matches {
		encoded, err := fetchBlob(m.Path, userStoreOwner, userStoreRepo, userStoreRef)
		if err != nil {
			fmt.Printf("Error fetching data for %s (%s): %v\n", m.Username, m.UUID, err)
			continue
		}

		result = ImportUsersFromShareCode(encoded)
	}

	if len(result) == 0 {
		return ""
	}

	return result
}

// Matches:
//   owner/repo
//   owner/repo@ref
//   https://github.com/owner/repo
//   https://github.com/owner/repo.git
//   https://github.com/owner/repo@ref
//   git@github.com:owner/repo.git
var storeRepoRe = regexp.MustCompile(
	`^(?:https?://github\.com/|git@github\.com:)?` + // optional URL/SSH prefix
		`(?P<owner>[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,37}[a-zA-Z0-9])?)` + // GitHub username rules
		`/` +
		`(?P<repo>[a-zA-Z0-9._-]+?)` + // repo name (non-greedy, so .git/@ref strip cleanly)
		`(?:\.git)?` + // optional .git suffix
		`(?:@(?P<ref>[a-zA-Z0-9._/-]+))?` + // optional @ref
		`/?$`,
)

type StoreRepo struct {
	Owner string
	Repo  string
	Ref   string
}

const defaultRef = "master"

func ParseStoreRepo(input string) (StoreRepo, error) {
	m := storeRepoRe.FindStringSubmatch(input)
	if m == nil {
		return StoreRepo{}, fmt.Errorf("invalid repo format %q: expected owner/repo, owner/repo@ref, or a GitHub URL", input)
	}

	result := StoreRepo{Ref: defaultRef}
	for i, name := range storeRepoRe.SubexpNames() {
		switch name {
		case "owner":
			result.Owner = m[i]
		case "repo":
			result.Repo = m[i]
		case "ref":
			if m[i] != "" {
				result.Ref = m[i]
			}
		}
	}

	return result, nil
}


const defaultPlatform = "github"

type UserQuery struct {
	Username string
	Platform string
}

// ParseUserQueries splits a space-separated list of "username" or
// "username@platform" tokens. Platform defaults to "github" when omitted.
func ParseUserQueries(input string) ([]UserQuery, error) {
	fields := strings.Fields(input) // splits on any whitespace, collapses repeats — safer than strings.Split(input, " ")

	if len(fields) == 0 {
		return nil, fmt.Errorf("no usernames provided")
	}

	queries := make([]UserQuery, 0, len(fields))
	for _, f := range fields {
		username, platform, found := strings.Cut(f, "@")
		if username == "" {
			return nil, fmt.Errorf("invalid entry %q: missing username", f)
		}

		if !found {
			platform = defaultPlatform
		} else if platform == "" {
			return nil, fmt.Errorf("invalid entry %q: missing platform after '@'", f)
		} else if strings.Contains(platform, "@") {
			return nil, fmt.Errorf("invalid entry %q: too many '@' separators", f)
		}
		platform = strings.ToLower(platform)

		queries = append(queries, UserQuery{Username: username, Platform: platform})
	}

	return queries, nil
}
