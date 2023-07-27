package armotypes

import "github.com/armosec/armoapi-go/identifiers"

type AttackChainStatus string
type ProcessingStatus string

const (
	StatusActive AttackChainStatus = "active"
	StatusFixed  AttackChainStatus = "fixed"
	// StatusFixedSeen AttackChainStatus = "fixedSeen"

	ProcessingStatusProcessing ProcessingStatus = "processing"
	ProcessingStatusDone       ProcessingStatus = "done"
)

type AttackChain struct {
	AttackChainNodes AttackChainNode `json:"attackChainNodes,omitempty" bson:"attackChainNodes,omitempty"`
	AttackChainConfig
}

type AttackChainConfig struct {
	PortalBase       `json:",inline" bson:",inline"`
	Resource         identifiers.PortalDesignator `json:"resource,omitempty" bson:"resource,omitempty"`
	Description      string                       `json:"description,omitempty" bson:"description,omitempty"`
	CreationTime     string                       `json:"creationTime,omitempty" bson:"creationTime,omitempty"`
	AttackChainID    string                       `json:"attackChainID,omitempty" bson:"attackChainID,omitempty"` // name/cluster/resourceID
	ClusterName      string                       `json:"clusterName,omitempty" bson:"clusterName,omitempty"`
	CustomerGUID     string                       `json:"customerGUID,omitempty" bson:"customerGUID,omitempty"`
	LatestReportGUID string                       `json:"latestReportGUID,omitempty" bson:"latestReportGUID,omitempty"` // latest reportGUID in which this attack chain was identified
	UIStatus         *AttackChainUIStatus         `json:"uiStatus,omitempty" bson:"uiStatus,omitempty"`
}

type AttackChainNode struct {
	Name             string                         `json:"name" bson:"name,omitempty"`
	Description      string                         `json:"description" bson:"description,omitempty"`
	ControlIDs       []string                       `json:"controlIDs,omitempty" bson:"controlIDs,omitempty"` // failed/ignored controls that are associated to this attack chain node
	Vulnerabilities  []Vulnerabilities              `json:"vulnerabilities,omitempty" bson:"vulnerabilities,omitempty"`
	RelatedResources []identifiers.PortalDesignator `json:"relatedResources" bson:"relatedResources,omitempty"`
	NextNodes        []AttackChainNode              `json:"nextNodes,omitempty" bson:"nextNodes,omitempty"`
}

type Vulnerabilities struct {
	ContainersScanID string   `json:"containersScanID" bson:"containersScanID,omitempty"`
	ContainerName    string   `json:"containerName" bson:"containerName,omitempty"`
	Names            []string `json:"names" bson:"names,omitempty"` // CVE names
}

// struct for UI support. All strings are timestamps
type AttackChainUIStatus struct {
	// fields updated by the BE
	FirstSeen string `json:"firstSeen,omitempty" bson:"firstSeen,omitempty"` // timestamp of first scan in which the attack chain was identified
	// fields updated by the UI
	ViewedMainScreen string `json:"wasViewedMainScreen,omitempty" bson:"wasViewedMainScreen,omitempty"` // if the attack chain was viewed by the user// New badge
	ProcessingStatus string `json:"processingStatus,omitempty" bson:"processingStatus,omitempty"`       // "processing"/ "done"
}