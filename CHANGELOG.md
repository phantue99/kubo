
# CHANGELOG

## Overview
This repository is a fork of `ipfs/kubo`. The changes primarily focus on integrating **AIOZ Network / W3IPFS** specific features. Key modifications include:
*   **Dedicated Gateway Access:** Middleware to restrict or allow access to specific CIDs based on external API validation.
*   **Pinning Service Configuration:** Custom configuration structs and CLI flags to handle pinning service endpoints, API keys, and uploaders.
*   **Block Encryption:** Support for encrypting blocks with specific keys and prefixes.
*   **Image Optimization:** Integration of `aioz-image-optimizer` with configuration for `ffmpeg` and `exiftool`.
*   **Rate Limiting:** Implementation of IP and CID-based rate limiting for the gateway.
*   **DMCA Compliance:** Middleware to check CIDs against a DMCA blocklist API.
*   **Branding:** Customization of the directory index HTML (WebUI) to reflect AIOZ/W3IPFS branding.
*   **Dependency Forks:** Replacement of core libraries (`go-bitswap`, `go-blockservice`, `boxo`) with custom forks (`github.com/phantue99/...`).

## List of Changed Files
*   `.gitignore`
*   `.vscode/launch.json`
*   `Makefile`
*   `Rules.mk`
*   `assets/dir-index-html/dir-index.html`
*   `assets/dir-index-html/src/dir-index.html`
*   `cmd/ipfs/daemon.go`
*   `cmd/ipfs/init.go`
*   `cmd/ipfs/main.go`
*   `config/config.go`
*   `config/config_pinning_service.go`
*   `config/init.go`
*   `core/commands/block.go`
*   `core/commands/dag/import.go`
*   `core/commands/get.go`
*   `core/corehttp/corehttp.go`
*   `core/corehttp/gateway.go`
*   `docs/examples/kubo-as-a-library/go.mod`
*   `docs/examples/kubo-as-a-library/go.sum`
*   `go.mod`
*   `go.sum`
*   `mk/golang.mk`
*   `repo/fsrepo/migrations/ipfsfetcher/ipfsfetcher.go`
*   `test/dependencies/go.mod`
*   `test/dependencies/go.sum`

---

## Commit History

### Commit `2c68dc3f0a3af6d922c711e01bd2fa5c1bf9cf0e`
**Message:** update dir index html
**Date:** 2023-10-10 11:18:55
**Author:** PHANTUE\AIOZ
**Changes:**
*   **Modified:** `assets/dir-index-html/dir-index.html`
*   **Summary:** Completely overhauled the directory listing HTML template. Updated CSS styling, meta tags (OpenGraph, Twitter), and page structure to match specific branding requirements.

### Commit `095dd99599b1915da5f6808d886ee0cd329aea56`
**Message:** update go mod
**Date:** 2023-10-10 11:19:24
**Author:** PHANTUE\AIOZ
**Changes:**
*   **Modified:** `go.mod`, `go.sum`, `mk/golang.mk`
*   **Summary:**
    *   Updated Go version to 1.18.
    *   Added `replace` directives to use forked versions of `github.com/ipfs/go-bitswap` and `github.com/ipfs/go-blockservice` (pointing to `github.com/phantue99/...`).
    *   Updated `golang.org/x/crypto` and `golang.org/x/net` dependencies.

### Commit `7a4b2afc79805276bd2b69960c50f87fade2cd95`
**Message:** Merge remote-tracking branch 'origin/HEAD' into fork-merge
**Date:** 2023-10-10 15:27:06
**Author:** PHANTUE\AIOZ
**Changes:**
*   **Modified:** `cmd/ipfs/daemon.go`, `cmd/ipfs/init.go`, `core/commands/block.go`, `core/commands/get.go`, `go.mod`, `go.sum`
*   **Summary:**
    *   **Daemon:** Implemented `dedicatedGatewayMiddleware` to restrict gateway access based on CID verification via an external API.
    *   **Init:** Added CLI flags for `dedicated-gateway`, `pinning-service`, `api-key`, and `uploader-endpoint`.
    *   **Core:** Updated block retrieval logic to use `cmdutils.PathOrCidPath`.

### Commit `825124fd5c98554234316e22e87e55ccff38991b`
**Message:** add redis conn config
**Date:** 2023-10-10 18:36:42
**Author:** PHANTUE\AIOZ
**Changes:**
*   **Modified:** `cmd/ipfs/daemon.go`, `cmd/ipfs/init.go`, `config/config_pinning_service.go`, `go.mod`, `go.sum`
*   **Summary:**
    *   Added `RedisConns` (slice of strings) to the `ConfigPinningSerice` struct.
    *   Added `--redis-conn` CLI flag to the `init` command.
    *   Passed Redis configuration to `blockservice.InitBlockService`.
    *   Replaced `github.com/ipfs/boxo` with a fork (`github.com/phantue99/boxo`).

### Commit `322604b9f7a7deb01c15adf622bb6258dc7169e9`
**Message:** Update IPFS daemon and corehttp packages
**Date:** 2024-03-18 15:46:06
**Author:** tue.phan
**Changes:**
*   **Modified:** `cmd/ipfs/daemon.go`, `core/commands/dag/import.go`, `core/corehttp/corehttp.go`, `core/corehttp/gateway.go`, `go.mod`
*   **Summary:**
    *   **Refactor:** Moved `DedicatedGatewayMiddleware` and `getDedicatedGatewayAccess` from `daemon.go` to `core/corehttp/corehttp.go`.
    *   **Gateway:** Updated `corehttp.Serve` and `gateway.NewHandler` to utilize the dedicated gateway middleware and configuration.
    *   **DAG Import:** Commented out the default pinning logic in `dagImport`.

### Commit `2b93702142929b2b70216f8cee362fa87f2e5bef`
**Message:** Update dependencies and add rate limiting for IP and CID
**Date:** 2024-04-25 18:04:56
**Author:** tue.phan
**Changes:**
*   **Modified:** `assets/dir-index-html/dir-index.html`, `core/corehttp/corehttp.go`, `core/corehttp/gateway.go`, `go.mod`, `go.sum`, `repo/fsrepo/migrations/ipfsfetcher/ipfsfetcher.go`
*   **Summary:**
    *   **Rate Limiting:** Implemented in-memory rate limiting in `DedicatedGatewayMiddleware` using `golang.org/x/time/rate`. Limits are applied per IP (100 requests/min) and per CID (15 requests/min) when dedicated gateway mode is disabled.
    *   **Branding:** Updated `dir-index.html` title and meta tags to "AIOZ W3IPFS.STORAGE".
    *   **Dependencies:** Added `github.com/juju/ratelimit` and `golang.org/x/time`.

### Commit `7cba022f4bf80791bef071a9182d8962e2fb050c`
**Message:** chore: Update ConfigPinningSerice to ConfigPinningService
**Date:** 2024-05-04 15:32:35
**Author:** tue.phan
**Changes:**
*   **Modified:** `.vscode/launch.json`, `cmd/ipfs/daemon.go`, `cmd/ipfs/init.go`, `config/config.go`, `config/config_pinning_service.go`, `config/init.go`, `core/corehttp/corehttp.go`, `core/corehttp/gateway.go`, `go.mod`, `go.sum`
*   **Summary:**
    *   **Refactor:** Renamed the struct `ConfigPinningSerice` (typo fix) to `ConfigPinningService` across the codebase.
    *   **Config:** Added `AmqpConnect` field to `ConfigPinningService` and a corresponding `--amqp-connect` CLI flag.
    *   **Go Version:** Updated `go.mod` to Go 1.22.

### Commit `9dbb9017e05921c75b642d1c12d48b6e94e06c03`
**Message:** chore: Ignore debug binary files in ipfs command directory
**Date:** 2025-03-04 18:02:44
**Author:** tue.phan
**Changes:**
*   **Modified:** `.gitignore`
*   **Summary:** Added `cmd/ipfs/__debug_bin*` to `.gitignore` to prevent accidental commits of debug binaries.

### Commit `ae869a0b8c514ad1a010d8de07ef9636109cf945`
**Message:** Merge branch 'fork-merge'
**Date:** 2025-03-04 18:03:47
**Author:** tue.phan
**Changes:**
*   **Summary:** Merge commit.

### Commit `20222d6aa43d81806a322d26037997161ebc8b7b`
**Message:** Update IPFS daemon and corehttp packages
**Date:** 2024-05-17 11:03:55
**Author:** tue.phan
**Changes:**
*   **Modified:** `.vscode/launch.json`, `Makefile`, `Rules.mk`, `cmd/ipfs/daemon.go`, `cmd/ipfs/init.go`, `config/config_pinning_service.go`, `go.mod`, `go.sum`
*   **Summary:**
    *   **Config:** Changed `RedisConns` (slice) to `RedisConn` (string) in `ConfigPinningService`.
    *   **Build:** Added `update-ipfs-production` target to `Makefile` and `Rules.mk` for deploying to a remote server via SSH/SCP.
    *   **Dependencies:** Updated `boxo` fork version.

### Commit `252dd98b517d8e329c2849f2cd8ef29d9d6fe3ed`
**Message:** Update IPFS daemon and corehttp packages
**Date:** 2024-06-25 17:43:50
**Author:** tue.phan
**Changes:**
*   **Modified:** `.vscode/launch.json`, `Makefile`, `Rules.mk`, `go.mod`, `go.sum`
*   **Summary:**
    *   **Build:** Added `update-ipfs-staging` target to `Rules.mk`.
    *   **Dependencies:** Updated `boxo` fork version.
    *   **Cleanup:** Removed accidental binary files (`.car` and `__debug_bin`) from `cmd/ipfs/`.

### Commit `84d238954e9ea7046da5536f7a929413102fca20`
**Message:** update dmca
**Date:** 2024-09-05 16:13:37
**Author:** tue.phan
**Changes:**
*   **Modified:** `core/corehttp/corehttp.go`
*   **Summary:**
    *   **Feature:** Added `checkDmca` function to verify if a CID is blocked via an external API.
    *   **Middleware:** Integrated DMCA checks into `DedicatedGatewayMiddleware`. Returns `410 Gone` with a specific legal message if blocked.

### Commit `37fc4b793dbc7b5919b0387834807bdf3e2ae3e9`
**Message:** Update branding and dependencies for AIOZ Pin
**Date:** 2025-03-04 17:58:14
**Author:** tue.phan
**Changes:**
*   **Modified:** `assets/dir-index-html/dir-index.html`, `assets/dir-index-html/src/dir-index.html`, `core/corehttp/corehttp.go`, `go.mod`, `go.sum`
*   **Summary:**
    *   **Branding:** Updated directory index HTML to use "AIOZ Pin" branding, new logo URL, and updated support email in DMCA error messages.
    *   **Dependencies:** Updated Go toolchain to 1.24.1. Updated `boxo` fork and `golang.org/x/crypto`.

### Commit `53d61db58a86a458f2034cf245b756a94026f685`
**Message:** Feature: options BlockEncryptionKey, EncryptedBlockPrefix
**Date:** 2025-05-05 12:35:41
**Author:** Hung Lam
**Changes:**
*   **Modified:** `cmd/ipfs/daemon.go`, `cmd/ipfs/init.go`, `config/config_pinning_service.go`, `go.mod`, `go.sum`
*   **Summary:**
    *   **Feature:** Added support for block encryption.
    *   **Config:** Added `BlockEncryptionKey` and `EncryptedBlockPrefix` to `ConfigPinningService`.
    *   **CLI:** Added `--encrypt-block-key` and `--encrypted-block-prefix` flags to `init`.
    *   **Dependencies:** Updated `golang.org/x/sys`, `net`, `term`, `text`, `tools`.

### Commit `11843aa3af545a41f27bd3e3df48a38658a14d75`
**Message:** Feature: options BlockEncryptionKey, EncryptedBlockPrefix
**Date:** 2025-05-08 08:55:42
**Author:** Hung Lam
**Changes:**
*   **Modified:** `cmd/ipfs/init.go`, `docs/examples/kubo-as-a-library/go.mod`, `docs/examples/kubo-as-a-library/go.sum`
*   **Summary:**
    *   **Fix:** Correctly parsed and passed `encrypt-block-key` and `encrypted-block-prefix` options in `cmd/ipfs/init.go`.
    *   **Dependencies:** Updated dependencies in examples.

### Commit `55968f94561652992d59a1f19407b5dda66f5cb8`
**Message:** Merge pull request #1 from phantue99/feat/encrypt-block
**Date:** 2025-05-21 15:59:44
**Author:** Phan Dinh Tue
**Changes:**
*   **Summary:** Merge commit for encryption features.

### Commit `e05c7b41414827075de527a2afb2a808b40a66c6`
**Message:** chore: fix go mod
**Date:** 2025-05-21 18:09:12
**Author:** Hung Lam
**Changes:**
*   **Modified:** `docs/examples/kubo-as-a-library/go.mod`, `docs/examples/kubo-as-a-library/go.sum`, `go.mod`, `go.sum`
*   **Summary:**
    *   **Dependencies:** Added `github.com/lamgiahungaioz/aioz-image-optimizer`.
    *   **Update:** Updated `boxo` fork version.

### Commit `9363542f8fa2aa05fbb9d900c31152d5c2cc1afc`
**Message:** update boxo version
**Date:** 2025-05-08 14:28:19
**Author:** Hung Lam
**Changes:**
*   **Modified:** `go.mod`, `go.sum`
*   **Summary:** Updated `boxo` fork dependency.

### Commit `a3b7f66188c3390925e1d9f70d9260e205d498cc`
**Message:** update boxo version fix cdn not decrypt
**Date:** 2025-05-08 17:39:32
**Author:** Hung Lam
**Changes:**
*   **Modified:** `.gitignore`, `go.mod`, `go.sum`
*   **Summary:** Updated `boxo` fork dependency to address a CDN decryption issue. Added `.idea` to gitignore.

### Commit `0a28e6368f1ce51fb19ae5bfcbb2975e813624c1`
**Message:** update: update boxo version
**Date:** 2025-05-21 15:01:07
**Author:** Hung Lam
**Changes:**
*   **Modified:** `go.mod`, `go.sum`
*   **Summary:** Updated `boxo` fork dependency and set Go version to 1.23.4.

### Commit `e8871f1e5c48670b40b8e333a701f556996dbc0b`
**Message:** Merge remote-tracking branch 'origin/feat/encrypt-block'
**Date:** 2025-05-21 18:09:53
**Author:** phantue
**Changes:**
*   **Summary:** Merge commit.

### Commit `6c4a61f520a23d435c83a26713a8f7a7002a738a`
**Message:** fix(Makefile): add build command, fix update ipfs for prod and staging
**Date:** 2025-10-21 10:29:45
**Author:** Hung Lam
**Changes:**
*   **Modified:** `.gitignore`, `Makefile`, `go.mod`, `go.sum`
*   **Summary:**
    *   **Makefile:** Added a `build` target. Refactored `update-ipfs-production` and `update-ipfs-staging` to use variables for IPs and directories.
    *   **Dependencies:** Added `replace` directive for `aioz-image-optimizer` to point to a local/private path. Added `webp`, `gift`, and `imaging` libraries.
    *   **Gitignore:** Added `bin/ipfs`.

### Commit `e9520662904541dc0561eab87f1c844076eb0612`
**Message:** fix(img-optimization): set bin path for exiftool and ffmpeg
**Date:** 2025-10-21 11:26:12
**Author:** Hung Lam
**Changes:**
*   **Modified:** `cmd/ipfs/main.go`, `go.mod`, `go.sum`
*   **Summary:**
    *   **Main:** Explicitly set binary paths for `exiftool` and `ffmpeg` using `aiozimageoptimizer` in `main.go`.
    *   **Dependencies:** Updated `boxo` fork.

### Commit `8b92b41c5ce7805443694517b6b416d9b607196c`
**Message:** Merge pull request #4 from phantue99/fix/img-optimization-query-params
**Date:** 2025-10-21 13:42:10
**Author:** Phan Dinh Tue
**Changes:**
*   **Summary:** Merge commit.