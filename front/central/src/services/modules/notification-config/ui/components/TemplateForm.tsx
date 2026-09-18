"use client";

import { useState } from "react";
import { SAMPLE_VALUES, TemplateBubble } from "./TemplateBubble";
import { useToast } from "@/shared/providers/toast-provider";
import {
  CreateTemplateDTO,
  TemplateCategory,
  TemplateHeaderType,
  TemplateScope,
  UpdateTemplateDTO,
  WhatsappTemplate,
} from "../../domain/scheduled-types";
import {
  createTemplateAction,
  updateTemplateAction,
  uploadTemplateMediaAction,
} from "../../infra/actions/whatsapp-templates";
import { metaRuleError } from "../../domain/template-rules";

interface TemplateFormProps {
  businessId?: number;
  variableCatalog: Record<string, string>;
  scope?: TemplateScope;
  template?: WhatsappTemplate | null;
  linkedButtonTexts?: string[];
  asFlowResponse?: boolean;
  onSuccess: (created?: WhatsappTemplate) => void;
  onCancel: () => void;
}

interface VariableRow {
  position: number;
  source: string;
  fallback: string;
}

interface ButtonRow {
  text: string;
}

const PLACEHOLDER = /\{\{(\d+)\}\}/g;

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

const MAX_BODY = 1024;
const MAX_HEADER = 60;
const MAX_FOOTER = 60;
const MAX_BUTTONS = 3;
const MAX_BUTTON_TEXT = 25;
const OPT_OUT_TEXT = "Dejar de recibir";

const inputCls =
  "w-full rounded-lg border px-3 py-2.5 text-sm text-gray-900 outline-none transition-colors placeholder:text-gray-400 focus:border-[var(--color-primary)] dark:bg-gray-800 dark:text-white";

const VARIABLE_CLASS: Record<string, string> = {
  "customer.first_name": "Nombre",
  "customer.full_name": "Nombre",
  "business.name": "Nombre",
  "sender.name": "Nombre",
  "campaign.name": "Nombre",
  "customer.days_inactive": "Cantidad",
  "customer.total_orders": "Cantidad",
  "customer.last_product": "Texto",
};

const CLASS_ORDER = ["Nombre", "Cantidad", "Texto"];

const FLOW_BLOCKED_SOURCES = ["sender.name", "campaign.name"];

function groupSources(
  sources: Array<[string, string]>,
): Array<[string, Array<[string, string]>]> {
  const groups = new Map<string, Array<[string, string]>>();

  for (const entry of sources) {
    const clase = VARIABLE_CLASS[entry[0]] || "Texto";
    const current = groups.get(clase) || [];
    current.push(entry);
    groups.set(clase, current);
  }

  return CLASS_ORDER.filter((clase) => groups.has(clase)).map(
    (clase) => [clase, groups.get(clase) as Array<[string, string]>],
  );
}

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
    .replace(/[̀-ͯ]/g, "")
    .replace(/[^a-z0-9]+/g, "_")
    .replace(/^_+|_+$/g, "")
    .slice(0, 60);
}

function initialButtons(template?: WhatsappTemplate | null): ButtonRow[] {
  if (!template?.Buttons) return [];

  return template.Buttons.map((button) => (button.Text ?? "").trim())
    .filter((text) => text && text.toLowerCase() !== OPT_OUT_TEXT.toLowerCase())
    .map((text) => ({ text }));
}

function initialVariables(template?: WhatsappTemplate | null): VariableRow[] {
  if (!template?.Variables) return [];
  return template.Variables.map((variable) => ({
    position: variable.Position,
    source: variable.Source,
    fallback: variable.Fallback || "",
  }));
}

export function TemplateForm({
  businessId,
  variableCatalog,
  scope = "scheduled",
  template = null,
  linkedButtonTexts = [],
  asFlowResponse = false,
  onSuccess,
  onCancel,
}: TemplateFormProps) {
  const { showToast } = useToast();
  const isEdit = Boolean(template);
  const lockedByReview = template?.Status === "pending";

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [displayName, setDisplayName] = useState(template?.Name ?? "");
  const [category, setCategory] = useState<TemplateCategory>(template?.Category ?? "UTILITY");
  const [headerText, setHeaderText] = useState(template?.HeaderText ?? "");
  const [headerType, setHeaderType] = useState<TemplateHeaderType>(
    template?.HeaderType === "IMAGE" ? "IMAGE" : "TEXT",
  );
  const [headerMediaURL, setHeaderMediaURL] = useState(template?.HeaderMediaURL ?? "");
  const [uploading, setUploading] = useState(false);
  const [bodyText, setBodyText] = useState(template?.BodyText ?? "");
  const [footerText, setFooterText] = useState(template?.FooterText ?? "");
  const [variables, setVariables] = useState<VariableRow[]>(initialVariables(template));
  const [buttons, setButtons] = useState<ButtonRow[]>(initialButtons(template));

  const sources = orderSources(variableCatalog).filter(
    ([key]) => !asFlowResponse || !FLOW_BLOCKED_SOURCES.includes(key),
  );
  const maxButtons = category === "MARKETING" ? MAX_BUTTONS - 1 : MAX_BUTTONS;

  const currentButtonTexts = new Set(
    buttons.map((item) => item.text.trim().toLowerCase()).filter(Boolean),
  );
  const brokenLinks = linkedButtonTexts.filter(
    (text) => !currentButtonTexts.has(text.trim().toLowerCase()),
  );

  const placeholders = countPlaceholders(bodyText);
  const slug = slugify(displayName);
  const prefix = businessId ? `${businessId}_` : "";
  const templateName = isEdit ? template?.Name ?? "" : slug ? `${prefix}${slug}` : "";
  const mismatch = placeholders !== variables.length;
  const ruleError = bodyText.trim()
    ? metaRuleError({
        body: bodyText.trim(),
        header: headerType === "IMAGE" ? "" : headerText.trim(),
        footer: footerText.trim(),
        buttons: buttons.map((item) => item.text.trim()).filter(Boolean),
        variables: placeholders,
      })
    : null;

  const sampleFor = (position: number): string => {
    const variable = variables.find((item) => item.position === position);
    if (!variable) return "···";
    return variable.fallback.trim() || SAMPLE_VALUES[variable.source] || "ejemplo";
  };

  const bodyPreview =
    bodyText.replace(PLACEHOLDER, (_, position: string) => sampleFor(Number(position))) ||
    "Tu mensaje aparecerá aquí…";

  const previewButtons = [
    ...buttons.map((item) => item.text.trim()).filter(Boolean),
    ...(category === "MARKETING" ? [OPT_OUT_TEXT] : []),
  ];

  const insertVariable = (source: string) => {
    const next = variables.length + 1;
    setVariables([...variables, { position: next, source, fallback: "" }]);
    const separator = bodyText && !/\s$/.test(bodyText) ? " " : "";
    setBodyText(`${bodyText}${separator}{{${next}}}`.slice(0, MAX_BODY));
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

  const handleImageChange = async (file: File | null) => {
    if (!file) return;

    const payload = new FormData();
    payload.append("file", file);

    setUploading(true);
    const result = await uploadTemplateMediaAction(payload, businessId);
    setUploading(false);

    if (!result.success || !result.url) {
      setError(result.error || "No se pudo subir la imagen");
      return;
    }

    setError("");
    setHeaderMediaURL(result.url);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (!templateName) {
      setError("El nombre de la plantilla es obligatorio");
      return;
    }
    if (!bodyText.trim()) {
      setError("El cuerpo del mensaje es obligatorio");
      return;
    }
    if (mismatch) {
      setError(
        `El cuerpo usa ${placeholders} variable(s) y declaraste ${variables.length}`,
      );
      return;
    }
    if (variables.some((item) => !item.source)) {
      setError("Toda variable necesita un dato asociado");
      return;
    }
    if (ruleError) {
      setError(ruleError);
      return;
    }
    if (headerType === "IMAGE" && !headerMediaURL) {
      setError("Subí la imagen del encabezado o volvé a encabezado de texto");
      return;
    }
    const filledButtons = buttons.map((item) => item.text.trim()).filter(Boolean);
    const uniqueButtons = new Set(filledButtons.map((text) => text.toLowerCase()));

    if (uniqueButtons.size !== filledButtons.length) {
      setError("Hay dos botones con el mismo texto");
      return;
    }
    if (filledButtons.some((text) => text.toLowerCase() === OPT_OUT_TEXT.toLowerCase())) {
      setError(`"${OPT_OUT_TEXT}" lo agrega Meta solo en las plantillas de marketing`);
      return;
    }

    const mappedButtons = filledButtons.map((text) => ({ type: "QUICK_REPLY", text }));

    const mappedVariables = variables.map((item) => ({
      position: item.position,
      source: item.source,
      label: variableCatalog[item.source] || "",
      fallback: item.fallback.trim(),
    }));

    setLoading(true);

    if (isEdit && template) {
      const dto: UpdateTemplateDTO = {
        category,
        header_text: headerType === "IMAGE" ? "" : headerText.trim(),
        header_type: headerType,
        header_media_url: headerType === "IMAGE" ? headerMediaURL : "",
        body_text: bodyText.trim(),
        footer_text: footerText.trim(),
        variables: mappedVariables,
        buttons: mappedButtons,
      };

      const result = await updateTemplateAction(template.ID, dto, businessId);
      setLoading(false);

      if (!result.success) {
        setError(result.error || "No se pudo editar la plantilla");
        return;
      }

      showToast(
        template.Status === "approved"
          ? "Cambios enviados a Meta. La plantilla vuelve a revisión."
          : "Cambios guardados en el borrador.",
        "success",
      );
      onSuccess();
      return;
    }

    const dto: CreateTemplateDTO = {
      scope,
      name: templateName,
      language: "es",
      category,
      header_text: headerType === "IMAGE" ? "" : headerText.trim(),
      header_type: headerType,
      header_media_url: headerType === "IMAGE" ? headerMediaURL : "",
      body_text: bodyText.trim(),
      footer_text: footerText.trim(),
      variables: mappedVariables,
      buttons: mappedButtons,
    };

    const result = await createTemplateAction(dto, businessId);
    setLoading(false);

    if (!result.success) {
      setError(result.error || "No se pudo crear la plantilla");
      return;
    }

    showToast("Plantilla guardada como borrador.", "success");
    onSuccess(result.data);
  };

  if (lockedByReview) {
    return (
      <div className="space-y-4 px-7 py-6">
        <p className="text-sm text-gray-600 dark:text-gray-300">
          {
            "Esta plantilla está en revisión de Meta. Hasta que responda no se puede editar: si la cambiás ahora, se pierde la revisión en curso."
          }
        </p>
        <div className="flex justify-end">
          <button
            type="button"
            onClick={onCancel}
            className="rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-800 transition-colors hover:bg-gray-100 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
          >
            {"Cerrar"}
          </button>
        </div>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit} className="flex min-h-0 flex-1 flex-col">
      <div className="flex min-h-0 flex-1 flex-col lg:flex-row">
        <div className="flex min-h-0 flex-[1.3] flex-col gap-5 overflow-y-auto border-b border-gray-200 px-7 py-6 lg:border-b-0 lg:border-r dark:border-gray-700">
          <div className="flex flex-col gap-1.5">
            <div className="flex items-baseline justify-between gap-3">
              <label className="text-[13px] font-semibold text-gray-900 dark:text-white">
                {"Nombre de la plantilla"}
              </label>
              <span className="text-[12px] text-gray-400">
                {isEdit ? "Meta no permite cambiarlo" : "solo interno, no se puede cambiar después"}
              </span>
            </div>
            <input
              value={isEdit ? templateName : displayName}
              disabled={isEdit}
              onChange={(e) => setDisplayName(e.target.value)}
              placeholder="Vuelve pronto"
              className={`${inputCls} border-gray-300 disabled:bg-gray-50 disabled:text-gray-500 dark:border-gray-600`}
            />
            {!isEdit && templateName && (
              <span className="text-[12px] text-gray-400">
                {"Meta la registrará como "}
                <code className="rounded bg-gray-100 px-1 dark:bg-gray-700">{templateName}</code>
              </span>
            )}
          </div>

          <div className="flex flex-col gap-1.5">
            <label className="text-[13px] font-semibold text-gray-900 dark:text-white">
              {"Categoría"}
            </label>
            <div className="flex rounded-[9px] bg-gray-100 p-[3px] dark:bg-gray-700">
              {([
                { key: "UTILITY" as const, label: "Utilidad" },
                { key: "MARKETING" as const, label: "Marketing" },
              ]).map((option) => {
                const active = category === option.key;
                return (
                  <button
                    key={option.key}
                    type="button"
                    onClick={() => {
                      setCategory(option.key);
                      const limit =
                        option.key === "MARKETING" ? MAX_BUTTONS - 1 : MAX_BUTTONS;
                      setButtons((current) => current.slice(0, limit));
                    }}
                    style={active ? { color: "var(--color-primary)" } : {}}
                    className={`flex-1 rounded-[7px] px-3 py-2 text-[13px] transition-colors ${
                      active
                        ? "bg-white font-semibold shadow-sm dark:bg-gray-800"
                        : "font-medium text-gray-500 dark:text-gray-400"
                    }`}
                  >
                    {option.label}
                  </button>
                );
              })}
            </div>
            {category === "MARKETING" && (
              <span className="text-[12px] text-amber-600 dark:text-amber-400">
                {"Meta exige botón de baja en las de marketing. Se agrega solo."}
              </span>
            )}
          </div>

          <div className="flex flex-col gap-1.5">
            <div className="flex items-baseline justify-between">
              <label className="text-[13px] font-semibold text-gray-900 dark:text-white">
                {"Encabezado "}
                <span className="font-normal text-gray-400">{"(opcional)"}</span>
              </label>
              {headerType === "TEXT" && (
                <span
                  className={`text-[12px] ${
                    headerText.length > 50 ? "text-red-500" : "text-gray-400"
                  }`}
                >
                  {`${headerText.length}/${MAX_HEADER}`}
                </span>
              )}
            </div>

            <div className="flex rounded-[9px] bg-gray-100 p-[3px] dark:bg-gray-700">
              {([
                { key: "TEXT" as const, label: "Texto" },
                { key: "IMAGE" as const, label: "Imagen" },
              ]).map((option) => {
                const active = headerType === option.key;
                return (
                  <button
                    key={option.key}
                    type="button"
                    onClick={() => setHeaderType(option.key)}
                    style={active ? { color: "var(--color-primary)" } : {}}
                    className={`flex-1 rounded-[7px] px-3 py-1.5 text-[12px] transition-colors ${
                      active
                        ? "bg-white font-semibold shadow-sm dark:bg-gray-800"
                        : "font-medium text-gray-500 dark:text-gray-400"
                    }`}
                  >
                    {option.label}
                  </button>
                );
              })}
            </div>

            {headerType === "TEXT" ? (
              <input
                value={headerText}
                maxLength={MAX_HEADER}
                onChange={(e) => setHeaderText(e.target.value)}
                placeholder="Hola"
                className={`${inputCls} border-gray-300 dark:border-gray-600`}
              />
            ) : (
              <div className="flex flex-col gap-1.5">
                <input
                  type="file"
                  accept="image/jpeg,image/png"
                  disabled={uploading}
                  onChange={(e) => handleImageChange(e.target.files?.[0] ?? null)}
                  className="text-[12px] file:mr-2 file:rounded-lg file:border file:border-gray-300 file:bg-white file:px-3 file:py-1.5 file:text-[12px] file:font-medium file:text-gray-700 dark:file:border-gray-600 dark:file:bg-gray-800 dark:file:text-gray-200"
                />
                {uploading && (
                  <span className="text-[12px] text-gray-400">{"Subiendo imagen..."}</span>
                )}
                {headerMediaURL && !uploading && (
                  <div className="flex items-center gap-2">
                    <span className="min-w-0 flex-1 truncate text-[12px] text-emerald-600 dark:text-emerald-400">
                      {"Imagen cargada"}
                    </span>
                    <button
                      type="button"
                      onClick={() => setHeaderMediaURL("")}
                      className="text-[12px] text-red-500 hover:underline"
                    >
                      {"Quitar"}
                    </button>
                  </div>
                )}
                <span className="text-[12px] text-gray-400">
                  {"JPG o PNG, hasta 5 MB. Meta revisa la imagen junto con el mensaje."}
                </span>
              </div>
            )}
          </div>

          <div className="flex flex-1 flex-col gap-1.5">
            <div className="flex items-center justify-between">
              <label className="text-[13px] font-semibold text-gray-900 dark:text-white">
                {"Mensaje"}
              </label>
              <span
                className={`text-[12px] ${
                  bodyText.length > 950 ? "text-red-500" : "text-gray-400"
                }`}
              >
                {`${bodyText.length}/${MAX_BODY} · ${placeholders} variable(s)`}
              </span>
            </div>
            <textarea
              value={bodyText}
              maxLength={MAX_BODY}
              rows={5}
              onChange={(e) => setBodyText(e.target.value)}
              placeholder="Hola {{1}}, te extrañamos…"
              className={`${inputCls} min-h-[110px] resize-y border-gray-300 dark:border-gray-600`}
            />

            <div className="mt-0.5 flex flex-col gap-2">
              {groupSources(sources).map(([clase, items]) => (
                <div key={clase} className="flex items-start gap-2">
                  <span className="mt-1.5 w-16 shrink-0 text-[11px] font-semibold uppercase tracking-wide text-gray-400">
                    {clase}
                  </span>
                  <div className="flex flex-wrap gap-2">
                    {items.map(([key, label]) => (
                      <button
                        key={key}
                        type="button"
                        onClick={() => insertVariable(key)}
                        className="rounded-full border border-dashed border-gray-300 bg-gray-100 px-2.5 py-1 text-[12px] font-medium text-gray-800 transition-colors hover:border-[var(--color-primary)] hover:text-[var(--color-primary)] dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
                      >
                        {`+ ${label}`}
                      </button>
                    ))}
                  </div>
                </div>
              ))}
            </div>
            <span className="text-[12px] text-gray-400">
              {asFlowResponse
                ? "El cliente se identifica por su teléfono. Si no está registrado se usa lo que pongas en el campo de respaldo."
                : "Tocá una variable para insertarla — no escribas las llaves a mano."}
            </span>

            {variables.length > 0 && (
              <div className="mt-1 flex flex-col gap-1.5 rounded-lg border border-gray-200 p-2.5 dark:border-gray-700">
                {variables.map((variable, index) => (
                  <div key={variable.position} className="flex items-center gap-2">
                    <code className="rounded bg-gray-100 px-1.5 py-0.5 text-[11px] dark:bg-gray-700">
                      {`{{${variable.position}}}`}
                    </code>
                    <span className="min-w-0 flex-1 truncate text-[12px] text-gray-600 dark:text-gray-300">
                      {variableCatalog[variable.source] || variable.source}
                    </span>
                    <input
                      value={variable.fallback}
                      onChange={(e) => updateVariable(index, { fallback: e.target.value })}
                      placeholder="Si falta el dato"
                      className="w-36 rounded-md border border-gray-300 px-2 py-1 text-[12px] outline-none focus:border-[var(--color-primary)] dark:border-gray-600 dark:bg-gray-800"
                    />
                    <button
                      type="button"
                      onClick={() => removeVariable(index)}
                      className="text-[12px] text-red-500 hover:underline"
                    >
                      {"Quitar"}
                    </button>
                  </div>
                ))}
              </div>
            )}

            <div className="mt-1 flex flex-col gap-1.5">
              <div className="flex items-baseline justify-between">
                <label className="text-[13px] font-semibold text-gray-900 dark:text-white">
                  {"Pie "}
                  <span className="font-normal text-gray-400">{"(opcional)"}</span>
                </label>
                <span
                  className={`text-[12px] ${
                    footerText.length > 50 ? "text-red-500" : "text-gray-400"
                  }`}
                >
                  {`${footerText.length}/${MAX_FOOTER}`}
                </span>
              </div>
              <input
                value={footerText}
                maxLength={MAX_FOOTER}
                onChange={(e) => setFooterText(e.target.value)}
                placeholder="Respondé BAJA para no recibir más"
                className={`${inputCls} border-gray-300 dark:border-gray-600`}
              />
            </div>

            <div className="mt-1 flex flex-col gap-1.5">
              <div className="flex items-baseline justify-between">
                <label className="text-[13px] font-semibold text-gray-900 dark:text-white">
                  {"Botones de respuesta "}
                  <span className="font-normal text-gray-400">{"(opcional)"}</span>
                </label>
                <span className="text-[12px] text-gray-400">
                  {`${buttons.length}/${maxButtons}`}
                </span>
              </div>

              {buttons.map((button, index) => (
                <div key={index} className="flex items-center gap-2">
                  <input
                    value={button.text}
                    maxLength={MAX_BUTTON_TEXT}
                    onChange={(e) =>
                      setButtons(
                        buttons.map((item, i) =>
                          i === index ? { ...item, text: e.target.value } : item,
                        ),
                      )
                    }
                    placeholder={"Saber más"}
                    className={`${inputCls} border-gray-300 py-1.5 text-[13px] dark:border-gray-600`}
                  />
                  <span className="w-10 shrink-0 text-right text-[11px] text-gray-400">
                    {`${button.text.length}/${MAX_BUTTON_TEXT}`}
                  </span>
                  <button
                    type="button"
                    onClick={() => setButtons(buttons.filter((_, i) => i !== index))}
                    className="text-[12px] text-red-500 hover:underline"
                  >
                    {"Quitar"}
                  </button>
                </div>
              ))}

              {buttons.length < maxButtons && (
                <button
                  type="button"
                  onClick={() => setButtons([...buttons, { text: "" }])}
                  className="self-start rounded-full border border-dashed border-gray-300 px-2.5 py-1 text-[12px] font-medium text-gray-800 transition-colors hover:border-[var(--color-primary)] hover:text-[var(--color-primary)] dark:border-gray-600 dark:text-gray-100"
                >
                  {"+ Agregar botón"}
                </button>
              )}

              <span className="text-[12px] text-gray-400">
                {category === "MARKETING"
                  ? "En marketing solo caben 2: el botón de baja ocupa el tercero."
                  : "Hasta 3 botones. El cliente responde tocándolos."}
              </span>
              <span className="text-[12px] text-gray-400">
                {"Qué responde cada botón se arma en el diagrama de flujo."}
              </span>

              {brokenLinks.length > 0 && (
                <div className="rounded-md border border-amber-300 bg-amber-50 px-3 py-2 text-[12px] text-amber-800 dark:border-amber-700 dark:bg-amber-950/30 dark:text-amber-300">
                  {`Al guardar se pierde la conexión de flujo de: ${brokenLinks.join(", ")}. Vas a tener que volver a enlazar la respuesta en el flujo.`}
                </div>
              )}
            </div>

            {mismatch && (
              <span className="text-[12px] text-red-600">
                {`El cuerpo usa ${placeholders} variable(s) y hay ${variables.length} declarada(s). Meta la rechaza si no coinciden.`}
              </span>
            )}
            {!mismatch && ruleError && (
              <span data-testid="template-rule-error" className="text-[12px] text-red-600">
                {ruleError}
              </span>
            )}
          </div>
        </div>

        <div className="flex min-h-0 flex-1 flex-col gap-3 bg-gray-50 p-6 dark:bg-gray-900/40">
          <span className="text-[11px] font-semibold uppercase tracking-[0.08em] text-gray-400">
            {"Así lo ve el cliente"}
          </span>
          <div className="flex flex-1 flex-col gap-2 rounded-xl bg-[#e9e2d9] p-4 dark:bg-[#2a2724]">
            <span className="self-center rounded-full bg-gray-50 px-2.5 py-0.5 text-[11px] text-gray-500 dark:bg-gray-800">
              {"Hoy"}
            </span>
            <TemplateBubble
              headerType={headerType}
              headerMediaURL={headerMediaURL}
              headerText={headerText}
              bodyText={bodyPreview}
              footerText={footerText}
              buttons={previewButtons}
            />
          </div>
          <span className="text-[12px] leading-snug text-gray-400">
            {"Los datos son de ejemplo. En el envío real se reemplazan por los del cliente."}
          </span>
        </div>
      </div>

      <div className="flex shrink-0 items-center justify-between gap-3 border-t border-gray-200 bg-gray-50 px-7 py-4 dark:border-gray-700 dark:bg-gray-900/40">
        {error ? (
          <span className="flex min-w-0 items-start gap-2 rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-[13px] font-medium text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-300">
            <svg
              className="mt-0.5 h-4 w-4 shrink-0"
              fill="none"
              stroke="currentColor"
              strokeWidth={2}
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M12 9v4m0 4h.01M10.29 3.86 1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"
              />
            </svg>
            <span>{error}</span>
          </span>
        ) : (
          <span className="text-[13px] text-gray-400">
            {"Queda como borrador. La enviás a revisión cuando el flujo esté listo."}
          </span>
        )}
        <div className="flex gap-2.5">
          <button
            type="button"
            onClick={onCancel}
            disabled={loading}
            className="rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-900 transition-colors hover:bg-gray-100 disabled:opacity-40 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
          >
            {"Cancelar"}
          </button>
          <button
            type="submit"
            disabled={loading || mismatch}
            style={{
              backgroundColor: "var(--color-primary)",
              color: "var(--color-on-primary, white)",
            }}
            className="rounded-lg px-4 py-2 text-sm font-semibold transition-opacity hover:opacity-90 disabled:opacity-40"
          >
            {loading ? "Guardando..." : isEdit ? "Guardar cambios" : "Guardar borrador"}
          </button>
        </div>
      </div>
    </form>
  );
}
