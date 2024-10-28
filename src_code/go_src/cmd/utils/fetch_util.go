package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func Fetch_AuthorFile(url string) {
	authrFileName := "authors"

	configDir, err := os.UserConfigDir()

	configDir = configDir + "/cocommit"

	fmt.Println("Would you like to place the author file in the default location? (%s) (y/n)", configDir)
	fmt.Println("WARNING: this will override the current author file if one is present in that directory with the name 'authors'")
	var response string
	_, err = fmt.Scanln(&response)
	if err != nil {
		fmt.Println("Error reading response")
	}
	
	var filepath string

	switch response {
	case "y", "Y", "yes", "Yes":
		filepath = configDir + "/" + authrFileName
	case "n", "N", "no", "No":
		fmt.Println("Please enter the path where you would like to save the author file")
		_, err = fmt.Scanln(&response)
		if err != nil {
			fmt.Println("Error reading response")
		}
		filepath = response + "/" + authrFileName 
	}

	os.Remove(filepath)

	// Fetch the author file from the given URL
	// Create the file
	out, err := os.Create(filepath)
	if err != nil  {
		fmt.Println("Error creating file: ", err)
		return
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching file: ", err)
		return
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil  {
		fmt.Println("Error writing file: ", err)
		return
	}

}