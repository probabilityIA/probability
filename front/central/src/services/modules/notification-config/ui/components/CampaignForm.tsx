"use client";

import { useEffect, useMemo, useState } from "react";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import { Button } from "@/shared/ui/button";
import { useToast } from "@/shared/providers/toast-provider";
import { WhatsappTemplate } from "../../domain/scheduled-types";
import {
  Campaign,
  CampaignAudiencePreview,
  CampaignAudienceType,
  CreateCampaignDTO,
} from "../../domain/campaign-types";
import {
  createCampaignAction,
  previewCampaignAudienceAction,
  updateCampaignAction,
} from "../../infra/actions/campaigns";
import { WhatsAppBubblePreview } from "./WhatsAppBubblePreview";
import { CampaignRules } from "./CampaignRules";

interface CampaignFormProps {
  businessId?: number;
  templates: WhatsappTemplate[];
  campaign?: Campaign | null;
  onSuccess: () => void;
  onCancel: () => void;
}

const SAMPLE_BY_SOURCE: Record<string, string> = {
  "customer.first_name": "Ana",
  "customer.full_name": "Ana Ramirez",
  "customer.days_inactive": "45",
  "customer.last_product": "Camiseta blanca",
  "customer.total_orders": "3",
  "business.name": "Mi Tienda",
  "campaign.name": "Ruta 30",
};

const TEMPLATE_STATUS_LABEL: Record<string, string> = {
  draft: "borrador",
  pending: "en revisión de Meta",
  paused: "pausada por Meta",
  disabled: "deshabilitada",
  failed: "falló el envío a Meta",
};

export function CampaignForm({
  businessId,
  templates,
  campaign,
  onSuccess,
  onCancel,
}: CampaignFormProps) {
  const { showToast } = useToast();

  const [name, setName] = useState(campaign?.Name || "");
  const [senderName, setSenderName] = useState(campaign?.SenderName || "");
  const [templateId, setTemplateId] = useState<number>(campaign?.WhatsappTemplateID || 0);
  const [audienceType, setAudienceType] = useState<CampaignAudienceType>(
    campaign?.AudienceType || "filtered_clients",
  );
  const [city, setCity] = useState(campaign?.AudienceParams?.City || "");
  const [createdFromDays, setCreatedFromDays] = useState(
    campaign?.AudienceParams?.CreatedFromDays || 0,
  );
  const [onlyWithoutOrder, setOnlyWithoutOrder] = useState(
    campaign?.AudienceParams?.OnlyWithoutOrder || false,
  );
  const [windowStart, setWindowStart] = useState(campaign?.SendWindowStart || "09:00");
  const [windowEnd, setWindowEnd] = useState(campaign?.SendWindowEnd || "19:00");
  const [dailyCap, setDailyCap] = useState(campaign?.DailySendCap || 250);
  const [batchSize, setBatchSize] = useState(campaign?.BatchSize || 50);

  const [preview, setPreview] = useState<CampaignAudiencePreview | null>(null);
  const [previewLoading, setPreviewLoading] = useState(false);
  const [saving, setSaving] = useState(false);

  const selectable = useMemo(
    () => templates.filter((template) => template.Status !== "rejected"),
    [templates],
  );

  const selected = useMemo(
    () => templates.find((template) => template.ID === templateId) || null,
    [templates, templateId],
  );

  const buildDTO = (): CreateCampaignDTO => ({
    whatsapp_template_id: templateId,
    name: name.trim(),
    sender_name: senderName.trim(),
    audience_type: audienceType,
    city: audienceType === "filtered_clients" ? city.trim() : "",
    created_from_days: audienceType === "filtered_clients" ? createdFromDays : 0,
    only_without_order: audienceType === "filtered_clients" ? onlyWithoutOrder : false,
    send_window_start: windowStart,
    send_window_end: windowEnd,
    daily_send_cap: dailyCap,
    batch_size: batchSize,
    variable_values: senderName.trim() ? { "sender.name": senderName.trim() } : undefined,
  });

  useEffect(() => {
    const loadPreview = async () => {
      setPreviewLoading(true);
      const result = await previewCampaignAudienceAction(
        { ...buildDTO(), whatsapp_template_id: templateId || 0 },
        businessId,
      );
      setPreviewLoading(false);
      if (result.success && result.data) {
        setPreview(result.data);
      }
    };

    const timer = setTimeout(loadPreview, 400);
    return () => clearTimeout(timer);
  }, [audienceType, city, createdFromDays, onlyWithoutOrder, businessId]);

  const sampleValues: Record<number, string> = {};
  for (const variable of selected?.Variables || []) {
    if (variable.Source === "sender.name") {
      sampleValues[variable.Position] = senderName.trim() || "tu nombre";
      continue;
    }
    sampleValues[variable.Position] =
      variable.Fallback || SAMPLE_BY_SOURCE[variable.Source] || "ejemplo";
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!name.trim()) {
      showToast("Ponele un nombre a la campaña", "error");
      return;
    }
    if (!templateId) {
      showToast("Elegí la plantilla que se va a enviar", "error");
      return;
    }

    setSaving(true);
    const result = campaign
      ? await updateCampaignAction(campaign.ID, buildDTO(), businessId)
      : await createCampaignAction(buildDTO(), businessId);
    setSaving(false);

    if (!result.success) {
      showToast(result.error || "No se pudo guardar la campaña", "error");
      return;
    }

    showToast(campaign ? "Campaña actualizada" : "Campaña creada como borrador", "success");
    onSuccess();
  };

  return (
    <form onSubmit={handleSubmit} className="grid gap-6 lg:grid-cols-[minmax(0,1fr)_320px]">
      <div className="space-y-4">
        <div>
          <Label htmlFor="campaign-name">{"Nombre de la campaña"}</Label>
          <p className="mb-1 text-xs text-gray-500">
            {"Solo lo ves vos, para diferenciarla de las demás. El cliente nunca lo lee."}
          </p>
          <Input
            id="campaign-name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Ruta 30 - septiembre"
          />
        </div>

        <div>
          <Label htmlFor="campaign-sender">{"Nombre de quien escribe"}</Label>
          <p className="mb-1 text-xs text-gray-500">
            {
              "Es opcional. Si la plantilla usa la variable «Nombre de quien escribe», se reemplaza por esto en cada mensaje. Así la misma plantilla la puede usar otra asesora con su nombre."
            }
          </p>
          <Input
            id="campaign-sender"
            value={senderName}
            onChange={(e) => setSenderName(e.target.value)}
            placeholder="Isabel Rojas"
          />
        </div>

        <div>
          <Label htmlFor="campaign-template">{"Plantilla"}</Label>
          <p className="mb-1 text-xs text-gray-500">
            {
              "Son tus plantillas de campaña. A la derecha ves exactamente el mensaje que le va a llegar al cliente. Para lanzar, Meta tiene que haberla aprobado."
            }
          </p>
          <select
            id="campaign-template"
            value={templateId}
            onChange={(e) => setTemplateId(Number(e.target.value))}
            className="w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 dark:border-gray-600 dark:bg-gray-700 dark:text-white"
          >
            <option value={0}>{"Seleccionar..."}</option>
            {selectable.map((template) => (
              <option key={template.ID} value={template.ID}>
                {template.Status === "approved"
                  ? template.Name
                  : `${template.Name} (${TEMPLATE_STATUS_LABEL[template.Status] || template.Status})`}
              </option>
            ))}
          </select>
          {selectable.length === 0 ? (
            <p className="mt-1 text-xs text-amber-600">
              {"Todavía no tenés plantillas de campaña. Creá una primero."}
            </p>
          ) : selected && selected.Status !== "approved" ? (
            <p className="mt-1 text-xs text-amber-600">
              {"Podés dejar la campaña armada, pero no vas a poder lanzarla hasta que Meta apruebe esta plantilla."}
            </p>
          ) : null}
        </div>

        <div className="rounded-lg border border-gray-200 p-3 dark:border-gray-600">
          <Label>{"A quién le llega"}</Label>
          <p className="mb-2 text-xs text-gray-500">
            {"Son los clientes que ya tenés cargados en el módulo de Clientes de este negocio."}
          </p>

          <div className="mb-3 flex gap-2">
            <button
              type="button"
              onClick={() => setAudienceType("all_clients")}
              className={`rounded-md px-3 py-1 text-xs font-medium ${
                audienceType === "all_clients"
                  ? "bg-[var(--color-primary)] text-white"
                  : "bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-300"
              }`}
            >
              {"Todos mis clientes"}
            </button>
            <button
              type="button"
              onClick={() => setAudienceType("filtered_clients")}
              className={`rounded-md px-3 py-1 text-xs font-medium ${
                audienceType === "filtered_clients"
                  ? "bg-[var(--color-primary)] text-white"
                  : "bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-300"
              }`}
            >
              {"Filtrar"}
            </button>
          </div>

          {audienceType === "filtered_clients" && (
            <div className="grid gap-3 sm:grid-cols-2">
              <div>
                <Label htmlFor="campaign-city">{"Ciudad"}</Label>
                <Input
                  id="campaign-city"
                  value={city}
                  onChange={(e) => setCity(e.target.value)}
                  placeholder="Bogota"
                />
              </div>
              <div>
                <Label htmlFor="campaign-days">{"Cargados en los últimos (días)"}</Label>
                <Input
                  id="campaign-days"
                  type="number"
                  min={0}
                  value={createdFromDays}
                  onChange={(e) => setCreatedFromDays(Number(e.target.value))}
                />
              </div>
              <label className="flex items-center gap-2 text-xs text-gray-600 dark:text-gray-300 sm:col-span-2">
                <input
                  type="checkbox"
                  checked={onlyWithoutOrder}
                  onChange={(e) => setOnlyWithoutOrder(e.target.checked)}
                />
                {"Solo los que todavía no me han comprado"}
              </label>
            </div>
          )}
        </div>

        <div className="grid gap-3 sm:grid-cols-2">
          <div>
            <Label htmlFor="campaign-start">{"Enviar desde"}</Label>
            <Input
              id="campaign-start"
              type="time"
              value={windowStart}
              onChange={(e) => setWindowStart(e.target.value)}
            />
          </div>
          <div>
            <Label htmlFor="campaign-end">{"Hasta"}</Label>
            <Input
              id="campaign-end"
              type="time"
              value={windowEnd}
              onChange={(e) => setWindowEnd(e.target.value)}
            />
          </div>
          <div>
            <Label htmlFor="campaign-cap">{"Máximo por día"}</Label>
            <p className="mb-1 text-xs text-gray-500">
              {"Meta permite hasta 1.000 conversaciones diarias en un número nuevo."}
            </p>
            <Input
              id="campaign-cap"
              type="number"
              min={1}
              max={1000}
              value={dailyCap}
              onChange={(e) => setDailyCap(Number(e.target.value))}
            />
          </div>
          <div>
            <Label htmlFor="campaign-batch">{"Mensajes por tanda"}</Label>
            <p className="mb-1 text-xs text-gray-500">
              {"Cada tanda sale cada pocos minutos, para no disparar todo junto."}
            </p>
            <Input
              id="campaign-batch"
              type="number"
              min={1}
              max={200}
              value={batchSize}
              onChange={(e) => setBatchSize(Number(e.target.value))}
            />
          </div>
        </div>

        <div className="flex justify-end gap-2 border-t pt-4">
          <Button type="button" variant="secondary" onClick={onCancel} disabled={saving}>
            {"Cancelar"}
          </Button>
          <Button type="submit" disabled={saving}>
            {saving ? "Guardando..." : campaign ? "Guardar cambios" : "Crear campaña"}
          </Button>
        </div>
      </div>

      <div className="space-y-4">
        <div className="rounded-lg border border-gray-200 p-3 dark:border-gray-600">
          <p className="mb-2 text-xs font-semibold uppercase tracking-wide text-gray-500">
            {"Audiencia"}
          </p>
          {previewLoading ? (
            <p className="text-xs text-gray-400">{"Calculando..."}</p>
          ) : preview ? (
            <div className="space-y-1 text-xs text-gray-600 dark:text-gray-300">
              <p>
                <span className="text-lg font-bold text-[var(--color-primary)]">
                  {preview.reachable}
                </span>{" "}
                {"personas van a recibir el mensaje"}
              </p>
              <p className="text-gray-400">
                {`${preview.total} coinciden con el filtro · ${preview.opted_out} se dieron de baja · ${preview.no_phone} sin celular válido`}
              </p>
              {preview.sample_names && preview.sample_names.length > 0 && (
                <p className="text-gray-400">
                  {`Por ejemplo: ${preview.sample_names.slice(0, 3).join(", ")}`}
                </p>
              )}
            </div>
          ) : (
            <p className="text-xs text-gray-400">{"Ajustá los filtros para ver el alcance."}</p>
          )}
        </div>

        {selected ? (
          <WhatsAppBubblePreview
            header={selected.HeaderText}
            body={selected.BodyText}
            footer={selected.FooterText}
            buttons={selected.Category === "MARKETING" ? ["Dejar de recibir"] : []}
            sampleValues={sampleValues}
          />
        ) : (
          <div className="rounded-xl border border-dashed border-gray-300 p-6 text-center text-xs text-gray-400 dark:border-gray-600">
            {"Elegí una plantilla para ver cómo le llega al cliente."}
          </div>
        )}

        <CampaignRules />
      </div>
    </form>
  );
}
