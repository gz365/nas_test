package client

import (
	"context"
	"errors"
	"fmt"
	"net"
	"regexp"
	"smb_gbox/conf"

	"github.com/hirochachacha/go-smb2"
)

type SMBClient struct {
	context context.Context
	conn    net.Conn
	dialer  *smb2.Dialer
	session *smb2.Session
	share   *smb2.Share
}

func NewClient(cf conf.Conf) (*SMBClient, error) {
	conn, err := net.Dial("tcp", cf.Host) // smb协议端口号445，ip地址看手机app里的设置
	if err != nil {

		return nil, err
	}

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     cf.Username,
			Password: cf.Password,
		},
	}

	s, err := d.Dial(conn)
	if err != nil {

		return nil, err
	}

	names, err := s.ListSharenames()
	if err != nil {
		return nil, err
	}

	shareName := ""
	reg := regexp.MustCompile("^" + cf.Username)
	for _, name := range names {
		if reg.MatchString(name) {
			shareName = name
			break
		}
	}

	if shareName == "" {
		return nil, errors.New("没获取到共享目录")
	}

	share, err := s.Mount(names[0])
	if err != nil {
		return nil, err
	}

	return &SMBClient{
			conn:    conn,
			dialer:  d,
			session: s,
			share:   share,
		},
		nil
}

func (c *SMBClient) Disconnect() {
	c.conn.Close()
	c.session.Logoff()
	c.share.Umount()
	fmt.Println("smbclient quit")
}

func (c *SMBClient) Upload() error {
	f, err := c.share.Create("0.text")
	if err != nil {
		return err
	}
	defer f.Close()
	b := []byte("hello world")
	_, err = f.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (c *SMBClient) ReadDir() error {
	fsArr, err := c.share.ReadDir(".")
	if err != nil {
		return err
	}
	if len(fsArr) < 1 {
		fmt.Println("空的")
	}
	for _, file := range fsArr {
		fmt.Println(file.Name())
	}
	return nil
}
