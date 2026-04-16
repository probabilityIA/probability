/**
 * Utilidad para reproducir sonidos de notificación
 */

/**
 * Reproduce un sonido de notificación usando Web Audio API
 * Genera un sonido suave y agradable para notificaciones
 */
export function playNotificationSound(): void {
    try {
        // Crear contexto de audio
        const audioContext = new (window.AudioContext || (window as any).webkitAudioContext)();
        
        // Frecuencias para un sonido de notificación agradable (dos tonos)
        const frequencies = [523.25, 659.25]; // Do y Mi en la escala musical
        const duration = 0.15; // Duración de cada tono en segundos
        const gainValue = 0.3; // Volumen (0.0 a 1.0)
        
        frequencies.forEach((frequency, index) => {
            const oscillator = audioContext.createOscillator();
            const gainNode = audioContext.createGain();
            
            // Conectar los nodos
            oscillator.connect(gainNode);
            gainNode.connect(audioContext.destination);
            
            // Configurar el oscilador
            oscillator.frequency.value = frequency;
            oscillator.type = 'sine'; // Onda sinusoidal (sonido suave)
            
            // Configurar el volumen con fade in/out para evitar clicks
            gainNode.gain.setValueAtTime(0, audioContext.currentTime);
            gainNode.gain.linearRampToValueAtTime(gainValue, audioContext.currentTime + 0.01);
            gainNode.gain.linearRampToValueAtTime(0, audioContext.currentTime + duration);
            
            // Programar el inicio del sonido con un pequeño delay entre tonos
            const startTime = audioContext.currentTime + (index * 0.1);
            oscillator.start(startTime);
            oscillator.stop(startTime + duration);
        });
    } catch (error) {
        // Si falla el Web Audio API, intentar con un elemento de audio simple
        console.warn('Error al reproducir sonido con Web Audio API:', error);
        playFallbackSound();
    }
}

/**
 * Método alternativo usando un elemento de audio HTML
 * Genera un sonido simple usando un oscilador básico
 */
function playFallbackSound(): void {
    try {
        // Crear un elemento de audio temporal
        const audio = new Audio();
        
        // Usar un data URI con un sonido generado
        // Generamos un sonido simple usando un oscilador
        const audioContext = new (window.AudioContext || (window as any).webkitAudioContext)();
        const oscillator = audioContext.createOscillator();
        const gainNode = audioContext.createGain();
        
        oscillator.connect(gainNode);
        gainNode.connect(audioContext.destination);
        
        oscillator.frequency.value = 800; // Frecuencia en Hz
        oscillator.type = 'sine';
        
        gainNode.gain.setValueAtTime(0, audioContext.currentTime);
        gainNode.gain.linearRampToValueAtTime(0.2, audioContext.currentTime + 0.01);
        gainNode.gain.linearRampToValueAtTime(0, audioContext.currentTime + 0.2);
        
        oscillator.start(audioContext.currentTime);
        oscillator.stop(audioContext.currentTime + 0.2);
    } catch (error) {
        console.warn('No se pudo reproducir el sonido de notificación:', error);
    }
}

/**
 * Reproduce un sonido de notificación más corto y discreto
 */
export function playShortNotificationSound(): void {
    try {
        const audioContext = new (window.AudioContext || (window as any).webkitAudioContext)();
        const oscillator = audioContext.createOscillator();
        const gainNode = audioContext.createGain();
        
        oscillator.connect(gainNode);
        gainNode.connect(audioContext.destination);
        
        oscillator.frequency.value = 800;
        oscillator.type = 'sine';
        
        gainNode.gain.setValueAtTime(0, audioContext.currentTime);
        gainNode.gain.linearRampToValueAtTime(0.15, audioContext.currentTime + 0.01);
        gainNode.gain.linearRampToValueAtTime(0, audioContext.currentTime + 0.1);
        
        oscillator.start(audioContext.currentTime);
        oscillator.stop(audioContext.currentTime + 0.1);
    } catch (error) {
        console.warn('No se pudo reproducir el sonido de notificación:', error);
    }
}

