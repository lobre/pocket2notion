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

	// 5 items clipped at the same time per request
	c.Capacity(5)

	// Clean if already loaded items
	c.Empty()

	// Prepare some items to be clipped
	c.Load(
		Clipping{title: "My title", url: "www.google.com"},
		Clipping{title: "My second title", url: "www.twitter.com"},
	)

	// Do the actual clipping request to Notion
	err := c.Save(parentBlock)
	if err != nil {
		t.Fatalf("could not clip: %v", err)
	}
}
