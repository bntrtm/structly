# Structly

![structly-example](https://github.com/user-attachments/assets/6effe862-21c1-472a-85ca-5b815580807c)

A powerful bubbletea component leveraging Go's `reflect` package to expose
structs as menus for input by CLI users. Enjoy the convenience of managing
a single type containing multiple fields of user input, optionally coupled
with related fields reserved for your business logic.

## Motivation

I first built Structly just as soon as I realized that I needed an efficient means of prompting users for form input through a CLI. Ever since, I have continued to improve its efficiencies and capabilities.

You might be wondering, _"Why not use Charm's `Huh?`?"_

But while [huh](https://github.com/charmbracelet/huh) is great, Structly aims to solve an albeit similar but separate set of problems.

Structly's methodology revolves (unsurprisingly) entirely around structs. This provides the following benefits:

- Fields are declared definitively, as opposed to functionally
- Better memory management potential, by virtue of struct layouts
- Struct tags become a useful tool for setting menu options for each field
- Rendering becomes as simple as providing an instance of the type
- Blacklisting fields from menu exposure means relevant information can be tightly coupled with user inputs
- Boilerplate is drastically reduced
  - Default values for fields are simply pulled from the struct instance provided at render time
  - Declaration of a menu is as familiar and succinct as defining a struct type and working with an instance of that, as opposed to calling a great amount of methods

I personally use Structly to expose CLI configurations to users so that they can modify those values easily through the CLI. Then, I can save the result simply by using `json` tags also set on the struct fields.

## Usage

To get started, import the `menu` package:

```bash
go get github.com/bntrtm/structly
```

```go
import "github.com/bntrtm/structly/menu"
```

### Generating Menus

Robust documentation has been provided under [the examples directory](/examples).
You ought to [start with the introductory example](/examples/introduction/); it will walk you
through the basics, and then point you to more advanced resources via
other examples as you read! That said, here's a quick rundown of the
Structly experience as laid out in the initial example:

1. **_Declare your struct._** You'll craft a struct just like any other; but perhaps with some Structly-exclusive struct tags and options that give you the best it has to offer.
2. **_Define a pointer to your struct._** Structly provides two fucntions you may use to generate a menu; each will expect a pointer to a struct as its first argument.
3. **_Generate your model._** Use one of the aforementioned functions, potentially passing in custom options, or an "Exception List" for blacklisting or whitelisting fields.
4. **_Use your model._** You'll render your menu in the terminal by running the model through a Bubble Tea program! Users will edit the fields you exposed from the struct, and when they're finished, you have their input stored within the same instance you defined in Step 2.
5. **_Do your thing._** You have the input you wanted. Now put it to good use!

> [!NOTE]
>
> As of right now, the only types compatible for user-editable fields are:
>
> - Strings
> - Integers
> - Booleans
>
> Be sure that any fields typed incompatibly are blacklisted using Structly's
> exception logic.

### Structly is a Bubble

At the end of the day, Structly is providing you with a Bubble Tea model.
This makes Structly at least as predictable as any other "bubble."

If you aren't familiar with Charm's Bubble Tea package and their greater
ecosystem, you can see the following resources
for more information:

- [Charm Bubble Tea](charm.land/bubbletea)
- [Charm Bubbles](charm.land/bubbles)
- [About Charm](charm.land/)

## Advanced Features

[Better Memory Management with the 'idx' tag](/examples/theIdxTag)
[Blacklisting fields from render with 'Exceptions'](/examples/exceptions/)
