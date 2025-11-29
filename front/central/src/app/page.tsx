export default function Home() {
  return (
    <div className="relative flex min-h-screen items-center justify-center overflow-hidden bg-gradient-to-br from-indigo-950 via-purple-900 to-pink-900">
      {/* Animated background elements */}
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute -top-40 -right-40 h-80 w-80 rounded-full bg-purple-500/30 blur-3xl animate-pulse"></div>
        <div className="absolute -bottom-40 -left-40 h-80 w-80 rounded-full bg-indigo-500/30 blur-3xl animate-pulse delay-1000"></div>
      </div>

      {/* Main content */}
      <main className="relative z-10 flex flex-col items-center justify-center px-6 text-center">
        <div className="mb-8 inline-block rounded-2xl bg-white/10 px-6 py-3 backdrop-blur-sm">
          <p className="text-sm font-medium text-purple-200">Sistema de Gestión</p>
        </div>

        <h1 className="mb-6 bg-gradient-to-r from-white via-purple-100 to-pink-100 bg-clip-text text-6xl font-bold tracking-tight text-transparent sm:text-7xl md:text-8xl">
          Bienvenidos a
          <br />
          <span className="bg-gradient-to-r from-purple-300 to-pink-300 bg-clip-text">
            Probability
          </span>
        </h1>

        <p className="mb-12 max-w-2xl text-lg text-purple-100 sm:text-xl">
          Plataforma centralizada para la gestión inteligente de procesos y análisis probabilístico
        </p>

        <div className="flex flex-col gap-4 sm:flex-row">
          <button className="group relative overflow-hidden rounded-full bg-white px-8 py-4 font-semibold text-purple-900 transition-all hover:scale-105 hover:shadow-2xl hover:shadow-purple-500/50">
            <span className="relative z-10">Comenzar</span>
            <div className="absolute inset-0 -z-0 bg-gradient-to-r from-purple-400 to-pink-400 opacity-0 transition-opacity group-hover:opacity-100"></div>
          </button>

          <button className="rounded-full border-2 border-white/30 bg-white/10 px-8 py-4 font-semibold text-white backdrop-blur-sm transition-all hover:scale-105 hover:border-white/50 hover:bg-white/20">
            Más Información
          </button>
        </div>

        {/* Feature cards */}
        <div className="mt-20 grid gap-6 sm:grid-cols-3">
          <div className="rounded-2xl bg-white/10 p-6 backdrop-blur-sm transition-all hover:bg-white/20">
            <div className="mb-3 text-4xl">📊</div>
            <h3 className="mb-2 font-semibold text-white">Análisis</h3>
            <p className="text-sm text-purple-200">Herramientas avanzadas de análisis probabilístico</p>
          </div>

          <div className="rounded-2xl bg-white/10 p-6 backdrop-blur-sm transition-all hover:bg-white/20">
            <div className="mb-3 text-4xl">🎯</div>
            <h3 className="mb-2 font-semibold text-white">Precisión</h3>
            <p className="text-sm text-purple-200">Resultados exactos y confiables</p>
          </div>

          <div className="rounded-2xl bg-white/10 p-6 backdrop-blur-sm transition-all hover:bg-white/20">
            <div className="mb-3 text-4xl">⚡</div>
            <h3 className="mb-2 font-semibold text-white">Velocidad</h3>
            <p className="text-sm text-purple-200">Procesamiento rápido y eficiente</p>
          </div>
        </div>
      </main>
    </div>
  );
}
