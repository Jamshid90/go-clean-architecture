package config

import (
	"github.com/BurntSushi/toml"
	"strings"
	"time"
)

type Config struct {
	Server struct {
		Protocol   string `toml:"protocol"`
		Host       string `toml:"host"`
		Port       string `toml:"port"`
		SSLCert    string `toml:"sslcert"`
		SSLPrivKey string `toml:"sslprivkey"`
	} `toml:"server"`
	Database struct {
		Host     string `toml:"host"`
		Port     string `toml:"port"`
		DBName   string `toml:"dbname"`
		User     string `toml:"user"`
		Password string `toml:"password"`
		Sslmode  string `toml:"sslmode"`
	} `toml:"database"`
	Context struct {
		Timeout     time.Duration `toml:"timeout"`
	} `toml:"context"`

	Jwt struct {
		Secret string `toml:"secret"`
		TTL    string `toml:"ttl"`
	} `toml:"jwt"`
}

func NewConfig(filePath string) (*Config, error) {
	var config *Config
	_, err := toml.DecodeFile(filePath, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func (c *Config) GetPsqlConnStr() string  {

	var conn []string

	if len(c.Database.Host) != 0 {
		conn = append(conn, "host="+c.Database.Host)
	}

	if len(c.Database.Port) != 0 {
		conn = append(conn, "port="+c.Database.Port)
	}

	if len(c.Database.User) != 0 {
		conn = append(conn, "user="+c.Database.User)
	}

	if len(c.Database.Password) != 0 {
		conn = append(conn, "password="+c.Database.Password)
	}

	if len(c.Database.DBName) != 0 {
		conn = append(conn, "dbname="+c.Database.DBName)
	}

	if len(c.Database.Sslmode) != 0 {
		conn = append(conn, "sslmode="+c.Database.Sslmode)
	}

	return strings.Join(conn, " ")
}
