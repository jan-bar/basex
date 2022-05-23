package basex

import (
	"crypto/md5"
	"crypto/rand"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

var machineID, objectIDCounter = initRandId()

// 参考项目: https://github.com/rs/xid/blob/5cbb911d27d5efc5f0be784aac766db82ebd067f/id.go#L113
func initRandId() ([]byte, uint32) {
	mw := md5.New()

	// 注意docker容器内pid都是1
	mw.Write(strconv.AppendInt(nil, int64(os.Getpid()), 10))

	// linux 下的 Platform MachineID
	data, err := os.ReadFile("/etc/machine-id")
	if err == nil {
		mw.Write(data)
	}

	// 在docker容器内包含docker容器id,不同容器不相同
	data, err = os.ReadFile("/proc/self/cpuset")
	if err == nil {
		mw.Write(data)
	}

	// 每个docker容器启动进程命令不同,不同模块不相同
	data, err = os.ReadFile("/proc/1/cmdline")
	if err == nil {
		mw.Write(data)
	}

	id := make([]byte, 4)
	_, err = rand.Reader.Read(id)
	if err == nil {
		mw.Write(id)
	}
	// 进程启动使用随机自增ID
	uid := uint32(id[0])<<24 | uint32(id[1])<<16 |
		uint32(id[2])<<8 | uint32(id[3])

	// 取md5的前4位
	copy(id, mw.Sum(nil))
	id[0] %= 0x9a // 确保 len(base62(src)) < 17
	return id, uid
}

const (
	base62Std = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	base62Len = len(base62Std)
)

func GetUniqueID(pre string) string {
	id := make([]byte, 16)
	// Machine, first 4 bytes of md5(xxx)
	id[0] = machineID[0]
	id[1] = machineID[1]
	id[2] = machineID[2]
	id[3] = machineID[3]
	// Timestamp, 4 bytes, big endian
	timestamp := time.Now().Unix()
	id[4] = byte(timestamp >> 24)
	id[5] = byte(timestamp >> 16)
	id[6] = byte(timestamp >> 8)
	id[7] = byte(timestamp)
	// Increment, 4 bytes, big endian
	count := atomic.AddUint32(&objectIDCounter, 1)
	id[8] = byte(count >> 24)
	id[9] = byte(count >> 16)
	id[10] = byte(count >> 8)
	id[11] = byte(count)

	// 将12字节数据转换到长度16的数组中
	digits := make([]int, 0, 16)
	for i := 0; i < 12; i++ {
		carry := int(id[i])

		for j := 0; j < len(digits); j++ {
			carry += digits[j] << 8
			digits[j] = carry % base62Len
			carry /= base62Len
		}

		for carry > 0 {
			digits = append(digits, carry%base62Len)
			carry /= base62Len
		}
	}

	cur := 0 // 生成16个字符的明文字符串
	for q := len(digits) - 1; q >= 0; q-- {
		id[cur] = base62Std[digits[q]]
		cur++
	}
	// len(digits) < 16,这时后面数据补零
	for cur < 16 {
		id[cur] = base62Std[0]
		cur++
	}
	// 将前缀和唯一ID拼接,返回最终编码
	return pre + string(id)
}
