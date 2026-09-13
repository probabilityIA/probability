"use client";

import { TemplateHeaderType, TemplateVariable } from "../../domain/scheduled-types";

const PLACEHOLDER = /\{\{(\d+)\}\}/g;

export const SAMPLE_VALUES: Record<string, string> = {
  "customer.first_name": "Ana",
  "customer.full_name": "Ana Ramirez",
  "customer.days_inactive": "45",
  "customer.last_product": "Camiseta blanca",
  "customer.total_orders": "3",
  "business.name": "Mi Tienda",
  "sender.name": "Tu nombre",
  "campaign.name": "Ruta 30",
};

export function fillPlaceholders(
  body: string,
  variables: TemplateVariable[] | null,
): string {
  const rows = variables ?? [];

  return body.replace(PLACEHOLDER, (_, position: string) => {
    const variable = rows.find((item) => item.Position === Number(position));
    if (!variable) return "···";
    return (variable.Fallback || "").trim() || SAMPLE_VALUES[variable.Source] || "ejemplo";
  });
}

interface TemplateBubbleProps {
  headerType?: TemplateHeaderType;
  headerMediaURL?: string;
  headerText?: string;
  bodyText: string;
  footerText?: string;
  buttons?: string[];
  className?: string;
  compact?: boolean;
}

export function TemplateBubble({
  headerType = "TEXT",
  headerMediaURL = "",
  headerText = "",
  bodyText,
  footerText = "",
  buttons = [],
  className = "max-w-[85%] self-start",
  compact = false,
}: TemplateBubbleProps) {
  return (
    <div
      className={`${className} rounded-[10px] rounded-bl-[3px] shadow-sm ${
        compact ? "p-2" : "p-3"
      }`}
      style={{ backgroundColor: "#dcf8c6" }}
    >
      {headerType === "IMAGE" && headerMediaURL && (
        <img
          src={headerMediaURL}
          alt="Encabezado"
          className={`mb-1.5 w-full rounded-[5px] object-cover ${
            compact ? "max-h-20" : "max-h-40"
          }`}
        />
      )}
      {headerType === "TEXT" && headerText.trim() && (
        <p
          className={`mb-0.5 font-bold text-[#111b21] ${compact ? "text-[12px]" : "text-sm"}`}
        >
          {headerText}
        </p>
      )}

      <p
        className={`whitespace-pre-wrap text-[#111b21] ${
          compact ? "text-[12px] leading-snug" : "text-sm leading-relaxed"
        }`}
      >
        {bodyText}
      </p>

      {footerText.trim() && (
        <p className={`mt-0.5 text-[#667781] ${compact ? "text-[10px]" : "text-[12px]"}`}>
          {footerText}
        </p>
      )}

      <p
        className={`mt-0.5 text-right text-[#667781] ${compact ? "text-[9px]" : "text-[10px]"}`}
      >
        {new Date().toLocaleTimeString("es-CO", { hour: "2-digit", minute: "2-digit" })}
      </p>

      {buttons.map((text) => (
        <div
          key={text}
          className={`border-t text-center font-medium text-[#00a5f4] ${
            compact ? "mt-1 pt-1 text-[11px]" : "mt-2 pt-2 text-[13px]"
          }`}
          style={{ borderColor: "rgba(17,27,33,0.12)" }}
        >
          {text}
        </div>
      ))}
    </div>
  );
}
