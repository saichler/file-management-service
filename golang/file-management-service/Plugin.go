package file_management_service

import (
	"github.com/saichler/file-management-service/golang/file-management-service/commands"
	message_handlers "github.com/saichler/file-management-service/golang/file-management-service/message-handlers"
	"github.com/saichler/file-management-service/golang/file-management-service/service"
	. "github.com/saichler/service-manager/golang/service-manager"
)

var Service, Commands, Handlers = NewService()

func NewService() (IService, IServiceCommands, IServiceMessageHandlers) {
	s := &service.FileManagerService{}
	h := &message_handlers.FileManagerHandlers{}
	h.Init(s)
	c := &commands.FileManagerCommands{}
	c.Init(s, h)
	return s, c, h
}
