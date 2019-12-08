package message_handlers

import (
	"errors"
	model2 "github.com/saichler/file-management-service/golang/file-management-service/model"
	fms "github.com/saichler/file-management-service/golang/file-management-service/service"
	. "github.com/saichler/messaging/golang/net/protocol"
	. "github.com/saichler/service-manager/golang/service-manager"
	. "github.com/saichler/utils/golang"
	"log"
)

type LS_MH struct {
	fs *fms.FileManagerService
}

func NewRlsMH(service IService) *LS_MH {
	lf := &LS_MH{}
	lf.fs = service.(*fms.FileManagerService)
	return lf
}

func (msgHandler *LS_MH) Init() {
}

func (msgHandler *LS_MH) Topic() string {
	return "ls"
}

func (msgHandler *LS_MH) Message(destination *ServiceID, data []byte, isReply bool) *Message {
	msg := msgHandler.fs.ServiceManager().NewMessage(msgHandler.Topic(), msgHandler.fs, destination, data, isReply)
	return msg
}

func (msgHandler *LS_MH) Handle(message *Message) {
	fr := &model2.FileRequest{}
	fr.FromBytes(message.Data())
	fd := model2.NewFileDescriptor(fr.Path(), fr.Dept(), fr.CalcHash())
	if fd == nil {
		fd = &model2.FileDescriptor{}
	}
	data := fd.ToBytes()
	msgHandler.fs.ServiceManager().Reply(message, data)
}

func (msgHandler *LS_MH) Request(args ...interface{}) (interface{}, error) {
	e := msgHandler.ValidateArgs(args)
	if e != nil {
		log.Println("Invalid Args: " + e.Error())
		return nil, e
	}
	path := args[0].(string)
	dept := args[1].(int)
	calcHash := args[2].(bool)
	peerID := args[3].(*ServiceID)

	if calcHash {
		Print("Fetching " + path + " Information with hashes, this may take a while...")
	} else {
		Print("Fetching " + path + " Information, this may take a min if its a directory...")
	}
	request := model2.NewFileRequest(path, dept, calcHash)
	response, e := msgHandler.fs.ServiceManager().Request(msgHandler.Topic(), msgHandler.fs, peerID, request.ToBytes(), false)
	Println("Done!")
	if e != nil {
		return nil, e
	}
	fd := model2.ReadFileDescriptor(NewByteSliceWithData(response, 0))
	return fd, nil
}

func (msgHandler *LS_MH) ValidateArgs(args []interface{}) error {
	if len(args) != 4 {
		return errors.New("4 Args needed:path, dept, calcHash, peerID")
	}
	_, ok := args[0].(string)
	if !ok {
		return errors.New("path must be string")
	}
	_, ok = args[1].(int)
	if !ok {
		return errors.New("dept must be int")
	}
	_, ok = args[2].(bool)
	if !ok {
		return errors.New("calcHash must be bool")
	}
	_, ok = args[3].(*ServiceID)
	if !ok {
		return errors.New("peerID must be *ServiceID")
	}
	return nil
}

func (msgHandler *LS_MH) FetchDescriptor(path string, dept int, calcHash bool, peerID *ServiceID) (*model2.FileDescriptor, error) {
	d, e := msgHandler.Request(path, dept, calcHash, peerID)
	if e != nil {
		return nil, e
	}
	return d.(*model2.FileDescriptor), e
}
