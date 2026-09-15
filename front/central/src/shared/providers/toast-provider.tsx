'use client';

import React, { createContext, useContext, useState, useCallback, useEffect } from 'react';
import { createPortal } from 'react-dom';
import { Toast, ToastType } from '../ui/toast';
import { isAccessDeniedMessage, notifyAccessDenied } from '../utils/access-denied';

interface ToastMessage {
    id: string;
    message: string;
    type: ToastType;
    duration?: number;
}

interface ToastContextType {
    showToast: (message: string, type: ToastType, duration?: number) => void;
}

const ToastContext = createContext<ToastContextType | undefined>(undefined);

export const ToastProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [toasts, setToasts] = useState<ToastMessage[]>([]);
    const [mounted, setMounted] = useState(false);

    useEffect(() => {
        setMounted(true);
    }, []);

    const showToast = useCallback((message: string, type: ToastType, duration = 6000) => {
        if (type === 'error') {
            const denied = isAccessDeniedMessage(message);
            if (denied) {
                notifyAccessDenied(denied);
                return;
            }
        }
        const id = Math.random().toString(36).substr(2, 9);
        setToasts((prev) => [...prev, { id, message, type, duration }]);
    }, []);

    const removeToast = useCallback((id: string) => {
        setToasts((prev) => prev.filter((toast) => toast.id !== id));
    }, []);

    return (
        <ToastContext.Provider value={{ showToast }}>
            {children}
            {mounted &&
                createPortal(
                    <div
                        className="fixed top-4 right-4 flex flex-col gap-2 pointer-events-none"
                        style={{ zIndex: 1000 }}
                    >
                        {toasts.map((toast) => (
                            <div key={toast.id} className="pointer-events-auto">
                                <Toast
                                    id={toast.id}
                                    message={toast.message}
                                    type={toast.type}
                                    duration={toast.duration}
                                    onClose={removeToast}
                                />
                            </div>
                        ))}
                    </div>,
                    document.body,
                )}
        </ToastContext.Provider>
    );
};

export const useToast = () => {
    const context = useContext(ToastContext);
    if (context === undefined) {
        throw new Error('useToast must be used within a ToastProvider');
    }
    return context;
};
