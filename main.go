package main

import (
	"log"
	"os"
	"os/exec"
	"bufio"
	"os/signal"
	"syscall"
	reaper "github.com/ramr/go-reaper"
	"github.com/RedHatInsights/haberdasher/logging"
	_ "github.com/RedHatInsights/haberdasher/emitters"
)

func main() {
	log.Println("Initializing haberdasher.")
	go reaper.Reap()
	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, syscall.SIGINT, syscall.SIGTERM)


	subcmdBin := os.Args[1]
	subcmdArgs := os.Args[2:len(os.Args)]

	emitterName, exists := os.LookupEnv("HABERDASHER_EMITTER")
	if !exists {
		emitterName = "stdout"
	}
	log.Println("Configured emitter:", emitterName)
	emitter := logging.Emitters[emitterName]

	subcmd := exec.Command(subcmdBin, subcmdArgs...)
	// pass through stdout
	subcmd.Stdout = os.Stdout
	subcmdErr, err := subcmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(subcmdErr)

	if err := subcmd.Start(); err != nil {
		log.Fatal(err)
	}

	go func() {
		for scanner.Scan() {
			go func() {
				logging.Emit(emitter, scanner.Text())
			}()
		}
	}()

	<-killSignal
	log.Println("Haberdasher shutting down.")
	err = emitter.Cleanup()
	if err != nil {
		log.Fatal("Error during shutdown:", err)
	} else {
		log.Println("Haberdasher shut down cleanly.")
	}
}