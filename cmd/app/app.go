package app

import (
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/config"
	"github.com/spf13/cobra"
	"os"
)

var (
	Command = &cobra.Command{
		Use:   fmt.Sprint("%s [command]", os.Args[0]),
		Short: "SmartOps Agent at your service.",
	}
	confFilePath string
)

const loggerName config.LoggerName = "CORE"

func init() {
	Command.PersistentFlags().StringVarP(&confFilePath, "cfgpath", "c", "", "path to directory containing smartops.yaml")
}
