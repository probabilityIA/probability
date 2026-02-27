'use client';

import { useState, useRef, useEffect } from 'react';

// ─── FAQ Data ────────────────────────────────────────────────────────────────
const FAQ_ITEMS = [
    {
        id: 'shopify',
        question: '¿Cómo integro con Shopify?',
        answer: `Para integrar tu tienda Shopify con ProbabilityIA:\n\n1. Ve a **Integraciones** en el menú lateral.\n2. Selecciona la categoría **E-Commerce** y haz clic en **Shopify**.\n3. Ingresa el dominio de tu tienda (ej: *mitienda.myshopify.com*).\n4. Serás redirigido a Shopify para autorizar el acceso.\n5. Una vez autorizado, ¡listo! Tus pedidos se sincronizarán automáticamente.\n\nSi tienes una Custom App de Shopify, puedes usar la opción de conexión personalizada.`,
    },
    {
        id: 'productos',
        question: '¿Cómo creo mis productos?',
        answer: `Para gestionar tus productos:\n\n1. Dirígete a **Productos** en el menú lateral.\n2. Haz clic en **Nuevo Producto** (botón superior derecho).\n3. Completa el nombre, descripción, precio y SKU.\n4. Agrega imágenes y variantes si las necesitas.\n5. Guarda con **Crear Producto**.\n\nTambién puedes importar productos masivamente desde tu integración de e-commerce con solo sincronizar.`,
    },
    {
        id: 'interrapidisimo',
        question: '¿Cómo solicito envíos con Interrapidísimo?',
        answer: `Para generar una guía de Interrapidísimo:\n\n1. Ve a **Envíos** en el menú lateral.\n2. Selecciona el pedido que deseas despachar.\n3. Haz clic en **Generar Guía**.\n4. En el paso de cotización verás las tarifas disponibles, incluyendo **Interrapidísimo**.\n5. Elige la tarifa y completa los datos del destinatario.\n6. Haz clic en **Generar Guía** para confirmar.\n\nLa guía PDF se descargará automáticamente.`,
    },
    {
        id: 'integraciones-generales',
        question: '¿Qué integraciones están disponibles?',
        answer: `ProbabilityIA cuenta con integraciones en varias categorías:\n\n📦 **E-Commerce:** Shopify, MercadoLibre, WooCommerce, VTEX, Tiendanube, Magento, Amazon, Falabella, Éxito.\n\n🚚 **Transporte:** EnvioClick (múltiples transportadoras), EnvíAme, MiPaquete.\n\n📄 **Facturación Electrónica:** Alegra, Siigo, Factus, Helisa, World Office, Softpymes.\n\n💬 **Mensajería:** WhatsApp Business.\n\n💳 **Pagos:** Nequi, Wompi, PayU, ePayco, Stripe, Bold, MercadoPago.`,
    },
    {
        id: 'pedidos',
        question: '¿Cómo gestiono mis pedidos?',
        answer: `Para ver y gestionar tus pedidos:\n\n1. Haz clic en **Pedidos** en el menú lateral.\n2. Verás todos los pedidos sincronizados desde tus integraciones.\n3. Puedes **filtrar** por estado, fecha o canal de venta.\n4. Haz clic en cualquier pedido para ver su detalle completo.\n5. Desde el detalle puedes: generar guía de envío, emitir factura electrónica o actualizar el estado.\n\nLos nuevos pedidos llegan en tiempo real vía webhooks.`,
    },
];

// ─── Suggestion Form ──────────────────────────────────────────────────────────
function SuggestionForm({ onBack }: { onBack: () => void }) {
    const [text, setText] = useState('');
    const [sent, setSent] = useState(false);

    const handleSend = () => {
        if (!text.trim()) return;
        setSent(true);
    };

    if (sent) {
        return (
            <div className="flex flex-col items-center justify-center gap-3 py-8 px-4 text-center">
                <div style={{ fontSize: 40 }}>🎉</div>
                <p style={{ color: '#7c3aed', fontWeight: 700, fontSize: 15 }}>¡Gracias por tu mensaje!</p>
                <p style={{ color: '#6b7280', fontSize: 13 }}>Tu sugerencia ha sido recibida. El equipo de ProbabilityIA la revisará pronto.</p>
                <button
                    onClick={onBack}
                    style={{ marginTop: 8, padding: '8px 20px', background: '#7c3aed', color: '#fff', border: 'none', borderRadius: 8, cursor: 'pointer', fontSize: 13, fontWeight: 600 }}
                >
                    Volver al inicio
                </button>
            </div>
        );
    }

    return (
        <div className="flex flex-col gap-3 px-1">
            <button onClick={onBack} style={{ alignSelf: 'flex-start', background: 'none', border: 'none', cursor: 'pointer', color: '#7c3aed', fontSize: 13, fontWeight: 600, padding: 0, display: 'flex', alignItems: 'center', gap: 4 }}>
                ← Volver
            </button>
            <p style={{ color: '#374151', fontSize: 13, fontWeight: 500 }}>Escribe tu sugerencia o petición y el equipo de ProbabilityIA la revisará:</p>
            <textarea
                value={text}
                onChange={e => setText(e.target.value)}
                placeholder="Ej: Me gustaría que hubiera una funcionalidad para..."
                rows={4}
                style={{ width: '100%', padding: '10px 12px', borderRadius: 10, border: '1.5px solid #e5e7eb', fontSize: 13, color: '#374151', resize: 'none', outline: 'none', fontFamily: 'inherit', background: '#f9fafb' }}
                onFocus={e => e.currentTarget.style.borderColor = '#a78bfa'}
                onBlur={e => e.currentTarget.style.borderColor = '#e5e7eb'}
            />
            <button
                onClick={handleSend}
                disabled={!text.trim()}
                style={{ padding: '10px', background: text.trim() ? 'linear-gradient(135deg, #7c3aed, #a855f7)' : '#e5e7eb', color: text.trim() ? '#fff' : '#9ca3af', border: 'none', borderRadius: 10, cursor: text.trim() ? 'pointer' : 'not-allowed', fontSize: 13, fontWeight: 700, transition: 'all 0.2s' }}
            >
                Enviar sugerencia ✉️
            </button>
        </div>
    );
}

// ─── Message Bubble ───────────────────────────────────────────────────────────
function MessageBubble({ text, isLina }: { text: string; isLina: boolean }) {
    // Render basic markdown: **bold**, *italic*, newlines
    const renderMarkdown = (raw: string) => {
        const lines = raw.split('\n');
        return lines.map((line, i) => {
            const parts = line.split(/(\*\*[^*]+\*\*|\*[^*]+\*)/g).map((part, j) => {
                if (part.startsWith('**') && part.endsWith('**')) {
                    return <strong key={j}>{part.slice(2, -2)}</strong>;
                }
                if (part.startsWith('*') && part.endsWith('*')) {
                    return <em key={j}>{part.slice(1, -1)}</em>;
                }
                return part;
            });
            return <span key={i}>{parts}{i < lines.length - 1 && <br />}</span>;
        });
    };

    return (
        <div style={{ display: 'flex', justifyContent: isLina ? 'flex-start' : 'flex-end', marginBottom: 8 }}>
            {isLina && (
                <div style={{ width: 28, height: 28, borderRadius: '50%', background: 'linear-gradient(135deg, #7c3aed, #a855f7)', display: 'flex', alignItems: 'center', justifyContent: 'center', marginRight: 8, flexShrink: 0, fontSize: 14 }}>
                    ✨
                </div>
            )}
            <div style={{
                maxWidth: '82%',
                padding: '10px 14px',
                borderRadius: isLina ? '4px 16px 16px 16px' : '16px 4px 16px 16px',
                background: isLina ? '#fff' : 'linear-gradient(135deg, #7c3aed, #a855f7)',
                color: isLina ? '#374151' : '#fff',
                fontSize: 13,
                lineHeight: 1.55,
                boxShadow: isLina ? '0 1px 4px rgba(0,0,0,0.08)' : '0 2px 8px rgba(124,58,237,0.3)',
                border: isLina ? '1px solid #f3f4f6' : 'none',
            }}>
                {renderMarkdown(text)}
            </div>
        </div>
    );
}

// ─── Main Component ───────────────────────────────────────────────────────────
interface LinaChatbotProps {
    userScope?: string;
    isSuperAdmin?: boolean;
}

type View = 'home' | 'faq_answer' | 'suggestion';
interface Message { text: string; isLina: boolean; }

export function LinaChatbot({ userScope, isSuperAdmin }: LinaChatbotProps) {
    // Solo mostrar para business o super_admin
    const shouldShow = userScope === 'business' || isSuperAdmin === true;
    const [open, setOpen] = useState(false);
    const [view, setView] = useState<View>('home');
    const [messages, setMessages] = useState<Message[]>([
        { text: '¡Hola! Soy **Lina**, tu asistente de **ProbabilityIA** 👋\n\nEstoy aquí para ayudarte. ¿Sobre qué tema tienes dudas?', isLina: true }
    ]);
    const [selectedFaq, setSelectedFaq] = useState<typeof FAQ_ITEMS[0] | null>(null);
    const messagesEndRef = useRef<HTMLDivElement>(null);
    const [pulse, setPulse] = useState(false);

    useEffect(() => {
        if (open) {
            setTimeout(() => messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' }), 50);
        }
    }, [open, messages, view]);

    // Pulse animation on first render
    useEffect(() => {
        if (!open) {
            const t = setInterval(() => setPulse(p => !p), 3000);
            return () => clearInterval(t);
        }
    }, [open]);

    const handleFaqClick = (faq: typeof FAQ_ITEMS[0]) => {
        setSelectedFaq(faq);
        setMessages(prev => [
            ...prev,
            { text: faq.question, isLina: false },
            { text: faq.answer, isLina: true },
        ]);
        setView('faq_answer');
    };

    const handleBack = () => {
        setView('home');
    };

    if (!shouldShow) return null;

    return (
        <>
            {/* ── Chat Window ── */}
            {open && (
                <div
                    style={{
                        position: 'fixed',
                        bottom: 90,
                        right: 24,
                        width: 360,
                        maxHeight: '70vh',
                        background: '#f9fafb',
                        borderRadius: 20,
                        boxShadow: '0 20px 60px rgba(0,0,0,0.18), 0 4px 20px rgba(124,58,237,0.15)',
                        display: 'flex',
                        flexDirection: 'column',
                        overflow: 'hidden',
                        zIndex: 9999,
                        animation: 'linaChatIn 0.25s cubic-bezier(.34,1.56,.64,1)',
                    }}
                >
                    {/* Header */}
                    <div style={{ background: 'linear-gradient(135deg, #5b21b6, #7c3aed, #a855f7)', padding: '16px 18px', display: 'flex', alignItems: 'center', gap: 12 }}>
                        <div style={{ width: 40, height: 40, borderRadius: '50%', background: 'rgba(255,255,255,0.2)', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 20, backdropFilter: 'blur(4px)' }}>
                            ✨
                        </div>
                        <div style={{ flex: 1 }}>
                            <div style={{ color: '#fff', fontWeight: 700, fontSize: 15, letterSpacing: 0.2 }}>Lina</div>
                            <div style={{ color: 'rgba(255,255,255,0.75)', fontSize: 12, display: 'flex', alignItems: 'center', gap: 5 }}>
                                <span style={{ width: 7, height: 7, borderRadius: '50%', background: '#4ade80', display: 'inline-block' }} />
                                Asistente de ProbabilityIA
                            </div>
                        </div>
                        <button
                            onClick={() => setOpen(false)}
                            style={{ background: 'rgba(255,255,255,0.15)', border: 'none', borderRadius: '50%', width: 32, height: 32, cursor: 'pointer', color: '#fff', fontSize: 16, display: 'flex', alignItems: 'center', justifyContent: 'center' }}
                            aria-label="Cerrar"
                        >
                            ✕
                        </button>
                    </div>

                    {/* Body */}
                    <div style={{ flex: 1, overflowY: 'auto', padding: '16px 14px 8px' }}>
                        {/* Messages */}
                        {messages.map((msg, i) => (
                            <MessageBubble key={i} text={msg.text} isLina={msg.isLina} />
                        ))}
                        <div ref={messagesEndRef} />
                    </div>

                    {/* Quick Actions */}
                    <div style={{ padding: '10px 14px 16px', background: '#f9fafb', borderTop: '1px solid #f3f4f6' }}>
                        {view === 'suggestion' ? (
                            <SuggestionForm onBack={() => { handleBack(); }} />
                        ) : view === 'faq_answer' ? (
                            <div style={{ display: 'flex', flexDirection: 'column', gap: 7 }}>
                                <p style={{ color: '#6b7280', fontSize: 12, fontWeight: 500, textAlign: 'center', marginBottom: 2 }}>¿Tienes otra pregunta?</p>
                                {FAQ_ITEMS.filter(f => f.id !== selectedFaq?.id).slice(0, 2).map(faq => (
                                    <button
                                        key={faq.id}
                                        onClick={() => handleFaqClick(faq)}
                                        style={{ padding: '9px 14px', background: '#fff', border: '1.5px solid #ede9fe', borderRadius: 10, cursor: 'pointer', fontSize: 12, color: '#5b21b6', fontWeight: 600, textAlign: 'left', transition: 'all 0.15s' }}
                                        onMouseEnter={e => { e.currentTarget.style.background = '#ede9fe'; }}
                                        onMouseLeave={e => { e.currentTarget.style.background = '#fff'; }}
                                    >
                                        {faq.question}
                                    </button>
                                ))}
                                <button
                                    onClick={handleBack}
                                    style={{ padding: '9px 14px', background: '#fff', border: '1.5px solid #e5e7eb', borderRadius: 10, cursor: 'pointer', fontSize: 12, color: '#6b7280', fontWeight: 600, textAlign: 'left', transition: 'all 0.15s' }}
                                    onMouseEnter={e => { e.currentTarget.style.background = '#f3f4f6'; }}
                                    onMouseLeave={e => { e.currentTarget.style.background = '#fff'; }}
                                >
                                    ← Ver todas las preguntas
                                </button>
                                <button
                                    onClick={() => setView('suggestion')}
                                    style={{ padding: '9px 14px', background: 'linear-gradient(135deg, #7c3aed, #a855f7)', border: 'none', borderRadius: 10, cursor: 'pointer', fontSize: 12, color: '#fff', fontWeight: 600, textAlign: 'center', transition: 'all 0.15s' }}
                                >
                                    ✉️ Enviar sugerencia o petición
                                </button>
                            </div>
                        ) : (
                            <div style={{ display: 'flex', flexDirection: 'column', gap: 7 }}>
                                <p style={{ color: '#6b7280', fontSize: 12, fontWeight: 500, textAlign: 'center', marginBottom: 2 }}>Selecciona una pregunta frecuente:</p>
                                {FAQ_ITEMS.map(faq => (
                                    <button
                                        key={faq.id}
                                        onClick={() => handleFaqClick(faq)}
                                        style={{ padding: '9px 14px', background: '#fff', border: '1.5px solid #ede9fe', borderRadius: 10, cursor: 'pointer', fontSize: 12, color: '#5b21b6', fontWeight: 600, textAlign: 'left', transition: 'all 0.15s' }}
                                        onMouseEnter={e => { e.currentTarget.style.background = '#ede9fe'; }}
                                        onMouseLeave={e => { e.currentTarget.style.background = '#fff'; }}
                                    >
                                        {faq.question}
                                    </button>
                                ))}
                                <div style={{ height: 1, background: '#f3f4f6', margin: '4px 0' }} />
                                <button
                                    onClick={() => setView('suggestion')}
                                    style={{ padding: '9px 14px', background: 'linear-gradient(135deg, #7c3aed, #a855f7)', border: 'none', borderRadius: 10, cursor: 'pointer', fontSize: 12, color: '#fff', fontWeight: 700, textAlign: 'center', transition: 'all 0.15s' }}
                                    onMouseEnter={e => { e.currentTarget.style.opacity = '0.9'; }}
                                    onMouseLeave={e => { e.currentTarget.style.opacity = '1'; }}
                                >
                                    ✉️ Enviar sugerencia o petición
                                </button>
                            </div>
                        )}
                    </div>
                </div>
            )}

            {/* ── Floating Button ── */}
            <button
                onClick={() => setOpen(o => !o)}
                aria-label="Abrir asistente Lina"
                style={{
                    position: 'fixed',
                    bottom: 24,
                    right: 24,
                    width: 58,
                    height: 58,
                    borderRadius: '50%',
                    background: 'linear-gradient(135deg, #5b21b6, #7c3aed, #a855f7)',
                    border: 'none',
                    cursor: 'pointer',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    boxShadow: '0 6px 24px rgba(124,58,237,0.5)',
                    zIndex: 9999,
                    transition: 'transform 0.2s, box-shadow 0.2s',
                    transform: open ? 'rotate(10deg) scale(1.05)' : pulse ? 'scale(1.08)' : 'scale(1)',
                    fontSize: open ? 22 : 26,
                }}
                onMouseEnter={e => { e.currentTarget.style.transform = 'scale(1.12)'; e.currentTarget.style.boxShadow = '0 8px 32px rgba(124,58,237,0.65)'; }}
                onMouseLeave={e => { e.currentTarget.style.transform = 'scale(1)'; e.currentTarget.style.boxShadow = '0 6px 24px rgba(124,58,237,0.5)'; }}
            >
                {open ? '✕' : '✨'}
                {/* Ping animation when closed */}
                {!open && (
                    <span style={{
                        position: 'absolute',
                        top: 0, right: 0,
                        width: 14, height: 14,
                        borderRadius: '50%',
                        background: '#4ade80',
                        border: '2px solid #fff',
                        animation: 'linaPing 2s ease-in-out infinite',
                    }} />
                )}
            </button>

            {/* ── Tooltip label ── */}
            {!open && (
                <div style={{
                    position: 'fixed',
                    bottom: 36,
                    right: 90,
                    background: '#fff',
                    padding: '6px 12px',
                    borderRadius: 20,
                    boxShadow: '0 4px 16px rgba(0,0,0,0.12)',
                    fontSize: 13,
                    fontWeight: 600,
                    color: '#5b21b6',
                    zIndex: 9998,
                    pointerEvents: 'none',
                    animation: 'linaFadeLabel 3s ease-in-out infinite',
                }}>
                    ¿Necesitas ayuda? 👋
                </div>
            )}

            {/* Animations */}
            <style>{`
                @keyframes linaChatIn {
                    from { opacity: 0; transform: translateY(20px) scale(0.95); }
                    to { opacity: 1; transform: translateY(0) scale(1); }
                }
                @keyframes linaPing {
                    0%, 100% { transform: scale(1); opacity: 1; }
                    50% { transform: scale(1.4); opacity: 0.5; }
                }
                @keyframes linaFadeLabel {
                    0%, 100% { opacity: 0; transform: translateX(8px); }
                    20%, 80% { opacity: 1; transform: translateX(0); }
                }
            `}</style>
        </>
    );
}
