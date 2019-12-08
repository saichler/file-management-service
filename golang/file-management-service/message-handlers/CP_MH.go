package message_handlers

import (
	"bytes"
	"errors"
	model2 "github.com/saichler/file-management-service/golang/file-management-service/model"
	fms "github.com/saichler/file-management-service/golang/file-management-service/service"
	. "github.com/saichler/messaging/golang/net/protocol"
	. "github.com/saichler/service-manager/golang/service-manager"
	. "github.com/saichler/utils/golang"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"time"
)

type CP_MH struct {
	fs *fms.FileManagerService
}

type FetchPartTask struct {
	cp         *CP_MH
	descriptor *model2.FileDescriptor
	part       int
}

type CopyFileJobListener struct {
	filename    string
	finishCount int
	parts       int
	lastReport  int64
}

func NewCpMH(service IService) *CP_MH {
	lf := &CP_MH{}
	lf.fs = service.(*fms.FileManagerService)
	return lf
}

func (cp *CP_MH) newFetchPartTask(descriptor *model2.FileDescriptor, part int) *FetchPartTask {
	fpt := &FetchPartTask{}
	fpt.descriptor = descriptor
	fpt.part = part
	fpt.cp = cp
	return fpt
}

func newCopyFileJobListener(filename string, parts int) *CopyFileJobListener {
	cpjl := &CopyFileJobListener{}
	cpjl.filename = filename
	cpjl.parts = parts
	cpjl.lastReport = time.Now().Unix()
	return cpjl
}

func (msgHandler *CP_MH) Init() {
}

func (msgHandler *CP_MH) Topic() string {
	return "cp"
}

func (msgHandler *CP_MH) Message(destination *ServiceID, data []byte, isReply bool) *Message {
	msg := msgHandler.fs.ServiceManager().NewMessage(msgHandler.Topic(), msgHandler.fs, destination, data, isReply)
	return msg
}

var total = 0

func (msgHandler *CP_MH) Handle(message *Message) {
	fileData := &model2.FileData{}
	fileData.Read(message.ByteSlice())
	fileData.LoadData()
	total += len(fileData.Data())
	msgHandler.fs.ServiceManager().Reply(message, fileData.ToBytes())
}

func (msgHandler *CP_MH) Request(args ...interface{}) (interface{}, error) {
	e := msgHandler.ValidateArgs(args)
	if e != nil {
		Println(e.Error())
		return nil, e
	}
	descriptor := args[0].(*model2.FileDescriptor)
	part := args[1].(int)

	fileData := model2.NewFileData(descriptor.SourcePath(), part, descriptor.Size())
	response, e := msgHandler.fs.ServiceManager().Request(msgHandler.Topic(), msgHandler.fs, msgHandler.fs.PeerServiceID(), fileData.ToBytes(), false)
	if e != nil {
		return nil, e
	}
	fileData.Read(NewByteSliceWithData(response, 0))
	return fileData, nil
}

func (msgHandler *CP_MH) ValidateArgs(args []interface{}) error {
	if len(args) != 2 {
		return errors.New("4 Args needed:path, dept, calcHash, peerID")
	}
	_, ok := args[0].(*model2.FileDescriptor)
	if !ok {
		return errors.New("descriptor must be *FileDescriptor")
	}
	_, ok = args[1].(int)
	if !ok {
		return errors.New("part must be int")
	}
	return nil
}

func (msgHandler *CP_MH) CopySmallFile(descriptor *model2.FileDescriptor) {
	d, e := msgHandler.Request(descriptor, 0)
	if e != nil {
		Println(e.Error())
		return
	}
	ioutil.WriteFile(descriptor.TargetPath(), d.(*model2.FileData).Data(), 777)
}

func (msgHandler *CP_MH) CopyLargeFile(descriptor *model2.FileDescriptor) {
	parts := descriptor.Parts()
	tasks := NewJob(5, newCopyFileJobListener(descriptor.TargetPath(), parts))
	for i := 0; i < parts; i++ {
		fpt := msgHandler.newFetchPartTask(descriptor, i)
		tasks.AddTask(fpt)
	}
	tasks.Run()
	assembleFile(descriptor)

}

func (msgHandler *CP_MH) CopyFile(descriptor *model2.FileDescriptor, asyncIfSmall bool) {
	parts := descriptor.Parts()
	if parts == 1 {
		if asyncIfSmall {
			go msgHandler.CopySmallFile(descriptor)
		} else {
			msgHandler.CopySmallFile(descriptor)
		}
	} else {
		msgHandler.CopyLargeFile(descriptor)
	}
}

func (task *FetchPartTask) Run() {
	d, e := task.cp.Request(task.descriptor, task.part)
	if e != nil {
		return
	}
	file, _ := os.Create(task.descriptor.TargetPath() + ".part-" + getPartString(task.part))
	file.Write(d.(*model2.FileData).Data())
	file.Close()
}

func (jl *CopyFileJobListener) Finished(task JobTask) {
	jl.finishCount++
	if time.Now().Unix()-jl.lastReport > 30 {
		p := float64(jl.finishCount) / float64(jl.parts) * 100
		Println("\n" + jl.filename + "(" + strconv.Itoa(int(p)) + "%).")
		jl.lastReport = time.Now().Unix()
	} else {
		Print(".")
	}
}

func getPartString(part int) string {
	str := strconv.Itoa(part)
	buff := bytes.Buffer{}
	for i := len(str); i < 5; i++ {
		buff.WriteString("0")
	}
	buff.WriteString(str)
	return buff.String()
}

func assembleFile(descriptor *model2.FileDescriptor) {
	parent := descriptor.TargetParent()
	if parent == nil {
		parent = descriptor.SourceParent()
	}

	dir := parent.TargetPath()
	filename := descriptor.TargetPath()

	filenames := make([]string, 0)
	files, _ := ioutil.ReadDir(dir)

	for _, f := range files {
		buff := bytes.Buffer{}
		buff.WriteString(dir)
		buff.WriteString("/")
		buff.WriteString(f.Name())
		part := buff.String()
		if isPartOfFile(part, filename) {
			filenames = append(filenames, part)
		}
	}

	sort.Slice(filenames, func(i, j int) bool {
		attValueA := filenames[i]
		attValueB := filenames[j]
		if attValueA < attValueB {
			return true
		}
		return false
	})

	file, _ := os.Create(filename)
	for _, fn := range filenames {
		data, _ := ioutil.ReadFile(fn)
		file.Write(data)
		os.Remove(fn)
	}
	file.Close()
}

func isPartOfFile(f, filename string) bool {
	if len(f) > len(filename) && f[0:len(filename)] == filename {
		return true
	}
	return false
}
