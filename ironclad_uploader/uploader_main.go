package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/remote_api"
)

type Config struct {
	Value []byte
}

func main() {
	const host = "ironclad-dot-go-marketplace.appspot.com"

	c, err := google.DefaultClient(context.Background(),
		"https://www.googleapis.com/auth/appengine.apis",
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/cloud-platform",
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx, err := remote_api.NewRemoteContext(host, c)
	if err != nil {
		log.Fatal(err)
	}

	upload := map[string]bool{}

	// see what's already there
	keys, err := datastore.NewQuery("Config").KeysOnly().GetAll(ctx, nil)
	if err != nil {
		log.Fatalf("read config: %s", err)
	}

	for _, key := range keys {
		upload[key.StringID()] = false
	}

	// see what we need to upload
	files, err := ioutil.ReadDir("config")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		upload[file.Name()] = true
	}

	for key, needed := range upload {
		if needed {
			contents, err := ioutil.ReadFile("config/" + key)
			if err != nil {
				log.Fatalf("%s: %s", key, err.Error())
			}

			secret := &Config{Value: contents}
			k := datastore.NewKey(ctx, "Config", key, 0, nil)
			if _, err := datastore.Put(ctx, k, secret); err != nil {
				log.Fatalf("%s: %s", key, err.Error())
			}
			fmt.Printf("%s: uploaded %d bytes\n", key, len(contents))
		} else {
			fmt.Printf("%s: (keeping existing value)\n", key)
		}
	}
}
