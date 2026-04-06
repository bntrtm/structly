# Field Exceptions

## Introduction

### About

In Structly, the term `exception` refers to a struct field marked as
blackisted or whitelisted from render. That is, when we have a struct
whose fields we wish to expose to a user, _except_ for one or more of
those fields (as they may be reserved only for business logic to handle),
we can blacklist them!

There are two methods available for blacklisting fields, and they each
exist to serve different purposes:

1. The `bl` Struct Tag
2. Exception Lists

## The `bl` Struct Tag

Using the `bl` struct tag on a single field in a struct will blacklist
it from all rendering. It is a means of indicating an exception at the
type level.

Using `bl` is often the most simple option for indicating exceptions,
but it does mean that we can _never_ expose that field to users for input.
It is simply impossible by design. The tag is most useful in cases where
one or more struct fields are declared solely for the purpose of being
coupled with user inputs for business logic to be run later.

When using the `bl` tag, its value never matters, so simply leave it empty
as a best practice.

```go
type s struct {
  s string `bl:""` // the 's' field will never render for user input
  i int
  b bool
}
```

## Exception Lists

Exception lists allow us to tell the Structly model that we don't want
one or more specified fields (as indicated to by name) to be rendered
_for a single instance of a model render._ For reference, see the following
struct and note that the functions provided for menu generation offer an
`exceptionList` parameter for passing variadic arguments:

```go
type s struct {
  s string
  i int
  b bool
}

func NewMenu(i any, exceptionList ...string) (Model, error) {...}

func NewMenuWithOptions(structlyPtr any, options *MenuOptions, exceptionList ...string) (Model, error) {}

```

Perhaps in some cases, we want to render all three of these fields
in a CLI, but in others, we don't. Rather than using the `bl` tag
as a permanent solution to the problem, we can simply pass an exception
list, appended with an indicator for blacklisting, to any of the
provided functions for menu generation:

```go
menu.NewMenu(&myStruct, "s", "i", menu.BlacklistIndicator)
```

The `BlacklistIndicator` value is just a string with the value `BL` under
the hood. This also works, but is not recommended:

```go
menu.NewMenu(&myStruct, "s", "i", "BL")
```

The indicator, which could just as well be the `WhitelistIndicator` (`WL`),
tells the Structly model which mode of exception should be used when
evaluating the list. This method does lose a bit of logical flow when reading
such function calls, however; it may be easy to assume that "BL" is the name
of a struct field that we want excepted, as opposed to what it is really there
for. With this in mind, Structly provides some convenience wrappers for providing
exception lists with this trouble abstracted away:

```go
model, err := menu.NewMenu(&myStruct, Black("s", "i")...)
```

```go
model, err := menu.NewMenu(&myStruct, White("b")...) // whitelisting 'b' yields same result as above
```

## Further Notes

### Non-Exclusivity

Both of the aforementioned methods work in tandem. Were you to the following:

```go
type s struct {
  s string `bl:""`
  i int
  b bool
}

menu.NewMenu(&myStruct, Black("i")...)
```

You would get a model whose render would consist only of the struct field `b`
exposed for user input.

### Tag Interoperability

The `bl` tag is **INCOMPATIBLE** with the `idx` tag for a single field. That is,
the tags will respect each other's presence so long as they never appear
together within the same full struct tag for a field.

```go
// This is VALID; Structly will not error out.
// The model will expose b before i, and s will not be exposed.
type s struct {
  s string  `bl:""`
  i int     `idx:"1"`
  b bool    `idx:"0"`
}

model, err := menu.NewMenu(&myStruct)

// This is INVALID; Structly will error out
type s struct {
  s string  `bl:"" idx:"2"`
  i int     `idx:"1"`
  b bool    `idx:"0"`
}
```

### All Exception Values Respected

You were made aware earlier that Structly evaluates the last element of
an exception list as the indicator for whether to define the exceptions
by whitelisting or blacklisting.

In case you wondered, you need never worry over whether any one of your
exceptions matches an indicator value. In the following examples, the struct
field `BL` _is_ properly rendered:

```go
type s struct {
  BL string
  s string
  i int
  b bool
}

exampleOne, err := menu.NewMenu(&myStruct, "s", "i", "BL")

exampleTwo, err := menu.NewMenu(&myStruct, Black("s", "i")...)
```
