package gometrix

type StatsdMetricsData struct {
	ServerHost string `yaml:"host"`
	ServerPort int64  `yaml:"port"`
	Prefix     string `yaml:"prefix"`
}
