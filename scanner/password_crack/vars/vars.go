package vars

import (
	"github.com/patrickmn/go-cache"
	"gopkg.in/cheggaaa/pb.v2"
	"strings"
	"sync"
	"time"
)

var (
	IpList = "ip_list.txt"
	UserDict = "user.dic"
	PassDict = "pass.dic"

	// 启动时间
	StartTime time.Time
	ResultFile = "password_crack.txt"

	// 超时时间
	TimeOut = 3 * time.Second
	ScanNum = 5000

	DebugMode bool
	// 弱口令扫描进度条
	ProgressBar *pb.ProgressBar
	// 检测端口是否开放的进度条
	ProgressBarActive *pb.ProgressBar
)

var (
	CacheService *cache.Cache

	PortNames = map[int]string{
		22:   "ssh",
		3306: "mysql",
		6379: "redis",
		1433: "mssql",
		5432: "postgresql",
		27017: "mongodb",
	}

	// 标记特定服务的用户破解成功不再尝试
	SuccessHash sync.Map
    SupportProtocols map[string]bool
)

func init() {
	CacheService = cache.New(cache.NoExpiration, cache.DefaultExpiration)

	SupportProtocols = make(map[string]bool)
	for _, proto = range PortNames {
		SupportProtocols[strings.ToUpper(proto)] = true
	}
}