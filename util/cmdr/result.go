package cmdr

import (
	"encoding/base64"
	"fmt"
	"unicode/utf8"

	sous "github.com/opentable/sous/lib"
)

type (
	// Result is a result of a CLI invokation.
	Result interface {
		// ExitCode is the exit code the program should exit with based on this
		// result.
		SetTraceID(traceID sous.TraceID)
		GetTraceID() sous.TraceID
		ExitCode() int
	}
	Tipper interface {
		UserTip() string
	}
	// SuccessResult is a successful result.
	SuccessResult struct {
		// Data is the real return value of this function, it will be printed to
		// stdout by default, for consumption by other commands/pipelines etc.
		Data    []byte
		TraceID sous.TraceID
	}
)

func (s *SuccessResult) GetTraceID() sous.TraceID        { return s.TraceID }
func (s *SuccessResult) SetTraceID(traceID sous.TraceID) { s.TraceID = traceID }
func (s *SuccessResult) ExitCode() int                   { return EX_OK }

func (s SuccessResult) String() string {
	if utf8.Valid(s.Data) {
		return string(s.Data)
	}
	return base64.StdEncoding.EncodeToString(s.Data)
}

func Success(v ...interface{}) Result {
	return &SuccessResult{Data: []byte(fmt.Sprintln(v...))}
}

func SuccessData(d []byte) Result { return &SuccessResult{Data: d} }

func Successf(format string, v ...interface{}) Result {
	return &SuccessResult{Data: []byte(fmt.Sprintf(format+"\n", v...))}
}
