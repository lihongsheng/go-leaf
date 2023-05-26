package pkg

import (
	"github.com/stretchr/testify/assert"
	"message-center/internal/conf"
	"testing"
)

func TestGrpcClientConn(t *testing.T) {
	addr := "grpcmessage.center.xmslol.cn:9010"
	logger := NewZapFileLog(conf.Log{Path: ""})

	cn, err := GrpcClientConn(addr, logger)
	t.Log(cn)
	assert.NoError(t, err)
}
