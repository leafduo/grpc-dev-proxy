package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/leafduo/grpc-dev-proxy/handler"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hugo",
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
				  love by spf13 and friends in Go.
				  Complete documentation is available at http://hugo.spf13.com`,
	Run: func(cmd *cobra.Command, args []string) {
		mux := http.NewServeMux()
		mux.HandleFunc("/", handler.HandleRequest)

		http.ListenAndServe("127.0.0.1:2333", mux)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
