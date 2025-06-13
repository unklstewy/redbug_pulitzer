package reporting

// CommandAPI represents a documented API command (duplicated from protocol package to avoid import cycle)
type CommandAPI struct {
	Command        string // Command type identifier
	HexValue       string // Raw hex representation
	ASCIIValue     string // ASCII representation
	Description    string // Human-readable description
	ResponseType   string // Type of response
	ResponseHex    string // Raw hex of response
	ResponseASCII  string // ASCII of response
	FrequencyCount int    // Number of times observed
	TimingAverage  string // Average response time
	DataCategory   string // For write operations: data category
	SuccessRate    string // Percentage of successful responses
}
