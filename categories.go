package ironclad

import (
	"encoding/json"
)

type Category int

const (
	Invalid Category = iota
	General
	Appliances
	Books
	Electronics
	Jobs
	Furniture
	Housing
	Cars
	Bikes
)

var AllCategories = []Category{
	General,
	Appliances,
	Books,
	Electronics,
	Jobs,
	Furniture,
	Housing,
	Cars,
	Bikes,
}

func ParseCategory(c string) Category {
	return map[string]Category{
		"":            General,
		"appliances":  Appliances,
		"books":       Books,
		"electronics": Electronics,
		"jobs":        Jobs,
		"furniture":   Furniture,
		"housing":     Housing,
		"cars":        Cars,
		"bikes":       Bikes,
	}[c]
}

func (c Category) String() string {
	return map[Category]string{
		Invalid:     "(invalid category)",
		General:     "General",
		Appliances:  "Appliances",
		Books:       "Books",
		Electronics: "Electronics",
		Jobs:        "Jobs",
		Furniture:   "Furniture",
		Housing:     "Housing",
		Cars:        "Cars",
		Bikes:       "Bikes",
	}[c]
}

func (c Category) ForURL() string {
	return map[Category]string{
		General:     "",
		Appliances:  "appliances",
		Books:       "books",
		Electronics: "electronics",
		Jobs:        "jobs",
		Furniture:   "furniture",
		Housing:     "housing",
		Cars:        "cars",
		Bikes:       "bikes",
	}[c]
}

func (c *Category) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.ForURL())
}

func (c *Category) UnmarshalJSON(b []byte) error {
	v := ""
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	*c = ParseCategory(v)
	return nil
}
