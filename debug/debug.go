package debug

import (
	"fmt"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

const rootNamespace = "content-sdk-go"

var echoLogger echo.Logger

// SetEchoLogger allows applications to inject an Echo logger used by this package
func SetEchoLogger(l echo.Logger) {
	echoLogger = l
}

func debug(debugModule string, format string, a ...any) {
	debug := os.Getenv("DEBUG")
	if debug == "true" || strings.Contains(debug, debugModule) {
		if echoLogger != nil {
			echoLogger.Debugf("[%s] %s", debugModule, fmt.Sprintf(format, a...))
			return
		}
		fmt.Printf("[%s] %s\n", debugModule, fmt.Sprintf(format, a...))
	}
}

func Common(format string, a ...any) {
	debug(rootNamespace+"/common", format, a...)
}

func Form(format string, a ...any) {
	debug(rootNamespace+"/form", format, a...)
}

func Http(format string, a ...any) {
	debug(rootNamespace+"/http", format, a...)
}

func Layout(format string, a ...any) {
	debug(rootNamespace+"/layout", format, a...)
}

func Dictionary(format string, a ...any) {
	debug(rootNamespace+"/dictionary", format, a...)
}

func Editing(format string, a ...any) {
	debug(rootNamespace+"/editing", format, a...)
}

func Sitemap(format string, a ...any) {
	debug(rootNamespace+"/sitemap", format, a...)
}

func Multisite(format string, a ...any) {
	debug(rootNamespace+"/multisite", format, a...)
}

func Robots(format string, a ...any) {
	debug(rootNamespace+"/robots", format, a...)
}

func Redirects(format string, a ...any) {
	debug(rootNamespace+"/redirects", format, a...)
}

func Locale(format string, a ...any) {
	debug(rootNamespace+"/locale", format, a...)
}

func ErrorPages(format string, a ...any) {
	debug(rootNamespace+"/metadata", format, a...)
}

func Proxy(format string, a ...any) {
	debug(rootNamespace+"/proxy", format, a...)
}
