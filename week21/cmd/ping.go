package cmd

import (
	"fmt"
	"net"
	"time"
	"github.com/spf13/cobra"
)

var pingCmd = &cobra.Command{
	Use:   "ping <host>",
	Short: "Check TCP connectivity to a host",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		host := args[0]
		port, _ := cmd.Flags().GetString("port")
		timeout, _ := cmd.Flags().GetDuration("timeout")
		addr := net.JoinHostPort(host, port)

		start := time.Now()
		conn, err := net.DialTimeout("tcp", addr, timeout)
		elapsed := time.Since(start)

		if err != nil {
			return fmt.Errorf("❌ %s unreachable: %w", addr, err)
		}
		conn.Close()
		fmt.Printf("✅ %s reachable in %v\n", addr, elapsed.Round(time.Millisecond))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
	pingCmd.Flags().StringP("port",    "p", "80",           "TCP port to check")
	pingCmd.Flags().Duration("timeout", 5*time.Second,     "connection timeout")
}