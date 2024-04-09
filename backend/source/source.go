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

var sources map[config.SourceType]Source

// Initializes the source package.
func Initialize() {
	registerSource(config.SourceTypeLocal, NewLocalSource())
}

// Returns the source which is currently configured in the config.
func GetCurrent() Source {
	return sources[config.GetSourceType()]
}

// Registers the given source with the given type.
func registerSource(sourceType config.SourceType, source Source) {
	sources[sourceType] = source
}
