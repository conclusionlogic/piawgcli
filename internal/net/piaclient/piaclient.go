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
package piaclient

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/go-resty/resty/v2"
	"k8s.io/klog/v2"
)

//https://github.com/go-resty/resty
type PiaClient interface {
	GetRegions() (PiaRegions, error)
	GetRegionById(id string) (PiaRegion, error)
	CreateTunnel(piaId string, piaPassword string, piaRegionId string) (io.Reader, error)
}

type PiaRegion struct {
	Id   string
	Name string
	Dns  string
	Ping uint16
}

type PiaRegions struct {
	Regions []PiaRegion
}

type piaClientImpl struct {
	regionUrl string
	http      *resty.Client
}

type UnknownRegionError struct {
	errMsg string
}

func (err UnknownRegionError) Error() string {
	return err.errMsg
}

func newUnknownRegionError(msg string) UnknownRegionError {
	return UnknownRegionError{
		errMsg: msg,
	}
}

func New(serverListUrl string) PiaClient {
	http := resty.New()
	c := piaClientImpl{
		regionUrl: serverListUrl,
		http:      http,
	}
	return c
}

func (clnt piaClientImpl) GetRegions() (PiaRegions, error) {
	resp, err := clnt.http.R().Get(clnt.regionUrl)
	if err != nil {
		return PiaRegions{}, fmt.Errorf("region url fetch failed: %w", err)
	}
	return parsePiaRegionJsonBody(resp.String())
}

func (clnt piaClientImpl) GetRegionById(id string) (PiaRegion, error) {
	regions, err := clnt.GetRegions()
	if err != nil {
		return PiaRegion{}, err
	}
	for _, r := range regions.Regions {
		if r.Id == id {
			return r, nil
		}
	}
	return PiaRegion{}, newUnknownRegionError(fmt.Sprintf("unknown region id: %s", id))
}

func (clnt piaClientImpl) CreateTunnel(piaId string, piaPwd string, piaRegionId string) (io.Reader, error) {
	return nil, nil
}

func parsePiaRegionJsonBody(payload string) (PiaRegions, error) {
	// the endpoint pads the json response with an undocumented signature blob of some kind so we must extract out only the json data in the response
	lastBrace := strings.LastIndex(payload, "}")
	klog.V(4).Infof("last brace: %d", lastBrace)
	body := []byte(payload[0 : lastBrace+1])
	klog.V(4).Infof("region payload: %s", body[:70])
	val := PiaRegions{}
	err := json.Unmarshal(body, &val)
	if err != nil {
		err = fmt.Errorf("region data parse failed: %w", err)
	}
	return val, err
}
