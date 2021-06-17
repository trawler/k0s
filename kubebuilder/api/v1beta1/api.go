/*
Copyright 2021 k0s authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package v1beta1

import (
	"fmt"
	"net"

	"github.com/asaskevich/govalidator"
)

var _ Validateable = (*K0sAPISpec)(nil)

// APISpec ...
type K0sAPISpec struct {
	Address         string            `yaml:"address"`
	Port            int               `yaml:"port"`
	K0sAPIPort      int               `yaml:"k0sApiPort,omitempty"`
	ExternalAddress string            `yaml:"externalAddress,omitempty"`
	SANs            []string          `yaml:"sans"`
	ExtraArgs       map[string]string `yaml:"extraArgs,omitempty"`
}

// DefaultAPISpec default settings for api
func DefaultAPISpec() *K0sAPISpec {
	// Collect all nodes addresses for sans
	// addresses, _ := util.AllAddresses()
	// publicAddress, _ := util.FirstPublicAddress()
	return &K0sAPISpec{
		Port:       6443,
		K0sAPIPort: 9443,
		SANs:       []string{"addresses"},
		Address:    "publicAddress",
		ExtraArgs:  make(map[string]string),
	}
}

// APIAddress ...
func (a *K0sAPISpec) APIAddress() string {
	if a.ExternalAddress != "" {
		return a.ExternalAddress
	}
	return a.Address
}

// APIAddressURL returns kube-apiserver external URI
func (a *K0sAPISpec) APIAddressURL() string {
	return a.getExternalURIForPort(a.Port)
}

// IsIPv6String returns if ip is IPv6.
func IsIPv6String(ip string) bool {
	netIP := net.ParseIP(ip)
	return netIP != nil && netIP.To4() == nil
}

// K0sControlPlaneAPIAddress returns the controller join APIs address
func (a *K0sAPISpec) K0sControlPlaneAPIAddress() string {
	return a.getExternalURIForPort(a.K0sAPIPort)
}

func (a *K0sAPISpec) getExternalURIForPort(port int) string {
	addr := a.Address
	if a.ExternalAddress != "" {
		addr = a.ExternalAddress
	}
	if IsIPv6String(addr) {
		return fmt.Sprintf("https://[%s]:%d", addr, port)
	}
	return fmt.Sprintf("https://%s:%d", addr, port)
}

// Sans return the given SANS plus all local adresses and externalAddress if given
func (a *K0sAPISpec) Sans() []string {
	// sans, _ := util.AllAddresses()
	sans := []string{"sans"}
	sans = append(sans, a.Address)
	sans = append(sans, a.SANs...)
	if a.ExternalAddress != "" {
		sans = append(sans, a.ExternalAddress)
	}

	return sans
}

// Validate validates APISpec struct
func (a *K0sAPISpec) Validate() []error {
	var errors []error

	for _, a := range a.Sans() {
		if govalidator.IsIP(a) {
			continue
		}
		if govalidator.IsDNSName(a) {
			continue
		}
		errors = append(errors, fmt.Errorf("%s is not a valid address for sans", a))
	}

	if !govalidator.IsIP(a.Address) {
		errors = append(errors, fmt.Errorf("spec.api.address: %q is not IP address", a.Address))
	}

	return errors
}
