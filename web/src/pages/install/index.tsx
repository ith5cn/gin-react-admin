import { useEffect, useMemo, useState } from 'react'
import { Button, Checkbox, Form, Input, InputNumber, Result, Select, Steps, message } from 'antd'
import { CheckCircle, Database, KeyRound, LoaderCircle, Server, ShieldCheck } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import { installCheckApi, installRunApi, installStatusApi, type InstallPayload } from '@/api/install'
import { useInstallBootstrap } from '@/routers/installBootstrapContext'

const defaultValues: InstallPayload = {
  mysql: {
    host: '127.0.0.1',
    port: '3306',
    user: 'root',
    password: '',
    dbname: 'ai_system',
  },
  redis: {
    mode: 'single',
    addr: '127.0.0.1:6379',
    addrs: '',
    password: '',
    db: 0,
  },
  sqlFiles: [],
  jwtSecret: 'gin-react-admin-change-me',
}

const Install = () => {
  const [form] = Form.useForm<InstallPayload>()
  const navigate = useNavigate()
  const { markInstalled } = useInstallBootstrap()
  const [current, setCurrent] = useState(0)
  const [loading, setLoading] = useState(true)
  const [checking, setChecking] = useState(false)
  const [installing, setInstalling] = useState(false)
  const [installed, setInstalled] = useState(false)
  const [sqlFiles, setSqlFiles] = useState<string[]>([])

  useEffect(() => {
    installStatusApi()
      .then((res) => {
        const files = res.data.sqlFiles || []
        setInstalled(res.data.installed)
        setSqlFiles(files)
        form.setFieldsValue({ ...defaultValues, sqlFiles: files })
      })
      .finally(() => setLoading(false))
  }, [form])

  const steps = useMemo(
    () => [
      { title: '连接配置', icon: <Server size={16} /> },
      { title: '数据导入', icon: <Database size={16} /> },
      { title: '安全设置', icon: <KeyRound size={16} /> },
      { title: '完成安装', icon: <ShieldCheck size={16} /> },
    ],
    [],
  )

  const checkConnection = async () => {
    const values = await form.validateFields(['mysql', 'redis'])
    setChecking(true)
    try {
      await installCheckApi(values)
      message.success('数据库和 Redis 连接正常')
      setCurrent(1)
    } finally {
      setChecking(false)
    }
  }

  const runInstall = async () => {
    const values = await form.validateFields()
    const payload: InstallPayload = {
      ...defaultValues,
      ...values,
      mysql: { ...defaultValues.mysql, ...values.mysql },
      redis: { ...defaultValues.redis, ...values.redis },
      sqlFiles: values.sqlFiles || [],
      jwtSecret: values.jwtSecret || defaultValues.jwtSecret,
    }
    setInstalling(true)
    try {
      await installRunApi(payload)
      markInstalled()
      setInstalled(true)
      setCurrent(3)
      message.success('安装完成')
    } finally {
      setInstalling(false)
    }
  }

  if (installed && !installing) {
    return (
      <div className="min-h-screen bg-[#f6f7f9] px-4 py-10">
        <Result
          status="success"
          title="系统已完成安装"
          subTitle="安装锁已写入，后端会按当前配置启动业务服务。"
          extra={<Button type="primary" onClick={() => navigate('/login')}>进入登录</Button>}
        />
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-[#f6f7f9] px-4 py-8 text-[#1f2328]">
      <div className="mx-auto max-w-[1040px]">
        <div className="mb-6 flex items-center justify-between border-b border-[#d8dee4] pb-4">
          <div>
            <h1 className="m-0 text-[26px] font-semibold leading-tight">系统安装向导</h1>
            <p className="mt-2 text-sm text-[#667085]">配置运行环境并导入初始化 SQL。</p>
          </div>
          {loading ? <LoaderCircle className="animate-spin text-[#1677ff]" size={24} /> : <CheckCircle className="text-[#1f8f4d]" size={24} />}
        </div>

        <div className="mb-6 rounded border border-[#d8dee4] bg-white p-4">
          <Steps current={current} items={steps} />
        </div>

        <Form form={form} layout="vertical" initialValues={defaultValues} disabled={loading || installing}>
          <div className={current === 0 ? 'grid gap-4 lg:grid-cols-2' : 'hidden'}>
              <section className="rounded border border-[#d8dee4] bg-white p-5">
                <h2 className="mb-4 text-base font-semibold">MySQL</h2>
                <div className="grid grid-cols-3 gap-3">
                  <Form.Item name={['mysql', 'host']} label="主机" rules={[{ required: true }]}>
                    <Input />
                  </Form.Item>
                  <Form.Item name={['mysql', 'port']} label="端口" rules={[{ required: true }]}>
                    <Input />
                  </Form.Item>
                  <Form.Item name={['mysql', 'dbname']} label="数据库" rules={[{ required: true }]}>
                    <Input />
                  </Form.Item>
                </div>
                <Form.Item name={['mysql', 'user']} label="用户名" rules={[{ required: true }]}>
                  <Input />
                </Form.Item>
                <Form.Item name={['mysql', 'password']} label="密码">
                  <Input.Password />
                </Form.Item>
              </section>

              <section className="rounded border border-[#d8dee4] bg-white p-5">
                <h2 className="mb-4 text-base font-semibold">Redis</h2>
                <div className="grid grid-cols-3 gap-3">
                  <Form.Item name={['redis', 'mode']} label="模式">
                    <Select options={[{ label: 'single', value: 'single' }, { label: 'cluster', value: 'cluster' }]} />
                  </Form.Item>
                  <Form.Item name={['redis', 'addr']} label="单机地址">
                    <Input />
                  </Form.Item>
                  <Form.Item name={['redis', 'db']} label="DB">
                    <InputNumber className="w-full" min={0} max={15} />
                  </Form.Item>
                </div>
                <Form.Item name={['redis', 'addrs']} label="集群地址">
                  <Input />
                </Form.Item>
                <Form.Item name={['redis', 'password']} label="密码">
                  <Input.Password />
                </Form.Item>
              </section>
          </div>

          <section className={current === 1 ? 'rounded border border-[#d8dee4] bg-white p-5' : 'hidden'}>
              <h2 className="mb-4 text-base font-semibold">初始化 SQL</h2>
              <Form.Item name="sqlFiles" rules={[{ required: true, message: '请选择至少一个 SQL 文件' }]}>
                <Checkbox.Group className="grid gap-3">
                  {sqlFiles.map((file) => (
                    <Checkbox key={file} value={file}>{file}</Checkbox>
                  ))}
                </Checkbox.Group>
              </Form.Item>
              {!sqlFiles.length && <div className="rounded border border-[#f0d98c] bg-[#fff8db] p-3 text-sm text-[#7a5d00]">未扫描到 SQL 文件，请把导出文件放到 server/sql 或 server/database 目录。</div>}
          </section>

          <section className={current === 2 ? 'rounded border border-[#d8dee4] bg-white p-5' : 'hidden'}>
              <h2 className="mb-4 text-base font-semibold">安全设置</h2>
              <Form.Item name="jwtSecret" label="JWT Secret" rules={[{ required: true }]}>
                <Input.Password />
              </Form.Item>
          </section>
        </Form>

        <div className="mt-6 flex justify-end gap-3">
          {current > 0 && current < 3 && <Button onClick={() => setCurrent(current - 1)}>上一步</Button>}
          {current === 0 && <Button type="primary" loading={checking} onClick={checkConnection}>检测连接</Button>}
          {current === 1 && <Button type="primary" onClick={() => setCurrent(2)} disabled={!sqlFiles.length}>下一步</Button>}
          {current === 2 && <Button type="primary" loading={installing} onClick={runInstall}>开始安装</Button>}
        </div>
      </div>
    </div>
  )
}

export default Install
