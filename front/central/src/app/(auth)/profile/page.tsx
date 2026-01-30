'use client';

import React, { useEffect, useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { TokenStorage } from '@/shared/config';
import { AvatarUpload, Button, Spinner, Alert } from '@/shared/ui';
import { getUserByIdAction, updateUserAction } from '@/services/auth/users/infra/actions';
import { getProductsAction } from '@/services/modules/products/infra/actions';
import { ChangePasswordForm } from '@/services/auth/login/ui';
import { Modal } from '@/shared/ui/modal';

// Tipos para el estado local
interface UserProfile {
    id: number;
    name: string;
    email: string;
    phone?: string;
    role?: string;
    avatarUrl?: string;
    isActive: boolean;
    businessRoleAssignments?: any[];
}

interface Product {
    id: string | number;
    name: string;
    description?: string;
    // Otros campos de producto
}

export default function ProfilePage() {
    const router = useRouter();
    const [user, setUser] = useState<UserProfile | null>(null);
    const [products, setProducts] = useState<Product[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [showChangePassword, setShowChangePassword] = useState(false);
    const [passwordSuccess, setPasswordSuccess] = useState(false);

    useEffect(() => {
        loadData();
    }, []);

    const loadData = async () => {
        try {
            setLoading(true);
            const tokenUser = TokenStorage.getUser();

            if (!tokenUser || !tokenUser.userId) {
                router.push('/login');
                return;
            }

            const userId = parseInt(tokenUser.userId);

            // 1. Obtener datos frescos del usuario
            const userResponse = await getUserByIdAction(userId);
            if (userResponse.success && userResponse.data) {
                const u = userResponse.data;

                // Normalizar datos
                const normalizedUser: UserProfile = {
                    id: u.id,
                    name: u.name,
                    email: u.email,
                    phone: u.phone,
                    role: tokenUser.role,
                    avatarUrl: (u as any).avatar_url || (u as any).avatarUrl,
                    isActive: (u as any).is_active || (u as any).isActive,
                    businessRoleAssignments: (u as any).business_role_assignments || (u as any).businessRoleAssignments
                };
                setUser(normalizedUser);
            } else {
                setError("No se pudo cargar la información del usuario.");
            }

            // 2. Obtener productos
            const productsResponse = await getProductsAction({ page: 1, page_size: 10 });
            if (productsResponse.success) {
                setProducts(productsResponse.data);
            }

        } catch (err: any) {
            console.error("Error al cargar perfil:", err);
            setError(err.message || "Error desconocido al cargar perfil");
        } finally {
            setLoading(false);
        }
    };

    const handleAvatarUpdate = async (file: File | null) => {
        if (!user) return;
        try {
            const userId = user.id;
            const updateData: any = {};

            if (file) {
                updateData.avatarFile = file;
            } else {
                updateData.remove_avatar = true;
            }

            const response = await updateUserAction(userId, updateData);

            if (response.success) {
                loadData();
                window.location.reload();
            } else {
                setError("Error al actualizar la imagen: " + (response.message || "Error desconocido"));
            }
        } catch (error: any) {
            console.error("Error updating avatar", error);
            setError("Excepción al actualizar imagen");
        }
    };

    const handlePasswordSuccess = () => {
        setPasswordSuccess(true);
        setTimeout(() => {
            setShowChangePassword(false);
            setPasswordSuccess(false);
        }, 2000);
    };

    if (loading) return <div className="flex w-full h-screen items-center justify-center"><Spinner size="lg" /></div>;
    if (!user) return <div className="p-8 text-center">Usuario no encontrado o sesión expirada.</div>;

    return (
        <div className="max-w-7xl mx-auto p-6 md:p-12 mb-20 animate-in fade-in duration-500">
            {/* Header Moderno Minimalista - Ahora en Card */}
            <div className="bg-white dark:bg-white/5 border border-gray-100 dark:border-white/10 rounded-3xl p-8 shadow-xl shadow-gray-200/50 dark:shadow-none mb-8 flex flex-col md:flex-row items-center md:items-start gap-8 transition-all hover:shadow-2xl hover:shadow-gray-200/50">
                {/* Avatar Section - Sin recortes para permitir botones de edición */}
                <div className="relative group shrink-0">

                    <AvatarUpload
                        currentAvatarUrl={user.avatarUrl}
                        onFileSelect={(file) => handleAvatarUpdate(file)}
                        onRemoveClick={() => handleAvatarUpdate(null)}
                        size="lg" // 128px
                        className="transition-transform duration-300 hover:scale-[1.02]"
                    />

                </div>

                {/* User Info */}
                <div className="flex-1 text-center md:text-left pt-2">
                    <div className="flex flex-col md:flex-row items-center md:items-start md:justify-between gap-4 mb-2">
                        <div className="flex flex-col md:flex-row items-center gap-4">
                            <h1 className="text-4xl font-bold text-gray-900 dark:text-white tracking-tight">
                                {user.name}
                            </h1>
                            <span className="px-3 py-1 rounded-full text-xs font-semibold bg-indigo-50 text-indigo-600 border border-indigo-100 dark:bg-indigo-900/30 dark:text-indigo-300 dark:border-indigo-800 self-center">
                                {user.role || 'Usuario'}
                            </span>
                        </div>
                    </div>
                    <p className="text-gray-500 dark:text-gray-400 text-lg mb-6 max-w-2xl font-light">
                        {user.email}
                    </p>


                </div>
            </div>

            {error && (
                <div className="mb-8 shadow-sm border-l-4 border-red-500">
                    <Alert type="error" onClose={() => setError(null)}>
                        {error}
                    </Alert>
                </div>
            )}

            {/* Bento Grid Layout */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                {/* Columna Principal (2/3) - Productos y Actividad */}
                <div className="md:col-span-2 space-y-6">
                    {/* Tarjeta de Productos - Estilo Bento */}
                    <section className="bg-white dark:bg-white/5 border border-gray-100 dark:border-white/10 rounded-3xl p-8 shadow-xl shadow-gray-200/50 dark:shadow-none transition-all hover:shadow-2xl hover:shadow-gray-200/50">
                        <div className="flex items-center justify-between mb-8">
                            <div className="flex items-center gap-3">
                                <div className="p-2 bg-indigo-50 dark:bg-indigo-900/30 rounded-xl text-indigo-600">
                                    <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                                    </svg>
                                </div>
                                <h3 className="text-xl font-bold text-gray-900 dark:text-white">Mis Productos</h3>
                            </div>
                            <Link
                                href="/products"
                                className="text-sm font-medium text-gray-500 hover:text-indigo-600 transition-colors flex items-center gap-1"
                            >
                                Ver todos
                                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                                </svg>
                            </Link>
                        </div>

                        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                            {products.length > 0 ? (
                                products.slice(0, 4).map((prod: any) => (
                                    <div key={prod.id} className="group relative bg-gray-50 dark:bg-white/5 rounded-2xl p-4 transition-all hover:bg-indigo-50/50 dark:hover:bg-white/10 hover:scale-[1.02] cursor-pointer border border-transparent hover:border-indigo-100">
                                        <div className="flex items-start justify-between mb-2">
                                            <div className="h-10 w-10 rounded-xl bg-white dark:bg-white/10 shadow-sm flex items-center justify-center text-lg">
                                                📦
                                            </div>
                                            <span className="text-[10px] font-mono text-gray-400 bg-white dark:bg-black/20 px-2 py-1 rounded-md">
                                                ID: {prod.id}
                                            </span>
                                        </div>
                                        <h4 className="font-bold text-gray-900 dark:text-white mb-1 truncate">{prod.name}</h4>
                                        <p className="text-xs text-gray-500 line-clamp-2">{prod.description || 'Sin descripción disponible'}</p>
                                    </div>
                                ))
                            ) : (
                                <div className="col-span-2 py-12 text-center text-gray-400 bg-gray-50 rounded-2xl border border-dashed border-gray-200">
                                    No tienes productos activos
                                </div>
                            )}

                            <Link href="/products/new" className="flex flex-col items-center justify-center gap-2 bg-white dark:bg-white/5 border-2 border-dashed border-gray-200 dark:border-white/10 rounded-2xl p-4 transition-all hover:border-indigo-400 hover:bg-indigo-50/30 group cursor-pointer min-h-[120px]">
                                <div className="h-10 w-10 rounded-full bg-indigo-50 dark:bg-indigo-900/30 text-indigo-600 flex items-center justify-center transition-transform group-hover:scale-110">
                                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                                    </svg>
                                </div>
                                <span className="text-sm font-medium text-gray-600 dark:text-gray-300 group-hover:text-indigo-600">Crear Nuevo</span>
                            </Link>
                        </div>
                    </section>
                </div>

                {/* Columna Lateral (1/3) - Detalles y Configuración */}
                <div className="space-y-6">
                    {/* Tarjeta de Contacto */}
                    <div className="bg-white dark:bg-white/5 border border-gray-100 dark:border-white/10 rounded-3xl p-6 shadow-lg shadow-gray-100/50 dark:shadow-none">
                        <h3 className="text-lg font-bold text-gray-900 dark:text-white mb-6">Detalles de Contacto</h3>
                        <div className="space-y-4">
                            <div className="flex items-center gap-4 p-3 hover:bg-gray-50 dark:hover:bg-white/5 rounded-2xl transition-colors group">
                                <div className="h-10 w-10 rounded-full bg-blue-50 text-blue-500 flex items-center justify-center group-hover:bg-blue-100 transition-colors">
                                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                                    </svg>
                                </div>
                                <div className="flex-1 min-w-0">
                                    <p className="text-xs text-gray-500 font-medium uppercase tracking-wide">Email</p>
                                    <p className="text-sm font-semibold text-gray-900 dark:text-white truncate" title={user.email}>
                                        {user.email}
                                    </p>
                                </div>
                            </div>

                            <div className="flex items-center gap-4 p-3 hover:bg-gray-50 dark:hover:bg-white/5 rounded-2xl transition-colors group">
                                <div className="h-10 w-10 rounded-full bg-green-50 text-green-500 flex items-center justify-center group-hover:bg-green-100 transition-colors">
                                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
                                    </svg>
                                </div>
                                <div className="flex-1 min-w-0">
                                    <p className="text-xs text-gray-500 font-medium uppercase tracking-wide">Teléfono</p>
                                    <p className="text-sm font-semibold text-gray-900 dark:text-white">
                                        {user.phone || 'Sin registrar'}
                                    </p>
                                </div>
                            </div>
                        </div>
                    </div>

                    {/* Tarjeta de Seguridad */}
                    <div className="bg-gradient-to-br from-indigo-900 to-purple-900 text-white rounded-3xl p-6 shadow-xl relative overflow-hidden">
                        <div className="absolute top-0 right-0 -mr-8 -mt-8 w-32 h-32 bg-white/10 rounded-full blur-2xl"></div>
                        <div className="relative z-10">
                            <h3 className="text-lg font-bold mb-2">Seguridad de Cuenta</h3>
                            <p className="text-indigo-200 text-sm mb-6">Protege tu cuenta con autenticación segura.</p>

                            <button
                                onClick={() => setShowChangePassword(true)}
                                className="w-full py-3 bg-white/10 hover:bg-white/20 backdrop-blur-sm border border-white/20 rounded-xl text-sm font-semibold transition-all flex items-center justify-center gap-2"
                            >
                                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                                </svg>
                                Cambiar Contraseña
                            </button>
                        </div>
                    </div>
                </div>
            </div>

            {/* Modal de Cambio de Contraseña */}
            <Modal
                isOpen={showChangePassword}
                onClose={() => setShowChangePassword(false)}
                title="Actualizar Seguridad"
                size="md"
            >
                <div className="pt-2">
                    {passwordSuccess ? (
                        <div className="flex flex-col items-center justify-center py-8 text-center animate-in zoom-in duration-300">
                            <div className="h-16 w-16 bg-green-100 text-green-500 rounded-full flex items-center justify-center mb-4">
                                <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                                </svg>
                            </div>
                            <h3 className="text-xl font-bold text-gray-900 mb-2">¡Contraseña Actualizada!</h3>
                            <p className="text-gray-500">Tu cuenta está segura con tu nueva contraseña.</p>
                        </div>
                    ) : (
                        <ChangePasswordForm
                            onSuccess={handlePasswordSuccess}
                            onCancel={() => setShowChangePassword(false)}
                        />
                    )}
                </div>
            </Modal>
        </div>
    );
}
