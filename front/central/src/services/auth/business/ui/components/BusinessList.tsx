import React, { useState } from 'react';
import { Table } from '@/shared/ui/table';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { Select } from '@/shared/ui/select';
import { Spinner } from '@/shared/ui/spinner';
import { ConfirmModal } from '@/shared/ui/confirm-modal';
import { Alert } from '@/shared/ui/alert';
import { Modal } from '@/shared/ui/modal';
import { Business } from '../../domain/types';
import { BusinessForm } from './BusinessForm';
import { useBusinesses } from '../hooks/useBusinesses';

export const BusinessList: React.FC = () => {
    const {
        businesses,
        loading,
        error,
        page,
        setPage,
        totalPages,
        searchName,
        setSearchName,
        filterType,
        setFilterType,
        businessTypes,
        deleteBusiness,
        refresh,
        setError
    } = useBusinesses();

    const [showCreateModal, setShowCreateModal] = useState(false);
    const [editingBusiness, setEditingBusiness] = useState<Business | null>(null);
    const [deleteId, setDeleteId] = useState<number | null>(null);

    const handleDelete = async () => {
        if (deleteId) {
            const success = await deleteBusiness(deleteId);
            if (success) setDeleteId(null);
        }
    };

    const handleSave = () => {
        setShowCreateModal(false);
        setEditingBusiness(null);
        refresh();
    };

    const columns = [
        { label: 'ID', key: 'id' },
        { label: 'Name', key: 'name' },
        { label: 'Code', key: 'code' },
        {
            label: 'Type',
            key: 'business_type',
            render: (_: unknown, row: Business) => row.business_type?.name || row.business_type_id
        },
        {
            label: 'Active',
            key: 'is_active',
            render: (_: unknown, row: Business) => row.is_active ? 'Yes' : 'No'
        },
        {
            label: 'Actions',
            key: 'actions',
            render: (_: unknown, row: Business) => (
                <div className="flex gap-2">
                    <Button variant="secondary" size="sm" onClick={() => { setEditingBusiness(row); setShowCreateModal(true); }}>Edit</Button>
                    <Button variant="danger" size="sm" onClick={() => setDeleteId(row.id)}>Delete</Button>
                </div>
            ),
        },
    ];

    return (
        <div className="p-6 space-y-6">
            <div className="flex justify-between items-center">
                <h1 className="text-2xl font-bold">Businesses</h1>
                <Button onClick={() => { setEditingBusiness(null); setShowCreateModal(true); }}>Create Business</Button>
            </div>

            {error && <Alert type="error" onClose={() => setError(null)}>{error}</Alert>}

            <div className="flex gap-4 mb-4">
                <Input
                    placeholder="Search by name..."
                    value={searchName}
                    onChange={(e) => setSearchName(e.target.value)}
                    className="max-w-xs"
                />
                <Select
                    value={filterType}
                    onChange={(e) => setFilterType(e.target.value)}
                    options={[
                        { label: 'All Types', value: '' },
                        ...businessTypes.map(t => ({ label: t.name, value: String(t.id) }))
                    ]}
                    className="max-w-xs"
                />
            </div>

            {loading ? (
                <div className="flex justify-center p-8"><Spinner /></div>
            ) : (
                <Table
                    data={businesses}
                    columns={columns}
                    keyExtractor={(item) => item.id}
                />
            )}

            <div className="flex justify-center gap-2 mt-4">
                <Button disabled={page === 1} onClick={() => setPage(p => p - 1)}>Previous</Button>
                <span className="self-center">Page {page} of {totalPages}</span>
                <Button disabled={page === totalPages} onClick={() => setPage(p => p + 1)}>Next</Button>
            </div>

            <Modal
                isOpen={showCreateModal}
                onClose={() => { setShowCreateModal(false); setEditingBusiness(null); }}
                title={editingBusiness ? "Edit Business" : "Create Business"}
            >
                <BusinessForm
                    initialData={editingBusiness || undefined}
                    onSuccess={handleSave}
                    onCancel={() => { setShowCreateModal(false); setEditingBusiness(null); }}
                    businessTypes={businessTypes}
                />
            </Modal>

            <ConfirmModal
                isOpen={!!deleteId}
                title="Delete Business"
                message="Are you sure you want to delete this business? This action cannot be undone."
                onConfirm={handleDelete}
                onClose={() => setDeleteId(null)}
            />
        </div>
    );
};
