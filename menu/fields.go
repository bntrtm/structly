package menu

import (
	"fmt"
	"reflect"
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

// getOrderedFields accounts for the presence of `idx` and `bl`
// tags to provide a map that defines the order by which each
// struct fields ought be rendered in the terminal.
//
// If the `idx` tag is in use on the struct,it returns a map of `idx`
// tag values corresponding to the indeces of struct fields within the
// struct represented by the given reflect.Type, in the order they
// are declared. The presence or non-presence of the tags is validated
// when reading each field. Fields blacklisted at the type level with
// the `bl` tag are expected not to have an idx tag.
//
// If the first non-blacklisted field is found to have an `idx` tag,
// all others will be expected to have one as well. Likewise, if it
// does not have the tag, all others will be expected not to have the tag.
// When no `idx` tags are used, the map keys and values will match, unless
// offset by 1 after and for each instance where a field is found to be
// blacklisted with the `bl` tag.
//
// Where validation fails, a nil map and error are returned.
func getOrderedFields(t reflect.Type) (map[int]int, error) {
	wantIdx := struct {
		val   bool
		isSet bool
	}{val: false, isSet: false}

	idxTagVals := map[int]int{}
	blacklistCount := 0
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		_, isBlacklisted := field.Tag.Lookup("bl")
		tagValue, isIndexed := field.Tag.Lookup("idx")

		if isBlacklisted {
			if isIndexed {
				return nil, fmt.Errorf("incompatible struct tags; unexpected `idx` tag found on `bl`-tagged field %s", field.Name)
			}
			blacklistCount++
			continue
		}

		// NOTE: at this point, can't possibly be blacklisted
		if !wantIdx.isSet {
			wantIdx.val = (len(idxTagVals) == 0 && isIndexed)
			wantIdx.isSet = true
		}

		if wantIdx.val {
			if !isIndexed {
				return nil, fmt.Errorf("no `idx` tag found on struct field %s", field.Name)
			}
			idx, err := strconv.Atoi(tagValue)
			if err != nil || idx < 0 {
				return nil, fmt.Errorf("value for `idx` tag on field %s must be an integer >= 0", field.Name)
			}
			if _, ok := idxTagVals[idx]; ok {
				return nil, fmt.Errorf("value %d for `idx` tag on field %s already assigned to another field", idx, field.Name)
			}
			idxTagVals[idx] = i

		} else if isIndexed {
			return nil, fmt.Errorf("unexpected `idx` tag found on field %s", field.Name)
		} else {
			idxTagVals[i-blacklistCount] = i
		}

	}

	for i := 0; i < len(idxTagVals); i++ {
		if _, ok := idxTagVals[i]; !ok {
			if wantIdx.val {
				return nil, fmt.Errorf("expected to find idx value of %d on some field, but found none", i)
			}
			return nil, fmt.Errorf("expected sequential indeces for map, but index %d is missing", i)

		}
	}

	return idxTagVals, nil
}
