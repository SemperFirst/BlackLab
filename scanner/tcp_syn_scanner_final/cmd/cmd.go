package cmd

import (
	"github.com/urfave/cli"
	"SafeDp/scanner/tcp_syn_scanner_final/util"
)

var Scan = cli.Command{
	Name: "scan",
	Usage: "Start to scan",
	Description: "Start to scan",
	Action: util.Scan,
	Flags: []cli.Flag{
		stringFlag("iplist, i", "", "IP list"),
		stringFlag("port, p", "", "Port list"),
		stringFlag("mode, m", "", "Scan mode"),
		intFlag("timeout, t", 3, "Timeout"),
		intFlag("concurrency, c", 1000, "Concurrency"),
	},
}

func stringFlag(name, value, usage string) cli.StringFlag {
	return cli.StringFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

func boolFlag(name, usage string) cli.BoolFlag {
	return cli.BoolFlag{
		Name:  name,
		Usage: usage,
	}
}

func intFlag(name string, value int, usage string) cli.IntFlag {
	return cli.IntFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}
