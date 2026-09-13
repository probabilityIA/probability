"use client";

import { useCallback, useEffect, useState } from "react";
import { useToast } from "@/shared/providers/toast-provider";
import { Campaign, CampaignSend } from "../../domain/campaign-types";
import { Flow } from "../../domain/scheduled-types";
import {
  cancelCampaignAction,
  launchCampaignAction,
  listCampaignSendsAction,
  pauseCampaignAction,
  resumeCampaignAction,
} from "../../infra/actions/campaigns";
import { WhatsAppConversations } from "./WhatsAppConversations";

interface CampaignDetailProps {
  campaign: Campaign;
  businessId?: number;
  flows?: Flow[];
  onChanged: () => void;
  onBack: () => void;
}

const STATUS_LABEL: Record<string, { label: string; className: string }> = {
  draft: { label: "Borrador", className: "bg-gray-100 text-gray-600" },
  scheduled: { label: "Programada", className: "bg-blue-100 text-blue-700" },
  running: { label: "Enviando", className: "bg-green-100 text-green-700" },
  paused: { label: "Pausada", className: "bg-amber-100 text-amber-700" },
  completed: { label: "Terminada", className: "bg-purple-100 text-purple-700" },
  cancelled: { label: "Cancelada", className: "bg-red-100 text-red-700" },
};

const SEND_STATUS_LABEL: Record<string, string> = {
  pending: "En cola",
  queued: "Enviando",
  sent: "Enviado",
  delivered: "Entregado",
  read: "Leído",
  replied: "Respondió",
  failed: "Falló",
  skipped: "Omitido",
};

export function CampaignDetail({
  campaign,
  businessId,
  flows = [],
  onChanged,
  onBack,
}: CampaignDetailProps) {
  const { showToast } = useToast();
  const [tab, setTab] = useState<"results" | "chat" | "config">("results");
  const [sends, setSends] = useState<CampaignSend[]>([]);
  const [loading, setLoading] = useState(true);
  const [acting, setActing] = useState(false);

  const fetchSends = useCallback(async () => {
    setLoading(true);
    const result = await listCampaignSendsAction(campaign.ID, businessId, undefined, 1, 100);
    setLoading(false);
    if (result.success) setSends(result.data);
  }, [campaign.ID, businessId]);

  useEffect(() => {
    fetchSends();
  }, [fetchSends]);

  const runAction = async (
    action: (id: number, businessId?: number) => Promise<{ success: boolean; error?: string }>,
    okMessage: string,
  ) => {
    setActing(true);
    const result = await action(campaign.ID, businessId);
    setActing(false);

    if (!result.success) {
      showToast(result.error || "No se pudo completar la acción", "error");
      return;
    }

    showToast(okMessage, "success");
    onChanged();
    fetchSends();
  };

  const status = STATUS_LABEL[campaign.Status] || STATUS_LABEL.draft;

  const params = campaign.AudienceParams;
  const flowName = flows.find((flow) => flow.ID === campaign.FlowID)?.Name;
  const manualCount = params?.ClientIDs?.length || 0;

  const audienceLabel =
    campaign.AudienceType === "all_clients"
      ? "Todos los clientes"
      : manualCount > 0
        ? `${manualCount} elegidos a mano`
        : "Filtrada";

  const filterRows = [
    params?.State ? { label: "Departamento", value: params.State } : null,
    params?.City ? { label: "Ciudad", value: params.City } : null,
    params?.CreatedFromDays
      ? { label: "Registrado en los últimos", value: `${params.CreatedFromDays} días` }
      : null,
    params?.RegisteredBeforeDays
      ? { label: "Registrado hace más de", value: `${params.RegisteredBeforeDays} días` }
      : null,
    params?.OnlyWithoutOrder ? { label: "Nunca me ha comprado", value: "Sí" } : null,
    params?.MinOrders ? { label: "Con al menos", value: `${params.MinOrders} compras` } : null,
    params?.MinSpent
      ? { label: "Ha gastado al menos", value: `$${params.MinSpent.toLocaleString("es-CO")}` }
      : null,
    params?.LastPurchaseBeforeDays
      ? { label: "Sin comprar hace más de", value: `${params.LastPurchaseBeforeDays} días` }
      : null,
  ].filter((row): row is { label: string; value: string } => row !== null);

  const configGroups = [
    {
      title: "Qué se envía",
      rows: [
        {
          label: "Flujo",
          value: flowName || (campaign.FlowID ? `#${campaign.FlowID}` : "Sin flujo"),
        },
        { label: "Remitente", value: campaign.SenderName || "-" },
      ],
    },
    {
      title: "A quién le llega",
      rows: [
        { label: "Audiencia", value: audienceLabel },
        ...(filterRows.length > 0 ? filterRows : [{ label: "Filtros", value: "Sin filtros" }]),
        {
          label: "Tope de frecuencia",
          value: params?.ExcludeRecentDays
            ? `Máximo ${params.ExcludeRecentMax || 1} en ${params.ExcludeRecentDays} días`
            : "Sin tope",
        },
      ],
    },
    {
      title: "Cuándo sale",
      rows: [
        {
          label: "Arranque",
          value: campaign.ScheduledAt
            ? new Date(campaign.ScheduledAt).toLocaleString("es-CO")
            : "Al lanzarla",
        },
        {
          label: "Franja horaria",
          value: `${campaign.SendWindowStart} a ${campaign.SendWindowEnd}`,
        },
        { label: "Zona horaria", value: campaign.Timezone },
        { label: "Tope diario", value: String(campaign.DailySendCap) },
        { label: "Tamaño de tanda", value: String(campaign.BatchSize) },
      ],
    },
  ];

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="flex items-center gap-3">
          <button
            type="button"
            onClick={onBack}
            className="rounded-md bg-gray-100 px-2 py-1 text-xs text-gray-600 hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-300"
          >
            {"Volver"}
          </button>
          <div>
            <h3 className="text-sm font-semibold text-gray-900 dark:text-white">
              {campaign.Name}
            </h3>
            <span className={`inline-block rounded-full px-2 py-0.5 text-[10px] font-semibold ${status.className}`}>
              {status.label}
            </span>
          </div>
        </div>

        <div className="flex gap-2">
          {(campaign.Status === "draft" || campaign.Status === "scheduled") && (
            <button
              type="button"
              disabled={acting}
              onClick={() => runAction(launchCampaignAction, "Campaña lanzada")}
              className="rounded-lg px-3 py-1.5 text-xs font-medium text-white disabled:opacity-50"
              style={{ backgroundColor: "var(--color-primary)", color: "var(--color-on-primary, white)" }}
            >
              {"Lanzar"}
            </button>
          )}
          {(campaign.Status === "running" || campaign.Status === "scheduled") && (
            <button
              type="button"
              disabled={acting}
              onClick={() => runAction(pauseCampaignAction, "Campaña pausada")}
              className="rounded-lg bg-amber-50 px-3 py-1.5 text-xs font-medium text-amber-700 disabled:opacity-50"
            >
              {"Pausar"}
            </button>
          )}
          {campaign.Status === "paused" && (
            <button
              type="button"
              disabled={acting}
              onClick={() => runAction(resumeCampaignAction, "Campaña reanudada")}
              className="rounded-lg bg-green-50 px-3 py-1.5 text-xs font-medium text-green-700 disabled:opacity-50"
            >
              {"Reanudar"}
            </button>
          )}
          {campaign.Status !== "completed" && campaign.Status !== "cancelled" && (
            <button
              type="button"
              disabled={acting}
              onClick={() => runAction(cancelCampaignAction, "Campaña cancelada")}
              className="rounded-lg bg-red-50 px-3 py-1.5 text-xs font-medium text-red-600 disabled:opacity-50"
            >
              {"Cancelar"}
            </button>
          )}
        </div>
      </div>

      {campaign.ErrorMessage && (
        <div className="rounded-lg border border-red-200 bg-red-50 p-3 text-xs text-red-700 dark:border-red-900/40 dark:bg-red-900/10 dark:text-red-300">
          {campaign.ErrorMessage}
        </div>
      )}

      <div className="grid grid-cols-2 gap-2 sm:grid-cols-5">
        {[
          { label: "Audiencia", value: campaign.AudienceCount },
          { label: "Enviados", value: campaign.SentCount },
          { label: "En cola", value: campaign.QueuedCount },
          { label: "Respondieron", value: campaign.RepliedCount },
          { label: "Fallidos", value: campaign.FailedCount },
        ].map((metric) => (
          <div
            key={metric.label}
            className="rounded-lg border border-gray-200 p-2 text-center dark:border-gray-600"
          >
            <p className="text-lg font-bold text-gray-900 dark:text-white">{metric.value}</p>
            <p className="text-[10px] uppercase tracking-wide text-gray-500">{metric.label}</p>
          </div>
        ))}
      </div>

      <div className="flex gap-1 border-b border-gray-200 dark:border-gray-600">
        <button
          type="button"
          onClick={() => setTab("results")}
          className={`px-4 py-2 text-sm font-medium ${
            tab === "results"
              ? "border-b-2 border-[var(--color-primary)] text-[var(--color-primary)]"
              : "text-gray-500"
          }`}
        >
          {"Resultados del envío"}
        </button>
        <button
          type="button"
          onClick={() => setTab("chat")}
          className={`px-4 py-2 text-sm font-medium ${
            tab === "chat"
              ? "border-b-2 border-[var(--color-primary)] text-[var(--color-primary)]"
              : "text-gray-500"
          }`}
        >
          {"Conversaciones"}
        </button>
        <button
          type="button"
          onClick={() => setTab("config")}
          className={`px-4 py-2 text-sm font-medium ${
            tab === "config"
              ? "border-b-2 border-[var(--color-primary)] text-[var(--color-primary)]"
              : "text-gray-500"
          }`}
        >
          {"Configuración"}
        </button>
      </div>

      {tab === "config" ? (
        <div className="space-y-3">
          {configGroups.map((group) => (
            <div
              key={group.title}
              className="rounded-lg border border-gray-200 p-3 dark:border-gray-600"
            >
              <p className="mb-2 text-xs font-semibold text-gray-900 dark:text-white">
                {group.title}
              </p>
              <dl className="grid gap-x-6 gap-y-1.5 sm:grid-cols-2">
                {group.rows.map((row) => (
                  <div key={row.label} className="flex justify-between gap-3 text-xs">
                    <dt className="text-gray-500">{row.label}</dt>
                    <dd className="text-right font-medium text-gray-800 dark:text-gray-100">
                      {row.value}
                    </dd>
                  </div>
                ))}
              </dl>
            </div>
          ))}
        </div>
      ) : tab === "results" ? (
        loading ? (
          <p className="py-6 text-center text-sm text-gray-500">{"Cargando..."}</p>
        ) : sends.length === 0 ? (
          <p className="py-6 text-center text-sm text-gray-500">
            {"Todavía no hay envíos. Lanzá la campaña para empezar."}
          </p>
        ) : (
          <div className="overflow-x-auto rounded-lg border border-gray-200 dark:border-gray-600">
            <table className="w-full">
              <thead style={{ backgroundColor: "var(--color-primary)" }}>
                <tr>
                  {["Cliente", "Celular", "Estado", "Enviado", "Detalle"].map((column) => (
                    <th
                      key={column}
                      style={{ color: "var(--color-on-primary, white)" }}
                      className="px-3 py-2 text-left text-[10px] font-semibold uppercase"
                    >
                      {column}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {sends.map((send) => (
                  <tr
                    key={send.ID}
                    className="border-b border-gray-100 text-xs dark:border-gray-700"
                  >
                    <td className="px-3 py-2 text-gray-700 dark:text-gray-200">
                      {send.ClientName || "-"}
                    </td>
                    <td className="px-3 py-2 text-gray-500">{send.Phone}</td>
                    <td className="px-3 py-2 text-gray-700 dark:text-gray-200">
                      {SEND_STATUS_LABEL[send.Status] || send.Status}
                    </td>
                    <td className="px-3 py-2 text-gray-500">
                      {send.SentAt ? new Date(send.SentAt).toLocaleString("es-CO") : "-"}
                    </td>
                    <td className="px-3 py-2 text-red-500">{send.ErrorMessage || ""}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )
      ) : (
        <WhatsAppConversations businessId={businessId} campaignId={campaign.ID} />
      )}
    </div>
  );
}
