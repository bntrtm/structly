package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bntrtm/structly/menu"
	tea "github.com/charmbracelet/bubbletea"
)

// Welcome to Structly!
// This is an introductory example you can read through to get a
// feel for the basics. There's more you can read about that covers
// the way struct tags are leveraged for memory performance and
// field blacklisting, but this is a good place to get started!

// STEP 1: Declare the struct whose fields you wish to expose to users for input.

// applicationForm holds fields typical of a job application.
type applicationForm struct {
	FirstName string `smname:"First Name"`
	LastName  string `smname:"Last Name"`
	Email     string
	PhoneNo   int    `smname:"Phone"`
	Country   string `smname:"Country"`
	Location  string `smname:"Location (City)"`
	CanTravel bool   `smname:"Travel" smdes:"Can you travel for work?"`
}

func main() {
	// STEP 2: Choose custom options to apply to your menu, if you desire.
	// Ensure that if you use custom options, you use NewMenuOptions(),
	// so that their values are initialized to sane defaults before use!
	customMenuOptions := menu.NewMenuOptions()
	customMenuOptions.SetHeader("Apply for this job: ")

	// STEP 3: Define a struct to use.
	// Don't worry, if you need to provide a struct with non-zero
	// values, you can also do that! The Structly bubbletea model will
	// respect those values as defaults all the same.
	newApplication := applicationForm{}

	// STEP 4: Initialize a menu!
	// Provide a pointer to your struct.
	// If you chose not to define custom options, use the NewMenu function, instead.
	model, err := menu.NewMenuWithOptions(&newApplication, customMenuOptions)
	if err != nil {
		log.Fatalf("Trouble generating the application: %s", err)
	}
	// STEP 5: Pass your new Structly model to the bubbletea package to generate your menu!
	// After the bubbletea program exits, the instance you defined will contain the
	// user-input values.
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
