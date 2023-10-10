package config

type ConfigPinningSerice struct {
	Uploader           string
	PinningService     string
	BlockserviceApiKey string
	DedicatedGateway   bool
	RedisConns		 []string
}
