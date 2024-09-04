package client

type IClient interface {
	Run()
	Close()

	SendRequest(pkgType string, data []byte)

	SetUserId(id string)
	SetUserName(name string)

	GetUserId() string
	GetUserName() string
}
