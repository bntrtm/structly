# Introduction to Structly

Welcome to the introductory example for using Structly's menu package!
You should find some example code (with robust in-line documentation)
alongside this README file within the repository.

Let's walk through it!

### Step 1: Declare the struct you wish to expose to your CLI user

In this struct, we build out a list of potential fields on a theoretical
job application to illustrate the idea.

**What To Know**:

- The `smname` tag establishes the title and formatting of the field. If the tag is not present, the menu will fall back to the default name of the struct field itself. For example, you'll see in the above demonstration that the `Email` field renders as we would expect despite the lack of the `smname` tag.
- The `smdes` tag renders an optional description when the user hovers their cursor over the field.

```go
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
```

### Step 3: Choose custom options to apply to your menu, if you desire

There are a number of custom options we could apply to our menu, such as
changing the ibeam cursor ('|'') rendered during string input, or what the field
cursor might look like ('> ').

For now, we'll just write a custom header to render during form interaction.
Custom options _must_ be initialized to defaults before setting any values
on them. Thankfully, if we use the provided function `NewMenuOptions()`, this
is already done for us under the hood!

```go
 customMenuOptions := menu.NewMenuOptions()
 customMenuOptions.SetHeader("Apply for this job: ")
```

### Step 4: Define a struct to use during menu input

Of course, you'll need a struct to expose to your CLI users!
Here, we simply define an empty one, but don't worry: if you need to provide a struct
with non-zero values, you can also do that! The Structly model will respect those values as defaults all the same.

```go
 newApplication := applicationForm{}
```

### Step 5: Initialize your model

Provide a pointer to your struct instance, any custom options you've defined, and any applicable exception list.

```go
 model, err := menu.NewMenuWithOptions(&newApplication, customMenuOptions)
 if err != nil {
  log.Fatalf("Trouble generating the application: %s", err)
 }
```

### Step 6: Use the model with the bubbletea package to render your menu

The menu is a bubbletea model! That is, it implements the Bubble Tea `Model` interface.
We're now ready to run it through a new Bubble Tea program to expose the menu to users to capture their input! The result is the demo you saw above.

```go
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
```

You have now captured user input for one or more fields using the `structly` package!
Do what you need with these new values. In the demo seen below, our program prints the name of the applicant after applying.

![structly-example](https://github.com/user-attachments/assets/6effe862-21c1-472a-85ca-5b815580807c)

## What's Next?

Now that you've been introduced Structly, you ought to know about ways to improve the memory efficiency of the model you're generating. It may not be necessary right away, but because Structly uses the `reflect` package, it may behoove you to at least be aware of the tools Structly makes available to you to keep your program running as efficiently as possible.

Otherwise, blacklisting with Structly's 'exceptions' logic is also one of its most useful features!

[Read about the IDX Tag](/examples/theIdxTag)
[Read about Field Exceptions](/examples/exceptions)
