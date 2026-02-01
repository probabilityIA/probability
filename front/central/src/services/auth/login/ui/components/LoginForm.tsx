'use client';

import { useState, useTransition } from 'react';
import { loginAction, getRolesPermissionsAction } from '../../infra/actions';
import { TokenStorage } from '@/shared/config';
import { useRouter } from 'next/navigation';
import { EnvelopeIcon, LockClosedIcon, EyeIcon, EyeSlashIcon } from '@heroicons/react/24/outline';

export const LoginForm = () => {
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
                const response = await loginAction({ email, password });
                if (response.success) {
                    // ✅ NO guardar token (viene en cookie HttpOnly del backend)
                    // Solo guardar datos del usuario en sessionStorage
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

                    // Obtener roles y permisos (cookie se envía automáticamente)
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
                        }
                    } catch (permErr) {
                        console.warn('No se pudieron obtener los permisos:', permErr);
                        if (response.data.is_super_admin) {
                            TokenStorage.setPermissions({
                                is_super: true,
                                business_id: 0,
                                business_name: '',
                                role_id: 0,
                                role_name: 'Super Admin',
                                resources: [],
                            });
                        }
                    }

                    console.log('✅ Login exitoso, redirigiendo...');
                    router.push('/home');
                }
            } catch (err: any) {
                console.error(err);
                setError(err.message || 'Credenciales inválidas. Por favor intenta de nuevo.');
            }
        });
    };

        return (
        <div className="mb-10 mt-1 flex h-full w-full items-center justify-center px-2 md:mx-0 md:px-0 lg:mb-10 lg:items-center lg:justify-start text-gray-900 dark:text-white">
            <div className="mt-[1vh] w-full max-w-full flex-col items-center md:pl-4 lg:pl-0 xl:max-w-[420px]">
                <h1 className="mb-2.5 text-5xl font-bold text-navy-700 dark:text-black">¡Bienvenido!</h1>
                <p className="mb-9 ml-1 text-base text-gray-600">Inicia sesión con tu correo electrónico y contraseña.</p>

              

                <form onSubmit={handleSubmit} className="w-full">
                    <div className="mb-5">
                        <label htmlFor="email" className="mb-2 block text-sm font-medium text-gray-600">Correo</label>
                        <div className="relative">
                                        <div className="absolute inset-y-0 left-3 flex items-center pointer-events-none">
                                            <EnvelopeIcon className="w-5 h-5 text-gray-400" />
                                        </div>
                        <input
                            id="email"
                            name="email"
                            type="text"
                            placeholder="usuario@gmail.com"
                            required
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            className="w-full rounded-lg border border-gray-200 pl-10 px-4 py-3 text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-[rgba(124,58,237,0.12)] focus:border-[#7c3aed]"
                        />
                        </div>
                    </div>

                    <div className="mb-2">
                        <label htmlFor="password" className="mb-2 block text-sm font-medium text-gray-600">Contraseña</label>
                                                <div className="relative">
                                                    <div className="absolute inset-y-0 left-3 flex items-center pointer-events-none">
                                                        <LockClosedIcon className="w-5 h-5 text-gray-400" />
                                                    </div>
                                                                                <button
                                                                                    type="button"
                                                                                    aria-label={showPassword ? 'Ocultar contraseña' : 'Ver contraseña'}
                                                                                    onClick={() => setShowPassword((s) => !s)}
                                                                                    className="absolute inset-y-0 right-3 flex items-center focus:outline-none active:scale-95"
                                                                                >
                                                                                    {showPassword ? (
                                                                                        <EyeSlashIcon className={`w-5 h-5 transition-colors duration-150 ease-in-out ${showPassword ? 'text-[#7c3aed]' : 'text-gray-400'}`} />
                                                                                    ) : (
                                                                                        <EyeIcon className={`w-5 h-5 transition-colors duration-150 ease-in-out ${showPassword ? 'text-[#7c3aed]' : 'text-gray-400'}`} />
                                                                                    )}
                                                                                </button>
                                                <input
                            id="password"
                            name="password"
                                                        type={showPassword ? 'text' : 'password'}
                            placeholder="Min. 8 characters"
                            required
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                                                        className="w-full rounded-lg border border-gray-200 pl-10 pr-10 px-4 py-3 text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-[rgba(124,58,237,0.12)] focus:border-[#7c3aed]"
                        />
                        </div>
                    </div>

                    <div className="mb-8 flex items-center justify-between px-2">
                      
                        <a className="text-sm font-medium text-[#a78bfa] hover:text-[#6d28d9]" href=" ">¿Olvidó su contraseña?</a>
                    </div>

                    {error && (
                        <div className="mb-3 p-3 rounded-lg bg-red-50 text-red-500 text-sm">{error}</div>
                    )}

                    <button type="submit" disabled={isPending} className="linear w-full rounded-xl bg-[#7c3aed] py-3 text-base font-medium text-white transition duration-200 hover:bg-[#6d28d9] active:bg-[#5b21b6] dark:bg-[#7c3aed] dark:text-white dark:hover:bg-[#6d28d9]">
                        {isPending ? 'Iniciando Sesión...' : 'Iniciar Sesión'}
                    </button>

                    <div className="mt-4 text-center">
                       
                    </div>
                </form>
            </div>
        </div>
    );
};

