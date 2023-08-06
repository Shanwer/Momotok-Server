package system

type configModel struct {
	Server *serverModel `yaml:"server"`
}

type serverModel struct {
	Mode string `yaml:"mode"` // run mode

	Host string `yaml:"host"` // server host

	Port string `yaml:"port"` // server port

	EnableHttps bool `yaml:"enable_https"` // enable https

	TokenExpireSecond int `yaml:"token_expire_second"` // token expire second

	StaticFileUrl string `yaml:"static_file_url"` // system static file path

	DefaultMaxPerPage int `yaml:"default_max_per_page"` //每页加载默认数目

	DatabaseAddress string `yaml:"database_address"` //数据库链接语句

	DriverName string `yaml:"driver_name"` //DriverName

	SecretKey string `yaml:"secret_key"` //token secretkey
}

//two structs above get config info from config.yaml
