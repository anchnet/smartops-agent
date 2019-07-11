package config

import (
	"fmt"
	//"gitlab.51idc.com/smartops/smartops-agent/pkg/util/log"
	//log "github.com/cihub/seelog"
	"strings"

	"github.com/cihub/seelog"
)

// LoggerName specifies the name of an instantiated logger.
type LoggerName string

const logDateFormat = "2006-01-02 15:04:05 MST" // see time.Format for format syntax

// BuildCommonFormat returns the log common format seelog string
func BuildCommonFormat(loggerName LoggerName) string {
	return fmt.Sprintf("%%Date(%s) | %s | %%LEVEL | (%%ShortFilePath:%%Line in %%FuncShort) | %%Msg%%n", logDateFormat, loggerName)
}

// BuildJSONFormat returns the log JSON format seelog string
func BuildJSONFormat(loggerName LoggerName) string {
	return fmt.Sprintf("{&quot;agent&quot;:&quot;%s&quot;,&quot;time&quot;:&quot;%%Date(%s)&quot;,&quot;level&quot;:&quot;%%LEVEL&quot;,&quot;file&quot;:&quot;%%ShortFilePath&quot;,&quot;line&quot;:&quot;%%Line&quot;,&quot;func&quot;:&quot;%%FuncShort&quot;,&quot;msg&quot;:&quot;%%Msg&quot;}%%n", strings.ToLower(string(loggerName)), logDateFormat)
}

// SetupLogger sets up a logger with the specified logger name and log level
// if a non empty logFile is provided, it will also log to the file
// a non empty syslogURI will enable syslog, and format them following RFC 5424 if specified
// you can also specify to log to the console and in JSON format
func SetupLogger(loggerName LoggerName, logLevel, logFile string, logToConsole, jsonFormat bool) error {
	seelogLogLevel := strings.ToLower(logLevel)

	configTemplate := fmt.Sprintf(`<seelog minlevel="%s">`, seelogLogLevel)

	formatID := "common"
	if jsonFormat {
		formatID = "json"
	}

	configTemplate += fmt.Sprintf(`<outputs formatid="%s">`, formatID)

	if logToConsole {
		configTemplate += `<console />`
	}
	if logFile != "" {
		configTemplate += fmt.Sprintf(`<rollingfile type="size" filename="%s" maxsize="%d" maxrolls="%d" />`, logFile, Smartcat.GetSizeInBytes("log_file_max_size"), Smartcat.GetInt("log_file_max_rolls"))
	}
	configTemplate += fmt.Sprintf(`</outputs>
	<formats>
		<format id="json" format="%s"/>
		<format id="common" format="%s"/>
	</formats>
</seelog>`,
		BuildJSONFormat(loggerName),
		BuildCommonFormat(loggerName),
		strings.ToLower(string(loggerName)),
		loggerName,
	)

	logger, err := seelog.LoggerFromConfigAsString(configTemplate)
	if err != nil {
		return err
	}
	seelog.ReplaceLogger(logger)

	return nil
}

func parseShortFilePath(params string) seelog.FormatterFunc {
	return func(message string, level seelog.LogLevel, context seelog.LogContextInterface) interface{} {
		return extractShortPathFromFullPath(context.FullPath())
	}
}

func extractShortPathFromFullPath(fullPath string) string {
	// We want to trim the part containing the path of the project
	// ie DataDog/smartops-agent/
	slices := strings.Split(fullPath, "-agent/")
	return slices[len(slices)-1]
}

func init() {
	seelog.RegisterCustomFormatter("ShortFilePath", parseShortFilePath)
}
