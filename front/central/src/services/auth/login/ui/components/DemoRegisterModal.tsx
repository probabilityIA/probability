'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { demoRegisterAction } from '../../infra/actions';

interface DemoRegisterModalProps {
  onClose: () => void;
}

export const DemoRegisterModal = ({ onClose }: DemoRegisterModalProps) => {
  const router = useRouter();
  const [fullName, setFullName] = useState('');
  const [businessName, setBusinessName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [phone, setPhone] = useState('');
  const [channel, setChannel] = useState<'email' | 'whatsapp'>('email');
  const [loading, setLoading] = useState(false);
  const [done, setDone] = useState(false);
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    if (password.length < 6) {
      setError('La contrasena debe tener al menos 6 caracteres');
      return;
    }
    if (channel === 'whatsapp' && phone.trim().length < 7) {
      setError('Ingresa un telefono valido para verificar por WhatsApp');
      return;
    }
    setLoading(true);
    try {
      const result = await demoRegisterAction({
        full_name: fullName.trim(),
        business_name: businessName.trim(),
        email: email.trim(),
        password,
        phone: channel === 'whatsapp' ? phone.trim() : undefined,
        channel,
      });
      if (result.success) {
        if (channel === 'whatsapp') {
          router.push(`/verify-demo?email=${encodeURIComponent(email.trim())}`);
          return;
        }
        setDone(true);
        setMessage(result.message || 'Cuenta creada. Revisa tu correo para verificar tu cuenta.');
      } else {
        setError(result.error || 'No se pudo crear la demo');
      }
    } catch {
      setError('Error al conectar con el servidor');
    } finally {
      setLoading(false);
    }
  };

  const channelBtn = (value: 'email' | 'whatsapp', label: string, hint: string) => (
    <button
      type="button"
      onClick={() => setChannel(value)}
      className={`flex-1 rounded-lg border px-3 py-2 text-left ${
        channel === value
          ? 'border-indigo-500 ring-1 ring-indigo-500 bg-indigo-50 dark:bg-indigo-900/20'
          : 'border-gray-300 dark:border-gray-600'
      }`}
    >
      <span className="block text-sm font-semibold text-gray-900 dark:text-white">{label}</span>
      <span className="block text-xs text-gray-500 dark:text-gray-400">{hint}</span>
    </button>
  );

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4"
      onClick={onClose}
    >
      <div
        className="w-full max-w-md rounded-2xl bg-white dark:bg-gray-800 shadow-xl p-7 relative"
        onClick={(e) => e.stopPropagation()}
      >
        <button
          type="button"
          onClick={onClose}
          aria-label="Cerrar"
          className="absolute right-4 top-4 text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 text-xl leading-none"
        >
          x
        </button>

        <h2 className="text-xl font-bold text-gray-900 dark:text-white">Crea tu demo gratis</h2>
        <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
          Prueba Probability con datos simulados. Sin tarjeta, sin compromiso.
        </p>

        {done ? (
          <div className="mt-6">
            <div className="rounded-lg bg-green-50 dark:bg-green-900/30 border border-green-200 dark:border-green-800 p-4 text-sm text-green-800 dark:text-green-300">
              {message}
            </div>
            <button
              type="button"
              onClick={onClose}
              className="mt-6 w-full rounded-lg bg-indigo-600 py-2.5 text-sm font-semibold text-white hover:bg-indigo-500"
            >
              Entendido
            </button>
          </div>
        ) : (
          <form onSubmit={handleSubmit} className="mt-5 space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">Tu nombre</label>
              <input
                type="text"
                required
                value={fullName}
                onChange={(e) => setFullName(e.target.value)}
                placeholder="Juan Perez"
                className="mt-1 w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 px-3 py-2 text-gray-900 dark:text-white focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">Nombre del negocio</label>
              <input
                type="text"
                required
                value={businessName}
                onChange={(e) => setBusinessName(e.target.value)}
                placeholder="Mi Tienda"
                className="mt-1 w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 px-3 py-2 text-gray-900 dark:text-white focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">Correo electronico</label>
              <input
                type="email"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="tucorreo@ejemplo.com"
                className="mt-1 w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 px-3 py-2 text-gray-900 dark:text-white focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">Contrasena</label>
              <input
                type="password"
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Minimo 6 caracteres"
                className="mt-1 w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 px-3 py-2 text-gray-900 dark:text-white focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Como quieres verificar tu cuenta</label>
              <div className="flex gap-2">
                {channelBtn('email', 'Correo', 'Enlace por email')}
                {channelBtn('whatsapp', 'WhatsApp', 'Codigo al celular')}
              </div>
            </div>

            {channel === 'whatsapp' && (
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">Telefono (WhatsApp)</label>
                <input
                  type="tel"
                  inputMode="tel"
                  value={phone}
                  onChange={(e) => setPhone(e.target.value)}
                  placeholder="3001234567"
                  className="mt-1 w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 px-3 py-2 text-gray-900 dark:text-white focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                />
              </div>
            )}

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
              {loading ? 'Creando...' : 'Crear mi demo'}
            </button>
          </form>
        )}
      </div>
    </div>
  );
};
