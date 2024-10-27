package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/joshy-joy/essay-word-counter/config"
	"github.com/joshy-joy/essay-word-counter/jobs"
)

func shutdown(cancel context.CancelFunc) {
	// Capture system interrupt signals for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("Execution terminated")
	cancel()
}

func getFlags() {
	file := flag.String("file", config.Get().DefaultFilePath, "Optional: To set file path containing the url")
	count := flag.Int("top", config.Get().ResultLength, "Optional: To set result count")
	flag.Parse()
	config.SetFilePath(*file)
	config.SetTopN(*count)
}

// Main function with graceful shutdown support
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown on interrupt signal
	go shutdown(cancel)

	err := config.InitConfig()
	if err != nil {
		log.Fatal("error initializing configurations")
	}

	// get arguments from cmd
	getFlags()

	err = jobs.StartWorkerPool(ctx)
	if err != nil {
		log.Fatal("error running job")
	}
}
