package task

import (
	"context"
	"io"
)

type Template interface {
	SetLogger(writer io.Writer)
	Start()
	Stop()
	GetContext() context.Context
}
