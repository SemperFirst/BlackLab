package main

import (
	"os"
	"runtime"

	"SafeDp/scanner/password_crack/cmd"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "password_crack"
	app.Author = "SemperFi"
	app.Email = ""
	app.Version = "2025/01/01"
	app.Usage = "Weak password scanner"
	app.Commands = []cli.Command{cmd.Scan}
	app.Flags = append(app.Flags, cmd.Scan.Flags...)
	err := app.Run(os.Args)
	_ = err
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
