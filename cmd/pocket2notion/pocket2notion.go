package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lobre/pocket2notion/config"
)

type arguments struct {
	notionKey string
	pocketKey string
	limit     int
}

func main() {
	var args arguments
	flag.IntVar(&args.limit, "limit", -1, "Limit of Pocket items to import")
	flag.StringVar(&args.notionKey, "notionKey", "", "Limit of Pocket items to import")
	flag.StringVar(&args.pocketKey, "pocketKey", "", "Limit of Pocket items to import")
	flag.Parse()

	// Init config project
	config, err := config.NewProject("pocket2notion")
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	err = listPocketItems(config)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
