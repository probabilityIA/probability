"use client";

import { useCallback, useEffect, useState } from "react";
import { Modal } from "@/shared/ui/modal";
import { WhatsappTemplate } from "../../domain/scheduled-types";
import { Campaign } from "../../domain/campaign-types";
import { listCampaignsAction, deleteCampaignAction } from "../../infra/actions/campaigns";
import {
  getTemplateVariablesAction,
  listTemplatesAction,
} from "../../infra/actions/whatsapp-templates";
import { useToast } from "@/shared/providers/toast-provider";
import { CampaignForm } from "./CampaignForm";
import { CampaignDetail } from "./CampaignDetail";
import { CampaignRules } from "./CampaignRules";
import { TemplateBuilder } from "./TemplateBuilder";

interface CampaignsSectionProps {
  businessId?: number;
}

const STATUS_LABEL: Record<string, { label: string; className: string }> = {
  draft: { label: "Borrador", className: "bg-gray-100 text-gray-600" },
  scheduled: { label: "Programada", className: "bg-blue-100 text-blue-700" },
  running: { label: "Enviando", className: "bg-green-100 text-green-700" },
  paused: { label: "Pausada", className: "bg-amber-100 text-amber-700" },
  completed: { label: "Terminada", className: "bg-purple-100 text-purple-700" },
  cancelled: { label: "Cancelada", className: "bg-red-100 text-red-700" },
};

export function CampaignsSection({ businessId }: CampaignsSectionProps) {
  const { showToast } = useToast();

  const [campaigns, setCampaigns] = useState<Campaign[]>([]);
  const [templates, setTemplates] = useState<WhatsappTemplate[]>([]);
  const [variableCatalog, setVariableCatalog] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(true);

  const [selected, setSelected] = useState<Campaign | null>(null);
  const [editing, setEditing] = useState<Campaign | null>(null);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [isTemplateOpen, setIsTemplateOpen] = useState(false);

  const fetchAll = useCallback(async () => {
    setLoading(true);
    const [campaignsResult, templatesResult, catalogResult] = await Promise.all([
      listCampaignsAction(businessId, undefined, 1, 50),
      listTemplatesAction(businessId, "campaign", undefined, 1, 100),
      getTemplateVariablesAction(businessId),
    ]);
    setLoading(false);

    if (campaignsResult.success) setCampaigns(campaignsResult.data);
    if (templatesResult.success) setTemplates(templatesResult.data);
    if (catalogResult.success) setVariableCatalog(catalogResult.data);
  }, [businessId]);

  useEffect(() => {
    fetchAll();
  }, [fetchAll]);

  useEffect(() => {
    if (!selected) return;
    const fresh = campaigns.find((item) => item.ID === selected.ID);
    if (fresh) setSelected(fresh);
  }, [campaigns, selected]);

  const handleDelete = async (campaign: Campaign) => {
    const result = await deleteCampaignAction(campaign.ID, businessId);
    if (!result.success) {
      showToast(result.error || "No se pudo eliminar", "error");
      return;
    }
    showToast("Campaña eliminada", "success");
    if (selected?.ID === campaign.ID) setSelected(null);
    fetchAll();
  };

  if (selected) {
    return (
      <CampaignDetail
        campaign={selected}
        businessId={businessId}
        onChanged={fetchAll}
        onBack={() => setSelected(null)}
      />
    );
  }

  return (
    <div className="space-y-5">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h3 className="text-sm font-semibold text-gray-900 dark:text-white">
            {"Campañas de marketing"}
          </h3>
          <p className="text-xs text-gray-500 dark:text-gray-400">
            {
              "Mensajes comerciales que salen desde tu propio número a los clientes que tenés cargados. No dependen de una orden ni de un evento."
            }
          </p>
        </div>
        <div className="flex gap-2">
          <button
            type="button"
            onClick={() => setIsTemplateOpen(true)}
            className="rounded-lg border border-[var(--color-primary)]/30 bg-[var(--color-primary)]/10 px-3 py-1.5 text-xs font-medium text-[var(--color-primary)]"
          >
            {"Nueva plantilla"}
          </button>
          <button
            type="button"
            onClick={() => {
              setEditing(null);
              setIsFormOpen(true);
            }}
            className="rounded-lg px-3 py-1.5 text-xs font-medium text-white"
            style={{ backgroundColor: "var(--color-primary)", color: "var(--color-on-primary, white)" }}
          >
            {"Nueva campaña"}
          </button>
        </div>
      </div>

      {loading ? (
        <p className="py-8 text-center text-sm text-gray-500">{"Cargando..."}</p>
      ) : campaigns.length === 0 ? (
        <div className="rounded-lg border border-dashed border-gray-300 p-6 text-center dark:border-gray-600">
          <p className="text-sm text-gray-500 dark:text-gray-400">
            {"Todavía no creaste ninguna campaña."}
          </p>
          <p className="mt-1 text-xs text-gray-400">
            {"Necesitás una plantilla de campaña aprobada por Meta y tu número propio conectado."}
          </p>
        </div>
      ) : (
        <div className="space-y-2">
          {campaigns.map((campaign) => {
            const status = STATUS_LABEL[campaign.Status] || STATUS_LABEL.draft;
            return (
              <div
                key={campaign.ID}
                className="flex flex-wrap items-center gap-3 rounded-lg border border-gray-200 p-3 dark:border-gray-600"
              >
                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-2">
                    <p className="truncate text-sm font-medium text-gray-900 dark:text-white">
                      {campaign.Name}
                    </p>
                    <span
                      className={`rounded-full px-2 py-0.5 text-[10px] font-semibold ${status.className}`}
                    >
                      {status.label}
                    </span>
                  </div>
                  <p className="text-xs text-gray-500 dark:text-gray-400">
                    {`${campaign.SentCount} de ${campaign.AudienceCount} enviados · ${campaign.RepliedCount} respondieron`}
                  </p>
                </div>

                <div className="flex gap-2">
                  <button
                    type="button"
                    onClick={() => setSelected(campaign)}
                    className="rounded-md bg-gray-100 px-2 py-1 text-xs text-gray-600 hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-300"
                  >
                    {"Ver"}
                  </button>
                  {(campaign.Status === "draft" || campaign.Status === "scheduled") && (
                    <button
                      type="button"
                      onClick={() => {
                        setEditing(campaign);
                        setIsFormOpen(true);
                      }}
                      className="rounded-md bg-gray-100 px-2 py-1 text-xs text-gray-600 hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-300"
                    >
                      {"Editar"}
                    </button>
                  )}
                  {campaign.Status !== "running" && (
                    <button
                      type="button"
                      onClick={() => handleDelete(campaign)}
                      className="rounded-md bg-red-50 px-2 py-1 text-xs text-red-600 hover:bg-red-100"
                    >
                      {"Eliminar"}
                    </button>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      )}

      <CampaignRules />

      <Modal
        isOpen={isFormOpen}
        onClose={() => setIsFormOpen(false)}
        title={editing ? "Editar campaña" : "Nueva campaña"}
        size="4xl"
      >
        {isFormOpen && (
          <CampaignForm
            businessId={businessId}
            templates={templates}
            campaign={editing}
            onSuccess={() => {
              setIsFormOpen(false);
              fetchAll();
            }}
            onCancel={() => setIsFormOpen(false)}
          />
        )}
      </Modal>

      <Modal
        isOpen={isTemplateOpen}
        onClose={() => setIsTemplateOpen(false)}
        title="Nueva plantilla de campaña"
        size="4xl"
      >
        {isTemplateOpen && (
          <TemplateBuilder
            businessId={businessId}
            variableCatalog={variableCatalog}
            scope="campaign"
            onSuccess={() => {
              setIsTemplateOpen(false);
              fetchAll();
            }}
            onCancel={() => setIsTemplateOpen(false)}
          />
        )}
      </Modal>
    </div>
  );
}
