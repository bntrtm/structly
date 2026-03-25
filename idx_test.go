package gostructui

import (
	"maps"
	"reflect"
	"testing"
)

func TestGetStructIdxMap(t *testing.T) {
	tests := []struct {
		expected map[int]int
		input    any
		name     string
		wantErr  bool
	}{
		{
			name: "no idx tags returns empty with no error",
			input: struct {
				b bool
				s string
				i int
			}{},
			expected: map[int]int{},
			wantErr:  false,
		},
		{
			name: "idx tags returns indeces as specified per field",
			input: struct {
				b bool   `idx:"2"`
				s string `idx:"0"`
				i int    `idx:"1"`
			}{},
			expected: map[int]int{
				2: 0,
				0: 1,
				1: 2,
			},
			wantErr: false,
		},
		{
			name: "idx tags not starting at 0 returns nil with error",
			input: struct {
				b bool   `idx:"3"`
				s string `idx:"1"`
				i int    `idx:"2"`
			}{},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "idx tags out of sequence returns nil with error",
			input: struct {
				b bool   `idx:"3"`
				s string `idx:"0"`
				i int    `idx:"1"`
			}{},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "missing idx tag returns nil with error",
			input: struct {
				b bool `idx:"2"`
				s string
				i int `idx:"1"`
			}{},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "missing idx on first field enforces non-presence",
			input: struct {
				b bool
				s string `idx:"0"`
				i int    `idx:"1"`
			}{},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rType := reflect.TypeOf(tt.input)
			tags, err := getStructIdxMap(rType)
			if (err != nil) != tt.wantErr {
				t.Errorf("got unexpected error: %s", err)
			}
			if !maps.Equal(tags, tt.expected) {
				t.Errorf("expected: %v, got: %v", tt.expected, tags)
			}
		})
	}
}

func TestIDXMemoryLayout(t *testing.T) {
	type explicitOrderForm struct {
		string1 string //nolint
		bool1   bool   //nolint
		string2 string //nolint
		bool2   bool   //nolint
		string3 string //nolint
	}

	type idxOrderForm struct {
		string1 string `idx:"0"` //nolint
		string2 string `idx:"2"` //nolint
		string3 string `idx:"5"` //nolint
		bool1   bool   `idx:"1"` //nolint
		bool2   bool   `idx:"3"` //nolint
	}

	expType := reflect.TypeFor[explicitOrderForm]()
	idxType := reflect.TypeFor[idxOrderForm]()
	expTSize := expType.Size()
	t.Logf("explicitly ordered struct is %d bytes\n", expTSize)
	idxTSize := idxType.Size()
	t.Logf("idx tag-ordered struct is %d bytes\n", idxTSize)
	if idxTSize > expTSize {
		t.Errorf("expected size of form ordered by idx tags (%d) to be of lower size than explicit counterpart (%d)", idxTSize, expTSize)
	}
}
