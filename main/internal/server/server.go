package server

import (
	"crypto/tls"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tier2pool/tier2pool/internal/jsonrpc2"
	"github.com/tier2pool/tier2pool/internal/pool"
	"net"
	"net/http"
	"net/url"
)

type Server interface {
	Run() error
}

func New(cmd *cobra.Command) Server {
	if cmd.Flag("debug").Value.String() == "true" {
		logrus.SetLevel(logrus.DebugLevel)
	}
	s := &server{
		cmd:        cmd,
		httpServer: echo.New(),
	}
	s.httpServer.HideBanner = true
	s.httpServer.HidePort = true
	s.httpServer.HTTPErrorHandler = func(err error, c echo.Context) {
		logrus.Error(err)
		if err := c.NoContent(http.StatusServiceUnavailable); err != nil {
			logrus.Error(err)
		}
	}
	redirectURL, err := url.Parse(s.cmd.Flag("redirect").Value.String())
	if err != nil {
		logrus.Fatal(err)
		return nil
	}
	s.httpServer.Use(middleware.ProxyWithConfig(middleware.ProxyConfig{
		Skipper: func(c echo.Context) bool {
			authorization := c.Request().Header.Get("Authorization")
			skip := authorization == fmt.Sprintf("Basic %s", s.cmd.Flag("token").Value.String())
			if !skip {
				logrus.Infof("invalid client or password error from %s", c.RealIP())
			}
			return skip
		},
		Balancer: middleware.NewRandomBalancer([]*middleware.ProxyTarget{
			{
				URL: redirectURL,
			},
		}),
	}))
	s.httpServer.GET("/", s.handle)
	return s
}

type server struct {
	cmd        *cobra.Command
	upgrader   websocket.Upgrader
	httpServer *echo.Echo
}

func (s *server) Run() error {
	logrus.Info("server is running")
	return s.httpServer.StartTLS(
		s.cmd.Flag("listen").Value.String(),
		s.cmd.Flag("ssl-certificate").Value.String(),
		s.cmd.Flag("ssl-certificate-key").Value.String(),
	)
}

func (s *server) handle(c echo.Context) error {
	ws, err := s.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = ws.Close()
	}()
	logrus.Infof("new connection from %s\n", c.RealIP())
	defer func() {
		logrus.Infof("client %s disable\n", c.RealIP())
	}()
	var poolConn net.Conn
	u, err := url.Parse(c.Request().Header.Get("X-Pool"))
	if err != nil {
		return err
	}
	switch u.Scheme {
	case "", "tcp":
		poolConn, err = net.Dial("tcp", u.Host)
	case "tls", "ssl":
		poolConn, err = tls.Dial("tcp", u.Host, &tls.Config{})
	default:
		return fmt.Errorf("unsupported mining pool protocol: %s", u.Scheme)
	}
	if err != nil {
		return err
	}
	conn := jsonrpc2.NewConn(poolConn)
	defer func() {
		_ = conn.Close()
	}()
	return s.forward(conn, ws)
}

func (s *server) forward(conn jsonrpc2.Conn, ws *websocket.Conn) error {
	errCh := make(chan error)
	go func() {
		errCh <- s.receive(conn, ws)
	}()
	go func() {
		errCh <- s.send(conn, ws)
	}()
	select {
	case err := <-errCh:
		return err
	}
}

func (s *server) receive(conn jsonrpc2.Conn, ws *websocket.Conn) error {
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
			logrus.Debugf("<-- %s", string(data))
			if _, err := conn.Write(data); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported message type %d\n", messageType)
		}
	}
}

func (s *server) send(conn jsonrpc2.Conn, ws *websocket.Conn) error {
	buffer := pool.Get()
	defer pool.Put(buffer)
	for {
		n, err := conn.Read(buffer)
		if n > 0 {
			logrus.Debugf("--> %s", string(buffer[:n]))
			if err := ws.WriteMessage(websocket.TextMessage, buffer[:n]); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					return nil
				}
				return err
			}
		}
		if err != nil {
			return err
		}
	}
}
