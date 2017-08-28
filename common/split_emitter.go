package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type EmitterMap map[string]Emitter

type SplitEmitter struct {
	destDir          string
	basenameTemplate string
	emitters         EmitterMap
}

func NewSplitEmitter(filenameTemplate string) (Emitter, error) {
	if err := validateFilenameTemplate(filenameTemplate); err != nil {
		return nil, err
	}

	basedir := filepath.Dir(filenameTemplate)
	basename := filepath.Base(filenameTemplate)
	basename = basename[0:strings.LastIndex(basename, ".")]

	return &SplitEmitter{emitters: make(EmitterMap), destDir: basedir, basenameTemplate: basename}, nil
}

func (se *SplitEmitter) Receiver() EntityResult {
	return nil
}

func (se *SplitEmitter) Emit(entity EntityResult) error {
	entityType, ok := entity["$type"].(string)

	if !ok {
		return fmt.Errorf("Could not determine $type of entity %v", entity)
	}

	if emitter, err := se.findOrCreateEmitter(entityType); err == nil {
		return emitter.Emit(entity)
	} else {
		return err
	}
}

func (se *SplitEmitter) NextEmitter(current EntityResult, key string, isMultiValue bool) Emitter {
	return se
}

func (se *SplitEmitter) Finalize() error {
	for _, emitter := range se.emitters {
		if err := emitter.Finalize(); err != nil {
			return err
		}
	}

	return nil
}

func (se *SplitEmitter) findOrCreateEmitter(entityType string) (Emitter, error) {
	if "" == entityType {
		return nil, fmt.Errorf("Cannot generate emitter without entity type")
	}

	if emitter, isPresent := se.emitters[entityType]; isPresent {
		return emitter, nil
	}

	filename := filepath.Join(se.destDir, se.basenameTemplate+"-"+entityType+".json")
	if emitter, err := NewFlatEmitter(filename); err == nil {
		se.emitters[entityType] = emitter
		return emitter, nil
	} else {
		return nil, err
	}
}

func validateFilenameTemplate(filename string) error {
	if "" == filename {
		return fmt.Errorf("Split output requires a filename template")
	}

	downcased := strings.ToLower(filepath.Base(filename))

	if ".json" != filepath.Ext(downcased) {
		return fmt.Errorf("Split output filename template must have a `.json` extension")
	}

	if ".json" == downcased {
		return fmt.Errorf("split output filename must have a basename before the `.json` extension")
	}

	basedir := filepath.Dir(filename)

	dirExists, patherr := isDir(basedir)

	if patherr != nil {
		return fmt.Errorf("Failed to stat() %q: %v", basedir, patherr)
	}

	if !dirExists {
		return fmt.Errorf("Directory %q is not a directory", basedir)
	}

	return nil
}

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
