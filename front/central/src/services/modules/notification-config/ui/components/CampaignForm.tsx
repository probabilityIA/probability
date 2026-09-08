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
      showToast("Ponle un nombre a la campaña", "error");
      return;
    }
    if (!templateId) {
      showToast("Elige la plantilla que se va a enviar", "error");
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

  const band = "border-t border-gray-200 py-5 dark:border-gray-700";

  const step = (n: string, title: string, note: string, extra?: React.ReactNode) => (
    <div className="mb-4 flex items-center gap-3">
      <span className="flex h-5 w-5 flex-none items-center justify-center rounded-full bg-[var(--color-primary)] text-[11px] font-semibold text-white">
        {n}
      </span>
      <span className="text-sm font-semibold text-gray-900 dark:text-gray-100">{title}</span>
      <span className="text-xs text-gray-500">{note}</span>
      {extra ? <span className="ml-auto">{extra}</span> : null}
    </div>
  );

  return (
    <form onSubmit={handleSubmit} className="flex flex-col">
      <div className="grid gap-6 lg:grid-cols-[minmax(0,1fr)_360px]">
        <div>
          <div className="pb-5">
            {step("1", "El mensaje", "Qu\u00e9 se manda y qui\u00e9n lo firma")}

            <div className="grid gap-4 md:grid-cols-3">
              <div>
                <Label htmlFor="campaign-name">{"Nombre de la campa\u00f1a"}</Label>
                <Input
                  id="campaign-name"
                  className="mt-1.5"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="Ruta 30 - septiembre"
                />
                <p className="mt-1 text-xs text-gray-500">
                  {"Solo lo ves t\u00fa. El cliente nunca lo lee."}
                </p>
              </div>

              <div>
                <Label htmlFor="campaign-sender">{"Nombre de quien escribe"}</Label>
                <Input
                  id="campaign-sender"
                  className="mt-1.5"
                  value={senderName}
                  onChange={(e) => setSenderName(e.target.value)}
                  placeholder={"Escribe aqu\u00ed tu nombre"}
                />
                <p className="mt-1 text-xs text-gray-500">
                  {"Opcional. Reemplaza la variable \u00abNombre de quien escribe\u00bb de la plantilla."}
                </p>
              </div>

              <div>
                <Label htmlFor="campaign-template">{"Plantilla"}</Label>
                <select
                  id="campaign-template"
                  value={templateId}
                  onChange={(e) => setTemplateId(Number(e.target.value))}
                  className="mt-1.5 w-full rounded-lg border border-gray-300 bg-white px-4 py-3 text-sm text-gray-900 dark:border-gray-600 dark:bg-gray-700 dark:text-white"
                >
                  <option value={0}>{"Seleccionar..."}</option>
                  {selectable.map((template) => (
                    <option key={template.ID} value={template.ID}>
                      {template.Name}
                    </option>
                  ))}
                </select>
                {selectable.length === 0 ? (
                  <p className="mt-1 text-xs text-amber-600">
                    {"Todav\u00eda no tienes plantillas de campa\u00f1a. Crea una primero."}
                  </p>
                ) : selected && selected.Status !== "approved" ? (
                  <p className="mt-1.5 flex items-center gap-1.5 text-xs font-medium text-amber-700 dark:text-amber-500">
                    <svg
                      className="h-3.5 w-3.5"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      strokeWidth={2.2}
                      strokeLinecap="round"
                    >
                      <circle cx="12" cy="12" r="9" />
                      <path d="M12 7v5l3 2" />
                    </svg>
                    {TEMPLATE_STATUS_LABEL[selected.Status] || selected.Status}
                  </p>
                ) : (
                  <p className="mt-1 text-xs text-gray-500">
                    {"Meta tiene que haberla aprobado para poder lanzar."}
                  </p>
                )}
              </div>
            </div>
          </div>

          <div className={band}>
            {step(
              "2",
              "A qui\u00e9n le llega",
              "Clientes cargados en el m\u00f3dulo de Clientes",
              <span className="flex gap-1.5">
                <button
                  type="button"
                  onClick={() => setAudienceType("all_clients")}
                  className={`rounded-md px-3 py-1.5 text-xs font-medium ${
                    audienceType === "all_clients"
                      ? "bg-[var(--color-primary)] text-white"
                      : "bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-300"
                  }`}
                >
                  {"Todos"}
                </button>
                <button
                  type="button"
                  onClick={() => setAudienceType("filtered_clients")}
                  className={`rounded-md px-3 py-1.5 text-xs font-medium ${
                    audienceType === "filtered_clients"
                      ? "bg-[var(--color-primary)] text-white"
                      : "bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-300"
                  }`}
                >
                  {"Filtrar"}
                </button>
              </span>,
            )}

            {audienceType === "filtered_clients" ? (
              <div className="grid items-end gap-4 md:grid-cols-3">
                <div>
                  <Label htmlFor="campaign-city">{"Ciudad"}</Label>
                  <Input
                    id="campaign-city"
                    className="mt-1.5"
                    value={city}
                    onChange={(e) => setCity(e.target.value)}
                    placeholder={"Bogot\u00e1"}
                  />
                </div>
                <div>
                  <Label htmlFor="campaign-days">{"Cargados hace (d\u00edas)"}</Label>
                  <Input
                    id="campaign-days"
                    className="mt-1.5"
                    type="number"
                    min={0}
                    value={createdFromDays}
                    onChange={(e) => setCreatedFromDays(Number(e.target.value))}
                  />
                </div>
                <label className="flex items-center gap-2 pb-3 text-xs text-gray-600 dark:text-gray-300">
                  <input
                    type="checkbox"
                    checked={onlyWithoutOrder}
                    onChange={(e) => setOnlyWithoutOrder(e.target.checked)}
                  />
                  {"Solo los que todav\u00eda no me han comprado"}
                </label>
              </div>
            ) : (
              <p className="text-xs text-gray-500">
                {"Le llega a todos tus clientes con celular v\u00e1lido que no se hayan dado de baja."}
              </p>
            )}
          </div>

          <div className={band}>
            {step("3", "Ritmo de env\u00edo", "Sale por tandas dentro de la franja")}

            <div className="grid gap-4 md:grid-cols-4">
              <div>
                <Label htmlFor="campaign-start">{"Desde"}</Label>
                <Input
                  id="campaign-start"
                  className="mt-1.5"
                  type="time"
                  value={windowStart}
                  onChange={(e) => setWindowStart(e.target.value)}
                />
              </div>
              <div>
                <Label htmlFor="campaign-end">{"Hasta"}</Label>
                <Input
                  id="campaign-end"
                  className="mt-1.5"
                  type="time"
                  value={windowEnd}
                  onChange={(e) => setWindowEnd(e.target.value)}
                />
              </div>
              <div>
                <Label htmlFor="campaign-cap">{"M\u00e1ximo por d\u00eda"}</Label>
                <Input
                  id="campaign-cap"
                  className="mt-1.5"
                  type="number"
                  min={1}
                  max={1000}
                  value={dailyCap}
                  onChange={(e) => setDailyCap(Number(e.target.value))}
                />
                <p className="mt-1 text-xs text-gray-500">{"Meta permite hasta 1.000."}</p>
              </div>
              <div>
                <Label htmlFor="campaign-batch">{"Mensajes por tanda"}</Label>
                <Input
                  id="campaign-batch"
                  className="mt-1.5"
                  type="number"
                  min={1}
                  max={200}
                  value={batchSize}
                  onChange={(e) => setBatchSize(Number(e.target.value))}
                />
                <p className="mt-1 text-xs text-gray-500">{"Una tanda cada pocos minutos."}</p>
              </div>
            </div>
          </div>

          <div className={band}>
            <CampaignRules collapsible />
          </div>
        </div>

        <aside className="space-y-4 lg:sticky lg:top-0 lg:self-start">
          <div className="rounded-xl border border-gray-200 bg-gray-50 p-4 dark:border-gray-600 dark:bg-gray-900/30">
            {previewLoading ? (
              <p className="text-xs text-gray-400">{"Calculando..."}</p>
            ) : preview ? (
              <div>
                <div className="flex items-baseline gap-2">
                  <span className="text-3xl font-bold leading-none text-[var(--color-primary)]">
                    {preview.reachable}
                  </span>
                  <span className="text-xs text-gray-600 dark:text-gray-300">
                    {"personas reciben el mensaje"}
                  </span>
                </div>
                <div className="mt-3 flex gap-5 border-t border-gray-200 pt-3 dark:border-gray-700">
                  <div>
                    <div className="text-sm font-semibold text-gray-900 dark:text-gray-100">
                      {preview.total}
                    </div>
                    <div className="text-[11px] text-gray-500">{"coinciden"}</div>
                  </div>
                  <div>
                    <div className="text-sm font-semibold text-gray-900 dark:text-gray-100">
                      {preview.opted_out}
                    </div>
                    <div className="text-[11px] text-gray-500">{"de baja"}</div>
                  </div>
                  <div>
                    <div className="text-sm font-semibold text-gray-900 dark:text-gray-100">
                      {preview.no_phone}
                    </div>
                    <div className="text-[11px] text-gray-500">{"sin celular"}</div>
                  </div>
                </div>
              </div>
            ) : (
              <p className="text-xs text-gray-400">{"Ajusta los filtros para ver el alcance."}</p>
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
              {"Elige una plantilla para ver c\u00f3mo le llega al cliente."}
            </div>
          )}

          <p className="text-xs leading-relaxed text-gray-500">
            {"La vista previa usa datos de ejemplo. El nombre real de cada cliente entra al enviar."}
          </p>
        </aside>
      </div>

      <div className="mt-5 flex items-center justify-between gap-4 border-t border-gray-200 pt-4 dark:border-gray-700">
        <p className="text-xs text-gray-500">
          {"Se guarda como borrador. No sale ning\u00fan mensaje hasta que la lances."}
        </p>
        <div className="flex gap-2">
          <Button type="button" variant="secondary" onClick={onCancel} disabled={saving}>
            {"Cancelar"}
          </Button>
          <Button type="submit" disabled={saving}>
            {saving ? "Guardando..." : campaign ? "Guardar cambios" : "Crear campa\u00f1a"}
          </Button>
        </div>
      </div>
    </form>
  );
}
