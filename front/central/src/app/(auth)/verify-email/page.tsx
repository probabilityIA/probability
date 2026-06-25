'use client';

import { useState, useEffect, Suspense, useRef } from 'react';
import Link from 'next/link';
import { useSearchParams } from 'next/navigation';
import { verifyEmailAction } from '@/services/auth/login/infra/actions';

function VerifyEmailContent() {
  const searchParams = useSearchParams();
  const token = searchParams.get('token') || '';
  const [status, setStatus] = useState<'loading' | 'success' | 'error'>('loading');
  const [message, setMessage] = useState('');
  const ran = useRef(false);

  useEffect(() => {
    if (ran.current) return;
    ran.current = true;

    if (!token) {
      setStatus('error');
      setMessage('El enlace es invalido o esta incompleto.');
      return;
    }
    (async () => {
      const result = await verifyEmailAction(token);
      if (result.success) {
        setStatus('success');
        setMessage(result.message || 'Cuenta verificada. Ya puedes iniciar sesion.');
      } else {
        setStatus('error');
        setMessage(result.error || 'El enlace es invalido o ha expirado.');
      }
    })();
  }, [token]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 px-4">
      <div className="w-full max-w-md bg-white dark:bg-gray-800 rounded-2xl shadow-lg p-8 text-center">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Verificacion de cuenta</h1>

        {status === 'loading' && (
          <p className="mt-6 text-sm text-gray-500 dark:text-gray-400">Verificando tu cuenta...</p>
        )}

        {status === 'success' && (
          <div className="mt-6">
            <div className="rounded-lg bg-green-50 dark:bg-green-900/30 border border-green-200 dark:border-green-800 p-4 text-sm text-green-800 dark:text-green-300">
              {message}
            </div>
            <Link
              href="/login"
              className="mt-6 inline-block w-full rounded-lg bg-indigo-600 py-2.5 text-sm font-semibold text-white hover:bg-indigo-500"
            >
              Iniciar sesion
            </Link>
          </div>
        )}

        {status === 'error' && (
          <div className="mt-6">
            <div className="rounded-lg bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-800 p-4 text-sm text-red-700 dark:text-red-300">
              {message}
            </div>
            <Link href="/login" className="mt-6 inline-block text-sm font-medium text-indigo-600 hover:text-indigo-500">
              Volver al inicio de sesion
            </Link>
          </div>
        )}
      </div>
    </div>
  );
}

export default function VerifyEmailPage() {
  return (
    <Suspense fallback={
      <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900">
        <div className="text-gray-500 dark:text-gray-400">Cargando...</div>
      </div>
    }>
      <VerifyEmailContent />
    </Suspense>
  );
}
