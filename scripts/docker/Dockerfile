FROM postgres:16-bookworm

# ----- утилиты + бинарь приложения -------------------------------
RUN apt-get update && \
    apt-get install -y --no-install-recommends curl ca-certificates gosu && \
    curl -L https://github.com/Safecast/safecast-new-map/releases/download/latest/safecast-new-map_linux_amd64 \
      -o /usr/local/bin/safecast-new-map && \
    chmod +x /usr/local/bin/safecast-new-map && \
    rm -rf /var/lib/apt/lists/*

# ----- ENV по умолчанию -----------------------------------------
ENV DB_NAME=safecast_new_map \
    DB_USER=safecast_new_map \
    PGDATA=/var/lib/postgresql/data \
    DEFAULT_LAT=51.389 \
    DEFAULT_LON=30.099 \
    DEFAULT_ZOOM=4 \
    DEFAULT_LAYER="Google Satellite"

# ----- entrypoint + директории ----------------------------------
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

VOLUME ["/var/lib/postgresql/data", "/certs"]

EXPOSE 8765 80 443
ENTRYPOINT ["/entrypoint.sh"]

