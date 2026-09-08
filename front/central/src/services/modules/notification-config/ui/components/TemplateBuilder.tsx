"use client";

import { useState } from "react";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import { Button } from "@/shared/ui/button";
import { useToast } from "@/shared/providers/toast-provider";
import { CreateTemplateDTO, TemplateCategory, TemplateScope } from "../../domain/scheduled-types";
import { createTemplateAction } from "../../infra/actions/whatsapp-templates";
import { WhatsAppBubblePreview } from "./WhatsAppBubblePreview";

interface TemplateBuilderProps {
  businessId?: number;
  variableCatalog: Record<string, string>;
  scope?: TemplateScope;
  onSuccess: () => void;
  onCancel: () => void;
}

interface VariableRow {
  position: number;
  source: string;
  fallback: string;
}

const PLACEHOLDER = /\{\{(\d+)\}\}/g;

const SAMPLE_VALUES: Record<string, string> = {
  "customer.first_name": "Ana",
  "customer.full_name": "Ana Ramirez",
  "customer.days_inactive": "45",
  "customer.last_product": "Camiseta blanca",
  "customer.total_orders": "3",
  "business.name": "Mi Tienda",
  "sender.name": "Isabel Rojas",
  "campaign.name": "Ruta 30",
};

const STARTERS: Array<{ label: string; body: string; sources: string[] }> = [
  {
    label: "Reactivar cliente dormido",
    body: "Hola {{1}}, hace {{2}} dias que no te vemos por {{3}}. Tenemos novedades que te pueden gustar.",
    sources: ["customer.first_name", "customer.days_inactive", "business.name"],
  },
  {
    label: "Recordar ultimo producto",
    body: "Hola {{1}}, vimos que te llevaste {{2}}. Si te gusto, te esperamos de nuevo en {{3}}.",
    sources: ["customer.first_name", "customer.last_product", "business.name"],
  },
  {
    label: "Agradecer al cliente frecuente",
    body: "Hola {{1}}, gracias por tus {{2}} pedidos en {{3}}. Queriamos saludarte.",
    sources: ["customer.first_name", "customer.total_orders", "business.name"],
  },
];

const CAMPAIGN_STARTERS: Array<{ label: string; body: string; sources: string[] }> = [
  {
    label: "Ruta 30 - control DIAN",
    body:
      "Hola {{1}}, mucho gusto, soy {{2}} de Siigo. Estamos acompanando a las empresas frente a los controles de la DIAN a la facturacion electronica. Para orientarte, contame: como estas facturando hoy en tu negocio?",
    sources: ["customer.first_name", "sender.name"],
  },
  {
    label: "Ruta 30 - ahorro de tiempo",
    body:
      "Hola {{1}}, mucho gusto, soy {{2}}. Ayudo a empresas a reducir el tiempo que gastan en facturacion electronica, control de inventario y cobros usando Siigo. Hoy tenes esos procesos en una sola plataforma o usas varios sistemas?",
    sources: ["customer.first_name", "sender.name"],
  },
  {
    label: "Ruta 30 - validacion de sistema",
    body:
      "Hola {{1}}, mucho gusto, soy {{2}}. Estamos haciendo una campana de validacion para que las empresas no queden expuestas a los nuevos controles de la DIAN. Que sistema usas hoy para facturar en tu empresa?",
    sources: ["customer.first_name", "sender.name"],
  },
];

const SOURCE_ORDER = [
  "customer.first_name",
  "customer.full_name",
  "sender.name",
  "campaign.name",
  "customer.days_inactive",
  "customer.last_product",
  "customer.total_orders",
  "business.name",
];

function orderSources(catalog: Record<string, string>): Array<[string, string]> {
  const known = SOURCE_ORDER.filter((key) => catalog[key]).map(
    (key) => [key, catalog[key]] as [string, string],
  );
  const rest = Object.entries(catalog).filter(([key]) => !SOURCE_ORDER.includes(key));
  return [...known, ...rest];
}

function countPlaceholders(body: string): number {
  const found = new Set<string>();
  for (const match of body.matchAll(PLACEHOLDER)) {
    found.add(match[1]);
  }
  return found.size;
}

function slugify(value: string): string {
  return value
    .toLowerCase()
    .normalize("NFD")
    .replace(/[\u0300-\u036f]/g, "")
    .replace(/[^a-z0-9]+/g, "_")
    .replace(/^_+|_+$/g, "")
    .slice(0, 60);
}

export function TemplateBuilder({
  businessId,
  variableCatalog,
  scope = "scheduled",
  onSuccess,
  onCancel,
}: TemplateBuilderProps) {
  const { showToast } = useToast();
  const [loading, setLoading] = useState(false);

  const [displayName, setDisplayName] = useState("");
  const [category, setCategory] = useState<TemplateCategory>("MARKETING");
  const [headerText, setHeaderText] = useState("");
  const [bodyText, setBodyText] = useState("");
  const [footerText, setFooterText] = useState("");
  const [variables, setVariables] = useState<VariableRow[]>([]);

  const sources = orderSources(variableCatalog);

  const sampleValues: Record<number, string> = {};
  for (const variable of variables) {
    sampleValues[variable.position] =
      variable.fallback.trim() || SAMPLE_VALUES[variable.source] || "ejemplo";
  }

  const previewButtons =
    category === "MARKETING" ? ["Dejar de recibir"] : [];

  const applyStarter = (starter: (typeof STARTERS)[number]) => {
    setBodyText(starter.body);
    setVariables(
      starter.sources.map((source, index) => ({
        position: index + 1,
        source,
        fallback: "",
      })),
    );
  };
  const placeholders = countPlaceholders(bodyText);
  const templateName = slugify(displayName);
  const mismatch = placeholders !== variables.length;

  const insertVariable = () => {
    const next = variables.length + 1;
    setVariables([...variables, { position: next, source: sources[0]?.[0] || "", fallback: "" }]);
    setBodyText(`${bodyText}{{${next}}}`);
  };

  const updateVariable = (index: number, patch: Partial<VariableRow>) => {
    setVariables(variables.map((item, i) => (i === index ? { ...item, ...patch } : item)));
  };

  const removeVariable = (index: number) => {
    const removed = variables[index];
    setVariables(
      variables
        .filter((_, i) => i !== index)
        .map((item, i) => ({ ...item, position: i + 1 })),
    );
    setBodyText(bodyText.replace(`{{${removed.position}}}`, ""));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!templateName) {
      showToast("El nombre de la plantilla es obligatorio", "error");
      return;
    }
    if (!bodyText.trim()) {
      showToast("El cuerpo del mensaje es obligatorio", "error");
      return;
    }
    if (mismatch) {
      showToast(
        `El cuerpo usa ${placeholders} variable(s) y declaraste ${variables.length}`,
        "error",
      );
      return;
    }
    if (variables.some((item) => !item.source)) {
      showToast("Toda variable necesita un dato asociado", "error");
      return;
    }

    const dto: CreateTemplateDTO = {
      scope,
      name: templateName,
      language: "es",
      category,
      header_text: headerText.trim(),
      body_text: bodyText.trim(),
      footer_text: footerText.trim(),
      variables: variables.map((item) => ({
        position: item.position,
        source: item.source,
        label: variableCatalog[item.source] || "",
        fallback: item.fallback.trim(),
      })),
    };

    setLoading(true);
    const result = await createTemplateAction(dto, businessId);
    setLoading(false);

    if (!result.success) {
      showToast(result.error || "No se pudo crear la plantilla", "error");
      return;
    }

    showToast("Plantilla enviada a Meta. Queda en revisión.", "success");
    onSuccess();
  };

  return (
    <form onSubmit={handleSubmit} className="grid gap-6 lg:grid-cols-[minmax(0,1fr)_320px]">
      <div className="space-y-4">
      <div>
        <p className="mb-2 text-xs font-medium text-gray-600 dark:text-gray-300">
          {"\u00bfNo sab\u00e9s por d\u00f3nde empezar? Arranc\u00e1 de un ejemplo:"}
        </p>
        <div className="flex flex-wrap gap-2">
          {(scope === "campaign" ? CAMPAIGN_STARTERS : STARTERS).map((starter) => (
            <button
              key={starter.label}
              type="button"
              onClick={() => applyStarter(starter)}
              className="rounded-full border border-[var(--color-primary)]/30 bg-[var(--color-primary)]/10 px-3 py-1 text-xs text-[var(--color-primary)] hover:bg-[var(--color-primary)]/20"
            >
              {starter.label}
            </button>
          ))}
        </div>
      </div>

      <div>
        <Label htmlFor="template-name">{"Nombre de la plantilla"}</Label>
        <p className="mb-1 text-xs text-gray-500">
          {
            "Solo lo ves vos y Meta, el cliente nunca lo lee. Sirve para identificarla despu\u00e9s. Ojo: una vez enviada a Meta el nombre no se puede cambiar."
          }
        </p>
        <Input
          id="template-name"
          value={displayName}
          onChange={(e) => setDisplayName(e.target.value)}
          placeholder="Vuelve pronto"
        />
        {templateName && (
          <p className="mt-1 text-xs text-gray-500">
            {"Meta la registrará como: "}
            <code className="rounded bg-gray-100 px-1">{templateName}</code>
          </p>
        )}
      </div>

      <div>
        <Label htmlFor="template-category">{"Categoría"}</Label>
        <select
          id="template-category"
          value={category}
          onChange={(e) => setCategory(e.target.value as TemplateCategory)}
          className="w-full rounded-md border border-gray-300 px-3 py-2 text-sm"
        >
          <option value="MARKETING">{"Marketing (promocional)"}</option>
          <option value="UTILITY">{"Utilidad (transaccional)"}</option>
        </select>
        {category === "MARKETING" && (
          <p className="mt-1 text-xs text-amber-600">
            {
              "Meta exige botón de baja en las plantillas de marketing. Se agrega automáticamente."
            }
          </p>
        )}
      </div>

      <div>
        <Label htmlFor="template-header">{"Encabezado (opcional, máx. 60)"}</Label>
        <Input
          id="template-header"
          value={headerText}
          maxLength={60}
          onChange={(e) => setHeaderText(e.target.value)}
        />
      </div>

      <div>
        <div className="flex items-center justify-between">
          <Label htmlFor="template-body">{"Mensaje"}</Label>
          <button
            type="button"
            onClick={insertVariable}
            disabled={sources.length === 0}
            className="text-xs font-medium text-[var(--color-primary)] hover:underline disabled:opacity-40"
          >
            {"+ Insertar variable"}
          </button>
        </div>
        <p className="mb-1 text-xs text-gray-500">
          {
            "Es lo que le llega al cliente. Us\u00e1 \u0022Insertar variable\u0022 para meter datos que cambian en cada env\u00edo (su nombre, los d\u00edas sin comprar). No escribas {{1}} a mano."
          }
        </p>
        <textarea
          id="template-body"
          value={bodyText}
          maxLength={1024}
          rows={5}
          onChange={(e) => setBodyText(e.target.value)}
          className="w-full rounded-md border border-gray-300 px-3 py-2 text-sm"
          placeholder="Hola {{1}}, hace {{2}} días que no te vemos."
        />
        <p className="mt-1 text-xs text-gray-500">
          {`${bodyText.length}/1024 - variables en el texto: ${placeholders}`}
        </p>
      </div>

      {variables.length > 0 && (
        <div className="space-y-2 rounded-md border border-gray-200 p-3">
          <p className="text-sm font-medium">{"Qué dato va en cada variable"}</p>
          {variables.map((variable, index) => (
            <div key={variable.position} className="flex items-center gap-2">
              <code className="rounded bg-gray-100 px-2 py-1 text-xs">
                {`{{${variable.position}}}`}
              </code>
              <select
                value={variable.source}
                onChange={(e) => updateVariable(index, { source: e.target.value })}
                className="flex-1 rounded-md border border-gray-300 px-2 py-1 text-sm"
              >
                {sources.map(([key, label]) => (
                  <option key={key} value={key}>
                    {label}
                  </option>
                ))}
              </select>
              <Input
                value={variable.fallback}
                onChange={(e) => updateVariable(index, { fallback: e.target.value })}
                placeholder="Si falta el dato"
                className="w-40"
              />
              <button
                type="button"
                onClick={() => removeVariable(index)}
                className="text-xs text-red-500 hover:underline"
              >
                {"Quitar"}
              </button>
            </div>
          ))}
        </div>
      )}

      <div>
        <Label htmlFor="template-footer">{"Pie (opcional, máx. 60)"}</Label>
        <Input
          id="template-footer"
          value={footerText}
          maxLength={60}
          onChange={(e) => setFooterText(e.target.value)}
        />
      </div>

      {mismatch && (
        <p className="text-sm text-red-600">
          {`El cuerpo usa ${placeholders} variable(s) y declaraste ${variables.length}. Meta rechaza la plantilla si no coinciden.`}
        </p>
      )}

      <p className="text-xs text-gray-500">
        {
          "Al crearla se env\u00eda a Meta para revisi\u00f3n. Suele tardar entre unos minutos y 24 horas. Reci\u00e9n cuando quede aprobada se puede usar en una regla programada."
        }
      </p>

      <div className="flex justify-end gap-2 pt-2">
        <Button type="button" variant="outline" onClick={onCancel} disabled={loading}>
          {"Cancelar"}
        </Button>
        <Button type="submit" disabled={loading || mismatch}>
          {loading ? "Enviando..." : "Crear y enviar a Meta"}
        </Button>
      </div>
      </div>

      <div className="lg:sticky lg:top-2 lg:self-start">
        <WhatsAppBubblePreview
          header={headerText}
          body={bodyText}
          footer={footerText}
          buttons={previewButtons}
          sampleValues={sampleValues}
        />
        <p className="mt-2 text-[11px] text-gray-500">
          {"Los datos son de ejemplo. En el env\u00edo real se reemplazan por los del cliente."}
        </p>
      </div>
    </form>
  );
}
