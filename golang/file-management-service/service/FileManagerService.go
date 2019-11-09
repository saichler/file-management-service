package service

import (
	"github.com/saichler/console/golang/console/commands"
	"github.com/saichler/messaging/golang/net/protocol"
	. "github.com/saichler/service-manager/golang/service-manager"
)

type FileManagerService struct {
	id       uint16
	cid      *commands.ConsoleId
	sm       *ServiceManager
	sid      *protocol.ServiceID
	peerSID  *protocol.ServiceID
	peerDir  string
	localDir string
}

func (fm *FileManagerService) Topic() string {
	return "File Manager"
}

func (fm *FileManagerService) ID() uint16 {
	return fm.id
}

func (fm *FileManagerService) ConsoleId() *commands.ConsoleId {
	return fm.cid
}

func (fm *FileManagerService) ServiceManager() *ServiceManager {
	return fm.sm
}

func (fm *FileManagerService) ServiceID() *protocol.ServiceID {
	return fm.sid
}

func (fm *FileManagerService) Init(sm *ServiceManager, id uint16, sid *protocol.ServiceID, cid *commands.ConsoleId) {
	fm.sm = sm
	fm.id = id
	fm.cid = cid
	fm.sid = sid
	fm.peerSID = sid
	fm.peerDir = "/tmp"
	fm.localDir = "/tmp"
}

func (fm *FileManagerService) PeerServiceID() *protocol.ServiceID {
	return fm.peerSID
}

func (fm *FileManagerService) SetPeerServiceID(sid *protocol.ServiceID) {
	fm.peerSID = sid
}

func (fm *FileManagerService) PeerDir() string {
	return fm.peerDir
}

func (fm *FileManagerService) SetPeerDir(dir string) {
	fm.peerDir = dir
}

func (fm *FileManagerService) LocalDir() string {
	return fm.localDir
}

func (fm *FileManagerService) SetLocalDir(l string) {
	fm.localDir = l
}
