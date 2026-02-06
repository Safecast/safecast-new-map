# Patroni PostgreSQL High-Availability Cluster Setup Plan

## Overview

This plan covers deploying a 3-node Patroni cluster for a 100GB+ PostgreSQL database with automatic failover, streaming replication, and self-healing capabilities.

---

## Architecture

```
                    ┌──────────────┐
                    │   HAProxy    │
                    │  (Load Bal)  │
                    └──────┬───────┘
                           │
            ┌──────────────┼──────────────┐
            │              │              │
     ┌──────┴──────┐ ┌────┴────────┐ ┌───┴───────┐
     │   Node 1    │ │   Node 2    │ │  Node 3   │
     │  (Primary)  │ │  (Replica)  │ │ (Replica) │
     │ PostgreSQL  │ │ PostgreSQL  │ │ PostgreSQL│
     │  + Patroni  │ │  + Patroni  │ │ + Patroni │
     └──────┬──────┘ └─────┬───────┘ └─────┬─────┘
            │              │               │
            └──────────────┼───────────────┘
                           │
                    ┌──────┴───────┐
                    │    etcd      │
                    │  (3-node     │
                    │   cluster)   │
                    └──────────────┘
```

### Components

| Component | Purpose | Minimum Nodes |
|-----------|---------|---------------|
| PostgreSQL | Database engine | 3 |
| Patroni | HA orchestration, failover management | 3 (co-located with PostgreSQL) |
| etcd | Distributed consensus store (leader election) | 3 (can co-locate or separate) |
| HAProxy | Connection routing (write → primary, read → replicas) | 2 (for HA) |
| pgBackRest | Backup and WAL archiving | 1 (backup server) |

---

## Hardware Requirements

### Per PostgreSQL/Patroni Node

| Resource | Minimum | Recommended |
|----------|---------|-------------|
| CPU | 4 cores | 8 cores |
| RAM | 8 GB | 16 GB |
| Storage (data) | 200 GB SSD | 400 GB NVMe |
| Storage (WAL) | 50 GB SSD | 100 GB NVMe (separate disk) |
| Network | 1 Gbps | 10 Gbps |

> **Note:** For a 100GB database, `shared_buffers` should be ~4GB and total RAM should be at least 2x the active dataset size for effective caching.

### Per etcd Node (if separate)

| Resource | Minimum |
|----------|---------|
| CPU | 2 cores |
| RAM | 4 GB |
| Storage | 20 GB SSD (low-latency critical) |

---

## Phase 1: Infrastructure Preparation

### 1.1 Server Provisioning

Provision 3 servers (physical or virtual) across different failure domains:

```
node1.db.example.com  — 192.168.1.11  (Rack A / AZ-a)
node2.db.example.com  — 192.168.1.12  (Rack B / AZ-b)
node3.db.example.com  — 192.168.1.13  (Rack C / AZ-c)
```

### 1.2 OS Setup (all nodes)

- OS: Ubuntu 22.04 LTS or Rocky Linux 9
- Disable swap (or set `vm.swappiness=1`)
- Set timezone to UTC
- Enable NTP synchronization (chrony recommended)
- Configure firewall rules:

| Port | Service | Direction |
|------|---------|-----------|
| 5432 | PostgreSQL | Inbound from app/HAProxy |
| 8008 | Patroni REST API | Inbound from HAProxy and other nodes |
| 2379-2380 | etcd client/peer | Inbound from cluster nodes |

### 1.3 Kernel Tuning

Add to `/etc/sysctl.conf`:

```ini
vm.swappiness = 1
vm.overcommit_memory = 2
vm.overcommit_ratio = 80
net.core.somaxconn = 65535
net.ipv4.tcp_keepalive_time = 60
net.ipv4.tcp_keepalive_intvl = 10
net.ipv4.tcp_keepalive_probes = 6
```

---

## Phase 2: etcd Cluster Setup

### 2.1 Install etcd (all 3 nodes)

```bash
sudo apt install etcd
# or for Rocky Linux:
# sudo dnf install etcd
```

### 2.2 Configure etcd

Example for Node 1 (`/etc/etcd/etcd.conf.yml`):

```yaml
name: etcd1
data-dir: /var/lib/etcd/data
listen-client-urls: http://192.168.1.11:2379,http://127.0.0.1:2379
advertise-client-urls: http://192.168.1.11:2379
listen-peer-urls: http://192.168.1.11:2380
initial-advertise-peer-urls: http://192.168.1.11:2380
initial-cluster: >-
  etcd1=http://192.168.1.11:2380,
  etcd2=http://192.168.1.12:2380,
  etcd3=http://192.168.1.13:2380
initial-cluster-state: new
initial-cluster-token: patroni-etcd-cluster
```

Repeat for Node 2 and Node 3, adjusting `name`, IP addresses, and URLs.

### 2.3 Start and Verify

```bash
sudo systemctl enable --now etcd
etcdctl endpoint health --cluster
```

---

## Phase 3: PostgreSQL and Patroni Installation

### 3.1 Install PostgreSQL (all 3 nodes)

```bash
# Add PostgreSQL APT repo
sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" \
  > /etc/apt/sources.list.d/pgdg.list'
curl -fsSL https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo gpg --dearmor \
  -o /etc/apt/trusted.gpg.d/postgresql.gpg
sudo apt update
sudo apt install postgresql-16
```

Stop the default PostgreSQL instance — Patroni will manage it:

```bash
sudo systemctl stop postgresql
sudo systemctl disable postgresql
```

### 3.2 Install Patroni (all 3 nodes)

```bash
sudo apt install python3-pip python3-psycopg2
sudo pip3 install patroni[etcd] --break-system-packages
```

### 3.3 Configure Patroni

Create `/etc/patroni/patroni.yml` on each node.

Example for Node 1:

```yaml
scope: pg-cluster
namespace: /db/
name: node1

restapi:
  listen: 0.0.0.0:8008
  connect_address: 192.168.1.11:8008

etcd3:
  hosts:
    - 192.168.1.11:2379
    - 192.168.1.12:2379
    - 192.168.1.13:2379

bootstrap:
  dcs:
    ttl: 30
    loop_wait: 10
    retry_timeout: 10
    maximum_lag_on_failover: 1048576  # 1MB
    postgresql:
      use_pg_rewind: true
      use_slots: true
      parameters:
        max_connections: 200
        shared_buffers: 4GB
        effective_cache_size: 12GB
        maintenance_work_mem: 512MB
        wal_buffers: 64MB
        work_mem: 32MB
        max_wal_senders: 10
        max_replication_slots: 10
        hot_standby: "on"
        wal_level: replica
        wal_log_hints: "on"
        wal_keep_size: 4096  # MB
        archive_mode: "on"
        archive_command: >-
          pgbackrest --stanza=main archive-push %p
        logging_collector: "on"
        log_directory: /var/log/postgresql
        log_filename: postgresql-%Y-%m-%d.log
        log_min_duration_statement: 1000

  initdb:
    - encoding: UTF8
    - data-checksums

  pg_hba:
    - host replication replicator 192.168.1.0/24 scram-sha-256
    - host all all 0.0.0.0/0 scram-sha-256

  users:
    admin:
      password: "CHANGE_ME_ADMIN"
      options:
        - createrole
        - createdb
    replicator:
      password: "CHANGE_ME_REPLICATOR"
      options:
        - replication

postgresql:
  listen: 0.0.0.0:5432
  connect_address: 192.168.1.11:5432
  data_dir: /var/lib/postgresql/16/main
  bin_dir: /usr/lib/postgresql/16/bin
  pgpass: /tmp/pgpass0
  authentication:
    superuser:
      username: postgres
      password: "CHANGE_ME_POSTGRES"
    replication:
      username: replicator
      password: "CHANGE_ME_REPLICATOR"
    rewind:
      username: postgres
      password: "CHANGE_ME_POSTGRES"

tags:
  nofailover: false
  noloadbalance: false
  clonefrom: false
  nosync: false
```

For Node 2 and Node 3, change `name` and `connect_address` accordingly.

### 3.4 Create Systemd Service

Create `/etc/systemd/system/patroni.service`:

```ini
[Unit]
Description=Patroni PostgreSQL HA
After=syslog.target network.target etcd.service
Wants=network-online.target

[Service]
Type=simple
User=postgres
Group=postgres
ExecStart=/usr/local/bin/patroni /etc/patroni/patroni.yml
ExecReload=/bin/kill -s HUP $MAINPID
KillMode=process
TimeoutSec=30
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

### 3.5 Start the Cluster

```bash
# Start Node 1 first (becomes initial primary)
sudo systemctl enable --now patroni

# Then start Node 2 and Node 3
# They will automatically clone from the primary
sudo systemctl enable --now patroni
```

### 3.6 Verify Cluster Status

```bash
patronictl -c /etc/patroni/patroni.yml list
```

Expected output:

```
+--------+-----------+---------+---------+----+-----------+
| Member |   Host    |  Role   |  State  | TL | Lag in MB |
+--------+-----------+---------+---------+----+-----------+
| node1  | 192.168.1.11 | Leader  | running |  1 |           |
| node2  | 192.168.1.12 | Replica | running |  1 |       0.0 |
| node3  | 192.168.1.13 | Replica | running |  1 |       0.0 |
+--------+-----------+---------+---------+----+-----------+
```

---

## Phase 4: HAProxy Setup

### 4.1 Install HAProxy

On one or two dedicated load balancer nodes:

```bash
sudo apt install haproxy
```

### 4.2 Configure HAProxy

`/etc/haproxy/haproxy.cfg`:

```
global
    maxconn 1000
    log /dev/log local0

defaults
    log global
    mode tcp
    retries 3
    timeout client 30m
    timeout connect 4s
    timeout server 30m
    timeout check 5s

listen stats
    mode http
    bind *:7000
    stats enable
    stats uri /

# Primary (read-write) connections
listen postgresql-primary
    bind *:5000
    option httpchk GET /primary
    http-check expect status 200
    default-server inter 3s fall 3 rise 2 on-marked-down shutdown-sessions
    server node1 192.168.1.11:5432 maxconn 100 check port 8008
    server node2 192.168.1.12:5432 maxconn 100 check port 8008
    server node3 192.168.1.13:5432 maxconn 100 check port 8008

# Replicas (read-only) connections
listen postgresql-replicas
    bind *:5001
    balance roundrobin
    option httpchk GET /replica
    http-check expect status 200
    default-server inter 3s fall 3 rise 2 on-marked-down shutdown-sessions
    server node1 192.168.1.11:5432 maxconn 100 check port 8008
    server node2 192.168.1.12:5432 maxconn 100 check port 8008
    server node3 192.168.1.13:5432 maxconn 100 check port 8008
```

### 4.3 Application Connection Strings

```
# Read-write (always goes to primary)
postgresql://admin:password@haproxy-host:5000/mydb

# Read-only (load balanced across replicas)
postgresql://admin:password@haproxy-host:5001/mydb
```

---

## Phase 5: Backup Configuration (pgBackRest)

### 5.1 Install pgBackRest (all nodes + backup server)

```bash
sudo apt install pgbackrest
```

### 5.2 Configure pgBackRest

`/etc/pgbackrest/pgbackrest.conf` (on all PostgreSQL nodes):

```ini
[main]
pg1-path=/var/lib/postgresql/16/main
pg1-port=5432
pg1-user=postgres

[global]
repo1-path=/var/lib/pgbackrest
repo1-retention-full=2
repo1-retention-diff=7
repo1-cipher-type=aes-256-cbc
repo1-cipher-pass=CHANGE_ME_CIPHER_KEY
process-max=4
compress-type=zst
compress-level=3
start-fast=y
delta=y
```

### 5.3 Initialize and Schedule Backups

```bash
# Initialize the backup stanza
sudo -u postgres pgbackrest --stanza=main stanza-create

# Full backup (run weekly via cron)
sudo -u postgres pgbackrest --stanza=main --type=full backup

# Differential backup (run daily via cron)
sudo -u postgres pgbackrest --stanza=main --type=diff backup
```

### 5.4 Cron Schedule

```cron
# /etc/cron.d/pgbackrest
0 2 * * 0  postgres  pgbackrest --stanza=main --type=full backup
0 2 * * 1-6  postgres  pgbackrest --stanza=main --type=diff backup
```

---

## Phase 6: Monitoring

### 6.1 Key Metrics to Monitor

| Metric | Alert Threshold |
|--------|----------------|
| Replication lag | > 100MB or > 30 seconds |
| Patroni cluster state | Any node not `running` |
| etcd cluster health | Any member unhealthy |
| Disk usage | > 80% |
| Connection count | > 80% of `max_connections` |
| WAL accumulation | > 10GB |
| Failed backups | Any failure |

### 6.2 Recommended Monitoring Stack (all open source)

- **Prometheus** + **postgres_exporter** — metrics collection
- **Grafana** — dashboards and alerting
- **Patroni REST API** — query `http://node:8008/cluster` for cluster state

### 6.3 Simple Health Check Script

```bash
#!/bin/bash
# patroni-health-check.sh
CLUSTER_STATUS=$(patronictl -c /etc/patroni/patroni.yml list -f json)
LEADER_COUNT=$(echo "$CLUSTER_STATUS" | jq '[.[] | select(.Role=="Leader")] | length')
RUNNING_COUNT=$(echo "$CLUSTER_STATUS" | jq '[.[] | select(.State=="running")] | length')

if [ "$LEADER_COUNT" -ne 1 ]; then
    echo "CRITICAL: Expected 1 leader, found $LEADER_COUNT"
    exit 2
fi

if [ "$RUNNING_COUNT" -lt 2 ]; then
    echo "WARNING: Only $RUNNING_COUNT nodes running"
    exit 1
fi

echo "OK: Cluster healthy - $RUNNING_COUNT nodes, 1 leader"
exit 0
```

---

## Phase 7: Migration of Existing Data

### 7.1 Migration Steps

1. **Set up the Patroni cluster** with an empty database (Phases 1–4).
2. **Stop writes** to the source database (schedule maintenance window).
3. **Dump the source database:**
   ```bash
   pg_dump -h old-server -U postgres -Fc -j 4 mydb > mydb.dump
   ```
4. **Restore to the Patroni primary:**
   ```bash
   pg_restore -h haproxy-host -p 5000 -U admin -d mydb -j 4 mydb.dump
   ```
5. **Verify** data integrity and replication lag on all nodes.
6. **Switch application** connection strings to HAProxy endpoints.

> For 100GB, expect the dump/restore to take approximately 30–60 minutes depending on hardware and data complexity.

### 7.2 Minimal Downtime Alternative

Use **logical replication** to migrate with near-zero downtime:

1. Set up the Patroni cluster.
2. Configure logical replication from old server → Patroni primary.
3. Let it sync fully.
4. Switch application connections.
5. Drop logical replication.

---

## Operational Runbook

### Manual Failover

```bash
patronictl -c /etc/patroni/patroni.yml switchover --master node1 --candidate node2
```

### Restart a Node

```bash
patronictl -c /etc/patroni/patroni.yml restart pg-cluster node1
```

### Reinitialize a Failed Node

```bash
patronictl -c /etc/patroni/patroni.yml reinit pg-cluster node3
```

### Check Replication Status

```sql
-- Run on primary
SELECT client_addr, state, sent_lsn, write_lsn, flush_lsn,
       replay_lsn, sync_state
FROM pg_stat_replication;
```

---

## Security Checklist

- [ ] Change all default passwords in `patroni.yml`
- [ ] Enable TLS for PostgreSQL connections (`ssl = on`)
- [ ] Enable TLS for etcd communication
- [ ] Enable TLS for Patroni REST API
- [ ] Restrict `pg_hba.conf` to known application IPs
- [ ] Use SCRAM-SHA-256 authentication (not MD5)
- [ ] Enable data checksums (included in config above)
- [ ] Encrypt backups (included in pgBackRest config)
- [ ] Restrict firewall to required ports only
- [ ] Set up audit logging (`pgaudit` extension)

---

## Estimated Timeline

| Phase | Task | Duration |
|-------|------|----------|
| 1 | Infrastructure provisioning and OS setup | 1–2 days |
| 2 | etcd cluster | 0.5 day |
| 3 | PostgreSQL + Patroni | 1 day |
| 4 | HAProxy | 0.5 day |
| 5 | Backup configuration | 0.5 day |
| 6 | Monitoring setup | 1 day |
| 7 | Data migration + testing | 1–2 days |
| — | **Total** | **5–7 days** |

---

## References

- [Patroni Documentation](https://patroni.readthedocs.io/)
- [etcd Documentation](https://etcd.io/docs/)
- [pgBackRest Documentation](https://pgbackrest.org/)
- [PostgreSQL Replication Docs](https://www.postgresql.org/docs/16/high-availability.html)
- [HAProxy Documentation](https://www.haproxy.org/#docs)
