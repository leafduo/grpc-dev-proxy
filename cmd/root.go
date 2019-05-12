package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/leafduo/grpc-dev-proxy/handler"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "grpc-dev-proxy",
	Short: "grpc-dev-proxy is a proxy for debugging gRPC service more easily",
	Run: func(cmd *cobra.Command, args []string) {
		mux := http.NewServeMux()
		mux.HandleFunc("/", handler.HandleRequest)

		addr := fmt.Sprintf("127.0.0.1:%d", port)
		fmt.Printf("Listening on: %s\n", addr)
		err := http.ListenAndServe(addr, mux)
		if err != nil {
			panic(err)
		}

	},
}

var port int

func init() {
	rootCmd.PersistentFlags().IntVar(&port, "port", 2333, "listening port")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
