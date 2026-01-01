'use client';

import { LoginForm } from '@/services/auth/login/ui';
import { useSearchParams } from 'next/navigation';
import { useEffect, useState, Suspense } from 'react';

function LoginContent() {
  const searchParams = useSearchParams();
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  useEffect(() => {
    const error = searchParams.get('error');
    if (error === 'no_business') {
      setErrorMessage('Usuario no tiene negocio asignado. Contacte al administrador.');
    }
  }, [searchParams]);

  return (
    <div className="h-screen w-screen flex bg-white dark:bg-navy-900">
      {/* Left Side - Login Form */}
        <div className="w-full md:w-[50vw] h-full flex items-center justify-center px-6 sm:px-12 md:px-24 xl:px-32 relative">
        <div className="w-full max-w-md mx-auto">
          <LoginForm />

          {/* Footer Links (restored) */}
          <div className="mt-12 flex flex-wrap gap-4 sm:gap-6 text-sm font-medium text-navy-700 dark:text-gray-400 justify-center sm:justify-start">
            <a href="#" className="text-[#7c3aed] hover:text-[#6d28d9]">Términos</a>
            <a href="#" className="text-[#7c3aed] hover:text-[#6d28d9]">Planes</a>
            <a href="#" className="text-[#7c3aed] hover:text-[#6d28d9]">Contáctanos</a>
          </div>
        </div>
      </div>

      {/* Right Side - Dashboard Preview */}
        <div
          className="hidden md:flex md:w-[50vw] h-full items-center justify-center relative overflow-hidden rounded-l-3xl z-2"
          style={{ boxShadow: '-2px 0 50px rgba(0,0,0,0.12)' }}
        >
        
        <img
          src="/banner.webp"
          alt="características de probability"
          className="absolute inset-0 w-full h-full object-cover"
        />
      </div>
    </div>
  );
}

export default function LoginPage() {
  return (
    <Suspense fallback={
      <div className="min-h-screen flex items-center justify-center bg-white">
        <div className="text-gray-500">Cargando...</div>
      </div>
    }>
      <LoginContent />
    </Suspense>
  );
}
