import tinycolor from 'tinycolor2';

export function generateColorScale(baseColor: string): { [key: string]: string } {
  const base = tinycolor(baseColor);
  const scale: { [key: string]: string } = {};

  scale['50'] = tinycolor.mix(base, '#ffffff', 95).toHexString();
  scale['100'] = tinycolor.mix(base, '#ffffff', 90).toHexString();
  scale['200'] = tinycolor.mix(base, '#ffffff', 80).toHexString();
  scale['300'] = tinycolor.mix(base, '#ffffff', 60).toHexString();
  scale['400'] = tinycolor.mix(base, '#ffffff', 40).toHexString();
  scale['500'] = base.toHexString();
  scale['600'] = tinycolor.mix(base, '#000000', 20).toHexString();
  scale['700'] = tinycolor.mix(base, '#000000', 40).toHexString();
  scale['800'] = tinycolor.mix(base, '#000000', 60).toHexString();
  scale['900'] = tinycolor.mix(base, '#000000', 80).toHexString();

  return scale;
}

export function updateAllColorScales(
  primaryColor: string,
  secondaryColor: string,
  tertiaryColor: string,
  quaternaryColor: string
): void {
  const primaryScale = generateColorScale(primaryColor);
  const secondaryScale = generateColorScale(secondaryColor);
  const tertiaryScale = generateColorScale(tertiaryColor);
  const quaternaryScale = generateColorScale(quaternaryColor);

  const root = document.documentElement;

  root.style.setProperty('--color-primary', primaryColor);
  root.style.setProperty('--color-secondary', secondaryColor);
  root.style.setProperty('--color-tertiary', tertiaryColor);
  root.style.setProperty('--color-quaternary', quaternaryColor);

  Object.entries(primaryScale).forEach(([level, color]) => {
    root.style.setProperty(`--color-primary-${level}`, color);
  });

  Object.entries(secondaryScale).forEach(([level, color]) => {
    root.style.setProperty(`--color-secondary-${level}`, color);
  });

  Object.entries(tertiaryScale).forEach(([level, color]) => {
    root.style.setProperty(`--color-tertiary-${level}`, color);
  });

  Object.entries(quaternaryScale).forEach(([level, color]) => {
    root.style.setProperty(`--color-quaternary-${level}`, color);
  });
}
