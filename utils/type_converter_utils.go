package utils

import (
	"github.com/marc7806/notion-cache/types"
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
