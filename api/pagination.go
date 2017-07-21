package api

import (
	fmt "fmt"
	"github.com/cohousing/cohousing-api/domain"
)

type ObjectList struct {
	Objects []interface{} `json:"objects"`
	domain.DefaultHalResource
}

func addPaginationLinks(objectList *ObjectList, baseUrl string, currentPage, recordsPerPage, totalRecordCount int) {
	var firstPage int = 1

	var lastPage int = totalRecordCount / recordsPerPage
	if lastPage < 1 {
		lastPage = 1
	}

	var prevPage int = currentPage - 1
	var nextPage int = currentPage + 1

	objectList.AddLink(domain.REL_SELF, generatePaginationUrl(baseUrl, currentPage))

	if firstPage != lastPage {
		objectList.AddLink(domain.REL_FIRST, generatePaginationUrl(baseUrl, firstPage))
		objectList.AddLink(domain.REL_LAST, generatePaginationUrl(baseUrl, lastPage))
	}
	if prevPage >= firstPage && prevPage < lastPage {
		objectList.AddLink(domain.REL_PREV, generatePaginationUrl(baseUrl, prevPage))
	}
	if nextPage <= lastPage {
		objectList.AddLink(domain.REL_NEXT, generatePaginationUrl(baseUrl, nextPage))
	}
}

func generatePaginationUrl(baseUrl string, page int) string {
	if page > 1 {
		return fmt.Sprintf("%s?page=%d", baseUrl, page)
	} else {
		return baseUrl
	}
}

func getStartRecord(currentPage, recordsPerPage int) int {
	return (currentPage - 1) * recordsPerPage
}
