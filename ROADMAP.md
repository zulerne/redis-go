# Roadmap

Current state: TCP server, RESP protocol, basic commands (GET/SET with TTL, lists, streams), CI/CD, linting.

## v0.2 — Core Commands

Missing fundamental Redis commands that are expected in any implementation.

- [ ] `DEL <key> [key ...]` — delete one or more keys
- [ ] `EXPIRE <key> <seconds>` / `PEXPIRE <key> <ms>` — set TTL on existing key
- [ ] `TTL <key>` / `PTTL <key>` — get remaining TTL
- [ ] `EXISTS <key> [key ...]` — check if keys exist
- [ ] `KEYS <pattern>` — find keys matching glob pattern
- [ ] `XRANGE <key> <start> <end> [COUNT count]` — read stream entries

## v0.3 — Hash Data Type

Hash is the most commonly used Redis type after strings and lists. Required for the messenger project (user sessions, metadata).

- [ ] `HSET <key> <field> <value> [field value ...]`
- [ ] `HGET <key> <field>`
- [ ] `HGETALL <key>`
- [ ] `HDEL <key> <field> [field ...]`
- [ ] `HEXISTS <key> <field>`
- [ ] `HLEN <key>`

## v0.4 — Set Data Type

Sets are used for unique collections, tagging, relationships.

- [ ] `SADD <key> <member> [member ...]`
- [ ] `SMEMBERS <key>`
- [ ] `SISMEMBER <key> <member>`
- [ ] `SREM <key> <member> [member ...]`
- [ ] `SCARD <key>`
- [ ] `SINTER <key> [key ...]` / `SUNION <key> [key ...]`

## v0.5 — Persistence

RDB-like snapshot: periodic and on-demand save/load of the data store.

- [ ] `SAVE` — synchronous snapshot to disk
- [ ] `BGSAVE` — background snapshot
- [ ] Auto-load on startup
- [ ] Configurable save intervals

## v0.6 — Pub/Sub

Required for the messenger project — message delivery between service instances.

- [ ] `SUBSCRIBE <channel> [channel ...]`
- [ ] `PUBLISH <channel> <message>`
- [ ] `UNSUBSCRIBE [channel ...]`
- [ ] `PSUBSCRIBE <pattern>` — pattern-based subscription

## v0.7 — Docker & Integration Tests

Production-readiness: containerized deployment and end-to-end testing.

- [ ] Multi-stage Dockerfile
- [ ] docker-compose.yml
- [ ] Integration tests (TCP client → server → assert RESP responses)
- [ ] README architecture diagram (Mermaid)
