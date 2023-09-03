package postgresmodels

import (
	"encoding/json"
	"time"

	"github.com/armosec/armoapi-go/identifiers"
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

// TODO: add explicit column names, add validation

const (
	SummaryStatusSuccess = "Success"
	SummaryStatusPending = "Pending"
)

type Vulnerability struct {
	BaseModel
	Name          string `gorm:"primaryKey"`
	Severity      string
	SeverityScore int
	IsRCE         bool
	Links         pq.StringArray `gorm:"type:text[]"`
	Description   string
}

type VulnerabilityFinding struct {
	BaseModel
	VulnerabilityName string        `gorm:"primaryKey"`
	Vulnerability     Vulnerability `gorm:"foreignKey:VulnerabilityName"`
	ImageScanId       string        `gorm:"primaryKey"`
	Component         string        `gorm:"primaryKey"`
	ComponentVersion  string        `gorm:"primaryKey"`
	LayerHash         string        `gorm:"primaryKey"`
	FixAvailable      *bool
	FixedInVersion    string
	LayerIndex        *int
	LayerCommand      string
	IsRelevant        *bool
	RelevantLabel     string
	IsIgnored         *bool
	IgnoreRuleIds     pq.StringArray `gorm:"type:text[]"`
}

type VulnerabilityScanSummary struct {
	BaseModel
	ScanKind                   string
	ImageScanId                string `gorm:"primaryKey"`
	ContainerSpecId            string
	Timestamp                  time.Time
	CustomerGuid               string
	Wlid                       string
	Designators                datatypes.JSON
	ImageRegistry              string
	ImageRepository            string
	ImageTag                   string
	ImageHash                  string
	JobIds                     pq.StringArray `gorm:"type:text[]"`
	Status                     string
	Errors                     pq.StringArray               `gorm:"type:text[]"`
	Findings                   []VulnerabilityFinding       `gorm:"foreignKey:ImageScanId"`
	VulnerabilitySeverityStats []VulnerabilitySeverityStats `gorm:"foreignKey:ImageScanId"`
	IsStub                     *bool                        // if true, this is a stub scan summary, and the actual scan summary is not yet available. Should be deleted once we have the real one.
}

// ContextualVulnerabilityFinding is a VulnerabilityFinding with a VulnerabilityScanSummary, do not auto-migrate it
// uses only for retreiving data from db
type ContextualVulnerabilityFinding struct {
	VulnerabilityFinding     `gorm:"embedded"`
	VulnerabilityScanSummary VulnerabilityScanSummary `gorm:"foreignKey:ImageScanId"`
}

func (ContextualVulnerabilityFinding) TableName() string {
	return "vulnerability_findings"
}

func (v VulnerabilityScanSummary) GetDesignators() (*identifiers.PortalDesignator, error) {
	var designators *identifiers.PortalDesignator

	desigs, err := v.Designators.Value()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(desigs.(string)), &designators); err != nil {
		return nil, err
	}

	return designators, nil
}

type VulnerabilitySeverityStats struct {
	BaseModel
	ImageScanId                  string         `gorm:"primaryKey"`
	Severity                     string         `gorm:"primaryKey"`
	DayDate                      datatypes.Date `gorm:"primaryKey"`
	SeverityScore                int
	TotalCount                   int64
	RCEFixCount                  int64
	FixAvailableOfTotalCount     int64
	RelevantCount                int64
	FixAvailableForRelevantCount int64
	RCECount                     int64
	UrgentCount                  int64
	NeglectedCount               int64
	HealthStatus                 string
}
