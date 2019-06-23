package common

type RunnerBackup struct {
	app    *App
	config *ConfigBackup
}

func NewRunnerBackup(app *App, config *ConfigBackup) *RunnerBackup {
	return &RunnerBackup{
		app:    app,
		config: config,
	}
}

func (rb *RunnerBackup) Run() error {
	return nil
}
