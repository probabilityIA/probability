import { useState } from 'preact/hooks';

export default function Calculator() {
	const [pedidos, setPedidos] = useState(1000);
	const [ticket, setTicket] = useState(150000);
	const [tasaDevolucion, setTasaDevolucion] = useState(20);

	// Calcular ahorro estimado (simplificado)
	const ahorroMensual = Math.round((pedidos * ticket * tasaDevolucion / 100) * 0.21);

	const formatCurrency = (value: number) => {
		return new Intl.NumberFormat('es-CO', {
			style: 'currency',
			currency: 'COP',
			minimumFractionDigits: 0,
			maximumFractionDigits: 0,
		}).format(value);
	};

	return (
		<>
			<style>{`
				.slider::-webkit-slider-thumb {
					appearance: none;
					width: 24px;
					height: 24px;
					border-radius: 50%;
					background: #facc15;
					cursor: pointer;
					box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
				}

				.slider::-moz-range-thumb {
					width: 24px;
					height: 24px;
					border-radius: 50%;
					background: #facc15;
					cursor: pointer;
					border: none;
					box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
				}
			`}</style>
			<section id="calculadora" class="bg-white py-20 px-6">
				<div class="container mx-auto max-w-7xl">
					<div class="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
						{/* Left Column - Info */}
						<div>
						<h2 class="text-4xl md:text-5xl font-bold text-gray-900 mb-6">
							Calcula tu Ahorro Potencial
						</h2>
						<p class="text-lg text-gray-600 mb-8">
							Las tasas de devolución en Colombia oscilan entre el <strong class="text-gray-900">15% y el 30%</strong>. Descubre cuánto dinero estás perdiendo y cuánto podrías recuperar implementando ProbabilityIA.
						</p>
						<ul class="space-y-4">
							<li class="flex items-start gap-3">
								<div class="w-6 h-6 rounded-full bg-green-100 border-2 border-green-500 flex items-center justify-center flex-shrink-0 mt-0.5">
									<svg class="w-4 h-4 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7"></path>
									</svg>
								</div>
								<span class="text-gray-700">Reducción inmediata de logística inversa.</span>
							</li>
							<li class="flex items-start gap-3">
								<div class="w-6 h-6 rounded-full bg-green-100 border-2 border-green-500 flex items-center justify-center flex-shrink-0 mt-0.5">
									<svg class="w-4 h-4 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7"></path>
									</svg>
								</div>
								<span class="text-gray-700">Ahorro en costos operativos y reempaquetado.</span>
							</li>
							<li class="flex items-start gap-3">
								<div class="w-6 h-6 rounded-full bg-green-100 border-2 border-green-500 flex items-center justify-center flex-shrink-0 mt-0.5">
									<svg class="w-4 h-4 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7"></path>
									</svg>
								</div>
								<span class="text-gray-700">Mejora la satisfacción y retención del cliente.</span>
							</li>
						</ul>
					</div>

					{/* Right Column - Calculator */}
					<div class="bg-purple-900 rounded-2xl shadow-2xl p-8">
						<div class="flex items-center gap-3 mb-8">
							<div class="w-12 h-12 bg-yellow-400 rounded-lg flex items-center justify-center">
								<svg class="w-6 h-6 text-purple-900" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"></path>
								</svg>
							</div>
							<h3 class="text-2xl font-bold text-white">Simulador de ROI</h3>
						</div>

						{/* Pedidos Mensuales */}
						<div class="mb-6">
							<div class="flex justify-between items-center mb-2">
								<label class="text-white text-sm font-medium">Pedidos Mensuales</label>
								<span class="text-white font-semibold">{pedidos.toLocaleString('es-CO')}</span>
							</div>
							<input
								type="range"
								min="100"
								max="10000"
								step="100"
								value={pedidos}
								onInput={(e) => setPedidos(Number((e.target as HTMLInputElement).value))}
								class="w-full h-2 bg-purple-700 rounded-lg appearance-none cursor-pointer slider"
							/>
						</div>

						{/* Ticket Promedio */}
						<div class="mb-6">
							<div class="flex justify-between items-center mb-2">
								<label class="text-white text-sm font-medium">Ticket Promedio (COP)</label>
								<span class="text-white font-semibold">{formatCurrency(ticket)}</span>
							</div>
							<input
								type="range"
								min="50000"
								max="500000"
								step="10000"
								value={ticket}
								onInput={(e) => setTicket(Number((e.target as HTMLInputElement).value))}
								class="w-full h-2 bg-purple-700 rounded-lg appearance-none cursor-pointer slider"
							/>
						</div>

						{/* Tasa Devolución */}
						<div class="mb-8">
							<div class="flex justify-between items-center mb-2">
								<label class="text-white text-sm font-medium">Tasa Devolución Actual</label>
								<span class="text-white font-semibold">{tasaDevolucion}%</span>
							</div>
							<input
								type="range"
								min="5"
								max="40"
								step="1"
								value={tasaDevolucion}
								onInput={(e) => setTasaDevolucion(Number((e.target as HTMLInputElement).value))}
								class="w-full h-2 bg-purple-700 rounded-lg appearance-none cursor-pointer slider"
							/>
						</div>

						{/* Ahorro Estimado */}
						<div class="bg-purple-800 rounded-lg p-6 mb-6">
							<p class="text-white text-sm mb-2">Ahorro Estimado Mensual con ProbabilityIA</p>
							<div class="text-4xl font-bold text-green-400 mb-2">{formatCurrency(ahorroMensual)}</div>
							<p class="text-white text-xs opacity-75">*Basado en una reducción promedio del 21% en devoluciones.</p>
						</div>

						{/* CTA Button */}
						<button class="w-full bg-yellow-400 hover:bg-yellow-500 text-gray-900 font-bold py-4 rounded-lg flex items-center justify-center gap-2 transition-colors">
							Obtener Análisis Completo
							<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
							</svg>
						</button>
						</div>
					</div>
				</div>
			</section>
		</>
	);
}
