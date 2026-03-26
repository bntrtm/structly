package menu

// NewCursor returns a new cursor instance.
// The cursor must be assigned to a non-nil context to persist.
// If the cursor could not be created, the pointer returned will
// be nil.
func NewCursor(context []menuField, index int) *cursor {
	c := &cursor{
		index: -1,
		ctx:   nil,
	}
	c.focus(context, index)
	if !c.valid() {
		return nil
	}

	return c
}

type cursor struct {
	index int
	ctx   []menuField
}

// valid returns the whether or not the cursor is valid.
// A valid cursor has an index >= 0 that exists within
// a known, non-nil context.
func (c cursor) valid() bool {
	if c.index < 0 {
		return false
	}
	if len(c.ctx) == 0 || c.ctx == nil {
		return false
	}
	if c.index > len(c.ctx)-1 {
		return false
	}

	return true
}

// idx returns the index of the cursor.
func (c cursor) idx() int {
	return c.index
}

// focus assigns a new context for the cursor.
func (c *cursor) focus(context []menuField, index int) {
	if len(context) == 0 || context == nil {
		return
	}
	c.ctx = context
	c.set(index)
}

// set assigns the cursor index within its context.
func (c *cursor) set(i int) {
	if len(c.ctx) == 0 || c.ctx == nil {
		// TODO: This scenario should first try to switch context,
		// where possible.
		c.index = -1
		return
	}
	if last := len(c.ctx) - 1; i > last {
		c.index = last
		return
	} else if i < 0 {
		c.index = 0
	} else {
		c.index = i
	}
}

// incr increases the field index the user is focused on
func (c *cursor) incr() {
	if c.idx() > 0 {
		c.under().errBuf = ""
		c.index--
	}
}

// decrCursor decreases the index the user is focused on
func (c *cursor) decr() {
	if c.idx() < len(c.ctx)-1 {
		c.under().errBuf = ""
		c.index++
	}
}

// under returns the menu field under the cursor.
func (c cursor) under() *menuField {
	if c.ctx == nil {
		return nil
	}
	return &c.ctx[c.index]
}
