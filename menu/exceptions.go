package menu

import "fmt"

// these constants are used to indicate exception lists
// as blacklists or whitelists.
const (
	BlacklistIndicator = "BL"
	WhitelistIndicator = "WL"
)

// Black accepts one or more strings and returns a slice of strings
// with an exception indicator for blacklisting appended to it.
//
// It is a convenience wrapper used to define a blacklist to be passed to
// exception variadics in NewMenu functions. If the input list contains
// no elements, Black will panic, as the call is unnecessary.
func Black(list ...string) []string {
	if len(list) == 0 {
		panic("unnecessary call to define blacklist; no inputs provided")
	}
	return append(list, BlacklistIndicator)
}

// White accepts one or more strings and returns a slice of strings
// with an exception indicator for whitelisting appended to it.
//
// It is a convenience wrapper used to define a whitelist to be passed to
// exception variadics in NewMenu functions. If the input list contains
// no elements, White will panic, as the call is unnecessary.
func White(list ...string) []string {
	if len(list) == 0 {
		panic("unnecessary call to define whitelist; no inputs provided")
	}
	return append(list, WhitelistIndicator)
}

// validateExceptionList ensures that the input string slice,
// if not empty, has a last element matching a constant that
// would define the slice as a blacklist or whitelist.
//
// If a non-zero -length slice successfully validates, the
// constant is popped and returned alongside  the slice of
// remaining elements.
//
// Zero-length inputs are always valid, returning zero values.
func validateExceptionList(exceptions []string) (string, []string, error) {
	if len(exceptions) == 0 {
		return "", nil, nil
	}
	const errPre = "could not validate field exceptions as whitelist or blacklist; "
	lastIndex := len(exceptions) - 1
	if last := exceptions[lastIndex]; last == BlacklistIndicator || last == WhitelistIndicator {
		if lastIndex == 0 {
			return "", nil, fmt.Errorf(errPre + "no field exceptions were provided (unnecessary list)")
		}
		return last, exceptions[:lastIndex], nil
	}
	return "", nil, fmt.Errorf(errPre+"expected last element value as an exception indicator ('%s' or '%s')", WhitelistIndicator, BlacklistIndicator)
}
