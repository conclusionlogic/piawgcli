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
	"testing"

	"github.com/pkg/errors"
)

type ping struct {
	duration uint16
	err      error
}

func (p ping) ping(string, uint8) (uint16, error) {
	return p.duration, p.err
}

func TestPingInterpretations(t *testing.T) {
	var tests = []struct {
		input    uint16
		expected uint16
	}{
		{0, 1},
		{1, 1},
		{2, 2},
		{10, 10},
		{100, 100},
		{500, 500},
		{10000, 10000},
	}
	for i, tc := range tests {
		pingImpl := ping{
			duration: tc.input,
			err:      nil,
		}
		p := newPinger(pingImpl, pingImpl)
		d, err := p.Ping("foo", 1)
		if err != nil || d != tc.expected {
			t.Errorf("itr %d: expect %d, received %d", i, tc.expected, d)
		}
	}
}

func TestFailedPing(t *testing.T) {
	pingImpl := ping{
		duration: 0,
		err:      errors.Errorf("oops"),
	}
	p := newPinger(pingImpl, pingImpl)
	_, err := p.Ping("foo", 1)
	if err == nil || err.Error() != "oops" {
		t.Errorf("did not receive expeted error response")
	}
}
