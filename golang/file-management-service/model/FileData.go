package model

import (
	utils "github.com/saichler/utils/golang"
	"os"
)

const (
	MAX_PART_SIZE = 1024 * 1024 * 8
)

type FileData struct {
	path string
	part int
	size int64
	data []byte
}

func NewFileData(path string, part int, size int64) *FileData {
	fileData := &FileData{}
	fileData.path = path
	fileData.part = part
	fileData.size = size
	if fileData.data == nil {
		fileData.data = make([]byte, 0)
	}
	return fileData
}

func (fileData *FileData) LoadData() {
	loc := MAX_PART_SIZE * int64(fileData.part)
	dataSize := MAX_PART_SIZE
	if fileData.size-loc < MAX_PART_SIZE {
		dataSize = int(fileData.size - loc)
	}

	data := make([]byte, dataSize)
	file, _ := os.Open(fileData.path)
	_, e := file.Seek(loc, 0)
	if e != nil {
		panic(e)
	}

	_, e = file.Read(data)
	if e != nil {
		panic(e)
	}
	file.Close()
	fileData.data = data
}

func (fileData *FileData) Marshal() []byte {
	bs := utils.NewByteSlice()
	bs.AddInt(fileData.part)
	bs.AddInt64(fileData.size)
	bs.AddString(fileData.path)
	bs.AddByteSlice(fileData.data)
	return bs.Data()
}

func (fileData *FileData) Unmarshal(data []byte) {
	bs := utils.NewByteSliceWithData(data, 0)
	fileData.part = bs.GetInt()
	fileData.size = bs.GetInt64()
	fileData.path = bs.GetString()
	fileData.data = bs.GetByteSlice()
}

func (fileData *FileData) Part() int {
	return fileData.part
}

func (fileData *FileData) Path() string {
	return fileData.path
}

func (fileData *FileData) Data() []byte {
	return fileData.data
}
