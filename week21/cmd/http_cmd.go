package cmd

import (
	"fmt"
	"net/http"
	"time"
	"github.com/spf13/cobra"
)

var httpCmd = &cobra.Command{
	Use:   "http <url>",
	Short: "Inspect HTTP response headers and status",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]
		method, _ := cmd.Flags().GetString("method")
		client := &http.Client{Timeout: 10 * time.Second}

		start := time.Now()
		resp, err := client.Head(url)
		if err != nil { return err }
		defer resp.Body.Close()

		fmt.Printf("URL     : %s\n", url)
		fmt.Printf("Method  : %s\n", method)
		fmt.Printf("Status  : %s\n", resp.Status)
		fmt.Printf("Latency : %v\n", time.Since(start).Round(time.Millisecond))
		fmt.Println("Headers :")
		for k, v := range resp.Header {
			fmt.Printf("  %-30s %s\n", k+":", v[0])
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
	httpCmd.Flags().StringP("method", "m", "HEAD", "HTTP method")
}