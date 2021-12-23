package client

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tier2pool/tier2pool/internal/jsonrpc2"
	"github.com/tier2pool/tier2pool/internal/pool"
	"io"
	"net"
	"net/http"
)

type Client interface {
	Run() error
}

func New(cmd *cobra.Command) Client {
	if cmd.Flag("debug").Value.String() == "true" {
		logrus.SetLevel(logrus.DebugLevel)
	}
	return &client{
		cmd: cmd,
	}
}

type client struct {
	cmd      *cobra.Command
	listener net.Listener
}

func (c *client) Run() error {
	defer logrus.Info("client has exited")
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		return err
	}
	defer listener.Close()
	logrus.Info("client is running")
	for {
		conn, err := listener.Accept()
		if err != nil {
			logrus.Error(err)
			continue
		}
		go func() {
			if err := c.handle(jsonrpc2.NewConn(conn)); err != nil {
				logrus.Error(err)
			}
		}()
	}
}

func (c *client) handle(conn jsonrpc2.Conn) error {
	logrus.Infof("new connection from %s\n", conn.RemoteAddr())
	defer func() {
		_ = conn.Close()
		logrus.Infof("client %s disconnect\n", conn.RemoteAddr())
	}()
	ws, _, err := websocket.DefaultDialer.Dial(c.cmd.Flag("server").Value.String(), http.Header{
		"Authorization": []string{fmt.Sprintf("Basic %s", c.cmd.Flag("token").Value.String())},
		"X-Pool":        []string{c.cmd.Flag("pool").Value.String()},
	})
	if err != nil {
		if errors.Is(err, websocket.ErrBadHandshake) || errors.Is(err, io.ErrUnexpectedEOF) {
			logrus.Warn("invalid server or wrong password")
			return nil
		}
		return err
	}
	defer func() {
		_ = ws.Close()
	}()
	errCh := make(chan error)
	go func() {
		errCh <- c.receive(conn, ws)
	}()
	go func() {
		errCh <- c.send(conn, ws)
	}()
	select {
	case err := <-errCh:
		return err
	}
}

func (c *client) receive(conn jsonrpc2.Conn, ws *websocket.Conn) error {
	buffer := pool.Get()
	defer pool.Put(buffer)
	for {
		n, err := conn.Read(buffer)
		if n > 0 {
			logrus.Debugf("<-- %s", string(buffer[:n]))
			if err := ws.WriteMessage(websocket.TextMessage, buffer[:n]); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					return nil
				}
				return err
			}
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
	}
}

func (c *client) send(conn jsonrpc2.Conn, ws *websocket.Conn) error {
	for {
		messageType, data, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				return nil
			}
			return err
		}
		switch messageType {
		case websocket.TextMessage:
			if len(data) > 0 {
				logrus.Debugf("--> %s", string(data))
				if _, err := conn.Write(data); err != nil {
					return err
				}
			}
		default:
			return fmt.Errorf("unsupported message type %d\n", messageType)
		}
	}
}
