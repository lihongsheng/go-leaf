package pkg

import (
	"context"
	apiError "go-leaf/api/error"
	"sync"
	"time"
)

type Generate interface {
	Generate() (int64, error)
	Parse(id int64) (time int64, node int64, seq int64, err error)
}

const (
	Epoch int64 = 1288834974657
	// NodeBits holds the number of bits to use for Node
	// Remember, you have a total 22 bits to share between Node/Step
	NodeBits uint8 = 10
	// StepBits holds the number of bits to use for Step
	// Remember, you have a total 22 bits to share between Node/Step
	StepBits uint8 = 12
)

type Mode int8

const (
	// ModeNormal is
	ModeNormal Mode = iota
	// ModeWait is logger info level.
	ModeWait
	// ModeAuto is auto add mills.
	ModeAuto
)

// Snowflake
// +--------------------------------------------------------------------------+
// | 1 Bit Unused | 41 Bit Timestamp |  10 Bit NodeID  |   12 Bit Sequence ID |
// +--------------------------------------------------------------------------+
type Snowflake struct {
	nodeBits           uint8
	stepBits           uint8
	node               int64
	step               int64
	stepMax            int64
	epoch              int64
	lastTimestamp      int64
	preTimestamp       int64
	mode               Mode
	mutex              sync.Mutex
	timestampLeftShift uint8
	maxWaitTime        time.Duration
	maxWait            int
	wait               int
}

func NewSnowflake(nodeBits, stepBits uint8, node, epoch int64, mode Mode, maxWaitTime time.Duration, maxWait int) (Generate, error) {
	if (nodeBits + stepBits) > 22 {
		return nil, apiError.ErrorSystemInvalidConfError(" you have a total 22 bits to share between Node/Step")
	}
	if NodeBits < uint8(1) {
		nodeBits = NodeBits
	}
	if StepBits < uint8(1) {
		stepBits = StepBits
	}
	if node > (1<<nodeBits - 1) {
		return nil, apiError.ErrorSystemInvalidConfError(" node[%d] is out maxNode[%d]", node, 1<<nodeBits-1)
	}
	if maxWaitTime < 1 {
		maxWaitTime = 5 * time.Millisecond
	}
	return &Snowflake{
		nodeBits:           nodeBits,
		stepBits:           stepBits,
		node:               node,
		step:               0,
		stepMax:            -1 ^ (-1 << stepBits),
		epoch:              epoch,
		mode:               mode,
		mutex:              sync.Mutex{},
		timestampLeftShift: stepBits + nodeBits,
		maxWaitTime:        maxWaitTime,
		maxWait:            maxWait,
	}, nil
}

func (s *Snowflake) Generate() (int64, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	var err error
	currentTime := time.Now().UnixMilli()
	if currentTime < s.lastTimestamp { // 时钟回滚
		switch s.mode {
		case ModeNormal: // 直接返回错，让客户端重试
			return 0, apiError.ErrorSystemClockRollbackError("Clock moved backwards.  Refusing to generate id for %d milliseconds", s.lastTimestamp-currentTime)
		case ModeWait: // 等待模式
			currentTime, err = s.waitTimeMills(s.lastTimestamp)
			if err != nil {
				return 0, err
			}
		case ModeAuto: // 自增模式, 自动增加 lastTimestamp
			s.step = (s.step + 1) & s.stepMax
			if s.step == 0 {
				currentTime = s.timeMills(currentTime)
				currentTime = s.lastTimestamp + 1
			} else {
				currentTime = s.lastTimestamp
			}
		}
	} else if currentTime == s.lastTimestamp {
		s.step = (s.step + 1) & s.stepMax
		if s.step == 0 {
			currentTime = s.timeMills(s.lastTimestamp)
		}
	} else {
		s.step = 0
	}
	s.lastTimestamp = currentTime
	return ((s.lastTimestamp - s.epoch) << (s.nodeBits + s.stepBits)) | (s.node << s.stepBits) | s.step, nil
}

// waitTimeMills is wait time
// if time out maxWaitTime return error.
func (s *Snowflake) waitTimeMills(lastTime int64) (int64, error) {
	if s.wait >= s.maxWait {
		return 0, apiError.ErrorSystemClockRollbackError("exceed maxWait %d", s.maxWait)
	}
	s.wait++
	ctx, cancel := context.WithTimeout(context.Background(), s.maxWaitTime)
	defer cancel()
	currentTime := time.Now().UnixMilli()
	for {
		select {
		case <-ctx.Done():
			return 0, apiError.ErrorSystemClockRollbackError("wait time out")
		default:
			currentTime = time.Now().UnixMilli()
			if currentTime > lastTime {
				s.wait = 0
				return currentTime, nil
			}
		}
	}
}

// timeMills
func (s *Snowflake) timeMills(lastTime int64) int64 {
	currentTime := time.Now().UnixMilli()
	for currentTime <= lastTime {
		currentTime = time.Now().UnixMilli()
	}
	return currentTime
}

func (s *Snowflake) Parse(id int64) (time int64, node int64, seq int64, err error) {
	time = id>>(s.nodeBits+s.stepBits) + s.epoch
	node = (1<<s.nodeBits - 1) << s.stepBits
	node = id & node >> s.stepBits
	seq = id & s.stepMax
	return
}
