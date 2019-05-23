// Package config provides a simple way of interacting
// with configuration files of a project that are situated under
// $HOME/.config/<project>.
//
// This package only handles one level of configuration files under the
// project directory. It does not support multi directory levels.
//
// Here is a example snapshot of the filesystem of a project called "totem".
// ~❯ tree .config/totem
// .config/totem
// ├── config
// ├── keys.json
// ├── state.ini
// └── users.conf
//
// The default configuration file is named "config".
package config

import (
	"os"
	"os/user"
	"path/filepath"
)

const configDir = ".config"
const configDirPerm os.FileMode = 755
const defaultFile = "config"

// Project contains the information needed to compute paths
// for a project configuration files located under $HOME/.config/<Name>.
type Project struct {
	Name string

	usr *user.User
}

// NewProject creates a project and makes sure the corresponding
// folder exist on the filesystem at $HOME/.config/<name>.
func NewProject(name string) (*Project, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	p := Project{name, usr}

	// Create config folder if it does not exist
	err = os.MkdirAll(p.Path(), configDirPerm)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// FilePath returns the path of a configuration file
// under the current project.
func (p *Project) FilePath(name string) string {
	return ""
}

// FilePathDefault returns the path of a the defualt configuration file
// of the current project.
func (p *Project) FilePathDefault() string {
	return ""
}

// Open returns an os.File from a name corresponding to a
// configuration file under the current project's path.
// Warning, the file should be closed after usage.
// If the file does not exist, it will be created first.
func (p *Project) Open(name string) (*os.File, error) {
	return nil, nil
}

// OpenDefault returns an os.File from a name corresponding to a
// the default configuration of the current project.
// Warning, the file should be closed after usage.
// If the file does not exist, it will be created first.
func (p *Project) OpenDefault() (*os.File, error) {
	return nil, nil
}

// Remove deletes a configuration file from the current project.
func (p *Project) Remove(name string) error {
	return nil
}

// RemoveDefault deletes the default configuration file from the current project.
func (p *Project) RemoveDefault() error {
	return nil
}

// Truncate empties a configuration file from the current project.
func (p *Project) Truncate(name string) error {
	return nil
}

// TruncateDefault empties the default configuration file from the current project.
func (p *Project) TruncateDefault() error {
	return nil
}

// Destroy removes the project configuration folder
// and all the files it contains.
func (p *Project) Destroy() error {
	return os.RemoveAll(p.Path())
}

// Path return the path of the folder containing configuration files.
func (p *Project) Path() string {
	return filepath.Join(p.usr.HomeDir, configDir, p.Name)
}
