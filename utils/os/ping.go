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

	"k8s.io/klog/v2"
)

type Pinger interface {
	Ping(host string, samples uint8) (avgDuration uint16, err error)
}

func NewPinger() Pinger {
	return newPinger(unixPingerImpl{}, windowsPingerImpl{})
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
		klog.V(1).Info("PING: windows")
		duration, err = p.pingWin.ping(host, samples)
	default:
		klog.Infof("PING: unix [%s]", runtime.GOOS)
		duration, err = p.pingUnix.ping(host, samples)
	}
	if duration == 0 {
		duration = 1
		klog.V(1).Infof("%s: sub ms ping rounded up to 1ms", host)
	}
	return duration, err
}

func (p windowsPingerImpl) ping(host string, samples uint8) (uint16, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5000*time.Millisecond)
	defer cancel()
	cmdline := []string{"ping", "-n", fmt.Sprint(samples), host}
	klog.Infof("Executing command line: %v", cmdline)
	cmd := exec.CommandContext(ctx, cmdline[0], cmdline[1:]...)
	out, err := cmd.CombinedOutput()
	klog.V(1).Infof("Output:\n%s", string(out))
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
