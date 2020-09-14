package main

import (
	"log"
	"os"
	"os/exec"
	"bufio"
	reaper "github.com/ramr/go-reaper"
	"github.com/RedHatInsights/haberdasher/logging"
	_ "github.com/RedHatInsights/haberdasher/emitters"
)

func main() {
	go reaper.Reap()
	log.Println("Initializing haberdasher.")

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

	for scanner.Scan() {
		go func() {
			logging.Emit(emitter, scanner.Text())
		}()
	}
}