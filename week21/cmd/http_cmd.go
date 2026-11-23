package cmd

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
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
		method = strings.ToUpper(method)

		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			return err
		}

		client := &http.Client{Timeout: 10 * time.Second}
		start := time.Now()
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		fmt.Printf("URL     : %s\n", url)
		fmt.Printf("Method  : %s\n", method)
		fmt.Printf("Status  : %s\n", resp.Status)
		fmt.Printf("Latency : %v\n", time.Since(start).Round(time.Millisecond))
		fmt.Println("Headers :")

		keys := make([]string, 0, len(resp.Header))
		for k := range resp.Header {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Printf("  %-30s %s\n", k+":", resp.Header[k][0])
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
	httpCmd.Flags().StringP("method", "m", "HEAD", "HTTP method")
}
