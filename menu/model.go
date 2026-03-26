package menu

import (
	"errors"
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

// NewMenu creates a new struct menu from the given parameters.
// If customOptions are not provided, the menu will fall back to defaults.
func NewMenu(structObj any, fieldList []string, asBlacklist bool, customOptions *MenuOptions) (Model, error) {
	// if fieldList is empty, all fields are exposed to users; otherwise, it is used as a whitelist.
	// if bool parameter 'asBlacklist' is 'true', the fieldList is used as a blacklist instead of a whitelist.
	t := reflect.TypeOf(structObj)
	v := reflect.ValueOf(structObj)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
		v = v.Elem()
	} else {
		return Model{}, errors.New("structObj should be a pointer to struct, so as to have addressable fields")
	}
	if t.Kind() != reflect.Struct {
		return Model{}, fmt.Errorf("input structObj found not to be a struct")
	}
	newModel := Model{
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

	if customOptions != nil {
		newModel.options = *customOptions
	}

	orderedFields, err := getStructIdxMap(t)
	if err != nil {
		return Model{}, err
	}

	for i := 0; i < t.NumField(); i++ {
		var j int
		if len(orderedFields) == 0 {
			j = i
		} else {
			var ok bool
			j, ok = orderedFields[i]
			if !ok {
				return Model{}, fmt.Errorf("could not resolve struct field to display by declaration index %d", i)
			}
		}
		field := t.Field(j)

		if len(fieldList) != 0 {
			if asBlacklist {
				if slices.Contains(fieldList, field.Name) {
					continue
				}
			} else {
				if !(slices.Contains(fieldList, field.Name)) {
					continue
				}
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
			return Model{}, fmt.Errorf("could not parse struct")
		}
		newField.name = field.Name
		newField.smName = field.Tag.Get("smname")
		newField.smDes = field.Tag.Get("smdes")
		newModel.menuFields = append(newModel.menuFields, newField)
	}

	if len(newModel.menuFields) == 0 {
		return Model{}, fmt.Errorf("ERROR: No fields to expose to users in struct")
	}

	newModel.state.cursor = NewCursor(newModel.menuFields, 0)
	if newModel.state.cursor == nil {
		return newModel, fmt.Errorf("ERROR, but: len fields %d", len(newModel.menuFields))
	}

	return newModel, nil
}

func (m Model) ParseStruct(obj any) error {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("expected a pointer to a struct, got %v", v.Kind())
	}
	v = v.Elem()

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
