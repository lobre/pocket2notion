package clipper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const url = "https://www.notion.so/api/v3/addWebClipperURLs"
const defaultCapacity = 5

// Clipper is used for clipping webpages into Notion.
type Clipper struct {
	token    string
	capacity int
	items    []Item
}

// New creates a new Clipper with a Notion authentication token.
// This token can be retrieved from the "token_v2" cookie created
// when authenticated from notion.so website.
func New(token string) *Clipper {
	return &Clipper{token: token, capacity: defaultCapacity}
}

// A Item is a single webpage item with a title and url.
type Item struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type payload struct {
	DataType string `json:"type"`
	BlockID  string `json:"blockId"`
	Items    []Item `json:"items"`
	From     string `json:"from"`
}

// Capacity sets how many clippings are sent to Notion by http request.
// The default value is 5.
func (c *Clipper) Capacity(cap int) {
	c.capacity = cap
}

// Empty clears the list of previously loaded clippings.
func (c *Clipper) Empty() {
	c.items = nil
}

// Load some clippings into the clipper to prepare them to be clipped into Notion.
func (c *Clipper) Load(clippings ...Item) {
	for _, clipping := range clippings {
		c.items = append(c.items, clipping)
	}
}

// Save loaded clippings to Notion under the block that has
// the given blockID.
func (c *Clipper) Save(blockID string) error {

	// TODO we need to handle dashes in blockID

	payload := payload{
		DataType: "block",
		BlockID:  blockID,
		Items:    c.items,
		From:     "chrome",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	var client = &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	// Set authentication cookie
	req.Header.Set("Cookie", fmt.Sprintf("token_v2=%s;", c.token))

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("wrong return value from Notion: %d", res.StatusCode)
	}

	return nil
}
