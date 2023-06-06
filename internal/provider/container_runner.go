package provider

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/virtual-kubelet/virtual-kubelet/node/api"
	"golang.org/x/crypto/ssh"
)

type ContainerRunner struct {
	client *ssh.Client
	attach api.AttachIO
}

// NewContainerRunner creates a new ContainerRunner object.
func NewContainerRunner(sshClient *ssh.Client, attachIO api.AttachIO) *ContainerRunner {
	return &ContainerRunner{
		attach: attachIO,
		client: sshClient,
	}
}

// Exec executes a command on a remote server via SSH protocol.
func (cr *ContainerRunner) Exec(ctx context.Context, cmd []string) error {
	session, err := cr.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	if cr.attach.TTY() {
		// Set up terminal modes
		modes := ssh.TerminalModes{}
		err = session.RequestPty("Xterm", 120, 60, modes)
		if err != nil {
			return err
		}
	}

	sessionStdinPipe, err := session.StdinPipe()
	if err != nil {
		return err
	}
	defer sessionStdinPipe.Close()

	sessionStdoutPipe, err := session.StdoutPipe()
	if err != nil {
		return err
	}

	sessionStderrPipe, err := session.StderrPipe()
	if err != nil {
		return err
	}

	// goroutine that is responsible for listening 'resize' channel and update the session window measurements accordingly
	go func() {
		for {
			select {
			case size := <-cr.attach.Resize():
				session.WindowChange(int(size.Height), int(size.Width))
			case <-ctx.Done():
				return
			default:
				time.Sleep(time.Millisecond * 50)
			}
		}
	}()

	// Channel that is used for sending errors happened during stdout pipe read.
	e := make(chan error, 1)

	aout := cr.attach.Stdout()
	if aout != nil {
		defer aout.Close()
	}
	go func() {
		_, err := io.Copy(aout, sessionStdoutPipe)
		if err != nil {
			// io.EOF or an error
			e <- err
			return
		}
	}()
	go func() { io.Copy(aout, sessionStderrPipe) }()

	ain := cr.attach.Stdin()
	if ain != nil {
		go func() { io.Copy(sessionStdinPipe, ain) }()
	}

	done := false
	// sending the command
	go func() {
		if err := session.Run(strings.Join(cmd, " ") + "\n"); err != nil {
			e <- err
			return
		}
		done = true
	}()

mainLoop:
	for {
		if done {
			break
		}
		select {
		case <-ctx.Done():
			break mainLoop
		case err = <-e:
			break mainLoop
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}

	return err
}
