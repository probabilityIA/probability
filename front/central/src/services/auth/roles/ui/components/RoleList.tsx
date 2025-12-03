import React, { useState } from 'react';
import { Table } from '@/shared/ui/table';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { Spinner } from '@/shared/ui/spinner';
import { ConfirmModal } from '@/shared/ui/confirm-modal';
import { Alert } from '@/shared/ui/alert';
import { Modal } from '@/shared/ui/modal';
import { Role } from '../../domain/types';
import { RoleForm } from './RoleForm';
import { useRoles } from '../hooks/useRoles';

export const RoleList: React.FC = () => {
    const {
        roles,
        loading,
        error,
        searchName,
        setSearchName,
        filterScope,
        setFilterScope,
        filterBusinessType,
        setFilterBusinessType,
        filterLevel,
        setFilterLevel,
        filterIsSystem,
        setFilterIsSystem,
        deleteRole,
        refresh,
        setError
    } = useRoles();

    const [showCreateModal, setShowCreateModal] = useState(false);
    const [editingRole, setEditingRole] = useState<Role | null>(null);
    const [deleteId, setDeleteId] = useState<number | null>(null);

    const handleDelete = async () => {
        if (deleteId) {
            const success = await deleteRole(deleteId);
            if (success) setDeleteId(null);
        }
    };

    const handleSave = () => {
        setShowCreateModal(false);
        setEditingRole(null);
        refresh();
    };

    const columns = [
        { label: 'ID', key: 'id' },
        { label: 'Name', key: 'name' },
        { label: 'Code', key: 'code' },
        { label: 'Level', key: 'level' },
        { label: 'System', key: 'is_system', render: (_: unknown, row: Role) => row.is_system ? 'Yes' : 'No' },
        { label: 'Scope', key: 'scope_name' },
        { label: 'Business Type', key: 'business_type_name' },
        {
            label: 'Actions',
            key: 'actions',
            render: (_: unknown, row: Role) => (
                <div className="flex gap-2">
                    <Button variant="secondary" size="sm" onClick={() => { setEditingRole(row); setShowCreateModal(true); }}>Edit</Button>
                    <Button variant="danger" size="sm" onClick={() => setDeleteId(row.id)}>Delete</Button>
                </div>
            ),
        },
    ];

    return (
        <div className="p-6 space-y-6">
            <div className="flex justify-between items-center">
                <h1 className="text-2xl font-bold">Roles</h1>
                <Button onClick={() => { setEditingRole(null); setShowCreateModal(true); }}>Create Role</Button>
            </div>

            {error && <Alert type="error" onClose={() => setError(null)}>{error}</Alert>}

            <div className="flex gap-4 mb-4 flex-wrap">
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
                <Input
                    placeholder="Filter by Level"
                    type="number"
                    value={filterLevel}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFilterLevel(e.target.value)}
                    className="max-w-xs"
                />
                <select
                    value={filterIsSystem}
                    onChange={(e) => setFilterIsSystem(e.target.value)}
                    className="border rounded p-2"
                >
                    <option value="">All Types</option>
                    <option value="true">System</option>
                    <option value="false">Non-System</option>
                </select>
            </div>

            {loading ? (
                <div className="flex justify-center p-8"><Spinner /></div>
            ) : (
                <Table
                    data={roles}
                    columns={columns}
                    keyExtractor={(item) => item.id}
                />
            )}

            <Modal
                isOpen={showCreateModal}
                onClose={() => { setShowCreateModal(false); setEditingRole(null); }}
                title={editingRole ? "Edit Role" : "Create Role"}
            >
                <RoleForm
                    initialData={editingRole || undefined}
                    onSuccess={handleSave}
                    onCancel={() => { setShowCreateModal(false); setEditingRole(null); }}
                />
            </Modal>

            <ConfirmModal
                isOpen={!!deleteId}
                title="Delete Role"
                message="Are you sure you want to delete this role? This action cannot be undone."
                onConfirm={handleDelete}
                onClose={() => setDeleteId(null)}
            />
        </div>
    );
};
