package interpreter

import (
	fs "path/filepath"
)

/**
 * resolves a path, optionally relative to a basepath, to a canonical
 * path on the file system
 */
func resolve(path, basepath string) (string, error) {
	if fs.IsAbs(path) || basepath == "" {
		return canonical(path)
	}
	return canonical(fs.Join(basepath, path))
}

/**
 * Finds the true file system path
 */
func canonical(path string) (string, error) {
	if p1, e := fs.EvalSymlinks(path); e == nil {
		return fs.Abs(p1)
	} else {
		return "", e
	}
}

/**
 * Tries to resolve a basepath relative to the current directory
 * and respects symlinks, but doesn't try to resolve them (i.e. canonicalize).
 * That is to say, basedir() does not verify if the path really exists.
 */
func basedir(path, basepath string) (string, error) {
	if fs.IsAbs(path) || basepath == "" {
		return fs.Dir(path), nil
	}

	return fs.Rel(".", fs.Join(basepath, fs.Dir(path)))
}
