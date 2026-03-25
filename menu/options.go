package menu

type MenuOptions struct {
	Header         string // message to display above the struct menu
	NavCursorChar  string // cursor during navigation
	EditCursorChar string // cursor during edit
	IBeamChar      string // character shown right of text during edit
	TabAfterEntry  bool   // whether or not to jump to the next field after field value entry
}

// Defaults returns a copy of default MenuOption values.
func (m *MenuOptions) Defaults() MenuOptions {
	return MenuOptions{
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
func (m *MenuOptions) Init() {
	*m = m.Defaults()
}
