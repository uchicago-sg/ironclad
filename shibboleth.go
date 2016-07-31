package ironclad

import (
	"golang.org/x/net/context"
	"net/http"
	"strings"
)

type LoginResult struct {
	NeedsShib  bool
	NewSubject *Subject
}

func (l LoginResult) Template() string  { return "" }
func (l LoginResult) Subject() *Subject { return l.NewSubject }

func (l LoginResult) NewURL() string {
	return "/"
}

func LoginPage(s *Subject, c context.Context, r *http.Request) (Template, error) {
	email := r.FormValue("email")
	shib := strings.HasSuffix(email, "@uchicago.edu")

	if !shib {
		s = newSubject(email, email)
	}

	return &LoginResult{
		NeedsShib:  shib,
		NewSubject: s,
	}, nil
}

func LogoutPage(s *Subject, c context.Context, r *http.Request) (Template, error) {
	return &LoginResult{
		NewSubject: nil,
	}, nil
}
