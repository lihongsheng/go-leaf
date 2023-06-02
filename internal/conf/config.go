package conf

import "go-leaf/internal/types"

const (
	EnvLocal   = "local"
	EnvTest    = "test"
	EnvProduct = "product"
)

type Conf struct {
	// etcd的地址
	EtcdAdds []string `json:"etcd_adds"`
	// 是否强依赖时钟
	Server    Server    `json:"server"`
	JaegerUrl string    `json:"jaeger_url"`
	Env       string    `json:"env"`
	Log       Log       `json:"log"`
	Snowflake Snowflake `json:"snowflake"`
	Pprof     string    `json:"pprof"`
	Auth      Auth      `json:"auth"`
}

type Auth struct {
	User string `json:"user"`
	Pwd  string `json:"pwd"`
}

type Snowflake struct {
	// # 模式1 正常模式：时钟回退会报错，
	// 模式2：等待模式，时钟回退切小于多少毫秒内会等待。
	// 模式3：自动切换模式
	//      当时钟回退：
	//         1. 记录当前回退的时间 back_time
	//         2. 让当前计数器继续保持递增 current_time + (back_time-back_time) + sqp++
	//         3. 当时间正常的时候,大于等于 current_time的时候 使用正常时间【对与系统负载较大的时候，时间有可能永远回归不了正常】
	Mode     string         `json:"mode"`
	WaitTime types.Duration `json:"wait_time"`
	// 数据中心标志位的长度
	DataLen int `json:"data_len"`
	// 步长标志位的长度
	SepLen int `json:"sep_len"`
	// 开始时间微妙数
	StartTime int64 `json:"start_time"`
}

type Server struct {
	// http 接口
	HTTPPort string `json:"http_port"`
	// grpc接口
	GrpcPort string         `json:"grpc_port"`
	TimeOut  types.Duration `json:"time_out"`
	Name     string         `json:"name"`
	Version  string         `json:"version"`
}

type Log struct {
	// 地址为空不启用文件模式
	Path string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	// 文件最大（MB单位）
	MaxSize int32 `protobuf:"varint,2,opt,name=max_size,json=maxSize,proto3" json:"max_size,omitempty"`
	// 文件最大备份天数（MB单位）
	MaxBackUp int32 `protobuf:"varint,3,opt,name=max_back_up,json=maxBackUp,proto3" json:"max_back_up,omitempty"`
	// 切分的数据
	MaxAge int32 `protobuf:"varint,4,opt,name=max_age,json=maxAge,proto3" json:"max_age,omitempty"`
	// 是否压缩
	Compress bool `protobuf:"varint,5,opt,name=compress,proto3" json:"compress,omitempty"`
}

func (b Conf) IsProduct() bool {
	return b.Env == EnvProduct
}

func (b Conf) IsTest() bool {
	return b.Env == EnvTest
}

func (b Conf) IsLocal() bool {
	return b.Env == EnvLocal
}
