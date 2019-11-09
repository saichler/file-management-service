package commands

import (
	"bytes"
	. "github.com/saichler/console/golang/console/commands"
	message_handlers2 "github.com/saichler/file-management-service/golang/file-management-service/message-handlers"
	model2 "github.com/saichler/file-management-service/golang/file-management-service/model"
	"github.com/saichler/file-management-service/golang/file-management-service/service"
	. "github.com/saichler/service-manager/golang/service-manager"
	utils "github.com/saichler/utils/golang"
)

type Diff struct {
	service *service.FileManagerService
	ls      *message_handlers2.LS_MH
}

func NewDiff(sm IService, ls IMessageHandler) *Diff {
	sd := &Diff{}
	sd.service = sm.(*service.FileManagerService)
	sd.ls = ls.(*message_handlers2.LS_MH)
	return sd
}

func (cmd *Diff) Command() string {
	return "diff"
}

func (cmd *Diff) Description() string {
	return "Compare remote and local dirs"
}

func (cmd *Diff) Usage() string {
	return "diff"
}

func (cmd *Diff) ConsoleId() *ConsoleId {
	return cmd.service.ConsoleId()
}

func (cmd *Diff) RunCommand(args []string, id *ConsoleId) (string, *ConsoleId) {
	aside, e := cmd.ls.FetchDescriptor(cmd.service.PeerDir(), 100, false, cmd.service.PeerServiceID())
	if e != nil {
		utils.Println(e.Error())
		return "", nil
	}
	zside := model2.NewFileDescriptor(cmd.service.LocalDir()+"/"+aside.Name(), 100, false)
	aSideMissing, zSideMissing := diff(aside, zside)
	buff := bytes.Buffer{}
	buff.WriteString("Remote missing files:\n")
	for name, name2 := range aSideMissing {
		buff.WriteString(name)
		buff.WriteString(" - ")
		buff.WriteString(name2)
		buff.WriteString("\n")
	}
	buff.WriteString("Local missing files:\n")
	for name, name2 := range zSideMissing {
		buff.WriteString(name)
		buff.WriteString(" - ")
		buff.WriteString(name2)
		buff.WriteString("\n")
	}
	return buff.String(), nil
}

func diff(aside, zside *model2.FileDescriptor) (map[string]string, map[string]string) {
	zsideTarget := model2.NewFileDescriptor(aside.SourceParent().SourcePath(), 0, false)
	asideTarget := model2.NewFileDescriptor(zside.SourceParent().SourcePath(), 0, false)
	aside.SetTargetParent(asideTarget)
	zside.SetTargetParent(zsideTarget)

	aSideMissing := make(map[string]string)
	zSideMissing := make(map[string]string)
	deepDiff(aside, zside, aside.SourceRoot(), zside.SourceRoot(), aSideMissing, zSideMissing)
	return aSideMissing, zSideMissing
}

func deepDiff(aside, zside, aSideRoot, zSideRoot *model2.FileDescriptor, aSideMissing, zSideMissing map[string]string) {
	if aside == nil && zside != nil {
		aSideMissing[zside.SourcePath()] = zside.TargetPath()
		return
	}
	if aside != nil && zside == nil {
		zSideMissing[aside.SourcePath()] = aside.TargetPath()
		return
	}

	if aside.IsDir() {
		for _, aSideChild := range aside.Files() {
			path := aSideChild.TargetPath()
			zSideChild := zSideRoot.Get(path)
			deepDiff(aSideChild, zSideChild, aSideRoot, zSideRoot, aSideMissing, zSideMissing)
		}
	}

	if zside.IsDir() {
		for _, zSideChild := range zside.Files() {
			path := zSideChild.TargetPath()
			aSideChild := aSideRoot.Get(path)
			deepDiff(aSideChild, zSideChild, aSideRoot, zSideRoot, aSideMissing, zSideMissing)
		}
	}
}
