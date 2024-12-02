package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/zwpaper/llmock/pkg/server"
)

func main() {
	rootCmd := createRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "llmock",
		Short: "llmock is a mock server for embeddings",
	}

	rootCmd.AddCommand(createServeCommand())
	return rootCmd
}

func createServeCommand() *cobra.Command {
	var address string
	var port int
	var embedding bool

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the mock server",
		Run: func(cmd *cobra.Command, args []string) {
			addr := "0.0.0.0"

			srv := server.New(addr, port)

			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()

			if err := srv.Run(ctx); err != nil && err != http.ErrServerClosed {
				fmt.Printf("Server error: %v\n", err)
			}
		},
	}

	serveCmd.Flags().IntVarP(&port, "port", "p", 30888, "Port for the server")
	serveCmd.Flags().StringVarP(&address, "address", "a", "localhost", "Address for the server")
	serveCmd.Flags().BoolVarP(&embedding, "embedding", "e", true, "Start the embedding server")
	return serveCmd
}
