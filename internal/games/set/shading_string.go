// Code generated by "stringer -type=Shading"; DO NOT EDIT.

package set

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Filled-0]
	_ = x[Outline-1]
	_ = x[Stripe-2]
}

const _Shading_name = "FilledOutlineStripe"

var _Shading_index = [...]uint8{0, 6, 13, 19}

func (i Shading) String() string {
	if i >= Shading(len(_Shading_index)-1) {
		return "Shading(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Shading_name[_Shading_index[i]:_Shading_index[i+1]]
}
