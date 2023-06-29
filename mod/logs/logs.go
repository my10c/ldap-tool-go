//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package logs

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

type (
	LogConfig struct {
		LogsDir       string
		LogFile       string
		LogMaxSize    int
		LogMaxBackups int
		LogMaxAge     int
	}
)

var (
	prefixes = []string{"main."}

	removeLines = strings.NewReplacer(
		"\r\n", "\\r\\n",
		"\r", "\\r",
		"\n", "\\n")
)

func tidy(s string) string {
	return strings.TrimSpace(removeLines.Replace(s))
}

// get the caller function
func funcCaller(depth int) string {
	pc, _, _, ok := runtime.Caller(depth)
	info := runtime.FuncForPC(pc)
	if ok && info != nil {
		fname := runtime.FuncForPC(pc).Name()
		pos := strings.LastIndex(fname, "/")
		if pos > 1 {
			fname = fname[pos+1:]
		}
		fname = strings.TrimSuffix(fname, ".0")
		for _, prefix := range prefixes {
			if strings.HasPrefix(fname, prefix) {
				for depth > 3 {
					depth--
					pc, _, _, ok := runtime.Caller(depth)
					info := runtime.FuncForPC(pc)
					if !ok || info == nil {
						break
					}
					f := runtime.FuncForPC(pc).Name()
					pos := strings.LastIndex(f, "/")
					if pos > 1 {
						f = f[pos+1:]
					}
					if f == "" {
						break
					}
					fname = fname + " > " + f
				}
				return "(" + fname + ")"
			}
		}
	}
	return ""
}

// Get the function name in the call stack that matches a prefix.
func funcName() string {
	flast := funcCaller(3)
	for i := 4; i < 7; i++ {
		fname := funcCaller(i)
		if fname != "" {
			return fname
		}
	}
	return flast
}

func Log(msg string, level string) {
	var b strings.Builder
	b.Grow(128)
	b.WriteString(fmt.Sprintf("%s : ", level))
	b.WriteString(funcName())
	b.WriteString(" ")
	b.WriteString(tidy(msg))
	log.Printf("%s", b.String())
}

// func InitLogs(config *configurator.Config) {
func InitLogs(c *LogConfig) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(&lumberjack.Logger{
		Filename:   c.LogsDir + "/" + c.LogFile,
		MaxSize:    c.LogMaxSize,
		MaxBackups: c.LogMaxBackups,
		MaxAge:     c.LogMaxAge,
		Compress:   true,
	})
}
