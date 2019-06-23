package scp

type Config struct {
	Addr               string
	Port               int
	User               string
	Password           string
	FileName           string
	PathToKey          string
	PrivateKeyPassword string
}
