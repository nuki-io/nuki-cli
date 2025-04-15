package bleflows

import "go.nuki.io/nuki/nukictl/pkg/nukible"

type Flow struct {
	ble *nukible.NukiBle
}

func NewFlow(ble *nukible.NukiBle) *Flow {
	return &Flow{
		ble: ble,
	}
}
