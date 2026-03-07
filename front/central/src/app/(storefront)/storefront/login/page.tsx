'use client';

import { useState, useTransition } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { TokenStorage } from '@/shared/config';
import { loginServerAction, getRolesPermissionsAction } from '@/services/auth/login/infra/actions';
import { getActionError } from '@/shared/utils/action-result';
import { EnvelopeIcon, LockClosedIcon, EyeIcon, EyeSlashIcon } from '@heroicons/react/24/outline';
import { applyBusinessTheme } from '@/shared/utils/apply-business-theme';

export default function StorefrontLoginPage() {
    const router = useRouter();
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [showPassword, setShowPassword] = useState(false);
    const [isPending, startTransition] = useTransition();
    const [error, setError] = useState<string | null>(null);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);

        startTransition(async () => {
            try {
                const result = await loginServerAction(email, password);

                if (!result.success) {
                    throw new Error(result.error || 'Error al iniciar sesion');
                }

                const response = result.data;

                if (response.success) {
                    TokenStorage.setUser({
                        userId: response.data.user.id.toString(),
                        name: response.data.user.name,
                        email: response.data.user.email,
                        role: 'user',
                        avatarUrl: response.data.user.avatar_url,
                        is_super_admin: response.data.is_super_admin,
                        scope: response.data.scope,
                    });

                    if (response.data.businesses) {
                        TokenStorage.setBusinessesData(response.data.businesses);
                    }

                    if (!response.data.is_super_admin && response.data.businesses?.length > 0) {
                        applyBusinessTheme(response.data.businesses[0]);
                    }

                    // Get permissions
                    try {
                        const permissionsResponse = await getRolesPermissionsAction();
                        if (permissionsResponse.success && permissionsResponse.data) {
                            TokenStorage.setPermissions({
                                is_super: permissionsResponse.data.is_super,
                                business_id: permissionsResponse.data.business_id,
                                business_name: permissionsResponse.data.business_name,
                                role_id: permissionsResponse.data.role?.id || 0,
                                role_name: permissionsResponse.data.role?.name || '',
                                resources: permissionsResponse.data.resources || [],
                            });

                            // Check role level - if level 5 (cliente_final), go to storefront
                            const roleName = permissionsResponse.data.role?.name || '';
                            if (roleName === 'cliente_final') {
                                router.push('/storefront/catalogo');
                                return;
                            }
                        }
                    } catch {
                        // If permissions fail, still continue
                    }

                    // For non-cliente_final users, redirect to admin
                    router.push('/home');
                }
            } catch (err: any) {
                setError(getActionError(err, 'Credenciales invalidas. Intenta de nuevo.'));
            }
        });
    };

    return (
        <div className="min-h-screen flex items-center justify-center px-4">
            <div className="w-full max-w-md">
                <div className="text-center mb-8">
                    <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Iniciar Sesion</h1>
                    <p className="mt-2 text-gray-600 dark:text-gray-400">Ingresa con tu cuenta de cliente</p>
                </div>

                <form onSubmit={handleSubmit} className="bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-200 dark:border-gray-700 p-6 space-y-4">
                    <div>
                        <label htmlFor="email" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Correo</label>
                        <div className="relative">
                            <div className="absolute inset-y-0 left-3 flex items-center pointer-events-none">
                                <EnvelopeIcon className="w-5 h-5 text-gray-400" />
                            </div>
                            <input
                                id="email"
                                type="email"
                                value={email}
                                onChange={e => setEmail(e.target.value)}
                                placeholder="tu@email.com"
                                required
                                className="w-full pl-10 pr-4 py-2.5 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                            />
                        </div>
                    </div>

                    <div>
                        <label htmlFor="password" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Contrasena</label>
                        <div className="relative">
                            <div className="absolute inset-y-0 left-3 flex items-center pointer-events-none">
                                <LockClosedIcon className="w-5 h-5 text-gray-400" />
                            </div>
                            <input
                                id="password"
                                type={showPassword ? 'text' : 'password'}
                                value={password}
                                onChange={e => setPassword(e.target.value)}
                                placeholder="Tu contrasena"
                                required
                                className="w-full pl-10 pr-10 py-2.5 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                            />
                            <button
                                type="button"
                                onClick={() => setShowPassword(!showPassword)}
                                className="absolute inset-y-0 right-3 flex items-center"
                            >
                                {showPassword ? <EyeSlashIcon className="w-5 h-5 text-gray-400" /> : <EyeIcon className="w-5 h-5 text-gray-400" />}
                            </button>
                        </div>
                    </div>

                    {error && (
                        <div className="p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
                            <p className="text-sm text-red-600 dark:text-red-400">{error}</p>
                        </div>
                    )}

                    <button
                        type="submit"
                        disabled={isPending}
                        className="w-full py-2.5 bg-indigo-600 text-white font-medium rounded-lg hover:bg-indigo-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors"
                    >
                        {isPending ? 'Ingresando...' : 'Ingresar'}
                    </button>
                </form>

                <p className="mt-4 text-center text-sm text-gray-600 dark:text-gray-400">
                    No tienes cuenta?{' '}
                    <Link href="/storefront/registro" className="text-indigo-600 hover:text-indigo-700 font-medium">
                        Registrate
                    </Link>
                </p>
            </div>
        </div>
    );
}
