// 本文件由代码生成器生成，重新生成会覆盖手工修改。
import { forwardRef, useImperativeHandle, useState } from "react";
import { Checkbox, DatePicker, Form, Input, InputNumber, Modal, Radio, Select, Slider, Switch, TreeSelect, message } from "antd";
import dayjs from "dayjs";
import type { Dayjs } from "dayjs";
import WangEditor from "@/components/wang-editor";
import { aiArticleCreateApi, aiArticleUpdateApi } from "@/api/system/ai-article";

export interface AiarticleEditRef {
  open: (type?: "add" | "edit", data?: Record<string, unknown>) => void;
}

interface AiarticleEditProps {
  onSuccess?: () => void;
}

const initialFormData = {
  categoryId: undefined,
  title: "",
  author: "",
  image: [],
  describe: 2,
  content: "",
  views: undefined,
  sort: undefined,
  status: undefined,
  isLink: undefined,
  linkUrl: [],
  isHot: undefined,
};

const DATE_FIELDS: Array<{ name: string; format: string }> = [
  { name: "isHot", format: "YYYY-MM-DD" },
];

const ARRAY_FIELDS = ["image", "linkUrl"];

const toFormValues = (data: Record<string, unknown>) => {
  const next: Record<string, unknown> = { ...data };
  DATE_FIELDS.forEach(({ name }) => {
    next[name] = next[name] ? dayjs(String(next[name])) : undefined;
  });
  ARRAY_FIELDS.forEach((name) => {
    if (typeof next[name] === "string") {
      next[name] = next[name] === "" ? [] : String(next[name]).split(",");
    }
  });
  return next;
};

const toSubmitValues = (values: Record<string, unknown>) => {
  const next: Record<string, unknown> = { ...values };
  DATE_FIELDS.forEach(({ name, format }) => {
    next[name] = next[name] ? (next[name] as Dayjs).format(format) : undefined;
  });
  ARRAY_FIELDS.forEach((name) => {
    if (Array.isArray(next[name])) {
      next[name] = (next[name] as string[]).join(",");
    }
  });
  return next;
};

const AiarticleEdit = forwardRef<AiarticleEditRef, AiarticleEditProps>(({ onSuccess }, ref) => {
  const [visible, setVisible] = useState(false);
  const [mode, setMode] = useState<"add" | "edit">("add");
  const [loading, setLoading] = useState(false);
  const [form] = Form.useForm();
  const title = "测试文章菜单1" + (mode === "edit" ? " - 编辑" : " - 新增");

  const open = (type: "add" | "edit" = "add", data?: Record<string, unknown>) => {
    setMode(type);
    form.resetFields();
    form.setFieldsValue(type === "edit" && data ? toFormValues(data) : initialFormData);
    setVisible(true);
  };

  const close = () => setVisible(false);

  const handleSubmit = async () => {
    try {
      setLoading(true);
      const values = await form.validateFields();
      const payload = toSubmitValues(values);
      if (mode === "add") {
        await aiArticleCreateApi(payload);
      } else {
        await aiArticleUpdateApi(values.id as number, payload);
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
    <Modal open={visible} title={title} width={600} confirmLoading={loading} onOk={handleSubmit} onCancel={close}>
      <Form form={form} labelCol={{ span: 4 }} wrapperCol={{ span: 18 }}>
        <Form.Item name="id" hidden>
          <Input />
        </Form.Item>
        <Form.Item name="categoryId" label="分类" rules={[{ required: true, message: "请输入分类" }]}>
          <Input.Password placeholder="请输入分类" />
        </Form.Item>
        <Form.Item name="title" label="文章标题" rules={[{ required: true, message: "请输入文章标题" }]}>
          <WangEditor />
        </Form.Item>
        <Form.Item name="author" label="文章作者">
          <InputNumber style={{ width: "100%" }} placeholder="请输入文章作者" />
        </Form.Item>
        <Form.Item name="image" label="文章图片">
          <Select mode="tags" open={false} tokenSeparators={[","]} placeholder="输入后回车添加文章图片" />
        </Form.Item>
        <Form.Item name="describe" label="文章简介" getValueProps={(value) => ({ checked: value === 1 })} normalize={(value) => (value ? 1 : 2)}>
          <Switch />
        </Form.Item>
        <Form.Item name="content" label="文章内容">
          <Slider />
        </Form.Item>
        {/* 浏览次数 的数据来源需按业务补充 */}
        <Form.Item name="views" label="浏览次数">
          <Select allowClear placeholder="请选择浏览次数" options={[]} />
        </Form.Item>
        {/* 排序 未配置数据字典，选项需按业务补充 */}
        <Form.Item name="sort" label="排序">
          <Select allowClear placeholder="请选择排序" options={[]} />
        </Form.Item>
        {/* 状态 的树形数据来源需按业务补充 */}
        <Form.Item name="status" label="状态">
          <TreeSelect allowClear placeholder="请选择状态" treeData={[]} />
        </Form.Item>
        {/* 是否外链 未配置数据字典，选项需按业务补充 */}
        <Form.Item name="isLink" label="是否外链">
          <Radio.Group options={[]} />
        </Form.Item>
        {/* 链接地址 未配置数据字典，选项需按业务补充 */}
        <Form.Item name="linkUrl" label="链接地址">
          <Checkbox.Group options={[]} />
        </Form.Item>
        <Form.Item name="isHot" label="是否热门">
          <DatePicker style={{ width: "100%" }} placeholder="请选择是否热门" />
        </Form.Item>
      </Form>
    </Modal>
  );
});

AiarticleEdit.displayName = "AiarticleEdit";

export default AiarticleEdit;
