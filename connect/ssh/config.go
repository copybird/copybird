package ssh

type Config struct {
	LocalEndpointHost string
	LocalEndpointPort int

	ServerEndpointHost string
	ServerEndpointPort int

	RemoteEndpointHost string
	RemoteEndpointPort int

	RemoteUser string

	// Full path to id_rsa key
	KeyPath string
}
