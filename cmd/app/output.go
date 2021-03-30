package app

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"

	"github.com/anchnet/smartops-agent/pkg/executor"
	"github.com/anchnet/smartops-agent/pkg/http"
	"github.com/spf13/cobra"
)

var (
	outputCmd = &cobra.Command{
		Use:   "output",
		Short: "Stand output",
		RunE:  output,
	}
	outputFlagT  string
	outputFlagID string
)

func output(cmd *cobra.Command, args []string) error {
	if runtime.GOOS == "windows" {
		return fmt.Errorf("unsupported operator system %v", runtime.GOOS)
	}

	opf := executor.CustomMonitorCmdRet{
		ID:      outputFlagID,
		Type:    outputFlagT,
		Message: strings.Join(args, ""),
	}
	reqByte, err := json.Marshal(opf)
	if err != nil {
		fmt.Println(`{"id":"0","message":"err","type":"n"}`)
		return err
	}
	respByte, err := http.LocalMetric(reqByte)
	if err != nil {
		fmt.Println("Send local metric server error: ", err)
		return err
	}
	fmt.Println(string(respByte))
	return nil
}

func init() {
	outputCmd.Flags().StringVarP(&outputFlagT, "type", "t", "", "type s/n.  s: string, n: num")
	outputCmd.Flags().StringVarP(&outputFlagID, "id", "", "", "id")
	Command.AddCommand(outputCmd)
}
