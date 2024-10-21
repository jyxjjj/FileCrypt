package main

import (
	"errors"
	"fmt"
	"github.com/huin/goupnp/dcps/internetgateway2"
)

type RouterClient interface {
	AddPortMapping(
		NewRemoteHost string,
		NewExternalPort uint16,
		NewProtocol string,
		NewInternalPort uint16,
		NewInternalClient string,
		NewEnabled bool,
		NewPortMappingDescription string,
		NewLeaseDuration uint32,
	) (err error)

	GetExternalIPAddress() (
		NewExternalIPAddress string,
		err error,
	)
}

func PickRouterClient() (RouterClient, error) {
	var ip1Clients []*internetgateway2.WANIPConnection1
	var ip2Clients []*internetgateway2.WANIPConnection2
	var ppp1Clients []*internetgateway2.WANPPPConnection1
	var err error
	ip1Clients, _, err = internetgateway2.NewWANIPConnection1Clients()
	if err != nil {
		ip2Clients, _, err = internetgateway2.NewWANIPConnection2Clients()
		if err != nil {
			ppp1Clients, _, err = internetgateway2.NewWANPPPConnection1Clients()
		}
	}

	switch {
	case len(ip2Clients) == 1:
		return ip2Clients[0], nil
	case len(ip1Clients) == 1:
		return ip1Clients[0], nil
	case len(ppp1Clients) == 1:
		return ppp1Clients[0], nil
	default:
		return nil, errors.New("multiple or no services found")
	}
}
func main() {
	client, err := PickRouterClient()
	if err != nil {
		fmt.Println("Failed to find a UPnP enabled router: ", err)
		return
	}

	externalIP, _ := client.GetExternalIPAddress()
	fmt.Println("Our external IP address is: ", externalIP)

	_ = client.AddPortMapping(
		"",
		56789,
		"UDP",
		56789,
		// Internal address on the LAN we want to forward to.
		"192.168.50.3",
		// Enabled:
		true,
		// Informational description for the client requesting the port forwarding.
		"UPnP test forwarding",
		// How long should the port forward last for in seconds.
		// If you want to keep it open for longer and potentially across router
		// resets, you might want to periodically request before this elapses.
		180,
	)
}
