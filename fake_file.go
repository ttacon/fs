package fs

import (
	"os"
	"syscall"
	"time"
)

type fakeFile struct {
	fd                     int
	name                   string
	access, modify, change time.Time
	isDir                  bool
	rdwrFlag               int
	mode                   os.FileMode
	info                   os.FileInfo
	uid, gid               int
	pointsTo               string // for links
	currPos                int64
	content                []byte
}

func newFakeFile(
	name string,
	mode os.FileMode,
	info os.FileInfo,
	rdwrFlag int) *fakeFile {
	now := time.Now()
	return &fakeFile{
		name:     name,
		access:   now,
		modify:   now,
		change:   now,
		rdwrFlag: rdwrFlag,
		info:     info,
		mode:     mode,
		uid:      0, // TODO(ttacon): query os singleton?
		gid:      0, // TODO(ttacon): query os singleton?
	}
}

func newDir(name string, info os.FileInfo) *fakeFile {
	now := time.Now()
	return &fakeFile{
		name:   name,
		access: now,
		modify: now,
		change: now,
		isDir:  true,
		// rdwrFlag?
		info: info,
		mode: os.ModeDir,
		uid:  0, // TODO(ttacon): query os singleton?
		gid:  0, // TODO(ttacon): query os singleton?
	}
}

func (f *fakeFile) Chdir() error {
	return currOs.Chdir(f.name)
}

func (f *fakeFile) Chmod(mode os.FileMode) error {
	return currOs.Chmod(f.name, mode)
}

func (f *fakeFile) Chown(uid, gid int) error {
	return currOs.Chown(f.name, uid, gid)
}

const O_CLOSED = -1

func (f *fakeFile) Close() error {
	// TODO(ttacon): Can this return an error?
	f.rdwrFlag = O_CLOSED
	return nil
}

func (f *fakeFile) Fd() uintptr {
	return uintptr(f.fd)
}

func (f *fakeFile) Name() string {
	return f.name
}

func (f *fakeFile) Read(b []byte) (n int, err error) {
	desired := int64(len(b))
	toRead := desired

	if desired > int64(len(f.content))-f.currPos {
		toRead = int64(len(f.content)) - f.currPos
	}

	copy(b, f.content[f.currPos:f.currPos+toRead])
	f.currPos += toRead
	return int(toRead), nil
}

func (f *fakeFile) ReadAt(b []byte, off int64) (n int, err error) {
	desired := int64(len(b))
	toRead := desired

	if desired > int64(len(f.content))-off {
		toRead = int64(len(f.content)) - off
	}

	copy(b, f.content[off:off+toRead])
	f.currPos = off + toRead
	return int(toRead), nil
}

func (f *fakeFile) Readdir(n int) (fi []os.FileInfo, err error) {
	// TODO(ttacon): do it
	panic("UNIMPLEMENTED")
}

func (f *fakeFile) Readdirnames(n int) (names []string, err error) {
	// TODO(ttacon): do it
	panic("UNIMPLEMENTED")
}

func (f *fakeFile) Seek(offset int64, whence int) (ret int64, err error) {
	var newOffset int64
	if whence == 0 {
		newOffset = offset
	} else if whence == 1 {
		newOffset = f.currPos + offset
	} else if whence == 2 {
		newOffset = int64(len(f.content)) + offset
	}

	if newOffset < 0 {
		// aparently seeking past the end is an error?
		return newOffset, &os.PathError{
			Op:   "seek",
			Path: f.name,
			Err:  syscall.EINVAL,
		}
	}
	f.currPos = newOffset
	return f.currPos, nil
}

func (f *fakeFile) Stat() (fi os.FileInfo, err error) {
	return f, nil
}

func (f *fakeFile) Sync() (err error) {
	// nothing to do
	return nil
}

func (f *fakeFile) Truncate(size int64) error {
	// TODO(ttacon): do some testing with this function to see
	// how the real one behaves

	var newContents = make([]byte, size)
	if size > int64(len(f.content)) {
		size = int64(len(f.content))
	}
	copy(newContents, f.content[:size])
	f.content = newContents
	return nil
}

func (f *fakeFile) Write(b []byte) (n int, err error) {
	space := int64(len(f.content)) - f.currPos
	if space > 0 {
		copy(f.content[f.currPos:], b[:space])
		f.content = append(f.content, b[:space]...)
	} else {
		f.content = append(f.content, b...)
	}
	f.currPos = int64(len(f.content))
	return len(b), nil
}

func (f *fakeFile) WriteAt(b []byte, off int64) (n int, err error) {
	// TODO(ttacon): do it
	panic("UNIMPLEMENTED")
}

func (f *fakeFile) WriteString(s string) (ret int, err error) {
	return f.Write([]byte(s))
}

////////// to make life simpler with os.FileInfo
func (f *fakeFile) Size() int64 {
	return int64(len(f.content))
}

func (f *fakeFile) Mode() os.FileMode {
	return f.mode
}

func (f *fakeFile) ModTime() time.Time {
	return f.modify
}

func (f *fakeFile) IsDir() bool {
	return f.isDir
}

func (f *fakeFile) Sys() interface{} {
	// TODO(ttacon): don't feel like doing this one right now
	return nil
}
