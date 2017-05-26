package config

// I2PConfig : Config struct for I2P
type I2PConfig struct {
	Name    string `json:"name"`
	Addr    string `json:"samaddr"`
	Keyfile string `json:"keyfile"`
}
