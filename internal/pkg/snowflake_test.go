package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSnowflake_Generate_Auto(t *testing.T) {

	snowflake, err := NewSnowflake(10, 12, 1, 1685502428004, ModeAuto, 5*time.Millisecond, 200)

	assert.NoError(t, err)
	id, err := snowflake.Generate()
	assert.NoError(t, err)
	t.Log(id)
	t.Log(snowflake.Parse(id))
	t.Log("----------------")

	id, err = snowflake.Generate()
	assert.NoError(t, err)
	t.Log(id)
	t.Log(snowflake.Parse(id))
	t.Log("----------------")

	snowflake.lastTimestamp = snowflake.lastTimestamp + 5
	snowflake.step++
	id, err = snowflake.Generate()
	assert.NoError(t, err)
	t.Log(id)
	t.Log(snowflake.Parse(id))
	t.Log("----------------")

	time.Sleep(5 * time.Millisecond)
	id, err = snowflake.Generate()
	assert.NoError(t, err)
	t.Log(id)
	t.Log(snowflake.Parse(id))
}
