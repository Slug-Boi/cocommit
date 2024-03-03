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




