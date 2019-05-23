package config

import (
	"io/ioutil"
	"testing"
)

func TestConfig(t *testing.T) {
	// Create project
	p, err := NewProject("test-project")
	if err != nil {
		t.Errorf("Can't create project: %v", err)
	}

	// Open configuration file
	f, err := p.Open("test-config")
	if err != nil {
		t.Errorf("Can't open configuration file: %v", err)
	}
	defer f.Close()

	// Write to file
	_, err = f.WriteString("Hello World")
	if err != nil {
		t.Errorf("Can't write to configuration file: %v", err)
	}

	// Check file content
	b, err := ioutil.ReadFile(p.FilePath("test-config"))
	if err != nil {
		t.Errorf("Can't read configuration file: %v", err)
	}

	// Test content
	if string(b) != "Hello World" {
		t.Errorf("Configuration file content does not match with Hello World")
	}

	// Cleanup by destroying project
	err = p.Destroy()
	if err != nil {
		t.Errorf("Can't destroy project: %v", err)
	}
}
