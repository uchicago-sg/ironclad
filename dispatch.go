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

type NewSubject interface {
	Subject() *Subject
}

type Handler func(s *Subject, c context.Context, r *http.Request) (Template, error)

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// read environment from the request
	c := contextForRequest(r)

	s := (*Subject)(nil)
	if cookie, _ := r.Cookie("session"); cookie != nil {
		s, _ = ParseSubject(c, cookie.Value) // TODO: possibly log this
	}

	// invoke the handler
	t, err := h(s, c, r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// write out new session
	if t, ok := t.(NewSubject); ok {
		s = t.Subject()
	}

	ss, err := s.Serialize(c)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: ss,
		Path:  "/",
	})

	// possibly redirect somewhere else
	if t, ok := t.(Redirect); ok && t.NewURL() != "" {
		http.Redirect(w, r, t.NewURL(), 303)
		return
	}

	// handle empty responses
	if t == nil {
		http.NotFound(w, r)
		return
	}

	// actually render the template
	if err := templateAssets().ExecuteTemplate(w, t.Template(), t); err != nil {
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
	mux.Handle("/login", Handler(LoginPage))
	mux.Handle("/logout", Handler(LogoutPage))
	mux.Handle("/saml-callback", Handler(SAMLRedirect))
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
