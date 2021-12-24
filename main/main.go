package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tier2pool/tier2pool/main/internal/client"
	"github.com/tier2pool/tier2pool/main/internal/server"
)

var rootCmd = &cobra.Command{
	Use:  "tier2pool",
	Long: "A mining pool proxy tool, support BTC, ETH, ETC, XMR mining pool, etc. https://github.com/tier2pool/tier2pool",
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "tier2pool server",
	Run: func(cmd *cobra.Command, args []string) {
		srv := server.New(cmd)
		if err := srv.Run(); err != nil {
			logrus.Fatalln(err)
		}
	},
}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "tier2pool client",
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
	serverCmd.Flags().Bool("debug", false, "enable debug mode")
	serverCmd.Flags().String("listen", "0.0.0.0:443", "server listener address")
	serverCmd.Flags().String("redirect", "", "redirect url for invalid requests")
	serverCmd.Flags().String("ssl-certificate", "", "ssl certificate")
	serverCmd.Flags().String("ssl-certificate-key", "", "ssl certificate private key")
	if err := serverCmd.MarkFlagRequired("ssl-certificate"); err != nil {
		logrus.Fatal(err)
	}
	if err := serverCmd.MarkFlagRequired("ssl-certificate-key"); err != nil {
		logrus.Fatal(err)
	}
	if err := serverCmd.MarkFlagRequired("token"); err != nil {
		logrus.Fatal(err)
	}
	clientCmd.Flags().Bool("debug", false, "enable debug mode")
	clientCmd.Flags().String("server", "", "tier2pool server address")
	clientCmd.Flags().String("pool", "", "mining pool address")
	clientCmd.Flags().String("token", "", "server access token")
	clientCmd.Flags().String("listen", "127.0.0.1:1234", "client listener address")
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
