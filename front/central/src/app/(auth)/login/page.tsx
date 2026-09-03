'use client';

import { LoginForm } from '@/services/auth/login/ui';
import { useSearchParams } from 'next/navigation';
import { useEffect, Suspense, useState } from 'react';
import { CookieStorage } from '@/shared/utils';
import { useShopifyAuth } from '@/providers/ShopifyAuthProvider';
import { useRouter } from 'next/navigation';
import { ThemeToggle } from '@/shared/ui/theme-toggle';
import { LoginHeroImage, LOGIN_CLIPS } from '@/shared/ui/login-hero-image';
import { LoginBubbleCard } from '@/shared/ui/login-bubble-card';

const SLIDES = [
  {
    titulo: 'Vende por redes, nosotros hacemos el resto',
    detalle: 'Pedidos, env\u00edos y clientes en un solo lugar, sin saltar entre pesta\u00f1as.',
  },
  {
    titulo: 'Cada pedido, listo para salir',
    detalle: 'Cotiza y genera gu\u00edas con varias transportadoras desde el mismo panel.',
  },
  {
    titulo: 'Tus clientes, siempre informados',
    detalle: 'WhatsApp autom\u00e1tico en cada cambio de estado del env\u00edo.',
  },
  {
    titulo: 'Menos devoluciones, m\u00e1s entregas',
    detalle: 'Detecta a tiempo los pedidos que vienen de vuelta.',
  },
];

function LoginContent() {
  const searchParams = useSearchParams();
  const { isShopifyEmbedded, sessionToken, isLoading: isShopifyLoading } = useShopifyAuth();
  const router = useRouter();
  const [isAuthenticating, setIsAuthenticating] = useState(false);
  const [isDark, setIsDark] = useState(false);
  const [slide, setSlide] = useState(0);

  useEffect(() => {
    const error = searchParams.get('error');
    if (error === 'no_business') {
      console.warn('Usuario no tiene negocio asignado. Contacte al administrador.');
    }
  }, [searchParams]);

  useEffect(() => {
    const id = setTimeout(() => setSlide((i) => (i + 1) % LOGIN_CLIPS.length), 9000);
    return () => clearTimeout(id);
  }, [slide]);

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

  const copy = SLIDES[slide];

  return (
    <div className="h-screen w-screen relative overflow-hidden bg-[#140c2d]">
      {/* Fondo: los videos en carrusel, difuminados */}
      <LoginHeroImage index={slide} />

      {/* Theme Toggle */}
      <div className="fixed top-5 left-5 z-50">
        <ThemeToggle />
      </div>

      {/* Contenido centrado: mensaje a la izquierda, login a la derecha */}
      <div className="relative z-30 h-full w-full overflow-y-auto">
        <div className="mx-auto flex min-h-full w-full max-w-[1180px] flex-col items-center justify-center gap-10 px-6 py-12 lg:flex-row lg:items-center lg:justify-between lg:gap-16">

          <div className="w-full max-w-[520px] text-center lg:text-left">
            <h2
              key={copy.titulo}
              className="text-3xl sm:text-4xl font-extrabold leading-tight text-white drop-shadow-[0_2px_14px_rgba(0,0,0,0.45)]"
              style={{ animation: 'loginFadeUp 700ms ease-out' }}
            >
              {copy.titulo}
            </h2>
            <p
              key={copy.detalle}
              className="mt-4 text-base sm:text-lg text-white/85 drop-shadow-[0_1px_10px_rgba(0,0,0,0.4)]"
              style={{ animation: 'loginFadeUp 700ms ease-out 120ms both' }}
            >
              {copy.detalle}
            </p>

            <div className="mt-7 flex justify-center gap-2 lg:justify-start">
              {SLIDES.map((s2, i) => (
                <button
                  key={s2.titulo}
                  type="button"
                  onClick={() => setSlide(i)}
                  aria-label={`Ver mensaje ${i + 1}`}
                  className={`h-2 rounded-full transition-all ${i === slide ? 'w-7 bg-white' : 'w-2 bg-white/45 hover:bg-white/70'}`}
                />
              ))}
            </div>
          </div>

          <LoginBubbleCard isDark={isDark}>
            <LoginForm />
          </LoginBubbleCard>
        </div>
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
