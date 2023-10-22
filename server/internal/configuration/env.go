package configuration

import "log/slog"

type envModel struct {
	LogLevel slog.Level `env:"LOG_LEVEL" envDefault:"info"`

	RestAddress string `env:"REST_ADDRESS,notEmpty" envDefault:"0.0.0.0:8080"`

	PostgresHost     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	PostgresPort     string `env:"POSTGRES_PORT" envDefault:"5432"`
	PostgresDatabase string `env:"POSTGRES_DATABASE,notEmpty" envDefault:"mincer"`
	PostgresUsername string `env:"POSTGRES_USERNAME,notEmpty" envDefault:"mincer"`
	PostgresPassword string `env:"POSTGRES_PASSWORD_FILE,file,required"`

	NCPort             int        `env:"NC_PORT,notEmpty" envDefault:"12345"`
	NCPrivateKey       PrivateKey `env:"NC_PRIVATE_KEY_FILE,file,required"`
	NCMaxClients       int        `env:"NC_MAX_CLIENTS,notEmpty" envDefault:"256"`
	NCRequestPerSecond int        `env:"NC_REQUEST_PER_SECOND,notEmpty" envDefault:"60"`
}

func (c Config) LogLevel() slog.Level { return c.env.LogLevel }

func (c Config) RestAddress() string { return c.env.RestAddress }

func (c Config) PostgresHost() string     { return c.env.PostgresHost }
func (c Config) PostgresPort() string     { return c.env.PostgresPort }
func (c Config) PostgresDatabase() string { return c.env.PostgresDatabase }
func (c Config) PostgresUsername() string { return c.env.PostgresUsername }
func (c Config) PostgresPassword() string { return c.env.PostgresPassword }

func (c Config) NCPort() int             { return c.env.NCPort }
func (c Config) NCPrivateKey() []byte    { return c.env.NCPrivateKey }
func (c Config) NCMaxClients() int       { return c.env.NCMaxClients }
func (c Config) NCRequestPerSecond() int { return c.env.NCRequestPerSecond }
