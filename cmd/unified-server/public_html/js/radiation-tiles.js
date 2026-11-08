// Overview radiation tiles. Below zoom 13 the map paints one PNG per XYZ tile
// instead of a Leaflet marker per measurement. Close zoom, track view, spectra,
// and realtime sensors stay on the existing marker path.
(function () {
  var POINT_ZOOM = 13;

  function tilesWanted() {
    if (window.isTrackView || window.currentTrackID || window.currentTrackIDs) return false;
    if (!window.map || window.map.getZoom() >= POINT_ZOOM) return false;
    if (typeof shouldShowHomeOverlay === 'function' && !window.homeOverlayDismissed && shouldShowHomeOverlay()) return false;
    return true;
  }

  function dateQuery() {
    if (typeof loadDateRangeState !== 'function') return '';
    var range = loadDateRangeState();
    if (!range || range.length < 2 || range[0] == null || range[1] == null) return '';
    return '?from=' + encodeURIComponent(range[0]) + '&to=' + encodeURIComponent(range[1]);
  }

  function tileURL() {
    return '/tiles/aggregate/{z}/{x}/{y}.png' + dateQuery();
  }

  function doseColor(dose) {
    if (typeof getGradientColor === 'function') return getGradientColor(dose);
    if (dose >= 1) return '#ff0000';
    if (dose >= 0.3) return '#ff8c00';
    return '#1e90ff';
  }

  function boundsQuery() {
    var b = window.map.getBounds();
    var sw = b.getSouthWest();
    var ne = b.getNorthEast();
    return 'minLat=' + sw.lat + '&minLon=' + sw.lng + '&maxLat=' + ne.lat + '&maxLon=' + ne.lng;
  }

  window.installRadiationTiles = function (map) {
    var layer = null;
    var sparse = L.layerGroup();
    var note = null;
    var original = window.updateMarkers;

    function ensureNote() {
      var legend = document.getElementById('legend');
      if (!legend || note) return;
      note = document.createElement('div');
      note.id = 'radiation-tile-note';
      note.style.cssText = 'font-size:11px;line-height:1.3;padding:4px 6px 6px;max-width:160px;';
      note.textContent = 'Highest measurement in each pixel. Blank means no measurement at this zoom.';
      legend.appendChild(note);
    }

    function showTiles() {
      ensureNote();
      if (note) note.style.display = '';
      var url = tileURL();
      if (!layer) {
        layer = L.tileLayer(url, {
          minZoom: 0,
          maxZoom: POINT_ZOOM - 1,
          maxNativeZoom: POINT_ZOOM - 1,
          opacity: 1,
          zIndex: 250,
          updateWhenIdle: false
        });
        layer._radiationQuery = dateQuery();
      } else if (layer._radiationQuery !== dateQuery()) {
        layer._radiationQuery = dateQuery();
        layer.setUrl(url);
      }
      if (!map.hasLayer(layer)) layer.addTo(map);
      if (!map.hasLayer(sparse)) sparse.addTo(map);
      refreshSparse();
      var loading = document.getElementById('loadingOverlay');
      if (loading) loading.style.display = 'none';
    }

    function hideTiles() {
      if (layer && map.hasLayer(layer)) map.removeLayer(layer);
      if (map.hasLayer(sparse)) {
        sparse.clearLayers();
        map.removeLayer(sparse);
      }
      if (note) note.style.display = 'none';
    }

    var lastSparse = '';

    function refreshSparse() {
      if (!tilesWanted()) return;
      var q = boundsQuery();
      if (q === lastSparse) return;
      lastSparse = q;
      sparse.clearLayers();
      if (typeof loadSpeedFilterState !== 'function' || loadSpeedFilterState().spectrum !== false) {
        fetch('/api/markers/spectra?' + q).then(function (r) { return r.ok ? r.json() : []; }).then(drawSpectra).catch(function () {});
      }
      if (window.safecastRealtimeEnabled) {
        fetch('/tiles/live?' + q).then(function (r) { return r.ok ? r.json() : []; }).then(drawLive).catch(function () {});
      }
    }

    function drawSpectra(rows) {
      if (!tilesWanted() || !rows) return;
      rows.forEach(function (m) {
        if (m == null || m.lat == null || m.lon == null) return;
        var marker = L.marker([m.lat, m.lon], {
          icon: L.divIcon({
            className: '',
            html: '<div style="background:' + doseColor(m.doseRate) + ';width:18px;height:18px;border:2px solid #fff;display:flex;align-items:center;justify-content:center;font-size:11px;">S</div>',
            iconSize: [18, 18],
            iconAnchor: [9, 9]
          })
        });
        marker.on('click', function (ev) {
          L.DomEvent.stopPropagation(ev);
          if (typeof openSpectrumModal === 'function') openSpectrumModal(m.id, m);
        });
        sparse.addLayer(marker);
      });
    }

    function drawLive(rows) {
      if (!tilesWanted() || !rows) return;
      rows.forEach(function (m) {
        if (m == null || m.lat == null || m.lon == null) return;
        var marker = L.circleMarker([m.lat, m.lon], {
          radius: 7,
          color: '#ffffff',
          weight: 2,
          fillColor: doseColor(m.doseRate),
          fillOpacity: 1
        });
        marker.on('click', function (ev) {
          L.DomEvent.stopPropagation(ev);
        });
        if (m.deviceName || m.doseRate != null) {
          marker.bindTooltip((m.deviceName || 'Live') + ' ' + m.doseRate + ' µSv/h');
        }
        sparse.addLayer(marker);
      });
    }

    map.on('click', function (ev) {
      if (!tilesWanted()) return;
      var target = ev.originalEvent && ev.originalEvent.target;
      if (target && target.closest && target.closest('.leaflet-marker-icon, .leaflet-interactive')) {
        if (target.closest('.leaflet-marker-icon')) return;
      }
      var url = '/tiles/aggregate/probe?lat=' + ev.latlng.lat + '&lon=' + ev.latlng.lng + '&z=' + map.getZoom() + dateQuery().replace('?', '&');
      fetch(url).then(function (r) { return r.ok ? r.json() : null; }).then(function (hit) {
        if (hit && hit.found && hit.trackID && typeof viewTrack === 'function') viewTrack(hit.trackID);
      }).catch(function () {});
    });

    window.updateMarkers = function () {
      if (tilesWanted()) {
        showTiles();
        return;
      }
      hideTiles();
      if (typeof original === 'function') return original.apply(this, arguments);
    };
  };
})();
