package cmd

import (
	"fmt"
	"net"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		if viper.GetBool("verbose") {
			fmt.Printf("[verbose] dialing %s with timeout %v\n", addr, timeout)
		}

		start := time.Now()
		conn, err := net.DialTimeout("tcp", addr, timeout)
		elapsed := time.Since(start)
		if err != nil {
			return fmt.Errorf("[FAIL] %s unreachable: %w", addr, err)
		}
		conn.Close()
		fmt.Printf("[OK] %s reachable in %v\n", addr, elapsed.Round(time.Millisecond))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
	pingCmd.Flags().StringP("port", "p", "80", "TCP port to check")
	pingCmd.Flags().Duration("timeout", 5*time.Second, "connection timeout")
}
