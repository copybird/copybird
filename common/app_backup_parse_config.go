package common

//func (a *App) backupParseConfig(configPath string) (*ConfigBackup, error) {
//	configFile, err := os.Open(configPath)
//	if err != nil {
//		return nil, err
//	}
//	var cfg ConfigBackup
//	err = yaml.NewDecoder(configFile).Decode(&cfg)
//	if err != nil {
//		return nil, err
//	}
//	if cfg.Input == nil {
//		return nil, fmt.Errorf("need input module")
//	}
//	moduleInput := a.modulesBackup[cfg.Input.Type]
//	if moduleInput == nil {
//		return nil, fmt.Errorf("backup input module %s not found", cfg.Input.Type)
//	}
//	return &cfg, nil
//}
