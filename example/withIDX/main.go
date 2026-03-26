package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bntrtm/gostructui/menu"
	tea "github.com/charmbracelet/bubbletea"
)

// STEP 1: Establish the struct you wish to expose to the user.

// applicationForm holds fields typical of a job application.
type applicationForm struct {
	FirstName        string `smname:"First Name" idx:"0"`
	LastName         string `smname:"Last Name" idx:"1"`
	Email            string `idx:"2"`
	Country          string `smname:"Country" idx:"4"`
	Location         string `smname:"Location (City)" idx:"5"`
	BlacklistedField string `idx:"7"`
	PhoneNo          int    `smname:"Phone" idx:"3"`
	CanTravel        bool   `smname:"Travel" smdes:"Can you travel for work?" idx:"6"`
}

func main() {
	// STEP 2: Choose custom options to apply to your menu, if you desire.
	// Ensure that if you use custom options, you use NewMenuOptions(),
	// so that their values are initialized to sane defaults before use!
	customMenuOptions := menu.NewMenuOptions()
	customMenuOptions.SetHeader("Apply for this job: ")

	// STEP 3: Provide a struct to use.
	// Don't worry, if you need to provide a struct with non-zero
	// values, you can also do that! The tea model will keep those
	// values intact.
	newApplication := applicationForm{}
	// STEP 4: Initialize a menu!
	// Provide a pointer to your struct, blacklisted or
	// whitelisted fields, and any custom options.
	model, err := menu.NewMenu(&newApplication, []string{"BlacklistedField"}, true, customMenuOptions)
	if err != nil {
		log.Fatalf("Trouble generating the application: %s", err)
	}
	// STEP 5: Use the menu---a bubbletea model---with the bubbletea package!
	// Here, we capture the result (our struct with user-entered values)
	// as the tea.Model variable "entry".
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		log.Fatalf("Trouble generating the application: %s", err)
	}
	if model.EndState.QuitWithCancel {
		fmt.Printf("Canceled application.\n")
		os.Exit(0)
	} else {
		err = model.ParseStruct(&newApplication)
		if err != nil {
			log.Fatalf("Trouble getting data from the application: %s", err)
		}

		// Your struct is now full of user-entered values!
		// Do what you need with it.

		// newApplication: "Wow, I feel like a new struct!"
	}
	if newApplication.FirstName == "" {
		log.Fatal("ERROR: Missing First Name field!")
	}
	fmt.Printf("Thank you for applying, %s!\n", newApplication.FirstName)
	time.Sleep(time.Second * 5)
	os.Exit(0)
}
