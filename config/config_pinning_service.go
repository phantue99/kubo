package config

type ConfigPinningService struct {
	Uploader             string
	PinningService       string
	BlockserviceApiKey   string
	DedicatedGateway     bool
	RedisConn            string
	AmqpConnect          string
	BlockEncryptionKey   string
	EncryptedBlockPrefix string
	IpfsDomain           string
	SslCertPath          string
	SslKeyPath           string
}
