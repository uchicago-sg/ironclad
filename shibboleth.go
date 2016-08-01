package ironclad

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"strings"
	"sync"

	"github.com/russellhaering/gosaml2"
	"golang.org/x/net/context"
)

var (
	samlProvider = &saml2.SAMLServiceProvider{
		IdentityProviderSSOURL:      "https://shibboleth2.uchicago.edu/idp/profile/SAML2/Redirect/SSO",
		IdentityProviderIssuer:      "https://marketplace.uchicago.edu",
		AssertionConsumerServiceURL: "https://go-marketplace.appspot.com/saml-callback",
		SignAuthnRequests:           false,
		AudienceURI:                 "",
	}
	samlOnce sync.Once
)

type SAMLConfig struct {
	Roots []*x509.Certificate

	Certificate *x509.Certificate
	PrivateKey  *rsa.PrivateKey
}

func (s *SAMLConfig) TrustCertificate(data []byte) error {
	block, _ := pem.Decode([]byte(data))
	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}

	s.Roots = append(s.Roots, certificate)
	return nil
}

func (s *SAMLConfig) Use(certificate, key []byte) error {
	block, _ := pem.Decode(certificate)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}

	block, _ = pem.Decode(key)
	pk, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}

	s.Certificate, s.PrivateKey = cert, pk

	return nil
}

func conf(c context.Context, key string) []byte {
	b, err := getConfig(c, key)
	if err != nil {
		panic(err)
	}
	return b
}

func initSAML(c context.Context) {
	samlOnce.Do(func() {
		s := &SAMLConfig{}

		if err := s.TrustCertificate(conf(c, "uchicago.crt")); err != nil {
			panic(err)
		}

		if err := s.Use(conf(c, "marketplace.crt"), conf(c, "marketplace.key")); err != nil {
			panic(err)
		}

		samlProvider.IDPCertificateStore = s
		samlProvider.SPKeyStore = s
	})

	if samlProvider.IDPCertificateStore == nil {
		panic("no SAML IDP :(")
	}
}

func (s SAMLConfig) GetKeyPair() (*rsa.PrivateKey, []byte, error) {
	return s.PrivateKey, []byte(s.Certificate.Raw), nil
}

func (s SAMLConfig) Certificates() ([]*x509.Certificate, error) {
	return s.Roots, nil
}

type LoginResult struct {
	NeedsShib  bool
	NewSubject *Subject
}

func (l LoginResult) Template() string  { return "" }
func (l LoginResult) Subject() *Subject { return l.NewSubject }

func (l LoginResult) NewURL() string {
	if l.NeedsShib {
		url, err := samlProvider.BuildAuthURL("helloworld")
		if err != nil {
			panic(err)
		}
		return url
	}
	return "/"
}

func LoginPage(s *Subject, c context.Context, r *http.Request) (Template, error) {
	initSAML(c)

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

func SAMLRedirect(s *Subject, c context.Context, r *http.Request) (Template, error) {
	initSAML(c)

	info, err := samlProvider.RetrieveAssertionInfo(r.FormValue("SAMLResponse"))
	if err != nil {
		return nil, err
	}

	// extract name and email from request
	displayName := string(info.Values["urn:oid:2.16.840.1.113730.3.1.241"].Values[0])
	mail := string(info.Values["urn:oid:0.9.2342.19200300.100.1.3"].Values[0])

	s = newSubject(displayName, mail)

	return &LoginResult{
		NewSubject: s,
	}, nil
}

func LogoutPage(s *Subject, c context.Context, r *http.Request) (Template, error) {
	return &LoginResult{
		NewSubject: nil,
	}, nil
}
