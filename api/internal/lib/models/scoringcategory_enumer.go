// Code generated by "enumer -sql -transform snake-upper -trimprefix ScoringCategory -type ScoringCategory ./api/internal/lib/models"; DO NOT EDIT.

package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

const _ScoringCategoryName = "UNKNOWNPUB_GOLF_NINE_HOLEPUB_GOLF_FIVE_HOLEPUB_GOLF_CHALLENGES"

var _ScoringCategoryIndex = [...]uint8{0, 7, 25, 43, 62}

const _ScoringCategoryLowerName = "unknownpub_golf_nine_holepub_golf_five_holepub_golf_challenges"

func (i ScoringCategory) String() string {
	if i < 0 || i >= ScoringCategory(len(_ScoringCategoryIndex)-1) {
		return fmt.Sprintf("ScoringCategory(%d)", i)
	}
	return _ScoringCategoryName[_ScoringCategoryIndex[i]:_ScoringCategoryIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _ScoringCategoryNoOp() {
	var x [1]struct{}
	_ = x[ScoringCategoryUnknown-(0)]
	_ = x[ScoringCategoryPubGolfNineHole-(1)]
	_ = x[ScoringCategoryPubGolfFiveHole-(2)]
	_ = x[ScoringCategoryPubGolfChallenges-(3)]
}

var _ScoringCategoryValues = []ScoringCategory{ScoringCategoryUnknown, ScoringCategoryPubGolfNineHole, ScoringCategoryPubGolfFiveHole, ScoringCategoryPubGolfChallenges}

var _ScoringCategoryNameToValueMap = map[string]ScoringCategory{
	_ScoringCategoryName[0:7]:        ScoringCategoryUnknown,
	_ScoringCategoryLowerName[0:7]:   ScoringCategoryUnknown,
	_ScoringCategoryName[7:25]:       ScoringCategoryPubGolfNineHole,
	_ScoringCategoryLowerName[7:25]:  ScoringCategoryPubGolfNineHole,
	_ScoringCategoryName[25:43]:      ScoringCategoryPubGolfFiveHole,
	_ScoringCategoryLowerName[25:43]: ScoringCategoryPubGolfFiveHole,
	_ScoringCategoryName[43:62]:      ScoringCategoryPubGolfChallenges,
	_ScoringCategoryLowerName[43:62]: ScoringCategoryPubGolfChallenges,
}

var _ScoringCategoryNames = []string{
	_ScoringCategoryName[0:7],
	_ScoringCategoryName[7:25],
	_ScoringCategoryName[25:43],
	_ScoringCategoryName[43:62],
}

// ScoringCategoryString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func ScoringCategoryString(s string) (ScoringCategory, error) {
	if val, ok := _ScoringCategoryNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _ScoringCategoryNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to ScoringCategory values", s)
}

// ScoringCategoryValues returns all values of the enum
func ScoringCategoryValues() []ScoringCategory {
	return _ScoringCategoryValues
}

// ScoringCategoryStrings returns a slice of all String values of the enum
func ScoringCategoryStrings() []string {
	strs := make([]string, len(_ScoringCategoryNames))
	copy(strs, _ScoringCategoryNames)
	return strs
}

// IsAScoringCategory returns "true" if the value is listed in the enum definition. "false" otherwise
func (i ScoringCategory) IsAScoringCategory() bool {
	for _, v := range _ScoringCategoryValues {
		if i == v {
			return true
		}
	}
	return false
}

func (i ScoringCategory) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *ScoringCategory) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of ScoringCategory: %[1]T(%[1]v)", value)
	}

	val, err := ScoringCategoryString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}