package ironclad

import (
	"golang.org/x/net/context"
	"net/http"
	"strings"
)

type Template interface {
	Template() string
}

type Redirect interface {
	NewURL() string
}

type Handler func(c context.Context, r *http.Request) (Template, error)

var templates = templateAssets()

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := contextForRequest(r)
	t, err := h(c, r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if t, ok := t.(Redirect); ok && t.NewURL() != "" {
		http.Redirect(w, r, t.NewURL(), 303)
		return
	}

	if t == nil {
		http.NotFound(w, r)
		return
	}

	if err := templates.ExecuteTemplate(w, t.Template(), t); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func New() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", Handler(SearchListings))
	mux.Handle("/edit/", Handler(EditListing))
	mux.Handle("/view/", Handler(ViewListing))
	mux.Handle("/create", Handler(CreateListing))
	mux.Handle("/static/", http.FileServer(staticAssets()))
	return mux
}

func idFrom(r *http.Request) ID {
	b := strings.Split(r.URL.Path, "/")
	if len(b) >= 3 {
		return ID(b[2])
	} else {
		return ID("")
	}
}
