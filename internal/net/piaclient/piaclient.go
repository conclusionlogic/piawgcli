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
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/go-resty/resty/v2"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"k8s.io/klog/v2"
)

//https://github.com/go-resty/resty
type PiaClient interface {
	GetRegions() (PiaRegions, error)
	GetRegionById(id string) (PiaRegion, error)
	GetAuthToken(piaId string, piaPassword string, piaRegion PiaRegion) (string, error)
	CreateTunnel(piaId string, piaPassword string, piaRegionId string) (io.Reader, error)
}

type PiaRegion struct {
	Id      string
	Name    string
	Dns     string
	Servers PiaServers
	Ping    uint16
}

type PiaServers struct {
	Wg   []PiaServer
	Meta []PiaServer
}

type PiaServer struct {
	Ip string
	Cn string
}

type PiaRegions struct {
	Regions []PiaRegion
}

type piaClientImpl struct {
	regionUrl string
	http      map[string]*resty.Client
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
	c := piaClientImpl{
		regionUrl: serverListUrl,
		http:      make(map[string]*resty.Client),
	}
	c.http["_"] = resty.New()
	return c
}

func (clnt piaClientImpl) getDefaultHttp() *resty.Client {
	return clnt.http["_"]
}

func (clnt piaClientImpl) getHttpForRegion(region PiaRegion) *resty.Client {
	c := clnt.http[region.Id]
	if c == nil {
		c = resty.New().
			SetTLSClientConfig(&tls.Config{
				ServerName: region.Servers.Meta[0].Cn,
			}).
			SetRootCertificate("assets/pia.pem")
		clnt.http[region.Id] = c
	}
	return c
}

func (clnt piaClientImpl) GetAuthToken(id string, pwd string, region PiaRegion) (string, error) {
	url := fmt.Sprintf("https://%s/authv3/generateToken", region.Servers.Meta[0].Ip)
	resp, err := clnt.getHttpForRegion(region).R().
		SetBasicAuth(id, pwd).
		Get(url)
	if err != nil {
		err = fmt.Errorf("token fetch failed: %w", err)
		return "", err
	} else {
		httpStatus := resp.StatusCode()
		if httpStatus == 403 {
			return "", fmt.Errorf("invalid PIA credentials")
		} else if httpStatus < 200 || httpStatus > 299 {
			return "", fmt.Errorf("invalid auth token response: %d", httpStatus)
		}
		klog.V(4).Info(resp.String())
	}
	var jsonResp struct {
		Status string
		Token  string
	}
	err = json.Unmarshal(resp.Body(), &jsonResp)
	if err != nil {
		err = fmt.Errorf("json parse of auth token failed: %w", err)
	}
	if jsonResp.Status != "OK" {
		err = fmt.Errorf("invalid auth token response: %s [%d]", jsonResp.Status, resp.StatusCode())
	}
	return jsonResp.Token, err
}

func (clnt piaClientImpl) GetRegions() (PiaRegions, error) {
	resp, err := clnt.getDefaultHttp().R().Get(clnt.regionUrl)
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
		if r.Id == id && len(r.Servers.Meta[0].Ip) > 0 && len(r.Servers.Wg[0].Ip) > 0 {
			return r, nil
		}
	}
	return PiaRegion{}, newUnknownRegionError(fmt.Sprintf("unknown region id: %s", id))
}

func (clnt piaClientImpl) CreateTunnel(piaId string, piaPwd string, piaRegionId string) (io.Reader, error) {
	privKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, fmt.Errorf("wg key generation failed: %w", err)
	}
	pubKey := privKey.PublicKey()
	r, err := clnt.GetRegionById(piaRegionId)
	if err != nil {
		return nil, err
	}
	authToken, err := clnt.GetAuthToken(piaId, piaPwd, r)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("https://%s:1337/addKey", r.Servers.Wg[0].Ip)
	http := clnt.getHttpForRegion(r)
	resp, err := http.R().
		SetQueryParams(map[string]string{
			"pubkey": pubKey.String(),
			"pt":     authToken,
		}).Get(url)
	if err != nil {
		return nil, fmt.Errorf("addKey failed: %w", err)
	}
	return strings.NewReader(resp.String()), nil
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
