package main

import (
	"log"
	"os"

	"github.com/bjia56/objective-lol/pkg/lsp/server"
)

func main() {
	// Create and start the LSP server
	lspServer := server.NewServer()

	// Log to stderr since stdout is used for LSP communication
	log.SetOutput(os.Stderr)
	log.Println("Starting Objective-LOL Language Server...")

	// Start the server
	if err := lspServer.Start(); err != nil {
		log.Fatalf("Failed to start LSP server: %v", err)
	}
}
