#!/bin/bash
# Script para iniciar la app Flutter en emulador Android
# Uso: ./run_emulator.sh [dev|staging|prod]

set -e

ANDROID_SDK_ROOT="${ANDROID_SDK_ROOT:-$HOME/Android/Sdk}"
EMULATOR="$ANDROID_SDK_ROOT/emulator/emulator"
ADB="$ANDROID_SDK_ROOT/platform-tools/adb"
AVD_NAME="Pixel_7"
ENV="${1:-prod}"

echo "=== Probability Mobile ==="
echo "Ambiente: $ENV"
echo "AVD: $AVD_NAME"
echo ""

# Verificar que el emulador existe
if ! $EMULATOR -list-avds 2>/dev/null | grep -q "$AVD_NAME"; then
  echo "Error: AVD '$AVD_NAME' no encontrado."
  echo "AVDs disponibles:"
  $EMULATOR -list-avds
  exit 1
fi

# Verificar si ya hay un emulador corriendo
if $ADB devices 2>/dev/null | grep -q "emulator.*device"; then
  echo "Emulador ya está corriendo."
else
  echo "Iniciando emulador..."
  $EMULATOR -avd $AVD_NAME -gpu auto &
  disown

  echo "Esperando a que el emulador arranque..."
  $ADB wait-for-device
  # Esperar a que termine de bootear
  while [ "$($ADB shell getprop sys.boot_completed 2>/dev/null | tr -d '\r')" != "1" ]; do
    sleep 2
  done
  echo "Emulador listo."
fi

echo ""
echo "Iniciando app Flutter (APP_ENV=$ENV)..."
flutter run -d emulator-5554 --dart-define=APP_ENV=$ENV
