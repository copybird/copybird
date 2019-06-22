package pagerduty

import (
	pagerduty "github.com/PagerDuty/go-pagerduty"
)

const MODULE_NAME = "pagerduty"

type PagerDuty struct {
	Config *Config
	client *pagerduty.Client
}

func (pd *PagerDuty) GetName() string {
	return MODULE_NAME
}

func (pd *PagerDuty) GetConfig() interface{} {
	return Config{}
}

func (pd *PagerDuty) InitModule(_conf interface{}) error {
	conf := _conf.(Config)
	pd.Config = &conf
	pd.client = pagerduty.NewClient(conf.AuthToken)
	return nil
}

func (pd *PagerDuty) Run() error {
	_, err := pd.client.CreateIncident(pd.Config.From, &pagerduty.CreateIncident{Incident: pagerduty.CreateIncidentOptions{
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

func (pd *PagerDuty) Close() error {
	return nil
}
