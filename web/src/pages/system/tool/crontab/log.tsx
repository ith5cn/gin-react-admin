import { crontabLogListApi } from "@/api/system/crontab";
import Ith5Table, { type TableRef } from "@/components/ith5ui/ith5-table";
import { PlayCircleOutlined } from "@ant-design/icons";
import { Button, Drawer, Modal } from "antd";
import dayjs from "dayjs";
import { forwardRef, useImperativeHandle, useRef, useState } from "react";


export interface CrontabLogRef {
    open: (data: { id: number | string }) => void;
}

// forwardRef 渲染函数的参数顺序是 (props, ref)，ref 是第二个参数
const CrontabLog = forwardRef<CrontabLogRef>((_, ref) => {
    const [visible, setVisible] = useState(false);
    const tableRef = useRef<TableRef>(null)
    const [searchForm, setSearchForm] = useState({
        crontabId: ''
    })

    // 打开弹框
    // 注意：不要在这里手动调用 refresh()
    // 使用 extraSearchParams 后，Ith5Table 内部会监听它的变化并自动刷新
    const open = (data: { id: number | string }) => {
        setSearchForm({ crontabId: String(data.id) });
        setVisible(true);
    };



    // 暴露给父组件的方法
    useImperativeHandle(ref, () => ({
        open
    }));

    const options = {
        api: crontabLogListApi
    }

    return (
        <Drawer title="任务日志" open={visible} onClose={() => setVisible(false)} width={800}>
            <Ith5Table
                ref={tableRef}
                options={options}
                extraSearchParams={searchForm}
                operationBeforeExtend={(record) => (
                    <>
                        <Button
                            type="link"
                            size="small"
                            icon={<PlayCircleOutlined />}
                            onClick={async () => {
                                Modal.info({
                                    title: '任务日志',
                                    content: record.exceptionInfo,
                                    okText: '确定'
                                })
                            }}
                        >
                            查看
                        </Button>
                    </>
                )}
                columns={[
                    {
                        title: '执行时间',
                        dataIndex: 'createTime',
                        key: 'createTime',
                        render: (text: string) => {
                            return text ? dayjs(text).format('YYYY-MM-DD HH:mm:ss') : '-';
                        }
                    },
                    {
                        title: '执行目标',
                        dataIndex: 'target',
                        key: 'target',
                    },
                    {
                        title: '执行结果',
                        dataIndex: 'status',
                        type: 'dict',
                        dict: 'result',
                        key: 'status',
                    }
                ]}
            />
        </Drawer>
    );
});

export default CrontabLog;