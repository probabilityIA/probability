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
  const [revealed, setRevealed] = useState(false);

  useEffect(() => {
    const error = searchParams.get('error');
    if (error === 'no_business') {
      console.warn('Usuario no tiene negocio asignado. Contacte al administrador.');
    }
  }, [searchParams]);

  useEffect(() => {
    const reduce = window.matchMedia?.('(prefers-reduced-motion: reduce)').matches;
    if (reduce) {
      setRevealed(true);
      return;
    }
    const id = setTimeout(() => setRevealed(true), 2200);
    return () => clearTimeout(id);
  }, []);

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
      {/* El video arranca a pantalla completa y se corre a la derecha cuando entra el panel */}
      <div
        className={`absolute top-0 right-0 bottom-0 overflow-hidden ${isDark ? 'bg-[#14093a]' : 'bg-[#F0EEFF]'} ${revealed ? 'lg:left-[calc(min(42vw,560px)-56px)]' : 'left-0'}`}
        style={{ transition: 'left 1000ms cubic-bezier(0.16, 1, 0.3, 1)' }}
      >
        <LoginHeroImage />
      </div>

      {/* Theme Toggle */}
      <div className="fixed top-5 left-5 z-50">
        <ThemeToggle />
      </div>

      {/* Panel del login: entra desde la izquierda y le cede el resto al video */}
      <div className="absolute inset-0 z-30 flex pointer-events-none">
        <LoginBubbleCard isDark={isDark} visible={revealed}>
          <LoginForm />
        </LoginBubbleCard>
      </div>

      {/* Enlaces al pie, sobre el video */}
      <div className="absolute bottom-5 left-0 right-0 z-40 flex justify-center lg:justify-end lg:pr-[6vw] gap-6 text-xs font-semibold text-white/85 drop-shadow-[0_1px_6px_rgba(0,0,0,0.6)]">
        <a href="#" className="hover:text-white transition-colors">{'T\u00e9rminos'}</a>
        <a href="#" className="hover:text-white transition-colors">Planes</a>
        <a href="#" className="hover:text-white transition-colors">{'Cont\u00e1ctanos'}</a>
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
