package fs

import (
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var (
	host = "fs-" + strconv.FormatInt(rand.Int63(), 10)
)

type fakeOS struct {
	lock                  *sync.Mutex
	envLock               *sync.RWMutex
	Stdin, Stdout, Stderr File
	envVars               map[string]string
	files                 map[string]*fakeFile
	cwd                   string

	// ??? where should this go?
	tmpDir string

	// current user info
	uid, gid int

	// other info
	pagesize  int
	pid, ppid int
}

func FakeOS() OperatingSystem {
	tmpDir := string(filepath.Separator) + "tmp"
	d := &fakeOS{
		lock:    new(sync.Mutex),
		envLock: new(sync.RWMutex),
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		envVars: map[string]string{},
		files: map[string]*fakeFile{
			// TODO(ttacon): add cwd file?
			tmpDir: &fakeFile{ /*TODO(ttacon): add info*/ },
		},
		tmpDir: tmpDir,

		// TODO(ttacon): better values for these?
		uid: 501, // no idea what a good value for this is
		// maybe grab the real one?
		gid: 20, //  this is staff on macs? better value?

		// interestingly, for any go program pid = ppid +3
		pid:  18012,
		ppid: 18009,
	}
	return d
}

func (d *fakeOS) Chdir(dir string) error {
	d.lock.Lock()
	if _, ok := d.files[dir]; !ok {
		d.lock.Unlock()
		return &os.PathError{
			Op:   "chdir",
			Path: dir,
			Err:  syscall.Errno(syscall.ENOENT),
		}
	}

	d.cwd = dir
	d.lock.Unlock()
	return nil
}

func (d *fakeOS) Chmod(name string, mode os.FileMode) error {
	d.lock.Lock()
	// TODO(ttacon): need to be able to simulate if the current
	// user has permission to do this or not

	f, ok := d.files[name]
	if !ok {
		d.lock.Unlock()
		return &os.PathError{
			Op:   "chmod",
			Path: name,
			Err:  syscall.Errno(syscall.ENOENT),
		}
	}
	f.mode = mode
	d.lock.Unlock()
	return nil
}

func (d *fakeOS) Chown(name string, uid, gid int) error {
	d.lock.Lock()
	// TODO(ttacon): need to be able to simulate if the current
	// user has permission to do this or not

	f, ok := d.files[name]
	if !ok {
		d.lock.Unlock()
		return &os.PathError{
			Op:   "chmod",
			Path: name,
			Err:  syscall.Errno(syscall.ENOENT),
		}
	}

	f.uid, f.gid = uid, gid
	d.lock.Unlock()
	return nil
}

func (d *fakeOS) Chtimes(name string, atime time.Time, mtime time.Time) error {
	d.lock.Lock()
	// TODO(ttacon): need to be able to simulate if the current
	// user has permission to do this or not

	f, ok := d.files[name]
	if !ok {
		d.lock.Unlock()
		return &os.PathError{
			Op:   "chmod",
			Path: name,
			Err:  syscall.Errno(syscall.ENOENT),
		}
	}

	f.access, f.modify = atime, mtime
	d.lock.Unlock()
	return nil
}

func (d *fakeOS) Clearenv() {
	d.envLock.Lock()
	d.envVars = map[string]string{}
	d.envLock.Unlock()
}

func (d *fakeOS) Environ() []string {
	d.envLock.RLock()
	var (
		i     = 0
		toRet = make([]string, len(d.envVars))
	)

	for k, v := range d.envVars {
		toRet[i] = k + "=" + v
		i++
	}
	d.envLock.RUnlock()
	return toRet
}

func (d *fakeOS) Exit(code int) {
	// actually exit so just use os version
	os.Exit(code)
}

func (d *fakeOS) Expand(s string, mapping func(string) string) string {
	// this doesn't actually touch OS so we can use the os one
	return os.Expand(s, mapping)
}

func (d *fakeOS) ExpandEnv(s string) string {
	return d.Expand(s, d.Getenv)
}

func (d *fakeOS) Getegid() int {
	// TODO(ttacon): simulate egid vs gid?
	return d.gid
}

func (d *fakeOS) Getenv(key string) string {
	d.envLock.RLock()
	v := d.envVars[key]
	d.envLock.RUnlock()
	return v
}

func (d *fakeOS) Geteuid() int {
	// TODO(ttacon): same question as Getegid
	return d.uid
}

func (d *fakeOS) Getgid() int {
	return d.gid
}

func (d *fakeOS) Getgroups() ([]int, error) {
	panic("UNIMPLEMENTED")
}

func (d *fakeOS) Getpagesize() int {
	return d.pagesize
}

func (d *fakeOS) Getpid() int {
	return d.pid
}

func (d *fakeOS) Getppid() int {
	return d.ppid
}

func (d *fakeOS) Getuid() int {
	return d.uid
}

func (d *fakeOS) Getwd() (dir string, err error) {
	// is this err non-nil on permission switching induced issues?
	return d.cwd, nil
}

func (d *fakeOS) Hostname() (name string, err error) {
	return host, nil
}

func (d *fakeOS) IsExist(err error) bool {
	return os.IsExist(err)
}

func (d *fakeOS) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

func (d *fakeOS) IsPathSeparator(c uint8) bool {
	return os.IsPathSeparator(c)
}

func (d *fakeOS) IsPermission(err error) bool {
	return os.IsPermission(err)
}

func (d *fakeOS) Lchown(name string, uid, gid int) error {
	d.lock.Lock()
	// TODO(ttacon): need to be able to simulate if the current
	// user has permission to do this or not

	f, ok := d.files[name]
	if !ok {
		d.lock.Unlock()
		return &os.PathError{
			Op:   "chmod",
			Path: name,
			Err:  syscall.Errno(syscall.ENOENT),
		}
	}

	f.uid, f.gid = uid, gid
	d.lock.Unlock()
	return nil
}

func (d *fakeOS) Link(oldname, newname string) error {
	d.lock.Lock()
	// TODO(ttacon): need to be able to simulate if the current
	// user has permission to do this or not

	f, ok := d.files[oldname]
	if !ok {
		d.lock.Unlock()
		return &os.PathError{
			Op:   "link",
			Path: oldname,
			Err:  syscall.Errno(syscall.ENOENT),
		}
	}

	if _, ok := d.files[newname]; !ok {
		d.lock.Unlock()
		return &os.PathError{
			Op:   "link",
			Path: newname,
			Err:  syscall.Errno(syscall.EEXIST),
		}
	}

	d.files[newname] = f
	d.lock.Unlock()
	return nil
}

func (d *fakeOS) Mkdir(name string, perm os.FileMode) error {
	base := filepath.Dir(name)
	d.lock.Lock()
	// TODO(ttacon): need to be able to simulate if the current
	// user has permission to do this or not

	f, ok := d.files[base]
	if !ok {
		d.lock.Unlock()
		return &os.PathError{
			Op:   "mkdir",
			Path: base,
			Err:  syscall.Errno(syscall.ENOENT),
		}
	}

	now := time.Now()
	d.files[name] = &fakeFile{
		access: now,
		modify: now,
		change: now,
		isDir:  true,
		mode:   os.ModeDir,
		info:   nil, // TODO(ttacon): do it
		uid:    f.uid,
		gid:    f.gid,
	}
	d.lock.Unlock()
	return nil
}

func (d *fakeOS) MkdirAll(path string, perm os.FileMode) error {
	pieces := filepath.SplitList(path)
	d.lock.Lock()
	curr := string(filepath.Separator)
	now := time.Now()
	for _, piece := range pieces {
		curr = filepath.Join(curr, piece)
		if _, ok := d.files[curr]; ok {
			continue
		}
		d.files[curr] = &fakeFile{
			access: now,
			modify: now,
			change: now,
			isDir:  true,
			mode:   perm,
			info:   nil, // TODO(ttaco): do it
			uid:    d.uid,
			gid:    d.gid,
		}
	}
	d.lock.Unlock()
	return nil
}

func (d *fakeOS) Readlink(name string) (string, error) {
	d.lock.Lock()
	// TODO(ttacon): need to be able to simulate if the current
	// user has permission to do this or not

	f, ok := d.files[name]
	if !ok {
		d.lock.Unlock()
		return "", &os.PathError{
			Op:   "readlink",
			Path: name,
			Err:  syscall.Errno(syscall.ENOENT),
		}
	}

	if f.pointsTo == "" {
		d.lock.Unlock()
		return "", &os.PathError{
			Op:   "readlink",
			Path: name,
			Err:  syscall.Errno(syscall.EINVAL),
		}
	}

	toReturn := f.pointsTo
	d.lock.Unlock()
	return toReturn, nil
}

func (d *fakeOS) Remove(name string) error {
	d.lock.Lock()
	// TODO(ttacon): need to be able to simulate if the current
	// user has permission to do this or not

	_, ok := d.files[name]
	if !ok {
		d.lock.Unlock()
		return &os.PathError{
			Op:   "remove",
			Path: name,
			Err:  syscall.Errno(syscall.ENOENT),
		}
	}

	delete(d.files, name)
	d.lock.Unlock()
	return nil
}

func (d *fakeOS) RemoveAll(path string) error {
	panic("UNIMPLEMENTED")
}

func (d *fakeOS) Rename(oldname, newname string) error {
	d.lock.Lock()
	// TODO(ttacon): need to be able to simulate if the current
	// user has permission to do this or not

	f, ok := d.files[oldname]
	if !ok {
		d.lock.Unlock()
		return &os.PathError{
			Op:   "rename",
			Path: oldname,
			Err:  syscall.Errno(syscall.ENOENT),
		}
	}

	delete(d.files, oldname)
	d.files[newname] = f

	d.lock.Unlock()
	return nil
}

func (d *fakeOS) SameFile(fi1, fi2 os.FileInfo) bool {
	panic("UNIMPLEMENTED")
}

func (d *fakeOS) Setenv(key, value string) error {
	d.envLock.Lock()
	d.envVars[key] = value
	d.envLock.Unlock()
	return nil
}

func (d *fakeOS) Symlink(oldname, newname string) error {
	d.lock.Lock()
	// TODO(ttacon): need to be able to simulate if the current
	// user has permission to do this or not

	f, ok := d.files[oldname]
	if !ok {
		d.lock.Unlock()
		return &os.PathError{
			Op:   "symlink",
			Path: oldname,
			Err:  syscall.Errno(syscall.ENOENT),
		}
	}

	// TODO(ttacon): this needs to be able to differentiate between
	// hard and soft links (Link() vs Symlink())
	if _, ok := d.files[newname]; !ok {
		d.lock.Unlock()
		return &os.PathError{
			Op:   "symlink",
			Path: newname,
			Err:  syscall.Errno(syscall.EEXIST),
		}
	}

	now := time.Now()
	d.files[newname] = &fakeFile{
		access:   now,
		modify:   now,
		change:   now,
		isDir:    false,
		mode:     f.mode | os.ModeSymlink,
		info:     nil, // TODO(ttacon): do it
		uid:      f.uid,
		gid:      f.gid,
		pointsTo: oldname,
	}

	d.lock.Unlock()
	return nil
}

func (d *fakeOS) TempDir() string {
	return d.tmpDir
}

func (d *fakeOS) Truncate(name string, size int64) error {
	d.lock.Lock()
	// TODO(ttacon): need to be able to simulate if the current
	// user has permission to do this or not

	f, ok := d.files[name]
	if !ok {
		d.lock.Unlock()
		return &os.PathError{
			Op:   "symlink",
			Path: name,
			Err:  syscall.Errno(syscall.ENOENT),
		}
	}

	fSize := int64(len(f.content))

	if size > fSize {
		d.lock.Unlock()
		return nil
	}

	newContent := make([]byte, len(f.content), cap(f.content))
	copy(newContent, f.content[:size])
	f.content = newContent

	d.lock.Unlock()
	return nil
}

func (d *fakeOS) Create(name string) (file File, err error) {
	d.lock.Lock()
	// TODO(ttacon): need to be able to simulate if the current
	// user has permission to do this or not

	if _, ok := d.files[name]; ok {
		d.lock.Unlock()
		return nil, &os.PathError{
			Op:   "create",
			Path: name,
			Err:  syscall.Errno(syscall.EEXIST),
		}
	}

	now := time.Now()
	f := &fakeFile{
		access: now,
		modify: now,
		change: now,
		isDir:  true,
		mode:   os.ModePerm,
		info:   nil, // TODO(ttacon): do it
		uid:    d.uid,
		gid:    d.gid,
	}
	d.files[name] = f

	d.lock.Unlock()
	return f, nil
}

func (d *fakeOS) NewFile(fd uintptr, name string) File {
	// TODO(ttacon): swalllow fd?
	d.lock.Lock()
	// TODO(ttacon): need to be able to simulate if the current
	// user has permission to do this or not

	now := time.Now()
	f := &fakeFile{
		access: now,
		modify: now,
		change: now,
		isDir:  true,
		mode:   os.ModePerm,
		info:   nil, // TODO(ttacon): do it
		uid:    d.uid,
		gid:    d.gid,
	}
	d.files[name] = f

	d.lock.Unlock()
	return f
}

func (d *fakeOS) Open(name string) (file File, err error) {
	d.lock.Lock()
	// TODO(ttacon): need to be able to simulate if the current
	// user has permission to do this or not

	f, ok := d.files[name]
	if !ok {
		d.lock.Unlock()
		return nil, &os.PathError{
			Op:   "open",
			Path: name,
			Err:  syscall.Errno(syscall.ENOENT),
		}
	}

	d.lock.Unlock()
	return f, nil
}

func (d *fakeOS) OpenFile(name string, flag int, perm os.FileMode) (file File, err error) {
	d.lock.Lock()
	// TODO(ttacon): need to be able to simulate if the current
	// user has permission to do this or not

	if _, ok := d.files[name]; ok {
		d.lock.Unlock()
		return nil, &os.PathError{
			Op:   "newfile",
			Path: name,
			Err:  syscall.Errno(syscall.EEXIST),
		}
	}

	// TODO(ttacon): how is this different from Open()?
	now := time.Now()
	f := &fakeFile{
		access: now,
		modify: now,
		change: now,
		isDir:  true,
		mode:   os.ModePerm,
		info:   nil, // TODO(ttacon): do it
		uid:    d.uid,
		gid:    d.gid,
	}
	d.files[name] = f

	d.lock.Unlock()
	return f, nil
}

func (d *fakeOS) Pipe() (r File, w File, err error) {
	panic("UNIMPLEMENTED")
}

func (d *fakeOS) Lstat(name string) (fi os.FileInfo, err error) {
	panic("UNIMPLEMENTED")
}

func (d *fakeOS) Stat(name string) (fi os.FileInfo, err error) {
	d.lock.Lock()
	// TODO(ttacon): need to be able to simulate if the current
	// user has permission to do this or not

	f, ok := d.files[name]
	if !ok {
		d.lock.Unlock()
		return nil, &os.PathError{
			Op:   "chmod",
			Path: name,
			Err:  syscall.Errno(syscall.ENOENT),
		}
	}

	toReturn := f.info
	d.lock.Unlock()
	return toReturn, nil
}
