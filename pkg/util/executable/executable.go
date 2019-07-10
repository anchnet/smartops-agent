package executable

import (
	"github.com/kardianos/osext"
	"path/filepath"
)

func path() (string, error) {
	here, err := osext.Executable()
	if err != nil {
		return "", err
	}

	return filepath.EvalSymlinks(here)
}

// Folder returns the folder under which the executable is located,
// after having resolved all symlinks to the executable.
// Unlike os.Executable and osext.ExecutableFolder, Folder will
// resolve the symlinks across all platforms.
func Folder() (string, error) {
	p, err := path()
	if err != nil {
		return "", err
	}

	return filepath.Dir(p), nil
}
