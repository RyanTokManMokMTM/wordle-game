package client

type IClient interface {
	// Run start a client and connect to server
	Run()

	// Close disconnect the connection
	Close()
}
