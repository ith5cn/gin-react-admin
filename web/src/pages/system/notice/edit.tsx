import { forwardRef, useImperativeHandle, useState } from "react";
import { Col, Form, Input, message, Modal, Radio, Row } from "antd";

import { noticeCreateApi, noticeUpdateApi } from "@/api/system/notice";

export interface NoticeEditRef {
  open: (type?: "add" | "edit", data?: Record<string, unknown>) => void;
}

interface NoticeEditProps {
  onSuccess?: () => void;
}

const initialFormData = {
  title: "",
  type: 1,
  content: "",
  status: 1,
  remark: "",
};

const NoticeEdit = forwardRef<NoticeEditRef, NoticeEditProps>(({ onSuccess }, ref) => {
  const [visible, setVisible] = useState(false);
  const [mode, setMode] = useState<"add" | "edit">("add");
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);

  const title = "通知公告" + (mode === "edit" ? " - 编辑" : " - 新增");

  const open = (type: "add" | "edit" = "add", data?: Record<string, unknown>) => {
    setMode(type);
    form.resetFields();
    if (type === "edit" && data) {
      form.setFieldsValue(data);
    } else {
      form.setFieldsValue({ ...initialFormData });
    }
    setVisible(true);
  };

  const close = () => setVisible(false);

  const handleSubmit = async () => {
    try {
      setLoading(true);
      const values = await form.validateFields();
      if (mode === "add") {
        await noticeCreateApi(values);
      } else {
        await noticeUpdateApi(values.id, values);
      }
      message.success("操作成功");
      onSuccess?.();
      close();
    } catch (error) {
      if ((error as { errorFields?: unknown })?.errorFields) return;
    } finally {
      setLoading(false);
    }
  };

  useImperativeHandle(ref, () => ({ open }));

  return (
    <Modal
      open={visible}
      title={title}
      confirmLoading={loading}
      width={700}
      onOk={handleSubmit}
      onCancel={close}
      destroyOnClose
    >
      <Form form={form} labelCol={{ span: 4 }} wrapperCol={{ span: 20 }}>
        <Form.Item name="id" hidden>
          <Input />
        </Form.Item>
        <Row>
          <Col span={24}>
            <Form.Item name="title" label="标题" rules={[{ required: true, message: "请输入公告标题" }]}>
              <Input placeholder="请输入公告标题" maxLength={255} />
            </Form.Item>
          </Col>
          <Col span={24}>
            <Form.Item name="type" label="类型" rules={[{ required: true, message: "请选择类型" }]}>
              <Radio.Group
                options={[
                  { label: "通知", value: 1 },
                  { label: "公告", value: 2 },
                ]}
              />
            </Form.Item>
          </Col>
          <Col span={24}>
            <Form.Item name="status" label="状态">
              <Radio.Group
                options={[
                  { label: "正常", value: 1 },
                  { label: "停用", value: 2 },
                ]}
              />
            </Form.Item>
          </Col>
          <Col span={24}>
            <Form.Item name="content" label="内容" rules={[{ required: true, message: "请输入公告内容" }]}>
              <Input.TextArea rows={8} placeholder="请输入公告内容" />
            </Form.Item>
          </Col>
          <Col span={24}>
            <Form.Item name="remark" label="备注">
              <Input.TextArea rows={2} placeholder="请输入备注" />
            </Form.Item>
          </Col>
        </Row>
      </Form>
    </Modal>
  );
});

NoticeEdit.displayName = "NoticeEdit";

export default NoticeEdit;
