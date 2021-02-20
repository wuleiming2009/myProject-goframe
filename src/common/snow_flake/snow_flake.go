package snow_flake

import "time"

var lastTimeStamp int64
var sn, machineId, datacenterId uint64

func GenSnowFlake() uint64 {
	machineId = 1
	datacenterId = 2

	// 如果想让时间戳范围更长，也可以减去一个日期
	curTimeStamp := time.Now().UnixNano() / 1000000

	if curTimeStamp == lastTimeStamp {
		// 2的12次方 -1 = 4095，每毫秒可产生4095个ID
		if sn > 4095 {
			time.Sleep(time.Millisecond)
			curTimeStamp = time.Now().UnixNano() / 1000000
			sn = 0
		}
	} else {
		sn = 0
	}
	sn++
	lastTimeStamp = curTimeStamp
	// 应为时间戳后面有22位，所以向左移动22位
	curTimeStamp = curTimeStamp << 22
	machineId = machineId << 17
	datacenterId = datacenterId << 12
	// 通过与运算把各个部位连接在一起
	return uint64(curTimeStamp) | machineId | datacenterId | sn
}
