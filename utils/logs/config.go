package logs

// 日志配置项
type LoggerConfig struct {
	Mode        string   `json:"mode"` // 日志类型：console,file,kafka
	Env         string   `json:"env,default=pro"`
	Level       string   `json:"level"`                // 日志等级：INFO,DEBUG,ERROR等
	ServiceName string   `json:"serviceName,optional"` // 服务名称
	KafkaInfo   struct { // kafka配置
		Address []string `json:"address"` // kafka服务IP
		Topic   string   `json:"topic"`   // 主题
	} `json:"kafkaInfo,optional"`
	FileInfo struct { // 文件配置
		Path string `json:"path"` // 写文件路径
	} `json:"fileInfo,optional"`
}
