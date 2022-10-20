package flags

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bernylinville/wtf/cfg"
	goFlags "github.com/jessevdk/go-flags"
)

// Flags is the container for command line flag data
type Flags struct {
	Config  string `short:"c" long:"config" optional:"yes" description:"Path to config file"`
	Module  string `short:"m" long:"module" optional:"yes" description:"Display info about a specific module, i.e.: 'wtfutil -m=todo'"`
	Profile bool   `short:"p" long:"profile" optional:"yes" description:"Profile application memory usage"`
	Version bool   `short:"v" long:"version" description:"Show version info"`
	// Work-around go-flags misfeatures. If any sub-command is defined
	// then `wtf` (no sub-commands, the common usage), is warned about.
	Opt struct {
		Cmd  string   `positional-arg-name:"command"`
		Args []string `positional-arg-name:"args"`
	} `positional-args:"yes"` // 待解释

	hasCustom bool
}

var EXTRA = `
Commands:
  save-secret <service>
    service      Service URL or module name of secret.
  Save a secret into the secret store. The secret will be prompted for.
  Requires wtf.secretStore to be configured.  See individual modules for
  information on what service and secret means for their configuration,
  not all modules use secrets.
`

// NewFlags creates an instance of Flags
func NewFlags() *Flags {
	flags := Flags{}
	return &flags
}

// Parse parses the incoming flags
func (flags *Flags) Parse() {
	// goFlags.Default is a convenient default set of options which should cover most of the uses of the flags package.
	// Default = HelpFlag | PrintErrors | PassDoubleDash
	parser := goFlags.NewParser(&flags, goFlags.Default)
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*goFlags.Error); ok && flagsErr.Type == goFlags.ErrHelp {
			fmt.Println(EXTRA)
			os.Exit(0)
		}
	}

	// If we have a custom config, then we're done parsing parameters, we don't need to generate the default value
	flags.hasCustom = (len(flags.Config) > 0)
	if flags.hasCustom {
		return
	}

	// If no config file is explicitly passed in as a param then set the flag to the default config file
	configDir, err := cfg.WtfConfigDir()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	flags.Config = filepath.Join(configDir, "config.yml")
}
