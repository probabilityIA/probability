'use client';

import { useState, useTransition, useEffect } from 'react';
import { loginAction, getRolesPermissionsAction, loginServerAction } from '../../infra/actions';
import { TokenStorage } from '@/shared/config';
import { applyBusinessTheme, resetTheme } from '@/shared/utils/apply-business-theme';
import { useRouter } from 'next/navigation';
import { EnvelopeIcon, LockClosedIcon, EyeIcon, EyeSlashIcon } from '@heroicons/react/24/outline';
import { getActionError } from '@/shared/utils/action-result';

export const LoginForm = () => {
    const router = useRouter();
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [showPassword, setShowPassword] = useState(false);
    const [isPending, startTransition] = useTransition();
    const [error, setError] = useState<string | null>(null);
    const [isDark, setIsDark] = useState(false);

    useEffect(() => {
        const htmlElement = document.documentElement;
        setIsDark(htmlElement.classList.contains('dark'));

        const observer = new MutationObserver(() => {
            setIsDark(htmlElement.classList.contains('dark'));
        });

        observer.observe(htmlElement, { attributes: true });
        return () => observer.disconnect();
    }, []);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);

        startTransition(async () => {
            try {
                const isLocalDev = typeof window !== 'undefined' && window.location.hostname === 'localhost';
                let response;

                if (isLocalDev) {
                    const result = await loginServerAction(email, password);
                    if (!result.success) throw new Error(result.error || 'Error al iniciar sesión');
                    response = result.data;
                } else {
                    const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || '/api/v1';
                    const loginResponse = await fetch(`${baseUrl}/auth/login`, {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ email, password }),
                        credentials: 'include',
                    });

                    if (!loginResponse.ok) {
                        const errorData = await loginResponse.json();
                        throw new Error(errorData.error || errorData.message || 'Error al iniciar sesión');
                    }
                    response = await loginResponse.json();
                }

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
                    } else {
                        resetTheme();
                    }

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
                                subscription_status: permissionsResponse.data.subscription_status,
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
                                subscription_status: 'active',
                            });
                        }
                    }

                    router.push('/home');
                }
            } catch (err: any) {
                console.error(err);
                setError(getActionError(err, 'Credenciales inválidas. Por favor intenta de nuevo.'));
            }
        });
    };

    const formClass = isDark ? 'login-form-dark' : 'login-form-light';

    return (
        <div className={`w-full max-w-sm ${formClass}`}>
            {/* Logo */}
            <div className={isDark ? 'login-logo-dark' : 'login-logo-light'}>
                <img
                    src="/logo (2).png"
                    alt="ProbabilityIA Logo"
                    className="w-8 h-8"
                />
                <div className={isDark ? 'login-logo-text-dark' : 'login-logo-text-light'}>
                    ProbabilityIA
                </div>
            </div>

            {/* Header */}
            <div className={isDark ? 'login-header-dark' : 'login-header-light'}>
                <h1 className={isDark ? 'login-title-dark' : 'login-title-light'}>
                    {isDark ? 'Bienvenido de vuelta' : '¡Bienvenido!'}
                </h1>
                <p className={isDark ? 'login-subtitle-dark' : 'login-subtitle-light'}>
                    {isDark ? 'Ingresa tus credenciales para acceder al panel de gestión.' : 'Inicia sesión con tu correo electrónico y contraseña.'}
                </p>
            </div>

            {/* Form */}
            <form onSubmit={handleSubmit} className="w-full">
                {/* Email Field */}
                <div className={isDark ? 'login-form-group-dark' : 'login-form-group-light'}>
                    <label className={isDark ? 'login-label-dark' : 'login-label-light'}>
                        {isDark ? 'Correo electrónico' : 'Correo'}
                    </label>
                    <div className={isDark ? 'login-input-wrapper-dark' : 'login-input-wrapper-light'}>
                        <EnvelopeIcon className="w-5 h-5" />
                        <input
                            type="email"
                            placeholder={isDark ? 'correo@ejemplo.com' : 'usuario@gmail.com'}
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            required
                        />
                    </div>
                </div>

                {/* Password Field */}
                <div className={isDark ? 'login-form-group-dark' : 'login-form-group-light'}>
                    {isDark ? (
                        <label className="login-label-dark">Contraseña</label>
                    ) : (
                        <div className="login-label-row-light">
                            <label className="login-label-light">Contraseña</label>
                            <a href="#" className="login-forgot-light">¿Olvidó su contraseña?</a>
                        </div>
                    )}
                    <div className={isDark ? 'login-input-wrapper-dark' : 'login-input-wrapper-light'}>
                        <LockClosedIcon className="w-5 h-5" />
                        <input
                            type={showPassword ? 'text' : 'password'}
                            placeholder="Contraseña"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            required
                        />
                        <button
                            type="button"
                            onClick={() => setShowPassword(!showPassword)}
                            aria-label={showPassword ? 'Ocultar contraseña' : 'Ver contraseña'}
                        >
                            {showPassword ? (
                                <EyeSlashIcon className="w-5 h-5" />
                            ) : (
                                <EyeIcon className="w-5 h-5" />
                            )}
                        </button>
                    </div>
                </div>

                {/* Error */}
                {error && (
                    <div className={`p-3 rounded-lg text-sm mb-5 ${isDark ? 'bg-red-900/30 text-red-400 border border-red-800/50' : 'bg-red-50 text-red-600 border border-red-200'}`}>
                        {error}
                    </div>
                )}

                {/* Submit Button */}
                <button
                    type="submit"
                    disabled={isPending}
                    className={isDark ? 'login-button-dark' : 'login-button-light'}
                >
                    {isPending ? 'Iniciando Sesión...' : 'Iniciar Sesión'}
                </button>
            </form>

            {/* Footer */}
            <div className={isDark ? 'login-footer-dark' : 'login-footer-light'} style={{ marginTop: '32px' }}>
                <a href="#">Términos</a>
                <a href="#">Planes</a>
                <a href="#">Contáctanos</a>
            </div>
        </div>
    );
};

