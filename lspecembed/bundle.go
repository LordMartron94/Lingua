package lspecembed

import (
	"fmt"
	"lingua/lspec_std"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const stdRelPath = "lspec_std/std.lspec"

// File is one embedded .lspec (or dependency) relative to the materialization root.
type File struct {
	RelPath string
	Data    []byte
}

// Bundle is an embedded spec tree rooted at RootRelPath.
type Bundle struct {
	RootRelPath string
	Files       []File

	once    sync.Once
	root    string
	rootErr error
}

// BundleCreate validates and returns a bundle. Panics on invalid input.
func BundleCreate(rootRelPath string, files ...File) *Bundle {
	rootRelPath = normalizeRelPath(rootRelPath)
	if !strings.HasSuffix(rootRelPath, ".lspec") {
		panic(fmt.Sprintf("lspecembed: root %q must end with .lspec", rootRelPath))
	}

	seen := make(map[string]struct{}, len(files)+1)
	deduped := make([]File, 0, len(files))

	for _, file := range files {
		rel := normalizeRelPath(file.RelPath)
		validateRelPath(rel)
		if _, exists := seen[rel]; exists {
			panic(fmt.Sprintf("lspecembed: duplicate file %q", rel))
		}
		seen[rel] = struct{}{}
		deduped = append(deduped, File{RelPath: rel, Data: file.Data})
	}

	if _, ok := seen[rootRelPath]; !ok {
		panic(fmt.Sprintf("lspecembed: root %q not present in files", rootRelPath))
	}

	return &Bundle{
		RootRelPath: rootRelPath,
		Files:       deduped,
	}
}

// StdFile returns the shared lspec_std dependency for specs that import it.
func StdFile() File {
	return File{
		RelPath: stdRelPath,
		Data:    lspec_std.EmbeddedStdLspec,
	}
}

// BundleRootPath materializes the bundle once per process and returns the root .lspec path.
func BundleRootPath(bundle *Bundle) (string, error) {
	if bundle == nil {
		return "", fmt.Errorf("lspecembed: nil bundle")
	}
	bundle.once.Do(func() {
		bundle.root, bundle.rootErr = materializeBundle(bundle)
	})
	return bundle.root, bundle.rootErr
}

// ResolvePath returns the env override when set and readable, otherwise BundleRootPath.
func ResolvePath(envVar string, bundle *Bundle) (string, error) {
	if envVar != "" {
		if path := os.Getenv(envVar); path != "" {
			if _, err := os.Stat(path); err != nil {
				return "", fmt.Errorf("%s not usable (%q): %w", envVar, path, err)
			}
			return path, nil
		}
	}
	return BundleRootPath(bundle)
}

func materializeBundle(bundle *Bundle) (string, error) {
	root, err := os.MkdirTemp("", "lingua-lspec-*")
	if err != nil {
		return "", fmt.Errorf("lspecembed: create temp dir: %w", err)
	}

	for _, file := range bundle.Files {
		absPath := filepath.Join(root, filepath.FromSlash(file.RelPath))
		if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
			return "", fmt.Errorf("lspecembed: mkdir %q: %w", file.RelPath, err)
		}
		if err := os.WriteFile(absPath, file.Data, 0o644); err != nil {
			return "", fmt.Errorf("lspecembed: write %q: %w", file.RelPath, err)
		}
	}

	return filepath.Join(root, filepath.FromSlash(bundle.RootRelPath)), nil
}

func normalizeRelPath(path string) string {
	return filepath.ToSlash(filepath.Clean(path))
}

func validateRelPath(path string) {
	if filepath.IsAbs(path) {
		panic(fmt.Sprintf("lspecembed: absolute path %q not allowed", path))
	}
	if path == ".." || strings.HasPrefix(path, "../") || strings.Contains(path, "/../") {
		panic(fmt.Sprintf("lspecembed: path %q must not contain ..", path))
	}
}
