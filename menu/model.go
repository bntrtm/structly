package menu

import (
	"fmt"
	"log"
	"reflect"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
)

type state struct {
	cursor         *cursor // reference to cursor for handling navigation
	isEditingValue bool    // tracks state of field editing
}

type EndState struct {
	QuitWithCancel bool // can be used to communicate whether changes ought be saved
}

// Model is a menu that can be used to expose
// primitive struct fields to end users for input,
// as if they were elements of a menu.
type Model struct {
	// MENU STATE
	// fields which can be edited; populated dynamically
	menuFields []menuField
	options    MenuOptions
	state      *state
	EndState   EndState
}

func (m *Model) getFieldAtIndex(i int) *menuField {
	return &m.menuFields[i]
}

func (m *Model) getFieldUnderCursor() *menuField {
	return m.state.cursor.under()
}

// NewMenu attempts to validate the given interface object as
// a type compatible for rendering as a Structly menu, and, if
// successful, generates and returns a menu as a bubbletea model.
//
// Interface 'i' MUST be a pointer to a struct that satisfies
// all requirements of a struct compatible with rendering as a
// Structly menu.
//
// The optional 'exceptions' parameter may be provided one or more
// field names to blacklist from view within the resulting Menu
// instance. If used, the list's final element must match the value
// of either of the indicator constants used to define exception lists.
// The Black() and White() functions exist as convenience wrappers to
// provide this functionally.
func NewMenu(structlyPtr any, exceptions ...string) (Model, error) {
	v, err := validateStructPtr(structlyPtr)
	if err != nil {
		return Model{}, err
	}

	return generateNewMenu(v, nil, exceptions...)
}

// NewMenuWithOptions operates just as NewMenu does, but exposes
// a parameter for passing a list of options. Because a call of
// this function is necessarily deliberate, it will helpfully
// return an error if no options are passed in.
func NewMenuWithOptions(structlyPtr any, options *MenuOptions, list ...string) (Model, error) {
	m := Model{}

	if options == nil {
		return m, fmt.Errorf("new menu requested with options, but no options were provided")
	}

	v, err := validateStructPtr(structlyPtr)
	if err != nil {
		return m, err
	}

	return generateNewMenu(v, options, list...)
}

// generateNewMenu expects a reflect.Value validated as a struct value and
// generates a new menu model from the given parameters. If custom options
// are not provided, the menu will fall back to defaults.
func generateNewMenu(v reflect.Value, options *MenuOptions, exceptions ...string) (Model, error) {
	m := Model{
		menuFields: []menuField{},
		options:    *NewMenuOptions(),
		state: &state{
			cursor:         nil,
			isEditingValue: false,
		},
		EndState: EndState{
			QuitWithCancel: false,
		},
	}

	if options != nil {
		m.options = *options
	}

	t := v.Type()
	fields := getFields(t)
	orderedFields, err := getOrderedFields(fields)
	if err != nil {
		return m, err
	}

	exceptionListIndicator, exceptions, err := validateExceptionList(exceptions)
	if err != nil {
		return m, err
	}

	for i := 0; i < len(orderedFields); i++ {
		j, ok := orderedFields[i]
		if !ok {
			return m, fmt.Errorf("could not resolve struct field to display by declaration index %d", i)
		}
		field := fields[j]

		if len(exceptions) != 0 {
			switch exceptionListIndicator {
			case BlacklistIndicator:
				if slices.Contains(exceptions, field.Name) {
					continue
				}
			case WhitelistIndicator:
				if !(slices.Contains(exceptions, field.Name)) {
					continue
				}
			case "":
				panic("no indicator provided for exception list")
			default:
				panic("found unexpected indicator for exception list: " + exceptionListIndicator)
			}
		}

		fieldVal := v.FieldByName(field.Name)
		if !fieldVal.CanSet() {
			log.Printf("Warning: Field '%s' left unexposed (cannot be set; unexported or not addressable).\n", field.Name)
			continue
		}

		newField := menuField{}
		switch field.Type.Kind() {
		case reflect.String:
			newField.kind = FieldString
			newField.s = fieldVal.String()
		case reflect.Bool:
			newField.kind = FieldBool
			newField.b = fieldVal.Bool()
		case reflect.Int:
			newField.kind = FieldInt
			newField.i = int(fieldVal.Int())
		default:
			return m, fmt.Errorf("could not parse struct")
		}
		newField.name = field.Name
		newField.smName = field.Tag.Get("smname")
		newField.smDes = field.Tag.Get("smdes")
		m.menuFields = append(m.menuFields, newField)
	}

	if len(m.menuFields) == 0 {
		return m, fmt.Errorf("ERROR: No fields to expose to users in struct")
	}

	m.state.cursor = NewCursor(m.menuFields, 0)
	if m.state.cursor == nil {
		return m, fmt.Errorf("ERROR, but: len fields %d", len(m.menuFields))
	}

	return m, nil
}

// validateStructPtr takes in an interface and ensures that
// it is a pointer to a struct type before returning then
// returning the struct as a reflect.Value.
func validateStructPtr(i any) (reflect.Value, error) {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Pointer {
		return v, fmt.Errorf("input interface should be a pointer to a struct, so as to have addressable fields")
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return v, fmt.Errorf("input ptr found not to point to a struct")
	}

	return v, nil
}

func (m Model) ParseStruct(structlyPtr any) error {
	v, err := validateStructPtr(structlyPtr)
	if err != nil {
		return err
	}

	for _, f := range m.menuFields {
		field := v.FieldByName(f.name)

		if !field.IsValid() {
			log.Printf("Warning: Field '%s' not found in struct.\n", f.name)
			continue
		}
		if !field.CanSet() {
			log.Printf("Warning: Field '%s' cannot be set (unexported or not addressable).\n", f.name)
			continue
		}

		switch f.kind {
		case FieldString:
			field.SetString(f.s)
		case FieldBool:
			field.SetBool(f.b)
		case FieldInt:
			field.SetInt(int64(f.i))
		default:
			return fmt.Errorf("unsupported kind for field '%s': %v", f.name, f.kind)
		}
	}

	return nil
}

func (m Model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:

		// toggle edit mode on field if 'enter' key was pressed
		if msg.String() == "enter" {
			f := m.getFieldUnderCursor()
			if !m.state.isEditingValue {
				m.state.isEditingValue = true
			} else {
				f.commitEdit()
				m.state.isEditingValue = false
				if m.options.TabAfterEntry {
					m.state.cursor.decr()
				}
			}
		} else if msg.Type == tea.KeyBackspace {
			if m.state.isEditingValue {
				m.getFieldUnderCursor().handleBackspace()
			}
		} else {
			if m.state.isEditingValue {
				m.getFieldUnderCursor().handleChar(msg.String())
			} else {
				// Cool, what was the actual key pressed?
				switch msg.String() {

				case "s":
					return m, tea.Quit

				// These keys should exit the program.
				case "ctrl+c", "q":
					m.EndState.QuitWithCancel = true
					return m, tea.Quit

				// The "up" and "k" keys move the cursor up, or users may tab backward.
				case "up", "k", "shift+tab":
					m.state.cursor.incr()

				// The "down" and "j" keys move the cursor down, or users may tab forward.
				case "down", "j", "tab":
					m.state.cursor.decr()

				}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m Model) View() string {
	var s string
	// Add the header, if it exists
	if m.options.header != "" {
		s = m.options.header + "\n"
	}
	s += "\n"

	// for formatting, get longest field name
	maxFieldName := 0
	for _, field := range m.menuFields {
		if fieldName := field.getFieldName(); len(fieldName) > maxFieldName {
			maxFieldName = len(fieldName)
		}
	}

	// for formatting, get longest cursor string and build
	// the empty version of the cursor based on its length
	cursorEmpty := ""
	for _, cursor := range []string{m.options.NavCursorChar, m.options.EditCursorChar} {
		if len(cursor) > len(cursorEmpty) {
			cursorEmpty = ""
			for range cursor {
				cursorEmpty += " "
			}
		}
	}

	// Iterate over our fields
	for i, f := range m.menuFields {

		// Is the cursor pointing at this choice?
		cursor := "  " // no cursor
		if i == m.state.cursor.idx() {
			if m.state.isEditingValue {
				cursor = m.options.EditCursorChar
			} else {
				cursor = m.options.NavCursorChar
			}
		}

		// string represenation of field value
		value := f.render(m.state.isEditingValue && m.state.cursor.idx() == i, m.options.IBeamChar)
		s += fmt.Sprintf("%s ⟦ %-*s ⟧: %s\n", cursor, maxFieldName, f.getFieldName(), value)
	}

	// The footer
	s += "\n"
	if smDes := m.getFieldAtIndex(m.state.cursor.idx()).smDes; smDes != "" {
		s += smDes
	}
	s += "\n"

	s += "\nPress s to save and quit.\nPress q to quit without saving.\n"
	if f := m.getFieldUnderCursor(); f.errBuf != "" {
		s += fmt.Sprintf("ERROR: %s\n", f.errBuf)
	}

	// Send the UI for rendering
	return s
}
