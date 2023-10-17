package configuration

type envModel struct {
	RestAddress string `env:"REST_ADDRESS,notEmpty" envDefault:"0.0.0.0:8080"`

	PostgresHost     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	PostgresPort     string `env:"POSTGRES_PORT" envDefault:"5432"`
	PostgresDatabase string `env:"POSTGRES_DATABASE,notEmpty" envDefault:"mincer"`
	PostgresUsername string `env:"POSTGRES_USERNAME,notEmpty" envDefault:"mincer"`
	PostgresPassword string `env:"POSTGRES_PASSWORD_FILE,file,required"`

	NCAddress          string     `env:"NC_ADDRESS,notEmpty" envDefault:"192.168.0.17:12345"` // TODO default
	NCPrivateKey       PrivateKey `env:"NC_PRIVATE_KEY_FILE,file,required"`
	NCMaxClients       int        `env:"NC_MAX_CLIENTS,notEmpty" envDefault:"256"`
	NCRequestPerSecond int        `env:"NC_REQUEST_PER_SECOND,notEmpty" envDefault:"60"`
}

func (c Config) RestAddress() string { return c.env.RestAddress }

func (c Config) PostgresHost() string     { return c.env.PostgresHost }
func (c Config) PostgresPort() string     { return c.env.PostgresPort }
func (c Config) PostgresDatabase() string { return c.env.PostgresDatabase }
func (c Config) PostgresUsername() string { return c.env.PostgresUsername }
func (c Config) PostgresPassword() string { return c.env.PostgresPassword }

func (c Config) NCAddress() string       { return c.env.NCAddress }
func (c Config) NCPrivateKey() []byte    { return c.env.NCPrivateKey }
func (c Config) NCMaxClients() int       { return c.env.NCMaxClients }
func (c Config) NCRequestPerSecond() int { return c.env.NCRequestPerSecond }
