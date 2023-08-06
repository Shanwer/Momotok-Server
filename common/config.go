package common

type configModel struct {
	Server *serverModel `yaml:"server"`
}

//serverModel get server information from config.yml

type serverModel struct {
	Mode string `yaml:"mode"` // run mode

	Host string `yaml:"host"` // server host

	Port string `yaml:"port"` // server port

	EnableHttps bool `yaml:"enable_https"` // enable https

	TokenExpireSecond int `yaml:"token_expire_second"` // token expire second

	SystemStaticFilePath string `yaml:"system_static_file_path"` // system static file path

	VideoUrlPath string `yaml:"video_url_path"` // 视频url

	MaxPerPage int `yaml:"max_per_page"` //每页加载最大数目

	DefaultMaxPerPage int `yaml:"default_max_per_page"` //每页加载默认数目

	DatabaseAddress string `yaml:"database_address"` //数据库链接语句

	DriverName string `yaml:"driver_name"` //DriverName

	SecretKey string `yaml:"secret_key"` //token secretkey
}
