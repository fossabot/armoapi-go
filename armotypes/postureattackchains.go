package armotypes

type AttackChainStatus string
type ProcessingStatus string

const (
	StatusActive AttackChainStatus = "active"
	StatusFixed  AttackChainStatus = "fixed"
	// StatusFixedSeen AttackChainStatus = "fixedSeen"

	ProcessingStatusProcessing ProcessingStatus = "processing"
	ProcessingStatusDone       ProcessingStatus = "done"
)

type AttackChainType struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AttackChain struct {
	Type             *AttackChainType     `json:"type"`
	ClusterName      string               `json:"clusterName"`
	Resource         PortalDesignator     `json:"resource"`
	AttackChainID    string               `json:"attackChainID"` // name/cluster/resourceID
	CustomerGUID     string               `json:"customerGUID"`
	Status           AttackChainStatus    `json:"status"`
	AttackChainNodes AttackChainNode      `json:"attackChainNodes"`
	UIStatus         *AttackChainUIStatus `json:"uiStatus"`
	AttackChainScans *AttackChainScans    `json:"attackChainScans"`
}

type AttackChainNode struct {
	Name             string               `json:"name"`
	Description      string               `json:"description"`
	ControlIDs       []string             `json:"controlIDs,omitempty"` // failed/ignored controls that are associated to this attack chain node
	Vulnerabilities  []VulnerabilityNames `json:"vulnerabilitiesNames,omitempty"`
	RelatedResources []PortalDesignator   `json:"relatedResources"`
	NextNodes        []AttackChainNode    `json:"nextNodes,omitempty"`
}

type VulnerabilityNames struct {
	ContainersScanID string `json:"containersScanID"`
	Name             string `json:"name"` // CVE name
}

type AttackChainScans struct {
	LatestReportGUID string `json:"latestReportGUID"` // list of reportGUIDs in which this attack chain was identified
}

// struct for UI support. All strings are timestamps
type AttackChainUIStatus struct {
	// fields updated by the BE
	FirstSeen string `json:"firstSeen"` // timestamp of first scan in which the attack chain was identified
	// fields updated by the UI
	ViewedMainScreen string `json:"wasViewedMainScreen"` // if the attack chain was viewed by the user// New badge
	ProcessingStatus string `json:"processingStatus"`    // "processing"/"" //
}
