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
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"time"
)

type Pinger interface {
	Ping(host string, samples uint8) (avgDuration uint16, err error)
}

func NewPinger() Pinger {
	return pingerImpl{
		pingUnix: unixPingerImpl{},
		pingWin:  windowsPingerImpl{},
	}
}

func newPinger(unix unixPinger, win windowsPinger) Pinger {
	return pingerImpl{
		pingUnix: unix,
		pingWin:  win,
	}
}

type windowsPinger interface {
	ping(host string, samples uint8) (avgDuration uint16, err error)
}

type unixPinger interface {
	ping(host string, sampels uint8) (avgDuration uint16, err error)
}

type windowsPingerImpl struct{}
type unixPingerImpl struct{}

type pingerImpl struct {
	pingUnix unixPinger
	pingWin  windowsPinger
}

func (p pingerImpl) Ping(host string, samples uint8) (uint16, error) {
	var duration uint16
	var err error
	switch runtime.GOOS {
	case "windows":
		duration, err = p.pingWin.ping(host, samples)
	default:
		duration, err = p.pingUnix.ping(host, samples)
	}
	if duration == 0 {
		duration = 1
	}
	return duration, err
}

func (p windowsPingerImpl) ping(host string, samples uint8) (uint16, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5000*time.Millisecond)
	defer cancel()
	cmd := exec.CommandContext(ctx, "ping", "-n", fmt.Sprint(samples), host)
	out, err := cmd.CombinedOutput()
	var ping uint16
	if ctx.Err() != nil {
		err = ctx.Err()
	} else {
		re := regexp.MustCompile(`Average = (\d+)ms`)
		matches := re.FindStringSubmatch(string(out))
		if len(matches) > 1 {
			var val uint64
			val, err = strconv.ParseUint(matches[1], 10, 16)
			ping = uint16(val)
		}
	}
	return ping, err
}

func (p unixPingerImpl) ping(host string, samples uint8) (uint16, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5000*time.Millisecond)
	defer cancel()
	cmd := exec.CommandContext(ctx, "ping", "-c", string(samples), host)
	out, err := cmd.CombinedOutput()
	var ping uint16
	if ctx.Err() != nil {
		err = ctx.Err()
	} else {
		re := regexp.MustCompile(`= \d+(?:\.\d+)?\/(\d+)(?:\.\d+)?\/.+ ms`)
		matches := re.FindStringSubmatch(string(out))
		if len(matches) > 1 {
			var val uint64
			val, err = strconv.ParseUint(matches[1], 10, 16)
			ping = uint16(val)
		}
	}
	return ping, err
}
