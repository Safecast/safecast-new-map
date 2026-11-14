// Street colors stay under the radiation overlay. Roads, railways, borders,
// coastlines, rivers, and labels are drawn above it so a colored area still
// shows where you are. One MapLibre style holds all of those layers; the
// radiation PNG is a raster layer inserted between the two groups.
//
// Building footprints stay in the color group. An outline of every building
// above the dose colors hides the measurement at city zoom.
// Place names and airports stay above. Point-of-interest icons do not: the
// Positron and Dark styles omit them, and adding them back would cover the dose.
//
// REFERENCE is the only place those colors and widths are chosen.
// Earth is the ground under the radiation tiles. An empty earth color keeps
// whatever the OpenFreeMap style drew. "contrast" is the gray map with
// black-and-white lines, so a dose color cannot hide a road. "map" with
// base "bright" is OpenFreeMap's colored style, left as published.
// Above the radiation raster the stack is coast, rivers, roads and rails,
// borders, then labels. Open the map with ?tune=1 to try other values.
(function () {
  var REFERENCE = {
    earth: { land: '', water: '', park: '', building: '' },
    palette: 'contrast',
    base: 'theme',
    roadCasing: '#1b242b',
    roadInner: '#ffffff',
    minorRoad: '#161616',
    path: '#3e4c57',
    rail: '#161616',
    railDash: '#ffffff',
    boundary: '#1c2830',
    water: '#0d5c78',
    halo: 'rgba(255,255,255,0.92)',
    // 1 keeps the stroke the style already drew. Halo widths are the pale rim, in pixels.
    lineWidth: { road: 1, path: 1, rail: 1, boundary: 1, river: 1, coast: 1 },
    haloWidth: { road: 1.6, path: 1.1, boundary: 1.5, river: 1.3, coast: 1.2 },
    coastWidth: [0, 0.7, 6, 1.2, 12, 1.8],
    riverWidth: [3, 0.6, 8, 1.1, 13, 1.8],
    labelHalo: 1.3,
    show: { roads: true, rails: true, boundaries: true, coast: true, rivers: true, labels: true }
  };

  var LIGHT_STYLE = 'https://tiles.openfreemap.org/styles/positron';
  var DARK_STYLE = 'https://tiles.openfreemap.org/styles/dark';
  // Bright is the colored OpenFreeMap style: blue water, green parks, orange roads.
  var BRIGHT_STYLE = 'https://tiles.openfreemap.org/styles/bright';
  // OpenFreeMap's planet TileJSON asks for this credit.
  var ATTRIBUTION = '<a href="https://openfreemap.org" target="_blank">OpenFreeMap</a> <a href="https://www.openmaptiles.org/" target="_blank">&copy; OpenMapTiles</a> Data from <a href="https://www.openstreetmap.org/copyright" target="_blank">OpenStreetMap</a>';
  var GOOGLE_ATTR = '&copy; Google';
  // The Leaflet binding shows MapLibre zoom N at Leaflet zoom N+1.
  // These are MapLibre zooms: municipalities from about Leaflet 10, wards from 12.
  var ADMIN_LEVELS = [
    { level: 6, minzoom: 7, dash: [2, 1], width: ['interpolate', ['linear'], ['zoom'], 7, 0.8, 12, 1.5] },
    { level: 7, minzoom: 9, dash: [3, 1.5], width: ['interpolate', ['linear'], ['zoom'], 9, 0.6, 13, 1.3] },
    { level: 8, minzoom: 11, dash: [1, 1.4], width: ['interpolate', ['linear'], ['zoom'], 11, 0.5, 14, 1.1] }
  ];

  var styleCache = {};
  var generation = 0;
  var inflightSignature = null;
  var appliedSignature = '';
  var referenceApplied = null;
  // Unscaled stroke widths, so a slider change multiplies the style width once.
  var baseWidths = {};
  var glLayer = null;
  var gl = null;
  var map = null;
  var fallbackURLs = null;
  var radiationURL = '';
  var radiationVisible = false;
  var appliedRadiationURL = '';
  var googleAttrOn = false;
  var status = 'loading';
  var activeReference = cloneRef(REFERENCE);

  function cloneRef(src) {
    return JSON.parse(JSON.stringify(src));
  }

  function currentRef() {
    return activeReference;
  }

  // A partial preview copy keeps any key it leaves out at the file default.
  function normalizeRef(next) {
    var base = cloneRef(REFERENCE);
    if (!next || typeof next !== 'object') return base;
    Object.keys(base).forEach(function (key) {
      if (next[key] == null) return;
      if (key === 'palette') {
        if (next.palette === 'map' || next.palette === 'contrast') base.palette = next.palette;
        return;
      }
      if (key === 'base') {
        if (next.base === 'bright' || next.base === 'theme') base.base = next.base;
        return;
      }
      if (key === 'earth' || key === 'lineWidth' || key === 'haloWidth' || key === 'show') {
        Object.keys(base[key]).forEach(function (sub) {
          if (next[key][sub] != null) base[key][sub] = next[key][sub];
        });
        return;
      }
      base[key] = next[key];
    });
    return base;
  }

  function stopsToWidth(stops) {
    var expr = ['interpolate', ['linear'], ['zoom']];
    for (var i = 0; i < stops.length; i++) expr.push(stops[i]);
    return expr;
  }

  function isHalo(layer) {
    return (layer.id || '').indexOf('-safecast-halo') !== -1;
  }

  function baseId(layer) {
    return (layer.id || '').replace(/-safecast-halo$/, '');
  }

  // One name for the line and its halo, so a checkbox hides both.
  function roleOf(layer) {
    if (layer.type === 'symbol') return 'labels';
    if (layer.type !== 'line') return '';
    var id = isHalo(layer) ? baseId(layer) : (layer.id || '');
    var src = layer['source-layer'] || '';
    if (id === 'safecast-water-edge') return 'coast';
    if (src === 'waterway' || id === 'waterway') return 'rivers';
    if (src === 'boundary' || id.indexOf('boundary') !== -1) return 'boundaries';
    if (id.indexOf('railway') !== -1 || id.indexOf('transit') !== -1) return 'rails';
    if (src === 'transportation' || src === 'aeroway' || id.indexOf('highway') !== -1 || id.indexOf('aeroway') !== -1 || id.indexOf('road_') === 0) return 'roads';
    return '';
  }

  // Land, water, parks, and buildings are the fills under the dose tiles.
  function earthRole(layer) {
    var id = layer.id || '';
    var src = layer['source-layer'] || '';
    if (layer.type === 'background') return 'land';
    if (layer.type !== 'fill') return '';
    if (src === 'water' || id === 'water' || id.indexOf('water') === 0) return 'water';
    if (src === 'building' || id.indexOf('building') !== -1) return 'building';
    if (src === 'park' || id.indexOf('park') !== -1 || id.indexOf('wood') !== -1 || id.indexOf('grass') !== -1) return 'park';
    if (id.indexOf('residential') !== -1) return 'land';
    return '';
  }

  function paintEarthOn(layer) {
    var role = earthRole(layer);
    var color = role && currentRef().earth && currentRef().earth[role];
    if (!color) return;
    layer.paint = layer.paint || {};
    layer.paint[layer.type === 'background' ? 'background-color' : 'fill-color'] = color;
  }

  function visibility(ref, role) {
    if (!role || !ref.show || ref.show[role] !== false) return 'visible';
    return 'none';
  }

  function haloExtra(ref, role, id) {
    if (role === 'coast') return ref.haloWidth.coast;
    if (role === 'rivers') return ref.haloWidth.river;
    if (role === 'boundaries') return ref.haloWidth.boundary;
    if (id.indexOf('path') !== -1) return ref.haloWidth.path;
    return ref.haloWidth.road;
  }

  function strokeColor(layer, ref) {
    var id = layer.id || '';
    var role = roleOf(layer);
    if (role === 'rivers' || role === 'coast') return ref.water;
    if (role === 'boundaries') return ref.boundary;
    if (role === 'rails') return id.indexOf('dash') !== -1 ? ref.railDash : ref.rail;
    if (role !== 'roads') return '';
    if (id.indexOf('casing') !== -1) return ref.roadCasing;
    if (id.indexOf('dash') !== -1 || id.indexOf('inner') !== -1 || id === 'aeroway-runway') return ref.roadInner;
    if (id.indexOf('path') !== -1 || id.indexOf('pier') !== -1) return ref.path;
    return ref.minorRoad;
  }

  function markVisible(layer, role) {
    layer.layout = layer.layout || {};
    layer.layout.visibility = visibility(currentRef(), role);
  }

  function themeNow() {
    var attr = document.documentElement.dataset.theme;
    if (attr === 'dark' || attr === 'light') return attr;
    try {
      var stored = sessionStorage.getItem('themePreference');
      if (stored === 'dark' || stored === 'light') return stored;
    } catch (err) {}
    if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) return 'dark';
    return 'light';
  }

  function absoluteTileURL(url) {
    if (!url) return '';
    if (/^https?:\/\//i.test(url)) return url;
    // URL() percent-encodes the {z}/{x}/{y} placeholders MapLibre substitutes.
    if (url.charAt(0) !== '/') url = '/' + url;
    return window.location.origin + url;
  }

  function styleURL() {
    if (currentRef().base === 'bright') return BRIGHT_STYLE;
    return themeNow() === 'dark' ? DARK_STYLE : LIGHT_STYLE;
  }

  // Theme is ignored while the colored style is selected, so the light/dark
  // switch does not swap that trial back to the gray map.
  function styleSignature() {
    var ref = currentRef();
    var theme = ref.base === 'bright' ? 'bright' : themeNow();
    return theme + '|' + ref.palette;
  }

  function loadStyle() {
    var url = styleURL();
    if (styleCache[url]) return Promise.resolve(styleCache[url]);
    return fetch(url).then(function (response) {
      if (!response.ok) throw new Error('style ' + response.status);
      return response.json();
    }).then(function (json) {
      styleCache[url] = json;
      return json;
    });
  }

  function vectorSourceName(style) {
    var names = Object.keys(style.sources || {});
    for (var i = 0; i < names.length; i++) {
      if (style.sources[names[i]].type === 'vector') return names[i];
    }
    return 'openmaptiles';
  }

  function isAdminGet(expr) {
    return Array.isArray(expr) && expr[0] === 'get' && expr[1] === 'admin_level';
  }

  // Levels the style already strokes, so we don't draw a second prefecture line.
  function coveredAdminLevels(style) {
    var levels = {};
    function mark(filter) {
      if (!Array.isArray(filter) || !filter.length) return;
      var op = filter[0];
      if (op === '==' && isAdminGet(filter[1]) && typeof filter[2] === 'number') {
        levels[filter[2]] = true;
      }
      if (op === 'all' || op === 'any') {
        var ge = null;
        var le = null;
        for (var i = 1; i < filter.length; i++) {
          var child = filter[i];
          if (!Array.isArray(child)) continue;
          if ((child[0] === '>=' || child[0] === '>') && isAdminGet(child[1]) && typeof child[2] === 'number') {
            ge = child[0] === '>' ? child[2] + 1 : child[2];
          } else if ((child[0] === '<=' || child[0] === '<') && isAdminGet(child[1]) && typeof child[2] === 'number') {
            le = child[0] === '<' ? child[2] - 1 : child[2];
          }
          mark(child);
        }
        if (ge != null && le != null) {
          for (var n = ge; n <= le && n < 12; n++) levels[n] = true;
        }
        return;
      }
      for (var j = 1; j < filter.length; j++) mark(filter[j]);
    }
    (style.layers || []).forEach(function (layer) {
      if (layer['source-layer'] === 'boundary') mark(layer.filter);
    });
    return levels;
  }

  function cloneLineShell(layer, suffix) {
    var out = {
      id: layer.id + suffix,
      type: 'line',
      source: layer.source,
      'source-layer': layer['source-layer'],
      layout: { 'line-cap': 'round', 'line-join': 'round' }
    };
    if (layer.filter) out.filter = layer.filter;
    if (typeof layer.minzoom === 'number') out.minzoom = layer.minzoom;
    if (typeof layer.maxzoom === 'number') out.maxzoom = layer.maxzoom;
    return out;
  }

  // Zoom stops have to stay at the top of the expression. Wrapping an
  // interpolate in "+" is rejected by the style parser, so the extra pixels
  // are added to each output stop instead.
  function lineScale(ref, role, id) {
    var scale = ref.lineWidth || {};
    var value = 1;
    if (role === 'coast') value = scale.coast;
    else if (role === 'rivers') value = scale.river;
    else if (role === 'boundaries') value = scale.boundary;
    else if (role === 'rails') value = scale.rail;
    else if (role === 'roads' && (id || '').indexOf('path') !== -1) value = scale.path;
    else if (role === 'roads') value = scale.road;
    return typeof value === 'number' && value > 0 ? value : 1;
  }

  // Multiply each output stop. The zoom stops stay put, so a road still grows as you zoom in.
  function scaleWidth(width, factor) {
    if (factor === 1 || width == null) return width;
    if (typeof width === 'number') return Math.max(0.05, width * factor);
    if (Array.isArray(width) && (width[0] === 'interpolate' || width[0] === 'step')) {
      var copy = JSON.parse(JSON.stringify(width));
      var start = copy[0] === 'interpolate' ? 4 : 2;
      for (var i = start; i < copy.length; i += 2) {
        if (typeof copy[i] === 'number') copy[i] = Math.max(0.05, copy[i] * factor);
      }
      return copy;
    }
    return width;
  }

  function rememberWidth(id, width) {
    if (id && width != null && baseWidths[id] == null) baseWidths[id] = width;
  }

  function sizedWidth(id, role, fallback) {
    var ref = currentRef();
    var base = baseWidths[id];
    if (role === 'coast') base = stopsToWidth(ref.coastWidth);
    else if (role === 'rivers') base = stopsToWidth(ref.riverWidth);
    else if (base == null) base = fallback;
    return scaleWidth(base, lineScale(ref, role, id));
  }

  function wider(width, extra) {
    if (typeof width === 'number') return Math.max(0, width + extra);
    if (Array.isArray(width) && (width[0] === 'interpolate' || width[0] === 'step')) {
      var copy = JSON.parse(JSON.stringify(width));
      var start = copy[0] === 'interpolate' ? 4 : 2;
      for (var i = start; i < copy.length; i += 2) {
        if (typeof copy[i] === 'number') copy[i] = Math.max(0, copy[i] + extra);
      }
      return copy;
    }
    return ['interpolate', ['linear'], ['zoom'], 4, extra + 0.6, 14, extra + 1.6];
  }

  // A light edge keeps a dark stroke visible on navy dose colors, and the
  // dark stroke keeps it visible on yellow and white.
  function lightHalo(layer, extra) {
    var shell = cloneLineShell(layer, '-safecast-halo');
    var width = layer.paint && layer.paint['line-width'];
    shell.paint = {
      'line-color': currentRef().halo,
      'line-width': wider(width, extra)
    };
    markVisible(shell, roleOf(layer));
    return shell;
  }

  function strengthenLabel(layer, theme) {
    if (layer.type !== 'symbol') return;
    var paint = layer.paint || {};
    layer.paint = paint;
    markVisible(layer, 'labels');
    if (currentRef().palette === 'map') return;
    if (!paint['text-color']) return;
    // The halo color follows the theme. The width is the one volunteers tune.
    var width = currentRef().labelHalo;
    if (!paint['text-halo-color']) {
      paint['text-halo-color'] = theme === 'dark' ? 'rgba(0,0,0,0.8)' : 'rgba(255,255,255,0.92)';
      paint['text-halo-width'] = width;
      return;
    }
    if (typeof paint['text-halo-width'] === 'number' && paint['text-halo-width'] < width) {
      paint['text-halo-width'] = width;
    }
  }

  // Casings and dashes are already a pair. Only the single-stroke lines need a halo.
  function needsHalo(layer) {
    var src = layer['source-layer'];
    var id = layer.id || '';
    if (src === 'boundary' || src === 'waterway') return true;
    if (src !== 'transportation' && src !== 'aeroway') return false;
    if (id.indexOf('casing') !== -1 || id.indexOf('dash') !== -1) return false;
    if (id.indexOf('inner') !== -1 || id === 'aeroway-runway') return false;
    return true;
  }

  function pushReferenceLine(layer, lines) {
    var src = layer['source-layer'];
    var paint = layer.paint || {};
    layer.paint = paint;
    var role = roleOf(layer);
    rememberWidth(layer.id, paint['line-width']);
    var sized = sizedWidth(layer.id, role, paint['line-width']);
    if (sized != null && role) paint['line-width'] = sized;
    if (currentRef().palette === 'map') {
      markVisible(layer, role);
      lines.push(layer);
      return;
    }
    if (src === 'boundary' || src === 'waterway' || src === 'transportation' || src === 'aeroway') {
      var color = strokeColor(layer, currentRef());
      if (color) paint['line-color'] = color;
      markVisible(layer, role);
      if (needsHalo(layer)) lines.push(lightHalo(layer, haloExtra(currentRef(), role, layer.id || '')));
      lines.push(layer);
      return;
    }
    lines.push(layer);
  }

  function waterEdge(sourceName) {
    var layer = {
      id: 'safecast-water-edge',
      type: 'line',
      source: sourceName,
      'source-layer': 'water',
      filter: ['all',
        ['match', ['geometry-type'], ['Polygon', 'MultiPolygon'], true, false],
        ['!=', ['get', 'brunnel'], 'tunnel']
      ],
      paint: {
        'line-color': currentRef().water,
        'line-width': sizedWidth('safecast-water-edge', 'coast', null)
      }
    };
    markVisible(layer, 'coast');
    return [lightHalo(layer, currentRef().haloWidth.coast), layer];
  }

  function extraBoundaries(sourceName, covered) {
    var layers = [];
    ADMIN_LEVELS.forEach(function (spec) {
      if (covered[spec.level]) return;
      layers.push({
        id: 'safecast-boundary-' + spec.level,
        type: 'line',
        source: sourceName,
        'source-layer': 'boundary',
        minzoom: spec.minzoom,
        filter: ['all',
          ['==', ['get', 'admin_level'], spec.level],
          ['!=', ['get', 'maritime'], 1],
          ['!=', ['get', 'disputed'], 1]
        ],
        paint: {
          'line-color': currentRef().palette === 'map' ? '#9aa0a6' : currentRef().boundary,
          'line-dasharray': spec.dash,
          'line-width': spec.width
        }
      });
    });
    return layers;
  }

  function composeStyle(baseStyle, theme, radiationTiles) {
    baseWidths = {};
    var style = JSON.parse(JSON.stringify(baseStyle));
    var sourceName = vectorSourceName(style);
    var below = [];
    var lines = [];
    var symbols = [];
    style.layers.forEach(function (layer) {
      if (layer.type === 'background' || layer.type === 'fill' || layer.type === 'raster') {
        paintEarthOn(layer);
        below.push(layer);
      } else if (layer.type === 'symbol') {
        strengthenLabel(layer, theme);
        symbols.push(layer);
      } else {
        pushReferenceLine(layer, lines);
      }
    });
    var edge = currentRef().palette === 'map' ? [] : waterEdge(sourceName);
    lines = edge.concat(lines);
    extraBoundaries(sourceName, coveredAdminLevels(style)).forEach(function (layer) {
      pushReferenceLine(layer, lines);
    });

    style.sources['safecast-satellite'] = {
      type: 'raster',
      tiles: [0, 1, 2, 3].map(function (n) {
        return 'https://mt' + n + '.google.com/vt/lyrs=s&x={x}&y={y}&z={z}';
      }),
      tileSize: 256,
      maxzoom: 20
    };
    style.sources['safecast-radiation'] = {
      type: 'raster',
      tiles: [radiationTiles],
      tileSize: 256,
      maxzoom: 13
    };
    // Satellite sits on the land colors. Radiation sits on that. Lines and
    // labels stay last, so they remain visible over a colored measurement.
    var satellite = map && window.googleSatellite && map.hasLayer(window.googleSatellite);
    style.layers = below.concat([
      {
        id: 'safecast-satellite',
        type: 'raster',
        source: 'safecast-satellite',
        layout: { visibility: satellite ? 'visible' : 'none' },
        paint: { 'raster-fade-duration': 0 }
      },
      {
        id: 'safecast-radiation',
        type: 'raster',
        source: 'safecast-radiation',
        layout: { visibility: radiationVisible && radiationURL ? 'visible' : 'none' },
        paint: { 'raster-opacity': 1, 'raster-fade-duration': 0 }
      }
    ], lines, symbols);
    return style;
  }

  function placeholderTiles() {
    return window.location.origin + '/tiles/aggregate/{z}/{x}/{y}.png?speed=none';
  }

  function applySatelliteVisibility() {
    if (!gl || !gl.getLayer('safecast-satellite')) return;
    var satellite = !!(window.googleSatellite && map.hasLayer(window.googleSatellite));
    gl.setLayoutProperty('safecast-satellite', 'visibility', satellite ? 'visible' : 'none');
  }

  function applyRadiation() {
    if (!gl || !gl.getSource('safecast-radiation')) return;
    if (radiationURL && radiationURL !== appliedRadiationURL) {
      gl.getSource('safecast-radiation').setTiles([radiationURL]);
      appliedRadiationURL = radiationURL;
    }
    if (gl.getLayer('safecast-radiation')) {
      var show = radiationVisible && !!radiationURL;
      gl.setLayoutProperty('safecast-radiation', 'visibility', show ? 'visible' : 'none');
    }
  }

  function syncGoogleAttribution() {
    if (!map || !map.attributionControl) return;
    var want = status === 'vector' && !!(window.googleSatellite && map.hasLayer(window.googleSatellite));
    if (want && !googleAttrOn) {
      map.attributionControl.addAttribution(GOOGLE_ATTR);
      googleAttrOn = true;
    } else if (!want && googleAttrOn) {
      map.attributionControl.removeAttribution(GOOGLE_ATTR);
      googleAttrOn = false;
    }
  }

  function setPaint(id, key, value) {
    if (!gl || !gl.getLayer(id) || value == null) return;
    try { gl.setPaintProperty(id, key, value); } catch (err) {}
  }

  function setVis(id, value) {
    if (!gl || !gl.getLayer(id)) return;
    try { gl.setLayoutProperty(id, 'visibility', value); } catch (err) {}
  }

  function paintReference() {
    var ref = activeReference;
    var layers;
    try {
      layers = gl && gl.getStyle() && gl.getStyle().layers;
    } catch (err) {
      return;
    }
    if (!layers) return;
    var halos = [];
    var keepStyle = ref.palette === 'map';
    layers.forEach(function (layer) {
      if (isHalo(layer)) {
        halos.push(layer);
        return;
      }
      var ground = earthRole(layer);
      var earth = ground && ref.earth && ref.earth[ground];
      if (earth) setPaint(layer.id, layer.type === 'background' ? 'background-color' : 'fill-color', earth);
      var role = roleOf(layer);
      if (!role) return;
      if (role === 'labels') {
        if (!keepStyle && gl.getPaintProperty(layer.id, 'text-color')) setPaint(layer.id, 'text-halo-width', ref.labelHalo);
        setVis(layer.id, visibility(ref, 'labels'));
        return;
      }
      var currentWidth = gl.getPaintProperty(layer.id, 'line-width');
      rememberWidth(layer.id, currentWidth);
      var sized = sizedWidth(layer.id, role, currentWidth);
      if (sized != null) setPaint(layer.id, 'line-width', sized);
      if (keepStyle) {
        setVis(layer.id, visibility(ref, role));
        return;
      }
      var color = strokeColor(layer, ref);
      if (color) setPaint(layer.id, 'line-color', color);
      setVis(layer.id, visibility(ref, role));
    });
    if (keepStyle) return;
    halos.forEach(function (layer) {
      var id = baseId(layer);
      var role = roleOf(layer);
      setPaint(layer.id, 'line-color', ref.halo);
      setPaint(layer.id, 'line-width', wider(gl.getPaintProperty(id, 'line-width'), haloExtra(ref, role, id)));
      setVis(layer.id, visibility(ref, role));
    });
  }

  // The style builder and the preview panel share this path. A theme reload
  // calls it again so a trial is not wiped by setStyle. Switching between the
  // gray ink and the colored map rebuilds the style, because those add and
  // remove layers rather than only recoloring them.
  function earthWasCleared(prev, next) {
    var keys = ['land', 'water', 'park', 'building'];
    for (var i = 0; i < keys.length; i++) {
      if (prev.earth[keys[i]] && !next.earth[keys[i]]) return true;
    }
    return false;
  }

  function applyReference(next) {
    var prev = activeReference;
    if (next) activeReference = normalizeRef(next);
    // A cleared earth color has to come back from the style, which means a rebuild.
    if (activeReference.base !== prev.base || activeReference.palette !== prev.palette || earthWasCleared(prev, activeReference)) {
      appliedSignature = '';
      loadAndApply();
      return;
    }
    paintReference();
  }

  function firstEarthColor(role) {
    if (!gl || !gl.getStyle()) return '';
    var layers = gl.getStyle().layers || [];
    for (var i = 0; i < layers.length; i++) {
      if (earthRole(layers[i]) !== role) continue;
      var key = layers[i].type === 'background' ? 'background-color' : 'fill-color';
      var value = gl.getPaintProperty(layers[i].id, key);
      if (typeof value === 'string' && value) return value;
    }
    return '';
  }

  function onStyleLoad() {
    applySatelliteVisibility();
    applyRadiation();
    paintReference();
    syncGoogleAttribution();
    if (referenceApplied) referenceApplied();
  }

  function applyLoadedStyle(json, theme) {
    var baked = radiationURL || placeholderTiles();
    var style = composeStyle(json, theme, baked);
    appliedSignature = styleSignature();
    appliedRadiationURL = baked;
    if (!glLayer) {
      glLayer = L.maplibreGL({
        style: style,
        attribution: ATTRIBUTION,
        interactive: false,
        padding: 0.15,
        fadeDuration: 0
      });
      glLayer.addTo(map);
      gl = glLayer.getMaplibreMap();
      gl.on('style.load', onStyleLoad);
      status = 'vector';
      api.status = 'vector';
      document.documentElement.dataset.basemap = 'vector';
    } else if (gl) {
      gl.setStyle(style);
      // setStyle loads the new style asynchronously. style.load is what
      // reapplies a trial; listen once in case the standing listener misses it.
      gl.once('style.load', onStyleLoad);
    }
  }

  function enableRasterFallback() {
    if (status === 'fallback') return;
    status = 'fallback';
    api.status = 'fallback';
    document.documentElement.dataset.basemap = 'raster';
    syncGoogleAttribution();
    if (glLayer && map && map.hasLayer(glLayer)) map.removeLayer(glLayer);
    glLayer = null;
    gl = null;
    var osmURL = (fallbackURLs && fallbackURLs.osmURL) || 'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png';
    var satelliteURL = (fallbackURLs && fallbackURLs.satelliteURL) || 'https://mt1.google.com/vt/lyrs=s&x={x}&y={y}&z={z}';
    if (window.osmLayer && window.osmLayer.addLayer && !window.osmLayer.getLayers().length) {
      window.osmLayer.addLayer(L.tileLayer(osmURL, {
        maxZoom: 18,
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
      }));
    }
    if (window.googleSatellite && window.googleSatellite.addLayer && !window.googleSatellite.getLayers().length) {
      window.googleSatellite.addLayer(L.tileLayer(satelliteURL, {
        maxZoom: 20,
        attribution: GOOGLE_ATTR
      }));
    }
    if (typeof window.updateMarkers === 'function') window.updateMarkers();
  }

  function loadAndApply() {
    if (status === 'fallback') return;
    var signature = styleSignature();
    if (signature === appliedSignature && gl) {
      applySatelliteVisibility();
      applyRadiation();
      paintReference();
      syncGoogleAttribution();
      return;
    }
    if (signature === inflightSignature) return;
    inflightSignature = signature;
    var gen = ++generation;
    loadStyle().then(function (json) {
      if (gen !== generation) return;
      inflightSignature = null;
      if (styleSignature() !== signature) {
        loadAndApply();
        return;
      }
      try {
        applyLoadedStyle(json, themeNow());
      } catch (err) {
        console.warn('reference basemap rejected style', err);
        enableRasterFallback();
      }
    }).catch(function (err) {
      if (gen !== generation) return;
      inflightSignature = null;
      console.warn('reference basemap unavailable', err);
      enableRasterFallback();
    });
  }

  var api = {
    status: 'loading',
    // A fresh copy, so a preview edit cannot change the file defaults.
    reference: function () { return cloneRef(REFERENCE); },
    current: function () { return cloneRef(activeReference); },
    earthColor: function (role) {
      var chosen = activeReference.earth && activeReference.earth[role];
      return chosen || firstEarthColor(role);
    },
    whenApplied: function (fn) { referenceApplied = fn; },
    applyReference: function (next) { applyReference(next); },
    setRadiation: function (url, visible) {
      radiationVisible = !!visible;
      if (url) radiationURL = absoluteTileURL(url);
      if (status === 'fallback') return;
      applyRadiation();
    },
    setTheme: function () {
      loadAndApply();
    },
    syncBase: function () {
      applySatelliteVisibility();
      syncGoogleAttribution();
    }
  };

  window.installReferenceBasemap = function (leafletMap, urls) {
    map = leafletMap;
    fallbackURLs = urls || null;
    window.safecastBasemap = api;
    if (!window.maplibregl || typeof L.maplibreGL !== 'function') {
      enableRasterFallback();
      return;
    }
    loadAndApply();
  };
})();
