package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/lobre/pocket2notion/config"
	"github.com/motemen/go-pocket/api"
)

const notionTokenFile = "notion_token"

type arguments struct {
	pocketCountFilter     int
	pocketFavoritedFilter bool
	pocketArchivedFilter  bool
	pocketTagFilter       string
	pocketSearchFilter    string
	pocketSinceFilter     int

	notionTags    bool
	notionBlockID string
}

func main() {
	// redefine usage
	flag.Usage = func() {
		usage()
	}

	var args arguments

	// filters
	flag.IntVar(&args.pocketCountFilter, "count", 0, "Number of Pocket items to import")
	flag.BoolVar(&args.pocketFavoritedFilter, "favorited", false, "Only import favorited Pocket items")
	flag.BoolVar(&args.pocketArchivedFilter, "archived", false, "Only import archived Pocket items")
	flag.StringVar(&args.pocketTagFilter, "tag", "", "Only import Pocket items matching with tag")
	flag.StringVar(&args.pocketSearchFilter, "search", "", "Only import Pocket items matching with search")
	flag.IntVar(&args.pocketSinceFilter, "since", 0, "Only import Pocket items since a timestamp")

	flag.BoolVar(&args.notionTags, "notion-tags", true, "Append Pocket tags to Notion by appending them to the item title with a hashtag")

	flag.Parse()

	args.notionBlockID = flag.Arg(0)
	if args.notionBlockID == "" {
		usage()
		os.Exit(1)
	}

	// init config project
	config, err := config.NewProject("pocket2notion")
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	items, err := retrievePocketItems(config, args)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	err = pushToNotion(config, items, args.notionBlockID)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func pushToNotion(config *config.Project, items []api.Item, blockID string) error {
	token, err := loadStringFromConfig(config.FilePath(notionTokenFile))
	if err != nil {
		return err
	}
	fmt.Println(token)

	// c := clipper.New(token)

	// list items
	for _, item := range items {
		fmt.Printf("URL: %s\n", item.GivenURL)
		fmt.Printf("Title: %s\n", item.GivenTitle)

		var tags bytes.Buffer
		for tag := range item.Tags {
			tags.WriteString(fmt.Sprintf(" #%s", tag))
		}
		fmt.Printf("Tags: %v\n", tags.String())
	}

	// Prepare some items to be clipped
	// c.Load(
	// 	clipper.Item{Title: "My title", URL: "www.google.com"},
	// 	clipper.Item{Title: "My second title", URL: "www.twitter.com"},
	// )

	// Do the actual clipping request to Notion
	// err = c.Save(blockID)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func loadStringFromConfig(path string) (string, error) {
	consumerKey, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bytes.SplitN(consumerKey, []byte("\n"), 2)[0]), nil
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] NOTION_BLOCK_ID\n", os.Args[0])
	flag.PrintDefaults()
}
