'use client';

import { LoginForm } from '@/services/auth/login/ui';
import { useSearchParams } from 'next/navigation';
import { useEffect, Suspense, useState } from 'react';
import { CookieStorage } from '@/shared/utils';
import { useShopifyAuth } from '@/providers/ShopifyAuthProvider';
import { useRouter } from 'next/navigation';
import { LoginHeroImage, LOGIN_CLIPS } from '@/shared/ui/login-hero-image';
import { LoginBubbleCard } from '@/shared/ui/login-bubble-card';
import Image from 'next/image';

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
    const teniaDark = htmlElement.classList.contains('dark');
    htmlElement.classList.remove('dark');
    return () => {
      if (teniaDark) htmlElement.classList.add('dark');
    };
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
      <div className="h-screen w-screen flex flex-col items-center justify-center bg-white">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mb-4"></div>
        <p className="text-gray-600">Autenticando con Shopify...</p>
      </div>
    );
  }

  const copy = SLIDES[slide];

  return (
    <div className="h-screen w-screen relative overflow-hidden bg-[#140c2d]">
      <LoginHeroImage index={slide} />

      <div className="relative z-30 flex h-full w-full flex-col">
        <header className="shrink-0 bg-white/92 backdrop-blur-md">
          <div className="mx-auto flex h-16 w-full max-w-[1180px] items-center justify-between gap-4 px-6">
            <div className="flex items-center gap-3">
              <Image
                src="/logo-wordmark-v2.png"
                alt="ProbabilityIA"
                width={150}
                height={36}
                priority
                className="h-8 w-auto object-contain"
              />
              <span className="hidden text-sm text-gray-500 sm:inline">
                {'Log\u00edstica y ventas en un solo panel'}
              </span>
            </div>

            <div className="flex items-center gap-4">
              <a
                href="https://www.probabilityia.com.co"
                target="_blank"
                rel="noopener noreferrer"
                className="hidden text-sm font-medium text-gray-600 transition-colors hover:text-[#8B5CF6] sm:inline"
              >
                Sitio web
              </a>
              <a
                href="https://meet.brevo.com/probability-ia/reunion-de-30-minutos"
                target="_blank"
                rel="noopener noreferrer"
                className="rounded-full bg-[#8B5CF6] px-4 py-2 text-sm font-semibold text-white transition-colors hover:bg-[#7c3aed]"
              >
                Agendar demo
              </a>
            </div>
          </div>
        </header>

        <div className="relative min-h-0 flex-1 overflow-y-auto">
          <div
            aria-hidden="true"
            className="pointer-events-none absolute inset-y-0 right-0 hidden w-[34%] bg-white/92 backdrop-blur-md lg:block"
          />
          <div className="relative flex min-h-full w-full flex-col items-center justify-center gap-10 px-6 py-12 lg:flex-row lg:items-center lg:justify-start lg:gap-0 lg:px-0">

            <div className="w-full max-w-[520px] text-center lg:mx-auto lg:w-[66%] lg:max-w-none lg:px-10 lg:text-left">
              <div className="lg:mx-auto lg:max-w-[520px]">
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
            </div>

            <div className="relative z-10 flex w-full justify-center lg:w-[34%] lg:px-10">
              <div className="w-full max-w-[380px]">
                <LoginBubbleCard>
                  <LoginForm />
                </LoginBubbleCard>
              </div>
            </div>
          </div>
        </div>

        <footer className="shrink-0 bg-white/92 backdrop-blur-md">
          <div className="mx-auto flex w-full max-w-[1180px] flex-col items-center justify-between gap-2 px-6 py-3 text-xs text-gray-600 sm:flex-row sm:gap-6">
            <div className="flex flex-wrap items-center justify-center gap-x-5 gap-y-1">
              <a href="mailto:gerencia@probabilityapp.com" className="transition-colors hover:text-[#8B5CF6]">
                gerencia@probabilityapp.com
              </a>
              <a href="tel:+573138241302" className="transition-colors hover:text-[#8B5CF6]">
                +57 313 824 1302
              </a>
              <span className="hidden sm:inline">{'Lun a Vie 8:00 a.m. - 6:00 p.m.'}</span>
            </div>
            <p className="text-gray-500">
              {`\u00a9 ${new Date().getFullYear()} ProbabilityIA. Todos los derechos reservados.`}
            </p>
          </div>
        </footer>
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
