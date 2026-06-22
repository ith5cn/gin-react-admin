import { useEffect, useMemo, useState } from 'react'
import type { Key } from 'react'
import { Button, Drawer, Input, Modal, Space, Table, Tag, Tooltip, message } from 'antd'
import type { ColumnsType } from 'antd/es/table'
import { Database, Eye, RefreshCcw, RotateCcw, Trash2, Wrench } from 'lucide-react'
import {
  databaseClearFragmentApi,
  databaseOptimizeApi,
  databaseRecycleDestroyApi,
  databaseRecycleListApi,
  databaseRecycleRecoverApi,
  databaseTableColumnsApi,
  databaseTableListApi,
  type DatabaseColumnRecord,
  type DatabaseTableRecord,
  type RecycleRecord,
} from '@/api/data/database'

const formatBytes = (value?: number) => {
  const size = Number(value || 0)
  if (size <= 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const index = Math.min(Math.floor(Math.log(size) / Math.log(1024)), units.length - 1)
  return `${(size / 1024 ** index).toFixed(index === 0 ? 0 : 2)} ${units[index]}`
}

const DatabaseIndex = () => {
  const [loading, setLoading] = useState(false)
  const [tableName, setTableName] = useState('')
  const [rows, setRows] = useState<DatabaseTableRecord[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [limit, setLimit] = useState(10)
  const [selectedRowKeys, setSelectedRowKeys] = useState<Key[]>([])

  const [columnsOpen, setColumnsOpen] = useState(false)
  const [columnsLoading, setColumnsLoading] = useState(false)
  const [currentTable, setCurrentTable] = useState<DatabaseTableRecord | null>(null)
  const [columnRows, setColumnRows] = useState<DatabaseColumnRecord[]>([])

  const [recycleOpen, setRecycleOpen] = useState(false)
  const [recycleLoading, setRecycleLoading] = useState(false)
  const [recycleRows, setRecycleRows] = useState<RecycleRecord[]>([])
  const [recycleTotal, setRecycleTotal] = useState(0)
  const [recyclePage, setRecyclePage] = useState(1)
  const [recycleLimit, setRecycleLimit] = useState(10)
  const [recycleTable, setRecycleTable] = useState<DatabaseTableRecord | null>(null)
  const [selectedRecycleKeys, setSelectedRecycleKeys] = useState<Key[]>([])

  const selectedTables = useMemo(() => selectedRowKeys.map(String), [selectedRowKeys])

  const fetchTables = async (nextPage = page, nextLimit = limit) => {
    setLoading(true)
    try {
      const res = await databaseTableListApi({
        page: nextPage,
        limit: nextLimit,
        tableName,
      })
      setRows(res.data.list || [])
      setTotal(res.data.total || 0)
      setPage(nextPage)
      setLimit(nextLimit)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchTables(1, limit)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const requireSelectedTables = () => {
    if (!selectedTables.length) {
      message.warning('请先选择数据表')
      return false
    }
    return true
  }

  const confirmMaintain = (title: string, api: (tables: string[]) => Promise<any>) => {
    if (!requireSelectedTables()) return
    Modal.confirm({
      title,
      content: `将对 ${selectedTables.length} 张数据表执行操作，是否继续？`,
      onOk: async () => {
        await api(selectedTables)
        message.success('操作成功')
        setSelectedRowKeys([])
        fetchTables()
      },
    })
  }

  const openColumns = async (record: DatabaseTableRecord) => {
    setCurrentTable(record)
    setColumnsOpen(true)
    setColumnsLoading(true)
    try {
      const res = await databaseTableColumnsApi(record.tableName)
      setColumnRows(res.data || [])
    } finally {
      setColumnsLoading(false)
    }
  }

  const fetchRecycle = async (record = recycleTable, nextPage = recyclePage, nextLimit = recycleLimit) => {
    if (!record) return
    setRecycleLoading(true)
    try {
      const res = await databaseRecycleListApi({
        tableName: record.tableName,
        page: nextPage,
        limit: nextLimit,
      })
      setRecycleRows(res.data.list || [])
      setRecycleTotal(res.data.total || 0)
      setRecyclePage(nextPage)
      setRecycleLimit(nextLimit)
    } finally {
      setRecycleLoading(false)
    }
  }

  const openRecycle = async (record: DatabaseTableRecord) => {
    setRecycleTable(record)
    setRecycleOpen(true)
    setSelectedRecycleKeys([])
    await fetchRecycle(record, 1, recycleLimit)
  }

  const confirmRecycleAction = (type: 'recover' | 'destroy', ids: Array<string | number>) => {
    if (!recycleTable || !ids.length) {
      message.warning('请先选择回收站数据')
      return
    }
    const isDestroy = type === 'destroy'
    Modal.confirm({
      title: isDestroy ? '永久删除数据' : '恢复数据',
      content: isDestroy ? '永久删除后无法恢复，是否继续？' : '确定恢复选中的数据吗？',
      okButtonProps: { danger: isDestroy },
      onOk: async () => {
        if (isDestroy) {
          await databaseRecycleDestroyApi(recycleTable.tableName, ids)
        } else {
          await databaseRecycleRecoverApi(recycleTable.tableName, ids)
        }
        message.success('操作成功')
        setSelectedRecycleKeys([])
        fetchRecycle(recycleTable, recyclePage, recycleLimit)
      },
    })
  }

  const columns: ColumnsType<DatabaseTableRecord> = [
    { title: '表名称', dataIndex: 'tableName', width: 220, fixed: 'left' },
    { title: '表注释', dataIndex: 'tableComment', width: 220, render: (value) => value || '-' },
    { title: '表引擎', dataIndex: 'engine', width: 100 },
    { title: '数据更新时间', dataIndex: 'updateTime', width: 180, render: (value) => value || '-' },
    { title: '总行数', dataIndex: 'tableRows', width: 110 },
    {
      title: '碎片大小',
      dataIndex: 'fragmentSize',
      width: 120,
      render: (value) => <Tag color={Number(value) > 0 ? 'gold' : 'default'}>{formatBytes(value)}</Tag>,
    },
    { title: '数据大小', dataIndex: 'dataSize', width: 120, render: formatBytes },
    { title: '索引大小', dataIndex: 'indexSize', width: 120, render: formatBytes },
    { title: '字符集', dataIndex: 'tableCharset', width: 120 },
    { title: '创建时间', dataIndex: 'createTime', width: 180, render: (value) => value || '-' },
    {
      title: '操作',
      key: 'action',
      width: 220,
      fixed: 'right',
      render: (_value, record) => (
        <Space size={4}>
          <Button type="link" size="small" icon={<Eye size={14} />} onClick={() => openColumns(record)}>
            表结构
          </Button>
          <Button type="link" size="small" icon={<Trash2 size={14} />} onClick={() => openRecycle(record)}>
            回收站数据
          </Button>
        </Space>
      ),
    },
  ]

  const columnColumns: ColumnsType<DatabaseColumnRecord> = [
    { title: '字段名', dataIndex: 'columnName', width: 180 },
    { title: '字段注释', dataIndex: 'columnComment', width: 180, render: (value) => value || '-' },
    { title: '字段类型', dataIndex: 'columnType', width: 180 },
    { title: '允许为空', dataIndex: 'isNullable', width: 100 },
    { title: '键', dataIndex: 'columnKey', width: 90, render: (value) => value || '-' },
    { title: '默认值', dataIndex: 'columnDefault', width: 140, render: (value) => value ?? '-' },
    { title: '额外信息', dataIndex: 'extra', width: 180, render: (value) => value || '-' },
  ]

  const recycleColumns: ColumnsType<RecycleRecord> = [
    { title: '删除时间', dataIndex: 'deleteTime', width: 180, render: (value) => value || '-' },
    {
      title: '数据内容',
      dataIndex: 'content',
      render: (value) => (
        <pre className="m-0 max-h-[120px] overflow-auto rounded bg-[#f6f8fa] p-2 text-xs leading-5">
          {JSON.stringify(value, null, 2)}
        </pre>
      ),
    },
    {
      title: '操作',
      key: 'action',
      width: 180,
      fixed: 'right',
      render: (_value, record) => (
        <Space size={4}>
          <Button type="link" size="small" icon={<RotateCcw size={14} />} onClick={() => confirmRecycleAction('recover', [record.id])}>
            恢复
          </Button>
          <Button type="link" size="small" danger icon={<Trash2 size={14} />} onClick={() => confirmRecycleAction('destroy', [record.id])}>
            永久删除
          </Button>
        </Space>
      ),
    },
  ]

  return (
    <div className="p-4">
      <div className="mb-4 flex flex-wrap items-center justify-between gap-3">
        <Space>
          <Input.Search
            allowClear
            placeholder="搜索表名称"
            value={tableName}
            onChange={(event) => setTableName(event.target.value)}
            onSearch={() => fetchTables(1, limit)}
            style={{ width: 260 }}
          />
          <Tooltip title="刷新">
            <Button icon={<RefreshCcw size={16} />} onClick={() => fetchTables()} />
          </Tooltip>
        </Space>
        <Space>
          <Button icon={<Database size={16} />} onClick={() => confirmMaintain('清理表碎片', databaseClearFragmentApi)}>
            清理碎片
          </Button>
          <Button type="primary" icon={<Wrench size={16} />} onClick={() => confirmMaintain('优化数据表', databaseOptimizeApi)}>
            优化表
          </Button>
        </Space>
      </div>

      <Table
        rowKey="tableName"
        bordered
        size="small"
        loading={loading}
        columns={columns}
        dataSource={rows}
        scroll={{ x: 'max-content' }}
        rowSelection={{
          selectedRowKeys,
          onChange: setSelectedRowKeys,
        }}
        pagination={{
          current: page,
          pageSize: limit,
          total,
          showSizeChanger: true,
          showTotal: (value) => `共 ${value} 张表`,
        }}
        onChange={(pagination) => fetchTables(pagination.current || 1, pagination.pageSize || 10)}
      />

      <Modal
        title={currentTable ? `${currentTable.tableName} 表结构` : '表结构'}
        open={columnsOpen}
        width={980}
        footer={null}
        onCancel={() => setColumnsOpen(false)}
      >
        <Table
          rowKey="columnName"
          size="small"
          bordered
          loading={columnsLoading}
          columns={columnColumns}
          dataSource={columnRows}
          pagination={false}
          scroll={{ x: 'max-content', y: 520 }}
        />
      </Modal>

      <Drawer
        title={recycleTable ? `${recycleTable.tableName} 回收站数据` : '回收站数据'}
        open={recycleOpen}
        width={860}
        onClose={() => setRecycleOpen(false)}
        extra={
          <Space>
            <Button
              icon={<RotateCcw size={16} />}
              onClick={() => confirmRecycleAction('recover', selectedRecycleKeys as Array<string | number>)}
            >
              批量恢复
            </Button>
            <Button
              danger
              icon={<Trash2 size={16} />}
              onClick={() => confirmRecycleAction('destroy', selectedRecycleKeys as Array<string | number>)}
            >
              批量永久删除
            </Button>
          </Space>
        }
      >
        <Table
          rowKey="id"
          size="small"
          bordered
          loading={recycleLoading}
          columns={recycleColumns}
          dataSource={recycleRows}
          scroll={{ x: 'max-content' }}
          rowSelection={{
            selectedRowKeys: selectedRecycleKeys,
            onChange: setSelectedRecycleKeys,
          }}
          pagination={{
            current: recyclePage,
            pageSize: recycleLimit,
            total: recycleTotal,
            showSizeChanger: true,
            showTotal: (value) => `共 ${value} 条`,
          }}
          onChange={(pagination) => fetchRecycle(recycleTable, pagination.current || 1, pagination.pageSize || 10)}
        />
      </Drawer>
    </div>
  )
}

export default DatabaseIndex
