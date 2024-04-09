package source

import (
	"moria/config"
)

type Source interface {
	// Uploads the given file bytes to the source with the given filename.
	Upload(file []byte, filename string) error
	// Deletes the file with the given URL from the source.
	Delete(filename string) error
	// Returns the file bytes of the file with the given URL.
	Get(filename string) ([]byte, error)
}

var source Source

// Initializes the source package.
func Initialize() {
	switch config.GetSourceType() {
	case config.SourceTypeLocal:
		source = NewLocalSource()
	case config.SourceTypeRemoteFTP:
		source = NewFTPSource()
	case config.SourceTypeRemoteSFTP:
		source = NewFTPSource()
	case config.SourceTypeRemoteSSH:
		source = NewSSHSource()
	}
}

// Returns the source which is currently configured in the config.
func GetCurrent() Source {
	return source
}
