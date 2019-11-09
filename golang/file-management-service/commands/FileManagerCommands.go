package commands

import (
	"github.com/saichler/console/golang/console/commands"
	commands2 "github.com/saichler/service-management-service/golang/management-service/commands"
	. "github.com/saichler/service-manager/golang/service-manager"
)

type FileManagerCommands struct {
	commands map[string]commands.Command
}

func (c *FileManagerCommands) Init(service IService, mh IServiceMessageHandlers) {
	c.commands = make(map[string]commands.Command)
	c.addCommand(NewPeers(service))
	c.addCommand(NewSet(service))
	c.addCommand(NewLS(service, mh.Handler("ls")))
	c.addCommand(NewCD(service, mh.Handler("ls")))
	c.addCommand(NewDiff(service, mh.Handler("ls")))
	c.addCommand(NewLCD(service, mh.Handler("ls")))
	c.addCommand(NewSync(service, mh.Handler("cp"), mh.Handler("ls")))
	c.addCommand(NewCpCMD(service, mh.Handler("ls"), mh.Handler("cp")))
	c.addAlias(commands2.NewCD(service), "..")
}

func (c *FileManagerCommands) addCommand(cmd commands.Command) {
	c.commands[cmd.Command()] = cmd
}

func (c *FileManagerCommands) addAlias(cmd commands.Command, alias string) {
	c.commands[alias] = cmd
}

func (c *FileManagerCommands) Commands() map[string]commands.Command {
	return c.commands
}
