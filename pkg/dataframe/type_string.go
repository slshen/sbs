// Code generated by "stringer -type=Type,SummaryType"; DO NOT EDIT.

package dataframe

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Int-0]
	_ = x[Float-1]
	_ = x[String-2]
	_ = x[Invalid-3]
}

const _Type_name = "IntFloatStringInvalid"

var _Type_index = [...]uint8{0, 3, 8, 14, 21}

func (i Type) String() string {
	if i < 0 || i >= Type(len(_Type_index)-1) {
		return "Type(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[None-0]
	_ = x[Sum-1]
	_ = x[Average-2]
}

const _SummaryType_name = "NoneSumAverage"

var _SummaryType_index = [...]uint8{0, 4, 7, 14}

func (i SummaryType) String() string {
	if i < 0 || i >= SummaryType(len(_SummaryType_index)-1) {
		return "SummaryType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _SummaryType_name[_SummaryType_index[i]:_SummaryType_index[i+1]]
}
