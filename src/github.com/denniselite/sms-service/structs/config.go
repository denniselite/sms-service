package structs

type Config struct {
	Project string       `yaml:"project"`
	Listen  int          `yaml:"listen"`
	Rabbit  RabbitConfig `yaml:"rabbit"`
	Sms                    []SmsProvider `yaml:"sms"`
}

type RabbitConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type SmsProvider struct {
	Codes     []string          `yml:"codes"`
	Enable    bool              `yml:"enable"`
	Settings  map[string]string `yml:"settings"`
	Transport string            `yml:"transport"`
}