package commands

import (
	. "github.com/saichler/console/golang/console/commands"
	message_handlers2 "github.com/saichler/file-management-service/golang/file-management-service/message-handlers"
	model2 "github.com/saichler/file-management-service/golang/file-management-service/model"
	"github.com/saichler/file-management-service/golang/file-management-service/service"
	"github.com/saichler/security"
	. "github.com/saichler/service-manager/golang/service-manager"
	. "github.com/saichler/utils/golang"
	"log"
	"os"
	"strconv"
)

type Sync struct {
	service     *service.FileManagerService
	cp          *message_handlers2.CP_MH
	ls          *message_handlers2.LS_MH
	running     bool
	currentSize int64
}

func NewSync(sm IService, cp IMessageHandler, ls IMessageHandler) *Sync {
	sd := &Sync{}
	sd.service = sm.(*service.FileManagerService)
	sd.cp = cp.(*message_handlers2.CP_MH)
	sd.ls = ls.(*message_handlers2.LS_MH)
	return sd
}

func (cmd *Sync) Command() string {
	return "sync"
}

func (cmd *Sync) Description() string {
	return "Sync remote direcory to local."
}

func (cmd *Sync) Usage() string {
	return "sync"
}

func (cmd *Sync) ConsoleId() *ConsoleId {
	return cmd.service.ConsoleId()
}

func (cmd *Sync) RunCommand(args []string, id *ConsoleId) (string, *ConsoleId) {
	cmd.running = true
	descriptor, e := cmd.ls.FetchDescriptor(cmd.service.PeerDir(), 100, false, cmd.service.PeerServiceID())
	if e != nil {
		return e.Error(), nil
	}
	files, dirs := countFilesAndDirectories(descriptor)
	log.Println("Going to sync " + strconv.Itoa(files) + " files in " + strconv.Itoa(dirs) + " directories.")
	log.Print("Creating " + strconv.Itoa(dirs) + " on local target...")
	cmd.createDirectories(descriptor)
	log.Println("Done!")

	log.Println("Start downloading new files...")
	sr := model2.NewSyncReport()
	cmd.copyFiles(descriptor, sr)
	if !cmd.running {
		return "", nil
	}
	report := sr.Report(true)
	Println(report)
	Println("Start downloading size diff files...")
	sideDiff := sr.SizeDiff()
	sr = model2.NewSyncReport()
	for _, d := range sideDiff {
		cmd.copyFile(d, sr, true)
		if !cmd.running {
			return "", nil
		}
	}
	report = sr.Report(true)
	Println(report)
	return "Done", nil
}

func (cmd *Sync) copyFiles(descriptor *model2.FileDescriptor, sr *model2.SyncReport) {
	if descriptor.IsDir() {
		for _, file := range descriptor.Files() {
			cmd.copyFiles(file, sr)
			if !cmd.running {
				return
			}
		}
		return
	}

	if descriptor.Name() == "" {
		sr.AddErrored("File does not exist.", descriptor)
		return
	} else if descriptor.Size() == 0 {
		sr.AddErrored("File size is 0", descriptor)
		return
	}

	if !cmd.running {
		return
	}

	exists, err := os.Stat(descriptor.TargetPath())

	if !os.IsNotExist(err) {
		if descriptor.Size() == exists.Size() {
			sr.AddExist(descriptor)
			return
		} else {
			sr.AddSizeDiff(descriptor)
			return
		}
		/*
			request := model.NewFileRequest(descriptor.SourcePath(), 1, true)
			response := cmd.ls.Request(request, cmd.service.PeerServiceID())
			descriptor.SetHash(response.(*model.FileDescriptor).Hash())
			hash, _ := security.FileHash(descriptor.TargetPath())
			if hash == descriptor.Hash() {
				return
			}*/
	}

	cmd.copyFile(descriptor, sr, false)
	report := sr.Report(false)
	if report != "" {
		log.Println(report)
	}
}

func (cmd *Sync) copyFile(descriptor *model2.FileDescriptor, sr *model2.SyncReport, forceCopy bool) {
	if !cmd.running {
		return
	}

	Print(descriptor.TargetPath() + " (" + strconv.Itoa(int(descriptor.Size())) + "): ")

	if _, err := os.Stat(descriptor.TargetPath()); !forceCopy && !os.IsNotExist(err) {
		hash, _ := security.FileHash(descriptor.TargetPath())
		if hash == descriptor.Hash() {
			sr.AddExist(descriptor)
			return
		}
	}

	cmd.cp.CopyFile(descriptor, true)

	sr.AddCopied(descriptor)

	if descriptor.Hash() != "" {
		hash, _ := security.FileHash(descriptor.TargetPath())
		valid := hash == descriptor.Hash()
		if valid {
			Println("Done!")
		} else {
			Println("Corrupted!")
		}
	} else {
		Println("Done!")
	}
}

func (cmd *Sync) createDirectories(descriptor *model2.FileDescriptor) {
	local := model2.NewFileDescriptor(cmd.service.LocalDir(), 0, false)
	descriptor.SetTargetParent(local)
	createDirectories(descriptor)
}

func createDirectories(descriptor *model2.FileDescriptor) {
	if descriptor.Files() == nil || len(descriptor.Files()) == 0 {
		return
	}
	_, e := os.Stat(descriptor.TargetPath())
	if os.IsNotExist(e) {
		os.MkdirAll(descriptor.TargetPath(), 0777)
	}
	for _, child := range descriptor.Files() {
		createDirectories(child)
	}
}

func countFilesAndDirectories(fileDescriptor *model2.FileDescriptor) (int, int) {
	if fileDescriptor.Files() == nil || len(fileDescriptor.Files()) == 0 {
		return 1, 0
	}
	dirs := 1
	files := 0
	for _, child := range fileDescriptor.Files() {
		f, d := countFilesAndDirectories(child)
		dirs += d
		files += f
	}
	return files, dirs
}

func (cmd *Sync) Stop() {
	cmd.running = false
	Println("Stop signal was sent")
}
