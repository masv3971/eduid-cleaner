package model

type APIServer struct {
	Host string `yaml:"host" validate:"required"`
}

type Workers struct {
	SKV   Worker `yaml:"skv"`
	Ladok Worker `yaml:"ladok"`
}
type Worker struct {
	Periodicity     float64 `yaml:"periodicity"`
	SubWorkerAmount int     `yaml:"sub_worker_amount" validate:"required"`
}
type Redis struct {
	DB                  int      `yaml:"db" validate:"required"`
	Addr                string   `yaml:"host" validate:"required_without_all=SentinelHosts SentinelServiceName"`
	SentinelHosts       []string `yaml:"sentinel_hosts" validate:"required_without=Addr,omitempty,min=2,max=4"`
	SentinelServiceName string   `yaml:"sentinel_service_name" validate:"required_with=SentinelHosts"`
}
type Storage struct {
	Redis Redis `yaml:"redis"`
}
type RemoteAPI struct {
	URL string `yaml:"url" validate:"required,url"`
}

type Sunet struct {
	Auth  RemoteAPI `yaml:"auth"`
	AmAPI RemoteAPI `yaml:"am_api"`
}

type Log struct {
	Level string `yaml:"level"`
}

// Cfg is the main configuration structure for this application
type Cfg struct {
	APIServer  APIServer `yaml:"api_server"`
	Production bool      `yaml:"production"`
	HTTPProxy  string    `yaml:"http_proxy"`
	Workers    Workers   `yaml:"workers"`
	Log        Log       `yaml:"log"`
	Storage    Storage   `yaml:"storage"`
	Sunet      Sunet     `yaml:"sunet"`
}

// Config represent the complete config file structure
type Config struct {
	EduID struct {
		Worker struct {
			Cleaner Cfg `yaml:"cleaner"`
		} `yaml:"worker"`
	} `yaml:"eduid"`
}
