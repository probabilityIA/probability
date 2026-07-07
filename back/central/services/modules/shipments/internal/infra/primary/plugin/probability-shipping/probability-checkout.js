(function () {
    if (typeof ProbabilityCheckout === 'undefined') {
        return;
    }

    var cfg = ProbabilityCheckout;
    var apiBase = cfg.backendUrl + '/api/v1/woocommerce';
    var headers = { 'Content-Type': 'application/json', 'X-Probability-Token': cfg.token };
    var validateTimer = null;
    var lastValidated = '';

    function prefixes() {
        var out = [];
        if (document.getElementById('shipping_city') &&
            document.getElementById('ship-to-different-address-checkbox') &&
            document.getElementById('ship-to-different-address-checkbox').checked) {
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

    function ensureDatalist(input) {
        var id = 'probability-cities-' + input.id;
        var dl = document.getElementById(id);
        if (!dl) {
            dl = document.createElement('datalist');
            dl.id = id;
            document.body.appendChild(dl);
            input.setAttribute('list', id);
        }
        return dl;
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

    function ensureNote(input) {
        var id = 'probability-note-' + input.id;
        var note = document.getElementById(id);
        if (!note) {
            note = document.createElement('div');
            note.id = id;
            note.style.fontSize = '12px';
            note.style.marginTop = '4px';
            if (input.parentNode) input.parentNode.appendChild(note);
        }
        return note;
    }

    function loadCities(prefix) {
        var cityInput = document.getElementById(prefix + '_city');
        if (!cityInput) return;
        var dept = stateName(prefix);
        if (!dept) return;

        fetch(apiBase + '/dane/' + cfg.integrationId + '/cities?state=' + encodeURIComponent(dept), { headers: headers })
            .then(function (r) { return r.ok ? r.json() : null; })
            .then(function (data) {
                if (!data || !data.cities) return;
                var dl = ensureDatalist(cityInput);
                dl.innerHTML = '';
                data.cities.forEach(function (c) {
                    var opt = document.createElement('option');
                    opt.value = c.name;
                    dl.appendChild(opt);
                });
            })
            .catch(function () {});
    }

    function validate(prefix) {
        var cityInput = document.getElementById(prefix + '_city');
        var addrInput = document.getElementById(prefix + '_address_1');
        if (!cityInput || !addrInput) return;

        var address = addrInput.value || '';
        var city = cityInput.value || '';
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
                var note = ensureNote(cityInput);
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
            var cityInput = document.getElementById(prefix + '_city');
            var addrInput = document.getElementById(prefix + '_address_1');
            if (!cityInput) return;

            loadCities(prefix);

            if (stateSel && !stateSel.getAttribute('data-probability')) {
                stateSel.setAttribute('data-probability', '1');
                stateSel.addEventListener('change', function () { loadCities(prefix); scheduleValidate(prefix); });
            }
            if (!cityInput.getAttribute('data-probability')) {
                cityInput.setAttribute('data-probability', '1');
                cityInput.addEventListener('change', function () { scheduleValidate(prefix); });
                cityInput.addEventListener('blur', function () { scheduleValidate(prefix); });
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
