package commands

import (
	. "github.com/saichler/console/golang/console/commands"
	"github.com/saichler/file-management-service/golang/file-management-service/service"
	. "github.com/saichler/service-manager/golang/service-manager"
	"strings"
)

type LCD struct {
	service *service.FileManagerService
	mh      IMessageHandler
}

func NewLCD(sm IService, mh IMessageHandler) *LCD {
	sd := &LCD{}
	sd.service = sm.(*service.FileManagerService)
	sd.mh = mh
	return sd
}

func (cmd *LCD) Command() string {
	return "lcd"
}

func (cmd *LCD) Description() string {
	return "Change local directory"
}

func (cmd *LCD) Usage() string {
	return "lcd <dir>"
}

func (cmd *LCD) ConsoleId() *ConsoleId {
	return cmd.service.ConsoleId()
}

func (cmd *LCD) RunCommand(args []string, id *ConsoleId) (string, *ConsoleId) {

	if len(args) == 0 {
		return cmd.Usage(), nil
	}

	dir := args[0]

	if string(dir[0]) == "/" {
		cmd.service.SetLocalDir(args[0])
		return "", nil
	}
	if dir == ".." {
		path := cmd.service.LocalDir()
		index := strings.LastIndex(path, "/")
		if index > 0 {
			path = path[0:index]
			cmd.service.SetLocalDir(path)
		} else {
			path = "/"
			cmd.service.SetLocalDir(path)
		}
		return "", nil
	}
	return "", nil
}
