package pagerduty

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/copybird/copybird/core"
)

const MODULE_NAME = "pagerduty"

type PagerDuty struct {
	core.Module
	config *Config
	client *pagerduty.Client
}

func (m *PagerDuty) GetName() string {
	return MODULE_NAME
}

func (m *PagerDuty) GetConfig() interface{} {
	return &Config{}
}

func (m *PagerDuty) InitModule(_conf interface{}) error {
	conf := _conf.(Config)
	m.config = &conf
	m.client = pagerduty.NewClient(m.config.AuthToken)
	return nil
}

func (m *PagerDuty) Run() error {
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

func (m *PagerDuty) Close() error {
	return nil
}
