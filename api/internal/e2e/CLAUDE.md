# E2E Tests

## Design Philosophy

Tests are organized by **user flow**, not by entity or RPC. Each file represents
a coherent scenario that a user (player or admin) would perform. New RPC coverage
should be added at the natural point in an existing flow where a user would invoke
it, not as a standalone test — unless it genuinely represents a new flow.

## Adding a New Test

1. Identify which flow the new RPC belongs to (or whether it needs a new flow file)
2. Add the RPC call at the natural point in the flow where a user would invoke it
3. Use `seedEvent` / `seedPlayer` helpers for setup — don't duplicate inline SQL
4. Event keys must be unique per test function to avoid DB collisions
