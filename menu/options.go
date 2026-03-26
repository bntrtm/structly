package menu

import "strings"

type MenuOptions struct {
	header         string // message to display above the struct menu
	NavCursorChar  string // cursor during navigation
	EditCursorChar string // cursor during edit
	IBeamChar      string // character shown right of text during edit
	TabAfterEntry  bool   // whether or not to jump to the next field after field value entry
}

// NewMenuOptions returns a new Menu Options type,
// initialized with default values.
func NewMenuOptions() *MenuOptions {
	m := &MenuOptions{}
	m.init()
	return m
}

// Defaults returns a copy of default MenuOption values.
func (m MenuOptions) Defaults() MenuOptions {
	return MenuOptions{
		header:         "",
		IBeamChar:      "|",
		NavCursorChar:  "> ",
		EditCursorChar: ">>",
		TabAfterEntry:  true,
	}
}

// Init initializes the menu settings with default values.
// When using custom settings, this should be called first,
// before then overriding specific default values with
// those desired.
func (m *MenuOptions) init() {
	*m = m.Defaults()
}

// SetHeader sets the internal header to the value
// provided, trimming any leading or trailing whitespace.
func (m *MenuOptions) SetHeader(str string) {
	m.header = strings.TrimSpace(str)
}
