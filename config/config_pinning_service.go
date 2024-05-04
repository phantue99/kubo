package config

type ConfigPinningService struct {
	Uploader           string
	PinningService     string
	BlockserviceApiKey string
	DedicatedGateway   bool
	RedisConns         []string
	AmqpConnect        string
}
