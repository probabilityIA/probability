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
            })
            .catch(function () {});
    }

    function scheduleValidate() {
        if (validateTimer) clearTimeout(validateTimer);
        validateTimer = setTimeout(validate, 900);
    }

    function start() {
        if (!isBlockCheckout()) return;

        injectLogos();
        var observer = new MutationObserver(function () { injectLogos(); });
        observer.observe(document.body, { childList: true, subtree: true });

        if (window.wp && window.wp.data && window.wp.data.subscribe) {
            window.wp.data.subscribe(function () { scheduleValidate(); });
        }
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', start);
    } else {
        start();
    }
})();
