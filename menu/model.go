package menu

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
)

// TModelStructMenu is a bubbletea model that can be used to expose
// primitive struct fields to end users for input,
// as if they were elements of a menu.
type TModelStructMenu struct {
	// MENU STATE
	// fields which can be edited; populated dynamically
	menuFields     []menuField
	cursor         int  // which field our cursor is pointing at
	isEditingValue bool // tracks state of field editing
	QuitWithCancel bool // can be used to communicate whether changes ought be saved
	Options        MenuOptions
}

// incrCursor increases the field index the user is focused on
func (m *TModelStructMenu) incrCursor() {
	if m.cursor > 0 {
		m.getFieldUnderCursor().errBuf = ""
		m.cursor--
	}
}

// decrCursor decreases the field index the user is focused on
func (m *TModelStructMenu) decrCursor() {
	m.getFieldUnderCursor().errBuf = ""
	if m.cursor < len(m.menuFields)-1 {
		m.cursor++
	}
}

func (m *TModelStructMenu) getFieldAtIndex(i int) *menuField {
	return &m.menuFields[i]
}

func (m *TModelStructMenu) getFieldUnderCursor() *menuField {
	return m.getFieldAtIndex(m.cursor)
}

// InitialTModelStructMenu creates a new struct menu from the given parameters.
// If customSettings are not provided, the menu will fall back to defaults.
// If using custom menu settings, first initialize them with the setDefaults() method.
func InitialTModelStructMenu(structObj any, fieldList []string, asBlacklist bool, customSettings *MenuOptions) (TModelStructMenu, error) {
	// if fieldList is empty, all fields are exposed to users; otherwise, it is used as a whitelist.
	// if bool parameter 'asBlacklist' is 'true', the fieldList is used as a blacklist instead of a whitelist.
	t := reflect.TypeOf(structObj)
	v := reflect.ValueOf(structObj)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
		v = v.Elem()
	} else {
		return TModelStructMenu{}, errors.New("structObj should be a pointer to struct, so as to have addressable fields")
	}
	if t.Kind() != reflect.Struct {
		return TModelStructMenu{}, fmt.Errorf("input structObj found not to be a struct")
	}
	newModel := TModelStructMenu{
		isEditingValue: false,
		menuFields:     []menuField{},
		QuitWithCancel: false,
	}

	if customSettings != nil {
		newModel.Options = *customSettings
	} else {
		newModel.Options.Init()
	}
	orderedFields, err := getStructIdxMap(t)
	if err != nil {
		return TModelStructMenu{}, err
	}

	for i := 0; i < t.NumField(); i++ {
		var j int
		if len(orderedFields) == 0 {
			j = i
		} else {
			var ok bool
			j, ok = orderedFields[i]
			if !ok {
				return TModelStructMenu{}, fmt.Errorf("could not resolve struct field to display by declaration index %d", i)
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
			return TModelStructMenu{}, fmt.Errorf("could not parse struct")
		}
		newField.name = field.Name
		newField.smName = field.Tag.Get("smname")
		newField.smDes = field.Tag.Get("smdes")
		newModel.menuFields = append(newModel.menuFields, newField)
	}

	if len(newModel.menuFields) == 0 {
		return TModelStructMenu{}, fmt.Errorf("ERROR: No fields to expose to users in struct")
	}

	return newModel, nil
}

func (m TModelStructMenu) ParseStruct(obj any) error {
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

func (m TModelStructMenu) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m TModelStructMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:

		// toggle edit mode on field if 'enter' key was pressed
		if msg.String() == "enter" {
			f := m.getFieldUnderCursor()
			if !m.isEditingValue {
				m.isEditingValue = true
			} else {
				f.commitEdit()
				m.isEditingValue = false
				if m.Options.TabAfterEntry {
					m.decrCursor()
				}
			}
		} else if msg.Type == tea.KeyBackspace {
			if m.isEditingValue {
				m.getFieldUnderCursor().handleBackspace()
			}
		} else {
			if m.isEditingValue {
				m.getFieldUnderCursor().handleChar(msg.String())
			} else {
				// Cool, what was the actual key pressed?
				switch msg.String() {

				case "s":
					return m, tea.Quit

				// These keys should exit the program.
				case "ctrl+c", "q":
					m.QuitWithCancel = true
					return m, tea.Quit

				// The "up" and "k" keys move the cursor up, or users may tab backward.
				case "up", "k", "shift+tab":
					m.incrCursor()

				// The "down" and "j" keys move the cursor down, or users may tab forward.
				case "down", "j", "tab":
					m.decrCursor()

				}
			}
		}
	}

	// Return the updated TModelStructMenu to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m TModelStructMenu) View() string {
	var s string
	// Add the header, if it exists
	if m.Options.Header != "" {
		s = m.Options.Header + "\n"
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
	for _, cursor := range []string{m.Options.NavCursorChar, m.Options.EditCursorChar} {
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
		if m.cursor == i {
			if m.isEditingValue {
				cursor = m.Options.EditCursorChar
			} else {
				cursor = m.Options.NavCursorChar
			}
		}

		// string represenation of field value
		value := f.render(m.isEditingValue && m.cursor == i, m.Options.IBeamChar)
		s += fmt.Sprintf("%s ⟦ %-*s ⟧: %s\n", cursor, maxFieldName, f.getFieldName(), value)
	}

	// The footer
	s += "\n"
	if smDes := m.getFieldAtIndex(m.cursor).smDes; smDes != "" {
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
