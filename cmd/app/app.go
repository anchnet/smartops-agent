package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.51idc.com/smartops/smartops-agent/pkg/config"
	"os"
)

var (
	Command = &cobra.Command{
		Use:   fmt.Sprint("%s [command]", os.Args[0]),
		Short: "SmartOps Collector at your service.",
	}
	confFilePath string
)

const loggerName config.LoggerName = "CORE"

func init() {
	Command.PersistentFlags().StringVarP(&confFilePath, "cfgpath", "c", "", "path to directory containing smartcat.yaml")
}
