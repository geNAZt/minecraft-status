package data

type MCVersion struct {
	Name     string  `json:"name"`
	Protocol float64 `json:"protocol"`
}

type MCPlayers struct {
	Max    float64 `json:"max"`
	Online float64 `json:"online"`
}

type Favicon struct {
	Icon        string
	DisplayTime int32
}

type Status struct {
	Version     *MCVersion `json:"version"`
	Players     *MCPlayers `json:"players"`
	Favicon     string     `json:"favicon"`
	Favicons    []Favicon
}
