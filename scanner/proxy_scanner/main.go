package main

import (
	"SafeDp/scanner/proxy_scanner/cmd"
	"os"
	"runtime"

	"github.com/urfave/cli"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	app := cli.NewApp()
	app.Name = "ProxyScanner"
	app.Usage = "A SOCKS4/SOCKS4a/SOCKS5/HTTP/HTTPS proxy scanner"
	app.Version = "1.0.0"
	app.Commands = []cli.Command{cmd.Scan}
	app.Flags = append(app.Flags, cmd.Scan.Flags...)
	_ = app.Run(os.Args)
}
