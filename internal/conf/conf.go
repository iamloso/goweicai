package conf

// Bootstrap is the configuration structure
type Bootstrap struct {
	Server    *Server    `json:"server"`
	Scheduler *Scheduler `json:"scheduler"`
	Data      *Data      `json:"data"`
	Wencai    *Wencai    `json:"wencai"`
}

// Server is the server configuration
type Server struct {
	Http *Server_HTTP `json:"http"`
}

type Server_HTTP struct {
	Addr    string `json:"addr"`
	Timeout string `json:"timeout"`
}

// Scheduler is the scheduler configuration
type Scheduler struct {
	Cron       string `json:"cron"`
	RunOnStart bool   `json:"run_on_start"`
}

// Data is the data layer configuration
type Data struct {
	Database *Data_Database `json:"database"`
}

type Data_Database struct {
	Driver string `json:"driver"`
	Source string `json:"source"`
}

// Wencai is the wencai API configuration
type Wencai struct {
	Query  string `json:"query"`
	Cookie string `json:"cookie"`
}
