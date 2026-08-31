'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { tiendaLoginAction, tiendaRegisterAction } from '../../infra/actions';

interface CuentaFormProps {
    slug: string;
}

export function CuentaForm({ slug }: CuentaFormProps) {
    const router = useRouter();
    const [mode, setMode] = useState<'login' | 'register'>('login');
    const [form, setForm] = useState({ name: '', email: '', password: '', phone: '', dni: '' });
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');

    const set = (k: string, v: string) => setForm(prev => ({ ...prev, [k]: v }));

    const inputCls = 'w-full px-4 py-2.5 border border-gray-300 rounded-lg bg-white text-gray-900 focus:ring-2 focus:border-transparent';

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError('');

        const result = mode === 'login'
            ? await tiendaLoginAction(form.email, form.password)
            : await tiendaRegisterAction(slug, {
                name: form.name,
                email: form.email,
                password: form.password,
                phone: form.phone || undefined,
                dni: form.dni || undefined,
            });

        setLoading(false);
        if (result.success) {
            window.dispatchEvent(new Event('tienda-session-changed'));
            router.push(`/tienda/${slug}`);
            router.refresh();
        } else {
            setError(result.error || 'Ocurrió un error, intenta de nuevo');
        }
    };

    return (
        <div className="max-w-md mx-auto bg-white rounded-2xl shadow-sm border border-gray-100 p-8">
            <div className="flex rounded-lg bg-gray-100 p-1 mb-6">
                <button
                    type="button"
                    onClick={() => { setMode('login'); setError(''); }}
                    className={`flex-1 py-2 text-sm font-medium rounded-md transition-colors ${mode === 'login' ? 'bg-white shadow text-gray-900' : 'text-gray-500'}`}
                >
                    Iniciar sesion
                </button>
                <button
                    type="button"
                    onClick={() => { setMode('register'); setError(''); }}
                    className={`flex-1 py-2 text-sm font-medium rounded-md transition-colors ${mode === 'register' ? 'bg-white shadow text-gray-900' : 'text-gray-500'}`}
                >
                    Crear cuenta
                </button>
            </div>

            <form onSubmit={handleSubmit} className="space-y-4">
                {mode === 'register' && (
                    <input
                        type="text"
                        required
                        placeholder="Nombre completo"
                        value={form.name}
                        onChange={(e) => set('name', e.target.value)}
                        className={inputCls}
                    />
                )}
                <input
                    type="email"
                    required
                    placeholder="Correo electrónico"
                    value={form.email}
                    onChange={(e) => set('email', e.target.value)}
                    className={inputCls}
                />
                <input
                    type="password"
                    required
                    minLength={6}
                    placeholder="Contraseña"
                    value={form.password}
                    onChange={(e) => set('password', e.target.value)}
                    className={inputCls}
                />
                {mode === 'register' && (
                    <div className="grid grid-cols-2 gap-3">
                        <input
                            type="tel"
                            placeholder="Teléfono"
                            value={form.phone}
                            onChange={(e) => set('phone', e.target.value)}
                            className={inputCls}
                        />
                        <input
                            type="text"
                            placeholder="Documento (opcional)"
                            value={form.dni}
                            onChange={(e) => set('dni', e.target.value)}
                            className={inputCls}
                        />
                    </div>
                )}

                {error && (
                    <div className="flex items-start gap-2 text-sm text-red-700 bg-red-50 border border-red-200 rounded-lg px-3 py-2.5">
                        <svg className="w-4 h-4 mt-0.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        <span>{error}</span>
                    </div>
                )}

                <button
                    type="submit"
                    disabled={loading}
                    className="w-full py-3 rounded-lg text-white font-semibold transition-transform hover:scale-[1.01] disabled:opacity-50"
                    style={{ backgroundColor: 'var(--brand-secondary)' }}
                >
                    {loading ? 'Un momento...' : mode === 'login' ? 'Entrar' : 'Crear mi cuenta'}
                </button>
            </form>

            <p className="text-xs text-gray-400 mt-4 text-center">
                Con tu cuenta puedes comprar mas rapido y ver el estado de tus pedidos.
            </p>
            {mode === 'register' && (
                <p className="text-xs text-gray-400 mt-2 text-center">
                    Si ya tienes cuenta en otra tienda de Probability, usa el mismo correo y contrasena para vincularte.
                </p>
            )}
        </div>
    );
}
