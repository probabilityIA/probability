## Referencia de integración con Shopify (backend Probability)

Este documento resume y agrupa todo el **código base de la integración con Shopify** usado en este backend, para que puedas:

- Replicarlo en otro proyecto (Django u otro stack).
- Explicarle a otra IA qué hace cada pieza.
- Tener una guía rápida de los puntos a adaptar (nombres de modelos, URLs, versiones de API, etc.).

---

## 1. Modelos principales

### 1.1. `ShopifyIntegration`

Responsable de guardar la **integración por tienda Shopify** para cada usuario, incluyendo:

- `store_name` (por ejemplo `mitienda.myshopify.com`)
- `encrypted_access_token` (token de acceso encriptado con Fernet)
- `permissions`, `is_active`, `last_sync`

Código original:

```python
from django.db import models
from django.contrib.auth.models import User
from cryptography.fernet import Fernet
from django.conf import settings
import base64


class ShopifyIntegration(models.Model):
    """
    Modelo para almacenar las integraciones de Shopify de los usuarios
    """
    user = models.ForeignKey(User, on_delete=models.CASCADE, related_name='shopify_integrations')
    store_name = models.CharField(max_length=200, verbose_name='Nombre de la Tienda')
    encrypted_access_token = models.TextField(verbose_name='Token de Acceso Encriptado')
    is_active = models.BooleanField(default=True, verbose_name='Activa')
    permissions = models.JSONField(default=list, verbose_name='Permisos')
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    last_sync = models.DateTimeField(null=True, blank=True, verbose_name='Última Sincronización')
    
    class Meta:
        verbose_name = 'Integración de Shopify'
        verbose_name_plural = 'Integraciones de Shopify'
        ordering = ['-created_at']
        unique_together = ['user', 'store_name']
    
    def __str__(self):
        return f"{self.store_name} - {self.user.get_full_name()}"
    
    def set_access_token(self, token):
        """Encriptar y guardar el token de acceso"""
        try:
            key = self._get_encryption_key()
            fernet = Fernet(key)
            encrypted_token = fernet.encrypt(token.encode())
            self.encrypted_access_token = base64.b64encode(encrypted_token).decode()
        except Exception:
            # Si falla la encriptación, guardar en texto plano (solo para desarrollo)
            self.encrypted_access_token = token
    
    def get_access_token(self):
        """Desencriptar y obtener el token de acceso"""
        if not self.encrypted_access_token:
            return None
        try:
            # Intentar desencriptar primero
            key = self._get_encryption_key()
            fernet = Fernet(key)
            encrypted_token = base64.b64decode(self.encrypted_access_token.encode())
            decrypted_token = fernet.decrypt(encrypted_token)
            return decrypted_token.decode()
        except Exception:
            # Si falla la desencriptación, asumir que está en texto plano
            return self.encrypted_access_token
    
    def _get_encryption_key(self):
        """Obtener o generar la clave de encriptación"""
        key = getattr(settings, 'SHOPIFY_ENCRYPTION_KEY', None)
        if not key:
            # Usar una clave fija para desarrollo (NO usar en producción)
            key = b'probability_shopify_key_2025_development_only_not_production'
        else:
            # Convertir string a bytes si es necesario
            if isinstance(key, str):
                key = key.encode()
        return key
```

### 1.2. `Order`

Modelo que guarda cada **orden importada desde Shopify**, con campos de:

- Datos básicos de la orden (`platform_order_id`, `order_number`, `total_price`, `currency`, etc.).
- Datos de cliente y dirección.
- Información de estados (`financial_status`, `fulfillment_status`).
- Campos de analítica / scoring (`delivery_score`, `delivery_probability`, `recommended_carrier`, etc.).

Código original (solo parte relevante para integración, puedes simplificarlo en tu nuevo proyecto):

```python
class PlatformType(models.TextChoices):
    SHOPIFY = 'shopify', 'Shopify'


class Order(models.Model):
    integration = models.ForeignKey(
        ShopifyIntegration, on_delete=models.CASCADE, related_name='orders'
    )

    platform = models.CharField(
        max_length=32,
        choices=PlatformType.choices,
        default=PlatformType.SHOPIFY,
        help_text='Plataforma origen de la orden'
    )
    platform_order_id = models.CharField(max_length=128, blank=True, null=True)
    order_number = models.CharField(max_length=64, blank=True, null=True)

    total_price = models.DecimalField(max_digits=12, decimal_places=2, null=True, blank=True)
    currency = models.CharField(max_length=12, blank=True, null=True)
    payment_type = models.CharField(max_length=64, blank=True, null=True)

    customer_name = models.CharField(max_length=128, blank=True, null=True)
    customer_email = models.EmailField(blank=True, null=True)
    phone = models.CharField(max_length=32, blank=True, null=True)

    country = models.CharField(max_length=64, blank=True, null=True)
    province = models.CharField(max_length=64, blank=True, null=True)
    city = models.CharField(max_length=64, blank=True, null=True)
    address = models.CharField(max_length=255, blank=True, null=True)
    address_complement = models.CharField(max_length=255, blank=True, null=True)

    source_name = models.CharField(max_length=64, blank=True, null=True)
    financial_status = models.CharField(max_length=64, blank=True, null=True)
    fulfillment_status = models.CharField(max_length=64, blank=True, null=True)

    created_at = models.DateTimeField(null=True, blank=True)
    imported_at = models.DateTimeField(auto_now_add=True)

    purchase_count = models.PositiveIntegerField(default=0)
    delivery_score = models.FloatField(default=0.0)
    delivery_probability = models.CharField(max_length=32, blank=True, null=True)
    recommended_carrier = models.CharField(max_length=128, blank=True, null=True)

    customer_json = models.JSONField(default=dict, blank=True)
    line_items = models.JSONField(default=list, blank=True)
    raw_response = models.JSONField(default=dict, blank=True)

    class Meta:
        ordering = ['-created_at']
        unique_together = ['integration', 'platform', 'platform_order_id']
```

---

## 2. Serializers para Shopify

### 2.1. Serializer de integración

```python
from rest_framework import serializers
from .models import ShopifyIntegration, Order
from .services import get_static_negative_factors


class ShopifyIntegrationSerializer(serializers.ModelSerializer):
    """Serializer para mostrar integraciones de Shopify"""
    class Meta:
        model = ShopifyIntegration
        fields = ['id', 'store_name', 'is_active', 'permissions', 'created_at', 'updated_at', 'last_sync']
        read_only_fields = ['id', 'created_at', 'updated_at', 'last_sync']
```

### 2.2. Serializer para conectar tienda

```python
class ShopifyConnectSerializer(serializers.Serializer):
    """Serializer para conectar una tienda de Shopify"""
    store_name = serializers.CharField(max_length=200, help_text='Nombre de la tienda (ej: sutienda.myshopify.com)')
    access_token = serializers.CharField(help_text='Token de acceso de la API de Shopify')
    
    def validate_store_name(self, value):
        """Validar formato del nombre de la tienda"""
        if not value.endswith('.myshopify.com'):
            raise serializers.ValidationError("El nombre de la tienda debe terminar en .myshopify.com")
        return value
    
    def validate_access_token(self, value):
        """Validar formato del token de acceso"""
        if not value.startswith('shp'):
            raise serializers.ValidationError("El token de acceso debe comenzar con 'shp'")
        return value
```

### 2.3. Serializer para probar conexión

```python
class ShopifyTestConnectionSerializer(serializers.Serializer):
    """Serializer para probar la conexión con Shopify"""
    message = serializers.CharField(read_only=True)
    success = serializers.BooleanField(read_only=True)
    store_info = serializers.DictField(read_only=True)
```

### 2.4. Serializer para órdenes

```python
class OrderSerializer(serializers.ModelSerializer):
    """Serializer ligero para exponer órdenes importadas"""
    negative_factors = serializers.SerializerMethodField()

    class Meta:
        model = Order
        fields = [
            'id', 'platform', 'platform_order_id', 'order_number', 'total_price', 'currency',
            'customer_name', 'customer_email', 'payment_type', 'created_at', 'imported_at',
            'source_name', 'financial_status', 'fulfillment_status', 'line_items',
            'address', 'city', 'country',
            'delivery_probability', 'recommended_carrier', 'delivery_score',
            'negative_factors', 'raw_response',
        ]
        read_only_fields = fields

    def get_negative_factors(self, obj):
        return get_static_negative_factors(obj)
```

---

## 3. Servicios: lógica de integración con Shopify

### 3.1. Mapeo de payload Shopify → modelo `Order`

```python
from .models import Order, ShopifyIntegration
import logging


def create_or_update_order(integration: ShopifyIntegration, order_data: dict, platform: str = 'shopify'):
    """Crear o actualizar una Order a partir del payload de la plataforma."""
    try:
        platform_order_id = str(
            order_data.get('id')
            or order_data.get('order_number')
            or order_data.get('name')
            or order_data.get('order_id')
            or ''
        )

        customer = order_data.get('customer') or {}
        shipping = order_data.get('shipping_address') or order_data.get('shipping') or {}

        def pick(*keys):
            """Busca una clave en order_data, shipping o customer con pequeñas variaciones de nombre."""
            def variants(k: str):
                if not k:
                    return []
                vlist = [k, k.lower()]
                if '_' in k:
                    vlist.append(k.replace('_', ' '))
                    vlist.append(k.replace('_', ' ').lower())
                if ' ' in k:
                    vlist.append(k.replace(' ', '_'))
                    vlist.append(k.replace(' ', '_').lower())
                return vlist

            for k in keys:
                if not k:
                    continue
                for cand in variants(k):
                    v = order_data.get(cand)
                    if v or v == 0:
                        return v
                    v = shipping.get(cand) if isinstance(shipping, dict) else None
                    if v or v == 0:
                        return v
                    v = customer.get(cand) if isinstance(customer, dict) else None
                    if v or v == 0:
                        return v
            return None

        first = pick('first_name', 'firstname', 'given_name')
        last = pick('last_name', 'lastname', 'family_name')
        name = pick('Shipping_Name', 'shipping_name')
        if not name and (first or last):
            name = ((first or '') + ' ' + (last or '')).strip() or None

        email = pick('email', 'Email', 'customer_email', 'shipping_email')
        phone = pick('phone', 'Phone', 'shipping_phone', 'Shipping_Phone') or order_data.get('phone') or customer.get('phone')

        country = pick('country', 'Shipping_Country', 'shipping_country')
        province = pick('province', 'province_code', 'Shipping_State', 'shipping_state')
        city = pick('city', 'Shipping_City', 'shipping_city')
        address1 = pick('address1', 'Shipping_Street', 'shipping_street', 'address')
        address2 = pick('address2', 'Shipping_Street2', 'shipping_address2')

        defaults = {
            'platform': platform,
            'order_number': (
                order_data.get('name')
                or order_data.get('order_number')
                or order_data.get('number')
            ),
            'total_price': (
                order_data.get('total_price')
                or order_data.get('order_total')
                or order_data.get('total')
            ),
            'currency': (
                order_data.get('currency')
                or order_data.get('currency_code')
                or order_data.get('Currency')
            ),
            'payment_type': (
                order_data.get('payment_gateway_names')[0]
                if order_data.get('payment_gateway_names')
                else None
            ) or order_data.get('payment_method'),
            'customer_name': name,
            'customer_email': email,
            'phone': phone,
            'customer_source': 'shipping' if shipping else ('customer' if customer else None),
            'country': country,
            'province': province,
            'city': city,
            'address': address1,
            'address_complement': address2,
            'source_name': order_data.get('source_name') or order_data.get('source'),
            'financial_status': order_data.get('financial_status'),
            'fulfillment_status': order_data.get('fulfillment_status'),
            'created_at': (
                order_data.get('created_at')
                or order_data.get('createdDate')
                or order_data.get('date')
            ),
            'customer_json': customer or {},
            'line_items': order_data.get('line_items') or order_data.get('items') or [],
            'raw_response': order_data,
        }

        obj, created = Order.objects.update_or_create(
            integration=integration,
            platform_order_id=platform_order_id,
            defaults=defaults,
        )

        if email:
            purchase_count = Order.objects.filter(
                integration=integration,
                customer_email=email,
            ).count()
            obj.purchase_count = purchase_count
            obj.save(update_fields=['purchase_count'])

        # En el proyecto original aquí también se calcula un "score" de entrega.
        # Puedes omitirlo o adaptarlo a tu nueva lógica.

        return obj
    except Exception as e:
        logging.exception('Error creando/actualizando orden: %s', e)
        return None
```

### 3.2. Obtener órdenes desde la API de Shopify

```python
import requests
import logging


def fetch_orders_for_integration(integration: ShopifyIntegration, created_at_min: str = None):
    """Trae órdenes desde Shopify para la integración indicada.

    - created_at_min: ISO8601 string para filtrar órdenes (ej: 2025-01-01T00:00:00)
    - Usa paginación por header Link (`rel="next"`).
    """
    access_token = integration.get_access_token()
    if not access_token:
        raise ValueError('No access token available for integration')

    shop_url = f"https://{integration.store_name}"
    base_url = f"{shop_url}/admin/api/2024-10/orders.json"
    params = {
        'status': 'any',
        'limit': 250,
    }
    if created_at_min:
        params['created_at_min'] = created_at_min

    headers = {
        'X-Shopify-Access-Token': access_token,
        'Content-Type': 'application/json',
    }

    resp = requests.get(base_url, headers=headers, params=params, timeout=30)
    if resp.status_code != 200:
        logging.error('Error fetching orders: %s', resp.status_code)
        return

    data = resp.json()
    orders = data.get('orders', [])
    for o in orders:
        create_or_update_order(integration, o, platform='shopify')

    # Si Shopify devuelve más páginas, seguir el Link "next"
    link_header = resp.headers.get('Link', '')
    while True:
        next_url = _parse_link_header(link_header)
        if not next_url:
            break
        try:
            resp = requests.get(next_url, headers=headers, timeout=30)
            if resp.status_code != 200:
                break
            data = resp.json()
            orders = data.get('orders', [])
            for o in orders:
                create_or_update_order(integration, o, platform='shopify')
            link_header = resp.headers.get('Link', '')
        except Exception:
            logging.exception('Error paginando órdenes de Shopify')
            break
```

> Nota: `_parse_link_header` es un helper sencillo que extrae la URL con `rel="next"` del header `Link`.

---

## 4. Tareas Celery

### 4.1. Importación inicial / sincronizaciones

```python
from config.celery import app
from django.utils import timezone
import logging
from .models import ShopifyIntegration
from . import services


@app.task(bind=True, max_retries=3)
def fetch_orders_initial_task(self, integration_id, created_at_min=None):
    try:
        integration = ShopifyIntegration.objects.get(id=integration_id)
        services.fetch_orders_for_integration(integration, created_at_min=created_at_min)
        integration.last_sync = timezone.now()
        integration.save()
    except Exception as exc:
        logging.exception('Error en fetch_orders_initial_task')
        try:
            raise self.retry(exc=exc, countdown=60)
        except Exception:
            logging.exception('Retry failed for fetch_orders_initial_task')
```

### 4.2. Procesar órdenes recibidas por Webhook

```python
@app.task(bind=True, max_retries=3)
def process_order_webhook_task(self, integration_id, order_payload):
    try:
        integration = ShopifyIntegration.objects.get(id=integration_id)
        services.create_or_update_order(integration, order_payload, platform='shopify')
    except Exception as exc:
        logging.exception('Error en process_order_webhook_task')
        try:
            raise self.retry(exc=exc, countdown=30)
        except Exception:
            logging.exception('Retry failed for process_order_webhook_task')
```

---

## 5. Vistas / Endpoints de la integración Shopify

### 5.1. Listar integraciones y conectar una nueva tienda

```python
from rest_framework.generics import GenericAPIView, ListAPIView, DestroyAPIView
from rest_framework.permissions import IsAuthenticated, AllowAny
from rest_framework.response import Response
from rest_framework import status, serializers
from django.shortcuts import get_object_or_404
from django.utils import timezone
from datetime import timedelta
import requests
import logging

from .models import ShopifyIntegration, Order
from .serializers import (
    ShopifyIntegrationSerializer,
    ShopifyConnectSerializer,
    ShopifyTestConnectionSerializer,
    OrderSerializer,
)
from . import tasks


class ShopifyIntegrationListView(ListAPIView):
    """Listar integraciones de Shopify del usuario autenticado."""
    serializer_class = ShopifyIntegrationSerializer
    permission_classes = [IsAuthenticated]
    
    def get_queryset(self):
        return ShopifyIntegration.objects.filter(user=self.request.user)


class ShopifyConnectView(GenericAPIView):
    """Conectar una tienda de Shopify (guarda integración y dispara importación inicial)."""
    serializer_class = ShopifyConnectSerializer
    permission_classes = [IsAuthenticated]
    
    def post(self, request):
        serializer = self.get_serializer(data=request.data)
        if serializer.is_valid():
            store_name = serializer.validated_data['store_name']
            access_token = serializer.validated_data['access_token']
            
            if ShopifyIntegration.objects.filter(user=request.user, store_name=store_name).exists():
                return Response({'error': 'Ya tienes una integración con esta tienda'}, status=status.HTTP_400_BAD_REQUEST)
            
            test_result = self._test_shopify_connection(store_name, access_token)
            if not test_result['success']:
                return Response({'error': test_result['message']}, status=status.HTTP_400_BAD_REQUEST)
            
            integration = ShopifyIntegration.objects.create(
                user=request.user,
                store_name=store_name,
                permissions=test_result.get('permissions', []),
            )
            integration.set_access_token(access_token)
            integration.save()

            # Importación inicial (~90 días)
            try:
                three_months_ago = timezone.now() - timedelta(days=90)
                tasks.fetch_orders_initial_task.delay(integration.id, three_months_ago.isoformat())
            except Exception:
                logging.exception('No se pudo encolar la tarea de importación inicial')

            return Response(
                {
                    'message': '¡Conexión con Shopify establecida correctamente!',
                    'integration': ShopifyIntegrationSerializer(integration).data,
                    'store_info': test_result.get('store_info', {}),
                },
                status=status.HTTP_201_CREATED,
            )
        
        return Response(serializer.errors, status=status.HTTP_400_BAD_REQUEST)
    
    def _test_shopify_connection(self, store_name, access_token):
        """Llamada simple a /shop.json para validar token y tienda."""
        try:
            shop_url = f"https://{store_name}"
            api_url = f"{shop_url}/admin/api/2023-10/shop.json"
            headers = {
                'X-Shopify-Access-Token': access_token,
                'Content-Type': 'application/json',
            }
            response = requests.get(api_url, headers=headers, timeout=10)
            
            if response.status_code == 200:
                shop_data = response.json()
                return {
                    'success': True,
                    'message': 'Conexión exitosa',
                    'store_info': shop_data.get('shop', {}),
                    'permissions': ['read_products', 'read_orders', 'read_customers', 'read_analytics'],
                }
            elif response.status_code == 401:
                return {
                    'success': False,
                    'message': 'Token de acceso inválido o sin permisos.',
                }
            elif response.status_code == 404:
                return {
                    'success': False,
                    'message': 'Tienda no encontrada. Verifica el nombre de la tienda.',
                }
            else:
                return {
                    'success': False,
                    'message': f'Error al conectar con Shopify: {response.status_code}',
                }
        except Exception as e:
            return {'success': False, 'message': f'Error de conexión: {str(e)}'}
```

### 5.2. Probar, desconectar, sincronizar y listar órdenes

```python
class ShopifyTestConnectionView(GenericAPIView):
    """Probar la conexión de una integración existente."""
    serializer_class = ShopifyTestConnectionSerializer
    permission_classes = [IsAuthenticated]
    
    def post(self, request, pk):
        integration = get_object_or_404(ShopifyIntegration, id=pk, user=request.user)
        access_token = integration.get_access_token()
        if not access_token:
            return Response({'success': False, 'message': 'No se pudo obtener el token de acceso'}, status=status.HTTP_400_BAD_REQUEST)
        
        test_result = ShopifyConnectView()._test_shopify_connection(integration.store_name, access_token)
        if test_result['success']:
            integration.last_sync = timezone.now()
            integration.save()
        return Response(test_result, status=status.HTTP_200_OK)


class ShopifyDisconnectView(DestroyAPIView):
    """Desconectar una integración de Shopify."""
    permission_classes = [IsAuthenticated]
    
    def get_queryset(self):
        return ShopifyIntegration.objects.filter(user=self.request.user)
    
    def destroy(self, request, *args, **kwargs):
        instance = self.get_object()
        store_name = instance.store_name
        instance.delete()
        return Response({'message': f'Integración con {store_name} desconectada exitosamente'}, status=status.HTTP_200_OK)


class ShopifySyncOrdersView(GenericAPIView):
    """Disparar sincronización manual de órdenes para una integración."""
    permission_classes = [IsAuthenticated]

    def post(self, request, pk):
        integration = get_object_or_404(ShopifyIntegration, id=pk, user=request.user)
        if integration.last_sync:
            created_at_min = integration.last_sync.isoformat()
        else:
            created_at_min = (timezone.now() - timedelta(days=90)).isoformat()

        try:
            tasks.fetch_orders_initial_task.delay(integration.id, created_at_min)
        except Exception:
            logging.exception('No se pudo encolar la tarea de sincronización manual')
            return Response({'error': 'No se pudo encolar la tarea de sincronización'}, status=status.HTTP_500_INTERNAL_SERVER_ERROR)

        return Response({'message': 'Sincronización encolada. Las órdenes nuevas serán importadas en breve.'}, status=status.HTTP_202_ACCEPTED)
```

```python
from rest_framework.pagination import PageNumberPagination
from django.db.models import Q


class DynamicPagination(PageNumberPagination):
    page_size = 20
    page_size_query_param = 'page_size'
    max_page_size = 100


class ShopifyIntegrationOrdersView(ListAPIView):
    """Listar órdenes importadas para una integración Shopify (con filtros y búsqueda)."""
    serializer_class = OrderSerializer
    permission_classes = [IsAuthenticated]
    pagination_class = DynamicPagination

    def get_queryset(self):
        pk = self.kwargs.get('pk')
        integration = get_object_or_404(ShopifyIntegration, id=pk, user=self.request.user)
        qs = integration.orders.all()

        today = self.request.query_params.get('today')
        if today and today.lower() in ('1', 'true', 'yes'):
            today_start = timezone.now().replace(hour=0, minute=0, second=0, microsecond=0)
            qs = qs.filter(imported_at__gte=today_start)

        min_total = self.request.query_params.get('min_total')
        if min_total:
            try:
                qs = qs.filter(total_price__gte=float(min_total))
            except Exception:
                pass

        since = self.request.query_params.get('since')
        if since:
            try:
                qs = qs.filter(created_at__gte=since)
            except Exception:
                pass

        search = self.request.query_params.get('search')
        if search:
            qs = qs.filter(
                Q(order_number__icontains=search)
                | Q(platform_order_id__icontains=search)
                | Q(customer_name__icontains=search)
                | Q(customer_email__icontains=search)
                | Q(address__icontains=search)
                | Q(city__icontains=search)
            )

        return qs.order_by('-created_at')
```

---

## 6. Webhooks de Shopify (opcional pero recomendado)

### 6.1. Endpoint público para recibir webhooks de órdenes

```python
from rest_framework.views import APIView
from rest_framework.decorators import permission_classes
from rest_framework.permissions import AllowAny
from django.utils.decorators import method_decorator
from django.views.decorators.csrf import csrf_exempt
import hmac
import hashlib
import base64
import json
from django.conf import settings


def _compute_hmac(secret: str, body: bytes) -> str:
    mac = hmac.new(secret.encode(), body, hashlib.sha256)
    return base64.b64encode(mac.digest()).decode()


@method_decorator(csrf_exempt, name='dispatch')
class ShopifyOrdersWebhookView(APIView):
    """Recibir webhooks de Shopify (orders/create, orders/updated)."""
    permission_classes = [AllowAny]

    def post(self, request, *args, **kwargs):
        raw_body = request.body or b''
        header_hmac = request.META.get('HTTP_X_SHOPIFY_HMAC_SHA256') or request.headers.get('X-Shopify-Hmac-Sha256')
        app_secret = getattr(settings, 'SHOPIFY_APP_SECRET', None)

        if app_secret:
            computed = _compute_hmac(app_secret, raw_body)
            if not header_hmac or not hmac.compare_digest(computed, header_hmac):
                logging.warning('Shopify webhook HMAC verification failed')
                return Response({'detail': 'HMAC verification failed'}, status=status.HTTP_401_UNAUTHORIZED)

        try:
            payload = json.loads(raw_body.decode('utf-8')) if raw_body else {}
        except Exception:
            payload = {}

        shop_domain = request.META.get('HTTP_X_SHOPIFY_SHOP_DOMAIN') or request.headers.get('X-Shopify-Shop-Domain')
        integration = None
        if shop_domain:
            integration = ShopifyIntegration.objects.filter(store_name=shop_domain).first()

        if not integration:
            return Response({'detail': 'Integration not found for shop'}, status=status.HTTP_404_NOT_FOUND)

        if isinstance(payload, dict) and payload.get('order'):
            tasks.process_order_webhook_task.delay(integration.id, payload.get('order'))
        elif isinstance(payload, dict) and payload.get('orders'):
            for o in payload.get('orders', []):
                tasks.process_order_webhook_task.delay(integration.id, o)
        elif isinstance(payload, dict) and payload.get('id'):
            tasks.process_order_webhook_task.delay(integration.id, payload)

        return Response({'status': 'ok'}, status=status.HTTP_200_OK)
```

### 6.2. Registrar webhook desde el backend

```python
class ShopifyRegisterWebhookView(GenericAPIView):
    """Registrar (o asegurar) webhooks en la tienda Shopify para una integración."""
    permission_classes = [IsAuthenticated]

    def post(self, request, pk):
        integration = get_object_or_404(ShopifyIntegration, id=pk, user=request.user)
        access_token = integration.get_access_token()
        if not access_token:
            return Response({'error': 'No se pudo obtener token de la integración'}, status=status.HTTP_400_BAD_REQUEST)

        public_base = request.data.get('public_base_url') or getattr(settings, 'PUBLIC_WEBHOOK_BASE_URL', None)
        if not public_base:
            return Response({'error': 'Se requiere public_base_url o PUBLIC_WEBHOOK_BASE_URL'}, status=status.HTTP_400_BAD_REQUEST)

        address = f"{public_base.rstrip('/')}/api/webhooks/shopify/orders/"

        shop_url = f"https://{integration.store_name}"
        webhooks_url = f"{shop_url}/admin/api/2024-10/webhooks.json"
        headers = {'X-Shopify-Access-Token': access_token, 'Content-Type': 'application/json'}

        resp = requests.get(webhooks_url, headers=headers, timeout=15)
        if resp.status_code != 200:
            return Response({'error': f'No se pudieron listar webhooks: {resp.status_code}'}, status=status.HTTP_400_BAD_REQUEST)

        existing = resp.json().get('webhooks', [])
        for w in existing:
            if w.get('address') == address and w.get('topic') == 'orders/create':
                return Response({'message': 'Webhook ya registrado', 'webhook': w}, status=status.HTTP_200_OK)

        payload = {'webhook': {'topic': 'orders/create', 'address': address, 'format': 'json'}}
        resp2 = requests.post(webhooks_url, headers=headers, json=payload, timeout=15)
        if resp2.status_code in (200, 201):
            return Response({'message': 'Webhook registrado', 'webhook': resp2.json().get('webhook', {})}, status=status.HTTP_201_CREATED)

        return Response({'error': 'Error al crear webhook', 'status_code': resp2.status_code}, status=status.HTTP_400_BAD_REQUEST)
```

---

## 7. URLs de la integración Shopify

Ejemplo de configuración de rutas (Django):

```python
from django.urls import path
from . import views

urlpatterns = [
    # Integraciones Shopify
    path('integrations/shopify/', views.ShopifyIntegrationListView.as_view(), name='shopify-integrations-list'),
    path('integrations/shopify/connect/', views.ShopifyConnectView.as_view(), name='shopify-connect'),
    path('integrations/shopify/<int:pk>/test/', views.ShopifyTestConnectionView.as_view(), name='shopify-test-connection'),
    path('integrations/shopify/<int:pk>/disconnect/', views.ShopifyDisconnectView.as_view(), name='shopify-disconnect'),
    path('integrations/shopify/<int:pk>/sync-orders/', views.ShopifySyncOrdersView.as_view(), name='shopify-sync-orders'),
    path('integrations/shopify/<int:pk>/orders/', views.ShopifyIntegrationOrdersView.as_view(), name='shopify-integration-orders'),

    # Webhooks
    path('webhooks/shopify/orders/', views.ShopifyOrdersWebhookView.as_view(), name='shopify-orders-webhook'),
]
```

---

## 8. Cómo explicarle esto a otra IA para migrarlo a otro proyecto

- **Objetivo**: “Quiero una integración con Shopify que permita conectar una tienda (con token privado), importar órdenes (inicial y por sincronización manual/webhooks) y exponerlas por API.”
- **Tecnologías actuales**: Django, Django REST Framework, Celery, PostgreSQL.
- **Piezas mínimas que debe replicar en el nuevo proyecto**:
  - Un modelo `ShopifyIntegration` por usuario con `store_name` y `access_token` (idealmente encriptado).
  - Un modelo `Order` (o similar) con los campos que quieras usar.
  - Lógica para llamar a `https://{store_name}/admin/api/<version>/orders.json` con `X-Shopify-Access-Token`.
  - Un endpoint para conectar la tienda (guardar token e iniciar importación inicial).
  - Un mecanismo para sincronizaciones periódicas o manuales (tarea background).
  - (Opcional pero útil) webhooks de `orders/create` para recibir nuevas órdenes en tiempo real.

Con este `.md` una IA puede leer rápidamente qué hace cada parte y ayudarte a portar la funcionalidad a otro framework o limpiarla para una nueva versión.


