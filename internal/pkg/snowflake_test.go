package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSnowflake_Generate_Auto(t *testing.T) {

	snowflake, err := NewSnowflake(8, 12, 1, 1685502428004, ModeAuto, 5*time.Millisecond, 200)

	s, ok := snowflake.(*Snowflake)

	assert.NoError(t, err)
	id, err := snowflake.Generate()
	assert.NoError(t, err)
	t.Log(id)
	t.Log(snowflake.Parse(id))

	id, err = snowflake.Generate()
	assert.NoError(t, err)
	t.Log(id)
	t.Log(snowflake.Parse(id))
	// 模拟时钟回退
	if ok {
		s.lastTimestamp = s.lastTimestamp + 5
		s.step++
	}

	id, err = snowflake.Generate()
	assert.NoError(t, err)
	t.Log(id)
	t.Log(snowflake.Parse(id))
	// 等待时钟回归正常
	time.Sleep(5 * time.Millisecond)
	id, err = snowflake.Generate()
	assert.NoError(t, err)
	t.Log(id)
	t.Log(snowflake.Parse(id))
}
