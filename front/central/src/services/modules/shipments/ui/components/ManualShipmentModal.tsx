import React, { useState } from 'react';
import { Modal, Button, Input } from '@/shared/ui';
import { createShipmentAction } from '../../infra/actions';
import { useToast } from '@/shared/providers/toast-provider';

interface ManualShipmentModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSuccess: () => void;
}

export const ManualShipmentModal: React.FC<ManualShipmentModalProps> = ({ isOpen, onClose, onSuccess }) => {
    const { showToast } = useToast();
    const [loading, setLoading] = useState(false);
    const [formData, setFormData] = useState({
        client_name: '',
        destination_address: '',
        tracking_number: '',
        carrier: ''
    });

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value
        });
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        try {
            await createShipmentAction({
                client_name: formData.client_name,
                destination_address: formData.destination_address,
                tracking_number: formData.tracking_number,
                carrier: formData.carrier,
                status: 'pending'
            });
            showToast('Envío manual registrado exitosamente', 'success');
            onSuccess();
            onClose();
            setFormData({ client_name: '', destination_address: '', tracking_number: '', carrier: '' });
        } catch (error) {
            console.error(error);
            showToast('Error al registrar el envío', 'error');
        } finally {
            setLoading(false);
        }
    };

    return (
        <Modal isOpen={isOpen} onClose={onClose} title="Registrar Envío Manual" size="md">
            <form onSubmit={handleSubmit} className="space-y-4">
                <Input
                    id="client_name"
                    name="client_name"
                    label="Cliente"
                    value={formData.client_name}
                    onChange={handleChange}
                    required
                />
                <Input
                    id="destination_address"
                    name="destination_address"
                    label="Destino"
                    value={formData.destination_address}
                    onChange={handleChange}
                    required
                />
                <Input
                    id="tracking_number"
                    name="tracking_number"
                    label="Tracking"
                    value={formData.tracking_number}
                    onChange={handleChange}
                    required
                />
                <Input
                    id="carrier"
                    name="carrier"
                    label="Transportadora"
                    value={formData.carrier}
                    onChange={handleChange}
                    placeholder="Opcional"
                />

                <div className="flex justify-end gap-2 mt-6">
                    <Button type="button" variant="outline" onClick={onClose}>
                        Cancelar
                    </Button>
                    <Button type="submit" loading={loading} disabled={loading}>
                        Guardar
                    </Button>
                </div>
            </form>
        </Modal>
    );
};
