# TODO

## Backlog

- [ ] DB Connection Stability (Timeout / Reconnect / Pooling)
- [ ] Dev Tools
  - [ ] DB Seeds
  - [ ] Staging CLI
- [ ] GetScoreForVenue() RPC
- [ ] Bugs
  - [ ] Pending calculation is naive (n-1, even if not latest bar)
- [ ] Telemetry Improvements
  - [ ] Deploy Markers
  - [ ] Schedule Cache Info
- [ ] Static Content RPCs
- [ ] Perf
  - [ ] Cache Scoreboard
  - [ ] Cache, Async or Bulk Fetch VenueByKey
  - [ ] Convert CreateAdjustmentWithTemplate to bulk insert
- [ ] Test Coverage
  - [ ] VenueKey DBC
  - [ ] GetSchedule RPC / lib / E2E
  - [ ] ULID / UInt <-> DB Serialization
  - [ ] RPC Guards
  - [ ] Look into https://github.com/mfridman/tparse
- [ ] Score Submission Followups
  - [ ] Accept hidden adjustment templates
  - [ ] Idempotency Keys
- [ ] Cleanup
  - [ ] Rename DAO methods EntityVerbModifier

## Long-Term Ideas

- [ ] Push Notifications
  - [ ] RegisterDevice() RPC
  - [ ] Background Task
- [ ] Multi-Event
- [ ] Photo Upload
- [ ] SMS Version
- [ ] Dynamic Marketing Site
- [ ] Admin App
  - [ ] CrawlSpace Rebrand
  - [ ] Launch in App Store
