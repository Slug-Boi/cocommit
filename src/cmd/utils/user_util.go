package utils

import (
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"

	"encoding/base64"
	"encoding/json"
)

// This util file is used to handle users and their information

type User struct {
	Shortname string   `json:"shortname"`
	Longname  string   `json:"longname"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Ex        bool     `json:"ex"`
	Groups    []string `json:"groups"`
	From_git  bool     `json:"from_git,omitempty"`
	Platform  string   `json:"platform,omitempty"`
	uuid      string
}

type Author struct {
	Authors map[string]User
}

// purely used for editing the author file later
var Authors = Author{}

var Users = map[string]User{}
var DefExclude = []string{}
var Groups = map[string][]User{}

var Git_Users = map[string]User{}
var Git_Groups = map[string][]User{}

func userEquals(a, b User) bool {
	return a.Shortname == b.Shortname &&
		a.Longname == b.Longname &&
		a.Username == b.Username &&
		a.Email == b.Email &&
		a.Ex == b.Ex &&
		a.Platform == b.Platform &&
		slices.Equal(a.Groups, b.Groups)
}

func LookupAuthor(key string) (User, bool) {
	if user, ok := Authors.Authors[key]; ok {
		user.uuid = key
		return user, true
	}

	if user, ok := Users[key]; ok {
		return user, true
	}

	return User{}, false
}

func LookupAuthorID(user User) (string, bool) {
	for id, author := range Authors.Authors {
		if userEquals(author, user) {
			return id, true
		}
	}

	return "", false
}

func authorTokenMatches(user User, token string) bool {
	token = strings.TrimSpace(token)
	return strings.EqualFold(token, user.Shortname) ||
		strings.EqualFold(token, user.Longname) ||
		strings.EqualFold(token, user.Username) ||
		strings.EqualFold(token, user.uuid)
}

func ResolveAuthorToken(token string) (User, bool) {
	if user, ok := Authors.Authors[token]; ok {
		user.uuid = token
		return user, true
	}

	ids := make([]string, 0, len(Authors.Authors))
	for id := range Authors.Authors {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	var match User
	found := false
	for _, id := range ids {
		user := Authors.Authors[id]
		user.uuid = id
		if !authorTokenMatches(user, token) {
			continue
		}
		if found {
			return match, true
		}
		match = user
		found = true
	}

	return match, found
}

func ResolveGroupToken(token string) ([]User, bool) {
	for groupName, users := range Groups {
		if strings.EqualFold(groupName, token) {
			return users, true
		}
	}

	return nil, false
}

func ContainsUser(users []User, user User) bool {
	return slices.ContainsFunc(users, func(u User) bool {
		return u.Shortname == user.Shortname &&
			u.Longname == user.Longname &&
			u.Username == user.Username &&
			u.Email == user.Email &&
			u.Ex == user.Ex &&
			u.Platform == user.Platform &&
			slices.Equal(u.Groups, user.Groups)
	})
}

func CheckUserFields(user User) bool {
	if user.Shortname == "" || user.Longname == "" || user.Username == "" || user.Email == "" || user.Platform == "" {
		return false
	}
	return true
}

func Define_users(author_file string) {
	// wipe the users map
	Users = map[string]User{}
	DefExclude = []string{}
	Groups = map[string][]User{}

	var auth Author

	data, err := os.ReadFile(author_file)
	if err != nil {
		panic(fmt.Sprintf("Error reading author file: %v", err))
	}
	err = json.Unmarshal(data, &auth)
	if err != nil {
		panic(fmt.Sprintf("Error unmarshalling json: %v", err))
	}

	Authors = auth

	for s, usr := range auth.Authors {
		usr.uuid = s
		Users[usr.Shortname] = usr
		Users[usr.Longname] = usr
		if usr.Ex {
			DefExclude = append(DefExclude, s)
		}

		group_info := usr.Groups
		if len(group_info) > 0 {
			for _, group := range group_info {
				if Groups[group] == nil {
					Groups[group] = []User{usr}
				} else {
					usr_lst := Groups[group]
					usr_lst = append(usr_lst, usr)
					Groups[group] = usr_lst
				}
			}
		}
	}
}

func Define_git_users() {
	// wipe the git users map
	Git_Users = map[string]User{}
	Git_Groups = map[string][]User{}

	// get all authors from git
	git_authors := GitCheckAuthors()

	for _, usr := range git_authors {
		if _, ok := LookupAuthorID(usr); !ok {
			usr.From_git = true
			Git_Users[usr.Shortname] = usr
			Git_Users[usr.Longname] = usr

			group_info := usr.Groups
			if len(group_info) > 0 {
				for _, group := range group_info {
					if Git_Groups[group] == nil {
						Git_Groups[group] = []User{usr}
					} else {
						usr_lst := Git_Groups[group]
						usr_lst = append(usr_lst, usr)
						Git_Groups[group] = usr_lst
					}
				}
			}
		}
	}
}

func RemoveUser(short string) {
	usr, ok := LookupAuthor(short)
	if !ok {
		return
	}
	delete(Users, usr.Shortname)
	delete(Users, usr.Longname)
}

func TempAddUser(username, email string) {
	usr := User{Username: username, Email: email}

	Users[username] = usr
}

func SerealizeUsers(authors []string) string {
	var users []User
	for _, name := range authors {
		if usr, ok := LookupAuthor(name); ok {
			users = append(users, usr)
		}
	}

	bytes, err := json.Marshal(users)
	if err != nil {
		panic(err)
	}

	encoded := base64.StdEncoding.EncodeToString(bytes)

	return encoded
}

func UnserealizeUsers(encoded string) ([]string, []string) {
	users := []User{}

	raw, _ := base64.StdEncoding.DecodeString(encoded)
	json.Unmarshal(raw, &users)

	added_users, not_added := CreateMultipleAuthors(users)

	return added_users, not_added
}

func ImportUsersFromShareCode(args []string) string {
	var sb strings.Builder

	if len(args) > 0 {
		added_users, not_added := UnserealizeUsers(args[0])

		if len(added_users) == 0 {
			fmt.Println("\033[33mNo authors added (authors probably already existed or corrupted \"share code\")\033[0m")
			sb.WriteString("\033[33mNo authors added (authors probably already existed or corrupted \"share code\")\033[0m\n")
		}

		if len(added_users) != 0 {
			fmt.Println("\033[32mAuthors added:\033[0m")
			sb.WriteString("\033[32mAuthors added:\033[0m\n")

			for _, usr := range added_users {
				fmt.Println("\033[32m+\033[0m ", usr)
				sb.WriteString("\033[32m+\033[0m ")
				sb.WriteString(usr)
			}

			sb.WriteString("\n")
		}

		if len(not_added) != 0 {
			fmt.Println("\033[33mAlready existing authors (not added):\033[0m")
			sb.WriteString("\033[33mAlready existing authors (not added):\033[0m\n")

			for _, usr := range not_added {
				fmt.Println("\033[33m~\033[0m ", usr)
				sb.WriteString("\033[33m~\033[0m ")
				sb.WriteString(usr)
			}

			sb.WriteString("\n")
		}

	} else if len(args) == 0 {
		fmt.Println("\033[33mNo \"share code\", please run the flag with a valid \"share code\"\033[0m")
		sb.WriteString("\033[33mNo \"share code\", please run the flag with a valid \"share code\"\033[0m")
		os.Exit(0)
	}

	return sb.String()
}

func CLIAuthorInput(authors []string) []string {
	var selected []string
	excludeMode := []string{}

	// write the commit message to the string builder
	fst := authors[0]

	if strings.EqualFold(fst, "all") {
		selected = add_x_users_string_slice(excludeMode, selected)
		return selected
	} else if group, ok := ResolveGroupToken(fst); ok {
		excludeMode = group_selection(group, excludeMode)
		selected = add_x_users_string_slice(excludeMode, selected)
		return selected
	}

	for _, committer := range authors {
		if user, ok := ResolveAuthorToken(committer); ok {
			if id, ok := LookupAuthorID(user); ok {
				selected = append(selected, id)
			}
		} else if committer[0] == '^' { // Negations
			if usr, ok := ResolveAuthorToken(committer[1:]); ok {
				if id, ok := LookupAuthorID(usr); ok {
					excludeMode = append(excludeMode, id)
				}
			}
		} else {
			println(committer, "was unknown. User either not defined or name typed wrong")
		}
	}

	if len(excludeMode) > 0 {
		selected = add_x_users_string_slice(excludeMode, selected)
	}

	return selected
}

func add_x_users_string_slice(excludeMode, selected []string) []string {
	if len(DefExclude) > 0 {
		excludeMode = append(excludeMode, DefExclude...)
	}
	for key, user := range Authors.Authors {
		user.uuid = key
		if !slices.Contains(excludeMode, key) {
			selected = append(selected, key)
			excludeMode = append(excludeMode, key)
		}
	}
	return selected
}
