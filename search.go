package ironclad

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/net/context"
)

type Engine struct {
	sync.Mutex
}

func (e *Engine) Search(
	c context.Context,
	query string,
	category Category,
	order SortOrder,
) ([]Listing, map[Category]int, error) {

	e.Lock()
	defer e.Unlock()

	listings, err := scan(c)
	if err != nil {
		return nil, nil, err
	}

	results := []Listing{}
	categories := map[Category]int{}

	for _, listing := range listings {
		if listing.Matches(query) {
			if listing.Category == category {
				results = append(results, listing)
			}
			categories[listing.Category]++
		}
	}

	sort.Sort(byOrder{results, order})

	return results, categories, nil
}

func Paginate(listings []Listing, offset int) (prev, off, next int) {
	off = offset
	if off >= len(listings) {
		off = len(listings)
	}
	if off < 0 {
		off = 0
	}
	next = off + 100
	if next >= len(listings) {
		next = len(listings)
	}
	prev = off - 100
	if prev < 0 {
		prev = 0
	}
	return
}

func (l Listing) Matches(query string) bool {
	return strings.Contains(strings.ToLower(l.Title), strings.ToLower(query))
}

// define how results are shown

var engine Engine

type SearchResults struct {
	Previous, Offset, Next, Total int

	Listings       []Listing
	Categories     map[Category]int
	ChangeCategory bool

	Common
}

func (r SearchResults) Template() string { return "SearchResults.html" }

func (r SearchResults) NewURL() string {
	if r.ChangeCategory {
		return "/" + r.Common.URLSuffix()
	}
	return ""
}

func SearchListings(s *Subject, c context.Context, r *http.Request) (Template, error) {
	com := NewCommon(s, r)

	results, cats, err := engine.Search(c, com.Query, com.Category, com.Order)
	if err != nil {
		return nil, err
	}

	off, _ := strconv.ParseInt(r.FormValue("offset"), 10, 64)
	previous, offset, next := Paginate(results, int(off))

	page := SearchResults{
		Previous: previous,
		Offset:   offset,
		Next:     next,
		Total:    len(results),

		Listings:   results[offset:next],
		Categories: cats,
		Common:     com,
	}

	// if we could have shown results, but didn't, because we were
	// in a weird category, redirect to the right one
	if len(results) == 0 {
		for _, category := range AllCategories {
			if cats[category] > 0 {
				page.Common.Category = category
				page.ChangeCategory = true
				break
			}
		}
	}

	return page, nil
}
