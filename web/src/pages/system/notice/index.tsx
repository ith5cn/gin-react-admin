import { useRef } from "react";
import { Col, Form, Input, message, Select, Tag } from "antd";
import dayjs from "dayjs";

import Ith5Table, { type ColumnDef, type TableRef } from "@/components/ith5ui/ith5-table";
import { noticeDeleteApi, noticeListApi } from "@/api/system/notice";
import NoticeEdit, { type NoticeEditRef } from "./edit";

const noticeTypeOptions = [
  { label: "通知", value: 1 },
  { label: "公告", value: 2 },
];

const NoticeIndex = () => {
  const editRef = useRef<NoticeEditRef>(null);
  const tableRef = useRef<TableRef>(null);

  return (
    <>
      <Ith5Table
        ref={tableRef}
        searchFields={
          <>
            <Col span={6}>
              <Form.Item name="title" label="标题">
                <Input placeholder="请输入标题" allowClear />
              </Form.Item>
            </Col>
            <Col span={6}>
              <Form.Item name="type" label="类型">
                <Select options={noticeTypeOptions} placeholder="请选择类型" allowClear />
              </Form.Item>
            </Col>
          </>
        }
        options={{
          api: noticeListApi,
          add: {
            show: true,
            auth: ["system/notice/create"],
            func: () => editRef.current?.open("add"),
          },
          edit: {
            show: true,
            auth: ["system/notice/update"],
            func: (record: Record<string, unknown>) => editRef.current?.open("edit", record),
          },
          delete: {
            show: true,
            auth: ["system/notice/destroy"],
            func: async (record: { id: number }) => {
              const res = await noticeDeleteApi(record.id);
              if (res.code === 0) {
                message.success("删除成功");
                tableRef.current?.refresh();
              }
            },
          },
        }}
        columns={
          [
            { title: "标题", dataIndex: "title" },
            {
              title: "类型",
              dataIndex: "type",
              width: 100,
              render: (value: number) =>
                value === 1 ? <Tag color="blue">通知</Tag> : <Tag color="orange">公告</Tag>,
            },
            { title: "状态", dataIndex: "status", width: 100, type: "dict", dict: "status" },
            { title: "备注", dataIndex: "remark", width: 200 },
            {
              title: "创建时间",
              dataIndex: "createTime",
              width: 180,
              render: (text: string) => (text ? dayjs(text).format("YYYY-MM-DD HH:mm:ss") : "-"),
            },
          ] as ColumnDef[]
        }
      />
      <NoticeEdit ref={editRef} onSuccess={() => tableRef.current?.refresh()} />
    </>
  );
};

export default NoticeIndex;
