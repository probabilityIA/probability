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
        order_id: '',
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
                order_id: formData.order_id || undefined,
                client_name: formData.client_name,
                destination_address: formData.destination_address,
                tracking_number: formData.tracking_number,
                carrier: formData.carrier,
                status: 'pending'
            });
            showToast('Envío registrado exitosamente', 'success');
            onSuccess();
            onClose();
            setFormData({ order_id: '', client_name: '', destination_address: '', tracking_number: '', carrier: '' });
        } catch (error) {
            console.error(error);
            showToast('Error al registrar el envío', 'error');
        } finally {
            setLoading(false);
        }
    };

    return (
        <Modal isOpen={isOpen} onClose={onClose} title="Registrar Envío" size="md">
            <form onSubmit={handleSubmit} className="space-y-4">
                <div className="bg-blue-50 border border-blue-200 rounded-lg p-3">
                    <p className="text-xs text-blue-600 font-semibold mb-2">💡 Opción 1: Selecciona una orden</p>
                    <Input
                        id="order_id"
                        name="order_id"
                        label="ID de Orden"
                        placeholder="prob-0001 o UUID..."
                        value={formData.order_id}
                        onChange={handleChange}
                        helperText={formData.order_id ? '✓ Los datos del cliente se buscarán automáticamente' : 'Opcional - si la dejas vacía, completa los datos abajo'}
                    />
                </div>

                <div className="border-t pt-4">
                    <p className="text-xs text-gray-500 font-semibold mb-3">💬 Opción 2: Completa manualmente (si no usas orden)</p>
                    <div className="space-y-4">
                        <Input
                            id="client_name"
                            name="client_name"
                            label="Cliente"
                            value={formData.client_name}
                            onChange={handleChange}
                            required={!formData.order_id}
                            disabled={!!formData.order_id}
                            placeholder={formData.order_id ? 'Se obtiene de la orden' : 'Nombre del cliente'}
                        />
                        <Input
                            id="destination_address"
                            name="destination_address"
                            label="Destino"
                            value={formData.destination_address}
                            onChange={handleChange}
                            required={!formData.order_id}
                            disabled={!!formData.order_id}
                            placeholder={formData.order_id ? 'Se obtiene de la orden' : 'Dirección de entrega'}
                        />
                    </div>
                </div>

                <div className="border-t pt-4 space-y-4">
                    <p className="text-xs text-gray-500 font-semibold">📦 Detalles del Envío</p>
                    <Input
                        id="tracking_number"
                        name="tracking_number"
                        label="Número de Tracking"
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
                        placeholder="Ej: Servientrega, DHL, etc"
                    />
                </div>

                <div className="flex justify-end gap-2 mt-6 pt-4 border-t">
                    <Button type="button" variant="outline" onClick={onClose}>
                        Cancelar
                    </Button>
                    <Button type="submit" loading={loading} disabled={loading}>
                        Guardar Envío
                    </Button>
                </div>
            </form>
        </Modal>
    );
};
