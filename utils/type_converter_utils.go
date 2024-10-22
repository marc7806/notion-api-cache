package utils

import (
	"github.com/marc7806/notion-cache/types"
	"strconv"
)

func BinarySortDirection(sortDirection types.QuerySortDirection) int8 {
	switch sortDirection {
	case types.Asc:
		return 1
	case types.Desc:
		return -1
	default:
		return 1
	}
}

func StringToBool(str string) (bool, error) {
	boolVal, err := strconv.ParseBool(str)
	if err != nil {
		return false, err
	}
	return boolVal, nil
}
