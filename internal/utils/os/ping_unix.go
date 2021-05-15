// +build !windows

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
	"time"

	"k8s.io/klog/v2"
)

type pingerImpl struct{}

func (p pingerImpl) Ping(host string, samples uint8) (uint16, error) {
	// TODO ping timeout should be configurable
	ctx, cancel := context.WithTimeout(context.Background(), 5000*time.Millisecond)
	defer cancel()
	cmdline := []string{"ping", "-c", fmt.Sprint(samples), host}
	cmd := exec.CommandContext(ctx, cmdline[0], cmdline[1:]...)
	out, err := cmd.CombinedOutput()
	klog.V(4).Infof("Executing command line: %v [rc=%d]", cmdline, cmd.ProcessState.ExitCode())
	klog.V(5).Infof("Output:\n%s", string(out))
	var ping uint16
	if ctx.Err() != nil {
		return ping, ctx.Err()
	} else {
		if err != nil {
			if !klog.V(5).Enabled() {
				klog.Errorf("Output:\n%s", string(out))
			}
			return ping, fmt.Errorf("ping failed: %w", err)
		}
		ping, err = parsePingTimeUnix(string(out))
	}
	return ping, err
}
