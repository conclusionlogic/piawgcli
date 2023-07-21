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
package os

import (
	"testing"

	"github.com/pkg/errors"
)

type ping struct {
	duration uint16
	err      error
}

func (p ping) Ping(string, uint8) (uint16, error) {
	return p.duration, p.err
}

func newPinger(duration uint16, err error) abstractPinger {
	return abstractPinger{
		pinger: ping{
			duration: duration,
			err:      err,
		},
	}
}

func TestPingInterpretations(t *testing.T) {
	var tests = []struct {
		input    uint16
		expected uint16
	}{
		{0, 1}, // pings <1ms are rounded up to 1ms
		{1, 1},
		{2, 2},
		{10, 10},
		{100, 100},
		{500, 500},
		{10000, 10000},
	}
	for i, tc := range tests {
		p := newPinger(tc.input, nil)
		d, err := p.Ping("foo", 1)
		if err != nil || d != tc.expected {
			t.Errorf("itr %d: expect %d, received %d", i, tc.expected, d)
		}
	}
}

func TestFailedPing(t *testing.T) {
	_, err := newPinger(0, errors.Errorf("oops")).Ping("foo", 1)
	if err == nil || err.Error()[0:12] != "ping failed:" {
		t.Errorf("did not receive expected error response")
	}
}

func TestUnixPingParse(t *testing.T) {
	var tests = []struct {
		input    string
		expected uint16
	}{
		{"rtt min/avg/max/mdev = 9.617/9.909/10.584/0.344 ms", 9},
		{"rtt min/avg/max/mdev = 0.617/0.909/4.584/0.344 ms", 0},
		{"rtt min/avg/max/mdev = 0.617/1.909/4.584/1.344 ms", 1},
		{"rtt min/avg/max/mdev = 0.617/15.909/4.584/0.344 ms", 15},
		{"rtt min/avg/max/mdev = 0.617/232.909/4.584/0.344 ms", 232},
		{"rtt min/avg/max/mdev = 0.617/1000.909/4.584/0.344 ms", 1000},
	}

	for i, tc := range tests {
		result, err := parsePingTimeUnix(tc.input)
		if err != nil {
			t.Errorf("unexpected error [itr=%d]: %s", i, err.Error())
		}
		if result != tc.expected {
			t.Errorf("expected %d, got %d [itr=%d]", tc.expected, result, i)
		}
	}
}

func TestWindowsPingParse(t *testing.T) {
	var tests = []struct {
		input    string
		expected uint16
	}{
		{"    Minimum = 9ms, Maximum = 10ms, Average = 9ms", 9},
		{"    Minimum = 9ms, Maximum = 10ms, Average = 0ms", 0},
		{"    Minimum = 9ms, Maximum = 10ms, Average = 1ms", 1},
		{"    Minimum = 9ms, Maximum = 10ms, Average = 15ms", 15},
		{"    Minimum = 9ms, Maximum = 10ms, Average = 232ms", 232},
		{"    Minimum = 9ms, Maximum = 10ms, Average = 1000ms", 1000},
	}

	for i, tc := range tests {
		result, err := parsePingTimeWindows(tc.input)
		if err != nil {
			t.Errorf("unexpected error [itr=%d]: %s", i, err.Error())
		}
		if result != tc.expected {
			t.Errorf("expected %d, got %d [itr=%d]", tc.expected, result, i)
		}
	}
}
