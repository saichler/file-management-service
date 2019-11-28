package model

import utils "github.com/saichler/utils/golang"

type FileRequest struct {
	path     string
	calcHash bool
	dept     int
}

func NewFileRequest(path string, dept int, calcHash bool) *FileRequest {
	fr := &FileRequest{}
	fr.path = path
	fr.dept = dept
	fr.calcHash = calcHash
	return fr
}

func (fr *FileRequest) Path() string {
	return fr.path
}

func (fr *FileRequest) Dept() int {
	return fr.dept
}

func (fr *FileRequest) CalcHash() bool {
	return fr.calcHash
}

func (fr *FileRequest) Bytes(bs *utils.ByteSlice) []byte {
	bs.AddString(fr.path)
	bs.AddInt(fr.dept)
	bs.AddBool(fr.calcHash)
	return bs.Data()
}

func (fr *FileRequest) Object(bs *utils.ByteSlice) {
	fr.path = bs.GetString()
	fr.dept = bs.GetInt()
	fr.calcHash = bs.GetBool()
}
