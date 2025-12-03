import React, { useState } from 'react';
import { Table } from '@/shared/ui/table';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { Spinner } from '@/shared/ui/spinner';
import { ConfirmModal } from '@/shared/ui/confirm-modal';
import { Alert } from '@/shared/ui/alert';
import { Modal } from '@/shared/ui/modal';
import { Permission } from '../../domain/types';
import { PermissionForm } from './PermissionForm';
import { usePermissions } from '../hooks/usePermissions';

export const PermissionList: React.FC = () => {
    const {
        permissions,
        loading,
        error,
        searchName,
        setSearchName,
        filterScope,
        setFilterScope,
        filterBusinessType,
        setFilterBusinessType,
        deletePermission,
        refresh,
        setError
    } = usePermissions();

    const [showCreateModal, setShowCreateModal] = useState(false);
    const [editingPermission, setEditingPermission] = useState<Permission | null>(null);
    const [deleteId, setDeleteId] = useState<number | null>(null);

    const handleDelete = async () => {
        if (deleteId) {
            const success = await deletePermission(deleteId);
            if (success) setDeleteId(null);
        }
    };

    const handleSave = () => {
        setShowCreateModal(false);
        setEditingPermission(null);
        refresh();
    };

    const columns = [
        { label: 'ID', key: 'id' },
        { label: 'Name', key: 'name' },
        { label: 'Code', key: 'code' },
        { label: 'Resource', key: 'resource' },
        { label: 'Action', key: 'action' },
        { label: 'Scope', key: 'scope_name' },
        { label: 'Business Type', key: 'business_type_name' },
        {
            label: 'Actions',
            key: 'actions',
            render: (_: unknown, row: Permission) => (
                <div className="flex gap-2">
                    <Button variant="secondary" size="sm" onClick={() => { setEditingPermission(row); setShowCreateModal(true); }}>Edit</Button>
                    <Button variant="danger" size="sm" onClick={() => setDeleteId(row.id)}>Delete</Button>
                </div>
            ),
        },
    ];

    return (
        <div className="p-6 space-y-6">
            <div className="flex justify-between items-center">
                <h1 className="text-2xl font-bold">Permissions</h1>
                <Button onClick={() => { setEditingPermission(null); setShowCreateModal(true); }}>Create Permission</Button>
            </div>

            {error && <Alert type="error" onClose={() => setError(null)}>{error}</Alert>}

            <div className="flex gap-4 mb-4">
                <Input
                    placeholder="Search by name..."
                    value={searchName}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => setSearchName(e.target.value)}
                    className="max-w-xs"
                />
                <Input
                    placeholder="Filter by Scope ID"
                    type="number"
                    value={filterScope}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFilterScope(e.target.value)}
                    className="max-w-xs"
                />
                <Input
                    placeholder="Filter by Business Type ID"
                    type="number"
                    value={filterBusinessType}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFilterBusinessType(e.target.value)}
                    className="max-w-xs"
                />
            </div>

            {loading ? (
                <div className="flex justify-center p-8"><Spinner /></div>
            ) : (
                <Table
                    data={permissions}
                    columns={columns}
                    keyExtractor={(item) => item.id}
                />
            )}

            <Modal
                isOpen={showCreateModal}
                onClose={() => { setShowCreateModal(false); setEditingPermission(null); }}
                title={editingPermission ? "Edit Permission" : "Create Permission"}
            >
                <PermissionForm
                    initialData={editingPermission || undefined}
                    onSuccess={handleSave}
                    onCancel={() => { setShowCreateModal(false); setEditingPermission(null); }}
                />
            </Modal>

            <ConfirmModal
                isOpen={!!deleteId}
                title="Delete Permission"
                message="Are you sure you want to delete this permission? This action cannot be undone."
                onConfirm={handleDelete}
                onClose={() => setDeleteId(null)}
            />
        </div>
    );
};
