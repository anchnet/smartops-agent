package app

import (
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/config"
	"github.com/spf13/cobra"
	"os"
)

var (
	Command = &cobra.Command{
		Use:   fmt.Sprint("/opt/cloudops-agent/agent [command]", os.Args[0]),
		Short: "CloudOps Agent at your service.",
	}
	confFilePath string
)

const loggerName config.LoggerName = "CORE"

func init() {
	Command.PersistentFlags().StringVarP(&confFilePath, "cfgpath", "c", "", "path to directory containing cloudops.yaml")
}
