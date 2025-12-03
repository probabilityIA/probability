import React from 'react';
import { Table } from '@/shared/ui/table';
import { Button } from '@/shared/ui/button';
import { Spinner } from '@/shared/ui/spinner';
import { Alert } from '@/shared/ui/alert';
import { ConfiguredResource } from '../../domain/types';
import { useResourceConfig } from '../hooks/useResourceConfig';

interface ResourceConfigProps {
    businessId: number;
}

export const ResourceConfig: React.FC<ResourceConfigProps> = ({ businessId }) => {
    const {
        config,
        loading,
        error,
        actionLoading,
        toggleResource,
        setError
    } = useResourceConfig(businessId);

    const columns = [
        { label: 'ID', key: 'resource_id' },
        { label: 'Resource Name', key: 'resource_name' },
        {
            label: 'Status',
            key: 'is_active',
            render: (_: unknown, row: ConfiguredResource) => (
                <span className={`px-2 py-1 rounded text-xs ${row.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}`}>
                    {row.is_active ? 'Active' : 'Inactive'}
                </span>
            )
        },
        {
            label: 'Actions',
            key: 'actions',
            render: (_: unknown, row: ConfiguredResource) => (
                <Button
                    size="sm"
                    variant={row.is_active ? 'danger' : 'primary'}
                    onClick={() => toggleResource(row)}
                    disabled={actionLoading === row.resource_id}
                >
                    {actionLoading === row.resource_id ? <Spinner size="sm" /> : (row.is_active ? 'Deactivate' : 'Activate')}
                </Button>
            )
        }
    ];

    if (loading) return <Spinner />;
    if (error) return <Alert type="error" onClose={() => setError(null)}>{error}</Alert>;
    if (!config) return <div>No configuration found.</div>;

    return (
        <div className="space-y-4">
            <div className="flex justify-between items-center">
                <h2 className="text-xl font-semibold">Configured Resources</h2>
                <div className="text-sm text-gray-500">
                    Total: {config.total} | Active: {config.active} | Inactive: {config.inactive}
                </div>
            </div>
            <Table
                data={config.resources}
                columns={columns}
                keyExtractor={(item) => item.resource_id}
            />
        </div>
    );
};
