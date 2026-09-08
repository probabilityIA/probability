"use client";

import { useEffect, useState } from "react";
import { Modal } from "@/shared/ui/modal";
import { TemplatePreview } from "../../domain/scheduled-types";
import { previewTemplatesByEventAction } from "../../infra/actions/template-preview";

interface TemplatePreviewModalProps {
  isOpen: boolean;
  eventCode: string;
  eventName: string;
  businessId?: number;
  onClose: () => void;
}

const STATUS_STYLE: Record<string, string> = {
  APPROVED: "bg-green-100 text-green-700",
  PENDING: "bg-amber-100 text-amber-700",
  REJECTED: "bg-red-100 text-red-700",
  PAUSED: "bg-orange-100 text-orange-700",
  DISABLED: "bg-gray-200 text-gray-600",
};

export function TemplatePreviewModal({
  isOpen,
  eventCode,
  eventName,
  businessId,
  onClose,
}: TemplatePreviewModalProps) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [previews, setPreviews] = useState<TemplatePreview[]>([]);

  useEffect(() => {
    if (!isOpen || !eventCode) return;

    let cancelled = false;

    const load = async () => {
      setLoading(true);
      setError("");
      const result = await previewTemplatesByEventAction(eventCode, businessId);
      if (cancelled) return;
      if (!result.success) {
        setError(result.error || "No se pudo consultar Meta");
        setPreviews([]);
      } else {
        setPreviews(result.data);
      }
      setLoading(false);
    };

    load();

    return () => {
      cancelled = true;
    };
  }, [isOpen, eventCode, businessId]);

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={`Plantilla — ${eventName}`} size="2xl">
      <div className="space-y-4">
        <p className="text-xs text-gray-500 dark:text-gray-400">
          {"Texto tal como está aprobado hoy en Meta. Se consulta en vivo, no es una copia local."}
        </p>

        {loading && (
          <p className="py-6 text-center text-sm text-gray-500">{"Consultando a Meta..."}</p>
        )}

        {!loading && error && (
          <p className="rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-700">
            {error}
          </p>
        )}

        {!loading && !error && previews.length === 0 && (
          <p className="rounded-md border border-dashed border-gray-300 p-4 text-sm text-gray-500">
            {"Este evento todavía no tiene una plantilla asociada."}
          </p>
        )}

        {!loading &&
          previews.map((preview) => (
            <div
              key={preview.name}
              className="rounded-lg border border-gray-200 dark:border-gray-600 p-4"
            >
              <div className="mb-2 flex items-start justify-between gap-3">
                <div className="min-w-0">
                  <p className="truncate text-sm font-semibold">{preview.name}</p>
                  {preview.condition && (
                    <p className="text-xs text-gray-500">{preview.condition}</p>
                  )}
                </div>
                <span
                  className={`shrink-0 rounded-full px-2 py-1 text-[11px] ${
                    STATUS_STYLE[preview.status] || "bg-gray-100 text-gray-600"
                  }`}
                >
                  {preview.found_in_whatsapp ? preview.status : "No existe en Meta"}
                </span>
              </div>

              {preview.found_in_whatsapp ? (
                <>
                  <div className="rounded-lg bg-[#e7ffdb] p-3 text-sm whitespace-pre-wrap text-gray-800">
                    {preview.header && (
                      <p className="mb-1 font-semibold">{preview.header}</p>
                    )}
                    <p>{preview.body}</p>
                    {preview.footer && (
                      <p className="mt-2 text-xs text-gray-500">{preview.footer}</p>
                    )}
                  </div>

                  {preview.buttons && preview.buttons.length > 0 && (
                    <div className="mt-2 flex flex-wrap gap-2">
                      {preview.buttons.map((label) => (
                        <span
                          key={label}
                          className="rounded border border-gray-300 px-2 py-1 text-xs text-gray-600"
                        >
                          {label}
                        </span>
                      ))}
                    </div>
                  )}

                  <p className="mt-2 text-xs text-gray-500">
                    {`${preview.variables} variable(s) · ${preview.category} · ${preview.language}`}
                  </p>

                  {preview.status === "REJECTED" && preview.rejected_reason && (
                    <p className="mt-1 text-xs text-red-600">{preview.rejected_reason}</p>
                  )}
                </>
              ) : (
                <p className="text-xs text-gray-500">
                  {"Meta no tiene una plantilla con este nombre. Si el evento se dispara, el mensaje no sale."}
                </p>
              )}
            </div>
          ))}
      </div>
    </Modal>
  );
}
