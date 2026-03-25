package gostructui

import (
	"fmt"
	"reflect"
	"strconv"
)

// getStructIdxMap returns a map of `idx` tag values which
// correspond to the indeces of struct fields within the struct
// represented by the given reflect.Type, in the order they are
// declared. The presence or non-presence of the tags is
// validated when reading each field.
//
// If the first field is found to have an `idx` tag, all others
// will be expected to have one as well. Likewise, if it does not
// have the tag, all others will be expected to not have the tag.
//
// Where validation fails, a nil map and error are
// returned.
func getStructIdxMap(t reflect.Type) (map[int]int, error) {
	wantIdx := false
	idxTagVals := map[int]int{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		tagValue, ok := field.Tag.Lookup("idx")
		if ok && i == 0 {
			wantIdx = true
		}

		if wantIdx {
			if !ok {
				return nil, fmt.Errorf("no `idx` tag found on struct field %s", field.Name)
			}
			idx, err := strconv.Atoi(tagValue)
			if err != nil || idx < 0 {
				return nil, fmt.Errorf("value for `idx` tag on field %s must be an integer >= 0", field.Name)
			}
			if val, ok := idxTagVals[idx]; ok {
				return nil, fmt.Errorf("value %d for `idx` tag on field %s already assigned to another field", val, field.Name)
			}
			idxTagVals[idx] = i
		} else if ok {
			return nil, fmt.Errorf("unexpected `idx` tag found on field %s", field.Name)
		}

	}

	for i := 0; i < len(idxTagVals); i++ {
		if _, ok := idxTagVals[i]; !ok {
			return nil, fmt.Errorf("expected to find idx value of %d on some field, but found none", i)
		}
	}

	return idxTagVals, nil
}
