package conf

import "go-leaf/internal/types"

const (
	EnvLocal   = "local"
	EnvTest    = "test"
	EnvProduct = "product"
)

type Conf struct {
	// etcd的地址
	Etcd string `json:"etcd"`
	// 是否强依赖时钟
	TimeSwitch bool   `json:"time_switch"`
	Server     Server `json:"server"`
	JaegerUrl  string `json:"jaeger_url"`
	Env        string `json:"env"`
	Log        Log    `json:"log"`
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
