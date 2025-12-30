package conf

// Bootstrap is the configuration structure
type Bootstrap struct {
	Server    *Server    `json:"server"`
	Debug     bool       `json:"debug"` // 全局调试开关
	Scheduler *Scheduler `json:"scheduler"`
	Data      *Data      `json:"data"`
	Wencai    *Wencai    `json:"wencai"`
	BaseInfo  *BaseInfo  `json:"baseinfo"`
	ZtInfo    *ZtInfo    `json:"ztinfo"`
}

// Server is the server configuration
type Server struct {
	Http *Server_HTTP `json:"http"`
	Grpc *Server_GRPC `json:"grpc"`
}

type Server_HTTP struct {
	Addr    string `json:"addr"`
	Timeout string `json:"timeout"`
}

type Server_GRPC struct {
	Addr    string `json:"addr"`
	Timeout string `json:"timeout"`
}

// Scheduler is the scheduler configuration
type Scheduler struct {
	Stock      *SchedulerTask `json:"stock"`        // 股票涨停数据任务
	BaseInfo   *SchedulerTask `json:"base_info"`    // 基础数据任务
	ZtInfo     *SchedulerTask `json:"zt_info"`      // 涨停数据任务
	RunOnStart bool           `json:"run_on_start"` // 启动时立即执行
}

// SchedulerTask 单个定时任务配置
type SchedulerTask struct {
	Enabled bool   `json:"enabled"` // 是否启用
	Cron    string `json:"cron"`    // Cron 表达式
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

// BaseInfo is the base info API configuration
type BaseInfo struct {
	Query  string `json:"query"`
	Cookie string `json:"cookie"`
}

// ZtInfo is the zt info API configuration
type ZtInfo struct {
	Query  string `json:"query"`
	Cookie string `json:"cookie"`
}
