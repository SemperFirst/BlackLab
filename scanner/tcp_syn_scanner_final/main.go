package main 

import (
	"os"
	"runtime"

	"github.com/urfave/cli"
	"SafeDp/scanner/tcp_syn_scanner_final/cmd"
)

func main() {
	app := cli.NewApp()
	app.Name = "port_scanner"
	app.Author = "SemperFi"
	app.Usage = "A tcp/syn port scanner"
	app.Version = "1.0.0"
	app.Commands = []cli.Command{cmd.Scan}
	app.Flags = append(app.Flags, cmd.Scan.Flags...)
	err := app.Run(os.Args)
	_ = err
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
