package ironclad

import (
	"net/http"
)

type Common struct {
	Message  string
	Query    string
	Category Category
	Order    SortOrder
}

func NewCommon(r *http.Request) Common {
	return Common{
		Message:  "",
		Query:    r.FormValue("q"),
		Category: ParseCategory(r.FormValue("category")),
		Order:    ParseSortOrder(r.FormValue("order")),
	}
}

func (c Common) SuffixWith(query string, category Category, order SortOrder) string {
	if query == "" {
		query = c.Query
	}
	if category == Category(0) {
		category = c.Category
	}
	if order == SortOrder("") {
		order = c.Order
	}

	r := ""
	if query != "" {
		r += "q=" + query
	}
	if category != General {
		if r != "" {
			r += "&"
		}
		r += "category=" + category.ForURL()
	}
	if order != AgeNewToOld {
		if r != "" {
			r += "&"
		}
		r += "order=" + string(order)
	}
	if r != "" {
		return "?" + r
	}
	return ""
}

func (c Common) URLSuffix() string {
	return c.SuffixWith("", Category(0), SortOrder(""))
}

func (c Common) AllCategories() []Category {
	return AllCategories
}
