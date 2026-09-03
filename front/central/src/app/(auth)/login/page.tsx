'use client';

import { LoginForm } from '@/services/auth/login/ui';
import { useSearchParams } from 'next/navigation';
import { useEffect, Suspense, useState } from 'react';
import { CookieStorage } from '@/shared/utils';
import { useShopifyAuth } from '@/providers/ShopifyAuthProvider';
import { useRouter } from 'next/navigation';
import { ThemeToggle } from '@/shared/ui/theme-toggle';
import { LoginHeroImage } from '@/shared/ui/login-hero-image';
import { LoginBubbleCard } from '@/shared/ui/login-bubble-card';

function LoginContent() {
  const searchParams = useSearchParams();
  const { isShopifyEmbedded, sessionToken, isLoading: isShopifyLoading } = useShopifyAuth();
  const router = useRouter();
  const [isAuthenticating, setIsAuthenticating] = useState(false);
  const [isDark, setIsDark] = useState(false);

  useEffect(() => {
    const error = searchParams.get('error');
    if (error === 'no_business') {
      console.warn('Usuario no tiene negocio asignado. Contacte al administrador.');
    }
  }, [searchParams]);

  useEffect(() => {
    const htmlElement = document.documentElement;
    setIsDark(htmlElement.classList.contains('dark'));

    const observer = new MutationObserver(() => {
      setIsDark(htmlElement.classList.contains('dark'));
    });

    observer.observe(htmlElement, { attributes: true });
    return () => observer.disconnect();
  }, []);

  useEffect(() => {
    const authenticateWithShopify = async () => {
      if (isShopifyEmbedded && sessionToken) {
        setIsAuthenticating(true);
        try {
          const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || 'https://www.probabilityia.com.co/api/v1';
          const response = await fetch(`${baseUrl}/integrations/shopify/auth/login`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({ session_token: sessionToken }),
            credentials: 'include',
          });

          if (response.ok) {
            const data = await response.json();
            if (data.user) {
              CookieStorage.setUser(data.user);
            }

            console.log('✅ Login con Shopify exitoso, redirigiendo...');
            router.push('/home');
          } else {
            console.error('Fallo login con Shopify', response.status);
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
    <div className="h-screen w-screen relative overflow-hidden">
      {/* Portada a pantalla completa */}
      <div className={`absolute inset-0 ${isDark ? 'bg-[#14093a]' : 'bg-[#F0EEFF]'}`}>
        <LoginHeroImage />
      </div>

      {/* Velo para que el formulario se lea sobre cualquier foto */}
      <div
        className="absolute inset-0 pointer-events-none"
        style={{
          background: isDark
            ? 'radial-gradient(120% 90% at 50% 50%, rgba(10,6,26,0.55) 0%, rgba(10,6,26,0.78) 100%)'
            : 'radial-gradient(120% 90% at 50% 50%, rgba(255,255,255,0.15) 0%, rgba(60,40,110,0.32) 100%)',
        }}
      />

      {/* Theme Toggle */}
      <div className="fixed top-5 left-5 z-50">
        <ThemeToggle />
      </div>

      {/* Login emergiendo como burbuja */}
      <div className="absolute inset-0 z-30 flex items-center justify-center px-5 py-8 overflow-y-auto">
        <LoginBubbleCard isDark={isDark}>
          <LoginForm />
        </LoginBubbleCard>
      </div>
    </div>
  );
}

export default function LoginPage() {
  return (
    <Suspense fallback={
      <div className="min-h-screen flex items-center justify-center bg-white">
        <div className="text-gray-500 dark:text-gray-400">Cargando...</div>
      </div>
    }>
      <LoginContent />
    </Suspense>
  );
}
