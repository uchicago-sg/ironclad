package ironclad

import (
	"errors"
	"net/http"

	"golang.org/x/net/context"
)

type SingleListing struct {
	Redirect bool
	Edit     bool
	Listing  *Listing
	Err      error
	Common
}

var AccessDenied = errors.New("access denied")

func (e SingleListing) Template() string { return "SingleListing.html" }
func (e SingleListing) NewURL() string {
	if e.Redirect {
		if e.Edit {
			return "/edit/" + string(e.Listing.ID)
		} else {
			return "/view/" + string(e.Listing.ID)
		}
	} else {
		return ""
	}
}

func ViewListing(s *Subject, c context.Context, r *http.Request) (Template, error) {
	listing, err := lookup(c, idFrom(r))
	if err != nil || listing == nil {
		return nil, err
	}

	return SingleListing{
		Listing: listing,
		Common:  NewCommon(s, r),
	}, nil
}

func EditListing(s *Subject, c context.Context, r *http.Request) (Template, error) {
	listing, err := lookup(c, idFrom(r))
	if err != nil || listing == nil {
		return nil, err
	}

	if !s.CanEdit(listing) {
		return nil, AccessDenied
	}

	resp := SingleListing{
		Edit:    true,
		Listing: listing,
		Common:  NewCommon(s, r),
	}

	if r.Method == "POST" {
		listing.Title = r.FormValue("title")
		listing.Body = r.FormValue("body")

		if err := listing.NormalizeAndValidate(); err != nil {
			resp.Err = err
		} else if err := persist(c, listing); err != nil {
			resp.Err = err
		} else {
			resp.Edit = false // hide the edit form if we're done
			resp.Redirect = true
		}
	}

	return resp, nil
}

func CreateListing(s *Subject, c context.Context, r *http.Request) (Template, error) {
	if !s.CanCreate() {
		return nil, AccessDenied
	}

	listing := &Listing{
		Seller:   s.Subject,
		Category: ParseCategory(r.FormValue("category")),
		Seeking:  r.FormValue("seeking") != "",
	}

	if err := persist(c, listing); err != nil {
		return nil, err
	}

	return SingleListing{
		Edit:     true,
		Redirect: true,
		Listing:  listing,
		Common:   NewCommon(s, r),
	}, nil
}
