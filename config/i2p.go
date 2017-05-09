package config

// i2p connectivity config
type I2PConfig struct {
	Name    string `json: "name"`
	Addr    string `json: "samaddr"`
	Keyfile string `json: "keyfile"`
}
