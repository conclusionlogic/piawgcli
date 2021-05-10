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
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/ddb_db/piawgcli/internal/net/piaclient"
	"gitlab.com/ddb_db/piawgcli/internal/utils/os"
)

// https://github.com/stretchr/testify
// TODO table driven tests

var regions = []piaclient.PiaRegion{
	{Id: "a", Name: "f", Ping: 0},
	{Id: "b", Name: "b", Ping: 1},
	{Id: "c", Name: "e", Ping: 2},
	{Id: "d", Name: "c", Ping: 3},
	{Id: "e", Name: "d", Ping: 4},
	{Id: "f", Name: "a", Ping: 5},
}

var regionsMixedCase = []piaclient.PiaRegion{
	{Id: "a", Name: "f", Ping: 0},
	{Id: "B", Name: "b", Ping: 1},
	{Id: "c", Name: "E", Ping: 2},
	{Id: "D", Name: "C", Ping: 3},
	{Id: "e", Name: "d", Ping: 4},
	{Id: "F", Name: "A", Ping: 5},
}

func TestSortRegionsAscendingByPing(t *testing.T) {
	var c = ShowRegionsCmd{
		Ping:      true,
		SortOrder: "asc",
	}
	var a = showRegionsAction{
		cmd:      &c,
		appState: nil,
		pinger:   os.NewPinger(),
		pia:      nil,
	}
	a.sortRegions(regions)
	result := func(r []piaclient.PiaRegion) []uint16 {
		var vals []uint16
		for _, it := range r {
			vals = append(vals, it.Ping)
		}
		return vals
	}(regions)
	require.Equal(t, []uint16{0, 1, 2, 3, 4, 5}, result, "not equal")
}

func TestSortRegionsAscendingById(t *testing.T) {
	var c = ShowRegionsCmd{
		Ping:      false,
		SortBy:    "id",
		SortOrder: "asc",
	}
	var a = showRegionsAction{
		cmd:      &c,
		appState: nil,
		pinger:   os.NewPinger(),
		pia:      nil,
	}
	a.sortRegions(regions)
	result := func(r []piaclient.PiaRegion) []string {
		var vals []string
		for _, it := range r {
			vals = append(vals, it.Id)
		}
		return vals
	}(regions)
	require.Equal(t, []string{"a", "b", "c", "d", "e", "f"}, result, "not equal")
}

func TestSortRegionsAscendingByName(t *testing.T) {
	var c = ShowRegionsCmd{
		Ping:      false,
		SortBy:    "name",
		SortOrder: "asc",
	}
	var a = showRegionsAction{
		cmd:      &c,
		appState: nil,
		pinger:   os.NewPinger(),
		pia:      nil,
	}
	a.sortRegions(regions)
	result := func(r []piaclient.PiaRegion) []string {
		var vals []string
		for _, it := range r {
			vals = append(vals, it.Name)
		}
		return vals
	}(regions)
	require.Equal(t, []string{"a", "b", "c", "d", "e", "f"}, result, "not equal")
}

func TestSortRegionsDescendingByPing(t *testing.T) {
	var c = ShowRegionsCmd{
		Ping:      true,
		SortOrder: "desc",
	}
	var a = showRegionsAction{
		cmd:      &c,
		appState: nil,
		pinger:   os.NewPinger(),
		pia:      nil,
	}
	a.sortRegions(regions)
	result := func(r []piaclient.PiaRegion) []uint16 {
		var vals []uint16
		for _, it := range r {
			vals = append(vals, it.Ping)
		}
		return vals
	}(regions)
	require.Equal(t, []uint16{5, 4, 3, 2, 1, 0}, result, "not equal")
}

func TestSortRegionsDesccendingById(t *testing.T) {
	var c = ShowRegionsCmd{
		Ping:      false,
		SortBy:    "id",
		SortOrder: "desc",
	}
	var a = showRegionsAction{
		cmd:      &c,
		appState: nil,
		pinger:   os.NewPinger(),
		pia:      nil,
	}
	a.sortRegions(regions)
	result := func(r []piaclient.PiaRegion) []string {
		var vals []string
		for _, it := range r {
			vals = append(vals, it.Id)
		}
		return vals
	}(regions)
	require.Equal(t, []string{"f", "e", "d", "c", "b", "a"}, result, "not equal")
}

func TestSortRegionsDescendingByName(t *testing.T) {
	var c = ShowRegionsCmd{
		Ping:      false,
		SortBy:    "name",
		SortOrder: "desc",
	}
	var a = showRegionsAction{
		cmd:      &c,
		appState: nil,
		pinger:   os.NewPinger(),
		pia:      nil,
	}
	a.sortRegions(regions)
	result := func(r []piaclient.PiaRegion) []string {
		var vals []string
		for _, it := range r {
			vals = append(vals, it.Name)
		}
		return vals
	}(regions)
	require.Equal(t, []string{"f", "e", "d", "c", "b", "a"}, result, "not equal")
}
