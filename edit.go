package ironclad

import (
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
	listing := &Listing{
		Seller:   "bob@uchicago.edu",
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
