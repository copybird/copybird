package nats

import "time"

type Config struct {
	ClientID           string       
	NATSURL            string       
	ClusterID          string 
	ConnectTimeout     time.Duration
}

