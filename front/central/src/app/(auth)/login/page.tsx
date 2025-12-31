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
    <div className="min-h-screen flex bg-white dark:bg-navy-900">
      {/* Left Side - Login Form */}
      <div className="w-full lg:w-1/2 flex flex-col justify-center px-6 sm:px-12 lg:px-24 xl:px-32 relative">
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
      <div className="hidden lg:flex lg:w-1/2 flex-col justify-center px-12 relative overflow-hidden rounded-l-3xl">
        {/* Background Pattern/Gradient */}
  <div className="absolute inset-0 bg-linear-to-br from-navy-900 to-navy-800" style={{background: 'linear-gradient(135deg,#0b1530 0%,#1e3a8a 100%)'}} />

        <div className="relative z-10 max-w-lg mx-auto w-full space-y-8">
          {/* Logo */}
          <div className="flex items-center gap-3 mb-12">
            <div className="w-10 h-10 bg-white/10 rounded-lg flex items-center justify-center text-white font-bold text-xl">P</div>
            <h1 className="text-3xl font-bold text-white">Probability</h1>
          </div>

          {/* Cards Grid */}
          <div className="grid grid-cols-2 gap-4">
            {/* Stats Card */}
            <div className="bg-white rounded-xl p-5 shadow-lg">
              <h3 className="text-gray-800 font-semibold text-sm mb-4">Estadísticas por Departamento</h3>
              <div className="space-y-3">
                <div className="flex justify-between items-center">
                  <span className="text-[#7c3aed] font-bold text-lg">$8,035</span>
                  <span className="text-gray-400 text-xs">Ventas</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-[#7c3aed] font-bold text-lg">$4,684</span>
                  <span className="text-gray-400 text-xs">Marketing</span>
                </div>
              </div>
            </div>

            {/* Analysis Card */}
            <div className="bg-white rounded-xl p-5 shadow-lg">
              <h3 className="text-gray-800 font-semibold text-sm mb-2">Análisis Probabilístico</h3>
              <div className="mt-4">
                <p className="text-gray-800 font-bold text-sm">Modelos Predictivos</p>
                <p className="text-gray-400 text-xs mt-1">Análisis avanzado de datos</p>
              </div>
            </div>

            {/* Monthly Results Card (Full Width) */}
            <div className="col-span-2 bg-white rounded-xl p-5 shadow-lg">
              <h3 className="text-gray-800 font-semibold text-sm mb-4">Resultados del Mes</h3>
              <div className="flex justify-between items-end mb-2">
                <span className="text-[#7c3aed] font-bold text-2xl">$69,700</span>
                <span className="text-[#7c3aed] text-xs font-medium bg-[#f3e8ff] px-2 py-1 rounded">+2.2%</span>
              </div>
              <div className="w-full bg-gray-100 rounded-full h-1.5 mt-2">
                <div className="bg-[#7c3aed] h-1.5 rounded-full w-3/4"></div>
              </div>
            </div>
          </div>

          {/* Description Text */}
          <p className="text-blue-200 text-sm leading-relaxed mt-8">
            Probability es la plataforma que ayuda a empresas a optimizar sus decisiones mediante análisis probabilístico avanzado.
          </p>
        </div>
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
