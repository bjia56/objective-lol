package main

import (
	"log"
	"os"

	"github.com/bjia56/objective-lol/pkg/lsp/server"
)

func main() {
	// Create and start the LSP server
	lspServer, err := server.NewServer("")
	if err != nil {
		log.Fatalf("Failed to create LSP server: %v", err)
	}

	// Log to stderr since stdout is used for LSP communication
	log.SetOutput(os.Stderr)
	log.Println("Starting Objective-LOL Language Server...")

	// Start the server
	if err := lspServer.Start(); err != nil {
		log.Fatalf("Failed to start LSP server: %v", err)
	}
}
