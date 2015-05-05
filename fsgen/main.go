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
	currentFile = flag.String("f", "", "fsgen file to use")
)

// TODO(ttacon): parse the generate go files to get the info we need

var (
	title     = chalk.Magenta.Color("[fsgen]:")
	errPrompt = chalk.Bold.NewStyle().WithForeground(chalk.Red).Style("ERROR")
	dirColor  = chalk.Green.Color("%s")
	dirPrompt = chalk.Magenta.Color("[fsgen $ " + dirColor + "]")

	defOS   = fs.DefaultOS()
	origDir fs.File
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
	cwdEnd := filepath.Base(cwd)

	f, err := defOS.Open(cwd)
	if err != nil {
		errMessage("failed to open current directory, err: ", err)
		return
	}
	origDir = f
	defer origDir.Close()

	for {
		fmt.Print(fmt.Sprintf(dirPrompt, cwdEnd) + " > ")
		nextLine, err := in.ReadString('\n')
		if err != nil {
			errMessage("failed to read next line, err: ", err)
			return
		}

		nextLine = strings.TrimSpace(nextLine)
		if nextLine == "ls" {
			entities, err := f.Readdir(-1)
			if err != nil {
				errMessage("failed to read directory, err: ", err)
				return
			}

			var entries = make([]string, len(entities))
			for i, entity := range entities {
				entry := entity.Name()
				if entity.IsDir() {
					entry = chalk.Yellow.Color(entry + "/")
				}
				entries[i] = entry
			}

			fmt.Println(strings.Join(entries, "\n"))
		} else if nextLine == "help" {
			listCommands()
		} else if strings.HasPrefix(nextLine, "cd") {
			// NOTE(ttacon): ~ is not currently supported
			solelyDir := strings.TrimSpace(strings.TrimPrefix(nextLine, "cd"))
			dir := filepath.Join(cwd, solelyDir)
			newDir, err := defOS.Open(dir)
			if err != nil {
				dir = solelyDir
				newDir, err = defOS.Open(dir)
				if err != nil {
					errMessage("failed to change directory, err: ", err)
					continue
				}
			}

			cwd = newDir.Name()
			cwdEnd = filepath.Base(cwd)

			oldDir := f
			f = newDir
			if oldDir != origDir {
				oldDir.Close()
			}
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
