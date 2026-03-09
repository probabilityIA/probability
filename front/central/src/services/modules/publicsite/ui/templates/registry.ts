import { TemplateComponents } from './types';
import { defaultTemplate } from './default';
// Import custom templates here:
// import { modernTemplate } from './modern';

const templates: Record<string, TemplateComponents> = {
    default: defaultTemplate,
    // modern: modernTemplate,
};

/**
 * Returns the template component set for the given name.
 * Falls back to "default" if the template is not found.
 */
export function getTemplate(name: string): TemplateComponents {
    return templates[name] || templates['default'];
}

/**
 * Returns the list of available templates for the admin selector.
 */
export function getAvailableTemplates(): { id: string; name: string; description: string }[] {
    return [
        { id: 'default', name: 'Clasico', description: 'Diseno limpio y profesional con secciones modulares' },
        // { id: 'modern', name: 'Moderno', description: 'Diseno contemporaneo con animaciones' },
    ];
}
