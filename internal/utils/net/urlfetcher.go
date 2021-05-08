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
	"io"
	"net/http"
)

type UrlFetcher interface {
	FetchString(url string) (body string, err error)
	FetchBytes(url string) (body []byte, err error)
}

func NewUrlFetcher() UrlFetcher {
	return urlFetcherImpl{}
}

type urlFetcherImpl struct{}

func (u urlFetcherImpl) FetchString(url string) (string, error) {
	payload, err := u.FetchBytes(url)
	return string(payload), err
}

func (u urlFetcherImpl) FetchBytes(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return body, err
}
