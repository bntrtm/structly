package menu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestExceptionIndicatorConstants exists as a sanity
// check to ensure that the BlacklistIndicator and
// WhitelistIndicator constants are distinct.
func TestExceptionIndicatorConstants(t *testing.T) {
	if BlacklistIndicator == WhitelistIndicator {
		t.Errorf("expected distinct exception indicators, but both are declared as '%s'", BlacklistIndicator)
	}
}

func TestExceptionConvenienceWrappers(t *testing.T) {
	tests := []struct {
		name            string
		f               func(...string) []string
		expectIndicator string
	}{
		{
			name:            "Test blacklisting wrapper 'Black'",
			f:               Black,
			expectIndicator: BlacklistIndicator,
		},
		{
			name:            "Test whitelisting wrapper 'White'",
			f:               White,
			expectIndicator: WhitelistIndicator,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("panics to report unnecessary call", func(t *testing.T) {
				assert.Panics(t, func() { tc.f() })
			})
			t.Run("appends indicator", func(t *testing.T) {
				exceptionList := tc.f("fieldName1", "fieldName2")
				if exceptionList[len(exceptionList)-1] != tc.expectIndicator {
					t.Fail()
				}
			})
		})
	}
}

func TestValidateExceptionList(t *testing.T) {
	t.Run("empty list is validated", func(t *testing.T) {
		ind, elist, err := validateExceptionList(nil)
		if ind != "" {
			t.Errorf("empty list unexpectedly returned indicator: %s", ind)
		}
		if len(elist) != 0 {
			t.Errorf("empty exception list unexpectedly returned elements: %s", elist)
		}
		if err != nil {
			t.Errorf("expected no error, but got: %s", err)
		}
	})

	t.Run("errors with only indicator provided", func(t *testing.T) {
		_, _, err := validateExceptionList([]string{BlacklistIndicator})
		if err == nil {
			t.Errorf("expected error with only exception indicator provided, but got none")
		}
	})

	t.Run("errors with no indicator appended", func(t *testing.T) {
		_, _, err := validateExceptionList([]string{"fieldName1", "fieldName2"})
		if err == nil {
			t.Errorf("expected error with no indicator provided, but got none")
		}
	})

	t.Run("pops indicator", func(t *testing.T) {
		t.Run("with distinct last exception", func(t *testing.T) {
			ind, elist, err := validateExceptionList([]string{"fieldName1", "fieldName2", BlacklistIndicator})
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}
			if ind != BlacklistIndicator {
				t.Errorf("expected validation to return provided indicator, but got indicator as: %s", ind)
			}
			if ind == elist[len(elist)-1] {
				t.Errorf("expected validation to return exception list without indicator appended")
			}
		})
		// in an edge case where an exception happens to match an indicator,
		// we must ensure that it is not left out of the list
		t.Run("with identical last exception", func(t *testing.T) {
			for _, withInd := range []string{BlacklistIndicator, WhitelistIndicator} {
				ind, elist, err := validateExceptionList([]string{"fieldName1", withInd, withInd})
				if err != nil {
					t.Errorf("unexpected error: %s", err)
				}
				if ind != withInd {
					t.Errorf("expected validation to return provided indicator, but got indicator as: %s", ind)
				}
				if withInd != elist[len(elist)-1] {
					t.Errorf("expected validation to return exception matching indicator value, but failed")
				}
			}
		})
	})

	t.Run("validates blacklist exceptions", func(t *testing.T) {
		t.Run("with literal arguments", func(t *testing.T) {
			_, _, err := validateExceptionList([]string{"fieldName1", "fieldName2", BlacklistIndicator})
			if err != nil {
				t.Errorf("expected call to validate two blacklist exceptions, but got error: %s", err)
			}
		})

		t.Run("with convenience wrapper", func(t *testing.T) {
			_, _, err := validateExceptionList(Black("fieldName1", "fieldName2"))
			if err != nil {
				t.Errorf("expected to validate exceptions, but got error: %s", err)
			}
		})
	})

	t.Run("validates whitelist exceptions", func(t *testing.T) {
		t.Run("with literal arguments", func(t *testing.T) {
			_, _, err := validateExceptionList([]string{"fieldName1", "fieldName2", WhitelistIndicator})
			if err != nil {
				t.Errorf("expected call to validate two whitelist exceptions, but got error: %s", err)
			}
		})

		t.Run("with convenience wrapper", func(t *testing.T) {
			_, _, err := validateExceptionList(White("fieldName1", "fieldName2"))
			if err != nil {
				t.Errorf("expected to validate exceptions, but got error: %s", err)
			}
		})
	})
}
