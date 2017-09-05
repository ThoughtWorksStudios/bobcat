package emitter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type EmitterProvider interface {
	Get(key string) (Emitter, error)
}

type PerTypeEmitterProvider struct {
	basedir  string
	basename string
}

/**
 * Creates a FlatEmitter for a given entity type, backed by a file derived from
 * the entity type name.
 */
func (p *PerTypeEmitterProvider) Get(entityType string) (Emitter, error) {
	if emitter, err := FlatEmitterForFile(p.PathFromType(entityType)); err != nil {
		return nil, err
	} else {
		emitter.Init()
		return emitter, nil
	}
}

/** Constructor */
func NewPerTypeEmitterProvider(filenameTemplate string) (EmitterProvider, error) {
	p := &PerTypeEmitterProvider{}

	if err := p.validateFilenameTemplate(filenameTemplate); err != nil {
		return nil, err
	}

	_bn := filepath.Base(filenameTemplate)

	p.basedir = filepath.Dir(filenameTemplate)
	p.basename = _bn[0:strings.LastIndex(_bn, ".")]

	return p, nil
}

/**
 * Derives a filename (with path) from an entity type, based on the filename
 * template, by adding `-{_type}` to the filename before the `.json` extension.
 *
 * e.g. basedir/basename-_type.json
 */
func (p *PerTypeEmitterProvider) PathFromType(entityType string) string {
	return filepath.Join(p.basedir, p.basename+"-"+entityType+".json")
}

/**
 * Validates that the template is of the correct format (must end in .json,
 * and has a non-zero basename with optional dirname), and that the enclosing
 * directory exists.
 */
func (p *PerTypeEmitterProvider) validateFilenameTemplate(filenameTemplate string) error {
	if "" == filenameTemplate {
		return fmt.Errorf("You must provide a filename template, e.g. path/to/file.json")
	}

	downcased := strings.ToLower(filepath.Base(filenameTemplate))

	if ".json" != filepath.Ext(downcased) {
		return fmt.Errorf("Filename template must have a `.json` extension")
	}

	if ".json" == downcased {
		return fmt.Errorf("Filename template must have a basename before the `.json` extension")
	}

	basedir := filepath.Dir(filenameTemplate)

	dirExists, patherr := isDir(basedir)

	if patherr != nil {
		return fmt.Errorf("Failed to stat() %q: %v", basedir, patherr)
	}

	if !dirExists {
		return fmt.Errorf("Directory %q is not a directory", basedir)
	}

	return nil
}

/**
 * Tests if the path exists and is a directory
 */
func isDir(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err == nil {
		return stat.IsDir(), nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return (nil != stat && stat.IsDir()), err
}
