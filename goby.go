package main

import (
	"os"
	"runtime"

	cli "gopkg.in/urfave/cli.v1"

	"github.com/MessageDream/goby/cmd"
	"github.com/MessageDream/goby/module/setting"
)

const APP_VER = "0.0.1 Beta"

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	setting.AppVer = APP_VER
}

func main() {
	app := cli.NewApp()
	app.Name = "goby"
	app.Usage = "react native code push server"
	app.Version = APP_VER
	app.Commands = []cli.Command{
		cmd.CmdWeb,
	}
	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Run(os.Args)
}
