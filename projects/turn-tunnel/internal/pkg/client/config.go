package client

type Config struct {
	Turn     string `koanf:"turn"`
	Stun     string `koanf:"stun"`
	Username string `koanf:"username"`
	Password string `koanf:"password"`
}
