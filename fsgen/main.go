package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ttacon/chalk"
	"github.com/ttacon/fs"
)

var (
	interactive = flag.Bool("i", false, "run fsgen in interactive mode")
)

var (
	title     = chalk.Magenta.Color("[fsgen]:")
	errPrompt = chalk.Bold.NewStyle().WithForeground(chalk.Red).Style("ERROR")

	defOS = fs.DefaultOS()
)

func main() {
	flag.Parse()

	if *interactive {
		fmt.Println(title, "running in interactive mode")
		repl()
		return
	}
	fmt.Println(title, "only interactive mode is currently supported")
}

func repl() {
	in := bufio.NewReader(os.Stdin)
	cwd, err := defOS.Getwd()
	if err != nil {
		errMessage("failed to get cwd, err: ", err)
		return
	}

	f, err := defOS.Open(cwd)
	if err != nil {
		errMessage("failed to open current directory, err: ", err)
		return
	}

	for {
		// TODO(ttacon): make prompt reflect current directory?
		fmt.Print("> ")
		nextLine, err := in.ReadString('\n')
		if err != nil {
			errMessage("failed to read next line, err: ", err)
			return
		}

		nextLine = strings.TrimSpace(nextLine)
		if nextLine == "ls" {
			entities, err := f.Readdirnames(-1)
			if err != nil {
				errMessage("failed to read directory, err: ", err)
				return
			}

			fmt.Println(strings.Join(entities, "\n"))
		} else if nextLine == "help" {
			listCommands()
		} else if strings.HasPrefix(nextLine, "cd") {
			// TODO(ttacon): support .. and ~, or are they transparently
			// supported through filepath.Join?

			dir := strings.TrimSpace(strings.TrimPrefix(nextLine, "cd"))
			newDir, err := defOS.Open(dir)
			if err != nil {
				dir = filepath.Join(cwd, dir)
				newDir, err = defOS.Open(dir)
				if err != nil {
					errMessage("failed to change directory, err: ", err)
					continue
				} else {
					cwd = dir
				}
			} else {
				cwd = dir
			}

			oldDir := f
			f = newDir
			oldDir.Close()
		} else if strings.HasPrefix(nextLine, "sync") {

		} else if nextLine == "exit" || nextLine == "quit" {
			return
		}
	}
}

func listCommands() {
	fmt.Println(`
Available Commands:

help   = display this list
ls     = list current directory
cd x   = change directory to directory x (x must exist)
sync y = add y to the current fs state to be saved
exit   = quit
quit   = quit
`)
}

func errMessagef(fmtString string, args ...interface{}) {
	fmt.Println(title, errPrompt, fmt.Sprintf(fmtString, args...))
}

func errMessage(pieces ...interface{}) {
	args := append([]interface{}{title, errPrompt}, pieces...)
	fmt.Println(args...)
}
