package common

import (
	"math"
	"strconv"

	"github.com/astaxie/beego"
)

type Pagination struct {
	PageIndex int
	PageCount int
	PageTotal int
	Url       string
}

//输出分页
func PaginationRender(page Pagination) string {
	pageHtml := `总共 <b>` + strconv.Itoa(page.PageCount) + `</b> 页 <b>` + strconv.Itoa(page.PageTotal) + `</b> 条`
	pageHtml += `<ul class="pagination pagination-sm no-margin float-right">`

	if page.PageIndex == 1 {
		pageHtml += `<li class="page-item disabled"><a href="#" aria-label="Previous"><span class="page-link" >&laquo;</span></a></li>`
	} else {
		pageHtml += `<li class="page-item "><a href="` + page.Url + `?page_count=` + strconv.Itoa(page.PageCount) + `&page_index=` + strconv.Itoa(page.PageIndex-1) + `" aria-label="Previous"><span class="page-link" aria-hidden="true">&laquo;</span></a></li>`
	}

	for i := 1; i <= int(math.Ceil(float64(page.PageTotal)/float64(page.PageCount))); i++ {
		isActive := ""
		if page.PageIndex == i {
			isActive = " active "
		}
		pageHtml += `<li class="page-item ` + isActive + `"><a href="` + page.Url + `?page_count=` + strconv.Itoa(page.PageCount) + `&page_index=` + strconv.Itoa(i) + `"><span class="page-link">` + strconv.Itoa(i) + `</span></a></li>`
	}

	if page.PageIndex == int(math.Ceil(float64(page.PageTotal)/float64(page.PageCount))) {
		pageHtml += `<li class="page-item disabled" ><a href="#" aria-label="Next"><span class="page-link" aria-hidden="true">&raquo;</span></a></li>`
	} else {
		pageHtml += `<li class="page-item "><a href="` + page.Url + `?page_count=` + strconv.Itoa(page.PageCount) + `&page_index=` + strconv.Itoa(page.PageIndex+1) + `" aria-label="Next"><span class="page-link" aria-hidden="true">&raquo;</span></a></li>`
	}
	pageHtml += `</ul>`
	pageHtml += `<label class="control-label float-right" style="margin-right: 10px; font-weight: 100;">
                        <small>显示</small>&nbsp;<select class="input-sm grid-per-pager" name="per-page">`
	pageCountList := beego.AppConfig.Strings("page::pageCountList")
	for _, pageCount := range pageCountList {
		isSelected := ""
		if strconv.Itoa(page.PageCount) == pageCount {
			isSelected = " selected "
		}
		pageHtml += `       <option value="` + page.Url + `?page_count=` + pageCount + `" ` + isSelected + `>` + pageCount + `</option>`
	}
	pageHtml += `</select><small>条</small></label>`
	return pageHtml
}
