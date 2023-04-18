package main

import (
	"io"
	"os"
	"time"
	"fmt"
	"os/exec"
	"log"
)

func main() {
	t := time.Now()
	fmt.Println(t.Format("2006_01_02_15_04_05"))

	cmd := exec.Command("mycli", "-h", "127.0.0.1", "-u", "root", "-P", "4000", "--csv", "--execute", "show databases")
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	slurp, _ := io.ReadAll(stderr)
	fmt.Printf("%s\n", slurp)
	slurp, _ = io.ReadAll(stdout)
	fmt.Printf("%s\n", slurp)

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}
