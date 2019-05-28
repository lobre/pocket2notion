# pocket2notion

Currently, there is no way to export Pocket items into Notion because Notion does not provide any API so far.
This projects aims to give a way to append web pages to Notion using the lately introduced [Web Clipper](https://www.notion.so/Web-Clipper-ba54b19ecaeb466b8070b9e683c5fce1).

In order to do this, a `clipper` golang package has been implemented to serve as a client for interacting with the Notion Web Clipper. See below for more information.

Currently, the Notion Web Clipper does not give us a way to provide more information than just a title and a URL. So we have no way to migrate tags from Pocket to Notion. In order to still pass them, I decided that tags will be appended to the title with a hashtag before. You can disable this feature with the `--notion-tags=false` flag.

## Install

    go get github.com/lobre/pocket2notion/cmd/pocket2notion

## Usage

### 1. Get a Pocket consumer key 

Go to https://getpocket.com/developer/apps/, create an app and gather the consumer key.

### 2. Get your Notion authentication token

As there is no official Notion API, you will need to fetch the value of an authenticated token that you can find using Chrome Developer Tools while browsing authenticated on notion.so. You need to gather the content of the `token_v2` cookie.

### 3. Add tokens to configuration

    mkdir -p ~/.config/pocket2notion
    echo "MY_POCKET_CONSUMER_KEY" > ~/.config/pocket2notion/pocket_consumer_key
    echo "MY_NOTION_TOKEN" > ~/.config/pocket2notion/notion_token

### 4. Run pocket2notion

    ❯ pocket2notion -h
    Usage of ./pocket2notion:
    -archived
            Only import archived Pocket items
    -count int
            Number of Pocket items to import
    -favorited
            Only import favorited Pocket items
    -httptest.serve string
            if non-empty, httptest.NewServer serves on this address and blocks
    -notion-tags
            Append Pocket tags to Notion by appending them to the item title with a hashtag (default true)
    -search string
            Only import Pocket items matching with search
    -since int
            Only import Pocket items since a timestamp
    -tag string
            Only import Pocket items matching with tag

### Examples

    pocket2notion --count 1
    pocket2notion --count 5 --tag=politics

## Notion clipper package

The `clipper` package provides bindings to reproduce the HTTP request made by the clipper extension to add articles to Notion.

### **Disclaimer**
This package has been written by analysing the Chrome extension XHR call. This service exposed by Notion is not public and so there are no guarantees that parameters won't change.
There are even good chances that it will evoluate as Notion might add new features to the clipper in the future.

I have implemented this package to help me doing the switch from Pocket to Notion. I don't aim to keep it updated and to provide full bindings to the Notion service. But feel free to improve/fix it if it turns broken. 

Currently, the package does not support the "+New links database" option as in the extension popup. You can only add items to an already existing database or page. You will by the way need the blockId of this database/page. You can usually find it in the URL of the page when browsing it from the web version of Notion. The package does not either give a way to indicate the URL property that you want to use in your existing database in Notion. It will default to choose the first property of type `link`. 

You can find more information in the `clipper` package's Golang documentation.
