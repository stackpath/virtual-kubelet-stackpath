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
		modes := ssh.TerminalModes{
			ssh.ECHO:    0,
			ssh.ECHOCTL: 0,
		}
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

	// Goroutine responsible for listening to the 'resize' channel, updating the session
	// window measurements, and handling session closure if the context is canceled
	// or the terminal window is closed by the user.
	g.Go(func() error {
		for {
			select {
			case size := <-cr.attach.Resize():
				// If the height and width are both 0, it likely indicates that the terminal
				// window has been closed by the user. In this case, the SSH session is closed
				// to terminate any running commands.
				if size.Height == 0 && size.Width == 0 {
					session.Close()
					return nil
				}
				session.WindowChange(int(size.Height), int(size.Width))
			case <-ctx.Done():
				// If the context is canceled, the SSH session is closed to terminate any
				// running commands.
				session.Close()
				return ctx.Err()
			}
		}
	})

	if aout := cr.attach.Stdout(); aout != nil {
		defer aout.Close()
		g.Go(func() error {
			io.Copy(aout, sessionStdoutPipe)
			return nil
		})
	}

	if aerr := cr.attach.Stderr(); aerr != nil {
		defer aerr.Close()
		g.Go(func() error {
			io.Copy(aerr, sessionStderrPipe)
			return nil
		})
	}

	if ain := cr.attach.Stdin(); ain != nil {
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
