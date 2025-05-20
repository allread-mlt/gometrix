package gometrix

type StatsdMetricsData struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	Prefix string `yaml:"prefix"`
}

type LoggingMetricsData struct {
	Timeout     int    `yaml:"timeout"`
	LogFilePath string `yaml:"log_file_path"`
	MaxFiles    int    `yaml:"max_files"`
	MaxFileSize int    `yaml:"max_file_size"`
}
