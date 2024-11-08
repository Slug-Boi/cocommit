package utils

import (
	"bufio"
	"os"
	"strings"
)

// This util file is used to handle users and their information
type User struct {
	Username string
	Email    string
	Names    string
}

var Users = map[string]User{}
var DefExclude = []string{}
var Groups = map[string][]User{}

func Define_users(author_file string) {
	file, err := os.Open(author_file)
	if err != nil {
		print("File not found")
		os.Exit(2)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// eat a single input
	scanner.Scan()

	// reads the input of authors file and formats accordingly
	for scanner.Scan() {
		input_str := scanner.Text()
		group_info := []string{}
		if strings.Contains(input_str, ";;") {
			input := strings.Split(input_str, ";;")
			input_str = input[0]
			group_info = append(group_info, strings.Split(input[1], "|")...)
		}
		info := strings.Split(input_str, "|")
		if len(info) < 4 {
			if len(info) > 0 {
				println("Error: User ", info[0], " is missing information")
			} else {
				println("Error: Some user is missing information")
			}
			println("Please check the author file for proper syntax")
			os.Exit(1)
		}
		usr := User{Username: info[2], Email: info[3], Names: info[0] + "/" + info[1]}
		Users[info[0]] = usr
		Users[info[1]] = usr
		// Adds users with the ex tag to the defExclude list
		if len(info) == 5 {
			if info[4] == "ex" {
				DefExclude = append(DefExclude, info[2])
			}
		}
		if len(group_info) > 0 {
			// Group assignment
			for _, group := range group_info {
				if Groups[group] == nil {
					Groups[group] = []User{usr}
				} else {
					//TODO: Try and find a cleaner way of doing this
					usr_lst := Groups[group]
					usr_lst = append(usr_lst, usr)
					Groups[group] = usr_lst
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		os.Exit(2)
	}
}

func RemoveUser(short string) {
	usr := Users[short]
	split := strings.Split(usr.Names, "/")
	delete(Users, split[0])
	delete(Users, split[1])
}

func TempAddUser(username, email string) {
	usr := User{Username: username, Email: email}

	Users[username] = usr
}
