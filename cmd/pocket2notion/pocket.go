package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"

	"github.com/lobre/pocket2notion/config"
	"github.com/pkg/errors"

	"github.com/motemen/go-pocket/api"
	"github.com/motemen/go-pocket/auth"
)

const pocketConsumerKeyFile = "pocket_consumer_key"
const pocketAccessTokenFile = "pocket_auth.json"

func retrievePocketItems(config *config.Project, args arguments) ([]api.Item, error) {
	fmt.Println("Pocket authentication...")
	client, err := newPocketClient(config)
	if err != nil {
		return nil, errors.Wrap(err, "can't instanciate Pocket client")
	}

	options := &api.RetrieveOption{
		Sort:       api.SortNewest,
		DetailType: api.DetailTypeComplete,
	}

	if args.pocketCountFilter != 0 {
		options.Count = args.pocketCountFilter
	}

	if args.pocketFavoritedFilter {
		options.Favorite = api.FavoriteFilterFavorited
	}

	if args.pocketArchivedFilter {
		options.State = api.StateArchive
	}

	if args.pocketTagFilter != "" {
		options.Tag = args.pocketTagFilter
	}

	if args.pocketSearchFilter != "" {
		options.Search = args.pocketSearchFilter
	}

	if args.pocketSinceFilter != 0 {
		options.Since = args.pocketSinceFilter
	}

	fmt.Println("Fetching Pocket items...")
	res, err := client.Retrieve(options)
	if err != nil {
		return nil, errors.Wrap(err, "can't retrieve Pocket list")
	}

	// convert into slice
	items := []api.Item{}
	for _, item := range res.List {
		items = append(items, item)
	}

	sort.Sort(bySortID(items))

	if args.pocketDeleteOrg {
		for _, item := range items {
			delAct := &api.Action {
				Action: "delete",
				ItemID: item.ItemID,
			}
			client.Modify(delAct)
		}
	}

	return items, nil
}

func newPocketClient(config *config.Project) (*api.Client, error) {
	consumerKey, err := loadStringFromConfig(config.FilePath(pocketConsumerKeyFile))
	if err != nil {
		return nil, errors.Wrap(err, "can't get pocket consumer key from config")
	}

	accessToken, err := restorePocketAccessToken(config, consumerKey)
	if err != nil {
		return nil, errors.Wrap(err, "can't get pocket access token")
	}

	return api.NewClient(consumerKey, accessToken.AccessToken), nil
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

	redirectURL := ts.URL

	requestToken, err := auth.ObtainRequestToken(consumerKey, redirectURL)
	if err != nil {
		return nil, err
	}

	url := auth.GenerateAuthorizationURL(requestToken, redirectURL)
	fmt.Println(url)

	<-ch

	ts.Close()

	return auth.ObtainAccessToken(consumerKey, requestToken)
}

type bySortID []api.Item

func (s bySortID) Len() int           { return len(s) }
func (s bySortID) Less(i, j int) bool { return s[i].SortId < s[j].SortId }
func (s bySortID) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
