"use client";

import { useEffect, useMemo, useState } from "react";
import { DayPicker } from "react-day-picker";
import { es } from "date-fns/locale";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import { Button } from "@/shared/ui/button";
import { useToast } from "@/shared/providers/toast-provider";
import { Flow, WhatsappTemplate } from "../../domain/scheduled-types";
import { CustomerInfo } from "@/services/modules/customers/domain/types";
import { getCustomersAction } from "@/services/modules/customers/infra/actions";
import {
  Campaign,
  CampaignAudienceLocation,
  CampaignAudiencePreview,
  CampaignAudienceType,
  CampaignDeliveryMode,
  CampaignScheduleMode,
  CreateCampaignDTO,
} from "../../domain/campaign-types";
import {
  createCampaignAction,
  listCampaignAudienceLocationsAction,
  previewCampaignAudienceAction,
  updateCampaignAction,
} from "../../infra/actions/campaigns";
import { WhatsAppBubblePreview } from "./WhatsAppBubblePreview";
import { CampaignRules } from "./CampaignRules";

interface CampaignFormProps {
  businessId?: number;
  templates: WhatsappTemplate[];
  flows?: Flow[];
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

const pad = (value: number) => String(value).padStart(2, "0");

const localDate = (date: Date) =>
  `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}`;

const localTime = (date: Date) => `${pad(date.getHours())}:${pad(date.getMinutes())}`;

const parseLocalDate = (value: string) => {
  const [year, month, day] = value.split("-").map(Number);
  return new Date(year, (month || 1) - 1, day || 1);
};

const shortDate = (value: string) =>
  parseLocalDate(value).toLocaleDateString("es-CO", { day: "numeric", month: "short" });

type FilterKey =
  | "state"
  | "city"
  | "created_from_days"
  | "registered_before_days"
  | "only_without_order"
  | "min_orders"
  | "min_spent"
  | "last_purchase_before_days";

interface FilterDef {
  key: FilterKey;
  label: string;
  kind: "text" | "number" | "bool" | "select";
  suffix?: string;
  placeholder?: string;
}

const FILTER_DEFS: FilterDef[] = [
  { key: "state", label: "Departamento", kind: "select" },
  { key: "city", label: "Ciudad", kind: "select" },
  { key: "created_from_days", label: "Registrado en los \u00faltimos", kind: "number", suffix: "d\u00edas" },
  { key: "registered_before_days", label: "Registrado hace m\u00e1s de", kind: "number", suffix: "d\u00edas" },
  { key: "only_without_order", label: "Nunca me ha comprado", kind: "bool" },
  { key: "min_orders", label: "Con al menos", kind: "number", suffix: "compras" },
  { key: "min_spent", label: "Ha gastado al menos", kind: "number", suffix: "$" },
  {
    key: "last_purchase_before_days",
    label: "Sin comprar hace m\u00e1s de",
    kind: "number",
    suffix: "d\u00edas",
  },
];

export function CampaignForm({
  businessId,
  templates,
  flows = [],
  campaign,
  onSuccess,
  onCancel,
}: CampaignFormProps) {
  const { showToast } = useToast();

  const [name, setName] = useState(campaign?.Name || "");
  const [senderName, setSenderName] = useState(campaign?.SenderName || "");
  const [flowId, setFlowId] = useState<number>(campaign?.FlowID || 0);
  const [audienceMode, setAudienceMode] = useState<"all" | "filters" | "manual">(() => {
    if (!campaign) return "all";
    if (campaign.AudienceType === "all_clients") return "all";
    return (campaign.AudienceParams?.ClientIDs?.length || 0) > 0 ? "manual" : "filters";
  });

  const audienceType: CampaignAudienceType =
    audienceMode === "all" ? "all_clients" : "filtered_clients";
  const [customers, setCustomers] = useState<CustomerInfo[]>([]);
  const [customersLoading, setCustomersLoading] = useState(false);
  const [customerSearch, setCustomerSearch] = useState("");
  const [capEnabled, setCapEnabled] = useState(
    (campaign?.AudienceParams?.ExcludeRecentDays || 0) > 0,
  );
  const [capDays, setCapDays] = useState(campaign?.AudienceParams?.ExcludeRecentDays || 15);
  const [capMax, setCapMax] = useState(campaign?.AudienceParams?.ExcludeRecentMax || 1);

  const [filters, setFilters] = useState<Partial<Record<FilterKey, string>>>(() => {
    const params = campaign?.AudienceParams;
    const initial: Partial<Record<FilterKey, string>> = {};
    if (!params) return initial;
    if (params.State) initial.state = params.State;
    if (params.City) initial.city = params.City;
    if (params.CreatedFromDays) initial.created_from_days = String(params.CreatedFromDays);
    if (params.RegisteredBeforeDays)
      initial.registered_before_days = String(params.RegisteredBeforeDays);
    if (params.OnlyWithoutOrder) initial.only_without_order = "1";
    if (params.MinOrders) initial.min_orders = String(params.MinOrders);
    if (params.MinSpent) initial.min_spent = String(params.MinSpent);
    if (params.LastPurchaseBeforeDays)
      initial.last_purchase_before_days = String(params.LastPurchaseBeforeDays);
    return initial;
  });

  const activeFilters = FILTER_DEFS.filter((def) => filters[def.key] !== undefined);
  const availableFilters = FILTER_DEFS.filter((def) => filters[def.key] === undefined);

  const addFilter = (key: FilterKey) => {
    const def = FILTER_DEFS.find((item) => item.key === key);
    setFilters((current) => ({ ...current, [key]: def?.kind === "bool" ? "1" : "" }));
  };

  const removeFilter = (key: FilterKey) =>
    setFilters((current) => {
      const next = { ...current };
      delete next[key];
      return next;
    });

  const filterNumber = (key: FilterKey) => Number(filters[key] || 0) || 0;
  const [selectedIds, setSelectedIds] = useState<number[]>(
    campaign?.AudienceParams?.ClientIDs || [],
  );
  const [windowStart, setWindowStart] = useState(campaign?.SendWindowStart || "09:00");
  const [windowEnd, setWindowEnd] = useState(campaign?.SendWindowEnd || "19:00");
  const [dailyCap, setDailyCap] = useState(campaign?.DailySendCap || 250);
  const [batchSize, setBatchSize] = useState(campaign?.BatchSize || 50);

  const scheduledAt = campaign?.ScheduledAt ? new Date(campaign.ScheduledAt) : null;
  const [startMode, setStartMode] = useState<"now" | "scheduled">(
    scheduledAt ? "scheduled" : "now",
  );
  const [startDate, setStartDate] = useState(scheduledAt ? localDate(scheduledAt) : "");
  const [startTime, setStartTime] = useState(scheduledAt ? localTime(scheduledAt) : "09:00");

  const [scheduleMode, setScheduleMode] = useState<CampaignScheduleMode>(
    campaign?.ScheduleMode || "daily",
  );
  const [intervalDays, setIntervalDays] = useState(campaign?.IntervalDays || 7);
  const [sendDates, setSendDates] = useState<string[]>(campaign?.SendDates || []);
  const [deliveryMode, setDeliveryMode] = useState<CampaignDeliveryMode>(
    campaign?.DeliveryMode || "distribute",
  );
  const [limitOccurrences, setLimitOccurrences] = useState((campaign?.Occurrences || 0) > 0);
  const [occurrences, setOccurrences] = useState(campaign?.Occurrences || 4);

  const needsOccurrences = scheduleMode !== "dates" && deliveryMode === "repeat";
  const usesOccurrences = scheduleMode !== "dates" && (limitOccurrences || needsOccurrences);
  const sortedDates = useMemo(() => [...sendDates].sort(), [sendDates]);

  const [locations, setLocations] = useState<CampaignAudienceLocation[]>([]);

  const stateOptions = useMemo(() => {
    const totals = new Map<string, number>();
    locations.forEach((item) => {
      if (!item.state) return;
      totals.set(item.state, (totals.get(item.state) || 0) + item.clients);
    });
    return Array.from(totals.entries())
      .map(([value, clients]) => ({ value, clients }))
      .sort((a, b) => b.clients - a.clients);
  }, [locations]);

  const cityOptions = useMemo(() => {
    const chosenState = filters.state || "";
    return locations
      .filter((item) => !chosenState || item.state === chosenState)
      .map((item) => ({ value: item.city, clients: item.clients }));
  }, [locations, filters.state]);

  const [preview, setPreview] = useState<CampaignAudiencePreview | null>(null);
  const [previewLoading, setPreviewLoading] = useState(false);
  const [showAudience, setShowAudience] = useState(false);
  const [saving, setSaving] = useState(false);

  const selectedFlow = useMemo(
    () => flows.find((flow) => flow.ID === flowId) || null,
    [flows, flowId],
  );

  const effectiveTemplateId = selectedFlow?.RootTemplateID ?? 0;

  const selected = useMemo(
    () => templates.find((template) => template.ID === effectiveTemplateId) || null,
    [templates, effectiveTemplateId],
  );

  const buildDTO = (): CreateCampaignDTO => ({
    whatsapp_template_id: effectiveTemplateId,
    flow_id: flowId || null,
    name: name.trim(),
    sender_name: senderName.trim(),
    audience_type: audienceType,
    client_ids: audienceMode === "manual" ? selectedIds : [],
    city: audienceMode === "filters" ? (filters.city || "").trim() : "",
    state: audienceMode === "filters" ? (filters.state || "").trim() : "",
    created_from_days: audienceMode === "filters" ? filterNumber("created_from_days") : 0,
    registered_before_days:
      audienceMode === "filters" ? filterNumber("registered_before_days") : 0,
    only_without_order: audienceMode === "filters" && filters.only_without_order === "1",
    min_orders: audienceMode === "filters" ? filterNumber("min_orders") : 0,
    min_spent: audienceMode === "filters" ? filterNumber("min_spent") : 0,
    last_purchase_before_days:
      audienceMode === "filters" ? filterNumber("last_purchase_before_days") : 0,
    exclude_recent_days: capEnabled ? capDays : 0,
    exclude_recent_max: capEnabled ? capMax : 0,
    send_window_start: windowStart,
    send_window_end: windowEnd,
    daily_send_cap: dailyCap,
    batch_size: batchSize,
    schedule_mode: scheduleMode,
    interval_days: scheduleMode === "interval" ? intervalDays : 0,
    send_dates: scheduleMode === "dates" ? sortedDates : [],
    delivery_mode: deliveryMode,
    occurrences: usesOccurrences ? occurrences : 0,
    scheduled_at:
      scheduleMode !== "dates" && startMode === "scheduled" && startDate
        ? new Date(`${startDate}T${startTime || "00:00"}`).toISOString()
        : null,
    variable_values: senderName.trim() ? { "sender.name": senderName.trim() } : undefined,
  });

  useEffect(() => {
    const loadLocations = async () => {
      const result = await listCampaignAudienceLocationsAction(businessId);
      if (result.success && result.data) setLocations(result.data);
    };
    loadLocations();
  }, [businessId]);

  useEffect(() => {
    const loadPreview = async () => {
      setPreviewLoading(true);
      const result = await previewCampaignAudienceAction(
        { ...buildDTO(), whatsapp_template_id: effectiveTemplateId || 0 },
        businessId,
      );
      setPreviewLoading(false);
      if (result.success && result.data) {
        setPreview(result.data);
      }
    };

    const timer = setTimeout(loadPreview, 400);
    return () => clearTimeout(timer);
  }, [audienceMode, selectedIds, filters, capEnabled, capDays, capMax, businessId]);

  useEffect(() => {
    const loadCustomers = async () => {
      setCustomersLoading(true);
      try {
        const result = await getCustomersAction({
          business_id: businessId,
          page: 1,
          page_size: 200,
          search: customerSearch.trim() || undefined,
        });
        setCustomers(result.data || []);
      } catch {
        setCustomers([]);
      }
      setCustomersLoading(false);
    };

    const timer = setTimeout(loadCustomers, 350);
    return () => clearTimeout(timer);
  }, [businessId, customerSearch]);

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
      showToast("Ponle un nombre a la campa\u00f1a", "error");
      return;
    }
    if (!flowId) {
      showToast("Elige el flujo que se va a enviar", "error");
      return;
    }
    if (scheduleMode === "dates" && sendDates.length === 0) {
      showToast("Marca en el calendario los d\u00edas en que sale la campa\u00f1a", "error");
      return;
    }
    if (scheduleMode === "interval" && (!intervalDays || intervalDays < 1)) {
      showToast("Indica cada cu\u00e1ntos d\u00edas se env\u00eda", "error");
      return;
    }
    if (usesOccurrences && (!occurrences || occurrences < 1)) {
      showToast("Indica cu\u00e1ntas veces se env\u00eda", "error");
      return;
    }
    if (scheduleMode !== "dates" && startMode === "scheduled" && !startDate) {
      showToast("Eleg\u00ed la fecha en la que arranca la campa\u00f1a", "error");
      return;
    }
    if (!effectiveTemplateId) {
      showToast("Ese flujo no tiene plantilla inicial. Editalo y eleg\u00ed una.", "error");
      return;
    }

    setSaving(true);
    const result = campaign
      ? await updateCampaignAction(campaign.ID, buildDTO(), businessId)
      : await createCampaignAction(buildDTO(), businessId);
    setSaving(false);

    if (!result.success) {
      showToast(result.error || "No se pudo guardar la campa\u00f1a", "error");
      return;
    }

    showToast(campaign ? "Campa\u00f1a actualizada" : "Campa\u00f1a creada como borrador", "success");
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
                <Label htmlFor="campaign-template">{"Flujo"}</Label>

                    <select
                      id="campaign-template"
                      value={flowId}
                      onChange={(e) => setFlowId(Number(e.target.value))}
                      className="mt-1.5 w-full rounded-lg border border-gray-300 bg-white px-4 py-3 text-sm text-gray-900 dark:border-gray-600 dark:bg-gray-700 dark:text-white"
                    >
                      <option value={0}>{"Seleccionar..."}</option>
                      {flows.map((flow) => (
                        <option key={flow.ID} value={flow.ID}>
                          {flow.Name}
                        </option>
                      ))}
                    </select>
                    {flows.length === 0 ? (
                      <p className="mt-1 text-xs text-amber-600">
                        {"Todav\u00eda no tienes flujos. Crea uno en la pesta\u00f1a Flujos."}
                      </p>
                    ) : selectedFlow ? (
                      <p className="mt-1 text-xs text-gray-500">
                        {`Arranca con ${selectedFlow.RootTemplateName || "su plantilla inicial"} y sigue por los botones. ${selectedFlow.StepCount} paso(s).`}
                      </p>
                    ) : (
                      <p className="mt-1 text-xs text-gray-500">
                        {"El primer mensaje abre la conversaci\u00f3n y los botones encadenan el resto."}
                      </p>
                    )}
                    {selectedFlow && selectedFlow.PendingCount > 0 && (
                      <p className="mt-1 text-xs font-medium text-amber-700 dark:text-amber-500">
                        {`${selectedFlow.PendingCount} respuesta(s) sin aprobar: no vas a poder lanzar hasta que Meta las apruebe.`}
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
                {([
                  { key: "all" as const, label: "Todos" },
                  { key: "filters" as const, label: "Filtrar" },
                  { key: "manual" as const, label: "Elegir a mano" },
                ]).map((option) => (
                  <button
                    key={option.key}
                    type="button"
                    onClick={() => setAudienceMode(option.key)}
                    className={`rounded-md px-3 py-1.5 text-xs font-medium ${
                      audienceMode === option.key
                        ? "bg-[var(--color-primary)] text-white"
                        : "bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-300"
                    }`}
                  >
                    {option.label}
                  </button>
                ))}
              </span>,
            )}

            {audienceMode === "manual" ? (
              <div className="space-y-2">
                <div className="flex flex-wrap items-center gap-3">
                  <Input
                    className="max-w-xs"
                    value={customerSearch}
                    onChange={(e) => setCustomerSearch(e.target.value)}
                    placeholder={"Buscar cliente por nombre o celular"}
                  />
                  <span className="text-xs text-gray-500">
                    {`${selectedIds.length} elegido(s)`}
                  </span>
                  {selectedIds.length > 0 && (
                    <button
                      type="button"
                      onClick={() => setSelectedIds([])}
                      className="text-xs font-medium text-red-500 hover:underline"
                    >
                      {"Quitar todos"}
                    </button>
                  )}
                </div>

                <div className="max-h-64 overflow-y-auto rounded-lg border border-gray-200 dark:border-gray-700">
                  {customersLoading ? (
                    <p className="p-4 text-center text-xs text-gray-400">{"Cargando..."}</p>
                  ) : customers.length === 0 ? (
                    <p className="p-4 text-center text-xs text-gray-400">
                      {"No hay clientes cargados que coincidan."}
                    </p>
                  ) : (
                    <ul className="divide-y divide-gray-100 dark:divide-gray-700">
                      {customers.map((customer) => {
                        const checked = selectedIds.includes(customer.id);
                        return (
                          <li key={customer.id}>
                            <label className="flex cursor-pointer items-center gap-3 px-3 py-2 hover:bg-gray-50 dark:hover:bg-gray-700/40">
                              <input
                                type="checkbox"
                                checked={checked}
                                onChange={() =>
                                  setSelectedIds((current) =>
                                    checked
                                      ? current.filter((id) => id !== customer.id)
                                      : [...current, customer.id],
                                  )
                                }
                              />
                              <span className="min-w-0 flex-1 truncate text-xs text-gray-800 dark:text-gray-100">
                                {customer.name}
                              </span>
                              <span className="shrink-0 text-[11px] text-gray-400">
                                {customer.phone}
                              </span>
                            </label>
                          </li>
                        );
                      })}
                    </ul>
                  )}
                </div>
              </div>
            ) : audienceMode === "filters" ? (
              <div className="space-y-2">
                {activeFilters.length === 0 ? (
                  <p className="text-xs text-gray-500">
                    {"Sin filtros le llega a todos. Agreg\u00e1 uno para acotar la audiencia."}
                  </p>
                ) : (
                  <div className="space-y-2">
                    {activeFilters.map((def) => (
                      <div
                        key={def.key}
                        className="flex flex-wrap items-center gap-2 rounded-lg border border-gray-200 px-3 py-2 dark:border-gray-700"
                      >
                        <span className="text-xs text-gray-700 dark:text-gray-200">
                          {def.label}
                        </span>

                        {def.kind === "select" ? (
                          <select
                            value={filters[def.key] ?? ""}
                            onChange={(e) =>
                              setFilters((current) => ({
                                ...current,
                                [def.key]: e.target.value,
                              }))
                            }
                            className="min-w-48 rounded-md border border-gray-300 px-2 py-1 text-xs dark:border-gray-600 dark:bg-gray-800"
                          >
                            <option value="">
                              {def.key === "state" ? "Todos" : "Todas"}
                            </option>
                            {(def.key === "state" ? stateOptions : cityOptions).map(
                              (option) => (
                                <option key={option.value} value={option.value}>
                                  {`${option.value} (${option.clients})`}
                                </option>
                              ),
                            )}
                          </select>
                        ) : def.kind !== "bool" ? (
                          <input
                            type={def.kind === "number" ? "number" : "text"}
                            min={def.kind === "number" ? 0 : undefined}
                            value={filters[def.key] ?? ""}
                            placeholder={def.placeholder}
                            onChange={(e) =>
                              setFilters((current) => ({
                                ...current,
                                [def.key]: e.target.value,
                              }))
                            }
                            className="w-28 rounded-md border border-gray-300 px-2 py-1 text-xs dark:border-gray-600 dark:bg-gray-800"
                          />
                        ) : null}

                        {def.suffix && (
                          <span className="text-xs text-gray-400">{def.suffix}</span>
                        )}

                        <button
                          type="button"
                          onClick={() => removeFilter(def.key)}
                          className="ml-auto text-xs font-medium text-red-500 hover:underline"
                        >
                          {"Quitar"}
                        </button>
                      </div>
                    ))}
                  </div>
                )}

                {availableFilters.length > 0 && (
                  <select
                    value=""
                    onChange={(e) => {
                      if (e.target.value) addFilter(e.target.value as FilterKey);
                    }}
                    style={{ height: 32, fontSize: 12, padding: "0 8px", borderRadius: 8 }}
                    className="border border-dashed border-gray-300 bg-white text-gray-700 outline-none dark:border-gray-600 dark:bg-gray-800 dark:text-gray-200"
                  >
                    <option value="">{"+ Agregar filtro"}</option>
                    {availableFilters.map((def) => (
                      <option key={def.key} value={def.key}>
                        {def.label}
                      </option>
                    ))}
                  </select>
                )}
              </div>
            ) : (
              <p className="text-xs text-gray-500">
                {"Le llega a todos tus clientes con celular v\u00e1lido que no se hayan dado de baja."}
              </p>
            )}

            <div className="mt-4 rounded-lg border border-gray-200 p-3 dark:border-gray-700">
              <label className="flex items-center gap-2 text-xs font-medium text-gray-700 dark:text-gray-200">
                <input
                  type="checkbox"
                  checked={capEnabled}
                  onChange={(e) => setCapEnabled(e.target.checked)}
                />
                {"No repetirle a quien ya recibi\u00f3 campa\u00f1as hace poco"}
              </label>

              {capEnabled && (
                <div className="mt-2 flex flex-wrap items-center gap-2 text-xs text-gray-600 dark:text-gray-300">
                  {"Excluir a quien ya recibi\u00f3"}
                  <input
                    type="number"
                    min={1}
                    max={20}
                    value={capMax}
                    onChange={(e) => setCapMax(Number(e.target.value))}
                    className="w-16 rounded-md border border-gray-300 px-2 py-1 text-xs dark:border-gray-600 dark:bg-gray-800"
                  />
                  {"o m\u00e1s campa\u00f1as en los \u00faltimos"}
                  <input
                    type="number"
                    min={1}
                    max={365}
                    value={capDays}
                    onChange={(e) => setCapDays(Number(e.target.value))}
                    className="w-16 rounded-md border border-gray-300 px-2 py-1 text-xs dark:border-gray-600 dark:bg-gray-800"
                  />
                  {"d\u00edas"}
                </div>
              )}

              <p className="mt-2 text-[11px] text-gray-400">
                {
                  "Cuenta los mensajes de campa\u00f1as que s\u00ed salieron. Evita quemar la l\u00ednea escribi\u00e9ndole tres veces a la misma persona."
                }
              </p>
            </div>
          </div>

          <div className={band}>
            {step(
              "3",
              "Cu\u00e1ndo sale",
              "D\u00edas de env\u00edo, franja y ritmo de las tandas",
              scheduleMode === "dates" ? undefined : <span className="flex gap-1.5">
                <button
                  type="button"
                  onClick={() => setStartMode("now")}
                  className={`rounded-md px-3 py-1.5 text-xs font-medium ${
                    startMode === "now"
                      ? "bg-[var(--color-primary)] text-white"
                      : "bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-300"
                  }`}
                >
                  {"Al lanzarla"}
                </button>
                <button
                  type="button"
                  onClick={() => setStartMode("scheduled")}
                  className={`rounded-md px-3 py-1.5 text-xs font-medium ${
                    startMode === "scheduled"
                      ? "bg-[var(--color-primary)] text-white"
                      : "bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-300"
                  }`}
                >
                  {"Programar"}
                </button>
              </span>,
            )}

            <div className="mb-4 flex flex-wrap gap-1.5">
              {([
                { key: "daily" as const, label: "Todos los d\u00edas" },
                { key: "interval" as const, label: "Cada cierto tiempo" },
                { key: "dates" as const, label: "Fechas del calendario" },
              ]).map((option) => (
                <button
                  key={option.key}
                  type="button"
                  onClick={() => setScheduleMode(option.key)}
                  className={`rounded-md px-3 py-1.5 text-xs font-medium ${
                    scheduleMode === option.key
                      ? "bg-[var(--color-primary)] text-white"
                      : "bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-300"
                  }`}
                >
                  {option.label}
                </button>
              ))}
            </div>

            {scheduleMode === "interval" && (
              <div className="mb-4 flex flex-wrap items-center gap-2 text-xs text-gray-600 dark:text-gray-300">
                {"Sale cada"}
                <input
                  type="number"
                  min={1}
                  max={365}
                  value={intervalDays}
                  onChange={(e) => setIntervalDays(Number(e.target.value))}
                  className="w-16 rounded-md border border-gray-300 px-2 py-1 text-xs dark:border-gray-600 dark:bg-gray-800"
                />
                {"d\u00edas, contando desde el d\u00eda de arranque. Con 7 es una vez por semana."}
              </div>
            )}

            {scheduleMode === "dates" && (
              <div className="mb-4 grid gap-4 md:grid-cols-[auto_minmax(0,1fr)]">
                <div className="rounded-lg border border-gray-200 p-2 dark:border-gray-700">
                  <DayPicker
                    mode="multiple"
                    locale={es}
                    selected={sendDates.map(parseLocalDate)}
                    onSelect={(dates) => setSendDates((dates || []).map(localDate))}
                    disabled={{ before: new Date() }}
                    classNames={{
                      months: "relative",
                      month_caption: "mb-2 flex h-8 items-center justify-center text-sm font-semibold capitalize text-gray-900 dark:text-gray-100",
                      nav: "absolute inset-x-0 top-0 flex h-8 items-center justify-between",
                      button_previous: "rounded p-1 text-gray-600 hover:bg-gray-100 disabled:opacity-30 dark:text-gray-300 dark:hover:bg-gray-700",
                      button_next: "rounded p-1 text-gray-600 hover:bg-gray-100 disabled:opacity-30 dark:text-gray-300 dark:hover:bg-gray-700",
                      chevron: "h-4 w-4 fill-current",
                      weekday: "w-9 pb-1 text-[10px] font-semibold uppercase text-gray-400",
                      day: "p-0.5 text-center",
                      day_button: "h-8 w-8 rounded-md text-xs text-gray-800 hover:bg-gray-100 dark:text-gray-100 dark:hover:bg-gray-700",
                      selected: "[&>button]:bg-[var(--color-primary)] [&>button]:font-semibold [&>button]:text-white",
                      today: "[&>button]:ring-1 [&>button]:ring-[var(--color-primary)]",
                      disabled: "[&>button]:cursor-not-allowed [&>button]:opacity-30",
                      outside: "opacity-40",
                    }}
                  />
                </div>
                <div>
                  <p className="text-xs text-gray-600 dark:text-gray-300">
                    {sortedDates.length === 0
                      ? "Toca los d\u00edas en que debe salir. Solo se env\u00eda en esas fechas."
                      : `${sortedDates.length} fecha(s) elegida(s):`}
                  </p>
                  {sortedDates.length > 0 && (
                    <div className="mt-2 flex flex-wrap gap-1.5">
                      {sortedDates.map((date) => (
                        <span
                          key={date}
                          className="flex items-center gap-1 rounded-full bg-gray-100 px-2 py-0.5 text-[11px] text-gray-700 dark:bg-gray-700 dark:text-gray-200"
                        >
                          {shortDate(date)}
                          <button
                            type="button"
                            onClick={() => setSendDates((current) => current.filter((item) => item !== date))}
                            className="text-gray-400 hover:text-red-500"
                            aria-label={`Quitar ${date}`}
                          >
                            {"\u00d7"}
                          </button>
                        </span>
                      ))}
                    </div>
                  )}
                  <p className="mt-3 text-[11px] text-gray-400">
                    {"Cada fecha arranca dentro de la franja horaria que elijas abajo."}
                  </p>
                </div>
              </div>
            )}

            <div className="mb-4 grid gap-2 md:grid-cols-2">
              {([
                {
                  key: "distribute" as const,
                  title: "Repartir la audiencia",
                  note: "Cada d\u00eda de env\u00edo sigue con los siguientes de la lista, hasta el m\u00e1ximo diario. A cada cliente le llega una sola vez.",
                },
                {
                  key: "repeat" as const,
                  title: "Repetir a toda la audiencia",
                  note: "Cada d\u00eda de env\u00edo vuelve a escribirle a todos, hasta el m\u00e1ximo diario. Si no alcanzan ese d\u00eda, se omiten en esa vuelta.",
                },
              ]).map((option) => (
                <button
                  key={option.key}
                  type="button"
                  onClick={() => setDeliveryMode(option.key)}
                  className={`rounded-lg border p-3 text-left ${
                    deliveryMode === option.key
                      ? "border-[var(--color-primary)] bg-[var(--color-primary)]/5"
                      : "border-gray-200 dark:border-gray-700"
                  }`}
                >
                  <span className="block text-xs font-semibold text-gray-900 dark:text-gray-100">
                    {option.title}
                  </span>
                  <span className="mt-1 block text-[11px] text-gray-500">{option.note}</span>
                </button>
              ))}
            </div>

            {deliveryMode === "repeat" && capEnabled && (
              <p className="mb-4 rounded-lg bg-amber-50 p-2 text-[11px] text-amber-700 dark:bg-amber-900/20 dark:text-amber-400">
                {"Ojo: el tope de frecuencia del paso 2 tambi\u00e9n cuenta los env\u00edos de esta campa\u00f1a y puede dejar por fuera a quien recibi\u00f3 la vuelta anterior."}
              </p>
            )}

            {scheduleMode !== "dates" && (
              <div className="mb-4 flex flex-wrap items-center gap-2 text-xs text-gray-600 dark:text-gray-300">
                {needsOccurrences ? (
                  <span>{"Se env\u00eda"}</span>
                ) : (
                  <label className="flex items-center gap-2">
                    <input
                      type="checkbox"
                      checked={limitOccurrences}
                      onChange={(e) => setLimitOccurrences(e.target.checked)}
                    />
                    {"Terminar despu\u00e9s de"}
                  </label>
                )}
                {usesOccurrences ? (
                  <>
                    <input
                      type="number"
                      min={1}
                      max={365}
                      value={occurrences}
                      onChange={(e) => setOccurrences(Number(e.target.value))}
                      className="w-16 rounded-md border border-gray-300 px-2 py-1 text-xs dark:border-gray-600 dark:bg-gray-800"
                    />
                    {scheduleMode === "daily" ? "d\u00edas de env\u00edo" : "env\u00edos"}
                  </>
                ) : (
                  <span className="text-gray-400">
                    {"(si no, sigue hasta que le llegue a toda la audiencia)"}
                  </span>
                )}
              </div>
            )}

            {scheduleMode === "dates" ? null : startMode === "scheduled" ? (
              <div className="mb-4 grid gap-4 md:grid-cols-4">
                <div>
                  <Label htmlFor="campaign-start-date">{"Fecha de inicio"}</Label>
                  <Input
                    id="campaign-start-date"
                    className="mt-1.5"
                    type="date"
                    min={localDate(new Date())}
                    value={startDate}
                    onChange={(e) => setStartDate(e.target.value)}
                  />
                </div>
                <div>
                  <Label htmlFor="campaign-start-time">{"Hora de inicio"}</Label>
                  <Input
                    id="campaign-start-time"
                    className="mt-1.5"
                    type="time"
                    value={startTime}
                    onChange={(e) => setStartTime(e.target.value)}
                  />
                </div>
                <p className="self-end pb-3 text-xs text-gray-500 md:col-span-2">
                  {"La primera tanda sale a esa hora, siempre que caiga dentro de la franja."}
                </p>
              </div>
            ) : (
              <p className="mb-4 text-xs text-gray-500">
                {"Empieza apenas la lances, dentro de la franja horaria que elijas abajo."}
              </p>
            )}

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

                <dl className="mt-3 space-y-1 border-t border-gray-200 pt-3 text-[11px] dark:border-gray-700">
                  <div className="flex justify-between">
                    <dt className="text-gray-500">{"coinciden"}</dt>
                    <dd className="font-semibold text-gray-900 dark:text-gray-100">
                      {preview.total}
                    </dd>
                  </div>
                  <div className="flex justify-between">
                    <dt className="text-gray-500">{"reciben"}</dt>
                    <dd className="font-semibold text-gray-900 dark:text-gray-100">
                      {preview.reachable}
                    </dd>
                  </div>
                  <div className="flex justify-between">
                    <dt className="text-amber-600 dark:text-amber-400">{"quedan fuera"}</dt>
                    <dd className="font-semibold text-amber-600 dark:text-amber-400">
                      {Math.max(0, preview.total - preview.reachable)}
                    </dd>
                  </div>
                  <div className="flex justify-between pl-3">
                    <dt className="text-gray-400">{"sin celular"}</dt>
                    <dd className="text-gray-500">{preview.no_phone}</dd>
                  </div>
                  <div className="flex justify-between pl-3">
                    <dt className="text-gray-400">{"dados de baja"}</dt>
                    <dd className="text-gray-500">{preview.opted_out}</dd>
                  </div>
                </dl>

                {(preview.clients?.length ?? 0) > 0 && (
                  <>
                    <button
                      type="button"
                      onClick={() => setShowAudience((current) => !current)}
                      className="mt-3 text-xs font-medium text-[var(--color-primary)] hover:underline"
                    >
                      {showAudience ? "Ocultar clientes" : "Ver clientes"}
                    </button>

                    {showAudience && (
                      <ul className="mt-2 max-h-56 divide-y divide-gray-100 overflow-y-auto rounded-lg border border-gray-200 dark:divide-gray-700 dark:border-gray-700">
                        {(preview.clients || []).map((client) => (
                          <li
                            key={client.client_id}
                            className="flex items-center gap-2 px-2.5 py-1.5"
                          >
                            <span className="min-w-0 flex-1 truncate text-[11px] text-gray-800 dark:text-gray-100">
                              {client.name}
                            </span>
                            <span className="shrink-0 text-[10px] text-gray-400">
                              {client.phone}
                            </span>
                          </li>
                        ))}
                      </ul>
                    )}

                    {preview.reachable > (preview.clients?.length ?? 0) && showAudience && (
                      <p className="mt-1 text-[10px] text-gray-400">
                        {`Se muestran los primeros ${preview.clients?.length ?? 0}.`}
                      </p>
                    )}
                  </>
                )}
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
