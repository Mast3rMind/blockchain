package blockchain

type Config struct {
	BindAddr      string
	AdvertiseAddr string
	DataDir       string
	SeedNodes     []string
}

func DefaultConfig() *Config {
	c := &Config{
		BindAddr: "127.0.0.1:9119",
		DataDir:  "blockchain-data",
	}
	c.AdvertiseAddr = c.BindAddr
	return c
}

func WriteConfiguration(dataDir string, keypair *Keypair) error {
	// TODO
	return nil
}
func OpenConfiguration(dataDir string) (*Keypair, error) {
	// TODO
	return nil, nil
}
