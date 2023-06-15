package mocks

import (
	"io"

	"github.com/virtual-kubelet/virtual-kubelet/node/api"
)

type MockExecIO struct {
	tty      bool
	input    io.Reader
	output   io.WriteCloser
	err      io.WriteCloser
	chResize chan api.TermSize
}

func NewMockExecIO(tty bool, input io.Reader, output io.WriteCloser, err io.WriteCloser, chResize chan api.TermSize) *MockExecIO {
	return &MockExecIO{
		tty:      tty,
		input:    input,
		output:   output,
		err:      err,
		chResize: chResize,
	}
}

func (e *MockExecIO) TTY() bool {
	return e.tty
}

func (e *MockExecIO) Stdin() io.Reader {
	return e.input
}

func (e *MockExecIO) Stdout() io.WriteCloser {
	return e.output
}

func (e *MockExecIO) Stderr() io.WriteCloser {
	return e.err
}

func (e *MockExecIO) Resize() <-chan api.TermSize {
	return e.chResize
}

type MockWriterCloser struct {
	Data     []byte
	isClosed bool
}

func (w *MockWriterCloser) Write(p []byte) (n int, err error) {
	w.Data = append(w.Data, p...)
	return
}

func (w *MockWriterCloser) Close() error {
	w.isClosed = true
	return nil
}
