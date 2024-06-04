package utils

func Pagination(total, perPage, page int32) (startIndex, endIndex int32) {
	if total == 0 || page == 0 { return }
	startIndex = (page - 1) * perPage
	endIndex = startIndex + perPage
	if startIndex >= total {
		startIndex, endIndex = 0, 0
		return
	}
	if endIndex >= total {
		endIndex = total
	}
	return
}