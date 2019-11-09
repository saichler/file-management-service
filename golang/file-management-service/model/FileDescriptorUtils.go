package model

import (
	"bytes"
	"github.com/saichler/security"
	utils "github.com/saichler/utils/golang"
	"io/ioutil"
	"os"
	"strings"
)

func createEmpty(path string) *FileDescriptor {
	index := strings.Index(path, "/")
	descriptor := &FileDescriptor{}
	if index == 0 {
		descriptor.name = "/"
	} else if index == -1 {
		descriptor.name = path
		return descriptor
	} else {
		descriptor.name = path[0:index]
	}
	descriptor.files = make(map[string]*FileDescriptor)
	child := createEmpty(path[index+1:])
	child.sourceParent = descriptor
	descriptor.files[child.name] = child
	return descriptor

}

func NewFileDescriptor(path string, dept int, calcHash bool) *FileDescriptor {
	root := createEmpty(path)
	descriptor := root.Get(path)
	if calcHash {
		pool := utils.NewWorkerPool(50)
		pool.Start()
		fill(descriptor, dept, 0, pool)
		pool.WaitForEmptyQueue()
		pool.Stop()
		return descriptor
	}
	fill(descriptor, dept, 0, nil)
	return descriptor
}

func fill(descriptor *FileDescriptor, dept, current int, pool *utils.WorkerPool) {
	file, e := os.Stat(descriptor.SourcePath())
	if e != nil {
		return
	}

	descriptor.size = file.Size()

	if file.IsDir() && current < dept {
		descriptor.files = make(map[string]*FileDescriptor)
		files, e := ioutil.ReadDir(descriptor.sourcePath)
		if e == nil {
			for _, file := range files {
				child := &FileDescriptor{}
				child.sourceParent = descriptor
				child.name = file.Name()
				fill(child, dept, current+1, pool)
				descriptor.files[child.name] = child
			}
		}
	} else if !file.IsDir() {
		descriptor.parts = int(descriptor.size/MAX_PART_SIZE) + 1
		if pool != nil {
			hashJob := &HashJob{}
			hashJob.descriptor = descriptor
			pool.AddTask(hashJob)
		}
	}
}

func (fileDescriptor *FileDescriptor) SetTargetName(n string) {
	fileDescriptor.targetName = n
}

func (fileDescriptor *FileDescriptor) SetTargetParent(descriptor *FileDescriptor) {
	fileDescriptor.targetParent = descriptor
}

func (fileDescriptor *FileDescriptor) TargetParent() *FileDescriptor {
	return fileDescriptor.targetParent
}

func (fileDescriptor *FileDescriptor) TargetRoot() *FileDescriptor {
	if fileDescriptor.sourceParent == nil {
		return fileDescriptor
	}
	if fileDescriptor.targetParent != nil {
		return fileDescriptor.targetParent.TargetRoot()
	}
	return fileDescriptor.sourceParent.TargetRoot()
}

func (fileDescriptor *FileDescriptor) TargetPath() string {
	if fileDescriptor.targetPath != "" {
		return fileDescriptor.targetPath
	}
	if fileDescriptor.sourceParent == nil {
		return "/"
	}
	buff := &bytes.Buffer{}
	fileDescriptor._targetPath(buff)
	fileDescriptor.targetPath = buff.String()
	return fileDescriptor.targetPath
}

func (fileDescriptor *FileDescriptor) _targetPath(buff *bytes.Buffer) {
	if fileDescriptor.sourceParent == nil {
		return
	}
	if fileDescriptor.targetParent != nil {
		fileDescriptor.targetParent._targetPath(buff)
	} else {
		fileDescriptor.sourceParent._targetPath(buff)
	}
	buff.WriteString("/")
	if fileDescriptor.targetName != "" {
		buff.WriteString(fileDescriptor.targetName)
	} else {
		buff.WriteString(fileDescriptor.name)
	}
}

type HashJob struct {
	descriptor *FileDescriptor
}

func (hashJob *HashJob) Run() {
	hashJob.descriptor.hash, _ = security.FileHash(hashJob.descriptor.SourcePath())
}
