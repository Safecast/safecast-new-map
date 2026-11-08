/* Safecast kriging overlay. The grid comes from /api/kriging/grid.
   Suppressed cells are not drawn: a blank area is unsupported, not safe.
   Markers stay above this pane. */
(function (global) {
  'use strict';

  var enabled = false;
  var mapRef = null;
  var optsRef = null;
  var canvas = null;
  var grid = null;
  var abort = null;
  var timer = 0;
  var statusEl = null;

  function text(key) {
    if (optsRef && typeof optsRef.translate === 'function') {
      return optsRef.translate(key);
    }
    return key;
  }

  function setStatus(message) {
    if (!statusEl) return;
    if (!message) {
      statusEl.style.display = 'none';
      statusEl.textContent = '';
      return;
    }
    statusEl.style.display = 'block';
    statusEl.textContent = message;
  }

  function ensurePane(map) {
    if (map.getPane('kriging')) return;
    var pane = map.createPane('kriging');
    pane.style.zIndex = '250';
    pane.style.pointerEvents = 'none';
  }

  function ensureCanvas(map) {
    ensurePane(map);
    if (canvas) return;
    canvas = L.DomUtil.create('canvas', 'leaflet-kriging-layer');
    canvas.style.pointerEvents = 'none';
    map.getPane('kriging').appendChild(canvas);
  }

  function clearCanvas() {
    grid = null;
    if (!canvas || !mapRef) return;
    var size = mapRef.getSize();
    var ctx = canvas.getContext('2d');
    ctx.setTransform(1, 0, 0, 1, 0, 0);
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    canvas.style.width = size.x + 'px';
    canvas.style.height = size.y + 'px';
  }

  function rgbOf(color) {
    var match = /rgb\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*\)/.exec(color || '');
    if (!match) return null;
    return [Number(match[1]), Number(match[2]), Number(match[3])];
  }

  function halfSpans(data) {
    var halfLat = Math.abs(data.maxLat - data.minLat) / data.rows / 2;
    var halfLon = Math.abs(data.maxLon - data.minLon) / data.cols / 2;
    if (data.rows > 1 && data.cells.length > data.cols) {
      halfLat = Math.abs(data.cells[data.cols].lat - data.cells[0].lat) / 2;
    }
    if (data.cols > 1) {
      var dLon = data.cells[1].lon - data.cells[0].lon;
      if (dLon > 180) dLon -= 360;
      if (dLon < -180) dLon += 360;
      halfLon = Math.abs(dLon) / 2;
    }
    return { halfLat: halfLat, halfLon: halfLon };
  }

  function draw() {
    if (!enabled || !mapRef || !grid || !grid.cells) return;
    ensureCanvas(mapRef);
    var size = mapRef.getSize();
    var dpr = global.devicePixelRatio || 1;
    var topLeft = mapRef.containerPointToLayerPoint([0, 0]);
    L.DomUtil.setPosition(canvas, topLeft);
    canvas.width = Math.max(1, Math.round(size.x * dpr));
    canvas.height = Math.max(1, Math.round(size.y * dpr));
    canvas.style.width = size.x + 'px';
    canvas.style.height = size.y + 'px';
    var ctx = canvas.getContext('2d');
    ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
    ctx.clearRect(0, 0, size.x, size.y);

    var span = halfSpans(grid);
    var sill = grid.partialSill || 1;
    for (var i = 0; i < grid.cells.length; i++) {
      var cell = grid.cells[i];
      if (!cell || cell.suppressed || cell.medianUSvH == null) continue;
      var alpha = 1 - cell.logVariance / sill;
      if (alpha <= 0) continue;
      if (alpha > 1) alpha = 1;
      var color = optsRef.color ? optsRef.color(cell.medianUSvH) : null;
      var rgb = rgbOf(color);
      if (!rgb) continue;
      var sw = mapRef.latLngToContainerPoint([cell.lat - span.halfLat, cell.lon - span.halfLon]);
      var ne = mapRef.latLngToContainerPoint([cell.lat + span.halfLat, cell.lon + span.halfLon]);
      var x = Math.min(sw.x, ne.x);
      var y = Math.min(sw.y, ne.y);
      ctx.globalAlpha = alpha;
      ctx.fillStyle = 'rgb(' + rgb[0] + ',' + rgb[1] + ',' + rgb[2] + ')';
      ctx.fillRect(x, y, Math.abs(ne.x - sw.x) + 1, Math.abs(ne.y - sw.y) + 1);
    }
    ctx.globalAlpha = 1;
  }

  function predictionTime() {
    if (optsRef && typeof optsRef.getPredictionTime === 'function') {
      var t = optsRef.getPredictionTime();
      if (isFinite(t)) return Math.round(t);
    }
    return Math.floor(Date.now() / 1000);
  }

  function load() {
    if (!enabled || !mapRef) return;
    if (abort) abort.abort();
    abort = new AbortController();
    var signal = abort.signal;
    var bounds = mapRef.getBounds();
    var size = mapRef.getSize();
    var params = new URLSearchParams();
    params.set('minLat', String(bounds.getSouth()));
    params.set('minLon', String(bounds.getWest()));
    params.set('maxLat', String(bounds.getNorth()));
    params.set('maxLon', String(bounds.getEast()));
    params.set('time', String(predictionTime()));
    params.set('cols', String(Math.max(8, Math.min(48, Math.round(size.x / 16)))));
    params.set('rows', String(Math.max(8, Math.min(48, Math.round(size.y / 16)))));

    fetch('/api/kriging/grid?' + params.toString(), { signal: signal })
      .then(function (res) {
        return res.json().then(function (data) {
          return { status: res.status, data: data };
        });
      })
      .then(function (msg) {
        if (signal.aborted || !enabled) return;
        var data = msg.data || {};
        if (data.refusal) {
          clearCanvas();
          setStatus(text('kriging_zoom_in'));
          return;
        }
        if (msg.status !== 200 || !data.cells) {
          clearCanvas();
          setStatus(text('kriging_unavailable'));
          return;
        }
        setStatus('');
        grid = data;
        draw();
      })
      .catch(function (err) {
        if (err && err.name === 'AbortError') return;
        if (!enabled) return;
        clearCanvas();
        setStatus(text('kriging_unavailable'));
      });
  }

  function schedule() {
    if (!enabled) return;
    global.clearTimeout(timer);
    timer = global.setTimeout(load, 200);
  }

  function onZoomStart() {
    clearCanvas();
  }

  function setEnabled(on) {
    enabled = on;
    if (!mapRef) return;
    if (!on) {
      if (abort) abort.abort();
      global.clearTimeout(timer);
      clearCanvas();
      if (canvas && canvas.parentNode) canvas.parentNode.removeChild(canvas);
      canvas = null;
      setStatus('');
      return;
    }
    ensureCanvas(mapRef);
    schedule();
  }

  function attach(map, opts) {
    mapRef = map;
    optsRef = opts || {};
    var Ctrl = L.Control.extend({
      options: { position: 'topleft' },
      onAdd: function () {
        var div = L.DomUtil.create('div', 'leaflet-control-layers');
        div.id = 'kriging-layer-ctrl';
        div.style.padding = '6px 10px';
        var label = text('kriging_layer');
        var tip = text('kriging_layer_tooltip');
        div.innerHTML =
          '<label style="display:block;cursor:pointer;white-space:nowrap;">' +
          '<input id="kriging-layer-toggle" type="checkbox"/> ' +
          '<span></span></label>' +
          '<div id="kriging-status" style="display:none;max-width:200px;white-space:normal;font-size:11px;margin-top:4px;"></div>';
        div.querySelector('span').textContent = label;
        div.title = tip;
        statusEl = div.querySelector('#kriging-status');
        L.DomEvent.disableClickPropagation(div);
        L.DomEvent.disableScrollPropagation(div);
        div.querySelector('input').addEventListener('change', function (ev) {
          setEnabled(ev.target.checked);
        });
        return div;
      }
    });
    map.addControl(new Ctrl());
    map.on('moveend zoomend resize', schedule);
    map.on('zoomstart', onZoomStart);
  }

  global.SafecastKriging = {
    attach: attach,
    refresh: function () {
      if (enabled) schedule();
    }
  };
})(window);
