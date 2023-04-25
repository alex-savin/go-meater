package meater

// Cook .
type Cook struct {
	ID          string      `json:"id"`             // 64 chars (sometimes 32)
	Name        string      `json:"name,omitempty"` // Predefined names
	State       string      `json:"state"`          // [ Not Started | Configured | Started | Ready For Resting | Resting | Slightly Underdone | Finished | Slightly Overdone | OVERCOOK! ]
	Temperature TempCook    `json:"temperature"`
	CookingTime CookingTime `json:"time"`
}

// TempCook .
type TempCook struct {
	Target float64 `json:"target,omitempty"` // Target temperature
	Peak   float64 `json:"peak,omitempty"`   // Peak temperature reached during cook
}

// CookingTime .
type CookingTime struct {
	Elapsed   int `json:"elapsed"`   // Time since the start of cook in seconds. Default: 0
	Remaining int `json:"remaining"` // Remaining time in seconds. When unknown/calculating default is used. Default: -1
}

// GetStatus .
func (c *Cook) GetStatus() string {
	return c.State
}

// GetName .
func (c *Cook) GetName() string {
	return c.Name
}

// GetTemps .
func (c *Cook) GetTemps() (float64, float64) {
	return c.Temperature.Target, c.Temperature.Peak
}

// GetTime .
func (c *Cook) GetTime() (int, int) {
	return c.CookingTime.Elapsed, c.CookingTime.Remaining
}
