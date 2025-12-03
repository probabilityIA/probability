import React, { useState } from 'react';
import { Table } from '@/shared/ui/table';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { Spinner } from '@/shared/ui/spinner';
import { ConfirmModal } from '@/shared/ui/confirm-modal';
import { Alert } from '@/shared/ui/alert';
import { Modal } from '@/shared/ui/modal';
import { User } from '../../domain/types';
import { UserForm } from './UserForm';
import { useUsers } from '../hooks/useUsers';

export const UserList: React.FC = () => {
    const {
        users,
        pagination,
        loading,
        error,
        page,
        setPage,
        // pageSize,
        // setPageSize,
        searchName,
        setSearchName,
        searchEmail,
        setSearchEmail,
        filterIsActive,
        setFilterIsActive,
        filterRoleId,
        setFilterRoleId,
        filterBusinessId,
        setFilterBusinessId,
        deleteUser,
        refresh,
        setError
    } = useUsers();

    const [showCreateModal, setShowCreateModal] = useState(false);
    const [editingUser, setEditingUser] = useState<User | null>(null);
    const [deleteId, setDeleteId] = useState<number | null>(null);

    const handleDelete = async () => {
        if (deleteId) {
            const success = await deleteUser(deleteId);
            if (success) setDeleteId(null);
        }
    };

    const handleSave = () => {
        setShowCreateModal(false);
        setEditingUser(null);
        refresh();
    };

    const columns = [
        { label: 'ID', key: 'id' },
        {
            label: 'Avatar',
            key: 'avatar_url',
            render: (_: unknown, row: User) => row.avatar_url ? <img src={row.avatar_url} alt={row.name} className="w-8 h-8 rounded-full" /> : <div className="w-8 h-8 rounded-full bg-gray-200" />
        },
        { label: 'Name', key: 'name' },
        { label: 'Email', key: 'email' },
        { label: 'Phone', key: 'phone' },
        { label: 'Active', key: 'is_active', render: (_: unknown, row: User) => row.is_active ? 'Yes' : 'No' },
        {
            label: 'Actions',
            key: 'actions',
            render: (_: unknown, row: User) => (
                <div className="flex gap-2">
                    <Button variant="secondary" size="sm" onClick={() => { setEditingUser(row); setShowCreateModal(true); }}>Edit</Button>
                    <Button variant="danger" size="sm" onClick={() => setDeleteId(row.id)}>Delete</Button>
                </div>
            ),
        },
    ];

    return (
        <div className="p-6 space-y-6">
            <div className="flex justify-between items-center">
                <h1 className="text-2xl font-bold">Users</h1>
                <Button onClick={() => { setEditingUser(null); setShowCreateModal(true); }}>Create User</Button>
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
                    placeholder="Search by email..."
                    value={searchEmail}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => setSearchEmail(e.target.value)}
                    className="max-w-xs"
                />
                <select
                    value={filterIsActive}
                    onChange={(e) => setFilterIsActive(e.target.value)}
                    className="border rounded p-2"
                >
                    <option value="">All Status</option>
                    <option value="true">Active</option>
                    <option value="false">Inactive</option>
                </select>
                <Input
                    placeholder="Filter by Role ID"
                    type="number"
                    value={filterRoleId}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFilterRoleId(e.target.value)}
                    className="max-w-xs"
                />
                <Input
                    placeholder="Filter by Business ID"
                    type="number"
                    value={filterBusinessId}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFilterBusinessId(e.target.value)}
                    className="max-w-xs"
                />
            </div>

            {loading ? (
                <div className="flex justify-center p-8"><Spinner /></div>
            ) : (
                <>
                    <Table
                        data={users}
                        columns={columns}
                        keyExtractor={(item) => item.id}
                    />
                    {pagination && (
                        <div className="flex justify-between items-center mt-4">
                            <Button
                                disabled={!pagination.has_prev}
                                onClick={() => setPage(page - 1)}
                                variant="secondary"
                            >
                                Previous
                            </Button>
                            <span>Page {pagination.current_page} of {pagination.last_page}</span>
                            <Button
                                disabled={!pagination.has_next}
                                onClick={() => setPage(page + 1)}
                                variant="secondary"
                            >
                                Next
                            </Button>
                        </div>
                    )}
                </>
            )}

            <Modal
                isOpen={showCreateModal}
                onClose={() => { setShowCreateModal(false); setEditingUser(null); }}
                title={editingUser ? "Edit User" : "Create User"}
            >
                <UserForm
                    initialData={editingUser || undefined}
                    onSuccess={handleSave}
                    onCancel={() => { setShowCreateModal(false); setEditingUser(null); }}
                />
            </Modal>

            <ConfirmModal
                isOpen={!!deleteId}
                title="Delete User"
                message="Are you sure you want to delete this user? This action cannot be undone."
                onConfirm={handleDelete}
                onClose={() => setDeleteId(null)}
            />
        </div>
    );
};
