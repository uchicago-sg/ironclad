// +build appengine

package ironclad

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// accessing data
func lookup(c context.Context, id ID) (*Listing, error) {
	if id == "" {
		return nil, nil
	}
	k := datastore.NewKey(c, "Listing", string(id), 0, nil)
	var listing Listing
	if err := datastore.Get(c, k, &listing); err != nil {
		return nil, err
	}
	return &listing, nil
}

func allocateID() (ID, error) {
	data := make([]byte, 30)
	_, err := rand.Read(data)
	if err != nil {
		return "", err
	}
	return ID(base64.RawURLEncoding.EncodeToString(data)), nil
}

func persist(c context.Context, l *Listing) error {
	if l.ID == "" {
		var err error
		l.ID, err = allocateID()
		if err != nil {
			return err
		}
	}

	k := datastore.NewKey(c, "Listing", string(l.ID), 0, nil)
	key, err := datastore.Put(c, k, l)
	l.ID = ID(key.StringID())
	return err
}

func scan(c context.Context) (listings []Listing, err error) {
	q := datastore.NewQuery("Listing").Project("Title", "Category", "Seller", "Seeking")
	keys, err := q.GetAll(c, &listings)
	for i := range listings {
		listings[i].ID = ID(keys[i].StringID())
	}
	return
}

// preparing routing
func contextForRequest(r *http.Request) context.Context {
	return appengine.NewContext(r)
}

func staticAssets() http.FileSystem {
	return nil
}

func templateAssets() *template.Template {
	return template.Must(template.New("").Funcs(
		template.FuncMap{
			"many": func(n int) []bool { return make([]bool, n) },
			"add":  func(a, b int) int { return a + b },
			"sub":  func(a, b int) int { return a - b },
		}).ParseGlob("tmpl/*.html"))
}

func init() {
	http.Handle("/", New())
}

type Secret struct {
	Value []byte
}

func sharedSecret(c context.Context) []byte {
	s := Secret{}
	k := datastore.NewKey(c, "Secret", "secret", 0, nil)
	err := datastore.Get(c, k, &s)

	if err == datastore.ErrNoSuchEntity {
		s.Value = make([]byte, 32)
		if _, err := rand.Read(s.Value); err != nil {
			panic(err)
		}
		if _, err := datastore.Put(c, k, &s); err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	return s.Value
}
