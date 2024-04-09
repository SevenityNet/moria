package source

import (
	"errors"
	"io"
	"moria/config"
	"time"

	"github.com/melbahja/goph"
)

type SSHSource struct {
	path       string
	host       string
	port       int
	user       string
	auth       goph.Auth
	closeOnEnd bool
	client     *goph.Client
}

func NewSSHSource() Source {
	auth, err := loadAuth()
	if err != nil {
		panic(err)
	}

	src := &SSHSource{
		path:       config.GetSourceRemoteSSHPath(),
		host:       config.GetSourceRemoteHost(),
		port:       config.GetSourceRemotePort(),
		user:       config.GetSourceRemoteSSHUser(),
		closeOnEnd: config.IsSourceRemoteSSHCloseOnEnd(),
		auth:       auth,
	}

	if !src.closeOnEnd {
		client, err := src.getSSHClient(30)
		if err != nil {
			panic(err)
		}

		src.client = client
	}

	return src
}

func (s *SSHSource) Upload(file []byte, filename string) error {
	client := s.client

	if s.closeOnEnd {
		c, err := s.getSSHClient(5)
		if err != nil {
			return err
		}

		defer c.Close()

		client = c
	}

	ftp, err := client.NewSftp()
	if err != nil {
		return err
	}
	defer ftp.Close()

	remote, err := ftp.Create(s.path + filename)
	if err != nil {
		return err
	}
	defer remote.Close()

	_, err = remote.Write(file)
	if err != nil {
		return err
	}

	return nil
}

func (s *SSHSource) Delete(filename string) error {
	client := s.client

	if s.closeOnEnd {
		c, err := s.getSSHClient(5)
		if err != nil {
			return err
		}

		defer c.Close()

		client = c
	}

	ftp, err := client.NewSftp()
	if err != nil {
		return err
	}

	err = ftp.Remove(s.path + filename)
	if err != nil {
		return err
	}

	return nil
}

func (s *SSHSource) Get(filename string) ([]byte, error) {
	client := s.client

	if s.closeOnEnd {
		c, err := s.getSSHClient(5)
		if err != nil {
			return nil, err
		}

		defer c.Close()

		client = c
	}

	ftp, err := client.NewSftp()
	if err != nil {
		return nil, err
	}

	remote, err := ftp.Open(s.path + filename)
	if err != nil {
		return nil, err
	}

	defer remote.Close()

	local, err := io.ReadAll(remote)
	if err != nil {
		return nil, err
	}

	return local, nil
}

func (s *SSHSource) getSSHClient(seconds int) (*goph.Client, error) {
	client, err := goph.NewConn(&goph.Config{
		Auth:    s.auth,
		User:    s.user,
		Addr:    s.host,
		Port:    uint(s.port),
		Timeout: time.Duration(seconds) * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func loadAuth() (goph.Auth, error) {
	if config.GetSourceRemoteSSHAuth() == "password" {
		return goph.Password(config.GetSourceRemoteSSHPass()), nil
	} else if config.GetSourceRemoteSSHAuth() == "key" {
		auth, err := goph.Key(config.GetSourceRemoteSSHKey(), config.GetSourceRemoteSSHKeyPass())
		if err != nil {
			return nil, err
		}

		return auth, nil
	}

	return nil, errors.New("invalid SSH auth type")
}
