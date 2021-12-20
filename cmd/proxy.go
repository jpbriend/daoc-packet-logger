package cmd

import (
	proxy "fr.jpbriend/daoc-packet-logger/internal"
	"github.com/spf13/cobra"
)

var ListenPort int
var RemoteHost string
var RemotePort int

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "this command will start dumping network packets intercepted between DAoC client and server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return proxy.Start(ListenPort, RemoteHost, RemotePort)
	},
}

func init() {
	proxyCmd.Flags().StringVar(&RemoteHost, "remoteHost", "localhost", "IP or domain name of the remote DAoC server")
	proxyCmd.MarkFlagRequired("remoteHost")
	proxyCmd.Flags().IntVar(&RemotePort, "remotePort", 10500, "port (as an int) of the remote DAoC server")
	proxyCmd.MarkFlagRequired("remotePort")
	proxyCmd.Flags().IntVar(&ListenPort, "listenPort", 7777, "port the proxy is listening for the DAoc client connections")
	rootCmd.AddCommand(proxyCmd)
}
