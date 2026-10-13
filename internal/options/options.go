package options

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/ayayaakasvin/localsend-auto-download/internal/domain"
)

func ParseFlags() (domain.FileDownloader, error) {
	if len(os.Args) < 2 {
		return nil, errors.New(InvalidCommand)
	}

	fs := flag.NewFlagSet("localsend-auto-download", flag.ExitOnError)

	downloaderFlagset := downloaderFlags{}
	downloaderFlagset.Register(fs)

	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	if downloaderFlagset.help {
		fs.SetOutput(os.Stdout)
		printHelp(fs)
		return nil, errors.New(HelpCommand)
	}

	if err := downloaderFlagset.Validate(); err != nil {
		return nil, err
	}

	return downloaderFlagset.Build()
}

func printHelp(fs *flag.FlagSet) {
	fmt.Fprintln(fs.Output(), `Usage:
  localsend-auto-download [flags]

Description:
  Downloads files shared by a LocalSend device on the local network.

Flags:`)
	fs.PrintDefaults()
	fmt.Fprintln(fs.Output(), `
Example:
  localsend-auto-download --host 192.168.1.20`)
}
