"use client";

import { useState } from "react";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import { Button } from "@/shared/ui/button";
import { useToast } from "@/shared/providers/toast-provider";
import {
  CreateScheduledRuleDTO,
  WhatsappTemplate,
} from "../../domain/scheduled-types";
import { createScheduledRuleAction } from "../../infra/actions/scheduled-rules";

interface ScheduledRuleFormProps {
  businessId?: number;
  whatsappTypeId: number;
  templates: WhatsappTemplate[];
  onSuccess: () => void;
  onCancel: () => void;
}

export function ScheduledRuleForm({
  businessId,
  whatsappTypeId,
  templates,
  onSuccess,
  onCancel,
}: ScheduledRuleFormProps) {
  const { showToast } = useToast();
  const [loading, setLoading] = useState(false);

  const approved = templates.filter((item) => item.Status === "approved");

  const [name, setName] = useState("");
  const [templateId, setTemplateId] = useState<number | null>(approved[0]?.ID ?? null);
  const [days, setDays] = useState(30);
  const [windowStart, setWindowStart] = useState("09:00");
  const [windowEnd, setWindowEnd] = useState("19:00");
  const [cooldown, setCooldown] = useState(30);
  const [dailyCap, setDailyCap] = useState(500);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!name.trim()) {
      showToast("La regla necesita un nombre", "error");
      return;
    }
    if (!templateId) {
      showToast("Elegí una plantilla aprobada", "error");
      return;
    }
    if (days <= 0) {
      showToast("Los días sin comprar deben ser mayores a cero", "error");
      return;
    }

    const dto: CreateScheduledRuleDTO = {
      notification_type_id: whatsappTypeId,
      whatsapp_template_id: templateId,
      name: name.trim(),
      segment_type: "customers_inactive",
      days_without_purchase: days,
      send_window_start: windowStart,
      send_window_end: windowEnd,
      cooldown_days: cooldown,
      daily_send_cap: dailyCap,
      frequency_minutes: 1440,
      requires_opt_in: true,
      enabled: true,
    };

    setLoading(true);
    const result = await createScheduledRuleAction(dto, businessId);
    setLoading(false);

    if (!result.success) {
      showToast(result.error || "No se pudo crear la regla", "error");
      return;
    }

    showToast("Regla programada creada", "success");
    onSuccess();
  };

  if (approved.length === 0) {
    return (
      <div className="space-y-3 rounded-md border border-amber-200 bg-amber-50 p-4">
        <p className="text-sm text-amber-800">
          {
            "No hay plantillas aprobadas todavía. Creá una y esperá a que Meta la apruebe antes de programar envíos."
          }
        </p>
        <Button type="button" variant="outline" onClick={onCancel}>
          {"Entendido"}
        </Button>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <Label htmlFor="rule-name">{"Nombre de la regla"}</Label>
        <p className="mb-1 text-xs text-gray-500">
          {"Solo para identificarla en esta lista. El cliente no lo ve."}
        </p>
        <Input
          id="rule-name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Reactivar clientes dormidos"
        />
      </div>

      <div>
        <Label htmlFor="rule-template">{"Plantilla"}</Label>
        <p className="mb-1 text-xs text-gray-500">
          {
            "El mensaje que se env\u00eda. Solo aparecen las que Meta ya aprob\u00f3: las que est\u00e1n en revisi\u00f3n o rechazadas no se pueden usar."
          }
        </p>
        <select
          id="rule-template"
          value={templateId ?? ""}
          onChange={(e) => setTemplateId(Number(e.target.value))}
          className="w-full rounded-md border border-gray-300 px-3 py-2 text-sm"
        >
          {approved.map((template) => (
            <option key={template.ID} value={template.ID}>
              {template.Name}
            </option>
          ))}
        </select>
      </div>

      <div>
        <Label htmlFor="rule-days">{"Días sin comprar"}</Label>
        <Input
          id="rule-days"
          type="number"
          min={1}
          value={days}
          onChange={(e) => setDays(Number(e.target.value))}
        />
      </div>

      <div>
        <p className="mb-1 text-xs text-gray-500">
          {
            "Franja horaria en la que se permite enviar (hora de Colombia). Fuera de ella la regla espera al d\u00eda siguiente, para no escribirle a nadie de madrugada."
          }
        </p>
      </div>

      <div className="grid grid-cols-2 gap-3">
        <div>
          <Label htmlFor="rule-start">{"Enviar desde"}</Label>
          <Input
            id="rule-start"
            type="time"
            value={windowStart}
            onChange={(e) => setWindowStart(e.target.value)}
          />
        </div>
        <div>
          <Label htmlFor="rule-end">{"Enviar hasta"}</Label>
          <Input
            id="rule-end"
            type="time"
            value={windowEnd}
            onChange={(e) => setWindowEnd(e.target.value)}
          />
        </div>
      </div>

      <div className="grid grid-cols-2 gap-3">
        <div>
          <Label htmlFor="rule-cooldown">{"No repetir antes de (días)"}</Label>
          <Input
            id="rule-cooldown"
            type="number"
            min={1}
            value={cooldown}
            onChange={(e) => setCooldown(Number(e.target.value))}
          />
        </div>
        <div>
          <Label htmlFor="rule-cap">{"Máximo de envíos por día"}</Label>
          <Input
            id="rule-cap"
            type="number"
            min={1}
            value={dailyCap}
            onChange={(e) => setDailyCap(Number(e.target.value))}
          />
        </div>
      </div>

      <p className="text-xs text-gray-500">
        {
          "Solo se envía a clientes que aceptaron recibir mensajes de marketing. La plantilla lleva botón de baja."
        }
      </p>

      <div className="flex justify-end gap-2 pt-2">
        <Button type="button" variant="outline" onClick={onCancel} disabled={loading}>
          {"Cancelar"}
        </Button>
        <Button type="submit" disabled={loading}>
          {loading ? "Creando..." : "Crear regla"}
        </Button>
      </div>
    </form>
  );
}
