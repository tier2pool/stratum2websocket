package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tier2pool/tier2pool/command/client"
	"github.com/tier2pool/tier2pool/command/server"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
}

func main() {
	command := &cobra.Command{
		Use:  "tier2pool",
		Long: "A mining pool proxy tool, support BTC, ETH, ETC, XMR mining pool, etc. https://github.com/tier2pool/tier2pool",
	}
	command.AddCommand(client.NewCommand())
	command.AddCommand(server.NewCommand())
	if err := command.Execute(); err != nil {
		logrus.Fatalln(err)
	}
}
