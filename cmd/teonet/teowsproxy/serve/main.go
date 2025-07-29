// Copyright 2023-2025 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The teonet-ws-proxy application is a WebSocket proxy server for Teonet.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/teonet-go/teoproxy/teonet/server"
)

const (
	appShort    = "teowsproxy"
	appName     = "Teonet WebSocket Proxy"
	appVersion  = "0.1.0"
	defaultPort = "8080"
)

// main is the entry point of the program.
//
// It initializes and runs the Teonet WebSocket proxy server. It's designed to
// run in containerized environments like Google Cloud Run, listening on the
// port specified by the PORT environment variable.
//
// publish one container running all time in Google Cloud Run:
//
//	gcloud run deploy teo-ws-proxy --source . --region=europe-north1 --no-cpu-throttling --min-instances=1 --max-instances=1 --timeout=3600 --session-affinity --concurrency=1000
//
// publish any containers running when thry need in Google Cloud Run:
//
//	gcloud run deploy teo-ws-proxy --source . --region=europe-north1 --timeout=3600 --session-affinity --concurrency=1000
func main() {

	// --- Configuration ---
	var monitor string
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	addr := ":" + port
	flag.StringVar(&monitor, "monitor", "", "teonet monitor address")
	flag.Parse()

	// --- Teonet Proxy Server Initialization ---
	teoServer, err := server.New(appShort, &server.TeonetMonitor{
		Addr:       monitor,
		AppName:    appName,
		AppShort:   appShort,
		AppVersion: appVersion,
	})
	if err != nil {
		log.Fatalf("Create teonet proxy server error: %v", err)
	}

	// --- HTTP Server Setup ---
	mux := http.NewServeMux()

	// Health check handler for Cloud Run
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Teonet WebSocket Proxy is running")
	})

	// The main WebSocket handler
	mux.HandleFunc("/ws", teoServer.HandleWebSocket)

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// --- Graceful Shutdown ---
	go func() {
		log.Printf("Starting Teonet WebSocket Proxy on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server gracefully stopped")
}
