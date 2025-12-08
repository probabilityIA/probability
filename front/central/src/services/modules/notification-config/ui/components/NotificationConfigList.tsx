"use client";

import { useState, useEffect } from "react";
import { NotificationConfig, ConfigFilter } from "../../domain/types";
import { notificationConfigUseCases } from "../../app/use-cases";
import { Button } from "@/shared/ui/button";
import { useToast } from "@/shared/providers/toast-provider";
import { Badge } from "@/shared/ui/badge";

interface NotificationConfigListProps {
    onEdit: (config: NotificationConfig) => void;
    onCreate: () => void;
}

export function NotificationConfigList({ onEdit, onCreate }: NotificationConfigListProps) {
    const [configs, setConfigs] = useState<NotificationConfig[]>([]);
    const [loading, setLoading] = useState(true);
    const { showToast } = useToast();

    const fetchConfigs = async () => {
        setLoading(true);
        const response = await notificationConfigUseCases.listConfigs();
        if (response.success && response.data) {
            setConfigs(response.data);
        } else {
            showToast(response.error || "Error al cargar configuraciones", "error");
        }
        setLoading(false);
    };

    useEffect(() => {
        fetchConfigs();
    }, []);

    const handleDelete = async (id: number) => {
        if (confirm("¿Estás seguro de eliminar esta configuración?")) {
            const response = await notificationConfigUseCases.deleteConfig(id);
            if (response.success) {
                showToast("Configuración eliminada correctamente", "success");
                fetchConfigs();
            } else {
                showToast(response.error || "Error al eliminar configuración", "error");
            }
        }
    };

    if (loading) {
        return <div className="p-4 text-center">Cargando...</div>;
    }

    return (
        <div className="space-y-4">
            <div className="flex justify-between items-center">
                <h2 className="text-2xl font-bold">Configuración de Notificaciones</h2>
                <Button onClick={onCreate}>
                    <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                    </svg>
                    Nueva Configuración
                </Button>
            </div>

            <div className="border rounded-lg overflow-hidden">
                <table className="w-full text-sm text-left">
                    <thead className="bg-gray-50 text-gray-700 uppercase">
                        <tr>
                            <th className="px-6 py-3">ID</th>
                            <th className="px-6 py-3">Business ID</th>
                            <th className="px-6 py-3">Evento</th>
                            <th className="px-6 py-3">Canales</th>
                            <th className="px-6 py-3">Estado</th>
                            <th className="px-6 py-3">Acciones</th>
                        </tr>
                    </thead>
                    <tbody>
                        {configs.length === 0 ? (
                            <tr>
                                <td colSpan={6} className="px-6 py-4 text-center text-gray-500">
                                    No hay configuraciones encontradas
                                </td>
                            </tr>
                        ) : (
                            configs.map((config) => (
                                <tr key={config.id} className="bg-white border-b hover:bg-gray-50">
                                    <td className="px-6 py-4">{config.id}</td>
                                    <td className="px-6 py-4">{config.business_id}</td>
                                    <td className="px-6 py-4 font-medium">{config.event_type}</td>
                                    <td className="px-6 py-4">
                                        <div className="flex gap-1 flex-wrap">
                                            {config.channels.map((channel) => (
                                                <Badge key={channel} type="secondary">
                                                    {channel}
                                                </Badge>
                                            ))}
                                        </div>
                                    </td>
                                    <td className="px-6 py-4">
                                        <Badge type={config.enabled ? "success" : "secondary"}>
                                            {config.enabled ? "Activo" : "Inactivo"}
                                        </Badge>
                                    </td>
                                    <td className="px-6 py-4">
                                        <div className="flex gap-2">
                                            <Button variant="outline" size="sm" onClick={() => onEdit(config)}>
                                                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                                                </svg>
                                            </Button>
                                            <Button variant="outline" size="sm" onClick={() => handleDelete(config.id)} className="text-red-500 hover:text-red-700">
                                                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                                                </svg>
                                            </Button>
                                        </div>
                                    </td>
                                </tr>
                            ))
                        )}
                    </tbody>
                </table>
            </div>
        </div>
    );
}
