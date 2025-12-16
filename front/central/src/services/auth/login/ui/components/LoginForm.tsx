'use client';

import { useState, useTransition } from 'react';
import { loginAction, getRolesPermissionsAction } from '../../infra/actions';
import { TokenStorage } from '@/shared/config';
import { useRouter } from 'next/navigation';

export const LoginForm = () => {
    const router = useRouter();
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [isPending, startTransition] = useTransition();
    const [error, setError] = useState<string | null>(null);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);

        startTransition(async () => {
            try {
                const response = await loginAction({ email, password });
                if (response.success) {
                    // Guardar token y datos del usuario
                    TokenStorage.setSessionToken(response.data.token);
                    TokenStorage.setUser({
                        userId: response.data.user.id.toString(),
                        name: response.data.user.name,
                        email: response.data.user.email,
                        role: 'user',
                        avatarUrl: response.data.user.avatar_url,
                        is_super_admin: response.data.is_super_admin,
                        scope: response.data.scope
                    });

                    if (response.data.businesses) {
                        TokenStorage.setBusinessesData(response.data.businesses);
                    }

                    // Obtener roles y permisos
                    try {
                        const permissionsResponse = await getRolesPermissionsAction(response.data.token);
                        if (permissionsResponse.success && permissionsResponse.data) {
                            TokenStorage.setPermissions({
                                is_super: permissionsResponse.data.is_super,
                                business_id: permissionsResponse.data.business_id,
                                business_name: permissionsResponse.data.business_name,
                                role_id: permissionsResponse.data.role?.id || 0,
                                role_name: permissionsResponse.data.role?.name || '',
                                resources: permissionsResponse.data.resources || []
                            });
                        }
                    } catch (permErr) {
                        console.warn('No se pudieron obtener los permisos:', permErr);
                        // Si es super admin, no necesita permisos del endpoint
                        if (response.data.is_super_admin) {
                            TokenStorage.setPermissions({
                                is_super: true,
                                business_id: 0,
                                business_name: '',
                                role_id: 0,
                                role_name: 'Super Admin',
                                resources: []
                            });
                        }
                    }

                    router.push('/home');
                }
            } catch (err: any) {
                console.error(err);
                setError(err.message || 'Credenciales inválidas. Por favor intenta de nuevo.');
            }
        });
    };

    return (
        <div className="w-full">
            <div className="mb-8 sm:mb-10">
                <h2 className="text-2xl sm:text-3xl font-bold text-gray-900 mb-2">
                    Iniciar Sesión
                </h2>
            </div>

            <form onSubmit={handleSubmit} className="space-y-6">
                {/* Email Field */}
                <div className="space-y-2">
                    <label htmlFor="email" className="block text-sm font-bold text-gray-400">
                        Correo Electrónico
                    </label>
                    <input
                        id="email"
                        name="email"
                        type="email"
                        autoComplete="email"
                        required
                        className="block w-full px-4 py-3 bg-white border border-gray-200 rounded-lg text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-green-500/20 focus:border-green-500 transition-all"
                        placeholder="tu@email.com"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                    />
                </div>

                {/* Password Field */}
                <div className="space-y-2">
                    <label htmlFor="password" className="block text-sm font-bold text-gray-400">
                        Contraseña
                    </label>
                    <input
                        id="password"
                        name="password"
                        type="password"
                        autoComplete="current-password"
                        required
                        className="block w-full px-4 py-3 bg-white border border-gray-200 rounded-lg text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-green-500/20 focus:border-green-500 transition-all"
                        placeholder="Tu contraseña"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                    />
                </div>

                {/* Error Message */}
                {error && (
                    <div className="p-3 rounded-lg bg-red-50 text-red-500 text-sm">
                        {error}
                    </div>
                )}

                {/* Submit Button */}
                <button
                    type="submit"
                    disabled={isPending}
                    className="w-full flex justify-center py-3 px-4 border border-transparent text-base font-bold rounded-lg text-white bg-[#4ade80] hover:bg-green-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 transition-all shadow-sm disabled:opacity-50 disabled:cursor-not-allowed mt-8"
                >
                    {isPending ? 'Iniciando sesión...' : 'Continuar'}
                </button>

                {/* Links */}
                <div className="space-y-4 mt-8 pt-4">
                    <div className="text-sm text-gray-500">
                        ¿No tienes cuenta? <a href="#" className="text-green-500 hover:underline font-medium">Regístrate</a>
                    </div>
                    <div className="text-sm text-gray-500">
                        ¿Olvidaste tu contraseña? <a href="#" className="text-green-500 hover:underline font-medium">Restablecer contraseña</a>
                    </div>
                </div>
            </form>
        </div>
    );
};
