package app

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

//./agent copyfile --path ./peak.png --mode 0667 https://www.baidu.com/img/flexible/logo/pc/peak-result.png

var (
	copyfileCmd = &cobra.Command{
		Use:   "copyfile",
		Short: "copy file from remote server",
		RunE:  copyfile,
	}

	copyfilePath string
	copyfileMode string
)

func copyfile(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("have to have an argument")
	}

	if copyfilePath == "" {
		return errors.New("file path not exist")
	}

	err := downloadFile(copyfilePath, args[0])
	return err
}

func init() {
	copyfileCmd.Flags().StringVarP(&copyfilePath, "path", "p", "", "file path")
	copyfileCmd.Flags().StringVarP(&copyfileMode, "mode", "m", "", "file mode")
	Command.AddCommand(copyfileCmd)
}

func downloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	err = chmod(filepath, copyfileMode)
	return err
}

func chmod(filepath string, mode string) error {
	var fmode os.FileMode
	i, err := strconv.ParseUint(mode, 0, 32)
	if err != nil {
		fmode = 0666
	} else {
		fmode = os.FileMode(i)
	}
	err = os.Chmod(filepath, fmode)
	return err
}
