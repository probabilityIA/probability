'use client';

import { LoginForm } from '@/services/auth/login/ui';
import { useSearchParams } from 'next/navigation';
import { useEffect, Suspense, useState } from 'react';
import Image from 'next/image';
import { CookieStorage } from '@/shared/utils';
import { useShopifyAuth } from '@/providers/ShopifyAuthProvider';
import { useRouter } from 'next/navigation';

function LoginContent() {
  const searchParams = useSearchParams();
  const { isShopifyEmbedded, sessionToken, isLoading: isShopifyLoading } = useShopifyAuth();
  const router = useRouter();
  const [isAuthenticating, setIsAuthenticating] = useState(false);

  useEffect(() => {
    const error = searchParams.get('error');
    if (error === 'no_business') {
      console.warn('Usuario no tiene negocio asignado. Contacte al administrador.');
    }
  }, [searchParams]);

  useEffect(() => {
    const authenticateWithShopify = async () => {
      if (isShopifyEmbedded && sessionToken) {
        setIsAuthenticating(true);
        try {
          const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || 'https://app.probabilityia.com.co/api/v1';
          const response = await fetch(`${baseUrl}/integrations/shopify/auth/login`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({ session_token: sessionToken }),
            credentials: 'include', // ✅ IMPORTANTE: incluir cookies
          });

          if (response.ok) {
            const data = await response.json();
            // ✅ Cookie ya está seteada en el navegador por el backend
            // Solo guardar datos del usuario
            if (data.user) {
              CookieStorage.setUser(data.user);
            }

            console.log('✅ Login con Shopify exitoso, redirigiendo...');
            router.push('/home');
          } else {
            console.error('Fallo login con Shopify', response.status);
            // Si falla, mostramos login normal o error
            setIsAuthenticating(false);
          }
        } catch (error) {
          console.error('Error autenticando con Shopify', error);
          setIsAuthenticating(false);
        }
      }
    };

    authenticateWithShopify();
  }, [isShopifyEmbedded, sessionToken, router]);

  if (isShopifyEmbedded && (isShopifyLoading || isAuthenticating)) {
    return (
      <div className="h-screen w-screen flex flex-col items-center justify-center bg-white dark:bg-navy-900">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mb-4"></div>
        <p className="text-gray-600 dark:text-gray-300">Autenticando con Shopify...</p>
      </div>
    );
  }

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

        <Image
          src="/banner.webp"
          alt="características de probability"
          fill
          className="object-cover"
          priority
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
