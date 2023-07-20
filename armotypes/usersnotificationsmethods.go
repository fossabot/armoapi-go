package armotypes

import (
	"fmt"

	"golang.org/x/exp/slices"
)

func (ap *NotificationParams) SetDriftPercentage(percentage int) {
	ap.DriftPercentage = &percentage
}

func (ap *NotificationParams) SetMinSeverity(severity int) {
	ap.MinSeverity = &severity
}

func (ac *AlertChannel) GetAlertConfig(notificationType NotificationType) *AlertConfig {
	for _, alert := range ac.Alerts {
		if alert.NotificationType == notificationType {
			return &alert
		}
	}
	return nil
}

func (ac *AlertChannel) AddAlertConfig(config AlertConfig) error {
	if config.NotificationType == "" {
		return fmt.Errorf("notification type is required")
	}
	for _, alert := range ac.Alerts {
		if alert.NotificationType == config.NotificationType {
			return fmt.Errorf("alert config for notification type %s already exists", config.NotificationType)
		}
	}
	ac.Alerts = append(ac.Alerts, config)
	return nil
}

func (ac *AlertChannel) IsInScope(cluster, namespace string) bool {
	if ac.Scope == nil {
		return true
	}
	return ac.Scope.IsInScope(cluster, namespace)
}

func (ac *AlertConfig) IsEnabled() bool {
	if ac.Disabled == nil {
		return true
	}
	return !*ac.Disabled
}

func (ac *AlertScope) IsInScope(cluster, namespace string) bool {
	if ac.Cluster == "" {
		//no scope defined
		return true
	}
	if ac.Cluster != cluster {
		return false
	}
	if namespace == "" {
		return true
	}
	if len(ac.Namespaces) == 0 {
		return true
	}
	return slices.Contains(ac.Namespaces, namespace)
}

func (nc *NotificationsConfig) IsInScope(cluster, namespace string) bool {
	for _, typesChannels := range nc.AlertChannels {
		for _, alertChannel := range typesChannels {
			if alertChannel.IsInScope(cluster, namespace) {
				return true
			}
		}
	}
	return false
}

func (nc *NotificationsConfig) GetProviderChannels(provider ChannelProvider) []AlertChannel {
	if nc.AlertChannels == nil {
		return nil
	}
	return nc.AlertChannels[provider]
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

func (nci *NotificationConfigIdentifier) Validate() error {
	if slices.Contains(notificationTypes, nci.NotificationType) {
		return nil
	}
	if nci.NotificationType == "" {
		return fmt.Errorf("notification type is required")
	}
	return fmt.Errorf("invalid notification type: %s", nci.NotificationType)
}
