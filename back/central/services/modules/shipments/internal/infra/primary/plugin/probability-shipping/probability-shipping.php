<?php
/**
 * Plugin Name: Probability Shipping
 * Description: Cotiza tarifas de transportadoras (EnvioClick, etc.) en el checkout consultando la API de Probability.
 * Version: 1.1.0
 * Author: Probability
 * Requires Plugins: woocommerce
 */

if (!defined('ABSPATH')) {
    exit;
}

add_action('woocommerce_shipping_init', function () {
    if (class_exists('Probability_Shipping_Method')) {
        return;
    }

    class Probability_Shipping_Method extends WC_Shipping_Method {

        public function __construct($instance_id = 0) {
            $this->id                 = 'probability_shipping';
            $this->instance_id        = absint($instance_id);
            $this->method_title       = 'Probability (Transportadoras)';
            $this->method_description = 'Tarifas en tiempo real desde la API de Probability.';
            $this->supports           = array('shipping-zones', 'instance-settings', 'settings');

            $this->init();
        }

        public function init() {
            $this->init_form_fields();
            $this->init_settings();

            $this->enabled = $this->get_option('enabled', 'yes');
            $this->title   = $this->get_option('title', 'Envio');

            add_action('woocommerce_update_options_shipping_' . $this->id, array($this, 'process_admin_options'));
        }

        public function init_form_fields() {
            $this->instance_form_fields = array(
                'enabled' => array(
                    'title'   => 'Habilitar',
                    'type'    => 'checkbox',
                    'label'   => 'Habilitar cotizacion con Probability',
                    'default' => 'yes',
                ),
                'title' => array(
                    'title'       => 'Titulo',
                    'type'        => 'text',
                    'description' => 'Nombre del metodo que ve el cliente.',
                    'default'     => 'Envio',
                    'desc_tip'    => true,
                ),
                'connection_key' => array(
                    'title'       => 'Clave de conexion',
                    'type'        => 'textarea',
                    'description' => 'Pega aqui la Clave de conexion que aparece en Probability (Integraciones -> WooCommerce). Con esto queda todo configurado.',
                    'default'     => '',
                    'css'         => 'height: 70px;',
                ),
                'fallback_cost' => array(
                    'title'       => 'Costo de respaldo',
                    'type'        => 'text',
                    'description' => 'Si la API no devuelve tarifas, usar este costo (vacio = sin opcion de envio).',
                    'default'     => '',
                    'desc_tip'    => true,
                ),
                'backend_url' => array(
                    'title'       => 'URL del backend (avanzado)',
                    'type'        => 'text',
                    'description' => 'Solo si no usas la Clave de conexion.',
                    'default'     => '',
                    'desc_tip'    => true,
                ),
                'integration_id' => array(
                    'title'       => 'Integration ID (avanzado)',
                    'type'        => 'text',
                    'description' => 'Solo si no usas la Clave de conexion.',
                    'default'     => '',
                    'desc_tip'    => true,
                ),
                'token' => array(
                    'title'       => 'Token (avanzado)',
                    'type'        => 'text',
                    'description' => 'Solo si no usas la Clave de conexion.',
                    'default'     => '',
                    'desc_tip'    => true,
                ),
            );
        }

        private function b64url_decode($data) {
            $pad = strlen($data) % 4;
            if ($pad > 0) {
                $data .= str_repeat('=', 4 - $pad);
            }
            return base64_decode(strtr($data, '-_', '+/'));
        }

        private function resolve_config() {
            $cfg = array('url' => '', 'integration_id' => '', 'token' => '');

            $key = trim($this->get_option('connection_key'));
            if ($key !== '') {
                $decoded = json_decode($this->b64url_decode($key), true);
                if (is_array($decoded)) {
                    $cfg['url']            = isset($decoded['url']) ? $decoded['url'] : '';
                    $cfg['integration_id'] = isset($decoded['integration_id']) ? (string) $decoded['integration_id'] : '';
                    $cfg['token']          = isset($decoded['token']) ? $decoded['token'] : '';
                }
            }

            if ($cfg['url'] === '') {
                $cfg['url'] = trim($this->get_option('backend_url'));
            }
            if ($cfg['integration_id'] === '') {
                $cfg['integration_id'] = trim($this->get_option('integration_id'));
            }
            if ($cfg['token'] === '') {
                $cfg['token'] = trim($this->get_option('token'));
            }

            return $cfg;
        }

        public function calculate_shipping($package = array()) {
            $cfg = $this->resolve_config();

            if (empty($cfg['url']) || empty($cfg['integration_id']) || empty($cfg['token'])) {
                $this->maybe_fallback();
                return;
            }

            $endpoint = rtrim($cfg['url'], '/') . '/api/v1/woocommerce/shipping-rates/' . rawurlencode($cfg['integration_id']);

            $body = $this->build_request_body($package);

            $response = wp_remote_post($endpoint, array(
                'timeout' => 15,
                'headers' => array(
                    'Content-Type'        => 'application/json',
                    'X-Probability-Token' => $cfg['token'],
                ),
                'body' => wp_json_encode($body),
            ));

            if (is_wp_error($response)) {
                $this->maybe_fallback();
                return;
            }

            $code = wp_remote_retrieve_response_code($response);
            if ($code !== 200) {
                $this->maybe_fallback();
                return;
            }

            $data  = json_decode(wp_remote_retrieve_body($response), true);
            $rates = isset($data['rates']) && is_array($data['rates']) ? $data['rates'] : array();

            if (empty($rates)) {
                $this->maybe_fallback();
                return;
            }

            foreach ($rates as $rate) {
                if (empty($rate['id']) || !isset($rate['cost'])) {
                    continue;
                }

                $label = isset($rate['label']) ? $rate['label'] : $this->title;
                if (!empty($rate['delivery_days'])) {
                    $label .= ' (' . intval($rate['delivery_days']) . ' dias habiles)';
                }

                $this->add_rate(array(
                    'id'        => $this->id . ':' . sanitize_title($rate['id']),
                    'label'     => $label,
                    'cost'      => floatval($rate['cost']),
                    'meta_data' => isset($rate['meta_data']) ? $rate['meta_data'] : array(),
                ));
            }
        }

        private function maybe_fallback() {
            $fallback = trim($this->get_option('fallback_cost'));
            if ($fallback === '') {
                return;
            }
            $this->add_rate(array(
                'id'    => $this->id . ':fallback',
                'label' => $this->title,
                'cost'  => floatval($fallback),
            ));
        }

        private function build_request_body($package) {
            $dest       = isset($package['destination']) ? $package['destination'] : array();
            $country    = isset($dest['country']) ? $dest['country'] : '';
            $state_code = isset($dest['state']) ? $dest['state'] : '';
            $state_name = $this->resolve_state_name($country, $state_code);

            $customer = WC()->customer;
            $name     = '';
            $phone    = '';
            $email    = '';
            if ($customer) {
                $name = trim($customer->get_shipping_first_name() . ' ' . $customer->get_shipping_last_name());
                if ($name === '') {
                    $name = trim($customer->get_billing_first_name() . ' ' . $customer->get_billing_last_name());
                }
                $phone = $customer->get_billing_phone();
                $email = $customer->get_billing_email();
            }

            $contents = array();
            if (!empty($package['contents'])) {
                foreach ($package['contents'] as $item) {
                    $product = isset($item['data']) ? $item['data'] : null;
                    if (!$product) {
                        continue;
                    }
                    $qty        = isset($item['quantity']) ? intval($item['quantity']) : 1;
                    $weight     = $product->get_weight();
                    $weight_g   = $weight !== '' ? floatval(wc_get_weight($weight, 'g')) : 0;
                    $unit_price = floatval($product->get_price());

                    $contents[] = array(
                        'name'         => $product->get_name(),
                        'sku'          => $product->get_sku(),
                        'quantity'     => $qty,
                        'weight_grams' => $weight_g,
                        'price'        => $unit_price,
                    );
                }
            }

            return array(
                'destination' => array(
                    'country'   => $country,
                    'state'     => $state_name,
                    'city'      => isset($dest['city']) ? $dest['city'] : '',
                    'postcode'  => isset($dest['postcode']) ? $dest['postcode'] : '',
                    'address_1' => isset($dest['address']) ? $dest['address'] : '',
                    'address_2' => isset($dest['address_2']) ? $dest['address_2'] : '',
                    'name'      => $name,
                    'phone'     => $phone,
                    'email'     => $email,
                ),
                'contents' => $contents,
                'currency' => get_woocommerce_currency(),
            );
        }

        private function resolve_state_name($country, $state_code) {
            if (empty($state_code) || empty($country)) {
                return $state_code;
            }
            $states = WC()->countries->get_states($country);
            if (is_array($states) && isset($states[$state_code])) {
                return $states[$state_code];
            }
            return $state_code;
        }
    }
});

add_filter('woocommerce_shipping_methods', function ($methods) {
    $methods['probability_shipping'] = 'Probability_Shipping_Method';
    return $methods;
});
