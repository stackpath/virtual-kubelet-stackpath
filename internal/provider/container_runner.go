package provider

import (
	"context"
	"io"
	"strings"

	"github.com/virtual-kubelet/virtual-kubelet/node/api"
	"golang.org/x/crypto/ssh"
	"golang.org/x/sync/errgroup"
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

	ctx, cancel := context.WithCancel(ctx)
	g, ctx := errgroup.WithContext(ctx)

	//goroutine that is responsible for listening to the 'resize' channel and updating the session window measurements
	g.Go(func() error {
		for {
			select {
			case size := <-cr.attach.Resize():
				session.WindowChange(int(size.Height), int(size.Width))
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})

	aout := cr.attach.Stdout()
	if aout != nil {
		defer aout.Close()
	}
	g.Go(func() error {
		_, err = io.Copy(aout, sessionStdoutPipe)
		if err != nil {
			cancel()
			return err
		}
		return nil
	})

	g.Go(func() error {
		io.Copy(aout, sessionStderrPipe)
		return nil
	})

	ain := cr.attach.Stdin()

	if ain != nil {
		g.Go(func() error {
			io.Copy(sessionStdinPipe, ain)
			return nil
		})
	}

	// sending the command
	g.Go(func() error {
		err := session.Run(strings.Join(cmd, " ") + "\n")
		cancel()
		return err
	})

	err = g.Wait()
	return err
}
