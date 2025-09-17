package apiv2

import "net/netip"

type DeviceID int

type CustomerID int

type Device struct {
	DeviceID   DeviceID
	Name       string
	StdName    string
	CustomerID CustomerID
	PrimaryIP  netip.Addr
	// TODO Status           enum
}
