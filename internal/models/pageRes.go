package models

import (
	"math"
)

type PageResponse struct {
	CurrentPageNumber int   `json:"currentPage"`
	TotalSize         int   `json:"totalSize"`
	TotalPages        []int `json:"totalPages"`
	MoveNext          *int  `json:"moveNext"`
	MovePrevious      *int  `json:"movePrevious"`
}

func (p *PageResponse) PagesCount(currentPage int, pageSize int, totalSize int) {
	pageLimit := 10
	p.CurrentPageNumber = currentPage
	p.TotalSize = totalSize

	pages := int(totalSize / pageSize)
	if pages == 0 {
		pages = 1
	} else {
		_pages := float64(totalSize) / float64(pageSize)
		pages = int(math.Ceil(_pages))
	}
	if currentPage > 10 {
		//如果目前頁數大於10 則取該區段的 +1 數字 ex 26/10 * 10 + 1 = 21   , mod  = 0 則不用加1
		currentPage = int(currentPage/pageLimit) * 10
		if currentPage%pageLimit > 0 {
			currentPage += 1
		} else {
			currentPage = int((currentPage-1)/pageLimit)*10 + 1 //  mod 相等 需停留在同一區段頁碼 故 -1
		}
	} else {
		currentPage = 1
		p.MovePrevious = nil
	}

	var totalPages []int
	tmpLimist := pageLimit
	for i := 1; i <= pages; i++ {
		if tmpLimist == 0 {
			break
		}
		totalPages = append(totalPages, i)
		tmpLimist--
	}

	//是否有下10筆
	nextLast := totalPages[len(totalPages)-1]
	nextLimit := int(totalSize / pageLimit)
	if nextLast < nextLimit {
		next := nextLast + 1
		p.MoveNext = &next
	}

	//是否有前10筆
	previousFirst := totalPages[0]
	previousLimit := pageLimit
	if previousFirst > previousLimit {
		previous := previousFirst - 1
		p.MovePrevious = &previous
	}
	p.TotalPages = totalPages
}
