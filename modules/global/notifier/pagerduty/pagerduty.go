package pagerduty

import (
	"context"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/copybird/copybird/core"
)

const GROUP_NAME = "global"
const TYPE_NAME = "notifier"
const MODULE_NAME = "pagerduty"

type GlobalNotifierPagerDuty struct {
	core.Module
	config *Config
	client *pagerduty.Client
}

func (m *GlobalNotifierPagerDuty) GetGroup() core.ModuleGroup {
	return GROUP_NAME
}

func (m *GlobalNotifierPagerDuty) GetType() core.ModuleType {
	return TYPE_NAME
}

func (m *GlobalNotifierPagerDuty) GetName() string {
	return MODULE_NAME
}

func (m *GlobalNotifierPagerDuty) GetConfig() interface{} {
	return &Config{}
}

func (m *GlobalNotifierPagerDuty) InitModule(_conf interface{}) error {
	conf := _conf.(Config)
	m.config = &conf
	m.client = pagerduty.NewClient(m.config.AuthToken)
	return nil
}

func (m *GlobalNotifierPagerDuty) Run(ctx context.Context) error {
	_, err := m.client.CreateIncident(m.config.From, &pagerduty.CreateIncident{Incident: pagerduty.CreateIncidentOptions{
		Type:  "dump creation status",
		Title: "Test",
		Service: pagerduty.APIReference{
			ID:   "P4B73MT",
			Type: "service_reference",
		},
		Body: pagerduty.APIDetails{
			Type:    "dump creation failed",
			Details: "Error message that goes along with fail",
		},
	}})
	return err
}

func (m *GlobalNotifierPagerDuty) Close() error {
	return nil
}
