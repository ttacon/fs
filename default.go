package fs

import (
	"os"
	"time"
)

type defaultOS struct {
	Stdin, Stdout, Stderr *File
}

func DefaultOS() FileSystem {
	d := &defaultOS{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	return d
}

func (d *defaultOS) Chdir(dir string) error {

}

func (d *defaultOS) Chmod(name string, mode FileMode) error {

}

func (d *defaultOS) Chown(name string, uid, gid int) error {

}

func (d *defaultOS) Chtimes(name string, atime time.Time, mtime time.Time) error {

}

func (d *defaultOS) IsExist(err error) bool {

}

func (d *defaultOS) IsNotExist(err error) bool {

}

func (d *defaultOS) IsPathSeparator(c uint8) bool {

}

func (d *defaultOS) IsPermission(err error) bool {

}

func (d *defaultOS) Lchown(name string, uid, gid int) error {

}

func (d *defaultOS) Link(oldname, newname string) error {
	return os.Link(oldname, newname)
}

func (d *defaultOS) Mkdir(name string, perm FileMode) error {
	return os.Mkdir(name, perm)
}

func (d *defaultOS) MkdirAll(path string, perm FileMode) error {
	return os.MkdirAll(path, perm)
}

func (d *defaultOS) Readlink(name string) (string, error) {
	return os.Readlink(name)
}

func (d *defaultOS) Remove(name string) error {
	return os.Remove(name)
}

func (d *defaultOS) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (d *defaultOS) Rename(oldname, newname string) error {
	return os.Rename(oldname, newname)
}

func (d *defaultOS) SameFile(fi1, fi2 FileInfo) bool {
	return os.SameFile(fi1, fi2)
}

func (d *defaultOS) Setenv(key, value string) error {
	return os.Setenv(key, value)
}

func (d *defaultOS) Symlink(oldname, newname string) error {
	return os.Symlink()
}

func (d *defaultOS) TempDir() string {
	return os.TempDir()
}

func (d *defaultOS) Truncate(name string, size int64) error {
	return os.Truncate(name, size)
}

//
func (d *defaultOS) Create(name string) (file *File, err error) {
	return os.Create(name)
}

func (d *defaultOS) NewFile(fd uintptr, name string) *File {
	return os.NewFile(fd, name)
}

func (d *defaultOS) Open(name string) (file *File, err error) {
	return os.Open(name)
}

func (d *defaultOS) OpenFile(name string, flag int, perm FileMode) (file *File, err error) {
	return os.OpenFile(name, flag, perm)
}

func (d *defaultOS) Pipe() (r *File, w *File, err error) {
	return os.Pipe()
}

//
func (d *defaultOS) Lstat(name string) (fi FileInfo, err error) {
	return os.Lstat(name)
}

func (d *defaultOS) Stat(name string) (fi FileInfo, err error) {
	return os.Stat(name)
}
