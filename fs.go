package fs

import (
	"errors"
	"syscall"
	"time"
)

//    Flags to Open wrapping those of the underlying system. Not all flags may
//    be implemented on a given system.
const (
	O_RDONLY int = syscall.O_RDONLY // open the file read-only.
	O_WRONLY int = syscall.O_WRONLY // open the file write-only.
	O_RDWR   int = syscall.O_RDWR   // open the file read-write.
	O_APPEND int = syscall.O_APPEND // append data to the file when writing.
	O_CREATE int = syscall.O_CREAT  // create a new file if none exists.
	O_EXCL   int = syscall.O_EXCL   // used with O_CREATE, file must not exist
	O_SYNC   int = syscall.O_SYNC   // open for synchronous I/O.
	O_TRUNC  int = syscall.O_TRUNC  // if possible, truncate file when opened.
)

//     Seek whence values.
const (
	SEEK_SET int = 0 // seek relative to the origin of the file
	SEEK_CUR int = 1 // seek relative to the current offset
	SEEK_END int = 2 // seek relative to the end
)

const (
	PathSeparator     = '/' // OS-specific path separator
	PathListSeparator = ':' // OS-specific path list separator
)

const DevNull = "/dev/null"

var (
	ErrInvalid    = errors.New("invalid argument")
	ErrPermission = errors.New("permission denied")
	ErrExist      = errors.New("file already exists")
	ErrNotExist   = errors.New("file does not exist")
)

type FileSystem interface {
	Chdir(dir string) error
	Chmod(name string, mode FileMode) error
	Chown(name string, uid, gid int) error
	Chtimes(name string, atime time.Time, mtime time.Time) error
	IsExist(err error) bool
	IsNotExist(err error) bool
	IsPathSeparator(c uint8) bool
	IsPermission(err error) bool
	Lchown(name string, uid, gid int) error
	Link(oldname, newname string) error
	Mkdir(name string, perm FileMode) error
	MkdirAll(path string, perm FileMode) error
	Readlink(name string) (string, error)
	Remove(name string) error
	RemoveAll(path string) error
	Rename(oldname, newname string) error
	SameFile(fi1, fi2 FileInfo) bool
	Setenv(key, value string) error
	Symlink(oldname, newname string) error
	TempDir() string
	Truncate(name string, size int64) error

	//
	Create(name string) (file *File, err error)
	NewFile(fd uintptr, name string) *File
	Open(name string) (file *File, err error)
	OpenFile(name string, flag int, perm FileMode) (file *File, err error)
	Pipe() (r *File, w *File, err error)

	//
	Lstat(name string) (fi FileInfo, err error)
	Stat(name string) (fi FileInfo, err error)
}

type File interface {
	Chdir() error
	Chmod(mode FileMode) error
	Chown(uid, gid int) error
	Close() error
	Fd() uintptr
	Name() string
	Read(b []byte) (n int, err error)
	ReadAt(b []byte, off int64) (n int, err error)
	Readdir(n int) (fi []FileInfo, err error)
	Readdirnames(n int) (names []string, err error)
	Seek(offset int64, whence int) (ret int64, err error)
	Stat() (fi FileInfo, err error)
	Sync() (err error)
	Truncate(size int64) error
	Write(b []byte) (n int, err error)
	WriteAt(b []byte, off int64) (n int, err error)
	WriteString(s string) (ret int, err error)
}

//    A FileInfo describes a file and is returned by Stat and Lstat.
type FileInfo interface {
	Name() string       // base name of the file
	Size() int64        // length in bytes for regular files; system-dependent for others
	Mode() FileMode     // file mode bits
	ModTime() time.Time // modification time
	IsDir() bool        // abbreviation for Mode().IsDir()
	Sys() interface{}   // underlying data source (can return nil)
}

type FileMode uint32

const (
	// The single letters are the abbreviations
	// used by the String method's formatting.
	ModeDir        FileMode = 1 << (32 - 1 - iota) // d: is a directory
	ModeAppend                                     // a: append-only
	ModeExclusive                                  // l: exclusive use
	ModeTemporary                                  // T: temporary file (not backed up)
	ModeSymlink                                    // L: symbolic link
	ModeDevice                                     // D: device file
	ModeNamedPipe                                  // p: named pipe (FIFO)
	ModeSocket                                     // S: Unix domain socket
	ModeSetuid                                     // u: setuid
	ModeSetgid                                     // g: setgid
	ModeCharDevice                                 // c: Unix character device, when ModeDevice is set
	ModeSticky                                     // t: sticky

	// Mask for the type bits. For regular files, none will be set.
	ModeType = ModeDir | ModeSymlink | ModeNamedPipe | ModeSocket | ModeDevice

	ModePerm FileMode = 0777 // permission bits
)

func (m FileMode) IsDir() bool {

}

func (m FileMode) IsRegular() bool {

}

func (m FileMode) Perm() FileMode {

}

func (m FileMode) String() string {

}

type LinkError struct {
	Op  string
	Old string
	New string
	Err error
}

func (e *LinkError) Error() string

type PathError struct {
	Op   string
	Path string
	Err  error
}

func (e *PathError) Error() string {

}
