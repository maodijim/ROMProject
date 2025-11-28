package task

import (
	"io"
)

type Template interface {
	SetLogger(writer io.Writer)
	Start()
	Stop()
}
