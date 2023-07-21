/*
piawgcli
Copyright (C) 2021-2023  Derek Battams <derek@battams.ca>

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
package actions

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	_ "embed"

	"gitlab.com/ddb_db/piawgcli/internal/appstate"
	"gitlab.com/ddb_db/piawgcli/internal/net/piaclient"
	"k8s.io/klog/v2"
)

type CreateConfigCmd struct {
	PiaId        string `required help:"PIA user id" placeholder:"ID"`
	PiaPassword  string `required help:"PIA password" placeholder:"PWD"`
	PiaRegionId  string `required help:"PIA region id to connect to; use show-regions command to get the region id" placeholder:"ID"`
	IgnorePiaDns bool   `help:"Do not set DNS servers to PIA servers in generated configuration"`
	Output       string `help:"write wg config to file instead of stdout" placeholder:"FILE"`
}

//go:embed assets/wg.conf.tmpl
var wgConfTmpl string

// TODO break this down (verify tmpl output, etc)
func (cmd *CreateConfigCmd) Run(state *appstate.State) error {
	pia := piaclient.New(state.ServerList)
	piaInterface, err := pia.CreateTunnel(cmd.PiaId, cmd.PiaPassword, cmd.PiaRegionId)
	if err != nil {
		return err
	}
	if cmd.IgnorePiaDns {
		klog.V(4).Info("ignoring PIA DNS servers")
		piaInterface.DnsServers = nil
	}

	var output *os.File
	if len(cmd.Output) > 0 {
		klog.V(4).Infof("writing config to %s", cmd.Output)
		output, err = os.Create(cmd.Output)
		if err != nil {
			return err
		}
		defer output.Close()
	} else {
		klog.V(4).Info("writing config to stdout")
		output = os.Stdout
	}

	var result string
	result, err = processTemplate(wgConfTmpl, piaInterface)
	if err != nil {
		return fmt.Errorf("template processing failed: %w", err)
	}

	_, err = output.WriteString(result)
	if err != nil {
		return fmt.Errorf("io error writing output: %w", err)
	}

	return nil
}

func processTemplate(tmplSource string, bindings interface{}) (string, error) {
	tmpl, err := template.New("wgconf").
		Funcs(template.FuncMap{"join": strings.Join}).
		Parse(tmplSource)
	if err != nil {
		return "", fmt.Errorf("wg template parsing failed: %w", err)
	}
	output := strings.Builder{}
	err = tmpl.Execute(&output, bindings)
	return output.String(), err
}
