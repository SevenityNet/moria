package source

import (
	"bytes"
	"fmt"
	"moria/config"
	"time"

	"github.com/jlaffaye/ftp"
)

type FTPSource struct {
	path       string
	host       string
	port       int
	user       string
	pass       string
	closeOnEnd bool
	client     *ftp.ServerConn
}

func NewFTPSource() Source {
	src := &FTPSource{
		path:       config.GetSourceRemoteFTPPath(),
		host:       config.GetSourceRemoteHost(),
		port:       config.GetSourceRemotePort(),
		user:       config.GetSourceRemoteFTPUser(),
		pass:       config.GetSourceRemoteFTPPass(),
		closeOnEnd: config.IsSourceRemoteFTPCloseOnEnd(),
	}

	if !src.closeOnEnd {
		client, err := getFTPClient(30)
		if err != nil {
			panic(err)
		}

		src.client = client
	}

	return src
}

func (s *FTPSource) Upload(file []byte, filename string) error {
	client := s.client

	if s.closeOnEnd {
		c, err := getFTPClient(5)
		if err != nil {
			return err
		}

		defer c.Quit()

		client = c
	}

	if s.path != "" {
		err := client.ChangeDir(s.path)
		if err != nil {
			return err
		}
	}

	err := client.Stor(s.path+filename, bytes.NewReader(file))
	if err != nil {
		return err
	}

	return nil
}

func (s *FTPSource) Delete(filename string) error {
	client := s.client

	if s.closeOnEnd {
		c, err := getFTPClient(5)
		if err != nil {
			return err
		}

		defer c.Quit()

		client = c
	}

	if s.path != "" {
		err := client.ChangeDir(s.path)
		if err != nil {
			return err
		}
	}

	err := client.Delete(s.path + filename)
	if err != nil {
		return err
	}

	return nil
}

func (s *FTPSource) Get(filename string) ([]byte, error) {
	client := s.client

	if s.closeOnEnd {
		c, err := getFTPClient(5)
		if err != nil {
			return nil, err
		}

		defer c.Quit()

		client = c
	}

	if s.path != "" {
		err := client.ChangeDir(s.path)
		if err != nil {
			return nil, err
		}
	}

	r, err := client.Retr(s.path + filename)
	if err != nil {
		return nil, err
	}

	defer r.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(r)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func getFTPClient(seconds int) (*ftp.ServerConn, error) {
	client, err := ftp.Dial(fmt.Sprintf("%s:%d", config.GetSourceRemoteHost(), config.GetSourceRemotePort()), ftp.DialWithTimeout(time.Duration(seconds)*time.Second))
	if err != nil {
		return nil, err
	}

	err = client.Login(config.GetSourceRemoteFTPUser(), config.GetSourceRemoteFTPPass())
	if err != nil {
		return nil, err
	}

	return client, nil
}
