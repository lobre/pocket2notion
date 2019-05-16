package main

import (
	"flag"
	"fmt"
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
	fmt.Println(args)
}
