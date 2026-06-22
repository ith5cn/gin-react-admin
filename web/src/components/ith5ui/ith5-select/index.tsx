import { Select } from 'antd';
import type { SelectProps } from 'antd';
import { useEffect } from 'react';
import useDictStore from '@/store/useDictStore';

interface Ith5SelectProps extends Omit<SelectProps, 'options'> {
    dict: string;
}

const Ith5Select: React.FC<Ith5SelectProps> = ({ dict, ...rest }) => {
    const { dictMap, isLoaded, fetchDictAll } = useDictStore();

    useEffect(() => {
        if (!isLoaded) fetchDictAll();
    }, [isLoaded, fetchDictAll]);

    const options = (dictMap[dict] ?? []).map((item) => ({
        label: item.label,
        value: item.value,
    }));

    return <Select allowClear options={options} {...rest} />;
};

export default Ith5Select;
