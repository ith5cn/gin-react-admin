import React from 'react';
import { useAuthStore } from '@/store/auth';

interface Props {
    code: string;
    children: React.ReactNode;
}

export const HasPermission: React.FC<Props> = ({ code, children }) => {
    const permissionCodes = useAuthStore((state) => state.codes);

    const hasAccess = permissionCodes.includes(code) || permissionCodes.includes('*:*:*');

    return hasAccess ? <>{children}</> : null;
};
