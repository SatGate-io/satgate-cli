# SatGate + lnget: Full Agent Commerce

Set up a paid API endpoint and test it with lnget.

## Server Side (SatGate CLI)

```bash
# 1. Check gateway is running with L402 (Charge mode)
satgate mode
# Should show: 💲 Charge for your paid routes

# 2. Mint a token for the API operator
satgate mint --agent "api-operator" --budget 0  # Unlimited for the operator

# 3. Monitor incoming payments
satgate spend
satgate tokens
```

## Client Side (lnget)

```bash
# Install lnget
npx -y @lightninglabs/lnget https://your-gateway.example.com/api/premium-data

# lnget automatically:
# 1. Gets 402 Payment Required from SatGate
# 2. Pays the Lightning invoice
# 3. Retries with L402 auth header
# 4. Returns the data

# Set a spending ceiling
lnget --max-cost 500 https://your-gateway.example.com/api/expensive-query
```

## Verify on Server

```bash
# See the payment in spend report
satgate spend
# → Shows the lnget payment attributed to the route

satgate report threats
# → No threats (legitimate L402 payment)
```

**Result**: Agent B paid Agent A's API. No signup, no API key exchange, no invoice processing. Lightning settled it in milliseconds. SatGate recorded the transaction, attributed the cost, and enforced the pricing policy.
