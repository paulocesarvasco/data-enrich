package constants

type Region int

const (
	Africa Region = iota
	Americas
	Asia
	Europa
	Oceania
)

// Converts region name to a string in lowercase
func (r Region) String() string {
	switch r {
	case Africa:
		return "africa"
	case Americas:
		return "americas"
	case Asia:
		return "asia"
	case Europa:
		return "europa"
	case Oceania:
		return "oceania"
	default:
		return ""
	}
}

const NUM_RECORDS int = 10
