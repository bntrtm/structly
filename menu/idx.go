package menu

import (
	"fmt"
	"reflect"
	"strconv"
)

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
// When no `idx` tags are used, the map keys and values will match.
//
// Where validation fails, a nil map and error are returned.
func getOrderedFields(t reflect.Type) (map[int]int, error) {
	wantIdx := struct {
		val   bool
		isSet bool
	}{val: false, isSet: false}

	idxTagVals := map[int]int{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		_, isBlacklisted := field.Tag.Lookup("bl")
		tagValue, isIndexed := field.Tag.Lookup("idx")

		if isBlacklisted {
			if isIndexed {
				return nil, fmt.Errorf("incompatible struct tags; unexpected `idx` tag found on `bl`-tagged field %s", field.Name)
			}
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
			if val, ok := idxTagVals[idx]; ok {
				return nil, fmt.Errorf("value %d for `idx` tag on field %s already assigned to another field", val, field.Name)
			}
			idxTagVals[idx] = i
		} else if isIndexed {
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
