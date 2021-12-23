package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tier2pool/tier2pool/main/internal/client"
	"github.com/tier2pool/tier2pool/main/internal/server"
)

var rootCmd = &cobra.Command{
	Use: "tier2pool",
}

var serverCmd = &cobra.Command{
	Use: "server",
	Run: func(cmd *cobra.Command, args []string) {
		srv := server.New(cmd)
		if err := srv.Run(); err != nil {
			logrus.Fatalln(err)
		}
	},
}

var clientCmd = &cobra.Command{
	Use: "client",
	Run: func(cmd *cobra.Command, args []string) {
		cli := client.New(cmd)
		if err := cli.Run(); err != nil {
			logrus.Fatalln(err)
		}
	},
}

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	serverCmd.Flags().String("token", "", "server access token")
	serverCmd.Flags().String("listen", "0.0.0.0:443", "")
	serverCmd.Flags().String("redirect", "", "")
	serverCmd.Flags().String("ssl-certificate", "", "")
	serverCmd.Flags().String("ssl-certificate-key", "", "")
	if err := serverCmd.MarkFlagRequired("ssl-certificate"); err != nil {
		logrus.Fatal(err)
	}
	if err := serverCmd.MarkFlagRequired("ssl-certificate-key"); err != nil {
		logrus.Fatal(err)
	}
	if err := serverCmd.MarkFlagRequired("token"); err != nil {
		logrus.Fatal(err)
	}
	clientCmd.Flags().String("server", "", "")
	clientCmd.Flags().String("pool", "", "")
	clientCmd.Flags().String("token", "", "")
	clientCmd.Flags().String("listen", "127.0.0.1:1234", "")
	if err := clientCmd.MarkFlagRequired("pool"); err != nil {
		logrus.Fatal(err)
	}
	if err := clientCmd.MarkFlagRequired("server"); err != nil {
		logrus.Fatal(err)
	}
	if err := clientCmd.MarkFlagRequired("token"); err != nil {
		logrus.Fatal(err)
	}
	rootCmd.AddCommand(serverCmd, clientCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err)
	}
}
