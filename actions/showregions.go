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
package actions

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/jamesrr39/semaphore"
	"gitlab.com/ddb_db/piawgcli/context"
	"gitlab.com/ddb_db/piawgcli/utils/net"
	"gitlab.com/ddb_db/piawgcli/utils/os"
	"k8s.io/klog/v2"
)

type ShowRegionsCmd struct {
	Ping      bool   `optional help:"ping each region and sort results by ping time" default:"0"`
	SortBy    string `optional help:"sort results by given field" enum:"id,name" default:"name"`
	SortOrder string `optional help:"sort results ASCending or DESCending order" enum:"asc,desc" default:"asc"`
	Search    string `optional help:"find regions containing search term"`
	Threads   uint8  `optional help:"max number of worker threads for pinging regions" default:"8"`
	Samples   uint8  `optional help:"number of samples to take when pinging regions" default:"3"`
}

func (cmd *ShowRegionsCmd) Run(ctx *context.Context) error {
	action := showRegionsAction{
		pinger:     os.NewPinger(),
		urlFetcher: net.NewUrlFetcher(),
		cmd:        cmd,
		ctx:        ctx,
	}
	return action.run()
}

type piaRegion struct {
	Id   string
	Name string
	Dns  string
	Ping uint16
}

type piaRegions struct {
	Regions []piaRegion
}

type showRegionsAction struct {
	cmd        *ShowRegionsCmd
	ctx        *context.Context
	pinger     os.Pinger
	urlFetcher net.UrlFetcher
}

func (action showRegionsAction) run() error {
	cmd := action.cmd
	pia, err := action.parseServerList()
	if err != nil {
		return err
	}
	if len(cmd.Search) > 0 {
		klog.V(5).Infof("applying region filter: %s", cmd.Search)
		pia.Regions = action.filter(pia.Regions)
	}
	if cmd.Ping {
		klog.V(5).Info("pinging regions")
		action.pingRegions(pia.Regions)
	}
	action.sortRegions(pia.Regions)
	action.printRegions(pia.Regions)
	return nil
}

func (action showRegionsAction) printRegions(regions []piaRegion) {
	fmt.Printf("%-24s %-18s %-9s\n", "NAME", "ID", "PING (ms)")
	fmt.Printf("%43s\n", strings.Repeat("=", 53))
	for _, r := range regions {
		ping := fmt.Sprint(r.Ping)
		if r.Ping == 0 {
			ping = ""
		}
		fmt.Printf("%-24s %-18s %9s\n", r.Name, r.Id, ping)
	}
}

func (action showRegionsAction) sortRegions(regions []piaRegion) {
	cmd := action.cmd
	sort.Slice(regions,
		func(i, j int) bool {
			if cmd.SortOrder != "asc" {
				klog.V(5).Info("sort order: desc")
				tmp := i
				i = j
				j = tmp
			} else {
				klog.V(5).Info("sort order: asc")
			}
			if cmd.Ping {
				klog.V(5).Info("sort key: ping")
				return regions[i].Ping < regions[j].Ping
			} else if cmd.SortBy == "name" {
				klog.V(5).Info("sort key: name")
				return regions[i].Name < regions[j].Name
			} else {
				klog.V(5).Info("sort name: id")
				return regions[i].Id < regions[j].Id
			}
		})
}

func (action showRegionsAction) parseServerList() (piaRegions, error) {
	payload, err := action.urlFetcher.FetchString(action.ctx.ServerList)
	if err != nil {
		return piaRegions{}, err
	}
	body := action.extractJsonBody(payload)
	var pia piaRegions
	err = json.Unmarshal(body, &pia)
	return pia, err
}

func (action showRegionsAction) extractJsonBody(payload string) []byte {
	// the endpoint pads the json response with a signature blob of some kind so we must extract out only the json data in the response
	lastBrace := strings.LastIndex(payload, "}")
	klog.V(4).Infof("last brace: %d", lastBrace)
	body := []byte(payload[0 : lastBrace+1])
	klog.V(4).Infof("region payload: %s", body[:70])
	return body
}

func (action showRegionsAction) pingRegions(regions []piaRegion) {
	sem := semaphore.NewSemaphore(uint(action.cmd.Threads))
	for i := range regions {
		offset := i
		sem.Add()
		go func() {
			defer sem.Done()
			regions[offset] = action.doPing(regions[offset])
		}()
	}
	klog.V(4).Infof("waiting on ~%d workers", sem.CurrentlyRunning())
	sem.Wait()
}

func (action showRegionsAction) isMatch(r piaRegion) bool {
	searchTerm := action.cmd.Search
	var searchName, searchId, searchPredicate string
	if !action.ctx.CaseSensitive {
		searchName = strings.ToLower(r.Name)
		searchId = strings.ToLower(r.Id)
		searchPredicate = strings.ToLower(searchTerm)
	} else {
		searchName = r.Name
		searchId = r.Id
		searchPredicate = searchTerm
	}
	klog.V(4).Infof("n=%s, i=%s, p=%s", searchName, searchId, searchPredicate)
	return strings.Contains(searchName, searchPredicate) || strings.Contains(searchId, searchPredicate)
}

func (action showRegionsAction) doPing(r piaRegion) piaRegion {
	ping, err := action.pinger.Ping(r.Dns, action.cmd.Samples)
	if err != nil {
		klog.Errorf("ping failed: %s\n%v", r.Name, err)
		ping = 10000
	}
	region := piaRegion{Id: r.Id, Name: r.Name, Ping: ping, Dns: r.Dns}
	klog.V(5).Infof("region pinged: %v", region)
	return region
}

func (action showRegionsAction) filter(regions []piaRegion) []piaRegion {
	var filtered []piaRegion
	for _, r := range regions {
		if action.isMatch(r) {
			filtered = append(filtered, r)
			klog.V(4).Infof("region %s: matched filter", r.Name)
		} else {
			klog.V(4).Infof("region %s: did not match filter", r.Name)
		}
	}
	return filtered
}
