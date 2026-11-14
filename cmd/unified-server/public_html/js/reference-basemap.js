// Street colors stay under the radiation overlay. Roads, railways, borders,
// coastlines, rivers, and labels are drawn above it so a colored area still
// shows where you are. One MapLibre style holds all of those layers; the
// radiation PNG is a raster layer inserted between the two groups.
//
// Building footprints stay in the color group. An outline of every building
// above the dose colors hides the measurement at city zoom.
// Place names and airports stay above. Point-of-interest icons do not: the
// Positron and Dark styles omit them, and adding them back would cover the dose.
(function () {
  var LIGHT_STYLE = 'https://tiles.openfreemap.org/styles/positron';
  var DARK_STYLE = 'https://tiles.openfreemap.org/styles/dark';
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
  var inflightTheme = null;
  var appliedTheme = null;
  var glLayer = null;
  var gl = null;
  var map = null;
  var fallbackURLs = null;
  var radiationURL = '';
  var radiationVisible = false;
  var appliedRadiationURL = '';
  var googleAttrOn = false;
  var status = 'loading';

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

  function loadStyle(theme) {
    var url = theme === 'dark' ? DARK_STYLE : LIGHT_STYLE;
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
      'line-color': 'rgba(255,255,255,0.92)',
      'line-width': wider(width, extra)
    };
    return shell;
  }

  function strengthenLabel(layer, theme) {
    if (layer.type !== 'symbol') return;
    var paint = layer.paint || {};
    layer.paint = paint;
    if (!paint['text-color']) return;
    if (!paint['text-halo-color']) {
      paint['text-halo-color'] = theme === 'dark' ? 'rgba(0,0,0,0.8)' : 'rgba(255,255,255,0.92)';
      paint['text-halo-width'] = 1.3;
      return;
    }
    if (typeof paint['text-halo-width'] === 'number' && paint['text-halo-width'] < 1.3) {
      paint['text-halo-width'] = 1.3;
    }
  }

  function pushReferenceLine(layer, lines) {
    var src = layer['source-layer'];
    var id = layer.id || '';
    var paint = layer.paint || {};
    layer.paint = paint;
    if (src === 'boundary') {
      lines.push(lightHalo(layer, 1.5));
      paint['line-color'] = '#1c2830';
      lines.push(layer);
      return;
    }
    if (src === 'waterway') {
      paint['line-color'] = '#0d5c78';
      paint['line-width'] = ['interpolate', ['linear'], ['zoom'], 3, 0.6, 8, 1.1, 13, 1.8];
      lines.push(lightHalo(layer, 1.3));
      lines.push(layer);
      return;
    }
    if (src === 'transportation' || src === 'aeroway') {
      if (id.indexOf('casing') !== -1) {
        paint['line-color'] = '#1b242b';
        lines.push(layer);
        return;
      }
      if (id.indexOf('dash') !== -1) {
        paint['line-color'] = '#ffffff';
        lines.push(layer);
        return;
      }
      if (id.indexOf('inner') !== -1 || id === 'aeroway-runway') {
        paint['line-color'] = '#ffffff';
        lines.push(layer);
        return;
      }
      // Minor roads, paths, and rails have no separate casing. Give them one
      // so the line still reads on both pale and dark dose colors.
      paint['line-color'] = id.indexOf('path') !== -1 || id.indexOf('pier') !== -1 ? '#3e4c57' : '#161616';
      lines.push(lightHalo(layer, id.indexOf('path') !== -1 ? 1.1 : 1.6));
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
        'line-color': '#0d5c78',
        'line-width': ['interpolate', ['linear'], ['zoom'], 0, 0.7, 6, 1.2, 12, 1.8]
      }
    };
    return [lightHalo(layer, 1.2), layer];
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
          'line-color': '#1c2830',
          'line-dasharray': spec.dash,
          'line-width': spec.width
        }
      });
    });
    return layers;
  }

  function composeStyle(baseStyle, theme, radiationTiles) {
    var style = JSON.parse(JSON.stringify(baseStyle));
    var sourceName = vectorSourceName(style);
    var below = [];
    var lines = [];
    var symbols = [];
    style.layers.forEach(function (layer) {
      if (layer.type === 'background' || layer.type === 'fill' || layer.type === 'raster') {
        below.push(layer);
      } else if (layer.type === 'symbol') {
        strengthenLabel(layer, theme);
        symbols.push(layer);
      } else {
        pushReferenceLine(layer, lines);
      }
    });
    var edge = waterEdge(sourceName);
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

  function onStyleLoad() {
    applySatelliteVisibility();
    applyRadiation();
    syncGoogleAttribution();
  }

  function applyLoadedStyle(json, theme) {
    var baked = radiationURL || placeholderTiles();
    var style = composeStyle(json, theme, baked);
    appliedTheme = theme;
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
    var requested = themeNow();
    if (requested === appliedTheme && gl) {
      applySatelliteVisibility();
      applyRadiation();
      syncGoogleAttribution();
      return;
    }
    if (requested === inflightTheme) return;
    inflightTheme = requested;
    var gen = ++generation;
    loadStyle(requested).then(function (json) {
      if (gen !== generation) return;
      inflightTheme = null;
      if (themeNow() !== requested) {
        loadAndApply();
        return;
      }
      try {
        applyLoadedStyle(json, requested);
      } catch (err) {
        console.warn('reference basemap rejected style', err);
        enableRasterFallback();
      }
    }).catch(function (err) {
      if (gen !== generation) return;
      inflightTheme = null;
      console.warn('reference basemap unavailable', err);
      enableRasterFallback();
    });
  }

  var api = {
    status: 'loading',
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
