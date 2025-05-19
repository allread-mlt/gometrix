package gometrix

type StatsdMetricsData struct {
	ServerHost string `yaml:"host"`
	ServerPort int64  `yaml:"port"`
	Prefix     string `yaml:"prefix"`
}

type LoggingMetricsData struct {
	Timeout     int    `yaml:"timeout"`
	LogFilePath string `yaml:"log_file_path"`
	MaxFiles    int    `yaml:"max_files"`
	MaxFileSize int    `yaml:"max_file_size"`
}

