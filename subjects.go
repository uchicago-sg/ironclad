package ironclad

import (
	"fmt"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/net/context"
)

// generating and removing subjects
var (
	hmacSecret []byte
	hmacOnce   sync.Once
)

type Subject struct {
	Name string `json:"cn"`
	jwt.StandardClaims
}

func ParseSubject(c context.Context, tokenString string) (*Subject, error) {
	hmacOnce.Do(func() {
		hmacSecret = sharedSecret(c)
	})

	if tokenString == "" {
		return nil, nil
	}

	subject := &Subject{}
	_, err := jwt.ParseWithClaims(tokenString, subject,
		func(token *jwt.Token) (interface{}, error) {
			// important! otherwise bad things happen
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v",
					token.Header["alg"])
			}
			return hmacSecret, nil
		})

	if err != nil {
		return nil, err
	} else if err := subject.Valid(); err != nil {
		return nil, err
	} else {
		return subject, nil
	}
}

func (s *Subject) Serialize(c context.Context) (string, error) {
	hmacOnce.Do(func() {
		hmacSecret = sharedSecret(c)
	})

	if s == nil {
		return "", nil
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, s).SignedString(hmacSecret)
}

func newSubject(name, email string) *Subject {
	return &Subject{
		Name: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
			Subject:   email,
		},
	}
}
