package model

import (
	"bytes"
	utils "github.com/saichler/utils/golang"
	"strings"
)

type FileDescriptor struct {
	name         string
	targetName   string
	sourceParent *FileDescriptor
	targetParent *FileDescriptor
	size         int64
	hash         string
	parts        int
	files        map[string]*FileDescriptor
	sourcePath   string
	targetPath   string
}

func (fileDescriptor *FileDescriptor) Name() string {
	return fileDescriptor.name
}

func (fileDescriptor *FileDescriptor) Hash() string {
	return fileDescriptor.hash
}

func (fileDescriptor *FileDescriptor) SetHash(hash string) {
	fileDescriptor.hash = hash
}

func (fileDescriptor *FileDescriptor) SourceParent() *FileDescriptor {
	return fileDescriptor.sourceParent
}

func (fileDescriptor *FileDescriptor) Size() int64 {
	return fileDescriptor.size
}

func (fileDescriptor *FileDescriptor) Parts() int {
	return fileDescriptor.parts
}

func (fileDescriptor *FileDescriptor) SetPart(p int) {
	fileDescriptor.parts = p
}

func (fileDescriptor *FileDescriptor) IsDir() bool {
	return !(fileDescriptor.files == nil || len(fileDescriptor.files) == 0)
}

func (fileDescriptor *FileDescriptor) Files() map[string]*FileDescriptor {
	return fileDescriptor.files
}

func (fileDescriptor *FileDescriptor) SourceRoot() *FileDescriptor {
	if fileDescriptor.sourceParent == nil {
		return fileDescriptor
	}
	return fileDescriptor.sourceParent.SourceRoot()
}

func (fileDescriptor *FileDescriptor) Get(path string) *FileDescriptor {
	index := strings.Index(path, "/")
	if index == 0 {
		path = path[1:]
		index = strings.Index(path, "/")
	}
	if index == -1 {
		index = len(path)
	}
	name := path[0:index]
	child := fileDescriptor.files[name]
	if child != nil && index != len(path) {
		return child.Get(path[index+1:])
	}
	return child
}

func (fileDescriptor *FileDescriptor) SourcePath() string {
	if fileDescriptor.sourcePath != "" {
		return fileDescriptor.sourcePath
	}
	if fileDescriptor.sourceParent == nil {
		return "/"
	}
	buff := &bytes.Buffer{}
	fileDescriptor._sourcePath(buff)
	fileDescriptor.sourcePath = buff.String()
	return fileDescriptor.sourcePath
}

func (fileDescriptor *FileDescriptor) _sourcePath(buff *bytes.Buffer) {
	if fileDescriptor.sourceParent == nil {
		return
	}
	fileDescriptor.sourceParent._sourcePath(buff)
	buff.WriteString("/")
	buff.WriteString(fileDescriptor.name)
}

func (fileDescriptor *FileDescriptor) Bytes() []byte {
	bs := utils.NewByteSlice()
	fileDescriptor.AddBytes(bs)
	return bs.Data()
}

func (fileDescriptor *FileDescriptor) AddBytes(bs *utils.ByteSlice) {
	bs.AddString(fileDescriptor.SourcePath())
	fileDescriptor.SourceRoot().bytes(bs)
}

func (fileDescriptor *FileDescriptor) bytes(bs *utils.ByteSlice) {
	bs.AddString(fileDescriptor.name)
	bs.AddInt64(fileDescriptor.size)
	bs.AddString(fileDescriptor.hash)
	bs.AddInt(fileDescriptor.parts)
	if fileDescriptor.files == nil {
		bs.AddInt(0)
	} else {
		bs.AddInt(len(fileDescriptor.files))
		for name, fdChild := range fileDescriptor.files {
			bs.AddString(name)
			fdChild.bytes(bs)
		}
	}
}

func Object(bs *utils.ByteSlice) *FileDescriptor {
	path := bs.GetString()
	root := &FileDescriptor{}
	root.object(bs)
	child := root.Get(path)
	return child
}

func (fileDescriptor *FileDescriptor) object(bs *utils.ByteSlice) {
	fileDescriptor.name = bs.GetString()
	fileDescriptor.size = bs.GetInt64()
	fileDescriptor.hash = bs.GetString()
	fileDescriptor.parts = bs.GetInt()
	fileDescriptor.files = make(map[string]*FileDescriptor)
	size := bs.GetInt()
	for i := 0; i < size; i++ {
		key := bs.GetString()
		child := &FileDescriptor{}
		child.object(bs)
		fileDescriptor.files[key] = child
		child.sourceParent = fileDescriptor
	}
}
