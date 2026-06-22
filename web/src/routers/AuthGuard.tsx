import React from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useAuthStore } from '@/store/auth';

interface AuthGuardProps {
    children: React.ReactNode;
}

const AuthGuard: React.FC<AuthGuardProps> = ({ children }) => {
    const token = useAuthStore((state) => state.token);
    const location = useLocation();

    // 没有 token 直接去登录页
    if (!token) {
        return <Navigate to="/login" state={{ from: location }} replace />;
    }

    // 验证通过，渲染子组件
    return <>{children}</>;
};

export default AuthGuard;
