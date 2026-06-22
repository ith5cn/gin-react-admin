import { forwardRef, useImperativeHandle, useState } from "react";
import { Form, Input, InputNumber, Modal, message } from "antd";
import { aiArticleCreateApi, aiArticleUpdateApi } from "@/api/system/ai-article";

export interface AiarticleEditRef {
  open: (type?: "add" | "edit", data?: Record<string, any>) => void;
}

interface AiarticleEditProps {
  onSuccess?: () => void;
}

const initialFormData = {
  categoryId: undefined,
  title: "",
  author: "",
  image: "",
  describe: "",
  content: "",
  views: undefined,
  sort: undefined,
  status: undefined,
  isLink: undefined,
  linkUrl: "",
  isHot: undefined,
};

const AiarticleEdit = forwardRef<AiarticleEditRef, AiarticleEditProps>(({ onSuccess }, ref) => {
  const [visible, setVisible] = useState(false);
  const [mode, setMode] = useState<"add" | "edit">("add");
  const [loading, setLoading] = useState(false);
  const [form] = Form.useForm();
  const title = "测试文章菜单" + (mode === "edit" ? " - 编辑" : " - 新增");

  const open = (type: "add" | "edit" = "add", data?: Record<string, any>) => {
    setMode(type);
    form.resetFields();
    form.setFieldsValue(type === "edit" && data ? data : initialFormData);
    setVisible(true);
  };

  const close = () => setVisible(false);

  const handleSubmit = async () => {
    try {
      setLoading(true);
      const values = await form.validateFields();
      if (mode === "add") {
        await aiArticleCreateApi(values);
      } else {
        await aiArticleUpdateApi(values.id, values);
      }
      message.success("操作成功");
      onSuccess?.();
      close();
    } catch (error: any) {
      if (error?.errorFields) return;
    } finally {
      setLoading(false);
    }
  };

  useImperativeHandle(ref, () => ({ open }));

  return (
    <Modal open={visible} title={title} width={720} confirmLoading={loading} onOk={handleSubmit} onCancel={close}>
      <Form form={form} labelCol={{ span: 4 }} wrapperCol={{ span: 18 }}>
        <Form.Item name="id" hidden>
          <Input />
        </Form.Item>
        <Form.Item name="categoryId" label="分类">
          <InputNumber style={{ width: "100%" }} placeholder="请输入分类" />
        </Form.Item>
        <Form.Item name="title" label="文章标题">
          <Input placeholder="请输入文章标题" />
        </Form.Item>
        <Form.Item name="author" label="文章作者">
          <Input placeholder="请输入文章作者" />
        </Form.Item>
        <Form.Item name="image" label="文章图片">
          <Input placeholder="请输入文章图片" />
        </Form.Item>
        <Form.Item name="describe" label="文章简介">
          <Input placeholder="请输入文章简介" />
        </Form.Item>
        <Form.Item name="content" label="文章内容">
          <Input.TextArea rows={4} placeholder="请输入文章内容" />
        </Form.Item>
        <Form.Item name="views" label="浏览次数">
          <InputNumber style={{ width: "100%" }} placeholder="请输入浏览次数" />
        </Form.Item>
        <Form.Item name="sort" label="排序">
          <InputNumber style={{ width: "100%" }} placeholder="请输入排序" />
        </Form.Item>
        <Form.Item name="status" label="状态">
          <InputNumber style={{ width: "100%" }} placeholder="请输入状态" />
        </Form.Item>
        <Form.Item name="isLink" label="是否外链">
          <InputNumber style={{ width: "100%" }} placeholder="请输入是否外链" />
        </Form.Item>
        <Form.Item name="linkUrl" label="链接地址">
          <Input placeholder="请输入链接地址" />
        </Form.Item>
        <Form.Item name="isHot" label="是否热门">
          <InputNumber style={{ width: "100%" }} placeholder="请输入是否热门" />
        </Form.Item>
      </Form>
    </Modal>
  );
});

AiarticleEdit.displayName = "AiarticleEdit";

export default AiarticleEdit;
