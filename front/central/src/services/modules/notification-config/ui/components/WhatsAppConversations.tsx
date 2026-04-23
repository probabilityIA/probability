'use client';

import { useState, useEffect, useCallback, useRef } from 'react';
import { useSSE } from '@/shared/hooks/use-sse';
import {
  listConversationsAction,
  getConversationMessagesAction,
  sendManualReplyAction,
  pauseAIAction,
  resumeAIAction,
} from '../../infra/actions';
import type {
  ConversationSummary,
  ConversationMessage,
  ConversationDetailResponse,
} from '../../domain/types';

interface WhatsAppConversationsProps {
  businessId?: number;
}

// ─── Helpers ───────────────────────────────────────────

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
};

const statusIcon: Record<string, string> = {
  sent: '\u2713',
  delivered: '\u2713\u2713',
  read: '\u2713\u2713',
  failed: '\u2717',
};

// ─── Component ─────────────────────────────────────────

export function WhatsAppConversations({ businessId }: WhatsAppConversationsProps) {
  // Conversation list state
  const [conversations, setConversations] = useState<ConversationSummary[]>([]);
  const [listLoading, setListLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [totalPages, setTotalPages] = useState(0);
  const [page, setPage] = useState(1);
  const pageSize = 20;

  // Filters
  const [stateFilter, setStateFilter] = useState('');
  const [phoneSearch, setPhoneSearch] = useState('');

  // Selected conversation detail
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [detail, setDetail] = useState<ConversationDetailResponse | null>(null);
  const [detailLoading, setDetailLoading] = useState(false);

  // Mobile view toggle
  const [showChat, setShowChat] = useState(false);

  // Composer state
  const [replyText, setReplyText] = useState('');
  const [sending, setSending] = useState(false);
  const [sendError, setSendError] = useState<string | null>(null);

  // AI control state — true = AI is active (human cannot write)
  const [aiPaused, setAiPaused] = useState(false);
  const [aiToggling, setAiToggling] = useState(false);

  const messagesEndRef = useRef<HTMLDivElement>(null);
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  // Ref para tracking de conversación activa sin reconectar el SSE
  const selectedIdRef = useRef<string | null>(null);
  useEffect(() => { selectedIdRef.current = selectedId; }, [selectedId]);

  // ─── Fetch conversations ────────────────
  const fetchConversations = useCallback(async () => {
    setListLoading(true);
    try {
      const res = await listConversationsAction({
        business_id: businessId ?? 0,
        state: stateFilter || undefined,
        phone: phoneSearch || undefined,
        page,
        page_size: pageSize,
      });
      if (res.success) {
        setConversations(res.data);
        setTotal(res.total);
        setTotalPages(res.total_pages);
      }
    } catch {
      setConversations([]);
    } finally {
      setListLoading(false);
    }
  }, [businessId, stateFilter, phoneSearch, page]);

  useEffect(() => { setPage(1); }, [stateFilter, phoneSearch, businessId]);
  useEffect(() => { fetchConversations(); }, [fetchConversations]);

  // ─── Fetch conversation detail ────────────────
  const openConversation = useCallback(async (convId: string) => {
    setSelectedId(convId);
    setShowChat(true);
    setDetailLoading(true);
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
  }, [businessId]);

  // Scroll to bottom when messages load
  useEffect(() => {
    if (detail && messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [detail]);

  // Reset composer and AI state when switching conversations
  useEffect(() => {
    setReplyText('');
    setSendError(null);
    setAiPaused(false); // assume AI is active until user toggles
  }, [selectedId]);

  // Ventana activa: el cliente escribió en las últimas 24h
  // Buscamos el último mensaje inbound en el historial
  const isWindowActive = (() => {
    if (!detail?.messages?.length) return false;
    const lastInbound = [...detail.messages].reverse().find(m => m.direction === 'inbound');
    if (!lastInbound) return false;
    const diffHours = (Date.now() - new Date(lastInbound.created_at).getTime()) / 3600000;
    return diffHours < 24;
  })();

  // ─── AI toggle ────────────────
  const handleToggleAI = useCallback(async () => {
    if (!selectedId || !detail || aiToggling) return;
    setAiToggling(true);
    const phone = detail.phone_number;
    const biz = businessId ?? 0;

    if (aiPaused) {
      // Resume: reactivate AI
      const res = await resumeAIAction(selectedId, phone, biz);
      if (res.success) setAiPaused(false);
    } else {
      // Pause: human takes control
      const res = await pauseAIAction(selectedId, phone, biz);
      if (res.success) setAiPaused(true);
    }
    setAiToggling(false);
  }, [selectedId, detail, businessId, aiPaused, aiToggling]);

  // ─── Send manual reply ────────────────
  const handleSend = useCallback(async () => {
    if (!replyText.trim() || !selectedId || !detail) return;
    setSending(true);
    setSendError(null);

    const phone = detail.phone_number;
    const biz = businessId ?? 0;
    const text = replyText.trim();

    // Optimistic: agregar mensaje al chat inmediatamente
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
      // Revertir optimistic
      setDetail(prev => prev ? {
        ...prev,
        messages: prev.messages.filter(m => m.id !== optimisticMsg.id),
      } : prev);
      setReplyText(text);
    }
    setSending(false);
  }, [replyText, selectedId, detail, businessId]);

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  // ─── SSE: actualizaciones en tiempo real ────────────────
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
          // Si es la conversación activa, agregar el mensaje al chat
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
          }
          // Siempre refrescar la lista para actualizar preview y contador
          fetchConversations();
        } else if (type === 'whatsapp.conversation_started') {
          fetchConversations();
        } else if (type === 'whatsapp.message_status_updated' && data.message_id && data.status) {
          // Actualizar estado visual del mensaje (✓✓ azul para leído, etc.)
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
        // ignorar errores de parse
      }
    },
  });

  // ─── Group messages by date ────────────────
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

  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden">
      {/* Header */}
      <div className="p-4 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 rounded-full bg-green-500 flex items-center justify-center">
            <svg className="w-4 h-4 text-white" viewBox="0 0 24 24" fill="currentColor">
              <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347z" />
              <path d="M12 2C6.477 2 2 6.477 2 12c0 1.89.525 3.66 1.438 5.168L2 22l4.832-1.438A9.955 9.955 0 0012 22c5.523 0 10-4.477 10-10S17.523 2 12 2zm0 18a8 8 0 01-4.243-1.214l-.252-.149-2.868.852.852-2.868-.149-.252A8 8 0 1112 20z" />
            </svg>
          </div>
          <h3 className="text-sm font-medium text-gray-900 dark:text-white">Conversaciones WhatsApp</h3>
          <span className="text-xs text-gray-400 dark:text-gray-500">({total})</span>
        </div>
        {/* Back button on mobile when chat is open */}
        {showChat && (
          <button
            onClick={() => setShowChat(false)}
            className="lg:hidden px-2 py-1 text-xs text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded transition-colors"
          >
            &larr; Volver
          </button>
        )}
      </div>

      <div className="flex" style={{ height: '600px' }}>
        {/* ═══ LEFT: Conversation List ═══ */}
        <div className={`w-full lg:w-[360px] border-r border-gray-200 dark:border-gray-700 flex flex-col ${showChat ? 'hidden lg:flex' : 'flex'}`}>
          {/* Filters */}
          <div className="p-3 border-b border-gray-200 dark:border-gray-700 space-y-2">
            <input
              type="text"
              value={phoneSearch}
              onChange={(e) => setPhoneSearch(e.target.value)}
              placeholder="Buscar por telefono..."
              className="w-full px-3 py-1.5 text-xs border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-400 dark:placeholder-gray-500 focus:outline-none focus:ring-1 focus:ring-green-500"
            />
            <select
              value={stateFilter}
              onChange={(e) => setStateFilter(e.target.value)}
              className="w-full px-3 py-1.5 text-xs border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:outline-none focus:ring-1 focus:ring-green-500"
            >
              <option value="">Todos los estados</option>
              <option value="START">Inicio</option>
              <option value="AWAITING_CONFIRMATION">Esperando confirmacion</option>
              <option value="COMPLETED">Completada</option>
            </select>
          </div>

          {/* List */}
          <div className="flex-1 overflow-y-auto">
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
              <div className="flex flex-col items-center justify-center h-full text-gray-400 dark:text-gray-500 p-4">
                <svg className="w-12 h-12 mb-2 opacity-30" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                </svg>
                <p className="text-xs">No hay conversaciones</p>
              </div>
            ) : (
              conversations.map((conv) => {
                const isActive = selectedId === conv.id;
                const state = stateLabel[conv.current_state] || { label: conv.current_state, color: 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300' };
                return (
                  <button
                    key={conv.id}
                    onClick={() => openConversation(conv.id)}
                    className={`w-full text-left p-3 border-b border-gray-100 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-750 transition-colors ${isActive ? 'bg-green-50 dark:bg-green-900/20 border-l-2 border-l-green-500' : ''}`}
                  >
                    <div className="flex gap-3">
                      {/* Avatar */}
                      <div className="w-10 h-10 rounded-full bg-gradient-to-br from-green-400 to-green-600 flex items-center justify-center text-white text-xs font-bold shrink-0">
                        {conv.phone_number.slice(-2)}
                      </div>

                      <div className="flex-1 min-w-0">
                        {/* Row 1: phone + time */}
                        <div className="flex items-center justify-between">
                          <span className="text-sm font-medium text-gray-900 dark:text-white truncate">
                            {conv.phone_number}
                          </span>
                          <span className="text-[10px] text-gray-400 dark:text-gray-500 shrink-0 ml-2">
                            {timeAgo(conv.last_activity)}
                          </span>
                        </div>

                        {/* Row 2: last message preview */}
                        <div className="flex items-center gap-1 mt-0.5">
                          {conv.last_message_direction === 'outbound' && (
                            <span className={`text-[10px] shrink-0 ${conv.last_message_status === 'read' ? 'text-blue-500' : conv.last_message_status === 'failed' ? 'text-red-400' : 'text-gray-400'}`}>
                              {statusIcon[conv.last_message_status] || ''}
                            </span>
                          )}
                          <p className="text-xs text-gray-500 dark:text-gray-400 truncate">
                            {conv.last_message_content || 'Sin mensajes'}
                          </p>
                        </div>

                        {/* Row 3: order + state badge + count */}
                        <div className="flex items-center gap-1.5 mt-1">
                          {conv.order_number && (
                            <span className="text-[10px] text-gray-400 dark:text-gray-500 font-mono">
                              #{conv.order_number}
                            </span>
                          )}
                          <span className={`inline-block px-1.5 py-0 rounded-full text-[9px] font-medium ${state.color}`}>
                            {state.label}
                          </span>
                          <span className="text-[10px] text-gray-400 dark:text-gray-500 ml-auto">
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

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="p-2 border-t border-gray-200 dark:border-gray-700 flex items-center justify-between">
              <span className="text-[10px] text-gray-400 dark:text-gray-500">{page}/{totalPages}</span>
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

        {/* ═══ RIGHT: Chat View ═══ */}
        <div className={`flex-1 flex flex-col ${showChat ? 'flex' : 'hidden lg:flex'}`}>
          {!selectedId ? (
            /* Empty state */
            <div className="flex-1 flex flex-col items-center justify-center text-gray-300 dark:text-gray-600">
              <svg className="w-20 h-20 mb-4 opacity-30" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
              </svg>
              <p className="text-sm">Selecciona una conversacion</p>
            </div>
          ) : detailLoading ? (
            /* Loading */
            <div className="flex-1 flex items-center justify-center">
              <div className="flex flex-col items-center gap-3">
                <div className="w-8 h-8 border-2 border-green-500 border-t-transparent rounded-full animate-spin" />
                <span className="text-xs text-gray-400 dark:text-gray-500">Cargando mensajes...</span>
              </div>
            </div>
          ) : detail ? (
            <>
              {/* Chat header */}
              <div className="px-4 py-3 border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-750 flex items-center gap-3">
                <div className="w-8 h-8 rounded-full bg-gradient-to-br from-green-400 to-green-600 flex items-center justify-center text-white text-xs font-bold shrink-0">
                  {detail.phone_number.slice(-2)}
                </div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium text-gray-900 dark:text-white">{detail.phone_number}</p>
                  <div className="flex items-center gap-2">
                    {detail.order_number && (
                      <span className="text-[10px] text-gray-400 dark:text-gray-500 font-mono">Orden #{detail.order_number}</span>
                    )}
                    {(() => {
                      const s = stateLabel[detail.current_state] || { label: detail.current_state, color: 'bg-gray-100 text-gray-600' };
                      return (
                        <span className={`inline-block px-1.5 py-0 rounded-full text-[9px] font-medium ${s.color}`}>
                          {s.label}
                        </span>
                      );
                    })()}
                  </div>
                </div>
                {/* AI toggle */}
                <div className="flex items-center gap-2 shrink-0">
                  <span className="text-[10px] text-gray-400 dark:text-gray-500">{detail.messages.length} msg</span>
                  <button
                    onClick={handleToggleAI}
                    disabled={aiToggling}
                    title={aiPaused ? 'Reactivar IA' : 'Pausar IA y tomar control'}
                    className={`flex items-center gap-1 px-2 py-1 rounded-full text-[10px] font-medium transition-colors disabled:opacity-50 ${
                      aiPaused
                        ? 'bg-red-100 text-red-700 hover:bg-red-200 dark:bg-red-900/40 dark:text-red-300'
                        : 'bg-green-100 text-green-700 hover:bg-green-200 dark:bg-green-900/40 dark:text-green-300'
                    }`}
                  >
                    {aiToggling ? (
                      <svg className="w-3 h-3 animate-spin" fill="none" viewBox="0 0 24 24">
                        <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z" />
                      </svg>
                    ) : (
                      <span>{aiPaused ? '🔴' : '🟢'}</span>
                    )}
                    <span>IA {aiPaused ? 'Pausada' : 'Activa'}</span>
                  </button>
                </div>
              </div>

              {/* Messages area */}
              <div
                className="flex-1 overflow-y-auto px-4 py-3 space-y-1"
                style={{ backgroundImage: 'url("data:image/svg+xml,%3Csvg width=\'60\' height=\'60\' viewBox=\'0 0 60 60\' xmlns=\'http://www.w3.org/2000/svg\'%3E%3Cg fill=\'none\' fill-rule=\'evenodd\'%3E%3Cg fill=\'%239C92AC\' fill-opacity=\'0.03\'%3E%3Cpath d=\'M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z\'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E")' }}
              >
                {groupedMessages.map((group) => (
                  <div key={group.date}>
                    {/* Date separator */}
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
                                ? 'bg-green-100 dark:bg-green-900/40 text-gray-800 dark:text-green-100 rounded-tr-none'
                                : 'bg-white dark:bg-gray-700 text-gray-800 dark:text-gray-200 rounded-tl-none border border-gray-100 dark:border-gray-600'
                            }`}
                          >
                            {/* Template label */}
                            {msg.template_name && isOutbound && (
                              <p className="text-[9px] text-green-600 dark:text-green-400 font-medium mb-0.5">
                                {msg.template_name}
                              </p>
                            )}

                            {/* Content */}
                            <p className="whitespace-pre-wrap break-words leading-relaxed">
                              {msg.content || (msg.template_name ? `[Plantilla: ${msg.template_name}]` : '[Sin contenido]')}
                            </p>

                            {/* Footer: time + status */}
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
                          </div>
                        </div>
                      );
                    })}
                  </div>
                ))}
                <div ref={messagesEndRef} />
              </div>

              {/* Bottom bar */}
              <div className="border-t border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 px-3 py-2">
                {!aiPaused ? (
                  /* AI is active — show blocked state */
                  <div className="flex flex-col items-center justify-center py-2 gap-1">
                    <p className="text-[11px] text-gray-500 dark:text-gray-400">
                      La IA está respondiendo. Pausa la IA para escribir manualmente.
                    </p>
                    <button
                      onClick={handleToggleAI}
                      disabled={aiToggling}
                      className="px-3 py-1 text-[11px] bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-300 rounded-full hover:bg-red-200 dark:hover:bg-red-900/60 transition-colors disabled:opacity-50"
                    >
                      {aiToggling ? 'Pausando...' : '🔴 Pausar IA y escribir'}
                    </button>
                  </div>
                ) : (
                  /* AI paused — composer active */
                  <>
                    {!isWindowActive && (
                      <p className="text-[10px] text-amber-500 dark:text-amber-400 mb-1 text-center">
                        ⚠ El cliente no ha respondido en 24h — Meta puede rechazar el mensaje
                      </p>
                    )}
                    {sendError && (
                      <p className="text-[10px] text-red-500 mb-1 text-center">{sendError}</p>
                    )}
                    <div className="flex items-end gap-2">
                      <textarea
                        ref={textareaRef}
                        value={replyText}
                        onChange={e => setReplyText(e.target.value)}
                        onKeyDown={handleKeyDown}
                        placeholder="Escribe un mensaje... (Enter para enviar, Shift+Enter para nueva línea)"
                        rows={1}
                        disabled={sending}
                        className="flex-1 resize-none bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-gray-100 text-sm rounded-xl px-3 py-2 outline-none focus:ring-2 focus:ring-green-500 placeholder-gray-400 disabled:opacity-50 max-h-32 overflow-y-auto"
                        style={{ fieldSizing: 'content' } as React.CSSProperties}
                      />
                      <button
                        onClick={handleSend}
                        disabled={sending || !replyText.trim()}
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
                        Ventana activa — el cliente respondio recientemente
                      </p>
                    )}
                  </>
                )}
              </div>
            </>
          ) : (
            <div className="flex-1 flex items-center justify-center text-gray-400">
              <p className="text-sm">Error al cargar la conversacion</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
