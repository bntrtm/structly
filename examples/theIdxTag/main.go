package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bntrtm/structly/menu"
	tea "github.com/charmbracelet/bubbletea"
)

// STEP 1: Declare your struct...

// applicationForm holds fields typical of a job application.
type applicationForm struct {
	FirstName   string `smname:"First Name" idx:"0"`
	LastName    string `smname:"Last Name" idx:"1"`
	Email       string `idx:"2"`
	Country     string `smname:"Country" idx:"4"`
	Location    string `smname:"Location (City)" idx:"5"`
	CoverLetter string `idx:"7" smdes:"Tell us a little bit about you!"`
	PhoneNo     int    `smname:"Phone" idx:"3"`
	CanTravel   bool   `smname:"Travel" smdes:"Can you travel for work?" idx:"6"`
}

func main() {
	// Declare options here, if you please...

	// STEP 2: Define your struct...
	newApplication := applicationForm{}

	// STEP 3: Initialize a menu...
	model, err := menu.NewMenu(&newApplication, menu.Black("BlacklistedField")...)
	if err != nil {
		log.Fatalf("Trouble generating the application: %s", err)
	}

	// STEP 5: Generate your menu!
	// Thanks to the power of the 'idx' tag, our fields will render in the order
	// that we specified in the declaration above.
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
	time.Sleep(time.Second * 3)
	os.Exit(0)
}
