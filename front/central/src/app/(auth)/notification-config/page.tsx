"use client";

import { useState } from "react";
import { ConfigListTable } from "@/services/modules/notification-config/ui";
import { NotificationConfigForm } from "@/services/modules/notification-config/ui/components/NotificationConfigForm";
import { NotificationEventTypeList } from "@/services/modules/notification-config/ui/components/NotificationEventTypeList";
import { NotificationEventTypeForm } from "@/services/modules/notification-config/ui/components/NotificationEventTypeForm";
import { NotificationTypeList } from "@/services/modules/notification-config/ui/components/NotificationTypeList";
import { NotificationTypeForm } from "@/services/modules/notification-config/ui/components/NotificationTypeForm";
import { Modal } from "@/shared/ui/modal";
import { NotificationConfig, NotificationEventType, NotificationType } from "@/services/modules/notification-config/domain/types";

type TabType = "configs" | "event-types" | "channels";

export default function NotificationConfigPage() {
    // Estado de tabs
    const [activeTab, setActiveTab] = useState<TabType>("configs");

    // Estado para Configuraciones
    const [isConfigModalOpen, setIsConfigModalOpen] = useState(false);
    const [selectedConfig, setSelectedConfig] = useState<NotificationConfig | undefined>(undefined);
    const [configRefreshKey, setConfigRefreshKey] = useState(0);

    // Estado para Tipos de Eventos
    const [isEventTypeModalOpen, setIsEventTypeModalOpen] = useState(false);
    const [selectedEventType, setSelectedEventType] = useState<NotificationEventType | undefined>(undefined);
    const [eventTypeRefreshKey, setEventTypeRefreshKey] = useState(0);

    // Estado para Canales (Notification Types)
    const [isChannelModalOpen, setIsChannelModalOpen] = useState(false);
    const [selectedChannel, setSelectedChannel] = useState<NotificationType | undefined>(undefined);
    const [channelRefreshKey, setChannelRefreshKey] = useState(0);

    // Handlers para Configuraciones
    const handleCreateConfig = () => {
        setSelectedConfig(undefined);
        setIsConfigModalOpen(true);
    };

    const handleEditConfig = (config: NotificationConfig) => {
        setSelectedConfig(config);
        setIsConfigModalOpen(true);
    };

    const handleConfigSuccess = () => {
        setIsConfigModalOpen(false);
        setConfigRefreshKey((prev) => prev + 1);
    };

    // Handlers para Tipos de Eventos
    const handleCreateEventType = () => {
        setSelectedEventType(undefined);
        setIsEventTypeModalOpen(true);
    };

    const handleEditEventType = (eventType: NotificationEventType) => {
        setSelectedEventType(eventType);
        setIsEventTypeModalOpen(true);
    };

    const handleEventTypeSuccess = () => {
        setIsEventTypeModalOpen(false);
        setEventTypeRefreshKey((prev) => prev + 1);
    };

    // Handlers para Canales
    const handleCreateChannel = () => {
        setSelectedChannel(undefined);
        setIsChannelModalOpen(true);
    };

    const handleEditChannel = (channel: NotificationType) => {
        setSelectedChannel(channel);
        setIsChannelModalOpen(true);
    };

    const handleChannelSuccess = () => {
        setIsChannelModalOpen(false);
        setChannelRefreshKey((prev) => prev + 1);
    };

    return (
        <div className="min-h-screen bg-gray-50 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
                <h1 className="text-2xl sm:text-3xl font-bold text-gray-900">Configuración de Notificaciones</h1>
            </div>

            {/* Tabs */}
            <div className="mb-6">
                <div className="border-b border-gray-200">
                    <nav className="-mb-px flex space-x-8" aria-label="Tabs">
                        <button
                            onClick={() => setActiveTab("configs")}
                            className={`
                                whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm
                                ${activeTab === "configs"
                                    ? "border-blue-500 text-blue-600"
                                    : "border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300"
                                }
                            `}
                        >
                            Configuraciones
                        </button>
                        <button
                            onClick={() => setActiveTab("event-types")}
                            className={`
                                whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm
                                ${activeTab === "event-types"
                                    ? "border-blue-500 text-blue-600"
                                    : "border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300"
                                }
                            `}
                        >
                            Tipos de Eventos
                        </button>
                        <button
                            onClick={() => setActiveTab("channels")}
                            className={`
                                whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm
                                ${activeTab === "channels"
                                    ? "border-blue-500 text-blue-600"
                                    : "border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300"
                                }
                            `}
                        >
                            Canales
                        </button>
                    </nav>
                </div>
            </div>

            {/* Contenido de tabs */}
            {activeTab === "configs" && (
                <ConfigListTable
                    onEdit={handleEditConfig}
                    onCreate={handleCreateConfig}
                    refreshKey={configRefreshKey}
                />
            )}

            {activeTab === "event-types" && (
                <NotificationEventTypeList
                    onEdit={handleEditEventType}
                    onCreate={handleCreateEventType}
                    refreshKey={eventTypeRefreshKey}
                />
            )}

            {activeTab === "channels" && (
                <NotificationTypeList
                    onEdit={handleEditChannel}
                    onCreate={handleCreateChannel}
                    refreshKey={channelRefreshKey}
                />
            )}

            {/* Modal para Configuraciones */}
            <Modal
                isOpen={isConfigModalOpen}
                onClose={() => setIsConfigModalOpen(false)}
                title={selectedConfig ? "Editar Configuración" : "Nueva Configuración"}
                size="lg"
            >
                <NotificationConfigForm
                    config={selectedConfig}
                    onSuccess={handleConfigSuccess}
                    onCancel={() => setIsConfigModalOpen(false)}
                />
            </Modal>

            {/* Modal para Tipos de Eventos */}
            <Modal
                isOpen={isEventTypeModalOpen}
                onClose={() => setIsEventTypeModalOpen(false)}
                title={selectedEventType ? "Editar Tipo de Evento" : "Nuevo Tipo de Evento"}
            >
                <NotificationEventTypeForm
                    eventType={selectedEventType}
                    onSuccess={handleEventTypeSuccess}
                    onCancel={() => setIsEventTypeModalOpen(false)}
                />
            </Modal>

            {/* Modal para Canales */}
            <Modal
                isOpen={isChannelModalOpen}
                onClose={() => setIsChannelModalOpen(false)}
                title={selectedChannel ? "Editar Canal" : "Nuevo Canal"}
            >
                <NotificationTypeForm
                    notificationType={selectedChannel}
                    onSuccess={handleChannelSuccess}
                    onCancel={() => setIsChannelModalOpen(false)}
                />
            </Modal>
        </div>
    );
}
