package commands

import (
	. "github.com/saichler/console/golang/console/commands"
	"github.com/saichler/file-management-service/golang/file-management-service/service"
	. "github.com/saichler/service-manager/golang/service-manager"
	"strings"
)

type CD struct {
	service *service.FileManagerService
	mh      IMessageHandler
}

func NewCD(sm IService, mh IMessageHandler) *CD {
	sd := &CD{}
	sd.service = sm.(*service.FileManagerService)
	sd.mh = mh
	return sd
}

func (cmd *CD) Command() string {
	return "cd"
}

func (cmd *CD) Description() string {
	return "Change directory"
}

func (cmd *CD) Usage() string {
	return "cd <dir>"
}

func (cmd *CD) ConsoleId() *ConsoleId {
	return cmd.service.ConsoleId()
}

func (cmd *CD) RunCommand(args []string, id *ConsoleId) (string, *ConsoleId) {

	if len(args) == 0 {
		return cmd.Usage(), nil
	}

	dir := args[0]

	if string(dir[0]) == "/" {
		cmd.service.SetPeerDir(args[0])
		id.SetSuffix(":" + dir)
		return "", nil
	}
	if dir == ".." {
		path := cmd.service.PeerDir()
		index := strings.LastIndex(path, "/")
		if index > 0 {
			path = path[0:index]
			cmd.service.SetPeerDir(path)
			id.SetSuffix(":" + path)
		} else {
			path = "/"
			cmd.service.SetPeerDir(path)
			id.SetSuffix(":" + path)
		}
		return "", nil
	}
	return "", nil
}
