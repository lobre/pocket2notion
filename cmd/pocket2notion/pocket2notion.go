package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lobre/pocket2notion/config"
)

type arguments struct {
	pocketCountFilter     int
	pocketFavoritedFilter bool
	pocketArchivedFilter  bool
	pocketTagFilter       string
	pocketSearchFilter    string
	pocketSinceFilter     int

	notionTags bool
}

func main() {
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

	// list items
	for _, item := range items {
		fmt.Printf("URL: %s\n", item.GivenURL)
		fmt.Printf("Title: %s\n", item.GivenTitle)

		tags := []string{}
		for tag, _ := range item.Tags {
			tags = append(tags, tag)
		}
		fmt.Printf("Tags: %v\n", tags)
	}
}
