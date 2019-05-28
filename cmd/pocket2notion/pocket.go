package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/lobre/pocket2notion/config"
	"github.com/pkg/errors"

	"github.com/motemen/go-pocket/api"
	"github.com/motemen/go-pocket/auth"
)

const pocketConsumerKeyFile = "pocket_consumer_key"
const pocketAccessTokenFile = "pocket_auth.json"

func listPocketItems(config *config.Project) error {
	client, err := newPocketClient(config)
	if err != nil {
		return errors.Wrap(err, "can't instanciate Pocket client")
	}

	options := &api.RetrieveOption{}
	res, err := client.Retrieve(options)
	if err != nil {
		return errors.Wrap(err, "can't retrieve Pocket list")
	}

	fmt.Println(res)

	return nil
}

func newPocketClient(config *config.Project) (*api.Client, error) {
	consumerKey, err := getPocketConsumerKey(config)
	if err != nil {
		return nil, errors.Wrap(err, "can't get pocket consumer key from config")
	}

	accessToken, err := restorePocketAccessToken(config, consumerKey)
	if err != nil {
		return nil, errors.Wrap(err, "can't get pocket access token")
	}

	return api.NewClient(consumerKey, accessToken.AccessToken), nil
}

func getPocketConsumerKey(config *config.Project) (string, error) {
	consumerKey, err := ioutil.ReadFile(config.FilePath(pocketConsumerKeyFile))
	if err != nil {
		return "", err
	}
	return string(bytes.SplitN(consumerKey, []byte("\n"), 2)[0]), nil
}

func restorePocketAccessToken(config *config.Project, consumerKey string) (*auth.Authorization, error) {
	accessToken := &auth.Authorization{}

	// Try to load access token from config
	file, err := config.Open(pocketAccessTokenFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(accessToken)

	if err != nil {
		log.Println(err)

		accessToken, err = obtainPocketAccessToken(consumerKey)
		if err != nil {
			return nil, err
		}

		// Save token to config
		err = json.NewEncoder(file).Encode(accessToken)
		if err != nil {
			return nil, err
		}
	}

	return accessToken, nil
}

func obtainPocketAccessToken(consumerKey string) (*auth.Authorization, error) {
	ch := make(chan struct{})
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/favicon.ico" {
				http.Error(w, "Not Found", 404)
				return
			}

			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintln(w, "Authorized.")
			ch <- struct{}{}
		}))
	defer ts.Close()

	redirectURL := ts.URL

	requestToken, err := auth.ObtainRequestToken(consumerKey, redirectURL)
	if err != nil {
		return nil, err
	}

	url := auth.GenerateAuthorizationURL(requestToken, redirectURL)
	fmt.Println(url)

	<-ch

	return auth.ObtainAccessToken(consumerKey, requestToken)
}
