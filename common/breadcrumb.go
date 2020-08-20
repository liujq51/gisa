package common

type Breadcrumb struct {
	Label string
	Url   string
}

type Breadcrumbs struct {
	items []Breadcrumb
}

//添加Breadcrumb
func (this *Breadcrumbs) AddBreadcrumbs(label, url string) {
	if this.items == nil {
		this.items = make([]Breadcrumb, 0)
	}

	breadcrumb := Breadcrumb{}
	breadcrumb.Label = label
	breadcrumb.Url = url

	this.items = append(this.items, breadcrumb)
}

//获取breadcrumbs的html代码
func (this *Breadcrumbs) ShowBreadcrumbs() string {
	var html = `<ol class="breadcrumb float-sm-right">`
	for _, item := range this.items {
		html += `<li class="breadcrumb-item">`

		if item.Url == "" {
			html += item.Label
		} else {
			html += "<a href='" + item.Url + "'>"
			html += item.Label
			html += "</a>"
		}

		html += "</li>"
	}

	html += "</ol>"

	return html
}
