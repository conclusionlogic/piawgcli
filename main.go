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
	"github.com/alecthomas/kong"
	"gitlab.com/ddb_db/piawgcli/actions"
	"gitlab.com/ddb_db/piawgcli/context"
)

var cli struct {
	CaseSensitive bool                   `help:"case sensitive searching" default:"1" negatable`
	Debug         bool                   `help:"debug mode"`
	ServerList    string                 `hidden help:"PIA server list source" default:"https://serverlist.piaservers.net/vpninfo/servers/v4"`
	ShowRegions   actions.ShowRegionsCmd `cmd help:"Show available regions"`
}

func main() {
	ctx := kong.Parse(&cli)
	err := ctx.Run(&context.Context{
		CaseSensitive: cli.CaseSensitive,
		Debug:         cli.Debug,
		ServerList:    cli.ServerList})
	ctx.FatalIfErrorf(err)
}
