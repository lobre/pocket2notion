package clipper

import (
	"os"
	"testing"
)

const tokenEnv = "NOTION_TOKEN"
const parentBlockEnv = "NOTION_PARENT_BLOCK"

func TestClipping(t *testing.T) {
	token, ok := os.LookupEnv(tokenEnv)
	if !ok {
		t.Fatalf("the Notion token is not defined")
	}

	parentBlock, ok := os.LookupEnv(parentBlockEnv)
	if !ok {
		t.Fatalf("The Notion parent block is not defined")
	}

	c := New(token)

	// Prepare some items to be clipped
	c.Load(
		Item{Title: "My title", URL: "www.google.com"},
		Item{Title: "My second title", URL: "www.twitter.com"},
	)

	// Do the actual clipping request to Notion
	err := c.Save(parentBlock)
	if err != nil {
		t.Fatalf("could not clip: %v", err)
	}
}
