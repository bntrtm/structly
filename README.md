# gostructui

![gostructui-example](https://github.com/user-attachments/assets/6effe862-21c1-472a-85ca-5b815580807c)

A Go Library of bubbletea models leveraging the `reflect` package to expose structs
as forms and menus directly to CLI users, allowing them to edit fields with primitive types.

## Motivation

I built `gostructui` just as soon as I realized that I needed an easy, all-in-one method to ask for
user form input through a CLI. We shouldn't always have to ask CLI users for _one thing at a time._
I personally like to expose config structs to CLI users so that they can set those values easily
through the CLI, and then I save the result.

## Usage

Right now, the only user-editable fields are:

- Strings
- Integers
- Booleans

The repo contains an example of how to use the package within `./example/default/main.go`. Let's walk through it!

### Step 1: Import the menu package

```bash
go get github.com/bntrtm/gostructui
```

```go
import "github.com/bntrtm/gostructui/menu"
```

### Step 2: Establish the struct you wish to expose to the user

In this struct, we build out a list of potential fields on a theoretical
job application to illustrate the idea.

**What To Know**:

- The `smname` tag establishes the title and formatting of the field. If the tag is not present,
  the menu will fall back to the default name of the struct field itself. For example, you'll
  see in the above demonstration that the `Email` field renders as we would expect despite the
  lack of the `smname` tag.
- The `smdes` tag renders an optional description when the user hovers their cursor over the field.
- We'll discuss the `BlacklistedField` bit in a minute. It will illustrate another feature!

```go
// applicationForm holds fields typical of a job application.
type applicationForm struct {
 FirstName        string `smname:"First Name"`
 LastName         string `smname:"Last Name"`
 Email            string
 PhoneNo          int    `smname:"Phone"`
 Country          string `smname:"Country"`
 Location         string `smname:"Location (City)"`
 CanTravel        bool   `smname:"Travel" smdes:"Can you travel for work?"`
 BlacklistedField string
}
```

### Step 3: Choose custom options to apply to your menu, if you desire

There are a number of custom options we could apply to our menu, such as
changing the ibeam cursor rendered during string input, or what the field
cursor might look like.

Here, we write a custom header to render during form interaction.
Because we're using custom options, we will have to initialize them
before setting any of the values on them.
_Never forget to do this! Zero values for menu options are NOT the defaults._

```go
 customMenuOptions := &menu.MenuOptions{}
 customMenuOptions.Init()
 customMenuOptions.Header = "Apply for this job: "
```

### Step 4: Provide a struct to use during menu input

Of course, you'll need a struct to expose to your CLI users!
Here, we simply declare an empty one, but don't worry: if you need to provide a struct
with non-zero values, you can also do that! The bubbletea model will keep those values
intact, showing them to users as existing values within the field.

```go
 newApplication := applicationForm{}
```

### Step 5: Initialize a menu

Provide a pointer to your struct, a list of fields used as a whitelist or blacklist, and any
custom options.

Hey, there's our `BlacklistedField` option we set earlier! See how our
argument passed to the `asBlacklist` parameter is set to `true`? It means that any fields
with the names given within the string slice to the left will be hidden from users. You can
see it in the demo above; the field doesn't show up!

```go
configEditMenu, err := gostructui.InitialTModelStructMenu(&newApplication, []string{"BlacklistedField"}, true, customMenuOptions)
 if err != nil {
  log.Fatal("Trouble generating the application.")
 }
```

### Step 6: Use the menu with the bubbletea package

The menu is a bubbletea model! That is, it implements the bubbletea package!
We're now ready to run it through bubbletea and expose the menu to users to capture
their input! The result is the demo you saw above.

```go
p := tea.NewProgram(configEditMenu)
 if entry, err := p.Run(); err != nil {
  log.Fatal("Trouble generating the application.")
 } else {
  if entry.(gostructui.TModelStructMenu).QuitWithCancel {
   fmt.Printf("Canceled application.\n")
   os.Exit(0)
  } else {
   err = entry.(gostructui.TModelStructMenu).ParseStruct(&newApplication)
   if err != nil {
    log.Fatal("Trouble generating the application.")
   }

   // newApplication: "Wow, I feel like a new struct!"
  }
  if newApplication.FirstName == "" {
   log.Fatal("ERROR: Missing First Name field!")
  }
  fmt.Printf("Thank you for applying, %s!\n", newApplication.FirstName)
  time.Sleep(time.Second * 5)
  os.Exit(0)
 }
```

You have now captured user input for one or more fields using the `gostructui` package!
Do what you need with these new values. In the demo, our program
prints the name of the applicant after applying.

## Advanced Features

### The `idx` Tag

Great! We can declare the shape of user input using a struct. But what if I'm a
sucker for memory performance? If we must declare the fields of an input struct
in the order by which we want them displayed to users, we may face a tradeoff in
the form of a suboptimal memory layout on that struct.

See, for example, our `applicationForm` struct from earlier. Imagine if we
wanted to display another bool-type option, `ConsentToSMS`, after the `PhoneNo`
field, but before the `Country` field. Because struct fields in Go sit in
memory within a contiguous block, we end up with needless padding in two places,
generated by Go for the sake of making each field easily addressable by the CPU:

```go
// assume we're on a 64-bit system for this example
type applicationForm struct {
  // ...
  PhoneNo          int
  ConsentToSMS     bool   // 1 byte
  // [ PADDING ]          // 7 bytes
  Country          string
  // ...
  CanTravel        bool   // 1 byte
  // [ PADDING ]          // 7 bytes
  BlacklistedField string
}
```

Go is adding 7 bytes of padding after each boolean-type struct field
by design, but we know that were we to pair those boolean fields together
(preferably at the end of the struct), our memory layout would be far more
optimal: Go would be adding 6 bytes of padding, rather than 14! Yet, because
of our selfish desire for a more sensible user form that asks for a user's
consent to receive SMS text messages _after_ inquiring about their phone number,
we end up wasting a total of 8 precious bytes!

Is there anyone who could help us?

Enter the `idx` tag! This tag allows us to declare our struct fields in whatever
more memory-performant way we desire, while telling `gostructui`'s bubbletea
model what order we actually want to display them in.

```go
// applicationForm holds fields typical of a job application.
type applicationForm struct {
 FirstName        string `idx:"0"`
 LastName         string `idx:"1"`
 Email            string `idx:"2"`
 Country          string `idx:"5"`
 Location         string `idx:"6"`
 BlacklistedField string `idx:"8"`
 PhoneNo          int    `idx:"3"`
 ConsentToSMS     bool   `idx:"4"`
 CanTravel        bool   `idx:"7"`
}
```

Sure, the order of display may become a bit less readable at a glance by us
developers within the source code, but we'll have to accept _some_ tradeoff
in the end, right? It's just nice to be able to choose which tradeoff we're
willing to accept.

The rules to using the `idx` tag are strict, but simple:

- **All or Nothing**: if you use it on one field, use it on all fields.
- **Start at 0**: the values must start at 0
- **Keep sequence**: the values must not break sequence.
  - No values may be skipped.
  - No two values may be the same.

To observe the `idx` tag in action, run the example at `./example/withIDX/main.go`!
You'll find that it has the same behavior as seen in the `default` example,
despite the reordered fields under the `applicationForm` struct.

This is the power of the `idx` tag!

You can also see a demonstration of just the behavior itself by running the
subtest under `idx_test.go` dedicated to that using the following command in
your terminal:

```bash
go test -run TestIDXMemoryLayout -v
```
