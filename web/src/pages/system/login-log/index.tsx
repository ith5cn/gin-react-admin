import { useRef } from "react";
import { Col, DatePicker, Form, Input, Select, Tag, message } from "antd";
import dayjs, { type Dayjs } from "dayjs";
import Ith5Table, { type ColumnDef, type TableRef } from "@/components/ith5ui/ith5-table";
import { loginLogDeleteApi, loginLogListApi } from "@/api/system/login-log";

interface LoginLogRecord {
  id: number;
  username?: string;
  status: number;
  ip?: string;
  ipLocation?: string;
  os?: string;
  browser?: string;
  message?: string;
  loginTime?: string;
}

const statusOptions = [
  { label: "成功", value: 1 },
  { label: "失败", value: 2 },
];

const normalizeSearchParams = (params: Record<string, unknown>) => {
  const next = { ...params };
  const loginTime = next.loginTime;
  if (Array.isArray(loginTime) && loginTime.length === 2) {
    next.loginTime = (loginTime as Dayjs[]).map((item) => item.format("YYYY-MM-DD HH:mm:ss")).join(",");
  }
  return next;
};

const formatDateTime = (value?: string) => value ? dayjs(value).format("YYYY-MM-DD HH:mm:ss") : "-";

const LoginLogIndex = () => {
  const tableRef = useRef<TableRef>(null);

  return (
    <Ith5Table
      ref={tableRef}
      searchFields={
        <>
          <Col span={6}>
            <Form.Item name="username" label="登录用户">
              <Input placeholder="请输入登录用户" allowClear />
            </Form.Item>
          </Col>
          <Col span={6}>
            <Form.Item name="status" label="登录状态">
              <Select placeholder="请选择登录状态" options={statusOptions} allowClear />
            </Form.Item>
          </Col>
          <Col span={6}>
            <Form.Item name="ip" label="登录 IP">
              <Input placeholder="请输入登录 IP" allowClear />
            </Form.Item>
          </Col>
          <Col span={8}>
            <Form.Item name="loginTime" label="登录时间">
              <DatePicker.RangePicker showTime style={{ width: "100%" }} />
            </Form.Item>
          </Col>
        </>
      }
      options={{
        api: (params) => loginLogListApi(normalizeSearchParams(params)),
        delete: {
          show: true,
          auth: ["system/login-log/destroy"],
          confirmText: "确定要删除这条登录日志吗？",
          func: async (record: LoginLogRecord) => {
            await loginLogDeleteApi(record.id);
            message.success("删除成功");
            tableRef.current?.refresh();
          },
        },
      }}
      columns={[
        { title: "登录用户", dataIndex: "username", width: 140, render: (value?: string) => value || "-" },
        {
          title: "登录状态",
          dataIndex: "status",
          width: 110,
          render: (value: number) => value === 1 ? <Tag color="success">成功</Tag> : <Tag color="error">失败</Tag>,
        },
        { title: "登录 IP", dataIndex: "ip", width: 150, render: (value?: string) => value || "-" },
        { title: "登录地点", dataIndex: "ipLocation", width: 140, render: (value?: string) => value || "-" },
        { title: "操作系统", dataIndex: "os", width: 120, render: (value?: string) => value || "-" },
        { title: "浏览器", dataIndex: "browser", width: 120, render: (value?: string) => value || "-" },
        { title: "登录信息", dataIndex: "message", width: 150, render: (value?: string) => value || "-" },
        { title: "登录时间", dataIndex: "loginTime", width: 180, render: formatDateTime },
      ] as ColumnDef[]}
    />
  );
};

export default LoginLogIndex;
