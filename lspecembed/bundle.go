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

	mu      sync.Mutex
	tempDir string
	root    string
	rootErr error
	refs    int
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

// BundleAcquire materializes the bundle when needed and returns the root .lspec path.
// Call release when finished so the temporary directory is removed.
func BundleAcquire(bundle *Bundle) (rootPath string, release func(), err error) {
	if bundle == nil {
		return "", func() {}, fmt.Errorf("lspecembed: nil bundle")
	}

	bundle.mu.Lock()
	defer bundle.mu.Unlock()

	if bundle.refs == 0 {
		bundle.root, bundle.tempDir, bundle.rootErr = materializeBundle(bundle)
	}
	if bundle.rootErr != nil {
		return "", func() {}, bundle.rootErr
	}

	bundle.refs++
	path := bundle.root
	return path, func() { bundleRelease(bundle) }, nil
}

// BundleRootPath is equivalent to BundleAcquire without releasing the temp directory.
// Prefer BundleAcquire and call release to avoid leaving materialized trees in the system temp folder.
func BundleRootPath(bundle *Bundle) (string, error) {
	path, _, err := BundleAcquire(bundle)
	return path, err
}

// ResolvePath returns the env override when set and readable, otherwise acquires the bundle.
// When the embedded bundle is materialized, call release after the path is no longer needed.
func ResolvePath(envVar string, bundle *Bundle) (path string, release func(), err error) {
	if envVar != "" {
		if path := os.Getenv(envVar); path != "" {
			if _, err := os.Stat(path); err != nil {
				return "", func() {}, fmt.Errorf("%s not usable (%q): %w", envVar, path, err)
			}
			return path, func() {}, nil
		}
	}
	return BundleAcquire(bundle)
}

func bundleRelease(bundle *Bundle) {
	bundle.mu.Lock()
	defer bundle.mu.Unlock()

	if bundle.refs <= 0 {
		return
	}
	bundle.refs--
	if bundle.refs > 0 {
		return
	}

	if bundle.tempDir != "" {
		_ = os.RemoveAll(bundle.tempDir)
	}
	bundle.tempDir = ""
	bundle.root = ""
	bundle.rootErr = nil
}

func materializeBundle(bundle *Bundle) (rootPath string, tempDir string, err error) {
	tempDir, err = os.MkdirTemp("", "lingua-lspec-*")
	if err != nil {
		return "", "", fmt.Errorf("lspecembed: create temp dir: %w", err)
	}

	for _, file := range bundle.Files {
		absPath := filepath.Join(tempDir, filepath.FromSlash(file.RelPath))
		if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
			_ = os.RemoveAll(tempDir)
			return "", "", fmt.Errorf("lspecembed: mkdir %q: %w", file.RelPath, err)
		}
		if err := os.WriteFile(absPath, file.Data, 0o644); err != nil {
			_ = os.RemoveAll(tempDir)
			return "", "", fmt.Errorf("lspecembed: write %q: %w", file.RelPath, err)
		}
	}

	return filepath.Join(tempDir, filepath.FromSlash(bundle.RootRelPath)), tempDir, nil
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
