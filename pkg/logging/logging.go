// MIT License

// Copyright (c) 2018 Masayuki Izumi

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package logging

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggingMode represents a logging configuration specification.
type LoggingMode int

// LoggingMode values
const (
	LoggingNop LoggingMode = iota
	LoggingVerbose
	LoggingDebug
)

var (
	logging = LoggingNop

	MyCapicalColorLevelEncoder = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		var colorPrefix string
		switch l {
		case zapcore.InfoLevel:
			colorPrefix = "[white:blue]"
		case zapcore.WarnLevel:
			colorPrefix = "[white:yellow]"
		case zapcore.ErrorLevel:
			colorPrefix = "[white:red]"
		case zapcore.FatalLevel:
			colorPrefix = "[white:red]"
		default:
			colorPrefix = ""
		}
		enc.AppendString(colorPrefix + l.CapitalString() + "[-:-:-]")
	}
)

// AddLoggingFlags sets "--debug" and "--verbose" flags to the given *cobra.Command instance.
func AddLoggingFlags(cmd *cobra.Command) {
	var (
		debugEnabled, verboseEnabled bool
	)

	cmd.PersistentFlags().BoolVar(
		&debugEnabled,
		"debug",
		false,
		fmt.Sprintf("Debug level output"),
	)
	cmd.PersistentFlags().BoolVarP(
		&verboseEnabled,
		"verbose",
		"v",
		false,
		fmt.Sprintf("Verbose level output"),
	)

	cobra.OnInitialize(func() {
		switch {
		case debugEnabled:
			Debug()
		case verboseEnabled:
			Verbose()
		}
	})
}

// Debug sets a debug logger in global.
func Debug() {
	logging = LoggingDebug
	l := newDebugLogger(os.Stderr)
	replaceLogger(l)
}

// Verbose sets a verbose logger in global.
func Verbose() {
	logging = LoggingVerbose
	l := newVerboseLogger(os.Stderr)
	replaceLogger(l)
}

// IsDebug returns true if a debug logger is used.
func IsDebug() bool { return logging == LoggingDebug }

// IsVerbose returns true if a verbose logger is used.
func IsVerbose() bool { return logging == LoggingVerbose }

// Logging returns a current logging mode.
func Logging() LoggingMode { return logging }

func newDebugLogger(w io.Writer) *zap.Logger {
	encoder := zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, zapcore.AddSync(w), zapcore.DebugLevel)

	return zap.New(core)
}

func newVerboseLogger(w io.Writer) *zap.Logger {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = MyCapicalColorLevelEncoder
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	core := zapcore.NewCore(encoder, zapcore.AddSync(w), zapcore.InfoLevel)
	return zap.New(core)
}

func replaceLogger(l *zap.Logger) {
	l.Sync()
	zap.ReplaceGlobals(l)
}

func SetLoggerOutputToTview(tview io.Writer) {
	if IsDebug() {
		replaceLogger(newDebugLogger(tview))
	} else if IsVerbose() {
		replaceLogger(newVerboseLogger(tview))
	}
}
