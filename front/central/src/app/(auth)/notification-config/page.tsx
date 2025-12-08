"use client";

import { useState } from "react";
import { NotificationConfigList } from "@/services/modules/notification-config/ui/components/NotificationConfigList";
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
        <div className="p-6">
            <NotificationConfigList
                key={refreshKey}
                onCreate={handleCreate}
                onEdit={handleEdit}
            />

            <Modal
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                title={selectedConfig ? "Editar Configuración" : "Nueva Configuración"}
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
