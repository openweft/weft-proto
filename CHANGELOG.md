# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project aims to adhere to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

Nothing yet.

## [v0.22.0] — 2026-08-18

First tag since v0.21.0. Additions only against that tag: no message or RPC was
renamed or removed, and no field number was reused. (The `labels` → `properties`
rename recorded under Unreleased shipped in v0.21.0 and is not part of this
release.)

### Added
- `VMInfo.status` + `SetVMStatus` RPC (`SetVMStatusRequest{uuid, status}` /
  `SetVMStatusResponse`) : the operator's administrative intent for a VM, as
  distinct from its observed runtime state. This is what weft-tui needs to show
  and flip that intent.
- `VMInfo.restart_count` + `VMInfo.max_restarts` : the k8s-style RESTARTS
  column, so a supervisor's retry budget is visible on the wire rather than
  inferred.
- `RestartVMRequest.host_uuid` : lets a restart be dispatched across hosts
  instead of only to the one holding the connection.
- `RegisterMicroVMRequest` / `RegisterMicroVMOp` : `cpu` and `mem_mb`
  workload-shape fields, and an `image` field carrying the OCI reference.
- `UninstallPlugin` RPC (`UninstallPluginRequest{name, instance_uuid}` /
  `UninstallPluginResponse`) : the counterpart to `InstallPlugin`.
- `EnablePlugin` / `DisablePlugin` RPCs and `PluginInstance.disabled` : a soft
  admin-state toggle that leaves the install side-effects in place.
- `Flavor.uuid` promoted to a stable identifier.

### Changed
- The `go` directive moves to 1.26.4, the fleet floor.
- The `wire-compat` CI gate was dropped during development, deliberately: while
  the schema moved quickly every legitimate refactor flagged as breaking. This
  tag is the stable baseline that commit asked for, so the gate can be
  reinstated against **v0.22.0** rather than against a floating latest.

## [v0.21.0] — 2026-06-22

### Added
- `VMState` extended with five lifecycle states that previously had no wire
  representation: `VM_STATE_CREATED` (5) for a vmDir provisioned and
  registered but never started, `VM_STATE_STARTING` (6) between an accepted
  `StartVM` and the guest's first heartbeat, `VM_STATE_STOPPING` (7) while
  qemu/vz drains, `VM_STATE_ZOMBIE` (8) for a process that died without
  going through `StopVM`, and `VM_STATE_DELETING` (9) so the
  operator-visible event survives long enough to flush. `VM_STATE_ZOMBIE`
  is what `zombiegc` surfaces as `weft_vm_zombies`.

## [v0.20.0] — 2026-06-22

### Added
- Host CPU / RAM / GPU inventory on `HostInfo` and `HostRegistration`:
  `cpu_count`, `memory_mib`, and `repeated GPU gpus`.
- `GPU` message — `vendor` / `model` / `memory_gib` / `mig_capable`, the
  last for datacenter SKUs that support MIG slicing.

## [v0.19.0] — 2026-06-22

### Added
- Host OS, kernel, network and storage facts on `HostInfo` and
  `HostRegistration`: `os_id`, `os_version`, `os_pretty`, `kernel_version`,
  `repeated NetworkInterface network_interfaces` and
  `repeated StorageMount storage_mounts`.

## [v0.18.0] — 2026-06-22

### Added
- `HostInfo.agent_version` + `HostInfo.driver_versions`, mirrored on
  `HostRegistration`, so the control plane can see what each host is
  actually running rather than inferring it.
- `GetClusterInfoResponse.control_plane_host_uuids`.

## [v0.17.0] — 2026-06-22

### Added
- `GetClusterInfo` / `SetClusterName` RPCs +
  `GetClusterInfoResponse{cluster_name, local_host_uuid}`.

### Fixed
- CI: the `wire-compat` job (buf breaking vs the latest tag) had a YAML
  heredoc parse error and was not running.

## [v0.16.0] — 2026-06-20

### Added
- `GuestHello.reported_cid` (uint32, field 5) — guest-side AF_VSOCK CID
  readback, so the host learns the CID the guest actually got rather than
  the one it asked for.

## [v0.15.0] — 2026-06-20

### Added
- `SetPodSpec` / `GetPodSpec` RPCs on `WeftAgent`. The spec travels as a
  protojson-encoded `spec_json` blob rather than an imported guestv1
  message, to keep guestv1 out of weft.proto.
- `GetMicroVMMetrics` RPC + `MicroVMMetricsResponse` — per-VM cpu / mem /
  net / disk + uptime snapshot.

## [v0.14.0] — 2026-06-20

### Added
- Generated Go bindings for `agentv1` and `guestv1` are now committed
  alongside `weftv1`, so consumers no longer regenerate them locally.
- `UpdateContainer` message (`container_id`, `command`, `env`) in
  guest.proto.

## [v0.13.0] — 2026-06-19

### Added
- `SetProjectTenant` RPC + `ProjectInfo.tenant_uuid` (field 4), completing
  the project surface alongside List / Create / Rename / Delete.

## [v0.12.0] — 2026-06-19

### Added
- `RestartVM` RPC + `VMInfo.host_uuid` (field 13).

### Changed (breaking)
- `labels` → `properties` across the host/VM scheduler surface. Field
  numbers are preserved, so binary payloads stay wire-compatible; the RPC
  and message names change: `SetHostLabels` → `SetHostProperties`,
  `SetVMLabels` → `SetVMProperties`, `HostInfo.labels` (10) →
  `properties`, `VMInfo.labels` (12) → `properties`,
  `RegisterHostRequest.labels` (10), `HostRegistration.labels`
  (agent.proto, 10) and `PodSpec.labels` (guest.proto, 5) likewise. One
  name across the openweft stack, matching `cluster.hcl`'s
  `Host.properties`.

## [v0.11.6] — 2026-06-14

### Added
- `PortInfo` message + `ListPortsForVM` RPC : read-only view of a
  VM's NIC bindings (MAC/IP/security-groups/created-at + reserved
  ingress/egress Mbps for the upcoming portqos persistence). Powers
  the webui Network panel and the (future) Port detail drawer.

## [v0.11.5] — 2026-06-14

### Added
- `FloatingIPInfo.rate_limit_pps` + `MapFloatingIPRequest.rate_limit_pps` :
  anti-DDoS PPS cap on inbound traffic to a floating IP. 0 means
  no cap ; >100k clamps. Persisted via the weft adapter, consumed
  by the floatingipnat reconciler on the host side.

## [v0.11.4] — 2026-06-14

### Fixed
- Regenerated the `AttestationService` stubs. The v0.11.3 merge kept the
  `.proto` change but carried stale `.pb.go` files, so the service was
  declared and not generated.

## [v0.11.3] — 2026-06-14

### Added
- `AttestationService` — node admission over opaque attest blobs:
  `Enroll`, `CompleteEnroll`, `RequestAdmission`, `Admit`. The payloads
  travel as `AttestMsg{payload, ak_name}` / `AttestResult{ok, reason}` /
  `AdmitResult{granted, reason, ak_name}` so weft-proto carries no
  attestation types of its own; `ak_name` is the AK Name key used to route
  `CompleteEnroll` and `Admit`.

## [v0.11.2] — 2026-06-14

### Added
- `NetworkInfo.external_mode` / `vlan` / `parent_interface` for L2 and
  VLAN floating IPs.

## [v0.11.1] — 2026-06-09

### Added
- `TriggerZombieSweep` RPC, returning the same `GetZombieReportResponse`.
  Closes the `weft instance gc --apply` loop.

## [v0.11.0] — 2026-06-09

### Added
- `GetZombieReport` RPC + `ZombieEntry` — uuid / name / project / host /
  kind / reason / deployment type, plus `detected_at_unix_ns` and
  `host_down_since_unix_ns`.

## [v0.10.1] — 2026-06-09

### Added
- `SetVMLabels` RPC + `VMInfo.labels` (field 12). Renamed to
  `SetVMProperties` / `properties` in v0.12.0.

## [v0.10.0] — 2026-06-08

### Added
- `HealthProbe` on `SchedulingRule` — `NONE` / `HTTP` / `TCP` / `EXEC`
  with `http_path` / `http_port` / `http_method` / `http_status_ok`,
  `tcp_port`, `exec_command`, and `initial_delay_ms` to wait out VM boot
  before the first probe.
- `RespawnPolicy` on `SchedulingRule`.

### Fixed
- Doc: `MapFloatingIP`'s comment described Envoy; the actual data plane is
  the embedded Caddy.

## [v0.9.0] — 2026-06-05

### Added

- **VolumeProperty + Share extensions + Bucket + SSHKeyCatalogue +
  SchedulingRule + RegistryRemote** — closes the Tier 4-6 CLI-vs-webui
  parity gap, six new noun families on `WeftAgent` :

  - `VolumePropertyInfo` + 3 RPCs (`GetVolumeProperty` /
    `SetVolumeProperty` / `DeleteVolumeProperty`) — mirror of
    VMProperty addressed by `volume_uuid` for block volumes.
  - `ShareInfo` extensions : `GetShare` + `ResizeShare` close
    the v0.8 gap (list/create/delete already shipped).
  - `BucketInfo` + 6 RPCs (`ListBuckets` / `GetBucket` /
    `CreateBucket` / `DeleteBucket` / `GetBucketPolicy` /
    `SetBucketPolicy`) — S3 bucket catalogue ; data lives on the
    S3 endpoint (versitygw / CubeFS objectnode), the agent
    tracks credentials + mutable policy JSON.
  - `SSHKeyCatalogueEntry` + 4 RPCs (`ListSSHKeyCatalogue` /
    `AddSSHKeyCatalogue` / `RemoveSSHKeyCatalogue` /
    `ImportSSHKeyCatalogue`) — cluster-wide named SSH keys VMs
    reference at CreateVM time. Distinct from per-VM
    `weft instance sshkey`.
  - `SchedulingRuleInfo` + 4 RPCs (`ListSchedulingRules` /
    `CreateSchedulingRule` / `UpdateSchedulingRule` /
    `DeleteSchedulingRule`) — per [[openweft_nominal_binding]],
    rules carry selector + target_count + anti_affinity, the
    scheduler reconciles toward target_count.
  - `RegistryRemoteInfo` + 4 RPCs (`ListRegistryRemotes` /
    `SetRegistryRemote` (upsert) / `DeleteRegistryRemote` /
    `SearchRegistryRemote`) — OCI registry alias catalogue ;
    `credential_secret_ref` points at the secret store so the
    row never carries the raw token on the wire.

  22 new RPCs total. Mirror the existing inventory-noun pattern
  (UUID-keyed, partial-PATCH updates, cascade refusal surfaces
  blocking counts on the response when relevant).

## [v0.8.0] — 2026-06-05

### Added

- **Subnet + LoadBalancer + DNSZone + DNSRecord registries on `WeftAgent`**
  — closes the last Tier-3 CLI-vs-webui parity gap by exposing four
  network-plane noun families that were previously webui-only :

  - `SubnetInfo` + 5 RPCs (`ListSubnets` / `GetSubnet` /
    `CreateSubnet` / `UpdateSubnet` / `DeleteSubnet`) — per-network
    IP scopes, parent is `network_uuid`, immutable CIDR, mutable
    gateway + dns_servers. Update carries a `clear_dns_servers`
    bool to disambiguate "keep" vs "clear" on the wire.
  - `LoadBalancerInfo` + 6 RPCs (`ListLoadBalancers` /
    `GetLoadBalancer` / `CreateLoadBalancer` / `UpdateLoadBalancer`
    / `SetLoadBalancerBackends` / `DeleteLoadBalancer`) — project-
    scoped VIPs with `protocol` ∈ {`l4_tcp`, `l4_udp`, `l7_http`,
    `l7_https`} and a `LBBackend{address, weight}` repeated list.
    `SetLoadBalancerBackends` replaces the list atomically (clients
    GET-modify-PUT for single-entry adds/removes). `DeleteLoadBalancer`
    refuses while a FloatingIP still maps to the VIP and surfaces
    `blocked_by_fips` so the operator unmaps it first.
  - `DNSZoneInfo` + 5 RPCs (`ListDNSZones` / `GetDNSZone` /
    `CreateDNSZone` / `UpdateDNSZone` / `DeleteDNSZone`) —
    authoritative apex per project. SOA email + default TTL are
    mutable ; the zone's `records` count is server-derived.
    `DeleteDNSZone` refuses while records still attach and surfaces
    `blocked_by_records`.
  - `DNSRecordInfo` + 4 RPCs (`ListDNSRecords` / `CreateDNSRecord`
    / `UpdateDNSRecord` / `DeleteDNSRecord`) — zone children with
    `type` ∈ {`A`, `AAAA`, `CNAME`, `MX`, `TXT`, `SRV`}, optional
    per-record TTL, MX/SRV priority.

  20 new RPCs total. Mirror the existing inventory-noun pattern
  (UUID-keyed, partial-PATCH updates, cascade refusal surfaces
  blocking counts on the response).

  Why proto and not webui-only : the parity audit identified these
  four ressources as Tier 3 because the CLI couldn't drive
  end-to-end network plumbing without them. With the registry in
  the control plane, `weft subnet create`, `weft loadbalancer
  set-backends`, `weft dns-zone create`, `weft dns-record create`
  all flow through the same Unix socket as every other `weft <noun>
  <verb>` ; the webui can either continue holding its local view or
  migrate onto the live RPC at its own pace.

## [v0.7.0] — 2026-06-05

### Added

- **AvailabilityZone + Rack registry on `WeftAgent`** — elevates AZ
  and Rack from webui-only persistence (previously
  `resourceByID["azs"|"racks"]`) to first-class control-plane RPCs
  so the CLI + every other client reaches the same source of truth.
  - `AZInfo` message : uuid, code (immutable short id), name,
    region, status, created-at, server-derived racks + hosts
    counts.
  - `RackInfo` message : uuid, az_uuid (parent), code, name,
    status, height_u, created-at, server-derived hosts count.
  - 10 new RPCs : `ListAZs` / `GetAZ` / `CreateAZ` / `UpdateAZ` /
    `DeleteAZ` and `ListRacks` / `GetRack` / `CreateRack` /
    `UpdateRack` / `DeleteRack`. `Update*` are partial PATCHes
    (empty string fields = keep current). `Delete*` refuses when
    child rows still bind to the row being deleted ; the response
    surfaces the blocking-count so the operator sees exactly what
    needs draining.

  Why proto and not webui-only : the CLI parity audit surfaced
  AZ/Rack CRUD as a Tier 1 gap because the CLI couldn't drive
  bring-up of a fresh cluster. With the registry in the control
  plane, `weft az create` and `weft rack create` work over the
  Unix socket exactly like every other `weft <noun> <verb>`, and
  the webui can either continue maintaining its local view or
  migrate onto the live RPC.

## [v0.6.0] — 2026-06-05

> **No v0.6.0 tag was ever pushed.** The repository goes v0.5.0 → v0.7.0.
> Everything below shipped to consumers as part of **v0.7.0**:
> `RevertVolumeSnapshot` is absent from `v0.5.0:weft.proto` and present in
> `v0.7.0:weft.proto`. The section is kept, rather than folded into v0.7.0,
> because it is the only place this surface is written up.

### Added

- **Volume snapshot/backup RPC surface** on `WeftAgent` :
  - `RevertVolumeSnapshot` — rolls a block-backend volume back to a captured snapshot (driver dispatches via `weft-block` ; file-backend parents reject with FailedPrecondition).
  - `CreateVolumeBackup` / `ListVolumeBackups` / `DeleteVolumeBackup` / `RestoreVolumeBackup` — off-host backups of block-backend volumes to one of four target schemes (`oci://` recommended, `s3://` for versitygw / CubeFS objectnode, `sftp://` for sftpgo, `fs:///` for dev).
- Supporting messages : `RevertVolumeSnapshotRequest/Response`, `CreateVolumeBackupRequest/Response`, `ListVolumeBackupsRequest/Response`, `DeleteVolumeBackupRequest/Response`, `RestoreVolumeBackupRequest/Response`, `VolumeBackupInfo` (URL + volume + snapshot + project + size + state + error + created).
- `VolumeInfo.backend` (field 8) — surfaces the storage backend (`file` default, `block` for weft-block). Drives the dashboard's affordance gating on snapshot Revert + Backup actions (block-only).

## [v0.5.0] — 2026-06-02

### Added

- `ListFederationPeers` RPC on `WeftAgent` service + `FederationPeerInfo` / `ListFederationPeersRequest` / `ListFederationPeersResponse` messages — surfaces the in-process `federation.Poller` snapshot (peer name, region, weight, last-seen, classified status). Per [[openweft_pull_model]], the RPC reads the locally-cached pull state ; it does NOT trigger a remote pull on the hot path.
- `ListPluginCatalogue` + `ListInstalledPlugins` + `InstallPlugin` RPCs on `WeftAgent` service + supporting `PluginInput` / `PluginCatalogueEntry` / `PluginInstance` / request-response messages — exposes the `pluginstore.Manager` catalogue + installed-instance registry + idempotent install (returns the deterministic instance UUID).

## [v0.4.0] — 2026-06-02

### Added

- `SetHostCordoned` RPC on `WeftAgent` service + `SetHostCordonedRequest{uuid, cordoned}` / `SetHostCordonedResponse{}` messages — flips the per-host cordon flag (idempotent). Drives `weft host cordon` / `weft host uncordon`.
- `HostInfo.cordoned` (field 14) — surfaces the cordon flag in the registry. Independent of `state` ; a cordoned host stays Active + reachable but the scheduler drops it from candidate sets for new placements.
- `StartVMRequest.requested_gpus` (field 4) and `StartVMRequest.requested_pci` (field 5) — start-time passthrough requests layered on top of the VM's persisted passthrough config. Mirrors the admission-time surface added to `CreateVMRequest` / `RegisterMicroVMRequest` in v0.3.0.

## [v0.3.0] — 2026-06-02

### Added

- `GPURequest` message: `vendor`, `model`, `count`, optional `mig_slice` — mirrors the in-tree `weft/scheduling.GPURequest` struct.
- `PCIPassthroughRequest` message: `vendor_id`, `device_id`, `count` for non-GPU PCI passthrough (NIC, NVMe, FPGA, sound card).
- `CreateVMRequest.requested_gpus` (field 10) and `CreateVMRequest.requested_pci` (field 11) — admission-time passthrough shape, persisted on the VMRecord, enforced by tenant_quotas. Closes the `nil` gap noted in commit 2ca4fce8a.
- `RegisterMicroVMRequest.requested_gpus` (field 9) and `RegisterMicroVMRequest.requested_pci` (field 10) — same surface for the microVM boot path.

## [v0.2.0] — 2026-05-31

### Added

- VolumeSnapshot RPCs: `Create`, `List`, `Restore`, `Delete` (reflink-backed CoW snapshots).

## [v0.1.0] — 2026-05-30

### Added

- Initial proto schema imported from existing tree.
- `Flavors` service: cluster-wide compute catalogue RPCs (etcd-backed).
- `Scripts` service: provisioning-script catalogue RPCs (etcd-backed).
- `VMProperty` service: per-VM host-set annotation RPCs.
- `UEFIVar` service: per-VM firmware NVRAM editor RPCs.
- `VMSSHKey` service: per-VM runtime SSH-key RPCs.
- `CreateVMRequest`: `scheduling_rule` + `network` fields (pull/reconcile labels).

### Removed

- `CreateVMRequest.ssh_pub` (cloud-init era); tag 6 reserved.
