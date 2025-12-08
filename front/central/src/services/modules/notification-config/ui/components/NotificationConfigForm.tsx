"use client";

import { useState, useEffect } from "react";
import { NotificationConfig, CreateConfigDTO, UpdateConfigDTO } from "../../domain/types";
import { notificationConfigUseCases } from "../../app/use-cases";
import { Button } from "@/shared/ui/button";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import { useToast } from "@/shared/providers/toast-provider";
import { Checkbox } from "@/shared/ui/checkbox";

interface NotificationConfigFormProps {
    config?: NotificationConfig;
    onSuccess: () => void;
    onCancel: () => void;
}

export function NotificationConfigForm({ config, onSuccess, onCancel }: NotificationConfigFormProps) {
    const [loading, setLoading] = useState(false);
    const { showToast } = useToast();

    const [formData, setFormData] = useState<CreateConfigDTO>({
        business_id: 0,
        event_type: "",
        enabled: true,
        channels: ["sse"],
        description: "",
        filters: {},
    });

    useEffect(() => {
        if (config) {
            setFormData({
                business_id: config.business_id,
                event_type: config.event_type,
                enabled: config.enabled,
                channels: config.channels,
                description: config.description,
                filters: config.filters,
            });
        }
    }, [config]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);

        try {
            let response;
            if (config) {
                const updateDto: UpdateConfigDTO = {
                    enabled: formData.enabled,
                    channels: formData.channels,
                    description: formData.description,
                    filters: formData.filters,
                };
                response = await notificationConfigUseCases.updateConfig(config.id, updateDto);
            } else {
                response = await notificationConfigUseCases.createConfig(formData);
            }

            if (response.success) {
                showToast(config ? "Configuración actualizada" : "Configuración creada", "success");
                onSuccess();
            } else {
                showToast(response.error || "Error al guardar configuración", "error");
            }
        } catch (error) {
            showToast("Error inesperado", "error");
        } finally {
            setLoading(false);
        }
    };

    const handleChannelChange = (channel: string, checked: boolean) => {
        const currentChannels = formData.channels || [];
        if (checked) {
            setFormData({ ...formData, channels: [...currentChannels, channel] });
        } else {
            setFormData({ ...formData, channels: currentChannels.filter((c) => c !== channel) });
        }
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid gap-4">
                {!config && (
                    <>
                        <div className="grid gap-2">
                            <Label htmlFor="business_id">Business ID</Label>
                            <Input
                                id="business_id"
                                type="number"
                                value={formData.business_id}
                                onChange={(e) => setFormData({ ...formData, business_id: parseInt(e.target.value) || 0 })}
                                required
                            />
                        </div>
                        <div className="grid gap-2">
                            <Label htmlFor="event_type">Tipo de Evento</Label>
                            <Input
                                id="event_type"
                                value={formData.event_type}
                                onChange={(e) => setFormData({ ...formData, event_type: e.target.value })}
                                placeholder="ej. order.created"
                                required
                            />
                        </div>
                    </>
                )}

                <div className="grid gap-2">
                    <Label>Canales</Label>
                    <div className="flex gap-4">
                        <div className="flex items-center space-x-2">
                            <Checkbox
                                id="channel-sse"
                                checked={formData.channels?.includes("sse")}
                                onCheckedChange={(checked: boolean) => handleChannelChange("sse", checked)}
                            />
                            <label htmlFor="channel-sse" className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                                SSE
                            </label>
                        </div>
                        {/* Add more channels here if needed */}
                    </div>
                </div>

                <div className="flex items-center space-x-2">
                    <Checkbox
                        id="enabled"
                        checked={formData.enabled}
                        onCheckedChange={(checked: boolean) => setFormData({ ...formData, enabled: checked })}
                    />
                    <label htmlFor="enabled" className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                        Habilitado
                    </label>
                </div>

                <div className="grid gap-2">
                    <Label htmlFor="description">Descripción</Label>
                    <Input
                        id="description"
                        value={formData.description}
                        onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                    />
                </div>
            </div>

            <div className="flex justify-end gap-2 pt-4">
                <Button type="button" variant="outline" onClick={onCancel} disabled={loading}>
                    Cancelar
                </Button>
                <Button type="submit" disabled={loading}>
                    {loading ? "Guardando..." : config ? "Actualizar" : "Crear"}
                </Button>
            </div>
        </form>
    );
}
