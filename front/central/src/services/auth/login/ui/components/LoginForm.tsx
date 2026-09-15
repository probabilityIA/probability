'use client';

import { useState, useTransition } from 'react';
import { loginAction, loginServerAction } from '../../infra/actions';
import { TokenStorage } from '@/shared/config';
import { applyBusinessTheme, resetTheme } from '@/shared/utils/apply-business-theme';
import { useRouter, useSearchParams } from 'next/navigation';
import { EnvelopeIcon, LockClosedIcon, EyeIcon, EyeSlashIcon } from '@heroicons/react/24/outline';
import { getActionError } from '@/shared/utils/action-result';
import { DemoRegisterModal } from './DemoRegisterModal';
import { GoogleLogo, googleLoginUrl } from './GoogleButton';

const MEDIA_BASE =
    process.env.NEXT_PUBLIC_S3_BASE_URL ||
    'https://probability-media-assets.s3.us-east-1.amazonaws.com';

function destinoSeguro(next: string | null): string | null {
    if (!next) return null;
    if (!next.startsWith('/') || next.startsWith('//')) return null;
    return next;
}

export const LoginForm = () => {
    const router = useRouter();
    const searchParams = useSearchParams();
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [showPassword, setShowPassword] = useState(false);
    const [isPending, startTransition] = useTransition();
    const [error, setError] = useState<string | null>(null);
    const [showDemoModal, setShowDemoModal] = useState(false);
    const errorGoogle = searchParams.get('google_error');

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

                    const destino = destinoSeguro(searchParams.get('next'));
                    router.push(destino ?? (response.data.is_super_admin ? '/tickets' : '/home'));
                }
            } catch (err: any) {
                console.error(err);
                setError(getActionError(err, 'Credenciales inválidas. Por favor intenta de nuevo.'));
            }
        });
    };

    const formClass = 'login-form-light';

    return (
        <div className={`w-full max-w-sm ${formClass}`}>
            <div className="login-logo-light">
                <div className="login-logo-mark">
                    <img
                        src={`${MEDIA_BASE}/public/brand/logo-mark-badge.svg`}
                        alt="ProbabilityIA"
                    />
                </div>
                <div className="login-logo-text-light">
                    ProbabilityIA
                </div>
            </div>

            <div className="login-header-light">
                <h1 className="login-title-light">
                    {'¡Bienvenido!'}
                </h1>
            </div>

            <form onSubmit={handleSubmit} className="w-full">
                <div className="login-form-group-light">
                    <label className="login-label-light">
                        Correo
                    </label>
                    <div className="login-input-wrapper-light">
                        <EnvelopeIcon className="w-5 h-5" />
                        <input
                            type="email"
                            placeholder="usuario@gmail.com"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            required
                        />
                    </div>
                </div>

                <div className="login-form-group-light">
                    <div className="login-label-row-light">
                        <label className="login-label-light">Contraseña</label>
                        <a href="/forgot-password" className="login-forgot-light">¿Olvidó su contraseña?</a>
                    </div>
                    <div className="login-input-wrapper-light">
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

                {!error && errorGoogle && (
                    <div className="p-3 rounded-lg text-sm mb-5 bg-red-50 text-red-600 border border-red-200">
                        {errorGoogle}
                    </div>
                )}

                {error && (
                    <div className="p-3 rounded-lg text-sm mb-5 bg-red-50 text-red-600 border border-red-200">
                        {error}
                    </div>
                )}

                <button
                    type="submit"
                    disabled={isPending}
                    className="login-button-light"
                >
                    {isPending ? 'Iniciando Sesión...' : 'Iniciar Sesión'}
                </button>

                <div className="my-5 flex items-center gap-3">
                    <span className="h-px flex-1 bg-gray-200" />
                    <span className="text-xs font-medium uppercase tracking-wide text-gray-400">o</span>
                    <span className="h-px flex-1 bg-gray-200" />
                </div>

                <a
                    href={googleLoginUrl()}
                    className="flex w-full items-center justify-center gap-3 rounded-lg border border-gray-300 bg-white px-4 py-2.5 text-sm font-semibold text-gray-700 transition-colors hover:bg-gray-50"
                >
                    <GoogleLogo />
                    {'Continuar con Google'}
                </a>

                <div className="mt-4 text-center text-sm font-medium text-gray-800">
                    ¿No tienes cuenta?{' '}
                    <button
                        type="button"
                        onClick={() => setShowDemoModal(true)}
                        className="font-bold text-[#5b21b6] hover:text-[#4c1d95]"
                    >
                        Crea tu demo gratis
                    </button>
                </div>
            </form>

            {showDemoModal && <DemoRegisterModal onClose={() => setShowDemoModal(false)} />}

            <div className="login-footer-light" style={{ marginTop: '20px' }}>
                <a href="#">{'T\u00e9rminos'}</a>
                <a href="#">Planes</a>
                <a href="#">{'Cont\u00e1ctanos'}</a>
            </div>
        </div>
    );
};

