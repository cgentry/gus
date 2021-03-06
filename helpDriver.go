// helpDriver will handle all of the output of help to the user.
package main

import (
	"fmt"
	"os"

	"github.com/cgentry/gdriver"
	"github.com/cgentry/gus/cli"
	"github.com/cgentry/gus/library/encryption"
	"github.com/cgentry/gus/library/storage"
)

var helpStore = &cli.Command{
	Name:      "store",
	UsageLine: "gus store [driver-name]",
	Short:     "Display a list of what drivers are available",
	Long: `
Display all of the drivers that are compiled into this runtime. If
you add in the 'driver-name', it will list specific help for that driver.

Each driver may require different paramters. The driver will give you some
details, but you should refer to the documentation
`,
}
var helpEncrypt = &cli.Command{
	Name:      "encrypt",
	UsageLine: "gus encrypt [driver-name]",
	Short:     "Display a list of what drivers are available",
	Long: `
Display all of the drivers that are compiled into this runtime. If
you add in the 'driver-name', it will list specific help for that driver.

Each driver may require different paramters. The driver will give you some
details, but you should refer to the documentation
`,
}

func init() {
	helpStore.Run = runStore
	helpEncrypt.Run = runEncrypt
}

// Output any help that is required
func runStore(cmd *cli.Command, args []string) {
	listStore := gdriver.ListMembers(storage.DriverGroup)

	if len(args) == 0 {
		cli.RenderTemplate(os.Stdout, templateStorageList, listStore)
		return
	}
	if len(args) == 1 {
		if entry, ok := listStore[args[0]]; ok {
			cli.RenderTemplate(os.Stdout, templateStorageEntry, entry)
			return
		}
		fmt.Fprintf(os.Stderr, "'%s' is not a valid storage driver\n", args[0])
	} else {
		fmt.Fprintf(os.Stderr, "Only one parameter for store command\nUse 'gus help store' for more information\n")
	}

	return
}

const templateStorageList = `
List of storage drivers available:{{ range . }}
  {{ .Identity(gdriver.IdentityName) }}: {{ .Identity( gdriver.IdentityShort) }}{{ end }}

`
const templateStorageEntry = `
{{ .Identity(gdriver.IdentityName) }}: {{ .Identity(gdriver.IdentityShort) }}
{{ .Identity(gdriver.IdentityLong) }}
`

func runEncrypt(cmd *cli.Command, args []string) {
	listStore := gdriver.ListMembers(encryption.DriverGroup)

	if len(args) == 0 {
		cli.RenderTemplate(os.Stdout, templateEncryptionList, listStore)
		return
	}
	if len(args) == 1 {
		if entry, ok := listStore[args[0]]; ok {
			cli.RenderTemplate(os.Stdout, templateEncryptionEntry, entry)
			return
		}
		fmt.Fprintf(os.Stderr, "'%s' is not a valid encryption driver\n", args[0])
	} else {
		fmt.Fprintf(os.Stderr, "Only one parameter for encrypt command\nUse 'gus help encrypt' for more information\n")
	}

	return
}

const templateEncryptionList = `
List of encryption drivers available:{{ range . }}
  {{ .Id }}: {{ .ShortHelp }}{{ end }}

`
const templateEncryptionEntry = `
{{ .Id }}: {{ .ShortHelp }}
{{ .LongHelp }}
`
