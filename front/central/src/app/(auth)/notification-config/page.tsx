"use client";

import { useState, useEffect } from "react";
import { ConfigListTable } from "@/services/modules/notification-config/ui/components/ConfigListTable";
import { IntegrationRulesForm } from "@/services/modules/notification-config/ui/components/IntegrationRulesForm";
import { IntegrationPicker } from "@/services/modules/notification-config/ui/components/IntegrationPicker";
import { NotificationEventTypeList } from "@/services/modules/notification-config/ui/components/NotificationEventTypeList";
import { NotificationEventTypeForm } from "@/services/modules/notification-config/ui/components/NotificationEventTypeForm";
import { NotificationTypeList } from "@/services/modules/notification-config/ui/components/NotificationTypeList";
import { NotificationTypeForm } from "@/services/modules/notification-config/ui/components/NotificationTypeForm";
import { Modal } from "@/shared/ui/modal";
import { NotificationEventType, NotificationType } from "@/services/modules/notification-config/domain/types";
import { usePermissions } from "@/shared/contexts/permissions-context";
import { useBusinessesSimple } from "@/services/auth/business/ui/hooks/useBusinessesSimple";
import type { IntegrationSimple } from "@/services/integrations/core/domain/types";

type TabType = "configs" | "event-types" | "channels";

export default function NotificationConfigPage() {
    // Super admin business selection
    const { isSuperAdmin } = usePermissions();
    const { businesses, loading: loadingBusinesses } = useBusinessesSimple();
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);

    const requiresBusinessSelection = isSuperAdmin && selectedBusinessId === null;

    // Estado de tabs
    const [activeTab, setActiveTab] = useState<TabType>("configs");

    // Estado para flujo de configs: picker -> rules form
    const [isPickerModalOpen, setIsPickerModalOpen] = useState(false);
    const [isRulesModalOpen, setIsRulesModalOpen] = useState(false);
    const [selectedIntegration, setSelectedIntegration] = useState<IntegrationSimple | undefined>(undefined);
    const [configRefreshKey, setConfigRefreshKey] = useState(0);

    // Estado para Tipos de Eventos
    const [isEventTypeModalOpen, setIsEventTypeModalOpen] = useState(false);
    const [selectedEventType, setSelectedEventType] = useState<NotificationEventType | undefined>(undefined);
    const [eventTypeRefreshKey, setEventTypeRefreshKey] = useState(0);

    // Estado para Canales (Notification Types)
    const [isChannelModalOpen, setIsChannelModalOpen] = useState(false);
    const [selectedChannel, setSelectedChannel] = useState<NotificationType | undefined>(undefined);
    const [channelRefreshKey, setChannelRefreshKey] = useState(0);

    // Reset al cambiar de negocio
    useEffect(() => {
        setConfigRefreshKey((prev) => prev + 1);
        setSelectedIntegration(undefined);
        setIsPickerModalOpen(false);
        setIsRulesModalOpen(false);
    }, [selectedBusinessId]);

    // Handlers para Configs: picker flow
    const handleCreateConfig = () => {
        setIsPickerModalOpen(true);
    };

    const handlePickIntegration = (integration: IntegrationSimple) => {
        setSelectedIntegration(integration);
        setIsPickerModalOpen(false);
        setIsRulesModalOpen(true);
    };

    const handleConfigureIntegration = (integration: IntegrationSimple) => {
        setSelectedIntegration(integration);
        setIsRulesModalOpen(true);
    };

    const handleRulesSuccess = () => {
        setIsRulesModalOpen(false);
        setSelectedIntegration(undefined);
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

            {/* Selector de negocio para super admin */}
            {isSuperAdmin && (
                <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-6">
                    <label className="block text-sm font-medium text-blue-800 mb-2">
                        Seleccionar Negocio
                    </label>
                    <select
                        value={selectedBusinessId?.toString() ?? ''}
                        onChange={(e) => setSelectedBusinessId(e.target.value ? Number(e.target.value) : null)}
                        className="w-full max-w-md px-3 py-2 border border-blue-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                        disabled={loadingBusinesses}
                    >
                        <option value="">— Selecciona un negocio —</option>
                        {businesses.map((b) => (
                            <option key={b.id} value={b.id}>
                                {b.name} (ID: {b.id})
                            </option>
                        ))}
                    </select>
                </div>
            )}

            {/* Gate: bloquear hasta seleccionar negocio */}
            {requiresBusinessSelection ? (
                <div className="text-center py-16 text-gray-500">
                    Selecciona un negocio para ver las configuraciones de notificación
                </div>
            ) : (
                <>
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
                            onConfigure={handleConfigureIntegration}
                            onCreate={handleCreateConfig}
                            refreshKey={configRefreshKey}
                            selectedBusinessId={isSuperAdmin ? selectedBusinessId ?? undefined : undefined}
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

                    {/* Modal: Integration Picker */}
                    <Modal
                        isOpen={isPickerModalOpen}
                        onClose={() => setIsPickerModalOpen(false)}
                        title="Seleccionar Integración"
                    >
                        <IntegrationPicker
                            businessId={isSuperAdmin ? selectedBusinessId ?? undefined : undefined}
                            onSelect={handlePickIntegration}
                            onCancel={() => setIsPickerModalOpen(false)}
                        />
                    </Modal>

                    {/* Modal: Integration Rules Form */}
                    <Modal
                        isOpen={isRulesModalOpen}
                        onClose={() => setIsRulesModalOpen(false)}
                        title={`Reglas de Notificación — ${selectedIntegration?.name || ''}`}
                        size="4xl"
                    >
                        {selectedIntegration && (
                            <IntegrationRulesForm
                                integration={selectedIntegration}
                                businessId={isSuperAdmin ? (selectedBusinessId ?? 0) : 0}
                                onSuccess={handleRulesSuccess}
                                onCancel={() => setIsRulesModalOpen(false)}
                            />
                        )}
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
                            type={selectedChannel}
                            onSuccess={handleChannelSuccess}
                            onCancel={() => setIsChannelModalOpen(false)}
                        />
                    </Modal>
                </>
            )}
        </div>
    );
}
