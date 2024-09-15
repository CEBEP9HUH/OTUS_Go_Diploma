package cmdexecutor

import (
	"context"
)

type CmdExecutor interface {
	Run(ctx context.Context) (result string, err error)
}
