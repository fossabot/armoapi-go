package armotypes

import (
	"fmt"
	"time"

	"golang.org/x/exp/slices"
)

type WeeklyReport struct {
	ClustersScannedThisWeek             int                      `json:"clustersScannedThisWeek" bson:"clustersScannedThisWeek"`
	ClustersScannedPrevWeek             int                      `json:"clustersScannedPrevWeek" bson:"clustersScannedPrevWeek"`
	LinkToConfigurationScanningFiltered string                   `json:"linkToConfigurationScanningFiltered" bson:"linkToConfigurationScanningFiltered"`
	RepositoriesScannedThisWeek         int                      `json:"repositoriesScannedThisWeek" bson:"repositoriesScannedThisWeek"`
	RepositoriesScannedPrevWeek         int                      `json:"repositoriesScannedPrevWeek" bson:"repositoriesScannedPrevWeek"`
	LinkToRepositoriesScanningFiltered  string                   `json:"linkToRepositoriesScanningFiltered" bson:"linkToRepositoriesScanningFiltered"`
	RegistriesScannedThisWeek           int                      `json:"registriesScannedThisWeek" bson:"registriesScannedThisWeek"`
	RegistriesScannedPrevWeek           int                      `json:"registriesScannedPrevWeek" bson:"registriesScannedPrevWeek"`
	LinkToRegistriesScanningFiltered    string                   `json:"linkToRegistriesScanningFiltered" bson:"linkToRegistriesScanningFiltered"`
	Top5FailedControls                  []TopCtrlItem            `json:"top5FailedControls" bson:"top5FailedControls"`
	Top5FailedCVEs                      []TopVulItem             `json:"top5FailedCVEs" bson:"top5FailedCVEs"`
	ClustersScanned                     []ClusterResourceScanned `json:"clustersScanned" bson:"clustersScanned"`
	RepositoriesScanned                 []RepositoryScanned      `json:"repositoriesScanned" bson:"repositoriesScanned"`
	RegistriesScanned                   []RegistryScanned        `json:"registriesScanned" bson:"registriesScanned"`
}
type PushNotification struct {
	Misconfigurations Misconfigurations
	NewClusterAdmins  NewClusterAdmins
}

type NewClusterAdmins []NewClusterAdmin
type NewClusterAdmin struct {
	Resource    string
	Link        string
	ClusterName string
}

type Misconfigurations []Misconfiguration
type Misconfiguration struct {
	Name                      string
	Type                      ScanType
	Link                      string
	PercentageIncrease        uint64
	FrameworksComplianceDrift map[string]int
}
type ScanType string

const (
	ScanTypePosture      ScanType = "posture"
	ScanTypeRepositories ScanType = "repository"
)

type NotificationsConfig struct {
	//Map of unsubscribed user id to notification config identifier
	UnsubscribedUsers  map[string][]NotificationConfigIdentifier `json:"unsubscribedUsers,omitempty" bson:"unsubscribedUsers,omitempty"`
	LatestWeeklyReport *WeeklyReport                             `json:"latestWeeklyReport,omitempty" bson:"latestWeeklyReport,omitempty"`
	LatestPushReports  map[string]*PushReport                    `json:"latestPushReports,omitempty" bson:"latestPushReports,omitempty"`
	AlertChannels      map[CollaborationType][]AlertChannel      `json:"alertChannels,omitempty" bson:"alertChannels,omitempty"`
}

type AlertChannel struct {
	ChannelType             CollaborationType `json:"channelType,omitempty" bson:"channelType,omitempty"`
	CollaborationConfigGUID string            `json:"collaborationConfigId,omitempty" bson:"collaborationConfigId,omitempty"`
	Alerts                  []AlertConfig     `json:"notifications,omitempty" bson:"notifications,omitempty"`
}

func (ac *AlertChannel) GetAlertConfig(notificationType NotificationType) *AlertConfig {
	for _, alert := range ac.Alerts {
		if alert.NotificationType == notificationType {
			return &alert
		}
	}
	return nil
}

type AlertConfig struct {
	NotificationConfigIdentifier `json:",inline" bson:",inline"`
	Scope                        []AlertScope           `json:"scope,omitempty" bson:"scope,omitempty"`
	Parameters                   map[string]interface{} `json:"attributes,omitempty" bson:"attributes,omitempty"`
	Disabled                     *bool                  `json:"disabled,omitempty" bson:"disabled,omitempty"`
}

type AlertScope struct {
	Cluster    string   `json:"cluster,omitempty" bson:"cluster,omitempty"`
	Namespaces []string `json:"namespaces,omitempty" bson:"namespaces,omitempty"`
}

func (nc *NotificationsConfig) GetAlertConfigurations(notificationType NotificationType) []AlertConfig {
	alerts := make([]AlertConfig, 0)
	for _, typesChannels := range nc.AlertChannels {
		for _, alertChannel := range typesChannels {
			if config := alertChannel.GetAlertConfig(notificationType); config != nil {
				alerts = append(alerts, *config)
			}
		}
	}
	return alerts
}

func (nc *NotificationsConfig) AddLatestPushReport(report *PushReport) {
	if report == nil {
		return
	}
	if nc.LatestPushReports == nil {
		nc.LatestPushReports = make(map[string]*PushReport, 0)
	}
	nc.LatestPushReports[fmt.Sprintf("%s_%s", report.Cluster, report.ScanType)] = report
}

func (nc *NotificationsConfig) GetLatestPushReport(cluster string, scanType ScanType) *PushReport {
	if val, ok := nc.LatestPushReports[fmt.Sprintf("%s_%s", cluster, scanType)]; ok {
		return val
	}
	return nil
}

type PushReport struct {
	Cluster                   string             `json:"custer,omitempty" bson:"custer,omitempty"`
	ReportGUID                string             `json:"reportGUID,omitempty" bson:"reportGUID,omitempty"`
	ScanType                  ScanType           `json:"scanType" bson:"scanType"`
	Timestamp                 time.Time          `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
	FailedResources           uint64             `json:"failedResources,omitempty" bson:"failedResources,omitempty"`
	FrameworksComplianceScore map[string]float32 `json:"frameworksComplianceScore,omitempty" bson:"frameworksComplianceScore,omitempty"`
}

type NotificationConfigIdentifier struct {
	NotificationType NotificationType `json:"notificationType,omitempty" bson:"notificationType,omitempty"`
}

func (nci *NotificationConfigIdentifier) Validate() error {
	if slices.Contains(notificationTypes, nci.NotificationType) {
		return nil
	}
	if nci.NotificationType == "" {
		return fmt.Errorf("notification type is required")
	}
	return fmt.Errorf("invalid notification type: %s", nci.NotificationType)
}

type NotificationType string

const (
	NotificationTypeAll                 NotificationType = "all"
	NotificationTypePushPosture         NotificationType = "push"
	NotificationTypeWeekly              NotificationType = "weekly"
	NotificationTypeComplianceDrift     NotificationType = "complianceDrift"
	NotificationTypeNewClusterAdmin     NotificationType = "newClusterAdmin"
	NotificationTypeNewVulnerability    NotificationType = "newVulnerability"
	NotificationTypeVulnerabilityNewFix NotificationType = "vulnerabilityNewFix"
)

var notificationTypes = []NotificationType{NotificationTypeAll,
	NotificationTypePushPosture,
	NotificationTypeWeekly,
	NotificationTypeComplianceDrift,
	NotificationTypeNewClusterAdmin,
	NotificationTypeNewVulnerability,
	NotificationTypeVulnerabilityNewFix,
}

type RegistryScanned struct {
	Registry ResourceScanned `json:"registry" bson:"registry"`
}

type RepositoryScanned struct {
	ReportGUID string          `json:"reportGUID" bson:"reportGUID"`
	Repository ResourceScanned `json:"repository" bson:"repository"`
}

type ClusterResourceScanned struct {
	ShortName       string          `json:"shortName" bson:"shortName"`
	Cluster         ResourceScanned `json:"cluster" bson:"cluster"`
	ReportGUID      string          `json:"reportGUID" bson:"reportGUID"`
	FailedResources uint64          `json:"failedResources" bson:"failedResources"`
}

type ResourceScanned struct {
	Kind                         string                     `json:"kind" bson:"kind"`
	Name                         string                     `json:"name" bson:"name"`
	MapSeverityToSeverityDetails map[string]SeverityDetails `json:"mapSeverityToSeverityDetails" bson:"mapSeverityToSeverityDetails"`
}

type SeverityDetails struct {
	Severity              string `json:"severity" bson:"severity"`
	FailedResourcesNumber int    `json:"failedResourcesNumber" bson:"failedResourcesNumber"`
}

type NotificationPushEvent struct {
	CustomerGUID string           `json:"customerGUID"`
	EventName    string           `json:"eventName"`
	EventTime    time.Time        `json:"eventTime"`
	ReportGUID   string           `json:"reportGUID,omitempty"`
	ClusterName  string           `json:"clusterName,omitempty"`
	Designators  PortalDesignator `json:"designators,omitempty"`
}
