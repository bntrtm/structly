package menu

import (
	"fmt"
	"strconv"
)

type FieldKind int

const (
	FieldString FieldKind = iota
	FieldBool
	FieldInt
)

type menuField struct {
	editBuf string // buffer for editing this field
	errBuf  string // potential error from bad input

	name   string // name of the struct field
	smName string // description pulled from smname tag
	smDes  string // description pulled from smdes tag

	kind FieldKind // value assigned to field
	s    string    // possible string value
	i    int       // possible int value
	b    bool      // possible bool value
}

func (f *menuField) handleChar(char string) {
	switch f.kind {
	case FieldInt:
		if (char >= "0" && char <= "9") || (char == "-" && len(f.editBuf) == 0) {
			f.editBuf += string(char)
		}
	case FieldString:
		f.editBuf += string(char)
	case FieldBool:
		switch char {
		case "t", "1":
			f.b = true
		case "f", "0":
			f.b = false
		case "right", "left", "l", "h":
			f.b = !f.b
		}
	}
}

func (f *menuField) handleBackspace() {
	if len(f.editBuf) == 0 {
		return
	}
	f.editBuf = f.editBuf[:len(f.editBuf)-1]
}

func (f *menuField) render(editing bool, iBeamChar string) string {
	switch f.kind {
	case FieldInt:
		if editing {
			return f.editBuf + iBeamChar
		}
		return strconv.Itoa(f.i)
	case FieldString:
		if editing {
			return f.editBuf + iBeamChar
		}
		return f.s
	case FieldBool:
		if editing {
			if f.b {
				return "[t] ||  f "
			}
			return " t  || [f]"
		}
		return fmt.Sprintf("%v", f.b)
	default:
		return ""
	}
}

func (f *menuField) commitEdit() {
	switch f.kind {
	case FieldInt:
		if f.editBuf == "" || f.editBuf == "-" {
			f.i = 0
			return
		}
		v, err := strconv.Atoi(f.editBuf)
		if err != nil {
			f.errBuf = err.Error()
			return
		}
		f.i = v
	case FieldString:
		f.s = f.editBuf
	}

	f.editBuf = ""
	f.errBuf = ""
}

// getFieldName returns a name for the menu field.
// If an override name was provided via the smname tag
// (e.g. for human readability or foramtting), that will
// be returned. Otherwise, the name of the struct field
// is returned.
func (f *menuField) getFieldName() string {
	if f.smName != "" {
		return f.smName
	}
	return f.name
}
