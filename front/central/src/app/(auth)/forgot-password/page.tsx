'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { recoveryChannelsAction, forgotPasswordAction } from '@/services/auth/login/infra/actions';

type Step = 'email' | 'channel' | 'sent';

export default function ForgotPasswordPage() {
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [step, setStep] = useState<Step>('email');
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');
  const [whatsapp, setWhatsapp] = useState<{ available: boolean; masked_phone: string }>({ available: false, masked_phone: '' });

  const handleEmailSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const result = await recoveryChannelsAction(email.trim());
      setWhatsapp(result.whatsapp || { available: false, masked_phone: '' });
      setStep('channel');
    } catch {
      setError('Error al conectar con el servidor');
    } finally {
      setLoading(false);
    }
  };

  const handleEmailChannel = async () => {
    setError('');
    setLoading(true);
    try {
      const result = await forgotPasswordAction(email.trim(), 'email');
      if (result.success) {
        setMessage(result.message || 'Si el correo esta registrado, recibiras un enlace para restablecer tu contrasena.');
        setStep('sent');
      } else {
        setError(result.error || 'No se pudo procesar la solicitud');
      }
    } catch {
      setError('Error al conectar con el servidor');
    } finally {
      setLoading(false);
    }
  };

  const handleWhatsAppChannel = async () => {
    setError('');
    setLoading(true);
    try {
      const result = await forgotPasswordAction(email.trim(), 'whatsapp');
      if (result.success) {
        router.push(`/verify-code?email=${encodeURIComponent(email.trim())}`);
      } else {
        setError(result.error || 'No se pudo procesar la solicitud');
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
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Recuperar contrasena</h1>

        {step === 'email' && (
          <>
            <p className="mt-2 text-sm text-gray-500 dark:text-gray-400">
              Ingresa tu usuario (correo) para continuar.
            </p>
            <form onSubmit={handleEmailSubmit} className="mt-6 space-y-4">
              <div>
                <label htmlFor="email" className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                  Correo electronico
                </label>
                <input
                  id="email"
                  type="email"
                  required
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="tucorreo@ejemplo.com"
                  className="mt-1 w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 px-3 py-2 text-gray-900 dark:text-white focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
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
                {loading ? 'Verificando...' : 'Continuar'}
              </button>
            </form>
          </>
        )}

        {step === 'channel' && (
          <>
            <p className="mt-2 text-sm text-gray-500 dark:text-gray-400">
              Elige como quieres recibir las instrucciones para restablecer tu contrasena.
            </p>
            <div className="mt-6 space-y-3">
              <button
                type="button"
                onClick={handleEmailChannel}
                disabled={loading}
                className="w-full flex items-center gap-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 px-4 py-3 text-left hover:border-indigo-500 disabled:opacity-60"
              >
                <span className="text-xl">✉️</span>
                <span>
                  <span className="block text-sm font-semibold text-gray-900 dark:text-white">Correo electronico</span>
                  <span className="block text-xs text-gray-500 dark:text-gray-400">Recibe un enlace en tu correo</span>
                </span>
              </button>

              <button
                type="button"
                onClick={handleWhatsAppChannel}
                disabled={loading || !whatsapp.available}
                className="w-full flex items-center gap-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 px-4 py-3 text-left hover:border-green-500 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <img src="/integrations/whatsapp.png" alt="WhatsApp" className="w-7 h-7 object-contain" />
                <span>
                  <span className="block text-sm font-semibold text-gray-900 dark:text-white">WhatsApp</span>
                  <span className="block text-xs text-gray-500 dark:text-gray-400">
                    {whatsapp.available
                      ? `Recibe un codigo en ${whatsapp.masked_phone}`
                      : 'No hay un telefono asociado a esta cuenta'}
                  </span>
                </span>
              </button>
            </div>

            {error && (
              <div className="mt-4 rounded-lg bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-800 p-3 text-sm text-red-700 dark:text-red-300">
                {error}
              </div>
            )}

            <button
              type="button"
              onClick={() => { setStep('email'); setError(''); }}
              className="mt-6 text-sm font-medium text-indigo-600 hover:text-indigo-500"
            >
              Cambiar correo
            </button>
          </>
        )}

        {step === 'sent' && (
          <div className="mt-6">
            <div className="rounded-lg bg-green-50 dark:bg-green-900/30 border border-green-200 dark:border-green-800 p-4 text-sm text-green-800 dark:text-green-300">
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
