// Code generated by "stringer -type=Gender -linecomment -output=gender_string.go"; DO NOT EDIT.

package operation

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[UNIDENTIFIED-0]
	_ = x[MALE-1]
	_ = x[FEMALE-2]
}

const _Gender_name = "UNIDENTIFIEDMALEFEMALE"

var _Gender_index = [...]uint8{0, 12, 16, 22}

func (i Gender) String() string {
	if i < 0 || i >= Gender(len(_Gender_index)-1) {
		return "Gender(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Gender_name[_Gender_index[i]:_Gender_index[i+1]]
}