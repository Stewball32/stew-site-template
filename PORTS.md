# Port registry — host-wide

Single source of truth for localhost port allocation across **all** of Stewart's
sites on a shared host. Every service binds **loopback only** (`127.0.0.1`) and is
fronted by the cloudflared tunnel — these ports are never exposed directly.

> **Adding a site?** Claim the next free block below, record it here in the same
> commit, and pass it to `scripts/scaffold-site.sh --port-base`. Never reuse a
> reserved port.

## Block convention

Each project gets a **contiguous block of 10**, allocated as:

| Offset | Tier | Purpose |
| --- | --- | --- |
| `base + 0` | **prod** | production container (published to loopback) |
| `base + 1` | **test** (`pre`) | beta/preview container — the gate before prod |
| `base + 2` | **dev backend** | Air-reloaded Go/PocketBase (internal; proxied by Vite) |
| `base + 3` | **dev vite** | Vite HMR dev server — the tunnel's dev target |
| `base + 4..9` | _reserved_ | headroom (extra services, second bot, metrics, …) |

Only **prod**, **test**, and **dev vite** are ever routed through cloudflared.
The dev backend is internal — Vite proxies `/api` to it.

## Allocations

### 🔒 Reserved — do NOT reassign (pre-existing, grandfathered)

These predate the block convention and are **interleaved**; they are recorded
exactly as deployed. Leave them alone.

| Project | prod | test | dev backend | dev vite | Other |
| --- | --- | --- | --- | --- | --- |
| **norcal-halo-site** | `8098` | `8101` | `8100` | `5273` | — |
| **xemu-cartographer** | `8099` | `18099` | `19090` | `19099` | container base `3300` |

> ⚠️ `8098`–`8101` interleaves the two projects (halo owns 8098/8100/8101, cart
> owns 8099). Treat the whole `8098`–`8101` range plus `5273`, `18099`, `19090`,
> `19099` and `3300` as taken.

### ✅ Assigned — block convention

| Project | Block | prod | test | dev backend | dev vite | Public hosts |
| --- | --- | --- | --- | --- | --- | --- |
| **align-the-day** | `8110`–`8119` | `8110` | `8111` | `8112` | `8113` | `alignthe.day` (own apex) |
| **NautsLadder** | `8120`–`8129` | `8120` | `8121` | `8122` | `8123` | TBD |

Both previously defaulted to the template's `8090` and would have collided on one
host — these blocks resolve that. Apply during rollout (Phase 2).

### 🧪 Template placeholder

| Project | Block | prod | test | dev backend | dev vite |
| --- | --- | --- | --- | --- | --- |
| **stew-site-template** | `8090`–`8093` | `8090` | `8091` | `8092` | `8093` |

These are **placeholders only** — `scripts/scaffold-site.sh` overwrites them for
every new site. **Never deploy the raw template**; a scaffolded site must claim
its own block here first.

### Free blocks

`8130`–`8139`, `8140`–`8149`, `8150`–`8159`, … allocate upward in order.

Avoid: `3300`, `5273`, `8090`–`8093` (template), `8098`–`8101`, `18099`,
`19090`, `19099`, and anything already listed above.

**Also observed listening on this host (non-project, avoid):** `8750`, `8760`,
`8764`, `8765`. Not owned by these repos — verify before claiming anything nearby.

## Checking a block is actually free

```sh
# every port in a candidate block
for p in $(seq 8130 8139); do
  ss -tlnp 2>/dev/null | grep -q ":$p " && echo "$p BUSY" || echo "$p free"
done
```
