package meater

import "errors"

// Probe .
type Probe struct {
	ID          string    `json:"id"`             // 64 chars
	Temperature TempProbe `json:"temperature"`    // Internal & Ambient
	Cook        *Cook     `json:"cook,omitempty"` // Could be null
	UpdatedAt   int       `json:"updated_at"`     // Time data was last updated at as a UNIX timestamp
	client      *Client
}

// TempDevice .
type TempProbe struct {
	Internal float64 `json:"internal"` // Internal temperature
	Ambient  float64 `json:"ambient"`  // Ambient temperature. If ambient is less than internal, ambient will equal internal
}

// GetCooking .
func (p *Probe) GetCooking() (*Cook, error) {
	p.client.GetProbeByID(p.ID)
	if !p.HasCooking() {
		return nil, errors.New("no active cook")
	}
	return p.Cook, nil
}

// HasCooking .
func (p *Probe) HasCooking() bool {
	p.client.GetProbeByID(p.ID)
	return !isNil(p.Cook)
}

// GetReadings .
func (p *Probe) GetReadings() (float64, float64) {
	p.client.GetProbeByID(p.ID)
	return p.Temperature.Internal, p.Temperature.Ambient
}
