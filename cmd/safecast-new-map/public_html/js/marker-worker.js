/**
 * Marker Parser Web Worker
 * Handles JSON and MessagePack parsing off the main thread to keep UI responsive
 * during heavy marker data loads.
 */

// Load MessagePack library for binary decoding (~50KB)
let hasMsgpack = false;
try {
  importScripts('/js/msgpack.min.js');
  hasMsgpack = typeof msgpack !== 'undefined';
  if (hasMsgpack) {
    console.log('[Worker] MessagePack loaded successfully');
  }
} catch (e) {
  console.warn('[Worker] MessagePack not available, using JSON only');
}

self.onmessage = function(e) {
  const { type, data, streamId, index, format } = e.data;

  if (type === 'parse') {
    try {
      let marker;

      if (format === 'msgpack' && hasMsgpack) {
        // Parse MessagePack binary data
        marker = msgpack.decode(new Uint8Array(data));
      } else {
        // Parse JSON string
        marker = JSON.parse(data);
      }

      // Basic validation - reject malformed markers early
      if (typeof marker.lat !== 'number' || typeof marker.lon !== 'number') {
        return; // Silently skip invalid markers
      }

      // Validate and normalize required numeric fields
      if (typeof marker.doseRate !== 'number') {
        marker.doseRate = 0;
      }
      if (typeof marker.speed !== 'number') {
        marker.speed = 0;
      }
      if (typeof marker.date !== 'number') {
        marker.date = 0;
      }

      // Send parsed marker back to main thread
      self.postMessage({
        type: 'marker',
        marker: marker,
        streamId: streamId,
        index: index
      });
    } catch (err) {
      // Parse error - skip this marker silently
      // Uncomment for debugging:
      // self.postMessage({ type: 'error', error: err.message, streamId, index });
    }
  } else if (type === 'parseBatch') {
    // Batch parsing for better performance with multiple messages
    const results = [];
    const items = Array.isArray(data) ? data : [data];

    for (let i = 0; i < items.length; i++) {
      try {
        let marker;

        if (format === 'msgpack' && hasMsgpack) {
          marker = msgpack.decode(new Uint8Array(items[i]));
        } else {
          marker = JSON.parse(items[i]);
        }

        // Validate required fields
        if (typeof marker.lat === 'number' && typeof marker.lon === 'number') {
          // Normalize optional fields
          if (typeof marker.doseRate !== 'number') marker.doseRate = 0;
          if (typeof marker.speed !== 'number') marker.speed = 0;
          if (typeof marker.date !== 'number') marker.date = 0;

          results.push(marker);
        }
      } catch (err) {
        // Skip invalid entries silently
      }
    }

    self.postMessage({
      type: 'markerBatch',
      markers: results,
      streamId: streamId
    });
  } else if (type === 'ping') {
    // Health check - respond with worker status
    self.postMessage({
      type: 'pong',
      hasMsgpack: hasMsgpack
    });
  }
};
