package clipper

const defaultCapacity = 5

// Clipper is used for clipping webpages into Notion.
type Clipper struct {
	token     string
	capacity  int
	clippings []Clipping
}

// New creates a new Clipper with a Notion authentication token.
// This token can be retrieved from the "token_v2" cookie created
// when authenticated from notion.so website.
func New(token string) *Clipper {
	return &Clipper{token: token, capacity: defaultCapacity}
}

// A Clipping is a single webpage item with a title and url.
type Clipping struct {
	title string
	url   string
}

// Capacity sets how many clippings are sent to Notion by http request.
// The default value is 5.
func (c *Clipper) Capacity(cap int) {
	c.capacity = cap
}

// Empty clears the list of previously loaded clippings.
func (c *Clipper) Empty() {
	c.clippings = nil
}

// Load some clippings into the clipper to prepare them to be clipped into Notion.
func (c *Clipper) Load(clippings ...Clipping) {
	for _, clipping := range clippings {
		c.clippings = append(c.clippings, clipping)
	}
}

// Save loaded clippings to Notion under the block that has
// the given blockID.
func (c *Clipper) Save(blockID string) error {
	return nil
}
