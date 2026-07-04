import { useEffect, useState } from "react";
import { Card, Col, Form, Input, message, Row, Button } from "antd";

import { changePasswordApi, updateProfileApi } from "@/api/system/user";
import ImageUpload from "@/components/image-upload";
import { useAuthStore } from "@/store/auth";

type ProfileFormValues = {
  nickname?: string;
  phone?: string;
  email?: string;
  avatar?: string;
  signed?: string;
};

type PasswordFormValues = {
  oldPassword: string;
  newPassword: string;
  confirmPassword: string;
};

const ProfileIndex = () => {
  const userInfo = useAuthStore((state) => state.userInfo);
  const initUserContext = useAuthStore((state) => state.initUserContext);
  const [profileForm] = Form.useForm<ProfileFormValues>();
  const [passwordForm] = Form.useForm<PasswordFormValues>();
  const [savingProfile, setSavingProfile] = useState(false);
  const [savingPassword, setSavingPassword] = useState(false);

  useEffect(() => {
    if (userInfo) {
      profileForm.setFieldsValue({
        nickname: userInfo.nickname ?? "",
        phone: userInfo.phone ?? "",
        email: userInfo.email ?? "",
        avatar: userInfo.avatar ?? "",
        signed: userInfo.signed ?? "",
      });
    }
  }, [userInfo, profileForm]);

  const handleSaveProfile = async () => {
    try {
      setSavingProfile(true);
      const values = await profileForm.validateFields();
      await updateProfileApi(values);
      message.success("资料已更新");
      // 重新拉取用户上下文，让头像昵称等在整个界面同步生效。
      await initUserContext();
    } catch (error) {
      if ((error as { errorFields?: unknown })?.errorFields) return;
    } finally {
      setSavingProfile(false);
    }
  };

  const handleChangePassword = async () => {
    try {
      setSavingPassword(true);
      const values = await passwordForm.validateFields();
      await changePasswordApi({
        oldPassword: values.oldPassword,
        newPassword: values.newPassword,
      });
      message.success("密码修改成功");
      passwordForm.resetFields();
    } catch (error) {
      if ((error as { errorFields?: unknown })?.errorFields) return;
    } finally {
      setSavingPassword(false);
    }
  };

  return (
    <Row gutter={16}>
      <Col xs={24} lg={14}>
        <Card title="基本资料">
          <Form form={profileForm} labelCol={{ span: 4 }} wrapperCol={{ span: 18 }}>
            <Form.Item name="avatar" label="头像">
              <ImageUpload />
            </Form.Item>
            <Form.Item label="用户名">
              <Input value={userInfo?.username} disabled />
            </Form.Item>
            <Form.Item name="nickname" label="昵称">
              <Input placeholder="请输入昵称" maxLength={30} />
            </Form.Item>
            <Form.Item
              name="phone"
              label="手机号"
              rules={[{ pattern: /^1\d{10}$/, message: "手机号格式不正确" }]}
            >
              <Input placeholder="请输入手机号" />
            </Form.Item>
            <Form.Item name="email" label="邮箱" rules={[{ type: "email", message: "邮箱格式不正确" }]}>
              <Input placeholder="请输入邮箱" />
            </Form.Item>
            <Form.Item name="signed" label="个人签名">
              <Input.TextArea rows={3} placeholder="写点什么介绍自己" maxLength={200} />
            </Form.Item>
            <Form.Item wrapperCol={{ offset: 4 }}>
              <Button type="primary" loading={savingProfile} onClick={handleSaveProfile}>
                保存资料
              </Button>
            </Form.Item>
          </Form>
        </Card>
      </Col>
      <Col xs={24} lg={10}>
        <Card title="修改密码">
          <Form form={passwordForm} labelCol={{ span: 6 }} wrapperCol={{ span: 16 }}>
            <Form.Item
              name="oldPassword"
              label="原密码"
              rules={[{ required: true, message: "请输入原密码" }]}
            >
              <Input.Password placeholder="请输入原密码" />
            </Form.Item>
            <Form.Item
              name="newPassword"
              label="新密码"
              rules={[
                { required: true, message: "请输入新密码" },
                { min: 6, message: "密码至少 6 位" },
              ]}
            >
              <Input.Password placeholder="请输入新密码（至少 6 位）" />
            </Form.Item>
            <Form.Item
              name="confirmPassword"
              label="确认新密码"
              dependencies={["newPassword"]}
              rules={[
                { required: true, message: "请再次输入新密码" },
                ({ getFieldValue }) => ({
                  validator(_, value) {
                    if (!value || getFieldValue("newPassword") === value) {
                      return Promise.resolve();
                    }
                    return Promise.reject(new Error("两次输入的密码不一致"));
                  },
                }),
              ]}
            >
              <Input.Password placeholder="请再次输入新密码" />
            </Form.Item>
            <Form.Item wrapperCol={{ offset: 6 }}>
              <Button type="primary" loading={savingPassword} onClick={handleChangePassword}>
                修改密码
              </Button>
            </Form.Item>
          </Form>
        </Card>
      </Col>
    </Row>
  );
};

export default ProfileIndex;
