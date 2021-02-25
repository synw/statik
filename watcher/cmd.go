package watcher

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/synw/statik/conf"
)

func runBuildCmd(path string, verbose bool) error {
	p := strings.ReplaceAll(path, conf.WorkingDir+"/", "")
	println("Change detected in", p)
	cmds := []string{}
	for k, v := range conf.WatchBuilders {
		if strings.HasPrefix(p, k) {
			for _, c := range v {
				cmds = append(cmds, c)
			}
		}
	}
	for _, c := range cmds {
		l := strings.Split(c, " ")
		cmdName := l[0]
		cmdArgs := l[1:]
		if verbose {
			fmt.Println("Executing", cmdName, cmdArgs)
		}
		cmd := exec.Command(cmdName, cmdArgs...)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}
		// start the command after having set up the pipe
		if err := cmd.Start(); err != nil {
			return err
		}
		// read command's stdout line by line
		in := bufio.NewScanner(stdout)
		for in.Scan() {
			log.Printf(in.Text()) // write each line to your log, or anything you need
		}
		if err := in.Err(); err != nil {
			log.Printf("error: %s", err)
		}
		cmd.Wait()
	}
	return nil
}
