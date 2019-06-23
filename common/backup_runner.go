package common

type BackupRunner struct {
	app *App
	config *BackupConfig
}

func NewBackupRunner(app *App, config *BackupConfig) *BackupRunner {
	return &BackupRunner{
		app: app,
		config: config,
	}
}

func (br *BackupRunner) Run() error {

}


