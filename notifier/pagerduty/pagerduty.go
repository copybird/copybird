package pagerduty

const MODULE_NAME = "pagerduty"

type PagerDuty struct {
	Config *Config
}

func (pd *PagerDuty) GetName() string {
	return MODULE_NAME
}

func (pd *PagerDuty) GetConfig() interface{} {
	return Config{}
}

func (pd *PagerDuty) InitModule(_conf interface{}) error {
	//conf := _conf.(Config)
	return nil
}

func (pd *PagerDuty) Run() error {

	return nil
}

func (pd *PagerDuty) Close() error {
	return nil
}
