package notify

// Snapshot English text notification_events.snapshot English text JSONB English text。English text db English text
// English text，English text server English text。English text
// 「English text」English text：db English text，English text。
//
// English text：English text、English text、English text，
// English text**English text**English text——English text「English text low」English text
// English text。English text fan-out English text JOIN findings/tasks/assets English text。
type Snapshot struct {
	// English text：finding_created / finding_status_changed
	Kind      string  `json:"kind"`
	FindingID int64   `json:"finding_id"`
	TaskID    int64   `json:"task_id"`
	VulnClass string  `json:"vulnclass"`
	Name      string  `json:"name"`
	Severity  string  `json:"severity"`
	Summary   string  `json:"summary"`
	AssetIDs  []int64 `json:"asset_ids"`
	// English text kind=finding_status_changed English text。
	FromStatus string `json:"from_status,omitempty"`
	ToStatus   string `json:"to_status,omitempty"`
}

// Item English text，English text。
type Item struct {
	FindingID int64
	Name      string
	VulnClass string
	Severity  string
	Summary   string
	// Assets English text（English text/IP）。English text server English text——
	// English text，English text。
	Assets []string
	// DetailURL English text；English text public_base_url，English text。
	DetailURL string
	// English text；English text「English text → English text」。
	FromStatus string
	ToStatus   string
}

// IsStatusChange English text。
func (i Item) IsStatusChange() bool { return i.FromStatus != "" || i.ToStatus != "" }

// Title English text：English text name，English text vulnclass，
// English text——English text。
func (i Item) Title() string {
	if i.Name != "" {
		return i.Name
	}
	if i.VulnClass != "" {
		return i.VulnClass
	}
	return "(English text)"
}

// Message English text。
type Message struct {
	// English text 1；English text（digest）English text。
	// English text，English text。
	Items []Item
	// Batch=true English text（English text、English text）。
	Batch bool
	// WindowMinutes English text（English text），English text Batch=true English text「English text N English text」。
	// English text time.Since：English text，English text。
	WindowMinutes int
	// HomeURL English text（English text public_base_url）；English text。
	HomeURL string
}
