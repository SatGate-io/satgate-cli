# SatGate CLI

**Manage your API's economic firewall from the terminal.**

The server-side counterpart to [lnget](https://github.com/lightninglabs/lightning-agent-tools). They give agents wallets; we give API owners the cash register.

## Install

### Claude Code Plugin
```bash
claude plugin marketplace add satgate-io/satgate-cli
claude plugin install satgate-cli@satgate-io
```

### ClawHub
```bash
clawhub install satgate-io/satgate
```

### Shell Script
```bash
curl -fsSL https://raw.githubusercontent.com/SatGate-io/satgate-cli/main/scripts/install.sh | bash
```

### From Source
```bash
git clone https://github.com/SatGate-io/satgate-cli.git
cd satgate-cli
make build
./bin/satgate version
```

## Quick Start

```bash
# Configure (interactive)
./scripts/configure.sh

# Or set env vars
export SATGATE_GATEWAY=http://localhost:9090
export SATGATE_ADMIN_TOKEN=sgk_your_token

# Check connection
satgate ping
satgate status

# Mint a token for an agent
satgate mint --agent "my-bot" --budget 500 --expiry 30d --routes "/api/openai/*"

# Check spend
satgate spend

# List all tokens
satgate tokens

# Revoke a compromised agent
satgate revoke <token-id>
```

## Commands

| Command | Description |
|---------|-------------|
| `satgate status` | Gateway health, version, uptime |
| `satgate ping` | Liveness check (exit 0 = healthy) |
| `satgate mint` | Mint a new capability token |
| `satgate tokens` | List all tokens with spend/budget |
| `satgate token <id>` | Token detail view |
| `satgate revoke <id>` | Revoke a token (irreversible) |
| `satgate spend` | Spend summary (org-wide or per-agent) |
| `satgate report threats` | Security threat report |
| `satgate mode` | Current policy mode per route |
| `satgate version` | CLI version and build info |

## Safety

- **Target printing**: Every mutating command shows the gateway URL before executing
- **Interactive confirmation**: Destructive ops require `y/N` confirmation
- **`--dry-run`**: Preview what would happen without executing
- **`--yes`**: Skip prompts (for CI/scripting — use with care)

## Dual Surface Support

The CLI works with both self-hosted gateways and SatGate Cloud:

| Surface | Auth | URL |
|---------|------|-----|
| `gateway` | `X-Admin-Token` | `http://localhost:9090` |
| `cloud` | Session cookie | `https://cloud.satgate.io` |

Auto-detected from the gateway URL, or set explicitly:
```bash
export SATGATE_SURFACE=cloud
export SATGATE_SESSION_TOKEN=eyJ...
```

## SatGate + lnget

**lnget** (Lightning Labs) handles the client side — agents paying for L402-gated APIs.
**SatGate CLI** handles the server side — enforcement, attribution, governance.

Together: the complete agent commerce stack.

```
Agent + lnget  ──→  SatGate Gateway  ──→  Your API
   (pays)           (enforces)            (serves)
```

See [docs/guides/lnget-integration.md](https://github.com/SatGate-io/satgate/blob/main/docs/guides/lnget-integration.md) for the full integration guide.

## Build

```bash
make build          # Debug binary → bin/satgate
make release        # Stripped binaries for all platforms → bin/release/
make test           # Run tests
make clean          # Remove build artifacts
```

## License

Apache 2.0 — see [LICENSE](LICENSE)
