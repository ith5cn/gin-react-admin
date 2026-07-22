// 本文件由代码生成器生成，重新生成会覆盖手工修改。
import { useRef } from "react";
import { Col, Form, InputNumber, message } from "antd";
import Ith5Table, { type TableRef } from "@/components/ith5ui/ith5-table";
import { aiArticleDeleteApi, aiArticleListApi } from "@/api/system/ai-article";
import AiarticleEdit, { type AiarticleEditRef } from "./edit";

const AiarticleIndex = () => {
  const tableRef = useRef<TableRef>(null);
  const editRef = useRef<AiarticleEditRef>(null);

  return (
    <>
      <Ith5Table
        ref={tableRef}
        searchFields={
          <>
            <Col span={6}>
              <Form.Item name="status" label="状态">
                <InputNumber style={{ width: "100%" }} placeholder="请输入状态" />
              </Form.Item>
            </Col>
          </>
        }
        options={{
          api: aiArticleListApi,
          add: { show: true, auth: ["system/ai-article/create"], func: () => editRef.current?.open("add") },
          edit: { show: true, auth: ["system/ai-article/update"], func: (record: { id: number }) => editRef.current?.open("edit", record) },
          delete: {
            show: true,
            auth: ["system/ai-article/destroy"],
            func: async (record: { id: number }) => {
              await aiArticleDeleteApi(record.id);
              message.success("删除成功");
              tableRef.current?.refresh();
            },
          },
        }}
        columns={[
          { title: "文章标题", dataIndex: "title", key: "title" },
          { title: "排序", dataIndex: "sort", key: "sort" },
          { title: "状态", dataIndex: "status", key: "status" }
        ]}
      />
      <AiarticleEdit ref={editRef} onSuccess={() => tableRef.current?.refresh()} />
    </>
  );
};

export default AiarticleIndex;
