'use client';

import React, { memo, useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useNavbarActions } from '@/shared/contexts/navbar-context';
import { useInvoicingBusiness } from '@/shared/contexts/invoicing-business-context';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useToast } from '@/shared/providers/toast-provider';
import { SuperAdminBusinessSelector } from './super-admin-business-selector';
import { MyIntegrationsButton } from './my-integrations-button';
import {
    getConfigsAction,
    enableConfigAction,
    disableConfigAction,
} from '@/services/modules/invoicing/infra/actions';
import type { InvoicingConfig } from '@/services/modules/invoicing/domain/types';

export const InvoicingSubNavbar = memo(function InvoicingSubNavbar() {
    const pathname = usePathname();
    const { actionButtons } = useNavbarActions();
    const { selectedBusinessId, setSelectedBusinessId } = useInvoicingBusiness();
    const { permissions, isSuperAdmin } = usePermissions();
    const { showToast } = useToast();
    const [config, setConfig] = useState<InvoicingConfig | null>(null);
    const [toggling, setToggling] = useState(false);

    const isInModule = pathname.startsWith('/invoicing');

    const loadConfig = useCallback(async () => {
        try {
            const effectiveBusinessId = isSuperAdmin
                ? (selectedBusinessId ?? undefined)
                : (permissions?.business_id || undefined);
            if (!effectiveBusinessId) {
                setConfig(null);
                return;
            }
            const response = await getConfigsAction({ business_id: effectiveBusinessId });
            const configs = response.data || [];
            setConfig(configs.length > 0 ? configs[0] : null);
        } catch {
            setConfig(null);
        }
    }, [isSuperAdmin, selectedBusinessId, permissions?.business_id]);

    useEffect(() => {
        if (isInModule) {
            loadConfig();
        }
    }, [isInModule, loadConfig]);

    const handleToggle = async () => {
        if (!config || toggling) return;
        setToggling(true);
        const wasEnabled = config.enabled;
        setConfig(prev => prev ? { ...prev, enabled: !prev.enabled } : prev);
        try {
            if (wasEnabled) {
                await disableConfigAction(config.id);
                showToast('Facturacion desactivada', 'success');
            } else {
                await enableConfigAction(config.id);
                showToast('Facturacion activada', 'success');
            }
        } catch (error: any) {
            setConfig(prev => prev ? { ...prev, enabled: wasEnabled } : prev);
            showToast('Error al cambiar estado: ' + error.message, 'error');
        } finally {
            setToggling(false);
        }
    };

    if (!isInModule) {
        return null;
    }

    const isActive = (path: string) => pathname === path || pathname.startsWith(path + '/');

    const menuItems = [
        { href: '/invoicing/invoices', label: 'Facturas', icon: '🧾' },
    ];

    return (
        <div className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 shadow-sm sticky top-0 z-40">
            <div className="px-4 sm:px-6 lg:px-8 py-2">
                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3 flex-wrap">
                        {menuItems.map((item) => (
                            <Link
                                key={item.href}
                                href={item.href}
                                className={`px-4 py-3 text-base font-medium whitespace-nowrap transition-all rounded-lg flex items-center gap-3 ${
                                    isActive(item.href)
                                        ? 'bg-purple-200 dark:bg-purple-900/50 text-purple-900 dark:text-purple-200'
                                        : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 hover:text-gray-900 dark:hover:text-gray-100 hover:shadow-md hover:scale-105'
                                }`}
                            >
                                <span>{item.icon}</span>
                                {item.label}
                            </Link>
                        ))}
                        {config && (
                            <button
                                onClick={handleToggle}
                                disabled={toggling}
                                className={`px-4 py-2 text-sm font-semibold rounded-lg transition-all duration-200 ${
                                    config.enabled
                                        ? 'bg-green-500 hover:bg-green-600 text-white'
                                        : 'bg-red-500 hover:bg-red-600 text-white'
                                } ${toggling ? 'opacity-50 cursor-not-allowed' : 'hover:shadow-lg hover:scale-105'}`}
                                title={config.enabled ? 'Facturacion activa - clic para desactivar' : 'Facturacion inactiva - clic para activar'}
                            >
                                {config.enabled ? 'Facturacion Activa' : 'Facturacion Inactiva'}
                            </button>
                        )}
                    </div>
                    <div className="flex items-center gap-2 ml-4">
                        <MyIntegrationsButton businessId={selectedBusinessId} />
                        <SuperAdminBusinessSelector
                            value={selectedBusinessId}
                            onChange={setSelectedBusinessId}
                            variant="navbar"
                        />
                        {actionButtons}
                    </div>
                </div>
            </div>
        </div>
    );
});
