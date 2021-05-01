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
package net

import (
	"testing"
)

func TestFetchBytesWithBadUrl(t *testing.T) {
	_, err := NewUrlFetcher().FetchBytes("http://127.0.0.1:4000")
	if err == nil {
		t.Errorf("Invalid url expected to return an error")
	}
}

func TestFetchBytesWithValidUrl(t *testing.T) {
	payload, err := NewUrlFetcher().FetchBytes("http://google.com")
	if err != nil || len(payload) == 0 {
		t.Errorf("Unexpected error trying to download valid url")
	}
}
