package commands

import (
	"bytes"
	. "github.com/saichler/console/golang/console/commands"
	message_handlers2 "github.com/saichler/file-management-service/golang/file-management-service/message-handlers"
	"github.com/saichler/file-management-service/golang/file-management-service/service"
	. "github.com/saichler/service-manager/golang/service-manager"
)

type LS struct {
	service *service.FileManagerService
	lsMH    *message_handlers2.LS_MH
}

func NewLS(sm IService, mh IMessageHandler) *LS {
	sd := &LS{}
	sd.service = sm.(*service.FileManagerService)
	sd.lsMH = mh.(*message_handlers2.LS_MH)
	return sd
}

func (cmd *LS) Command() string {
	return "ls"
}

func (cmd *LS) Description() string {
	return "List the pre-set peer files at a pre-set location"
}

func (cmd *LS) Usage() string {
	return "ls"
}

func (cmd *LS) ConsoleId() *ConsoleId {
	return cmd.service.ConsoleId()
}

func (cmd *LS) RunCommand(args []string, id *ConsoleId) (string, *ConsoleId) {
	fd, e := cmd.lsMH.FetchDescriptor(cmd.service.PeerDir(), 1, false, cmd.service.PeerServiceID())
	if e != nil {
		return e.Error(), nil
	}
	buff := bytes.Buffer{}
	buff.WriteString("------------------------------------------------\n")
	buff.WriteString("Peer: " + cmd.service.PeerServiceID().String())
	buff.WriteString("\n")
	buff.WriteString("Directory: ")
	buff.WriteString(cmd.service.PeerDir())
	buff.WriteString("\n")

	for _, file := range fd.Files() {
		buff.WriteString("  - ")
		buff.WriteString(file.Name())
		buff.WriteString("\n")
	}
	id.SetSuffix(":" + cmd.service.PeerDir())
	return buff.String(), nil
}
