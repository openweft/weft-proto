# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project aims to adhere to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
