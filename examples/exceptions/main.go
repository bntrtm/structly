package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bntrtm/structly/menu"
	tea "github.com/charmbracelet/bubbletea"
)

// This example is provided to demonstrate the use of 'exceptions.'
// In the context of Structly, this term refers to how we can tell
// Structly which struct fields we don't want exposed to users.
//
// If you've not yet read through the 'introduction' example, it is
// best to start there. You can compare it with this example to
// get an understanding of where exception logic comes into play.

// STEP 1: Declare the struct whose fields you wish to expose to users for input.

// applicationForm holds fields typical of a job application.
type applicationForm struct {
	FirstName        string `smname:"First Name"`
	LastName         string `smname:"Last Name"`
	Email            string
	PhoneNo          int    `smname:"Phone"`
	Country          string `smname:"Country"`
	Location         string `smname:"Location (City)"`
	CanTravel        bool   `smname:"Travel" smdes:"Can you travel for work?"`
	BlacklistedField string `bl:""` //
	BlacklistMe      int
}

func main() {
	// Declare options here, if you please...

	// STEP 2: Define your struct...
	newApplication := applicationForm{}
	// STEP 3: Initialize a menu!
	//
	// This is the step wherein you may specify exception lists for your model.
	//
	// NOTE: Black() and White() exist as convenience wrappers to satisfy
	// validation logic for exceptions under the hood.
	model, err := menu.NewMenu(&newApplication, menu.Black("BlacklistMe")...)
	if err != nil {
		log.Fatalf("Trouble generating the application: %s", err)
	}
	// STEP 4: Run the bubbletea program.
	// Blacklisted fields will not show up!
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		log.Fatalf("Trouble generating the application: %s", err)
	} else {
		if model.EndState.QuitWithCancel {
			fmt.Printf("Canceled application.\n")
			os.Exit(0)
		} else {
			err = model.ParseStruct(&newApplication)
			if err != nil {
				log.Fatalf("Trouble getting data from the application: %s", err)
			}

			// Your struct is now full of user-input values!
			// Do what you need with it.

			// newApplication: "Wow, I feel like a new struct!"
		}
		if newApplication.FirstName == "" {
			log.Fatal("ERROR: Missing First Name field!")
		}
		fmt.Printf("Thank you for applying, %s!\n", newApplication.FirstName)
		time.Sleep(time.Second * 3)
		os.Exit(0)
	}
}
