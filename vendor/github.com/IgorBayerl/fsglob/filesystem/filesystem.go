// Package filesystem provides abstractions for filesystem operations to enable
// easier testing and mocking of file system interactions.
package filesystem

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// Platformer provides a way for filesystem implementations to indicate
// what platform or environment they are simulating. This is particularly
// useful for mocks that need to simulate different operating systems.
type Platformer interface {
	// Platform returns a string identifier for the platform being simulated
	// (e.g., "windows", "linux", "darwin").
	Platform() string
}

// Filesystem defines the interface for filesystem operations that can be
// implemented by both real filesystem implementations and mocks for testing.
// It wraps common filesystem operations from the os and filepath packages.
type Filesystem interface {
	// Stat returns a FileInfo describing the named file.
	// If there is an error, it will be of type *PathError.
	Stat(name string) (fs.FileInfo, error)

	// ReadDir reads the named directory and returns a list of directory entries
	// sorted by filename.
	ReadDir(name string) ([]fs.DirEntry, error)

	// Getwd returns a rooted path name corresponding to the current directory.
	// If the current directory can be reached via multiple paths (due to
	// symbolic links), Getwd may return any one of them.
	Getwd() (string, error)

	// Abs returns an absolute representation of path. If the path is not
	// absolute it will be joined with the current working directory to turn
	// it into an absolute path.
	Abs(path string) (string, error)

	// MkdirAll creates a directory named path, along with any necessary parents,
	// and returns nil, or else returns an error. The permission bits perm are
	// used for all directories that MkdirAll creates.
	MkdirAll(path string, perm fs.FileMode) error

	// Create creates or truncates the named file. If the file already exists,
	// it is truncated. If the file does not exist, it is created with mode 0666
	// (before umask). If successful, methods on the returned file can be used
	// for I/O; the associated file descriptor has mode O_RDWR.
	Create(path string) (io.WriteCloser, error)

	// Open opens the named file for reading. If successful, methods on the
	// returned file can be used for reading; the associated file descriptor
	// has mode O_RDONLY.
	Open(path string) (fs.File, error)

	// ReadFile reads the named file and returns the contents. A successful
	// call returns err == nil, not err == EOF.
	ReadFile(path string) ([]byte, error)

	// WriteFile writes data to the named file, creating it if necessary.
	// If the file does not exist, WriteFile creates it with permissions perm
	// (before umask); otherwise WriteFile truncates it before writing, without
	// changing permissions.
	WriteFile(path string, data []byte, perm fs.FileMode) error
}

// DefaultFS implements the Filesystem interface using the standard `os` and `filepath` packages.
// It represents the real, underlying filesystem of the host operating system.
// This is the concrete implementation used in production code.
type DefaultFS struct{}

// Stat returns a FileInfo describing the named file using os.Stat.
func (DefaultFS) Stat(name string) (fs.FileInfo, error) { return os.Stat(name) }

// ReadDir reads the named directory using os.ReadDir.
func (DefaultFS) ReadDir(name string) ([]fs.DirEntry, error) { return os.ReadDir(name) }

// Getwd returns the current working directory using os.Getwd.
func (DefaultFS) Getwd() (string, error) { return os.Getwd() }

// Abs returns an absolute representation of path using filepath.Abs.
func (DefaultFS) Abs(path string) (string, error) { return filepath.Abs(path) }

// MkdirAll creates a directory and any necessary parents using os.MkdirAll.
func (DefaultFS) MkdirAll(path string, perm fs.FileMode) error { return os.MkdirAll(path, perm) }

// Create creates or truncates the named file using os.Create.
func (DefaultFS) Create(path string) (io.WriteCloser, error) { return os.Create(path) }

// Open opens the named file for reading using os.Open.
func (DefaultFS) Open(path string) (fs.File, error) { return os.Open(path) }

// ReadFile reads the named file and returns the contents using os.ReadFile.
func (DefaultFS) ReadFile(path string) ([]byte, error) { return os.ReadFile(path) }

// WriteFile writes data to the named file using os.WriteFile.
func (DefaultFS) WriteFile(path string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(path, data, perm)
}
