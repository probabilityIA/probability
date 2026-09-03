(function () {
    if (typeof ProbabilityCheckoutBlocks === 'undefined') {
        return;
    }

    var cfg = ProbabilityCheckoutBlocks;
    var apiBase = cfg.backendUrl + '/api/v1/woocommerce';
    var headers = { 'Content-Type': 'application/json', 'X-Probability-Token': cfg.token };

    function isBlockCheckout() {
        return !!document.querySelector('.wc-block-checkout, .wp-block-woocommerce-checkout');
    }

    function carrierFromLabel(text) {
        if (!text) return '';
        var parts = text.split(/\s[-–]\s/);
        return (parts[0] || '').trim();
    }

    function injectLogos() {
        var container = document.querySelector('.wc-block-components-shipping-rates-control');
        if (!container) return;
        var labels = container.querySelectorAll('.wc-block-components-radio-control__label');
        labels.forEach(function (label) {
            if (label.getAttribute('data-probability-logo')) return;
            var carrier = carrierFromLabel(label.textContent || '');
            if (!carrier) return;
            label.setAttribute('data-probability-logo', '1');
            var img = document.createElement('img');
            img.src = apiBase + '/carrier-logo/' + encodeURIComponent(carrier);
            img.alt = '';
            img.style.height = '18px';
            img.style.width = 'auto';
            img.style.verticalAlign = 'middle';
            img.style.marginRight = '8px';
            img.onerror = function () { if (img.parentNode) img.parentNode.removeChild(img); };
            label.insertBefore(img, label.firstChild);
        });
    }

    var lastValidated = '';
    var validateTimer = null;
    var pmap = null;
    var pmarker = null;

    function showMap(lat, lng) {
        if (!window.L || !lat || !lng) return;
        var note = document.getElementById('probability-blocks-note');
        var container = document.getElementById('probability-map');
        if (!container) {
            container = document.createElement('div');
            container.id = 'probability-map';
            var caption = document.createElement('div');
            caption.style.fontSize = '13px';
            caption.style.color = '#555';
            caption.style.margin = '8px 0 4px';
            caption.textContent = 'Confirma en el mapa que el punto de entrega es correcto';
            var mapEl = document.createElement('div');
            mapEl.id = 'probability-map-canvas';
            mapEl.style.height = '220px';
            mapEl.style.borderRadius = '8px';
            mapEl.style.overflow = 'hidden';
            container.appendChild(caption);
            container.appendChild(mapEl);
            var anchor = note || document.querySelector('.wc-block-components-shipping-rates-control') || document.querySelector('.wc-block-checkout');
            if (!anchor || !anchor.parentNode) return;
            anchor.parentNode.insertBefore(container, anchor.nextSibling);
        }
        var canvas = document.getElementById('probability-map-canvas');
        if (!pmap) {
            pmap = window.L.map(canvas).setView([lat, lng], 16);
            window.L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
                maxZoom: 19,
                attribution: '&copy; OpenStreetMap'
            }).addTo(pmap);
            pmarker = window.L.marker([lat, lng]).addTo(pmap);
        } else {
            pmap.setView([lat, lng], 16);
            pmarker.setLatLng([lat, lng]);
        }
        setTimeout(function () { if (pmap) pmap.invalidateSize(); }, 150);
    }

    function shippingAddress() {
        try {
            var store = window.wp && window.wp.data && window.wp.data.select('wc/store/cart');
            if (!store) return null;
            var data = store.getCustomerData ? store.getCustomerData() : null;
            if (!data) return null;
            return data.shippingAddress || data.billingAddress || null;
        } catch (e) {
            return null;
        }
    }

    function ensureNote() {
        var note = document.getElementById('probability-blocks-note');
        if (!note) {
            var anchor = document.querySelector('.wc-block-components-shipping-rates-control')
                || document.querySelector('.wc-block-checkout__shipping-fields')
                || document.querySelector('.wc-block-checkout');
            if (!anchor) return null;
            note = document.createElement('div');
            note.id = 'probability-blocks-note';
            note.style.fontSize = '13px';
            note.style.margin = '8px 0';
            anchor.parentNode.insertBefore(note, anchor);
        }
        return note;
    }

    function validate() {
        if (!cfg.showMap || cfg.showMap === 'false') return;

        var addr = shippingAddress();
        if (!addr) return;
        var address = addr.address_1 || '';
        var city = addr.city || '';
        var state = addr.state || '';
        if (address.length < 4 || city.length < 3) return;

        var key = address + '|' + city + '|' + state;
        if (key === lastValidated) return;
        lastValidated = key;

        fetch(apiBase + '/validate-address/' + cfg.integrationId, {
            method: 'POST',
            headers: headers,
            body: JSON.stringify({ address: address, city: city, state: state })
        })
            .then(function (r) { return r.ok ? r.json() : null; })
            .then(function (res) {
                if (!res) return;
                var note = ensureNote();
                if (!note) return;
                if (res.confidence === 'high') {
                    note.style.color = '#1a7f37';
                    note.textContent = 'Direccion de envio validada';
                } else if (res.confidence === 'medium') {
                    note.style.color = '#9a6700';
                    note.textContent = 'Direccion reconocida, verifica que sea correcta';
                } else {
                    note.style.color = '#b35900';
                    note.textContent = 'No pudimos validar la direccion, revisa ciudad y direccion';
                }

                if (res.found && res.lat && res.lng) {
                    showMap(res.lat, res.lng);
                }
            })
            .catch(function () {});
    }

    function scheduleValidate() {
        if (validateTimer) clearTimeout(validateTimer);
        validateTimer = setTimeout(validate, 900);
    }


    // ---- Selector de ciudad ----------------------------------------------
    // La ciudad deja de ser texto libre: se busca contra el listado DANE y el
    // comprador elige. Se manda el codigo, asi el backend no adivina el
    // municipio a partir del nombre.

    var cityBox = null;
    var cityTimer = null;
    var cityChosen = { name: '', code: '' };
    var lastCityQuery = '';

    function cityInput() {
        return document.getElementById('shipping-city') || document.getElementById('billing-city');
    }

    function setNativeValue(input, value) {
        var setter = Object.getOwnPropertyDescriptor(window.HTMLInputElement.prototype, 'value').set;
        setter.call(input, value);
        input.dispatchEvent(new Event('input', { bubbles: true }));
        input.dispatchEvent(new Event('change', { bubbles: true }));
    }

    function hideCityBox() {
        if (cityBox) cityBox.style.display = 'none';
    }

    function cityHint(input, text) {
        var hint = document.getElementById('probability-city-hint');
        if (!text) { if (hint) hint.remove(); return; }
        if (!hint) {
            hint = document.createElement('p');
            hint.id = 'probability-city-hint';
            hint.style.cssText = 'margin:4px 0 0;font-size:13px;color:#b35900';
            input.parentNode.appendChild(hint);
        }
        hint.textContent = text;
    }

    function ensureCityBox(input) {
        if (cityBox && cityBox.parentNode) return cityBox;
        cityBox = document.createElement('ul');
        cityBox.id = 'probability-city-options';
        cityBox.style.cssText = [
            'position:absolute', 'z-index:9999', 'left:0', 'right:0', 'margin:2px 0 0',
            'padding:0', 'list-style:none', 'background:#fff', 'color:#1a1a1a',
            'border:1px solid #ccc', 'border-radius:6px', 'max-height:220px',
            'overflow-y:auto', 'box-shadow:0 4px 12px rgba(0,0,0,.12)', 'display:none'
        ].join(';');
        var wrapper = input.parentNode;
        if (getComputedStyle(wrapper).position === 'static') wrapper.style.position = 'relative';
        wrapper.appendChild(cityBox);
        return cityBox;
    }

    function pushDane(code) {
        if (window.wc && window.wc.blocksCheckout && window.wc.blocksCheckout.extensionCartUpdate) {
            window.wc.blocksCheckout.extensionCartUpdate({
                namespace: 'probability_shipping',
                data: { dane_code: code }
            }).catch(function () {});
        }
    }

    function renderCities(input, cities) {
        var box = ensureCityBox(input);
        box.innerHTML = '';
        if (!cities || !cities.length) {
            var empty = document.createElement('li');
            empty.textContent = 'Sin resultados: revisa el nombre de la ciudad';
            empty.style.cssText = 'padding:8px 10px;font-size:13px;color:#8a6d09';
            box.appendChild(empty);
            box.style.display = 'block';
            cityHint(input, 'No encontramos esa ciudad. Escribe el nombre completo.');
            return;
        }
        cityHint(input, 'Selecciona tu ciudad de la lista para calcular el envio.');
        cities.forEach(function (city) {
            var li = document.createElement('li');
            li.style.cssText = 'padding:8px 10px;font-size:13px;cursor:pointer';
            li.textContent = city.name + (city.state ? '  -  ' + city.state : '');
            li.addEventListener('mouseenter', function () { li.style.background = '#f0f0f0'; });
            li.addEventListener('mouseleave', function () { li.style.background = ''; });
            li.addEventListener('mousedown', function (e) {
                e.preventDefault();
                cityChosen = { name: city.name, code: city.code };
                lastCityQuery = city.name;
                setNativeValue(input, city.name);
                pushDane(city.code);
                cityHint(input, '');
                hideCityBox();
            });
            box.appendChild(li);
        });
        box.style.display = 'block';
    }

    function searchCities() {
        var input = cityInput();
        if (!input) return;

        var term = (input.value || '').trim();
        if (term === cityChosen.name && cityChosen.code) { hideCityBox(); cityHint(input, ''); return; }

        if (cityChosen.code) { cityChosen = { name: '', code: '' }; pushDane(''); }
        if (term.length < 3) { hideCityBox(); cityHint(input, ''); return; }
        if (term === lastCityQuery) return;
        lastCityQuery = term;

        var addr = shippingAddress() || {};
        var qs = 'q=' + encodeURIComponent(term) + '&limit=8';
        if (addr.state) qs += '&state=' + encodeURIComponent(addr.state);

        fetch(apiBase + '/dane/' + cfg.integrationId + '/cities?' + qs, { headers: headers })
            .then(function (r) { return r.ok ? r.json() : null; })
            .then(function (res) { renderCities(input, res && res.cities); })
            .catch(function () {});
    }

    function bindCityInput() {
        var input = cityInput();
        if (!input || input.getAttribute('data-probability-city') === '1') return;
        input.setAttribute('data-probability-city', '1');
        input.setAttribute('autocomplete', 'off');
        input.addEventListener('input', function () {
            if (cityTimer) clearTimeout(cityTimer);
            cityTimer = setTimeout(searchCities, 250);
        });
        input.addEventListener('blur', function () { setTimeout(hideCityBox, 150); });
    }

    var lastCodSynced = null;

    function syncCodFlag() {
        try {
            var paymentStore = window.wp && window.wp.data && window.wp.data.select('wc/store/payment');
            if (!paymentStore || !paymentStore.getActivePaymentMethod) return;
            var active = paymentStore.getActivePaymentMethod();
            if (!active) return;
            var isCod = active === 'cod';
            if (isCod === lastCodSynced) return;
            lastCodSynced = isCod;
            if (window.wc && window.wc.blocksCheckout && window.wc.blocksCheckout.extensionCartUpdate) {
                window.wc.blocksCheckout.extensionCartUpdate({
                    namespace: 'probability_shipping',
                    data: { cod: isCod }
                }).catch(function () {});
            }
        } catch (e) {}
    }

    function start() {
        if (!isBlockCheckout()) return;

        injectLogos();
        bindCityInput();
        var observer = new MutationObserver(function () { injectLogos(); bindCityInput(); });
        observer.observe(document.body, { childList: true, subtree: true });

        if (window.wp && window.wp.data && window.wp.data.subscribe) {
            window.wp.data.subscribe(function () { scheduleValidate(); syncCodFlag(); });
        }
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', start);
    } else {
        start();
    }
})();
