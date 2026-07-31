'use client';

import Link from 'next/link';
import { useEffect, useRef, useState } from 'react';
import {
    getTiendaSessionAction, updateTiendaProfileAction, uploadTiendaAvatarAction,
    listTiendaAddressesAction, addTiendaAddressAction, setTiendaPrimaryAddressAction, deleteTiendaAddressAction,
} from '../../infra/actions';
import { TiendaSession, TiendaAddress } from '../../domain/types';

const inputCls = 'w-full px-4 py-2.5 border border-gray-300 rounded-lg bg-white text-gray-900 text-sm focus:ring-2 focus:border-transparent';

export function PerfilForm({ slug }: { slug: string }) {
    const [session, setSession] = useState<TiendaSession | null>(null);
    const [state, setState] = useState<'loading' | 'anonymous' | 'ready'>('loading');
    const [form, setForm] = useState({ name: '', phone: '', dni: '' });
    const [saving, setSaving] = useState(false);
    const [message, setMessage] = useState('');

    const [avatar, setAvatar] = useState('');
    const [uploadingAvatar, setUploadingAvatar] = useState(false);
    const fileRef = useRef<HTMLInputElement>(null);

    const [addresses, setAddresses] = useState<TiendaAddress[]>([]);
    const [newAddress, setNewAddress] = useState({ street: '', city: '', state: '' });
    const [addingAddress, setAddingAddress] = useState(false);

    const loadAddresses = async () => {
        const list = await listTiendaAddressesAction(slug);
        setAddresses(list || []);
    };

    useEffect(() => {
        let cancelled = false;
        (async () => {
            const s = await getTiendaSessionAction(slug);
            if (cancelled) return;
            if (!s) {
                setState('anonymous');
                return;
            }
            setSession(s);
            setForm({ name: s.name, phone: s.phone, dni: s.dni });
            setAvatar(s.avatar_url || '');
            await loadAddresses();
            if (!cancelled) setState('ready');
        })();
        return () => { cancelled = true; };
    }, [slug]);

    const flash = (msg: string) => {
        setMessage(msg);
        setTimeout(() => setMessage(''), 3000);
    };

    const handleSaveProfile = async (e: React.FormEvent) => {
        e.preventDefault();
        setSaving(true);
        const result = await updateTiendaProfileAction(slug, form);
        setSaving(false);
        if (result?.success) {
            flash('Perfil actualizado');
            window.dispatchEvent(new Event('tienda-session-changed'));
        } else {
            flash('No se pudo actualizar el perfil');
        }
    };

    const handleAvatar = async (file: File | undefined) => {
        if (!file) return;
        setUploadingAvatar(true);
        const formData = new FormData();
        formData.append('image', file);
        const result = await uploadTiendaAvatarAction(slug, formData);
        setUploadingAvatar(false);
        if (result?.success && result.avatar_url) {
            setAvatar(result.avatar_url);
            flash('Foto actualizada');
            window.dispatchEvent(new Event('tienda-session-changed'));
        } else {
            flash('No se pudo subir la foto');
        }
        if (fileRef.current) fileRef.current.value = '';
    };

    const handleAddAddress = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!newAddress.street.trim()) return;
        setAddingAddress(true);
        const result = await addTiendaAddressAction(slug, {
            street: newAddress.street,
            city: newAddress.city,
            state: newAddress.state,
        });
        setAddingAddress(false);
        if (result?.success) {
            setNewAddress({ street: '', city: '', state: '' });
            await loadAddresses();
            flash('Direccion agregada');
        } else {
            flash('No se pudo agregar la direccion');
        }
    };

    const handleSetPrimary = async (id: number) => {
        const result = await setTiendaPrimaryAddressAction(slug, id);
        if (result?.success) {
            await loadAddresses();
            flash('Direccion principal actualizada');
        }
    };

    const handleDeleteAddress = async (id: number) => {
        const result = await deleteTiendaAddressAction(slug, id);
        if (result?.success) {
            await loadAddresses();
            flash('Direccion eliminada');
        }
    };

    if (state === 'loading') {
        return <p className="text-center text-gray-400 py-16">Cargando tu perfil...</p>;
    }

    if (state === 'anonymous' || !session) {
        return (
            <div className="text-center py-16 space-y-4">
                <p className="text-gray-500">Inicia sesion para ver tu perfil.</p>
                <Link
                    href={`/tienda/${slug}/cuenta`}
                    className="inline-block px-6 py-3 rounded-lg text-white font-medium"
                    style={{ backgroundColor: 'var(--brand-secondary)' }}
                >
                    Iniciar sesion
                </Link>
            </div>
        );
    }

    return (
        <div className="max-w-xl mx-auto space-y-6">
            {message && (
                <p className="text-sm text-center text-green-700 bg-green-50 border border-green-200 rounded-lg px-3 py-2">{message}</p>
            )}

            <div className="bg-white rounded-2xl border border-gray-100 shadow-sm p-6 space-y-4">
                <div className="flex items-center gap-4">
                    <div className="relative">
                        {avatar ? (
                            <img src={avatar} alt={session.name} className="w-20 h-20 rounded-full object-cover border border-gray-200" />
                        ) : (
                            <div className="w-20 h-20 rounded-full flex items-center justify-center text-white text-2xl font-bold" style={{ backgroundColor: 'var(--brand-secondary)' }}>
                                {session.name.charAt(0).toUpperCase()}
                            </div>
                        )}
                        <button
                            type="button"
                            onClick={() => fileRef.current?.click()}
                            disabled={uploadingAvatar}
                            className="absolute -bottom-1 -right-1 p-1.5 rounded-full text-white shadow"
                            style={{ backgroundColor: 'var(--brand-secondary)' }}
                            title="Cambiar foto"
                        >
                            <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
                                <circle cx="12" cy="13" r="3" strokeWidth={2} stroke="currentColor" fill="none" />
                            </svg>
                        </button>
                        <input ref={fileRef} type="file" accept="image/*" className="hidden" onChange={(e) => handleAvatar(e.target.files?.[0])} />
                    </div>
                    <div>
                        <p className="font-semibold text-gray-900">{session.name}</p>
                        <p className="text-sm text-gray-400">{session.email}</p>
                        {uploadingAvatar && <p className="text-xs text-gray-400">Subiendo foto...</p>}
                    </div>
                </div>

                <form onSubmit={handleSaveProfile} className="space-y-3">
                    <div>
                        <label className="block text-xs font-medium text-gray-500 mb-1">Nombre</label>
                        <input type="text" value={form.name} onChange={(e) => setForm(f => ({ ...f, name: e.target.value }))} className={inputCls} />
                    </div>
                    <div className="grid grid-cols-2 gap-3">
                        <div>
                            <label className="block text-xs font-medium text-gray-500 mb-1">Telefono</label>
                            <input type="tel" value={form.phone} onChange={(e) => setForm(f => ({ ...f, phone: e.target.value }))} className={inputCls} />
                        </div>
                        <div>
                            <label className="block text-xs font-medium text-gray-500 mb-1">Documento</label>
                            <input type="text" value={form.dni} onChange={(e) => setForm(f => ({ ...f, dni: e.target.value }))} className={inputCls} />
                        </div>
                    </div>
                    <button
                        type="submit"
                        disabled={saving}
                        className="w-full py-2.5 rounded-lg text-white font-semibold disabled:opacity-50"
                        style={{ backgroundColor: 'var(--brand-secondary)' }}
                    >
                        {saving ? 'Guardando...' : 'Guardar cambios'}
                    </button>
                </form>
            </div>

            <div className="bg-white rounded-2xl border border-gray-100 shadow-sm p-6 space-y-4">
                <h2 className="font-semibold text-gray-900">Direcciones de entrega</h2>

                {addresses.length === 0 && (
                    <p className="text-sm text-gray-400">Aun no tienes direcciones guardadas.</p>
                )}

                {addresses.map(a => (
                    <div key={a.id} className="flex items-start justify-between gap-3 border border-gray-100 rounded-lg p-3">
                        <label className="flex items-start gap-3 cursor-pointer flex-1">
                            <input
                                type="radio"
                                name="primary-address"
                                checked={a.is_primary}
                                onChange={() => handleSetPrimary(a.id)}
                                className="mt-1"
                            />
                            <span>
                                <span className="block text-sm text-gray-800">{a.street}</span>
                                <span className="block text-xs text-gray-400">
                                    {[a.city, a.state].filter(Boolean).join(', ')}
                                    {a.is_primary && <span className="ml-2 text-green-600 font-medium">Principal</span>}
                                </span>
                            </span>
                        </label>
                        <button
                            type="button"
                            onClick={() => handleDeleteAddress(a.id)}
                            className="text-gray-300 hover:text-red-500 text-xs shrink-0"
                        >
                            Eliminar
                        </button>
                    </div>
                ))}

                <form onSubmit={handleAddAddress} className="space-y-2 border-t border-gray-100 pt-4">
                    <input
                        type="text"
                        placeholder="Nueva direccion (calle, numero)"
                        value={newAddress.street}
                        onChange={(e) => setNewAddress(a => ({ ...a, street: e.target.value }))}
                        className={inputCls}
                    />
                    <div className="grid grid-cols-2 gap-2">
                        <input
                            type="text"
                            placeholder="Ciudad"
                            value={newAddress.city}
                            onChange={(e) => setNewAddress(a => ({ ...a, city: e.target.value }))}
                            className={inputCls}
                        />
                        <input
                            type="text"
                            placeholder="Departamento"
                            value={newAddress.state}
                            onChange={(e) => setNewAddress(a => ({ ...a, state: e.target.value }))}
                            className={inputCls}
                        />
                    </div>
                    <button
                        type="submit"
                        disabled={addingAddress || !newAddress.street.trim()}
                        className="w-full py-2 rounded-lg text-sm font-medium border border-gray-200 text-gray-700 hover:bg-gray-50 disabled:opacity-40"
                    >
                        {addingAddress ? 'Agregando...' : 'Agregar direccion'}
                    </button>
                </form>
            </div>
        </div>
    );
}
