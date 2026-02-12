# Mint and Delegate Tokens

Create a token hierarchy: department → team → individual agent.

```bash
# 1. Mint a department root token ($10,000 budget)
satgate mint --agent "engineering-dept" --budget 10000 --routes "/api/*"
# → Token ID: eng-dept-abc123

# 2. Delegate to a team ($3,000 from engineering's budget)
satgate mint --agent "ml-team" --budget 3000 --parent "eng-dept-abc123" --routes "/api/openai/*,/api/anthropic/*"
# → Token ID: ml-team-def456

# 3. Delegate to an individual agent ($500 from ML team's budget)
satgate mint --agent "training-bot" --budget 500 --parent "ml-team-def456" --routes "/api/openai/v1/chat/*" --expiry 7d

# 4. Check the hierarchy
satgate tokens

# 5. Check spend at any level
satgate spend --agent "ml-team"
```

Key principle: each delegation can only **restrict**, never **escalate**. The training-bot can't access Anthropic (parent restricted to OpenAI + Anthropic, but this token is OpenAI-only). It can't exceed $500. It expires in 7 days even if the parent doesn't.

This is cryptographic delegation — no server round-trip needed to add restrictions.
