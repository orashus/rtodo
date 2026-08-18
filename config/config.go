package config

type config struct {
	FILE_PATH string
}

func (c *config) Configure(filePath string) {
	c.FILE_PATH = filePath
}

var Config config = config{}
