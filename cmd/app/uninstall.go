// +build !windows

package app

import (
	"fmt"
	"github.com/anchnet/smartops-agent/cmd/common"
	"github.com/anchnet/smartops-agent/pkg/pidfile"
	"github.com/anchnet/smartops-agent/pkg/util/file"
	"github.com/spf13/cobra"
	"os"
	"runtime"
	"syscall"
)

var (
	uninstallCmd = &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall SmartOps Agent",
		RunE:  uninstall,
	}
)

func uninstall(cmd *cobra.Command, args []string) error {
	if runtime.GOOS == "windows" {
		return fmt.Errorf("unsupported operator system %v", runtime.GOOS)
	}

	pid := pidfile.ReadPID(common.RootDirectory + "/" + common.DefaultPidFile)
	if pid > 0 {
		fmt.Println("Stopping agent...")
		if err := syscall.Kill(pid, syscall.SIGKILL); err != nil {
			return fmt.Errorf("stop agent failed, %v", err)
		}
	}

	fmt.Println("Delete agent")
	if err := os.RemoveAll(common.RootDirectory); err != nil {
		return fmt.Errorf("delete agent failed, %v", err)
	}

	if file.IsExist(common.DefaultSystemdPath) {
		fmt.Println("Delete systemd config")
		if err := os.Remove(common.DefaultSystemdPath); err != nil {
			return fmt.Errorf("delete agent systemd config failed, %v", err)
		}
	}

	if file.IsExist(common.DefaultUpstartPath) {
		fmt.Println("Delete upstart config")
		if err := os.Remove(common.DefaultUpstartPath); err != nil {
			return fmt.Errorf("delete agent upstart confi failed, %v", err)
		}
	}
	fmt.Println("Uninstall success")

	return nil
}
func init() {
	Command.AddCommand(uninstallCmd)
}
