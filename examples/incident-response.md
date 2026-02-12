# Incident Response: Rogue Agent

An agent is burning through budget or accessing unexpected routes.

```bash
# 1. Check current spend — who's burning money?
satgate spend

# 2. Identify the rogue token
satgate tokens | grep -i "high-spend-bot"

# 3. Get full details
satgate token <token-id>

# 4. Revoke immediately
satgate revoke <token-id>
# ⚠️  This is instant and irreversible. The agent loses all access.

# 5. Check threat report
satgate report threats

# 6. Verify revocation
satgate tokens | grep <token-id>
# Should show ⛔ revoked
```

No org-wide API key rotation. No redeployment. One command, instant kill.

Compare to traditional API keys: if an agent's API key is compromised, you have to rotate the key and redeploy every service that uses it. With SatGate macaroon tokens, you revoke the specific token. Everything else keeps running.
