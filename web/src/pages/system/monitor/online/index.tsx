import { useRef } from "react";
import { Button, Col, Form, Input, message, Popconfirm } from "antd";
import { LogoutOutlined } from "@ant-design/icons";
import dayjs from "dayjs";

import Ith5Table, { type ColumnDef, type TableRef } from "@/components/ith5ui/ith5-table";
import { kickOnlineUserApi, onlineUserListApi, type OnlineUserItem } from "@/api/system/online";
import { useAuthStore } from "@/store/auth";

const OnlineUserIndex = () => {
  const tableRef = useRef<TableRef>(null);
  const codes = useAuthStore((state) => state.codes);
  const canKick = codes.includes("*") || codes.includes("system/online/kick");

  const handleKick = async (record: OnlineUserItem) => {
    const res = await kickOnlineUserApi(record.accessJti);
    if (res.code === 0) {
      message.success("已踢下线");
      tableRef.current?.refresh();
    }
  };

  return (
    <Ith5Table
      ref={tableRef}
      searchFields={
        <>
          <Col span={6}>
            <Form.Item name="username" label="用户名">
              <Input placeholder="请输入用户名" allowClear />
            </Form.Item>
          </Col>
          <Col span={6}>
            <Form.Item name="ip" label="登录IP">
              <Input placeholder="请输入登录IP" allowClear />
            </Form.Item>
          </Col>
        </>
      }
      options={{ api: onlineUserListApi }}
      operationBeforeExtend={(record: OnlineUserItem) =>
        canKick ? (
          <Popconfirm title={`确定将 ${record.username} 踢下线吗？`} onConfirm={() => handleKick(record)}>
            <Button type="link" size="small" danger icon={<LogoutOutlined />}>
              踢下线
            </Button>
          </Popconfirm>
        ) : null
      }
      columns={
        [
          { title: "用户名", dataIndex: "username", width: 160 },
          { title: "登录IP", dataIndex: "ip", width: 160 },
          { title: "操作系统", dataIndex: "os", width: 140 },
          { title: "浏览器", dataIndex: "browser", width: 140 },
          {
            title: "登录时间",
            dataIndex: "loginTime",
            width: 180,
            render: (text: string) => (text ? dayjs(text).format("YYYY-MM-DD HH:mm:ss") : "-"),
          },
        ] as ColumnDef[]
      }
    />
  );
};

export default OnlineUserIndex;
