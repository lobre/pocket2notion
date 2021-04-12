package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/lobre/pocket2notion/config"
	"github.com/lobre/pocket2notion/notion/clipper"
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
	pocketDeleteOrg       bool

	notionTags      bool
	notionBatchSize int
	notionBlockID   string

	listOnly bool
}

func main() {
	// redefine usage
	flag.Usage = func() {
		usage()
	}

	var args arguments

	// filters
	flag.IntVar(&args.pocketCountFilter, "count", 0, "Number of Pocket items to import (newest clipped items first)")
	flag.BoolVar(&args.pocketFavoritedFilter, "favorited", false, "Only import favorited Pocket items")
	flag.BoolVar(&args.pocketArchivedFilter, "archived", false, "Only import archived Pocket items")
	flag.StringVar(&args.pocketTagFilter, "tag", "", "Only import Pocket items matching with tag")
	flag.StringVar(&args.pocketSearchFilter, "search", "", "Only import Pocket items matching with search")
	flag.IntVar(&args.pocketSinceFilter, "since", 0, "Only import Pocket items since a timestamp")
	flag.BoolVar(&args.pocketDeleteOrg, "delete", false, "Delete original in Pocket")

	flag.IntVar(&args.notionBatchSize, "notion-batch", 5, "Import into Notion by batch of <n> per http call")
	flag.BoolVar(&args.notionTags, "notion-tags", true, "Append Pocket tags to Notion by appending them to the item title with a hashtag")

	flag.BoolVar(&args.listOnly, "list-only", false, "Don't import into Notion but just list Pocket items (NOTION_BLOCK_ID not required with this flag)")

	flag.Parse()

	args.notionBlockID = flag.Arg(0)
	if args.notionBlockID == "" && !args.listOnly {
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

	err = pushToNotion(config, args, items)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func pushToNotion(config *config.Project, args arguments, items []api.Item) error {
	token, err := loadStringFromConfig(config.FilePath(notionTokenFile))
	if err != nil {
		return err
	}

	fmt.Println("List of items items")

	clip := clipper.New(token)

	// process by batch
	for i := 0; i < len(items); i += args.notionBatchSize {
		j := i + args.notionBatchSize
		if j > len(items) {
			j = len(items)
		}

		for _, item := range items[i:j] {
			url := item.ResolvedURL
			if url == "" {
				url = item.GivenURL
			}

			title := item.ResolvedTitle
			if title == "" {
				title = item.GivenTitle
			}

			fmt.Printf("> URL: %s\n", url)
			fmt.Printf("  Title: %s\n", title)

			notionItem := clipper.Item{
				Title: title,
				URL:   url,
			}

			if args.notionTags {
				var tags bytes.Buffer
				for tag := range item.Tags {
					tags.WriteString(fmt.Sprintf(" #%s", tag))
				}
				fTags := strings.TrimSpace(tags.String())
				fmt.Printf("  Tags: %v\n", fTags)

				notionItem.Title = fmt.Sprintf("%s %s", notionItem.Title, fTags)
			}

			// prepare item to be clipped
			if !args.listOnly {
				clip.Load(notionItem)
			}
		}

		// do the actual clipping request to Notion
		if !args.listOnly {
			fmt.Println("Pushing batch into Notion")
			err = clip.Save(args.notionBlockID)
			if err != nil {
				return err
			}
		}
	}

	fmt.Println("Terminated")
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
