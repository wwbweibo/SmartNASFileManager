package dl

type Config struct {
	Scheme  string `json:"scheme" yaml:"scheme"`
	Host    string `json:"host" yaml:"host"`
	Port    int    `json:"port" yaml:"port"`
	TimeOut int    `json:"timeout" yaml:"timeout"`
}
