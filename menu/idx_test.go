package menu

import (
	"maps"
	"reflect"
	"testing"
)

func TestGetOrderedFields(t *testing.T) {
	type idxTest struct {
		expected map[int]int
		input    any
		name     string
		wantErr  bool
	}

	idxTestsIsolated := []idxTest{
		{
			name: "no idx validation returns empty with no error",
			input: struct {
				s string
				i int
				b bool
			}{},
			expected: map[int]int{},
			wantErr:  false,
		},
		{
			name: "idx validation returns indeces as specified per field",
			input: struct {
				s string `idx:"2"`
				i int    `idx:"0"`
				b bool   `idx:"1"`
			}{},
			expected: map[int]int{
				2: 0,
				0: 1,
				1: 2,
			},
			wantErr: false,
		},
		{
			name: "idx validation not starting at 0 returns nil with error",
			input: struct {
				s string `idx:"3"`
				i int    `idx:"1"`
				b bool   `idx:"2"`
			}{},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "idx validation out of sequence returns nil with error",
			input: struct {
				s string `idx:"3"`
				i int    `idx:"0"`
				b bool   `idx:"1"`
			}{},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "missing idx tag returns nil with error",
			input: struct {
				s string `idx:"2"`
				b bool
				i int `idx:"1"`
			}{},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "missing idx on first field enforces non-presence",
			input: struct {
				s string
				i int  `idx:"0"`
				b bool `idx:"1"`
			}{},
			expected: nil,
			wantErr:  true,
		},
	}
	idxTestsWithBlacklistTag := []idxTest{
		{
			name: "idx validation skips bl-tagged field (first)",
			input: struct {
				s string `bl:""`
				i int    `idx:"0"`
				b bool   `idx:"1"`
			}{},
			expected: map[int]int{
				1: 2,
				0: 1,
			},
			wantErr: false,
		},
		{
			name: "idx validation skips bl-tagged field (middle)",
			input: struct {
				s string `idx:"1"`
				i int    `bl:""`
				b bool   `idx:"0"`
			}{},
			expected: map[int]int{
				1: 0,
				0: 2,
			},
			wantErr: false,
		},
		{
			name: "idx validation skips bl-tagged field (last)",
			input: struct {
				s string `idx:"1"`
				i int    `idx:"0"`
				b bool   `bl:""`
			}{},
			expected: map[int]int{
				1: 0,
				0: 1,
			},
			wantErr: false,
		},
		{
			name: "errors with incompatible tags idx and bl",
			input: struct {
				s string `idx:"2" bl:""`
				i int    `idx:"0"`
				b bool   `bl:""`
			}{},
			expected: nil,
			wantErr:  true,
		},
	}

	tests := []struct {
		name  string
		batch []idxTest
	}{
		{
			// test that idx-tagged structs work in isolation
			name:  "idx tag logic (isolated)",
			batch: idxTestsIsolated,
		},
		{
			// test idx tag interoperability with bl tag
			name:  "idx tag logic (bl tag compatibility)",
			batch: idxTestsWithBlacklistTag,
		},
	}

	for _, tb := range tests {
		t.Run(tb.name, func(t *testing.T) {
			for _, tc := range tb.batch {
				t.Run(tc.name, func(t *testing.T) {
					rType := reflect.TypeOf(tc.input)
					tags, err := getOrderedFields(rType)
					if (err != nil) != tc.wantErr {
						t.Errorf("got unexpected error: %v", err)
					}
					if !maps.Equal(tags, tc.expected) {
						t.Errorf("expected: %v, got: %v", tc.expected, tags)
					}
				})
			}
		})
	}
}

func TestIDXMemoryLayout(t *testing.T) {
	type inoptimalOrderForm struct {
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

	expType := reflect.TypeFor[inoptimalOrderForm]()
	idxType := reflect.TypeFor[idxOrderForm]()
	inopTSize := expType.Size()
	idxTSize := idxType.Size()
	if idxTSize > inopTSize {
		t.Errorf("expected size of form ordered by idx tags (%d) to be of lower size than explicit counterpart (%d)", idxTSize, inopTSize)
	} else {
		t.Logf("idx-tagged struct (%d bytes) < inoptimal struct (%d bytes)", idxTSize, inopTSize)
	}
}
