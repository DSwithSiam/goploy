package ssh
package ssh

import (
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"fmt"
)

type SSHClient struct {
	client *ssh.Client
}

func NewSSHClient(user, addr, password string) (*SSHClient, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	c, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}
	return &SSHClient{client: c}, nil
}

func (s *SSHClient) Run(cmd string) error {
	sess, err := s.client.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()
	sess.Stdout = os.Stdout
	sess.Stderr = os.Stderr
	return sess.Run(cmd)
}

func (s *SSHClient) Upload(local io.Reader, remotePath string) error {
	sess, err := s.client.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()
	go func() {
		w, _ := sess.StdinPipe()
		fmt.Fprintln(w, "C0644", 0, remotePath)
		io.Copy(w, local)
		w.Close()
	}()
	return sess.Run("scp -t " + remotePath)
}

func (s *SSHClient) Close() error {
	return s.client.Close()
}
