/*
   piawgcli
   Copyright (C) 2021  Derek Battams <derek@battams.ca>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"gitlab.com/ddb_db/piawgcli/actions"
	"gitlab.com/ddb_db/piawgcli/context"
	"k8s.io/klog/v2"
)

var cli struct {
	CaseSensitive bool                   `help:"case sensitive searching" default:"1" negatable`
	Debug         uint8                  `help:"log verbosity; higher=more log output" short:"v" default:"0"`
	LogFile       string                 `help:"log output to file instead of stdout"`
	ServerList    string                 `hidden help:"PIA server list source" default:"https://serverlist.piaservers.net/vpninfo/servers/v4"`
	ShowRegions   actions.ShowRegionsCmd `cmd help:"Show available regions"`
}

func main() {
	ctx := kong.Parse(&cli)
	klog.InitFlags(nil)
	flag.Set("v", fmt.Sprintf("%d", cli.Debug))
	if len(cli.LogFile) > 0 {
		logFile := prepLogFile()
		defer logFile.Close()
		flag.Set("logtostderr", "false")
		flag.Set("log_file", cli.LogFile)
	}
	defer klog.Flush()
	flag.Parse()
	err := ctx.Run(&context.Context{
		CaseSensitive: cli.CaseSensitive,
		Debug:         uint8(cli.Debug),
		ServerList:    cli.ServerList})
	ctx.FatalIfErrorf(err)
}

func prepLogFile() *os.File {
	f, err := os.Create(cli.LogFile)
	if err != nil {
		klog.Fatalf("unable to create log file: %s", cli.LogFile)
	}
	return f
}
