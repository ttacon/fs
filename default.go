package fs

import (
	"os"
	"time"
)

type defaultOS struct {
	Stdin, Stdout, Stderr File
}

func DefaultOS() OperatingSystem {
	d := &defaultOS{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	return d
}

func (d *defaultOS) Chdir(dir string) error {
	return os.Chdir(dir)
}

func (d *defaultOS) Chmod(name string, mode os.FileMode) error {
	return os.Chmod(name, mode)
}

func (d *defaultOS) Chown(name string, uid, gid int) error {
	return os.Chown(name, uid, gid)
}

func (d *defaultOS) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return os.Chtimes(name, atime, mtime)
}

func (d *defaultOS) Clearenv() {
	os.Clearenv()
}

func (d *defaultOS) Environ() []string {
	return os.Environ()
}

func (d *defaultOS) Exit(code int) {
	os.Exit(code)
}

func (d *defaultOS) Expand(s string, mapping func(string) string) string {
	return os.Expand(s, mapping)
}

func (d *defaultOS) ExpandEnv(s string) string {
	return os.ExpandEnv(s)
}

func (d *defaultOS) Getegid() int {
	return os.Getegid()
}

func (d *defaultOS) Getenv(key string) string {
	return os.Getenv(key)
}

func (d *defaultOS) Geteuid() int {
	return os.Geteuid()
}

func (d *defaultOS) Getgid() int {
	return os.Getgid()
}

func (d *defaultOS) Getgroups() ([]int, error) {
	return os.Getgroups()
}

func (d *defaultOS) Getpagesize() int {
	return os.Getpagesize()
}

func (d *defaultOS) Getpid() int {
	return os.Getpid()
}

func (d *defaultOS) Getppid() int {
	return os.Getppid()
}

func (d *defaultOS) Getuid() int {
	return os.Getuid()
}

func (d *defaultOS) Getwd() (dir string, err error) {
	return os.Getwd()
}

func (d *defaultOS) Hostname() (name string, err error) {
	return os.Hostname()
}

func (d *defaultOS) IsExist(err error) bool {
	return os.IsExist(err)
}

func (d *defaultOS) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

func (d *defaultOS) IsPathSeparator(c uint8) bool {
	return os.IsPathSeparator(c)
}

func (d *defaultOS) IsPermission(err error) bool {
	return os.IsPermission(err)
}

func (d *defaultOS) Lchown(name string, uid, gid int) error {
	return os.Lchown(name, uid, gid)
}

func (d *defaultOS) Link(oldname, newname string) error {
	return os.Link(oldname, newname)
}

func (d *defaultOS) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

func (d *defaultOS) MkdirAll(path string, perm os.FileMode) error {
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

func (d *defaultOS) SameFile(fi1, fi2 os.FileInfo) bool {
	return os.SameFile(fi1, fi2)
}

func (d *defaultOS) Setenv(key, value string) error {
	return os.Setenv(key, value)
}

func (d *defaultOS) Symlink(oldname, newname string) error {
	return os.Symlink(oldname, newname)
}

func (d *defaultOS) TempDir() string {
	return os.TempDir()
}

func (d *defaultOS) Truncate(name string, size int64) error {
	return os.Truncate(name, size)
}

//
func (d *defaultOS) Create(name string) (file File, err error) {
	return os.Create(name)
}

func (d *defaultOS) NewFile(fd uintptr, name string) File {
	return os.NewFile(fd, name)
}

func (d *defaultOS) Open(name string) (file File, err error) {
	return os.Open(name)
}

func (d *defaultOS) OpenFile(name string, flag int, perm os.FileMode) (file File, err error) {
	return os.OpenFile(name, flag, perm)
}

func (d *defaultOS) Pipe() (r File, w File, err error) {
	return os.Pipe()
}

//
func (d *defaultOS) Lstat(name string) (fi os.FileInfo, err error) {
	return os.Lstat(name)
}

func (d *defaultOS) Stat(name string) (fi os.FileInfo, err error) {
	return os.Stat(name)
}
