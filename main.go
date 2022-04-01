package cmd

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

// Run returns the result of the command.
func Run(cmdCtx context.Context, param []string) ([]byte, error) {
	cmd := exec.CommandContext(cmdCtx, param[0], param[1:]...)
	if cmd.Stdout != nil {
		return nil, errors.New("exec: Stdout already set")
	}
	if cmd.Stderr != nil {
		return nil, errors.New("exec: Stderr already set")
	}

	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b

	// start
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	// wait
	waitErr := make(chan error)
	go func() {
		waitErr <- cmd.Wait()
		close(waitErr)
	}()

	select {
	case <-cmdCtx.Done():
		return nil, cmdCtx.Err()

	case err = <-waitErr:
		if err != nil {
			if cmd.Stderr != nil {
				err = errors.New(cmd.Stderr.(*bytes.Buffer).String())
			}
			return nil, err
		}
	}

	return b.Bytes(), err
}

// RunOut returns + outputs to the output stream the result of executing the command specified by the 'param' attribute.
func RunOut(cmdCtx context.Context, param []string) ([]byte, error) {
	cmd := exec.CommandContext(cmdCtx, param[0], param[1:]...)

	// out && err
	var stdOut, stdErr bytes.Buffer
	stdOutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	go cmdStd(cmdCtx, stdOutPipe, &stdOut, os.Stdout)

	stdErrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	go cmdStd(cmdCtx, stdErrPipe, &stdErr, os.Stderr)

	// start
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	// wait
	waitErr := make(chan error)
	go func() {
		waitErr <- cmd.Wait()
		close(waitErr)
	}()

	select {
	case <-cmdCtx.Done():
		return nil, cmdCtx.Err()

	case err = <-waitErr:
		if err != nil {
			if &stdErr != nil {
				err = errors.New(stdErr.String())
			}
			return nil, err
		}
	}

	return stdOut.Bytes(), nil
}

// cmdStd read from src (io.Reader) and write to dst (bytes.Buffer) and out (io.Writer).
func cmdStd(ctx context.Context, src io.Reader, dst *bytes.Buffer, out io.Writer) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	mw := io.MultiWriter(dst, out)

	b, err := ioutil.ReadAll(src)
	if err != nil {
		return
	}

	_, err = mw.Write(b)
	if err != nil {
		return
	}

	return
}
