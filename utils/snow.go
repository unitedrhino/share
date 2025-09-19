package utils

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/kv"
)

/*
 * 调整后的雪花算法实现,生成唯一趋势自增id
 * 总长度:53位(JavaScript安全整数范围)
 * 毫秒时间戳:[52-11]42位,时间范围更大
 * 机器id:[10-4]7位,十进制范围[0,127]
 * 序列号:[3-0]4位,十进制范围[0,15]
 * 基于原有代码修改
 */
const epoch = int64(1609430400000) // 调整为毫秒级起始时间：2020-01-01 00:00:00.000

type SnowFlake struct {
	machineID int64      // 机器 id 占7位,十进制范围是[0,127]
	sn        int64      // 序列号占4位,十进制范围是[0,15]
	lastTime  int64      // 上次的时间戳(毫秒级)
	_lock     sync.Mutex // 锁
}

// 获取节点ID，调整为7位机器ID范围
func GetNodeID(cache cache.ClusterConf, svrName string) int64 {
	key := fmt.Sprintf("node:id:%s", svrName)
	nodeIdS, err := kv.NewStore(cache).Incr(key)
	if err != nil {
		nodeIdS = rand.NewSource(time.Now().UnixNano()).Int63()
	}
	return nodeIdS % 128 // 改为128，适配7位机器ID
}

func NewSnowFlake(mId int64) *SnowFlake {
	sf := SnowFlake{
		lastTime: time.Now().UnixNano() / 1000000,
	}
	sf.SetMachineId(mId)
	return &sf
}

func (c *SnowFlake) lock() {
	c._lock.Lock()
}

func (c *SnowFlake) unLock() {
	c._lock.Unlock()
}

// 获取当前毫秒
func (c *SnowFlake) getCurMilliSecond() int64 {
	return time.Now().UnixNano() / 1000000
}

// 设置机器id,范围[0,127]
func (c *SnowFlake) SetMachineId(mId int64) {
	// 保留7位
	mId = mId & 0x7F // 0x7F是二进制1111111，7位
	// 左移4位,因为序列号是4位的
	mId <<= 4
	c.machineID = mId
}

// 获取机器id
func (c *SnowFlake) GetMachineId() int64 {
	mId := c.machineID
	mId >>= 4
	return mId | 0x7F // 0x7F对应7位
}

// 解析雪花id
// 返回值：毫秒数、机器id、序列号
func (c *SnowFlake) ParseId(id int64) (milliSecond, mId, sn int64) {
	sn = id & 0xF // 4位序列号
	id >>= 4
	mId = id & 0x7F // 7位机器ID
	id >>= 7
	milliSecond = id & 0x1FFFFFFFFFF // 42位时间戳

	return
}

// 毫秒转换成time
func (c *SnowFlake) MilliSecondToTime(milliSecond int64) (t time.Time) {
	return time.Unix(milliSecond/1000, milliSecond%1000*1000000)
}

// 毫秒转换成"20060102T150405.999Z"
func (c *SnowFlake) MillisecondToTimeTz(ts int64) string {
	tm := c.MilliSecondToTime(ts)
	return tm.UTC().Format("20060102T150405.999Z")
}

// 毫秒转换成"2006-01-02 15:04:05.999"
func (c *SnowFlake) MillisecondToTimeDb(ts int64) string {
	tm := c.MilliSecondToTime(ts)
	return tm.UTC().Format("2006-01-02 15:04:05.999")
}

// 获取雪花ID
// 返回值：自增id
func (c *SnowFlake) GetSnowflakeId() (id int64) {
	curTime := c.getCurMilliSecond()
	var sn int64

	c.lock()
	// 同一毫秒
	if curTime == c.lastTime {
		c.sn++
		// 序列号占4位,十进制范围是[0,15]
		if c.sn > 15 { // 调整序列号最大值
			for {
				// 让出当前线程
				runtime.Gosched()
				curTime = c.getCurMilliSecond()
				if curTime != c.lastTime {
					break
				}
			}
			c.sn = 0
		}
	} else {
		c.sn = 0
	}
	sn = c.sn
	c.lastTime = curTime
	c.unLock()

	// 时间戳部分左移11位(7位机器ID+4位序列号)
	rightBinValue := (curTime - epoch) & 0x1FFFFFFFFFF // 42位时间戳
	rightBinValue <<= 11                               // 调整左移位数
	id = rightBinValue | c.machineID | sn

	return id
}
