package paging

import "math"

func GetPagingData(totalCount int32, pageSize int32) int32 {
	return int32(math.Ceil(float64(totalCount) / float64(pageSize)))
}
