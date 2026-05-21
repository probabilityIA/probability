'use client';

import { useCallback, useEffect, useState } from 'react';
import { Modal } from '@/shared/ui';
import { ClientGroup } from '../../domain/types';
import { listClientGroupsAction } from '../../infra/actions';
import { ClientGroupsPanel } from './ClientGroupsPanel';
import { CatalogPriceTable } from './CatalogPriceTable';

interface CatalogPricingModalProps {
    isOpen: boolean;
    onClose: () => void;
    businessId?: number;
}

type Tab = 'groups' | 'prices';

export function CatalogPricingModal({ isOpen, onClose, businessId }: CatalogPricingModalProps) {
    const [activeTab, setActiveTab] = useState<Tab>('groups');
    const [groups, setGroups] = useState<ClientGroup[]>([]);
    const [loadingGroups, setLoadingGroups] = useState(false);

    const refreshGroups = useCallback(async () => {
        setLoadingGroups(true);
        const result = await listClientGroupsAction(businessId, '', 1);
        setGroups(result.data);
        setLoadingGroups(false);
    }, [businessId]);

    useEffect(() => {
        if (isOpen) {
            setActiveTab('groups');
            refreshGroups();
        }
    }, [isOpen, refreshGroups]);

    const tabClass = (tab: Tab) =>
        `px-5 py-2.5 text-sm font-bold rounded-lg transition-all duration-200 ${
            activeTab === tab
                ? 'btn-business-primary text-white'
                : 'text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'
        }`;

    return (
        <Modal isOpen={isOpen} onClose={onClose} title="Grupos de clientes para catalogos" size="4xl">
            <div className="p-4 space-y-4">
                <p className="text-sm text-gray-500 dark:text-gray-400">
                    Agrupa tus clientes en tipos y define un precio por producto para cada grupo o cliente.
                    El precio base es el del producto; subelo o bajalo por grupo.
                </p>

                <div className="flex gap-2 bg-gray-100 dark:bg-gray-800 p-1.5 rounded-xl w-fit">
                    <button onClick={() => setActiveTab('groups')} className={tabClass('groups')}>
                        Grupos y clientes
                    </button>
                    <button onClick={() => setActiveTab('prices')} className={tabClass('prices')}>
                        Precios del catalogo
                    </button>
                </div>

                {activeTab === 'groups' ? (
                    <ClientGroupsPanel
                        businessId={businessId}
                        groups={groups}
                        loading={loadingGroups}
                        onGroupsChanged={refreshGroups}
                    />
                ) : (
                    <CatalogPriceTable businessId={businessId} groups={groups} />
                )}
            </div>
        </Modal>
    );
}
