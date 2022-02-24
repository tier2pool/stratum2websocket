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

type Server struct {
	config     *Config
	upgrader   websocket.Upgrader
	httpServer *echo.Echo
}

func (s *Server) Run() error {
	logrus.Infoln("server is running")
	return s.httpServer.StartTLS(
		s.config.Listen,
		s.config.SSLCertificate,
		s.config.SSLCertificateKey,
	)
}

func (s *Server) initializeHttpServer() error {
	s.httpServer.HideBanner = true
	s.httpServer.HidePort = true
	s.httpServer.HTTPErrorHandler = func(err error, c echo.Context) {
		logrus.Errorln(err)
		if err := c.NoContent(http.StatusServiceUnavailable); err != nil {
			logrus.Errorln(err)
		}
	}
	redirectURL, err := url.Parse(s.config.Redirect)
	if err != nil {
		logrus.Fatal(err)
		return nil
	}
	s.httpServer.Use(middleware.ProxyWithConfig(middleware.ProxyConfig{
		Skipper: func(c echo.Context) bool {
			authorization := c.Request().Header.Get("Authorization")
			skip := authorization == fmt.Sprintf("Basic %s", s.config.Token)
			if !skip {
				logrus.Infof("invalid client or password error from %s\n", c.RealIP())
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
	return nil
}

func (s *Server) handle(c echo.Context) error {
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
	case "tcp":
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

func (s *Server) forward(conn jsonrpc2.Conn, ws *websocket.Conn) error {
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

func (s *Server) receive(conn jsonrpc2.Conn, ws *websocket.Conn) error {
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

func (s *Server) send(conn jsonrpc2.Conn, ws *websocket.Conn) error {
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

type Config struct {
	Debug             bool
	Token             string
	Listen            string
	SSLCertificate    string
	SSLCertificateKey string
	Redirect          string
}

func NewCommand() *cobra.Command {
	server := &Server{
		httpServer: echo.New(),
		config:     &Config{},
	}
	command := &cobra.Command{
		Use:  "server",
		Long: "tier2pool server",
		Run: func(cmd *cobra.Command, args []string) {
			if server.config.Debug {
				logrus.SetLevel(logrus.DebugLevel)
			}
			if err := server.initializeHttpServer(); err != nil {
				logrus.Fatalln(err)
			}
			if err := server.Run(); err != nil {
				logrus.Fatalln(err)
			}
		},
	}
	command.Flags().BoolVar(&server.config.Debug, "debug", false, "Enable debug mode")
	command.Flags().StringVar(&server.config.Token, "token", "", "Server access token")
	command.Flags().StringVar(&server.config.Redirect, "redirect", "", "Redirect url for invalid requests")
	command.Flags().StringVar(&server.config.Listen, "listen", "0.0.0.0:443", "Server listener address")
	command.Flags().StringVar(&server.config.SSLCertificate, "ssl-certificate", "", "SSL certificate")
	command.Flags().StringVar(&server.config.SSLCertificateKey, "ssl-certificate-key", "", "SSL certificate private key")
	if err := command.MarkFlagRequired("token"); err != nil {
		logrus.Fatal(err)
	}
	if err := command.MarkFlagRequired("ssl-certificate"); err != nil {
		logrus.Fatal(err)
	}
	if err := command.MarkFlagRequired("ssl-certificate-key"); err != nil {
		logrus.Fatal(err)
	}
	return command
}
