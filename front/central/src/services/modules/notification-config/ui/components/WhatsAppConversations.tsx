'use client';

import { useState, useEffect, useCallback, useRef } from 'react';
import Link from 'next/link';
import { useSearchParams } from 'next/navigation';
import { useSSE } from '@/shared/hooks/use-sse';
import {
  listConversationsAction,
  getConversationMessagesAction,
  markConversationReadAction,
  sendManualMediaAction,
  sendManualReplyAction,
} from '../../infra/actions';
import type {
  ConversationSummary,
  ConversationMessage,
  ConversationMessageMedia,
  ConversationDetailResponse,
} from '../../domain/types';

interface WhatsAppConversationsProps {
  businessId?: number;
  campaignId?: number;
  fillHeight?: boolean;
}

const BELL = '\u{1F514}';

function timeAgo(dateStr: string): string {
  const diff = Date.now() - new Date(dateStr).getTime();
  const mins = Math.floor(diff / 60000);
  if (mins < 1) return 'ahora';
  if (mins < 60) return `${mins}m`;
  const hours = Math.floor(mins / 60);
  if (hours < 24) return `${hours}h`;
  const days = Math.floor(hours / 24);
  if (days < 7) return `${days}d`;
  return new Date(dateStr).toLocaleDateString('es-CO', { day: '2-digit', month: 'short' });
}

function formatMessageTime(dateStr: string): string {
  return new Date(dateStr).toLocaleTimeString('es-CO', {
    hour: '2-digit',
    minute: '2-digit',
  });
}

function formatMessageDate(dateStr: string): string {
  const d = new Date(dateStr);
  const today = new Date();
  const yesterday = new Date();
  yesterday.setDate(yesterday.getDate() - 1);

  if (d.toDateString() === today.toDateString()) return 'Hoy';
  if (d.toDateString() === yesterday.toDateString()) return 'Ayer';
  return d.toLocaleDateString('es-CO', { day: '2-digit', month: 'long', year: 'numeric' });
}

const stateLabel: Record<string, { label: string; color: string }> = {
  START: { label: 'Inicio', color: 'bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300' },
  AWAITING_CONFIRMATION: { label: 'Esperando', color: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900 dark:text-yellow-300' },
  COMPLETED: { label: 'Completada', color: 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300' },
  HANDOFF_TO_HUMAN: { label: 'Con asesor', color: 'bg-orange-100 text-orange-700 dark:bg-orange-900 dark:text-orange-300' },
  AI_ACTIVE: { label: 'IA activa', color: 'bg-indigo-100 text-indigo-700 dark:bg-indigo-900 dark:text-indigo-300' },
};

const statusIcon: Record<string, string> = {
  sent: '\u2713',
  delivered: '\u2713\u2713',
  read: '\u2713\u2713',
  failed: '\u2717',
};

const CHAT_BG_STORAGE_KEY = 'wa_chat_background_v1';

const DOODLE_ICONS = `
  <g id="wa-star"><polygon points="10,1 12.5,7 19,7 13.8,11 15.8,18 10,14 4.2,18 6.2,11 1,7 7.5,7"/></g>
  <g id="wa-heart"><path d="M10,17 C4,12 1,8.5 1,5.5 C1,2.5 3.2,1 5.5,1 C7.3,1 9,2 10,3.8 C11,2 12.7,1 14.5,1 C16.8,1 19,2.5 19,5.5 C19,8.5 16,12 10,17 Z"/></g>
  <g id="wa-note"><path d="M9,2 L9,14"/><path d="M9,2 L15,4"/><ellipse cx="6.3" cy="15" rx="3" ry="2.1" transform="rotate(-15 6.3 15)"/></g>
  <g id="wa-bubble"><rect x="1" y="1" width="18" height="11" rx="4"/><path d="M6,12 L4,16 L9,12"/></g>
  <g id="wa-sun"><circle cx="9" cy="9" r="3.2"/><path d="M9,1 L9,3 M9,15 L9,17 M1,9 L3,9 M15,9 L17,9 M3.5,3.5 L5,5 M13,13 L14.5,14.5 M14.5,3.5 L13,5 M5,13 L3.5,14.5"/></g>
  <g id="wa-cloud"><path d="M6,13 C2.7,13 1,10.8 1,8.8 C1,6.8 2.6,5.2 4.5,5.2 C4.9,2.8 7,1 9.5,1 C12.2,1 14.4,3 14.7,5.6 C17,5.8 19,7.6 19,9.8 C19,11.8 17.3,13 15,13 Z"/></g>
  <g id="wa-plane"><path d="M1,10 L19,2 L12,18 L9,11 Z M1,10 L9,11"/></g>
  <g id="wa-camera"><rect x="1" y="4" width="18" height="11" rx="2"/><rect x="6" y="1" width="6" height="3" rx="1"/><circle cx="10" cy="9.5" r="3.2"/></g>
  <g id="wa-balloon"><ellipse cx="7" cy="7" rx="6" ry="7"/><path d="M7,14 L6,16 L8,17 L6.5,19"/></g>
  <g id="wa-flower"><circle cx="5" cy="5" r="2.6"/><circle cx="13" cy="5" r="2.6"/><circle cx="5" cy="13" r="2.6"/><circle cx="13" cy="13" r="2.6"/></g>
  <g id="wa-cup"><path d="M2,5 L2,13 A5,5 0 0 0 12,13 L12,5 Z"/><path d="M12,6 C15,6 15,11 12,11"/><path d="M4,2 C4.5,3 3.5,3.5 4,4.5 M7,2 C7.5,3 6.5,3.5 7,4.5"/></g>
  <g id="wa-umbrella"><path d="M1,10 A9,9 0 0 1 19,10"/><path d="M10,10 L10,17 C10,18.5 8,18.5 8,17"/><path d="M10,10 L10,2"/></g>
`;

const DOODLE_LAYOUT = [
  { id: 'wa-star', x: 12, y: 8 }, { id: 'wa-cloud', x: 100, y: 4 }, { id: 'wa-heart', x: 190, y: 14 },
  { id: 'wa-note', x: 60, y: 40 }, { id: 'wa-balloon', x: 150, y: 55 }, { id: 'wa-sun', x: 225, y: 60 },
  { id: 'wa-bubble', x: 5, y: 75 }, { id: 'wa-umbrella', x: 95, y: 90 }, { id: 'wa-flower', x: 200, y: 100 },
  { id: 'wa-plane', x: 30, y: 120 }, { id: 'wa-cup', x: 140, y: 130 }, { id: 'wa-camera', x: 225, y: 140 },
  { id: 'wa-star', x: 70, y: 160 }, { id: 'wa-heart', x: 10, y: 185 }, { id: 'wa-cloud', x: 170, y: 175 },
  { id: 'wa-balloon', x: 105, y: 200 }, { id: 'wa-note', x: 210, y: 210 }, { id: 'wa-bubble', x: 45, y: 225 },
  { id: 'wa-sun', x: 150, y: 230 }, { id: 'wa-flower', x: 5, y: 230 },
] as const;

function buildDoodlePattern(strokeColor: string): string {
  const uses = DOODLE_LAYOUT
    .map((d) => `<use href="#${d.id}" x="${d.x}" y="${d.y}"/>`)
    .join('');
  const svg = `<svg xmlns="http://www.w3.org/2000/svg" width="260" height="260" viewBox="0 0 260 260">` +
    `<defs>${DOODLE_ICONS}</defs>` +
    `<g fill="none" stroke="${strokeColor}" stroke-width="1.3" stroke-linecap="round" stroke-linejoin="round">${uses}</g>` +
    `</svg>`;
  return `url("data:image/svg+xml,${encodeURIComponent(svg)}")`;
}

const DEFAULT_DOODLE_STROKE = '#9C92AC';
const DOODLE_PATTERN = buildDoodlePattern(DEFAULT_DOODLE_STROKE);

const CHAT_BACKGROUNDS: { color: string | null; label: string }[] = [
  { color: null, label: 'Predeterminado' },
  { color: '#fdf3e7', label: 'Crema' },
  { color: '#f5e2c8', label: 'Arena' },
  { color: '#f0d3a8', label: 'Trigo' },
  { color: '#e8bd8a', label: 'Durazno' },
  { color: '#dba36c', label: 'Canela clara' },
  { color: '#c9853f', label: 'Ámbar' },
  { color: '#e0a15a', label: 'Miel' },
  { color: '#d99a4e', label: 'Mostaza clara' },
  { color: '#e08a5c', label: 'Coral' },
  { color: '#a83e3e', label: 'Ladrillo' },
  { color: '#8a8a3c', label: 'Oliva cálido' },
  { color: '#e4d9f5', label: 'Lavanda' },
  { color: '#c3a6e8', label: 'Lila' },
  { color: '#9b6fd1', label: 'Morado' },
];

function hexToRgb(hex: string): [number, number, number] {
  const clean = hex.replace('#', '');
  return [
    parseInt(clean.substring(0, 2), 16),
    parseInt(clean.substring(2, 4), 16),
    parseInt(clean.substring(4, 6), 16),
  ];
}

function relativeLuminance([r, g, b]: [number, number, number]): number {
  const toLinear = (c: number) => {
    const s = c / 255;
    return s <= 0.03928 ? s / 12.92 : Math.pow((s + 0.055) / 1.055, 2.4);
  };
  return 0.2126 * toLinear(r) + 0.7152 * toLinear(g) + 0.0722 * toLinear(b);
}

function wcagContrast(rgbA: [number, number, number], rgbB: [number, number, number]): number {
  const lA = relativeLuminance(rgbA);
  const lB = relativeLuminance(rgbB);
  const lighter = Math.max(lA, lB);
  const darker = Math.min(lA, lB);
  return (lighter + 0.05) / (darker + 0.05);
}

function hexIsDark(hex: string): boolean {
  const rgb = hexToRgb(hex);
  const contrastWithWhite = wcagContrast(rgb, [255, 255, 255]);
  const contrastWithBlack = wcagContrast(rgb, [17, 24, 39]);
  return contrastWithWhite > contrastWithBlack;
}

function shadeColor(hex: string, percent: number): string {
  const [r, g, b] = hexToRgb(hex);
  const target = percent > 0 ? 255 : 0;
  const p = Math.abs(percent) / 100;
  const mix = (channel: number) => Math.round((target - channel) * p + channel);
  const toHex = (channel: number) => channel.toString(16).padStart(2, '0');
  return `#${toHex(mix(r))}${toHex(mix(g))}${toHex(mix(b))}`;
}

function frameColor(hex: string): string {
  return hexIsDark(hex) ? shadeColor(hex, 22) : shadeColor(hex, -22);
}

function mutedTextColor(hex: string): string {
  const rgb = hexToRgb(hex);
  const base: [number, number, number] = hexIsDark(hex) ? [255, 255, 255] : [17, 24, 39];
  const mix = (channel: number, baseChannel: number) => Math.round(baseChannel * 0.7 + channel * 0.3);
  return `rgb(${mix(rgb[0], base[0])}, ${mix(rgb[1], base[1])}, ${mix(rgb[2], base[2])})`;
}

function lightAccentColor(hex: string): string {
  return shadeColor(hex, 55);
}

interface ChatBackgroundPref {
  color: string | null;
  doodles: boolean;
}

function loadChatBackgroundPref(): ChatBackgroundPref {
  if (typeof window === 'undefined') return { color: null, doodles: true };
  try {
    const raw = window.localStorage.getItem(CHAT_BG_STORAGE_KEY);
    if (!raw) return { color: null, doodles: true };
    const parsed = JSON.parse(raw);
    return {
      color: typeof parsed.color === 'string' ? parsed.color : null,
      doodles: typeof parsed.doodles === 'boolean' ? parsed.doodles : true,
    };
  } catch {
    return { color: null, doodles: true };
  }
}

function chatBackgroundStyle(pref: ChatBackgroundPref): React.CSSProperties {
  if (!pref.color) {
    return { backgroundImage: DOODLE_PATTERN };
  }
  if (!pref.doodles) {
    return { backgroundColor: pref.color };
  }
  const strokeColor = hexIsDark(pref.color) ? 'rgba(255,255,255,0.22)' : 'rgba(0,0,0,0.14)';
  return {
    backgroundColor: pref.color,
    backgroundImage: buildDoodlePattern(strokeColor),
  };
}

function ChatBackgroundModal({
  pref,
  onChange,
  onClose,
}: {
  pref: ChatBackgroundPref;
  onChange: (next: ChatBackgroundPref) => void;
  onClose: () => void;
}) {
  return (
    <div className="absolute inset-0 z-20 flex bg-white">
      <div className="flex w-full max-w-[260px] flex-col border-r border-gray-200">
        <div className="flex items-center gap-2 border-b border-gray-200 px-3 py-2.5">
          <button
            onClick={onClose}
            className="rounded-full p-1 text-gray-500 hover:bg-gray-100"
            title="Atrás"
          >
            <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M15 19l-7-7 7-7" />
            </svg>
          </button>
          <span className="text-xs font-medium text-gray-900">Establece el fondo del chat</span>
        </div>

        <label className="flex items-center gap-2 border-b border-gray-100 px-3 py-2 text-xs text-gray-700">
          <input
            type="checkbox"
            checked={pref.doodles}
            onChange={(e) => onChange({ ...pref, doodles: e.target.checked })}
            className="h-3.5 w-3.5 accent-green-500"
          />
          {'Añadir WhatsApp Doodles'}
        </label>

        <div className="grid flex-1 grid-cols-3 gap-1.5 overflow-y-auto p-2">
          {CHAT_BACKGROUNDS.map((bg) => {
            const isSelected = pref.color === bg.color;
            return (
              <button
                key={bg.label}
                onClick={() => onChange({ ...pref, color: bg.color })}
                title={bg.label}
                className={`aspect-square rounded-md border-2 transition-colors ${isSelected ? 'border-green-500' : 'border-transparent hover:border-gray-300'}`}
                style={bg.color ? { backgroundColor: bg.color } : { backgroundColor: '#3a3a3a' }}
              >
                {!bg.color && (
                  <span className="flex h-full w-full items-center justify-center text-center text-[8px] leading-tight text-white/80">
                    Predeterm.
                  </span>
                )}
              </button>
            );
          })}
        </div>
      </div>

      <div className="flex flex-1 flex-col">
        <div className="border-b border-gray-200 px-3 py-2.5 text-center">
          <span className="text-xs font-medium text-gray-900">Vista previa del fondo</span>
        </div>
        <div className="flex-1" style={chatBackgroundStyle(pref)} />
      </div>
    </div>
  );
}

const orderTagClass =
  'inline-flex max-w-[140px] items-center gap-1 truncate rounded-full border border-sky-200 bg-sky-50 px-1.5 py-0 text-[10px] font-medium text-sky-700 dark:border-sky-800 dark:bg-sky-900 dark:text-sky-300';

const campaignTagClass =
  'inline-flex max-w-[160px] items-center gap-1 truncate rounded-full border border-fuchsia-200 bg-fuchsia-50 px-1.5 py-0 text-[10px] font-medium text-fuchsia-700 dark:border-fuchsia-800 dark:bg-fuchsia-900 dark:text-fuchsia-300';

function OrderTag({ orderNumber, orderId, linked }: { orderNumber: string; orderId?: string; linked?: boolean }) {
  const label = `Orden ${orderNumber.startsWith('#') ? orderNumber : `#${orderNumber}`}`;
  if (linked && orderId) {
    return (
      <Link
        href={`/orders?order_id=${orderId}`}
        className={`${orderTagClass} hover:bg-sky-100 dark:hover:bg-sky-900/60`}
        title="Abrir la orden"
      >
        <span className="h-1.5 w-1.5 shrink-0 rounded-full bg-sky-500" />
        <span className="truncate">{label}</span>
      </Link>
    );
  }
  return (
    <span className={orderTagClass} title={label}>
      <span className="h-1.5 w-1.5 shrink-0 rounded-full bg-sky-500" />
      <span className="truncate">{label}</span>
    </span>
  );
}

const ATTACHMENT_ACCEPT =
  'image/jpeg,image/png,video/mp4,video/3gpp,audio/aac,audio/mp4,audio/mpeg,audio/amr,audio/ogg,application/pdf,.doc,.docx,.xls,.xlsx,.ppt,.pptx,.txt,.csv';

function formatBytes(size: number): string {
  if (!size) return '';
  if (size < 1024) return `${size} B`;
  if (size < 1024 * 1024) return `${Math.round(size / 1024)} KB`;
  return `${(size / (1024 * 1024)).toFixed(1)} MB`;
}

function MessageMediaView({ media, onLoad }: { media: ConversationMessageMedia; onLoad?: () => void }) {
  if (!media.available || !media.url) {
    return (
      <div className="mb-1 rounded-md border border-dashed border-gray-300 px-2 py-1.5 text-[11px] text-gray-500 dark:border-gray-600 dark:text-gray-400">
        {'Archivo no disponible'}
      </div>
    );
  }
  if (media.type === 'image' || media.type === 'sticker') {
    return (
      <a href={media.url} target="_blank" rel="noreferrer" className="mb-1 block">
        <img src={media.url} alt={media.filename || 'Imagen'} onLoad={onLoad} className="max-h-64 max-w-full rounded-md object-contain" />
      </a>
    );
  }
  if (media.type === 'video') {
    return <video src={media.url} controls onLoadedMetadata={onLoad} className="mb-1 max-h-64 max-w-full rounded-md" />;
  }
  if (media.type === 'audio') {
    return <audio src={media.url} controls className="mb-1 w-64 max-w-full" />;
  }
  return (
    <a
      href={media.url}
      target="_blank"
      rel="noreferrer"
      className="mb-1 flex items-center gap-2 rounded-md border border-gray-200 bg-white/70 px-2 py-1.5 hover:bg-white dark:border-gray-600 dark:bg-gray-800"
    >
      <span className="text-lg">{'\u{1F4C4}'}</span>
      <span className="min-w-0">
        <span className="block truncate text-xs font-medium">{media.filename || 'Documento'}</span>
        <span className="block text-[10px] text-gray-500">{formatBytes(media.size)}</span>
      </span>
    </a>
  );
}

function OptOutTag() {
  return (
    <span className="inline-flex items-center gap-1 rounded-full border border-red-200 bg-red-50 px-1.5 py-0 text-[10px] font-semibold text-red-700 dark:border-red-800 dark:bg-red-900 dark:text-red-300">
      <span className="h-1.5 w-1.5 shrink-0 rounded-full bg-red-500" />
      {'No quiere recibir m\u00e1s'}
    </span>
  );
}

function CampaignTag({ name }: { name: string }) {
  return (
    <span className={campaignTagClass} title={`Campa\u00f1a: ${name}`}>
      <span className="h-1.5 w-1.5 shrink-0 rounded-full bg-fuchsia-500" />
      <span className="truncate">{name}</span>
    </span>
  );
}

export function WhatsAppConversations({ businessId, campaignId, fillHeight = false }: WhatsAppConversationsProps) {
  const [conversations, setConversations] = useState<ConversationSummary[]>([]);
  const [listLoading, setListLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [totalPages, setTotalPages] = useState(0);
  const [unreadConversations, setUnreadConversations] = useState(0);
  const [page, setPage] = useState(1);
  const pageSize = 20;

  const [stateFilter, setStateFilter] = useState('');
  const [phoneSearch, setPhoneSearch] = useState('');

  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [detail, setDetail] = useState<ConversationDetailResponse | null>(null);
  const [detailLoading, setDetailLoading] = useState(false);

  const [showChat, setShowChat] = useState(false);
  const [chatBgPref, setChatBgPref] = useState<ChatBackgroundPref>({ color: null, doodles: true });
  const [showBgPicker, setShowBgPicker] = useState(false);

  useEffect(() => { setChatBgPref(loadChatBackgroundPref()); }, []);

  const updateChatBgPref = useCallback((next: ChatBackgroundPref) => {
    setChatBgPref(next);
    try {
      window.localStorage.setItem(CHAT_BG_STORAGE_KEY, JSON.stringify(next));
    } catch {
      return;
    }
  }, []);

  const [replyText, setReplyText] = useState('');
  const [sending, setSending] = useState(false);
  const [sendError, setSendError] = useState<string | null>(null);
  const [attachment, setAttachment] = useState<File | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const messagesContainerRef = useRef<HTMLDivElement>(null);
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const selectedIdRef = useRef<string | null>(null);
  useEffect(() => { selectedIdRef.current = selectedId; }, [selectedId]);
  const searchParams = useSearchParams();
  const requestedConversation = searchParams.get('conversation');
  const openedFromUrl = useRef<string | null>(null);

  const fetchConversations = useCallback(async () => {
    setListLoading(true);
    try {
      const res = await listConversationsAction({
        business_id: businessId ?? 0,
        state: stateFilter || undefined,
        phone: phoneSearch || undefined,
        campaign_id: campaignId,
        page,
        page_size: pageSize,
      });
      if (res.success) {
        setConversations(res.data);
        setTotal(res.total);
        setTotalPages(res.total_pages);
        setUnreadConversations(res.unread_conversations ?? 0);
      }
    } catch {
      setConversations([]);
    } finally {
      setListLoading(false);
    }
  }, [businessId, campaignId, stateFilter, phoneSearch, page]);

  useEffect(() => { setPage(1); }, [stateFilter, phoneSearch, businessId, campaignId]);
  useEffect(() => { fetchConversations(); }, [fetchConversations]);


  const openConversation = useCallback(async (convId: string) => {
    setSelectedId(convId);
    setShowChat(true);
    setDetailLoading(true);

    const target = conversations.find((c) => c.id === convId);
    if (target && target.unread_count > 0) {
      setUnreadConversations((n) => Math.max(0, n - 1));
      setConversations((prev) => prev.map((c) => (c.id === convId ? { ...c, unread_count: 0 } : c)));
    }
    markConversationReadAction(convId, businessId ?? 0);

    try {
      const res = await getConversationMessagesAction(convId, businessId ?? 0);
      if (res.success && res.data) {
        setDetail(res.data);
      }
    } catch {
      setDetail(null);
    } finally {
      setDetailLoading(false);
    }
  }, [businessId, conversations]);

  useEffect(() => {
    if (!requestedConversation || listLoading || openedFromUrl.current === requestedConversation) return;
    openedFromUrl.current = requestedConversation;
    void openConversation(requestedConversation);
  }, [requestedConversation, listLoading, openConversation]);

  const followLatest = useRef(true);

  const scrollToLatest = useCallback((behavior: ScrollBehavior = 'smooth') => {
    const container = messagesContainerRef.current;
    if (container) {
      container.scrollTo({ top: container.scrollHeight, behavior });
    }
  }, []);

  const handleMediaLoaded = useCallback(() => {
    if (followLatest.current) scrollToLatest('auto');
  }, [scrollToLatest]);

  useEffect(() => {
    if (detail) {
      followLatest.current = true;
      scrollToLatest();
    }
  }, [detail, scrollToLatest]);

  useEffect(() => {
    setReplyText('');
    setSendError(null);
    setAttachment(null);
  }, [selectedId]);

  const pickAttachment = (file: File | undefined) => {
    if (!file) return;
    const isImage = file.type === 'image/jpeg' || file.type === 'image/png';
    const limitMb = isImage ? 5 : 10;
    if (file.size > limitMb * 1024 * 1024) {
      setSendError(`El archivo supera ${limitMb} MB`);
      return;
    }
    setSendError(null);
    setAttachment(file);
  };

  const isWindowActive = (() => {
    if (!detail?.messages?.length) return false;
    const lastInbound = [...detail.messages].reverse().find(m => m.direction === 'inbound');
    if (!lastInbound) return false;
    const diffHours = (Date.now() - new Date(lastInbound.created_at).getTime()) / 3600000;
    return diffHours < 24;
  })();

  const handleSend = useCallback(async () => {
    if (attachment && selectedId && detail) {
      setSending(true);
      setSendError(null);
      const formData = new FormData();
      formData.append('file', attachment);
      formData.append('caption', replyText.trim());
      formData.append('phone_number', detail.phone_number);
      formData.append('business_id', String(businessId ?? 0));
      const res = await sendManualMediaAction(selectedId, formData);
      if (!res.success) {
        setSendError(res.error ?? 'No se pudo enviar el archivo');
      } else {
        setAttachment(null);
        setReplyText('');
        const refreshed = await getConversationMessagesAction(selectedId, businessId ?? 0);
        if (refreshed.success && refreshed.data) setDetail(refreshed.data);
      }
      setSending(false);
      return;
    }
    if (!replyText.trim() || !selectedId || !detail) return;
    setSending(true);
    setSendError(null);

    const phone = detail.phone_number;
    const biz = businessId ?? 0;
    const text = replyText.trim();

    const optimisticMsg: ConversationMessage = {
      id: `optimistic-${Date.now()}`,
      direction: 'outbound',
      message_id: '',
      template_name: '',
      content: text,
      status: 'sent',
      created_at: new Date().toISOString(),
    };
    setDetail(prev => prev ? {
      ...prev,
      messages: [...prev.messages, optimisticMsg],
    } : prev);
    setReplyText('');

    const res = await sendManualReplyAction(selectedId, phone, biz, text);
    if (!res.success) {
      setSendError(res.error ?? 'Error al enviar');
      setDetail(prev => prev ? {
        ...prev,
        messages: prev.messages.filter(m => m.id !== optimisticMsg.id),
      } : prev);
      setReplyText(text);
    }
    setSending(false);
  }, [replyText, selectedId, detail, businessId, attachment]);

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  useSSE({
    businessId: businessId ?? 0,
    eventTypes: [
      'whatsapp.message_received',
      'whatsapp.conversation_started',
      'whatsapp.message_status_updated',
    ],
    onMessage: (event: MessageEvent) => {
      try {
        const envelope = JSON.parse(event.data) as {
          type: string;
          timestamp: string;
          data: {
            conversation_id?: string;
            message_id?: string;
            content?: string;
            status?: string;
          };
        };
        const type = envelope.type || event.type;
        const data = envelope.data || {};

        if (type === 'whatsapp.message_received') {
          if (data.conversation_id && data.conversation_id === selectedIdRef.current) {
            const newMsg: ConversationMessage = {
              id: data.message_id || `sse-${Date.now()}`,
              direction: 'inbound',
              message_id: data.message_id || '',
              template_name: '',
              content: data.content || '',
              status: 'delivered',
              created_at: envelope.timestamp || new Date().toISOString(),
            };
            setDetail(prev => prev ? { ...prev, messages: [...prev.messages, newMsg] } : prev);
            markConversationReadAction(data.conversation_id, businessId ?? 0).finally(() => fetchConversations());
          } else {
            fetchConversations();
          }
        } else if (type === 'whatsapp.conversation_started') {
          fetchConversations();
        } else if (type === 'whatsapp.message_status_updated' && data.message_id && data.status) {
          setDetail(prev => {
            if (!prev) return prev;
            return {
              ...prev,
              messages: prev.messages.map(m =>
                m.message_id === data.message_id
                  ? { ...m, status: data.status as string }
                  : m
              ),
            };
          });
        }
      } catch {
        return;
      }
    },
  });

  const pressedButtons: Record<string, string> = {};
  if (detail) {
    detail.messages.forEach((msg, index) => {
      if (msg.direction !== 'outbound' || !msg.buttons?.length) return;
      const reply = detail.messages.slice(index + 1).find((next) => next.direction === 'inbound');
      if (!reply) return;
      const answer = reply.content.trim().toLowerCase();
      const pressed = msg.buttons.find((button) => button.text.trim().toLowerCase() === answer);
      if (pressed) pressedButtons[msg.id] = pressed.text;
    });
  }

  const groupedMessages: { date: string; messages: ConversationMessage[] }[] = [];
  if (detail) {
    let currentDate = '';
    for (const msg of detail.messages) {
      const d = formatMessageDate(msg.created_at);
      if (d !== currentDate) {
        currentDate = d;
        groupedMessages.push({ date: d, messages: [] });
      }
      groupedMessages[groupedMessages.length - 1].messages.push(msg);
    }
  }

  const containerStyle: React.CSSProperties = fillHeight
    ? { height: 'calc(100vh - 170px)', minHeight: 480 }
    : { height: 600 };

  const themed = !!chatBgPref.color;
  const forceDark = themed && hexIsDark(chatBgPref.color as string);
  const themedSurfaceStyle: React.CSSProperties | undefined = themed
    ? { backgroundColor: chatBgPref.color as string }
    : undefined;
  const themedBorderStyle: React.CSSProperties | undefined = themed
    ? { borderColor: frameColor(chatBgPref.color as string) }
    : undefined;
  const mutedStyle: React.CSSProperties | undefined = themed
    ? { color: mutedTextColor(chatBgPref.color as string) }
    : undefined;
  const lightAccentStyle: React.CSSProperties | undefined = themed
    ? { backgroundColor: lightAccentColor(chatBgPref.color as string), color: '#111827', borderColor: frameColor(chatBgPref.color as string) }
    : undefined;
  const scrollbarStyle: React.CSSProperties | undefined = themed
    ? ({ scrollbarColor: `${frameColor(chatBgPref.color as string)} ${lightAccentColor(chatBgPref.color as string)}` } as React.CSSProperties)
    : undefined;

  return (
    <div
      className={`relative flex flex-col overflow-hidden rounded-lg border border-gray-200 bg-white shadow-sm dark:border-gray-700 dark:bg-gray-800 ${forceDark ? 'dark' : ''}`}
      style={{ ...containerStyle, ...themedSurfaceStyle, ...themedBorderStyle }}
    >
      {showBgPicker && (
        <ChatBackgroundModal
          pref={chatBgPref}
          onChange={updateChatBgPref}
          onClose={() => setShowBgPicker(false)}
        />
      )}
      <div className="flex min-h-0 flex-1">
        <div
          className={`w-full lg:w-[360px] border-r border-gray-200 dark:border-gray-700 flex flex-col min-h-0 ${showChat ? 'hidden lg:flex' : 'flex'}`}
          style={themedBorderStyle}
        >
          <div className="p-3 border-b border-gray-200 dark:border-gray-700 space-y-2">
            <div className="flex items-center gap-2">
              <div className="w-6 h-6 rounded-full bg-green-500 flex items-center justify-center shrink-0">
                <svg className="w-3.5 h-3.5 text-white" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347z" />
                  <path d="M12 2C6.477 2 2 6.477 2 12c0 1.89.525 3.66 1.438 5.168L2 22l4.832-1.438A9.955 9.955 0 0012 22c5.523 0 10-4.477 10-10S17.523 2 12 2zm0 18a8 8 0 01-4.243-1.214l-.252-.149-2.868.852.852-2.868-.149-.252A8 8 0 1112 20z" />
                </svg>
              </div>
              <h3 data-tour="notif-conversations-list" className="text-sm font-medium text-gray-900 dark:text-white">Conversaciones</h3>
              <span className="text-xs text-gray-400 dark:text-gray-500" style={mutedStyle}>({total})</span>
              <div className="ml-auto flex items-center gap-1">
                {unreadConversations > 0 && (
                  <span className="rounded-full bg-green-500 px-2 py-0.5 text-[10px] font-semibold text-white">
                    {`${unreadConversations} sin leer`}
                  </span>
                )}
                <button
                  onClick={() => setShowBgPicker(true)}
                  className="rounded-full p-1.5 text-gray-500 transition-colors hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700"
                  title="Fondo del chat"
                >
                  <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14M14 8h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                  </svg>
                </button>
              </div>
            </div>
            <p className="text-[10px] text-gray-400 dark:text-gray-500" style={mutedStyle}>
              {'Los mensajes y archivos se eliminan autom\u00e1ticamente despu\u00e9s de 1 a\u00f1o.'}
            </p>
            <input
              type="text"
              value={phoneSearch}
              onChange={(e) => setPhoneSearch(e.target.value)}
              placeholder={'Buscar por tel\u00e9fono...'}
              className="w-full px-3 py-1.5 text-xs border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-400 dark:placeholder-gray-500 focus:outline-none focus:ring-1 focus:ring-green-500"
              style={lightAccentStyle}
            />
            <select
              value={stateFilter}
              onChange={(e) => setStateFilter(e.target.value)}
              className="w-full px-3 py-1.5 text-xs border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:outline-none focus:ring-1 focus:ring-green-500"
              style={lightAccentStyle}
            >
              <option value="">Todos los estados</option>
              <option value="START">Inicio</option>
              <option value="AWAITING_CONFIRMATION">Esperando {'confirmaci\u00f3n'}</option>
              <option value="COMPLETED">Completada</option>
            </select>
          </div>

          <div className="flex-1 overflow-y-auto" style={scrollbarStyle}>
            {listLoading ? (
              Array.from({ length: 6 }).map((_, i) => (
                <div key={i} className="p-3 border-b border-gray-100 dark:border-gray-700 animate-pulse">
                  <div className="flex gap-3">
                    <div className="w-10 h-10 rounded-full bg-gray-200 dark:bg-gray-600 shrink-0" />
                    <div className="flex-1 space-y-2">
                      <div className="h-3 bg-gray-200 dark:bg-gray-600 rounded w-3/4" />
                      <div className="h-3 bg-gray-200 dark:bg-gray-600 rounded w-1/2" />
                    </div>
                  </div>
                </div>
              ))
            ) : conversations.length === 0 ? (
              <div className="flex flex-col items-center justify-center h-full text-gray-400 dark:text-gray-500 p-4" style={mutedStyle}>
                <svg className="w-12 h-12 mb-2 opacity-30" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                </svg>
                <p className="text-xs">No hay conversaciones</p>
              </div>
            ) : (
              conversations.map((conv) => {
                const isActive = selectedId === conv.id;
                const isSystemAlert = conv.conversation_type === 'system_alert';
                const state = stateLabel[conv.current_state] || { label: conv.current_state, color: 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300' };
                const unread = conv.unread_count > 0;
                const optedOut = conv.opted_out;
                return (
                  <button
                    key={conv.id}
                    onClick={() => openConversation(conv.id)}
                    className={`w-full text-left p-3 border-b border-gray-100 dark:border-gray-700 hover:bg-black/5 dark:hover:bg-white/10 transition-colors ${isActive ? 'bg-black/5 dark:bg-white/10 border-l-2 border-l-green-500' : optedOut ? 'border-l-2 border-l-red-500 bg-red-50/60 dark:bg-red-500/10' : ''}`}
                  >
                    <div className="flex gap-3">
                      <div className={`w-10 h-10 rounded-full flex items-center justify-center text-white text-xs font-bold shrink-0 ${isSystemAlert ? 'bg-gradient-to-br from-purple-400 to-purple-600' : optedOut ? 'bg-gradient-to-br from-red-400 to-red-600' : 'bg-gradient-to-br from-green-400 to-green-600'}`}>
                        {isSystemAlert ? BELL : conv.phone_number.slice(-2)}
                      </div>

                      <div className="flex-1 min-w-0">
                        <div className="flex items-center justify-between">
                          <span className={`text-sm text-gray-900 dark:text-white truncate ${unread ? 'font-semibold' : 'font-medium'}`}>
                            {conv.customer_name || conv.phone_number}
                          </span>
                          <span className={`text-[10px] shrink-0 ml-2 ${unread ? 'font-semibold text-green-600 dark:text-green-400' : 'text-gray-400 dark:text-gray-500'}`} style={unread ? undefined : mutedStyle}>
                            {timeAgo(conv.last_activity)}
                          </span>
                        </div>
                        {conv.customer_name && (
                          <p className="text-[11px] text-gray-500 dark:text-gray-400 truncate" style={mutedStyle}>{conv.phone_number}</p>
                        )}

                        <div className="flex items-center gap-1 mt-0.5">
                          {conv.last_message_direction === 'outbound' && (
                            <span className={`text-[10px] shrink-0 ${conv.last_message_status === 'read' ? 'text-blue-500' : conv.last_message_status === 'failed' ? 'text-red-400' : 'text-gray-400'}`}>
                              {statusIcon[conv.last_message_status] || ''}
                            </span>
                          )}
                          <p className={`text-xs truncate ${unread ? 'font-medium text-gray-800 dark:text-gray-100' : 'text-gray-500 dark:text-gray-400'}`} style={unread ? undefined : mutedStyle}>
                            {conv.last_message_content || 'Sin mensajes'}
                          </p>
                          {unread && (
                            <span className="ml-auto flex h-4 min-w-4 shrink-0 items-center justify-center rounded-full bg-green-500 px-1 text-[10px] font-semibold text-white">
                              {conv.unread_count > 99 ? '99+' : conv.unread_count}
                            </span>
                          )}
                        </div>

                        <div className="flex flex-wrap items-center gap-1 mt-1">
                          {isSystemAlert ? (
                            <span className="inline-block px-1.5 py-0 rounded-full text-[9px] font-medium bg-purple-100 text-purple-700 dark:bg-purple-900 dark:text-purple-300">
                              Aviso del sistema
                            </span>
                          ) : (
                            <>
                              {conv.order_number && (
                                <OrderTag orderNumber={conv.order_number} orderId={conv.order_id} />
                              )}
                              {optedOut && <OptOutTag />}
                              {conv.campaign_name && <CampaignTag name={conv.campaign_name} />}
                              <span className={`inline-block px-1.5 py-0 rounded-full text-[9px] font-medium ${state.color}`}>
                                {state.label}
                              </span>
                            </>
                          )}
                          <span className="text-[10px] text-gray-400 dark:text-gray-500 ml-auto" style={mutedStyle}>
                            {conv.message_count} msg
                          </span>
                        </div>
                      </div>
                    </div>
                  </button>
                );
              })
            )}
          </div>

          {totalPages > 1 && (
            <div className="p-2 border-t border-gray-200 dark:border-gray-700 flex items-center justify-between">
              <span className="text-[10px] text-gray-400 dark:text-gray-500" style={mutedStyle}>{page}/{totalPages}</span>
              <div className="flex gap-1">
                <button
                  onClick={() => setPage((p) => Math.max(1, p - 1))}
                  disabled={page === 1}
                  className="px-2 py-0.5 text-[10px] border border-gray-300 dark:border-gray-600 rounded disabled:opacity-30 hover:bg-gray-50 dark:hover:bg-gray-700 text-gray-600 dark:text-gray-300 transition-colors"
                >
                  Ant
                </button>
                <button
                  onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                  disabled={page === totalPages}
                  className="px-2 py-0.5 text-[10px] border border-gray-300 dark:border-gray-600 rounded disabled:opacity-30 hover:bg-gray-50 dark:hover:bg-gray-700 text-gray-600 dark:text-gray-300 transition-colors"
                >
                  Sig
                </button>
              </div>
            </div>
          )}
        </div>

        <div className={`flex-1 flex flex-col min-h-0 ${showChat ? 'flex' : 'hidden lg:flex'}`}>
          {!selectedId ? (
            <div className="flex-1 flex flex-col items-center justify-center text-gray-300 dark:text-gray-600">
              <svg className="w-20 h-20 mb-4 opacity-30" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
              </svg>
              <p className="text-sm">Selecciona una {'conversaci\u00f3n'}</p>
            </div>
          ) : detailLoading ? (
            <div className="flex-1 flex items-center justify-center">
              <div className="flex flex-col items-center gap-3">
                <div className="w-8 h-8 border-2 border-green-500 border-t-transparent rounded-full animate-spin" />
                <span className="text-xs text-gray-400 dark:text-gray-500" style={mutedStyle}>Cargando mensajes...</span>
              </div>
            </div>
          ) : detail ? (
            <>
              {(() => {
                const isSystemAlert = detail.conversation_type === 'system_alert';
                return (
                  <div
                    className="px-4 py-2.5 border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-700 flex items-center gap-3"
                    style={{ ...themedSurfaceStyle, ...themedBorderStyle }}
                  >
                    <button
                      onClick={() => setShowChat(false)}
                      className="lg:hidden px-2 py-1 text-xs text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded transition-colors"
                    >
                      {'\u2190 Volver'}
                    </button>
                    <div className={`w-8 h-8 rounded-full flex items-center justify-center text-white text-xs font-bold shrink-0 ${isSystemAlert ? 'bg-gradient-to-br from-purple-400 to-purple-600' : detail.opted_out ? 'bg-gradient-to-br from-red-400 to-red-600' : 'bg-gradient-to-br from-green-400 to-green-600'}`}>
                      {isSystemAlert ? BELL : detail.phone_number.slice(-2)}
                    </div>
                    <div className="flex-1 min-w-0">
                      {detail.customer_name ? (
                        <>
                          <p className="text-sm font-medium text-gray-900 dark:text-white truncate">{detail.customer_name}</p>
                          <p className="text-[11px] text-gray-500 dark:text-gray-400" style={mutedStyle}>{detail.phone_number}</p>
                        </>
                      ) : (
                        <p className="text-sm font-medium text-gray-900 dark:text-white">{detail.phone_number}</p>
                      )}
                      <div className="flex flex-wrap items-center gap-1.5">
                        {isSystemAlert ? (
                          <span className="inline-block px-1.5 py-0 rounded-full text-[9px] font-medium bg-purple-100 text-purple-700 dark:bg-purple-900 dark:text-purple-300">
                            Aviso del sistema
                          </span>
                        ) : (
                          <>
                            {detail.order_number && (
                              <OrderTag orderNumber={detail.order_number} orderId={detail.order_id} linked />
                            )}
                            {detail.opted_out && <OptOutTag />}
                            {detail.campaign_name && <CampaignTag name={detail.campaign_name} />}
                            {(() => {
                              const s = stateLabel[detail.current_state] || { label: detail.current_state, color: 'bg-gray-100 text-gray-600' };
                              return (
                                <span className={`inline-block px-1.5 py-0 rounded-full text-[9px] font-medium ${s.color}`}>
                                  {s.label}
                                </span>
                              );
                            })()}
                          </>
                        )}
                      </div>
                    </div>
                    <div className="flex items-center gap-2 shrink-0">
                      <button
                        onClick={() => setShowBgPicker(true)}
                        className="rounded-full p-1.5 text-gray-500 transition-colors hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700"
                        title="Fondo del chat"
                      >
                        <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                          <path strokeLinecap="round" strokeLinejoin="round" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14M14 8h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                        </svg>
                      </button>
                      <span className="text-[10px] text-gray-400 dark:text-gray-500" style={mutedStyle}>{detail.messages.length} msg</span>
                    </div>
                  </div>
                );
              })()}

              <div
                ref={messagesContainerRef}
                onScroll={(e) => {
                  const el = e.currentTarget;
                  followLatest.current = el.scrollHeight - el.scrollTop - el.clientHeight < 120;
                }}
                className="flex-1 overflow-y-auto px-4 py-3 space-y-1 selection:bg-green-500 selection:text-white"
                style={{ ...chatBackgroundStyle(chatBgPref), ...scrollbarStyle }}
              >
                {groupedMessages.map((group) => (
                  <div key={group.date}>
                    <div className="flex items-center justify-center my-3">
                      <span className="px-3 py-0.5 bg-white dark:bg-gray-700 rounded-full text-[10px] text-gray-500 dark:text-gray-400 shadow-sm border border-gray-200 dark:border-gray-600">
                        {group.date}
                      </span>
                    </div>

                    {group.messages.map((msg) => {
                      const isOutbound = msg.direction === 'outbound';
                      return (
                        <div
                          key={msg.id}
                          className={`flex mb-1 ${isOutbound ? 'justify-end' : 'justify-start'}`}
                        >
                          <div
                            className={`relative max-w-[75%] px-3 py-1.5 rounded-lg text-xs shadow-sm ${
                              isOutbound
                                ? 'bg-green-100 dark:bg-green-800 text-gray-800 dark:text-green-100 rounded-tr-none'
                                : 'bg-white dark:bg-gray-700 text-gray-800 dark:text-gray-200 rounded-tl-none border border-gray-100 dark:border-gray-600'
                            }`}
                          >
                            {msg.template_name && isOutbound && (
                              <p className="text-[9px] text-green-600 dark:text-green-400 font-medium mb-0.5">
                                {msg.template_name}
                              </p>
                            )}

                            {msg.media && <MessageMediaView media={msg.media} onLoad={handleMediaLoaded} />}

                            {(msg.content || !msg.media) && (
                              <p className="whitespace-pre-wrap break-words leading-relaxed">
                                {msg.content || (msg.template_name ? `[Plantilla: ${msg.template_name}]` : '[Sin contenido]')}
                              </p>
                            )}

                            <div className={`flex items-center gap-1 mt-0.5 ${isOutbound ? 'justify-end' : 'justify-start'}`}>
                              <span className="text-[9px] text-gray-400 dark:text-gray-500">
                                {formatMessageTime(msg.created_at)}
                              </span>
                              {isOutbound && (
                                <span className={`text-[9px] ${msg.status === 'read' ? 'text-blue-500' : msg.status === 'failed' ? 'text-red-400' : 'text-gray-400'}`}>
                                  {statusIcon[msg.status] || ''}
                                </span>
                              )}
                            </div>

                            {isOutbound && msg.buttons && msg.buttons.length > 0 && (
                              <div className="-mx-3 -mb-1.5 mt-1.5 overflow-hidden rounded-b-lg border-t border-green-200 dark:border-green-800">
                                {msg.buttons.map((button) => {
                                  const pressed = pressedButtons[msg.id] === button.text;
                                  return (
                                    <div
                                      key={button.text}
                                      className={`flex items-center justify-center gap-1.5 border-b border-green-200 px-3 py-1.5 text-xs font-medium last:border-b-0 dark:border-green-800 ${
                                        pressed
                                          ? 'bg-sky-100 text-sky-700 dark:bg-sky-900 dark:text-sky-300'
                                          : 'bg-white/60 text-sky-600 dark:bg-gray-700 dark:text-sky-400'
                                      }`}
                                      title={pressed ? 'El cliente eligi\u00f3 esta opci\u00f3n' : undefined}
                                    >
                                      <span>{button.type === 'URL' ? '\u2197' : '\u21a9'}</span>
                                      <span>{button.text}</span>
                                      {pressed && <span className="text-[10px]">{'\u2713'}</span>}
                                    </div>
                                  );
                                })}
                              </div>
                            )}
                          </div>
                        </div>
                      );
                    })}
                  </div>
                ))}
              </div>

              <div
                className="border-t border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 px-3 py-2"
                style={{ ...themedSurfaceStyle, ...themedBorderStyle }}
              >
                {detail.conversation_type === 'system_alert' ? (
                  <div className="flex items-center justify-center py-2">
                    <p className="text-[11px] text-purple-600 dark:text-purple-400">
                      {`${BELL} Este es un aviso autom\u00e1tico del sistema, no admite respuesta.`}
                    </p>
                  </div>
                ) : (
                  <>
                    {!isWindowActive && (
                      <p className="text-[10px] text-amber-600 dark:text-amber-400 mb-1 text-center">
                        {"\u26a0 La ventana de 24h est\u00e1 cerrada: el cliente no ha escrito. Meta solo permite enviarle una plantilla aprobada."}
                      </p>
                    )}
                    {sendError && (
                      <p className="text-[10px] text-red-500 mb-1 text-center">{sendError}</p>
                    )}
                    {attachment && (
                      <div className="mb-1.5 flex items-center gap-2 rounded-lg border border-gray-200 bg-gray-50 px-2 py-1 text-xs dark:border-gray-600 dark:bg-gray-700">
                        <span>{attachment.type.startsWith('image/') ? '\u{1F5BC}' : '\u{1F4CE}'}</span>
                        <span className="min-w-0 flex-1 truncate text-gray-700 dark:text-gray-200">{attachment.name}</span>
                        <span className="shrink-0 text-[10px] text-gray-400">{formatBytes(attachment.size)}</span>
                        <button
                          type="button"
                          onClick={() => setAttachment(null)}
                          className="shrink-0 text-gray-400 hover:text-red-500"
                          aria-label="Quitar archivo"
                        >
                          {'\u00d7'}
                        </button>
                      </div>
                    )}
                    <div className="flex items-end gap-2">
                      <input
                        ref={fileInputRef}
                        type="file"
                        accept={ATTACHMENT_ACCEPT}
                        className="hidden"
                        onChange={(e) => {
                          pickAttachment(e.target.files?.[0]);
                          e.target.value = '';
                        }}
                      />
                      <button
                        type="button"
                        onClick={() => fileInputRef.current?.click()}
                        disabled={sending || !isWindowActive}
                        className="flex-shrink-0 rounded-full p-2 text-gray-500 transition-colors hover:bg-gray-100 disabled:cursor-not-allowed disabled:opacity-40 dark:text-gray-300 dark:hover:bg-gray-700"
                        title={isWindowActive ? 'Adjuntar archivo' : 'Solo se pueden enviar archivos con la ventana de 24h abierta'}
                      >
                        <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                          <path strokeLinecap="round" strokeLinejoin="round" d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13" />
                        </svg>
                      </button>
                      <textarea
                        ref={textareaRef}
                        value={replyText}
                        onChange={e => setReplyText(e.target.value)}
                        onKeyDown={handleKeyDown}
                        placeholder={'Escribe un mensaje... (Enter para enviar, Shift+Enter para nueva l\u00ednea)'}
                        rows={1}
                        disabled={sending}
                        className="flex-1 resize-none bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-gray-100 text-sm rounded-xl px-3 py-2 outline-none focus:ring-2 focus:ring-green-500 placeholder-gray-400 disabled:opacity-50 max-h-32 overflow-y-auto"
                        style={{ fieldSizing: 'content' } as React.CSSProperties}
                      />
                      <button
                        onClick={handleSend}
                        disabled={sending || (!replyText.trim() && !attachment)}
                        className="flex-shrink-0 bg-green-500 hover:bg-green-600 disabled:opacity-40 disabled:cursor-not-allowed text-white rounded-full p-2 transition-colors"
                        title="Enviar mensaje"
                      >
                        {sending ? (
                          <svg className="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z" />
                          </svg>
                        ) : (
                          <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
                            <path d="M2.01 21L23 12 2.01 3 2 10l15 2-15 2z" />
                          </svg>
                        )}
                      </button>
                    </div>
                    {isWindowActive && (
                      <p className="text-[9px] text-green-600 dark:text-green-400 mt-1 text-center">
                        {'Ventana activa \u2014 el cliente respondi\u00f3 recientemente'}
                      </p>
                    )}
                  </>
                )}
              </div>
            </>
          ) : (
            <div className="flex-1 flex items-center justify-center text-gray-400">
              <p className="text-sm">Error al cargar la {'conversaci\u00f3n'}</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
