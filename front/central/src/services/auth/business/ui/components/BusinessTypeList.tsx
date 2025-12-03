import React, { useState } from 'react';
import { Table } from '@/shared/ui/table';
import { Button } from '@/shared/ui/button';
import { Spinner } from '@/shared/ui/spinner';
import { Alert } from '@/shared/ui/alert';
import { BusinessType } from '../../domain/types';
import { Modal } from '@/shared/ui/modal';
import { BusinessTypeForm } from './BusinessTypeForm';
import { ConfirmModal } from '@/shared/ui/confirm-modal';
import { useBusinessTypes } from '../hooks/useBusinessTypes';

export const BusinessTypeList: React.FC = () => {
    const {
        types,
        loading,
        error,
        deleteType,
        refresh,
        setError
    } = useBusinessTypes();

    const [showModal, setShowModal] = useState(false);
    const [editingType, setEditingType] = useState<BusinessType | null>(null);
    const [deleteId, setDeleteId] = useState<number | null>(null);

    const handleDelete = async () => {
        if (deleteId) {
            const success = await deleteType(deleteId);
            if (success) setDeleteId(null);
        }
    };

    const handleSave = () => {
        setShowModal(false);
        setEditingType(null);
        refresh();
    }

    const columns = [
        { label: 'ID', key: 'id' },
        { label: 'Name', key: 'name' },
        { label: 'Code', key: 'code' },
        { label: 'Active', key: 'is_active', render: (_: unknown, row: BusinessType) => row.is_active ? 'Yes' : 'No' },
        {
            label: 'Actions',
            key: 'actions',
            render: (_: unknown, row: BusinessType) => (
                <div className="flex gap-2">
                    <Button variant="secondary" size="sm" onClick={() => { setEditingType(row); setShowModal(true); }}>Edit</Button>
                    <Button variant="danger" size="sm" onClick={() => setDeleteId(row.id)}>Delete</Button>
                </div>
            )
        }
    ];

    return (
        <div className="p-6 space-y-6">
            <div className="flex justify-between items-center">
                <h1 className="text-2xl font-bold">Business Types</h1>
                <Button onClick={() => { setEditingType(null); setShowModal(true); }}>Create Type</Button>
            </div>

            {error && <Alert type="error" onClose={() => setError(null)}>{error}</Alert>}

            {loading ? <Spinner /> : <Table data={types} columns={columns} keyExtractor={(item) => item.id} />}

            <Modal
                isOpen={showModal}
                onClose={() => { setShowModal(false); setEditingType(null); }}
                title={editingType ? "Edit Type" : "Create Type"}
            >
                <BusinessTypeForm
                    initialData={editingType || undefined}
                    onSuccess={handleSave}
                    onCancel={() => { setShowModal(false); setEditingType(null); }}
                />
            </Modal>

            <ConfirmModal
                isOpen={!!deleteId}
                title="Delete Business Type"
                message="Are you sure?"
                onConfirm={handleDelete}
                onClose={() => setDeleteId(null)}
            />
        </div>
    );
};
