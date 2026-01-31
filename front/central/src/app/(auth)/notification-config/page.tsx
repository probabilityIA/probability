"use client";

import { useState } from "react";
import { ConfigListTable } from "@/services/modules/notification-config/ui";
import { NotificationConfigForm } from "@/services/modules/notification-config/ui/components/NotificationConfigForm";
import { Modal } from "@/shared/ui/modal";
import { NotificationConfig } from "@/services/modules/notification-config/domain/types";

export default function NotificationConfigPage() {
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [selectedConfig, setSelectedConfig] = useState<NotificationConfig | undefined>(undefined);
    const [refreshKey, setRefreshKey] = useState(0);

    const handleCreate = () => {
        setSelectedConfig(undefined);
        setIsModalOpen(true);
    };

    const handleEdit = (config: NotificationConfig) => {
        setSelectedConfig(config);
        setIsModalOpen(true);
    };

    const handleSuccess = () => {
        setIsModalOpen(false);
        setRefreshKey((prev) => prev + 1);
    };

    return (
        <div className="min-h-screen bg-gray-50 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
                <h1 className="text-2xl sm:text-3xl font-bold text-gray-900">Configuración de Notificaciones</h1>
            </div>

            <ConfigListTable
                onEdit={handleEdit}
                onCreate={handleCreate}
                refreshKey={refreshKey}
            />

            <Modal
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                title={selectedConfig ? "Editar Configuración" : "Nueva Configuración"}
                size="lg"
            >
                <NotificationConfigForm
                    config={selectedConfig}
                    onSuccess={handleSuccess}
                    onCancel={() => setIsModalOpen(false)}
                />
            </Modal>
        </div>
    );
}
