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

func signalHandler(pid *int, emitter logging.Emitter, signalChan chan os.Signal) {
	var signalToSendChild syscall.Signal = syscall.SIGHUP
	for {
		signalReceived := <-signalChan
		log.Println("Signal received:", signalReceived)
		switch signalReceived {
		case syscall.SIGHUP:
			signalToSendChild = syscall.SIGHUP
		case syscall.SIGINT:
			signalToSendChild = syscall.SIGINT
		case syscall.SIGTERM:
			signalToSendChild = syscall.SIGTERM
		case syscall.SIGKILL:
			signalToSendChild = syscall.SIGKILL
		}
		log.Println("Sending signal to", *pid)
		syscall.Kill(*pid, signalToSendChild)
		log.Println("Trigering emitter shutdown")
		if err := emitter.Cleanup(); err != nil {
			log.Println("Error cleaning up emitter:", err)
		}
		os.Exit(0)
	}
}

func main() {
	log.Println("Initializing haberdasher.")

	emitterName, exists := os.LookupEnv("HABERDASHER_EMITTER")
	if !exists {
		emitterName = "stdout"
	}
	log.Println("Configured emitter:", emitterName)
	emitter := logging.Emitters[emitterName]

	go reaper.Reap()
	subcmdPid := -1
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go signalHandler(&subcmdPid, emitter, signalChan)

	subcmdBin := os.Args[1]
	subcmdArgs := os.Args[2:len(os.Args)]


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
	subcmdPid = subcmd.Process.Pid

	for scanner.Scan() {
		go func() {
			logging.Emit(emitter, scanner.Text())
		}()
	}
}