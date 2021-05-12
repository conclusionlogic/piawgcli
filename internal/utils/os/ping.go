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
package os

import (
	"fmt"
	"regexp"
	"strconv"
)

type Pinger interface {
	Ping(host string, samples uint8) (avgDuration uint16, err error)
}

type abstractPinger struct {
	pinger Pinger
}

func NewPinger() Pinger {
	return abstractPinger{
		pinger: pingerImpl{},
	}
}

func (p abstractPinger) Ping(host string, samples uint8) (uint16, error) {
	ping, err := p.pinger.Ping(host, samples)
	if err != nil {
		return 9999, fmt.Errorf("ping failed: %w", err)
	} else if ping == 0 {
		ping = 1
	}
	return ping, err
}

func parsePingTimeUnix(output string) (uint16, error) {
	var err error
	var ping uint16
	re := regexp.MustCompile(`= \d+(?:\.\d+)?\/(\d+)(?:\.\d+)?\/.+ ms`)
	matches := re.FindStringSubmatch(output)
	if len(matches) > 1 {
		var val uint64
		val, err = strconv.ParseUint(matches[1], 10, 16)
		ping = uint16(val)
	} else {
		err = fmt.Errorf("unable to find ping timings in output [unix]")
	}
	return ping, err
}

func parsePingTimeWindows(output string) (uint16, error) {
	var err error
	var ping uint16
	re := regexp.MustCompile(`Average = (\d+)ms`)
	matches := re.FindStringSubmatch(output)
	if len(matches) > 1 {
		var val uint64
		val, err = strconv.ParseUint(matches[1], 10, 16)
		ping = uint16(val)
	} else {
		err = fmt.Errorf("unable to find ping timings in output [windows]")
	}
	return ping, err
}
