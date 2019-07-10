package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	CaptureCmd = &cobra.Command{
		Use:   fmt.Sprint("%s [command]", os.Args[0]),
		Short: "SmartOps Collector at your service.",
	}
	confFilePath string
)

func init() {
	CaptureCmd.PersistentFlags().StringVarP(&confFilePath, "cfgpath", "c", "", "path to directory containing smartcat.yaml")
}
