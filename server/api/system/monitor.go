package system

import (
	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

// ServerMonitor 服务监控快照（主机/CPU/内存/磁盘/Go 运行时/Redis）。
func ServerMonitor(c *gin.Context) {
	successOrFail(c, systemService.ServerMonitorInfo(), nil)
}
