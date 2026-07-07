(function () {
    if (typeof ProbabilityCheckout === 'undefined') {
        return;
    }

    var cfg = ProbabilityCheckout;
    var apiBase = cfg.backendUrl + '/api/v1/woocommerce';
    var headers = { 'Content-Type': 'application/json', 'X-Probability-Token': cfg.token };
    var validateTimer = null;
    var lastValidated = '';
    var lastDept = {};
    var pmap = null;
    var pmarker = null;

    function showMap(prefix, lat, lng) {
        if (!window.L || !lat || !lng) return;
        var note = document.getElementById('probability-note-' + prefix);
        var container = document.getElementById('probability-map');
        if (!container) {
            container = document.createElement('div');
            container.id = 'probability-map';
            var caption = document.createElement('div');
            caption.style.fontSize = '12px';
            caption.style.color = '#555';
            caption.style.margin = '8px 0 4px';
            caption.textContent = 'Confirma en el mapa que el punto de entrega es correcto';
            var mapEl = document.createElement('div');
            mapEl.id = 'probability-map-canvas';
            mapEl.style.height = '200px';
            mapEl.style.borderRadius = '8px';
            mapEl.style.overflow = 'hidden';
            container.appendChild(caption);
            container.appendChild(mapEl);
            var host = note && note.parentNode ? note.parentNode : (cityField(prefix) ? cityField(prefix).parentNode : null);
            if (!host) return;
            if (note && note.nextSibling) host.insertBefore(container, note.nextSibling);
            else host.appendChild(container);
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

    function prefixes() {
        var out = [];
        var shipToDiff = document.getElementById('ship-to-different-address-checkbox');
        if (document.getElementById('shipping_city') && shipToDiff && shipToDiff.checked) {
            out.push('shipping');
        } else {
            out.push('billing');
        }
        return out;
    }

    function stateName(prefix) {
        var sel = document.getElementById(prefix + '_state');
        if (!sel) return '';
        if (sel.tagName === 'SELECT' && sel.selectedIndex >= 0) {
            return sel.options[sel.selectedIndex].text || sel.value || '';
        }
        return sel.value || '';
    }

    function cityField(prefix) {
        return document.getElementById(prefix + '_city');
    }

    function enhance(el, prefix) {
        var $ = window.jQuery;
        if (!$ || !$.fn) return;
        try {
            if ($.fn.selectWoo) {
                $(el).selectWoo({ width: '100%', placeholder: 'Selecciona tu ciudad', allowClear: false });
            } else if ($.fn.select2) {
                $(el).select2({ width: '100%', placeholder: 'Selecciona tu ciudad' });
            }
        } catch (e) {}
        $(el).on('change', function () { scheduleValidate(prefix); });
    }

    function destroyEnhance(el) {
        var $ = window.jQuery;
        if (!$ || !$.fn) return;
        try {
            if ($.fn.selectWoo && $(el).data('select2')) {
                $(el).selectWoo('destroy');
            } else if ($.fn.select2 && $(el).data('select2')) {
                $(el).select2('destroy');
            }
        } catch (e) {}
    }

    function buildOptions(select, cities, currentVal) {
        select.innerHTML = '';
        var empty = document.createElement('option');
        empty.value = '';
        empty.textContent = 'Selecciona tu ciudad';
        select.appendChild(empty);
        cities.forEach(function (c) {
            var opt = document.createElement('option');
            opt.value = c.name;
            opt.textContent = c.name;
            if (currentVal && c.name.toLowerCase() === currentVal.toLowerCase()) {
                opt.selected = true;
            }
            select.appendChild(opt);
        });
    }

    function loadCities(prefix, force) {
        var field = cityField(prefix);
        if (!field) return;
        var dept = stateName(prefix);
        if (!dept) return;
        if (!force && lastDept[prefix] === dept) return;
        lastDept[prefix] = dept;

        var currentVal = field.value || '';

        fetch(apiBase + '/dane/' + cfg.integrationId + '/cities?state=' + encodeURIComponent(dept), { headers: headers })
            .then(function (r) { return r.ok ? r.json() : null; })
            .then(function (data) {
                if (!data || !data.cities) return;
                var el = cityField(prefix);
                if (!el) return;

                if (el.tagName === 'SELECT' && el.getAttribute('data-probability-city')) {
                    destroyEnhance(el);
                    buildOptions(el, data.cities, currentVal);
                    enhance(el, prefix);
                    return;
                }

                var select = document.createElement('select');
                select.id = el.id;
                select.name = el.name;
                select.className = el.className;
                select.setAttribute('data-probability-city', '1');
                if (el.required) select.required = true;
                buildOptions(select, data.cities, currentVal);
                el.parentNode.replaceChild(select, el);
                enhance(select, prefix);
            })
            .catch(function () {});
    }

    function ensureHidden(name) {
        var el = document.querySelector('input[name="' + name + '"]');
        if (!el) {
            var form = document.querySelector('form.checkout') || document.querySelector('form[name="checkout"]');
            if (!form) return null;
            el = document.createElement('input');
            el.type = 'hidden';
            el.name = name;
            form.appendChild(el);
        }
        return el;
    }

    function ensureNote(prefix) {
        var field = cityField(prefix);
        if (!field) return null;
        var host = field.closest('.form-row') || field.parentNode;
        var id = 'probability-note-' + prefix;
        var note = document.getElementById(id);
        if (!note) {
            note = document.createElement('div');
            note.id = id;
            note.style.fontSize = '12px';
            note.style.marginTop = '4px';
            host.appendChild(note);
        }
        return note;
    }

    function validate(prefix) {
        var field = cityField(prefix);
        var addrInput = document.getElementById(prefix + '_address_1');
        if (!field || !addrInput) return;

        var address = addrInput.value || '';
        var city = field.value || '';
        var state = stateName(prefix);
        if (address.length < 4 || city.length < 3) return;

        var key = prefix + '|' + address + '|' + city + '|' + state;
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
                var note = ensureNote(prefix);
                if (!note) return;
                var dane = ensureHidden('probability_dane_code');
                var lat = ensureHidden('probability_lat');
                var lng = ensureHidden('probability_lng');
                if (dane) dane.value = res.dane_code || '';
                if (lat) lat.value = res.lat || '';
                if (lng) lng.value = res.lng || '';

                if (res.confidence === 'high') {
                    note.style.color = '#1a7f37';
                    note.textContent = 'Direccion validada';
                } else if (res.confidence === 'medium') {
                    note.style.color = '#9a6700';
                    note.textContent = 'Direccion reconocida, verifica que sea correcta';
                } else {
                    note.style.color = '#b35900';
                    note.textContent = 'No pudimos validar la direccion, revisa ciudad y direccion';
                }

                if (res.found && res.lat && res.lng) {
                    showMap(prefix, res.lat, res.lng);
                }
            })
            .catch(function () {});
    }

    function scheduleValidate(prefix) {
        if (validateTimer) clearTimeout(validateTimer);
        validateTimer = setTimeout(function () { validate(prefix); }, 700);
    }

    function wire() {
        prefixes().forEach(function (prefix) {
            var stateSel = document.getElementById(prefix + '_state');
            var addrInput = document.getElementById(prefix + '_address_1');
            if (!cityField(prefix)) return;

            loadCities(prefix, false);

            if (stateSel && !stateSel.getAttribute('data-probability')) {
                stateSel.setAttribute('data-probability', '1');
                if (window.jQuery) {
                    window.jQuery(stateSel).on('change', function () { loadCities(prefix, true); scheduleValidate(prefix); });
                } else {
                    stateSel.addEventListener('change', function () { loadCities(prefix, true); scheduleValidate(prefix); });
                }
            }
            if (addrInput && !addrInput.getAttribute('data-probability')) {
                addrInput.setAttribute('data-probability', '1');
                addrInput.addEventListener('blur', function () { scheduleValidate(prefix); });
            }
        });
    }

    if (window.jQuery) {
        window.jQuery(document.body).on('updated_checkout country_to_state_changed', function () { wire(); });
    }
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', wire);
    } else {
        wire();
    }
})();
