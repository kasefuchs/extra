package tunnel

type Config struct {
	MTU    int    `koanf:"mtu"`
	Listen string `koanf:"listen"`
	Remote string `koanf:"remote"`
}
