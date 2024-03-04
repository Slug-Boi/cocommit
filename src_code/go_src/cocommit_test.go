package main

import (
	"os"
	"os/exec"
	"testing"
)

func Test_emptyInput(t *testing.T) {
	authors := make(map[string]user)
	authors["test"] = user{username: "test", email: "test"} 
	if os.Getenv("BE_CRASHER") == "1" {
        NoInput([]string{}, authors)
        return
    }
    cmd := exec.Command(os.Args[0], "-test.run=Test_emptyInput")
    cmd.Env = append(os.Environ(), "BE_CRASHER=1")
    err := cmd.Run()
    if e, ok := err.(*exec.ExitError); ok && !e.Success() {
        return
    }
    t.Fatalf("process ran with err %v, want exit status 1", err)
}

func Test_usersInput(t *testing.T) {
	authors := make(map[string]user)
	authors["test"] = user{username: "test", email: "test"} 
	if os.Getenv("BE_CRASHER") == "1" {
        NoInput([]string{"users"}, authors)
        return
    }
    cmd := exec.Command(os.Args[0], "-test.run=Test_usersInput")
    cmd.Env = append(os.Environ(), "BE_CRASHER=1")
    err := cmd.Run()
    if e, ok := err.(*exec.ExitError); ok && !e.Success() {
        return
    }
    t.Fatalf("process ran with err %v, want exit status 1", err)
}
//TODO: Turn this into a fuzz test

func Test_commit_message(t *testing.T) {
    //authors := make(map[string]user)
	users["test"] = user{username: "test", email: "test"} 
    sb_author("test")
    commit := sb_build()
    if commit != "\nCo-authored-by: test <test>" {
        t.Fatalf("String built incorrectly. Strings did not match: Created -> %s Expected -> Co-authored-by: test <test>",commit)
    }
}
//TODO: Turn this into a fuzz test
func Test_add_all(t *testing.T) {
    for k := range users {
        delete(users, k)
    }
    users["test1"] = user{username: "test1", email: "test1"}
    users["test2"] = user{username: "test2", email: "test2"}  
    users["test3"] = user{username: "test3", email: "test3"} 
    all_flag = true
    add_x_users([]string{})

    commit := sb_build()

    if commit != "\nCo-authored-by: test <test>\nCo-authored-by: test1 <test1>\nCo-authored-by: test2 <test2>\nCo-authored-by: test3 <test3>" {
        t.Fatalf("String built incorrectly. Strings did not match: Created -> %s Expected -> Co-authored-by: test <test>",commit)
    }
}  






