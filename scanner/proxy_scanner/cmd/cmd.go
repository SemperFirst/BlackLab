package cmd

import (
	"github.com/urfave/cli"
	"SafeDp/scanner/proxy_scanner/proxy"
)

var Scan = cli.Command{
	Name: 	  "ProxyScan",
	Usage: 	  "Start to scan proxy",
	Description: "Start to scan proxy",
	Action:   proxy.Scan,
	Flags: []cli.Flag{
		boolFlag("debug, d", "debug mode"),
		intFlag("timeout, t", 5, "timeout"),
		intFlag("scan_num, n", 100, "scan num"),
		stringFlag("filename, f", "iplist.txt", "filename"),
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


