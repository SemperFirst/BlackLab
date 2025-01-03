package cmd

import (
	"SafeDp/scanner/password_crack/util"
	"github.com/urfave/cli"
)

var Scan = cli.Command{
	Name: "Scan",
	Usage: "Start to crack weak password",
	Description: "Start to crack weak password",
	Action: util.Scan,
	Flags: []cli.Flag{
}