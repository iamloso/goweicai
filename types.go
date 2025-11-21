package gowencai

// QueryOptions contains all options for the Get query
type QueryOptions struct {
	Query         string                 // 查询问句 (必填)
	SortKey       string                 // 排序字段
	SortOrder     string                 // 排序规则: asc/desc
	Page          int                    // 页号，默认1
	PerPage       int                    // 每页数据条数，默认100
	Loop          interface{}            // 是否循环分页: false/true/数字
	QueryType     string                 // 查询类型: stock/zhishu/fund等
	Retry         int                    // 重试次数，默认10
	Sleep         int                    // 请求间隔秒数，默认0
	Log           bool                   // 是否打印日志
	Pro           bool                   // 是否付费版
	Cookie        string                 // Cookie值 (必填)
	RequestParams map[string]interface{} // 额外的请求参数
	NoDetail      bool                   // 详情类问题不返回字典
	Find          []string               // 指定标的排在前面
	UserAgent     string                 // 自定义User-Agent
}

// RobotDataParams contains the parsed robot data
type RobotDataParams struct {
	Data      interface{}        `json:"data"`
	RowCount  int                `json:"row_count"`
	URL       string             `json:"url"`
	URLParams map[string]string  `json:"url_params"`
}

// Component represents a data component
type Component struct {
	ShowType    string                 `json:"show_type"`
	CID         string                 `json:"cid"`
	PUUID       string                 `json:"puuid"`
	UUID        string                 `json:"uuid"`
	Data        map[string]interface{} `json:"data"`
	Config      map[string]interface{} `json:"config"`
	TitleConfig map[string]interface{} `json:"title_config"`
	TabList     []interface{}          `json:"tab_list"`
}

// RobotResponse represents the response from get-robot-data endpoint
type RobotResponse struct {
	Data struct {
		Answer []struct {
			Txt []struct {
				Content interface{} `json:"content"`
			} `json:"txt"`
		} `json:"answer"`
	} `json:"data"`
}
