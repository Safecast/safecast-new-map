// Preview-only controls for the REFERENCE object in reference-basemap.js.
// To remove this panel, delete this file and the reference-tune.js script tag
// in map.html. Nothing else refers to it, and the map keeps the file defaults.
(function () {
  var STORE = 'safecastReferenceTune';

  var EARTH = [
    ['land', 'Land'],
    ['water', 'Water fill'],
    ['park', 'Park'],
    ['building', 'Building']
  ];
  var COLORS = [
    ['roadCasing', 'Road casing'],
    ['roadInner', 'Road center'],
    ['minorRoad', 'Minor road'],
    ['path', 'Path'],
    ['rail', 'Rail'],
    ['railDash', 'Rail dash'],
    ['boundary', 'Boundary'],
    ['water', 'Water line']
  ];
  var STROKES = [
    ['road', 'Road width'],
    ['path', 'Path width'],
    ['rail', 'Rail width'],
    ['boundary', 'Border width'],
    ['river', 'River width'],
    ['coast', 'Coast width']
  ];
  var WIDTHS = [
    ['road', 'Road halo'],
    ['path', 'Path halo'],
    ['boundary', 'Border halo'],
    ['river', 'River halo'],
    ['coast', 'Coast halo']
  ];
  var GROUPS = [
    ['roads', 'Roads'],
    ['rails', 'Rails'],
    ['boundaries', 'Borders'],
    ['coast', 'Coast'],
    ['rivers', 'Rivers'],
    ['labels', 'Labels']
  ];

  function savedCopy() {
    try {
      var raw = sessionStorage.getItem(STORE);
      return raw ? JSON.parse(raw) : null;
    } catch (err) {
      return null;
    }
  }

  function remember(ref) {
    try { sessionStorage.setItem(STORE, JSON.stringify(ref)); } catch (err) {}
  }

  function forget() {
    try { sessionStorage.removeItem(STORE); } catch (err) {}
  }

  function channel(n) {
    var value = Math.max(0, Math.min(255, Math.round(+n)));
    return ('0' + value.toString(16)).slice(-2);
  }

  function hexOf(color) {
    var text = String(color || '').trim();
    var short = text.match(/^#([0-9a-fA-F]{3})$/);
    if (short) {
      return '#' + short[1].replace(/(.)/g, '$1$1');
    }
    if (/^#[0-9a-fA-F]{6}$/.test(text)) return text.toLowerCase();
    var rgb = text.match(/rgba?\(\s*([\d.]+)[ ,]+([\d.]+)[ ,]+([\d.]+)/i);
    if (rgb) return '#' + channel(rgb[1]) + channel(rgb[2]) + channel(rgb[3]);
    var hsl = text.match(/hsla?\(\s*([\d.]+)\s*,\s*([\d.]+)%\s*,\s*([\d.]+)%/i);
    if (!hsl) return '#f2f3f0';
    var h = +hsl[1];
    var s = +hsl[2] / 100;
    var l = +hsl[3] / 100;
    var c = (1 - Math.abs(2 * l - 1)) * s;
    var hp = h / 60;
    var x = c * (1 - Math.abs(hp % 2 - 1));
    var r = 0;
    var g = 0;
    var b = 0;
    if (hp < 1) { r = c; g = x; }
    else if (hp < 2) { r = x; g = c; }
    else if (hp < 3) { g = c; b = x; }
    else if (hp < 4) { g = x; b = c; }
    else if (hp < 5) { r = x; b = c; }
    else { r = c; b = x; }
    var m = l - c / 2;
    return '#' + channel((r + m) * 255) + channel((g + m) * 255) + channel((b + m) * 255);
  }

  function rgbaParts(color) {
    var match = String(color || '').match(/rgba?\(\s*([\d.]+)\s*,\s*([\d.]+)\s*,\s*([\d.]+)(?:\s*,\s*([\d.]+))?\s*\)/i);
    if (!match) return { hex: '#ffffff', alpha: 0.92 };
    return {
      hex: '#' + channel(match[1]) + channel(match[2]) + channel(match[3]),
      alpha: match[4] == null ? 1 : +match[4]
    };
  }

  function rgba(hex, alpha) {
    var n = parseInt(String(hex).slice(1), 16);
    return 'rgba(' + ((n >> 16) & 255) + ',' + ((n >> 8) & 255) + ',' + (n & 255) + ',' + alpha + ')';
  }

  // The clipboard text is the object in the file, so a volunteer can paste it back.
  function formatReference(ref) {
    function quote(value) {
      return "'" + String(value).replace(/\\/g, '\\\\').replace(/'/g, "\\'") + "'";
    }
    function num(value) {
      return String(Math.round(Number(value) * 100) / 100);
    }
    var stroke = ref.lineWidth;
    var halo = ref.haloWidth;
    var show = ref.show;
    var earth = ref.earth;
    return [
      'var REFERENCE = {',
      '  earth: { land: ' + quote(earth.land) + ', water: ' + quote(earth.water) + ', park: ' + quote(earth.park) + ', building: ' + quote(earth.building) + ' },',
      '  palette: ' + quote(ref.palette) + ',',
      '  base: ' + quote(ref.base) + ',',
      '  roadCasing: ' + quote(ref.roadCasing) + ',',
      '  roadInner: ' + quote(ref.roadInner) + ',',
      '  minorRoad: ' + quote(ref.minorRoad) + ',',
      '  path: ' + quote(ref.path) + ',',
      '  rail: ' + quote(ref.rail) + ',',
      '  railDash: ' + quote(ref.railDash) + ',',
      '  boundary: ' + quote(ref.boundary) + ',',
      '  water: ' + quote(ref.water) + ',',
      '  halo: ' + quote(ref.halo) + ',',
      '  lineWidth: { road: ' + num(stroke.road) + ', path: ' + num(stroke.path) + ', rail: ' + num(stroke.rail) + ', boundary: ' + num(stroke.boundary) + ', river: ' + num(stroke.river) + ', coast: ' + num(stroke.coast) + ' },',
      '  haloWidth: { road: ' + num(halo.road) + ', path: ' + num(halo.path) + ', boundary: ' + num(halo.boundary) + ', river: ' + num(halo.river) + ', coast: ' + num(halo.coast) + ' },',
      '  coastWidth: [' + ref.coastWidth.join(', ') + '],',
      '  riverWidth: [' + ref.riverWidth.join(', ') + '],',
      '  labelHalo: ' + num(ref.labelHalo) + ',',
      '  show: { roads: ' + !!show.roads + ', rails: ' + !!show.rails + ', boundaries: ' + !!show.boundaries + ', coast: ' + !!show.coast + ', rivers: ' + !!show.rivers + ', labels: ' + !!show.labels + ' }',
      '};'
    ].join('\n');
  }

  function field(labelText) {
    var label = document.createElement('label');
    var span = document.createElement('span');
    span.textContent = labelText;
    label.appendChild(span);
    return label;
  }

  function mount(api) {
    var stored = savedCopy();
    if (stored) api.applyReference(stored);
    var working = api.current();

    var style = document.createElement('style');
    style.textContent = [
      '#reference-tune { position: fixed; top: 156px; left: 10px; z-index: 1100; box-sizing: border-box;',
      'width: 230px; max-height: calc(100vh - 176px); overflow: auto; padding: 12px;',
      'background: rgba(255,255,255,0.96); color: #1b242b; border: 1px solid #c5ced6;',
      'border-radius: 8px; box-shadow: 0 8px 24px rgba(0,0,0,0.18);',
      'font: 12px/1.4 system-ui, sans-serif; }',
      '#reference-tune h2 { margin: 0 0 4px; font-size: 13px; }',
      '#reference-tune p { margin: 0 0 8px; color: #4a5964; }',
      '#reference-tune label { display: flex; align-items: center; justify-content: space-between; gap: 8px; margin: 4px 0; }',
      '#reference-tune input[type="color"] { width: 36px; height: 22px; padding: 0; border: 1px solid #c5ced6; background: #fff; }',
      '#reference-tune input[type="range"] { width: 108px; }',
      '#reference-tune .actions, #reference-tune .palette { display: flex; gap: 8px; margin-top: 8px; }',
      '#reference-tune button { flex: 1; }',
      '#reference-tune button.on { font-weight: 700; outline: 2px solid #1b242b; }',
      '#reference-tune .note { min-height: 1.2em; margin: 6px 0 0; }'
    ].join('');
    document.head.appendChild(style);

    var panel = document.createElement('aside');
    panel.id = 'reference-tune';
    var title = document.createElement('h2');
    title.textContent = 'Reference';
    var intro = document.createElement('p');
    intro.textContent = 'Land and water sit under the dose colors. Contrast keeps the gray map and the black-and-white lines. Map colors loads the colored map and leaves its colors alone.';
    panel.appendChild(title);
    panel.appendChild(intro);

    function publish() {
      api.applyReference(working);
      remember(working);
    }

    var earthInputs = {};
    EARTH.forEach(function (spec) {
      var label = field(spec[1]);
      var input = document.createElement('input');
      input.type = 'color';
      input.setAttribute('data-earth', spec[0]);
      var shown = working.earth[spec[0]] || api.earthColor(spec[0]);
      input.value = hexOf(shown);
      input.addEventListener('input', function () {
        working.earth[spec[0]] = input.value;
        publish();
      });
      earthInputs[spec[0]] = input;
      label.appendChild(input);
      panel.appendChild(label);
    });

    function syncEarthInputs() {
      EARTH.forEach(function (spec) {
        if (working.earth[spec[0]]) return;
        var live = api.earthColor(spec[0]);
        if (live) earthInputs[spec[0]].value = hexOf(live);
      });
    }
    api.whenApplied(syncEarthInputs);

    var palette = document.createElement('div');
    palette.className = 'palette';
    var contrastBtn = document.createElement('button');
    contrastBtn.type = 'button';
    contrastBtn.textContent = 'Contrast';
    var mapBtn = document.createElement('button');
    mapBtn.type = 'button';
    mapBtn.textContent = 'Map colors';
    function markPalette() {
      var mapOn = working.palette === 'map';
      contrastBtn.classList.toggle('on', !mapOn);
      mapBtn.classList.toggle('on', mapOn);
      // The line inks belong to the contrast look. Width sliders stay usable on either map.
      var inputs = panel.querySelectorAll('input[type="color"]:not([data-earth]), input[data-halo]');
      for (var i = 0; i < inputs.length; i++) inputs[i].disabled = mapOn;
    }
    function watchEarth() {
      var waits = 0;
      var timer = setInterval(function () {
        syncEarthInputs();
        if (++waits > 20) clearInterval(timer);
      }, 200);
    }
    function usePalette(mode) {
      working.palette = mode;
      working.base = mode === 'map' ? 'bright' : 'theme';
      working.earth = { land: '', water: '', park: '', building: '' };
      markPalette();
      publish();
      watchEarth();
    }
    contrastBtn.addEventListener('click', function () { usePalette('contrast'); });
    mapBtn.addEventListener('click', function () { usePalette('map'); });
    markPalette();
    palette.appendChild(contrastBtn);
    palette.appendChild(mapBtn);
    panel.appendChild(palette);

    COLORS.forEach(function (spec) {
      var label = field(spec[1]);
      var input = document.createElement('input');
      input.type = 'color';
      input.value = hexOf(working[spec[0]]);
      input.addEventListener('input', function () {
        working[spec[0]] = input.value;
        publish();
      });
      label.appendChild(input);
      panel.appendChild(label);
    });

    var widthNote = document.createElement('p');
    widthNote.textContent = 'Width is the line. Halo is only the pale rim. 1 matches the map’s own width.';
    panel.appendChild(widthNote);

    STROKES.forEach(function (spec) {
      var label = field(spec[1]);
      var input = document.createElement('input');
      input.type = 'range';
      input.min = '0.5';
      input.max = '3';
      input.step = '0.1';
      input.value = String(working.lineWidth[spec[0]]);
      input.addEventListener('input', function () {
        working.lineWidth[spec[0]] = Number(input.value);
        publish();
      });
      label.appendChild(input);
      panel.appendChild(label);
    });

    var halo = rgbaParts(working.halo);
    var haloLabel = field('Halo');
    var haloColor = document.createElement('input');
    haloColor.type = 'color';
    haloColor.setAttribute('data-halo', '1');
    haloColor.value = halo.hex;
    var haloAlpha = document.createElement('input');
    haloAlpha.type = 'range';
    haloAlpha.setAttribute('data-halo', '1');
    haloAlpha.min = '0';
    haloAlpha.max = '1';
    haloAlpha.step = '0.01';
    haloAlpha.value = String(halo.alpha);
    function paintHalo() {
      working.halo = rgba(haloColor.value, Number(haloAlpha.value));
      publish();
    }
    haloColor.addEventListener('input', paintHalo);
    haloAlpha.addEventListener('input', paintHalo);
    haloLabel.appendChild(haloAlpha);
    haloLabel.appendChild(haloColor);
    panel.appendChild(haloLabel);

    WIDTHS.forEach(function (spec) {
      var label = field(spec[1]);
      var input = document.createElement('input');
      input.type = 'range';
      input.setAttribute('data-halo', '1');
      input.min = '0';
      input.max = '4';
      input.step = '0.1';
      input.value = String(working.haloWidth[spec[0]]);
      input.addEventListener('input', function () {
        working.haloWidth[spec[0]] = Number(input.value);
        publish();
      });
      label.appendChild(input);
      panel.appendChild(label);
    });

    var labelWidth = field('Label halo');
    var labelInput = document.createElement('input');
    labelInput.type = 'range';
    labelInput.setAttribute('data-halo', '1');
    labelInput.min = '0';
    labelInput.max = '4';
    labelInput.step = '0.1';
    labelInput.value = String(working.labelHalo);
    labelInput.addEventListener('input', function () {
      working.labelHalo = Number(labelInput.value);
      publish();
    });
    labelWidth.appendChild(labelInput);
    panel.appendChild(labelWidth);

    GROUPS.forEach(function (spec) {
      var label = field(spec[1]);
      var input = document.createElement('input');
      input.type = 'checkbox';
      input.checked = working.show[spec[0]] !== false;
      input.addEventListener('change', function () {
        working.show[spec[0]] = input.checked;
        publish();
      });
      label.appendChild(input);
      panel.appendChild(label);
    });
    markPalette();
    watchEarth();

    var note = document.createElement('p');
    note.className = 'note';
    var actions = document.createElement('div');
    actions.className = 'actions';
    var copy = document.createElement('button');
    copy.type = 'button';
    copy.textContent = 'Copy';
    var reset = document.createElement('button');
    reset.type = 'button';
    reset.textContent = 'Reset';
    copy.addEventListener('click', function () {
      var text = formatReference(working);
      function show(message) { note.textContent = message; }
      if (navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard.writeText(text).then(function () {
          show('Copied');
        }).catch(function () {
          show(text);
        });
        return;
      }
      show(text);
    });
    reset.addEventListener('click', function () {
      working = api.reference();
      forget();
      api.applyReference(working);
      panel.remove();
      style.remove();
      mount(api);
    });
    actions.appendChild(copy);
    actions.appendChild(reset);
    panel.appendChild(actions);
    panel.appendChild(note);
    document.body.appendChild(panel);
  }

  function start() {
    var params = new URLSearchParams(window.location.search);
    var tune = params.get('tune') === '1';
    var tries = 0;
    var timer = setInterval(function () {
      var api = window.safecastBasemap;
      if (!api || typeof api.applyReference !== 'function' || typeof api.reference !== 'function') {
        if (++tries > 200) clearInterval(timer);
        return;
      }
      clearInterval(timer);
      // A refresh keeps the trial in this tab. Another visitor never sees it.
      var stored = savedCopy();
      if (stored) api.applyReference(stored);
      if (tune) mount(api);
    }, 50);
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', start);
  } else {
    start();
  }
})();
