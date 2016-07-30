// +build !appengine

package ironclad

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"html/template"
	"net/http"
	"os"

	"golang.org/x/net/context"
)

// storing listings in JSON

type storage struct {
	Listings []Listing `json:"listings"`
}

func (s storage) Save() error {
	fp, err := os.Create("ironclad.json")
	if err != nil {
		return err
	}

	if err := json.NewEncoder(fp).Encode(s); err != nil {
		return err
	}
	return nil
}

func newStorage() (*storage, error) {
	fp, err := os.Open("ironclad.json")
	if err != nil {
		return nil, err
	}

	s := &storage{}
	err = json.NewDecoder(fp).Decode(s) // ignore errors
	return s, err
}

func mustStorage() *storage {
	s, err := newStorage()
	if err != nil {
		return &storage{Listings: []Listing{}}
	}
	return s
}

var memoryStorage = mustStorage()

// helper functions

func contextForRequest(r *http.Request) context.Context {
	return context.Background()
}

func staticAssets() http.FileSystem {
	return http.Dir(".")
}

func templateAssets() *template.Template {
	return template.Must(template.New("").Funcs(
		template.FuncMap{
			"many": func(n int) []bool { return make([]bool, n) },
		}).ParseGlob("tmpl/*.html"))
}

func allocateID() (ID, error) {
	data := make([]byte, 30)
	_, err := rand.Read(data)
	if err != nil {
		return "", err
	}
	return ID(base64.RawURLEncoding.EncodeToString(data)), nil
}

func lookup(c context.Context, id ID) (*Listing, error) {
	for _, l := range memoryStorage.Listings {
		if l.ID == id {
			return &l, nil
		}
	}

	return nil, nil
}

func persist(c context.Context, l *Listing) error {
	if l.ID == "" {
		var err error
		l.ID, err = allocateID()
		if err != nil {
			return err
		}
	}

	for idx, list := range memoryStorage.Listings {
		if l.ID == list.ID {
			memoryStorage.Listings[idx] = *l
			return nil
		}
	}

	memoryStorage.Listings = append(memoryStorage.Listings, *l)
	memoryStorage.Save()
	return nil
}

func scan(c context.Context) ([]Listing, error) {
	return memoryStorage.Listings, nil
}
