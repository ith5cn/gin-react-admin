import { useCallback, useEffect, useState } from "react";
import { Button, Card, Descriptions, Progress, Spin, Tag } from "antd";
import { ReloadOutlined } from "@ant-design/icons";

import { serverMonitorApi, type ServerMonitorData } from "@/api/system/monitor";

const formatBytes = (bytes: number) => {
  if (!bytes) return "0 B";
  const units = ["B", "KB", "MB", "GB", "TB"];
  let value = bytes;
  let unitIndex = 0;
  while (value >= 1024 && unitIndex < units.length - 1) {
    value /= 1024;
    unitIndex += 1;
  }
  return `${value.toFixed(2)} ${units[unitIndex]}`;
};

const formatDuration = (seconds: number) => {
  if (!seconds) return "-";
  const days = Math.floor(seconds / 86400);
  const hours = Math.floor((seconds % 86400) / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  if (days > 0) return `${days} 天 ${hours} 小时`;
  if (hours > 0) return `${hours} 小时 ${minutes} 分钟`;
  return `${minutes} 分钟`;
};

const percentStatus = (percent: number) => (percent >= 90 ? "exception" : "normal");

const ServerMonitorIndex = () => {
  const [data, setData] = useState<ServerMonitorData | null>(null);
  const [loading, setLoading] = useState(false);

  const fetchData = useCallback(async () => {
    setLoading(true);
    try {
      const res = await serverMonitorApi();
      setData(res.data);
    } catch {
      // 错误已由 request 拦截器统一处理
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return (
    <div className="flex flex-col gap-4">
      <div className="flex justify-end">
        <Button icon={<ReloadOutlined />} loading={loading} onClick={fetchData}>
          刷新
        </Button>
      </div>

      <Spin spinning={loading && !data}>
        <div className="grid grid-cols-1 gap-4 lg:grid-cols-3">
          <Card title="CPU" size="small">
            <div className="flex justify-center py-2">
              <Progress
                type="dashboard"
                percent={data?.cpu.usedPercent ?? 0}
                status={percentStatus(data?.cpu.usedPercent ?? 0)}
              />
            </div>
            <Descriptions column={1} size="small">
              <Descriptions.Item label="核心数">{data?.cpu.cores ?? "-"}</Descriptions.Item>
              <Descriptions.Item label="型号">{data?.cpu.modelName || "-"}</Descriptions.Item>
            </Descriptions>
          </Card>

          <Card title="内存" size="small">
            <div className="flex justify-center py-2">
              <Progress
                type="dashboard"
                percent={data?.memory.usedPercent ?? 0}
                status={percentStatus(data?.memory.usedPercent ?? 0)}
              />
            </div>
            <Descriptions column={1} size="small">
              <Descriptions.Item label="总内存">{formatBytes(data?.memory.total ?? 0)}</Descriptions.Item>
              <Descriptions.Item label="已使用">{formatBytes(data?.memory.used ?? 0)}</Descriptions.Item>
            </Descriptions>
          </Card>

          <Card title="磁盘（根分区）" size="small">
            <div className="flex justify-center py-2">
              <Progress
                type="dashboard"
                percent={data?.disk.usedPercent ?? 0}
                status={percentStatus(data?.disk.usedPercent ?? 0)}
              />
            </div>
            <Descriptions column={1} size="small">
              <Descriptions.Item label="总容量">{formatBytes(data?.disk.total ?? 0)}</Descriptions.Item>
              <Descriptions.Item label="已使用">{formatBytes(data?.disk.used ?? 0)}</Descriptions.Item>
            </Descriptions>
          </Card>
        </div>

        <div className="mt-4 grid grid-cols-1 gap-4 lg:grid-cols-3">
          <Card title="主机信息" size="small">
            <Descriptions column={1} size="small">
              <Descriptions.Item label="主机名">{data?.host.hostname || "-"}</Descriptions.Item>
              <Descriptions.Item label="操作系统">{data?.host.platform || "-"}</Descriptions.Item>
              <Descriptions.Item label="内核架构">{data?.host.arch || "-"}</Descriptions.Item>
              <Descriptions.Item label="开机时长">{formatDuration(data?.host.uptime ?? 0)}</Descriptions.Item>
            </Descriptions>
          </Card>

          <Card title="Go 运行时" size="small">
            <Descriptions column={1} size="small">
              <Descriptions.Item label="Go 版本">{data?.runtime.goVersion || "-"}</Descriptions.Item>
              <Descriptions.Item label="Goroutine 数">{data?.runtime.goroutines ?? "-"}</Descriptions.Item>
              <Descriptions.Item label="堆内存">{formatBytes(data?.runtime.heapAlloc ?? 0)}</Descriptions.Item>
              <Descriptions.Item label="GC 次数">{data?.runtime.gcCount ?? "-"}</Descriptions.Item>
              <Descriptions.Item label="启动时间">{data?.runtime.startTime || "-"}</Descriptions.Item>
              <Descriptions.Item label="运行时长">{formatDuration(data?.runtime.uptime ?? 0)}</Descriptions.Item>
            </Descriptions>
          </Card>

          <Card
            title={
              <span>
                Redis{" "}
                {data &&
                  (data.redis.available ? <Tag color="green">在线</Tag> : <Tag color="red">不可用</Tag>)}
              </span>
            }
            size="small"
          >
            <Descriptions column={1} size="small">
              <Descriptions.Item label="版本">{data?.redis.version || "-"}</Descriptions.Item>
              <Descriptions.Item label="模式">{data?.redis.mode || "-"}</Descriptions.Item>
              <Descriptions.Item label="连接数">{data?.redis.connectedClients || "-"}</Descriptions.Item>
              <Descriptions.Item label="内存占用">{data?.redis.usedMemoryHuman || "-"}</Descriptions.Item>
              <Descriptions.Item label="Key 数量">{data?.redis.keyCount ?? "-"}</Descriptions.Item>
              <Descriptions.Item label="运行时长">
                {formatDuration(Number(data?.redis.uptimeSeconds) || 0)}
              </Descriptions.Item>
            </Descriptions>
          </Card>
        </div>
      </Spin>
    </div>
  );
};

export default ServerMonitorIndex;
