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

func (fr *FileRequest) Marshal() []byte {
	bs := utils.NewByteSlice()
	bs.AddString(fr.path)
	bs.AddInt(fr.dept)
	bs.AddBool(fr.calcHash)
	return bs.Data()
}

func (fr *FileRequest) UnMarshal(data []byte) {
	bs := utils.NewByteSliceWithData(data, 0)
	fr.path = bs.GetString()
	fr.dept = bs.GetInt()
	fr.calcHash = bs.GetBool()
}
