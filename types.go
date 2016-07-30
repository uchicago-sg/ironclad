package ironclad

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type ID string

type Listing struct {
	ID          ID        `json:"id"`
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	Seller      string    `json:"seller"`
	Price       int       `json:"price"`
	Category    Category  `json:"category"`
	Seeking     bool      `json:"seeking"`
	LastUpdated time.Time `json:"lastUpdated"`

	Sublease  bool `json:"sublease"`
	Bedrooms  int  `json:"bedrooms"`
	Bathrooms int  `json:"bathrooms"`
}

var EmptyTitle = errors.New("empty title")
var EmptySeller = errors.New("empty seller")

func (l *Listing) NormalizeAndValidate() error {
	l.Title = strings.TrimSpace(l.Title)
	l.Seller = strings.TrimSpace(l.Seller)

	if l.Title == "" {
		return EmptyTitle
	} else if l.Seller == "" {
		return EmptySeller
	}
	return nil
}

func (l Listing) FormattedPrice() string {
	if l.Price == 0 {
		return "-"
	} else {
		val := fmt.Sprintf("$%0.0f", float64(l.Price)/100.0)
		if l.Price%100 != 0 {
			val = fmt.Sprintf("$%0.2f", float64(l.Price)/100.0)
		}
		if l.Category == Housing {
			if l.Price > 1000000 {
				val = fmt.Sprintf("$%0.0fk", float64(l.Price)/100000.0)
			} else {
				val += "/mo"
			}
		}
		return val
	}
}

func (l Listing) FormattedAge() string {
	u := l.LastUpdated
	if u.IsZero() || u.Year() < 2000 {
		return "-"
	} else {
		d := time.Since(u)
		v := "moments ago"
		if d.Hours() > 3*24 {
			v = fmt.Sprintf("%d %s", u.Day(), u.Month())
			if u.Year() != time.Now().Year() {
				v += fmt.Sprintf(" %d", u.Year())
			}
			return v
		} else if d.Hours() > 24 {
			v = fmt.Sprintf("%.0f day", d.Hours()/24)
		} else if d.Hours() > 1 {
			v = fmt.Sprintf("%.0f hour", d.Hours())
		} else if d.Minutes() > 1 {
			v = fmt.Sprintf("%.0f minute", d.Minutes())
		}

		if v[:2] != "1 " {
			v += "s"
		}
		return v
	}
}
