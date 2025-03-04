package config

type ConfigPinningService struct {
	Uploader           string
	PinningService     string
	BlockserviceApiKey string
	DedicatedGateway   bool
	RedisConn         string
	AmqpConnect        string
}
