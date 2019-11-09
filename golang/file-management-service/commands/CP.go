package commands

import (
	. "github.com/saichler/console/golang/console/commands"
	message_handlers2 "github.com/saichler/file-management-service/golang/file-management-service/message-handlers"
	model2 "github.com/saichler/file-management-service/golang/file-management-service/model"
	"github.com/saichler/file-management-service/golang/file-management-service/service"
	"github.com/saichler/security"
	. "github.com/saichler/service-manager/golang/service-manager"
	. "github.com/saichler/utils/golang"
	"os"
)

type CP struct {
	service *service.FileManagerService
	ls      *message_handlers2.LS_MH
	cp      *message_handlers2.CP_MH
}

func NewCpCMD(sm IService, ls, cp IMessageHandler) *CP {
	sd := &CP{}
	sd.service = sm.(*service.FileManagerService)
	sd.ls = ls.(*message_handlers2.LS_MH)
	sd.cp = cp.(*message_handlers2.CP_MH)
	return sd
}

func (cmd *CP) Command() string {
	return "cp"
}

func (cmd *CP) Description() string {
	return "Copy a file from the remote location"
}

func (cmd *CP) Usage() string {
	return "cp <remote path> <local path>"
}

func (cmd *CP) ConsoleId() *ConsoleId {
	return cmd.service.ConsoleId()
}

func (cmd *CP) RunCommand(args []string, id *ConsoleId) (string, *ConsoleId) {
	if len(args) < 2 {
		return cmd.Usage(), nil
	}
	rfilename := cmd.service.PeerDir() + "/" + args[0]

	descriptor, e := cmd.ls.FetchDescriptor(rfilename, 0, true, cmd.service.PeerServiceID())
	if e != nil {
		Println(e.Error())
		return "", nil
	}

	targetParent := model2.NewFileDescriptor(cmd.service.LocalDir(), 0, false)
	descriptor.SetTargetParent(targetParent)
	descriptor.SetTargetName(args[1])

	if _, err := os.Stat(descriptor.TargetPath()); !os.IsNotExist(err) {
		hash, _ := security.FileHash(descriptor.TargetPath())
		if hash == descriptor.Hash() {
			return "File " + descriptor.TargetPath() + " already exist in local dir", nil
		}
		Println("File already exist in local, overwrite (yes/no)?")
		resp, _ := Read()
		if resp != "yes" {
			return "Aborting", nil
		}
	}

	cmd.cp.CopyFile(descriptor, false)

	hash, _ := security.FileHash(descriptor.TargetPath())
	valid := hash == descriptor.Hash()
	if valid {
		return "Done!", nil
	}
	return "Corrupted", nil
}
