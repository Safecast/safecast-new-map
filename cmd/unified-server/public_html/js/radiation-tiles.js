// Overview radiation tiles. Through zoom 13 the map paints one PNG per XYZ tile
// instead of a Leaflet marker per measurement. Closer zoom, track view, spectra,
// and realtime sensors stay on the existing marker path.
(function () {
  var POINT_ZOOM = 14;

  function tilesWanted() {
    if (window.isTrackView || window.currentTrackID || window.currentTrackIDs) return false;
    if (!window.map || window.map.getZoom() >= POINT_ZOOM) return false;
    if (typeof shouldShowHomeOverlay === 'function' && !window.homeOverlayDismissed && shouldShowHomeOverlay()) return false;
    return true;
  }

  function speedToken() {
    // Same bands as the point layer. Planes are off unless the user checks them,
    // because those readings are higher and would otherwise paint the overview.
    if (typeof loadSpeedFilterState !== 'function') return 'cw';
    var st = loadSpeedFilterState();
    var s = '';
    if (st.plane) s += 'a';
    if (st.car) s += 'c';
    if (st.ped) s += 'w';
    return s || 'none';
  }

  function tileQuery() {
    var parts = [];
    if (typeof loadDateRangeState === 'function') {
      var range = loadDateRangeState();
      if (range && range.length >= 2 && range[0] != null && range[1] != null) {
        parts.push('from=' + encodeURIComponent(range[0]));
        parts.push('to=' + encodeURIComponent(range[1]));
      }
    }
    parts.push('speed=' + speedToken());
    return '?' + parts.join('&');
  }

  function tileURL() {
    return '/tiles/aggregate/{z}/{x}/{y}.png' + tileQuery();
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
    var original = window.updateMarkers;
    var lastSpan = '';
    var spanToken = 0;

    // The date slider on the marker map is filled from the dates of the
    // measurements in view. Overview tiles replace that stream, so ask for
    // the same span and hand it to the existing slider.
    function refreshDateSlider() {
      if (typeof syncDateSliderToMarkerRange !== 'function' || typeof window.__setDateSliderVisibility !== 'function') return;
      var q = boundsQuery();
      if (q === lastSpan) return;
      lastSpan = q;
      var token = ++spanToken;
      fetch('/tiles/aggregate/span?' + q).then(function (r) { return r.ok ? r.json() : null; }).then(function (span) {
        if (token !== spanToken || !tilesWanted()) return;
        if (!span || !span.ok) {
          syncDateSliderToMarkerRange({ hasHistoricalMarkers: false });
          return;
        }
        syncDateSliderToMarkerRange({ minTs: span.min, maxTs: span.max, hasHistoricalMarkers: true });
      }).catch(function () {});
    }

    // Point markers and the overview tiles are exclusive. Entering the overview
    // drops whatever the marker stream already drew, including markers from a
    // closer zoom. Leaving it removes the tiles before those markers return.
    function dropPointLayer() {
      window.lastViewportKey = null;
      if (typeof window.markerLoadToken === 'number') window.markerLoadToken++;
      if (typeof closeMarkerStreams === 'function') closeMarkerStreams();
      if (typeof cancelPendingFlush === 'function') cancelPendingFlush();
      if (typeof clearRenderedMarkers === 'function') clearRenderedMarkers();
    }

    var layerZoom = null;

    function showTiles() {
      dropPointLayer();
      refreshDateSlider();
      var url = tileURL();
      if (!layer) {
        layer = L.tileLayer(url, {
          minZoom: 0,
          maxZoom: POINT_ZOOM - 1,
          opacity: 1,
          zIndex: 250,
          updateWhenIdle: true,
          updateWhenZooming: false,
          keepBuffer: 0
        });
        layer._radiationQuery = tileQuery();
      } else if (layer._radiationQuery !== tileQuery()) {
        layer._radiationQuery = tileQuery();
        layer.setUrl(url);
      }
      var already = map.hasLayer(layer);
      if (!already) layer.addTo(map);
      // A zoom change can leave the previous zoom's tiles scaled in place,
      // which shows up as a hard-edged block of the wrong size.
      if (already && layerZoom !== map.getZoom()) layer.redraw();
      layerZoom = map.getZoom();
      if (!map.hasLayer(sparse)) sparse.addTo(map);
      refreshSparse();
      var loading = document.getElementById('loadingOverlay');
      if (loading) loading.style.display = 'none';
    }

    function hideTiles() {
      if (layer && map.hasLayer(layer)) map.removeLayer(layer);
      layerZoom = null;
      if (map.hasLayer(sparse)) {
        sparse.clearLayers();
        map.removeLayer(sparse);
      }
      lastSpan = '';
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
      var url = '/tiles/aggregate/probe?lat=' + ev.latlng.lat + '&lon=' + ev.latlng.lng + '&z=' + map.getZoom() + tileQuery().replace('?', '&');
      fetch(url).then(function (r) { return r.ok ? r.json() : null; }).then(function (hit) {
        if (hit && hit.found && hit.trackID && typeof viewTrack === 'function') viewTrack(hit.trackID);
      }).catch(function () {});
    });

    map.on('zoom', function () {
      if (map.getZoom() >= POINT_ZOOM) hideTiles();
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
