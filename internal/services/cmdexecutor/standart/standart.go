package standart

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/cmdexecutor"
)

type cmdExecutor struct {
	cmd  string
	args []string
}

func MakeStandardCmdExecutor(cmd string, args ...string) cmdexecutor.CmdExecutor {
	ret := &cmdExecutor{
		cmd:  cmd,
		args: []string{cmd},
	}
	ret.args = append(ret.args, args...)
	return ret
}

func (ce *cmdExecutor) Run(ctx context.Context) (string, error) {
	var out bytes.Buffer
	cmd := exec.Cmd{
		Path:   ce.cmd,
		Args:   ce.args,
		Stdout: &out,
		Stderr: &out,
	}
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("command executor run: %w", err)
	}
	waited := make(chan error)
	go func() {
		waited <- cmd.Wait()
	}()
	select {
	case <-ctx.Done():
		if err := cmd.Process.Signal(os.Interrupt); err != nil {
			cmd.Process.Kill()
		}
	case err := <-waited:
		if err != nil {
			return "", fmt.Errorf("command executor finished with an error: %w", err)
		}
	}

	return out.String(), nil
}
