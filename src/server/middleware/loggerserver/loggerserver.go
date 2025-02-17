// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE.gin file.

package loggerserver

import (
	"errors"
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/logger"
	"github.com/SongZihuan/huan-proxy/src/server/responsewriter"
	"io"
	"net/http"
	"time"
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

type LogServer struct {
	skip      map[string]struct{}
	isTerm    bool
	logWriter func(msg string)
}

func NewLogServer() *LogServer {
	return &LogServer{
		skip:      make(map[string]struct{}, 10),
		isTerm:    logger.IsInfoTermNotDumb(),
		logWriter: logger.InfoWrite,
	}
}

// LoggerConfig defines the config for Logger middleware.
type LoggerConfig struct {
	// Optional. Default value is gin.defaultLogFormatter
	Formatter LogFormatter

	// Output is a writer where logs are written.
	// Optional. Default value is gin.DefaultWriter.
	Output io.Writer

	// SkipPaths is an url path array which logs are not written.
	// Optional.
	SkipPaths []string
}

// LogFormatter gives the signature of the formatter function passed to LoggerWithFormatter
type LogFormatter func(params LogFormatterParams) string

// LogFormatterParams is the structure any formatter will be handed when time to log comes
type LogFormatterParams struct {
	Request *http.Request

	// TimeStamp shows the time after the server returns a response.
	TimeStamp time.Time
	// StatusCode is HTTP response code.
	StatusCode int
	// Latency is how much time the server cost to process a certain request.
	Latency time.Duration
	// RemoteAddr equals Context's RemoteAddr method.
	RemoteAddr string
	// Method is the HTTP method given to the request.
	Method string
	// Path is a path the client requests.
	Path string
	// isTerm shows whether gin's output descriptor refers to a terminal.
	isTerm bool
	// BodySize is the size of the Response Body
	BodySize int64
	// Keys are the keys set on the request's context.
	Keys map[string]any
}

// StatusCodeColor is the ANSI color for appropriately logging http status code to a terminal.
func (p *LogFormatterParams) StatusCodeColor() string {
	code := p.StatusCode

	switch {
	case code >= http.StatusContinue && code < http.StatusOK:
		return white
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return white
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return yellow
	default:
		return red
	}
}

// MethodColor is the ANSI color for appropriately logging http method to a terminal.
func (p *LogFormatterParams) MethodColor() string {
	method := p.Method

	switch method {
	case http.MethodGet:
		return blue
	case http.MethodPost:
		return cyan
	case http.MethodPut:
		return yellow
	case http.MethodDelete:
		return red
	case http.MethodPatch:
		return green
	case http.MethodHead:
		return magenta
	case http.MethodOptions:
		return white
	default:
		return reset
	}
}

// ResetColor resets all escape attributes.
func (p *LogFormatterParams) ResetColor() string {
	return reset
}

func (ls *LogServer) Formatter(param LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if ls.isTerm {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}
	return fmt.Sprintf("[Huan-Proxy] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.RemoteAddr,
		methodColor, param.Method, resetColor,
		param.Path,
	)
}

func (ls *LogServer) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// Start timer
	serverErr := false // 若为true则为server error（http.StatusInternalServerError）
	start := time.Now()
	path := r.URL.Path
	raw := r.URL.RawQuery

	writer := responsewriter.NewResponseWriter(w)

	// Process request
	next(writer, r)

	err := writer.WriteToResponse()
	if err != nil && !errors.Is(err, responsewriter.ErrHasWriter) {
		serverErr = true
		writer.ServerError()
		// 请求发生服务器故障，日志服务继续
	}

	param := LogFormatterParams{
		Request: r,
		isTerm:  ls.isTerm,
		Keys:    make(map[string]any),
	}

	// Stop timer
	param.TimeStamp = time.Now()
	param.Latency = param.TimeStamp.Sub(start)

	param.RemoteAddr = r.RemoteAddr
	param.Method = r.Method
	if serverErr {
		param.StatusCode = http.StatusInternalServerError
		param.BodySize = 0
	} else {
		param.StatusCode = writer.Status()
		param.BodySize = writer.Size()
	}

	if raw != "" {
		path = path + "?" + raw
	}

	param.Path = path

	if ls.logWriter != nil {
		ls.logWriter(ls.Formatter(param))
	}
}
