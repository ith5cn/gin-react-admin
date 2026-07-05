import request from "@/utils/request";

export interface ServerMonitorData {
  host: {
    hostname: string;
    os: string;
    platform: string;
    arch: string;
    uptime: number;
  };
  cpu: {
    cores: number;
    modelName: string;
    usedPercent: number;
  };
  memory: {
    total: number;
    used: number;
    usedPercent: number;
  };
  disk: {
    total: number;
    used: number;
    usedPercent: number;
  };
  runtime: {
    goVersion: string;
    goroutines: number;
    heapAlloc: number;
    gcCount: number;
    pid: number;
    startTime: string;
    uptime: number;
  };
  redis: {
    version: string;
    mode: string;
    connectedClients: string;
    usedMemoryHuman: string;
    uptimeSeconds: string;
    totalCommands: string;
    keyCount: number;
    available: boolean;
  };
}

/**
 * 服务监控快照
 */
export const serverMonitorApi = () => {
  return request.get<ServerMonitorData>("/system/monitor/server");
};
