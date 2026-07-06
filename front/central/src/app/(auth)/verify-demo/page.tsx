'use client';

import { useState, Suspense } from 'react';
import Link from 'next/link';
import { useSearchParams, useRouter } from 'next/navigation';
import { demoVerifyOtpAction } from '@/services/auth/login/infra/actions';

function VerifyDemoContent() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const email = searchParams.get('email') || '';

  const [code, setCode] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [done, setDone] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    if (!/^\d{6}$/.test(code)) {
      setError('El codigo debe tener 6 digitos');
      return;
    }
    setLoading(true);
    try {
      const result = await demoVerifyOtpAction(email, code);
      if (result.success) {
        setDone(true);
        setTimeout(() => router.push('/login'), 2500);
      } else {
        setError(result.error || 'Codigo invalido o expirado');
      }
    } catch {
      setError('Error al conectar con el servidor');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 px-4">
      <div className="w-full max-w-md bg-white dark:bg-gray-800 rounded-2xl shadow-lg p-8">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Verifica tu demo</h1>
        <p className="mt-2 text-sm text-gray-500 dark:text-gray-400">
          Escribe el codigo de 6 digitos que enviamos por WhatsApp para activar tu cuenta.
        </p>

        {done ? (
          <div className="mt-6 rounded-lg bg-green-50 dark:bg-green-900/30 border border-green-200 dark:border-green-800 p-4 text-sm text-green-800 dark:text-green-300">
            Cuenta verificada. Redirigiendo al inicio de sesion...
          </div>
        ) : (
          <form onSubmit={handleSubmit} className="mt-6 space-y-4">
            <div>
              <label htmlFor="code" className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                Codigo de verificacion
              </label>
              <input
                id="code"
                inputMode="numeric"
                autoComplete="one-time-code"
                maxLength={6}
                required
                value={code}
                onChange={(e) => setCode(e.target.value.replace(/\D/g, ''))}
                placeholder="000000"
                className="mt-1 w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 px-3 py-2 text-center text-2xl tracking-[0.5em] font-semibold text-gray-900 dark:text-white focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
              />
            </div>

            {error && (
              <div className="rounded-lg bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-800 p-3 text-sm text-red-700 dark:text-red-300">
                {error}
              </div>
            )}

            <button
              type="submit"
              disabled={loading}
              className="w-full rounded-lg bg-indigo-600 py-2.5 text-sm font-semibold text-white hover:bg-indigo-500 disabled:opacity-60"
            >
              {loading ? 'Verificando...' : 'Verificar y activar'}
            </button>
          </form>
        )}

        <Link href="/login" className="mt-6 inline-block text-sm font-medium text-indigo-600 hover:text-indigo-500">
          Volver al inicio de sesion
        </Link>
      </div>
    </div>
  );
}

export default function VerifyDemoPage() {
  return (
    <Suspense fallback={<div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900">Cargando...</div>}>
      <VerifyDemoContent />
    </Suspense>
  );
}
