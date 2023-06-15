package sshtest

import (
	"io"

	"github.com/gliderlabs/ssh"
)

type TestSSHServer struct {
	server *ssh.Server
}

// TestSSHServer creates a new TestSSHServer.
func NewTestSSHServer(addr string, username, pass string) *TestSSHServer {
	return &TestSSHServer{
		server: &ssh.Server{
			Addr: addr,
			PasswordHandler: func(ctx ssh.Context, password string) bool {
				if ctx.User() == username && pass == password {
					return true
				} else {
					return false
				}
			},
		},
	}
}

// ListenAndServe listens on the TCP network address srv.Addr
// and then calls Serve to handle incoming connections
func (fs *TestSSHServer) ListenAndServe() error {
	return fs.server.ListenAndServe()
}

// Close returns any error returned from closing
// the Server's underlying Listener(s).
func (fs *TestSSHServer) Close() error {
	return fs.server.Close()
}

// SetReturnString takes in a string and set it as
// the response from the server
func (fs *TestSSHServer) SetReturnString(str string) {
	fs.server.Handler = func(s ssh.Session) {
		io.WriteString(s, str)
	}
}
