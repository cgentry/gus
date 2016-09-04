package main

import (
	"github.com/cgentry/gus/cli"
	"github.com/cgentry/gus/library/encryption"
	"github.com/cgentry/gus/service/web"
)

var cmdService = &cli.Command{
	Name:      "service",
	UsageLine: "gus service [-c configfile]",
	Short:     "Run the program in service mode.",
	Long: `
Service will listen in on a port and wait for requests for user activity from
clients. Clients will call to register, authenticate, login and logout from
the system. Each request is made over HTTP but must use a PUT instead of a GET.

The single option, "-c" allows you to specify where to load the configuration
file from. The default configuration file is ` + DefaultConfigFilename + `.

`,
}

func init() {
	cmdService.Run = runService
	addCommonCommandFlags(cmdService)

}

func runService(cmd *cli.Command, args []string) {

	c, err := GetConfigFile()
	if err != nil {
		runtimeFail("Opening configuration file", err)
	}
	encryption.GetDriver(c.Encrypt.Name).Setup(c.Encrypt.Options)
	router := web.New(c)
	router.Register(web.RouteMap).Serve()

	return
}
