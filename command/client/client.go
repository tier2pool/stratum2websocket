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

type Client struct {
	config   *Config
	listener net.Listener
}

func (c *Client) Run() error {
	defer logrus.Infoln("client has exited")
	listener, err := net.Listen("tcp", c.config.Listen)
	if err != nil {
		return err
	}
	defer func() {
		_ = listener.Close()
	}()
	logrus.Infoln("client is running")
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

func (c *Client) handle(conn jsonrpc2.Conn) error {
	logrus.Infof("new connection from %s\n", conn.RemoteAddr())
	defer func() {
		_ = conn.Close()
		logrus.Infof("client %s disconnect\n", conn.RemoteAddr())
	}()
	ws, _, err := websocket.DefaultDialer.Dial(c.config.Server, http.Header{
		"Authorization": []string{fmt.Sprintf("Basic %s", c.config.Token)},
		"X-Pool":        []string{c.config.Pool},
	})
	if err != nil {
		if errors.Is(err, websocket.ErrBadHandshake) || errors.Is(err, io.ErrUnexpectedEOF) {
			logrus.Warnln("invalid server or wrong password")
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

func (c *Client) receive(conn jsonrpc2.Conn, ws *websocket.Conn) error {
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

func (c *Client) send(conn jsonrpc2.Conn, ws *websocket.Conn) error {
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

type Config struct {
	Debug  bool
	Server string
	Pool   string
	Token  string
	Listen string
}

func NewCommand() *cobra.Command {
	client := &Client{
		config: &Config{},
	}
	command := &cobra.Command{
		Use:  "client",
		Long: "tier2pool client",
		Run: func(cmd *cobra.Command, args []string) {
			if client.config.Debug {
				logrus.SetLevel(logrus.DebugLevel)
			}
			if err := client.Run(); err != nil {
				logrus.Fatalln(err)
			}
		},
	}
	command.Flags().BoolVar(&client.config.Debug, "debug", false, "enable debug mode")
	command.Flags().StringVar(&client.config.Server, "server", "", "tier2pool server address")
	command.Flags().StringVar(&client.config.Pool, "pool", "", "mining pool address")
	command.Flags().StringVar(&client.config.Token, "token", "", "server access token")
	command.Flags().StringVar(&client.config.Listen, "listen", "", "client listener address")
	if err := command.MarkFlagRequired("pool"); err != nil {
		logrus.Fatalln(err)
	}
	if err := command.MarkFlagRequired("server"); err != nil {
		logrus.Fatalln(err)
	}
	if err := command.MarkFlagRequired("token"); err != nil {
		logrus.Fatalln(err)
	}
	return command
}
