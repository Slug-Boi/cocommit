package utils

import (
	"fmt"
	"os"
	"slices"

	"encoding/json"
)

// This util file is used to handle users and their information

type User struct {
	Shortname string   `json:"shortname"`
	Longname string   `json:"longname"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Ex        bool     `json:"ex"`
	Groups    []string `json:"groups"`
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

func ContainsUser(users []User, user User) bool {
    return slices.ContainsFunc(users, func(u User) bool {
        return u.Shortname == user.Shortname && 
               u.Longname == user.Longname && 
               u.Username == user.Username && 
               u.Email == user.Email && 
               u.Ex == user.Ex &&
               slices.Equal(u.Groups, user.Groups)
    })
}

func CheckUserFields(user User) bool {
	if user.Shortname == "" || user.Longname == "" || user.Username == "" || user.Email == "" {
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
	
	for _, usr := range auth.Authors {
		Users[usr.Shortname] = usr
		Users[usr.Longname] = usr
		if usr.Ex {
			DefExclude = append(DefExclude, usr.Shortname)
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
		if _, ok := Users[usr.Shortname]; !ok {
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
	usr := Users[short]
	delete(Users, usr.Shortname)
	delete(Users, usr.Longname)
}

func TempAddUser(username, email string) {
	usr := User{Username: username, Email: email}

	Users[username] = usr
}
