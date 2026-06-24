# Change feed schema

The emitter writes three artifacts in the output directory:

- `changes.json`: pretty JSON snapshot of one `changes.ChangeBatch`. Validated
  against `internal/output/feed.schema.json` by `output.ValidateFeed`.
- `feed.jsonl`: append-only JSON Lines stream; one change record per line.
- `report.md`: deterministic markdown summary for humans.

The published JSON Schema lives at `internal/output/feed.schema.json` and is
embedded in the binary (exported as `output.FeedSchema`) so downstream consumers
can validate the feed with their own tooling.

## Record fields

Each `feed.jsonl` record and each entry in `changes.json.changes` has:

- `id`: the oasdiff rule/category id (for example `response-optional-property-added`).
  It is NOT unique per record. The same rule fires for many operations, and every
  synthetic added/removed operation shares a class id (`endpoint-added` /
  `operation-removed`). Do not key on `id`.
- `fingerprint`: a stable, unique-per-change key. Every record carries one
  (synthesized from kind+method+path when oasdiff does not supply it). Dedupe the
  feed on `fingerprint`.
- `kind`: the normalized change kind (for example `operation.added`).
- `severity`: one of the enum strings `info`, `warn`, or `breaking`.
- `method`, `path`, `operation_id`, `section`, `text`: change context (optional).

## Summary

`changes.json` includes a `summary` object with snake_case integer counts:

- `total`, `breaking`, `warn`, `info`
- `added`, `removed`, `modified`

## Delivery semantics

`feed.jsonl` is at-least-once: a clean run is deduped per spec window by the
`.feed_window` marker, but a process killed mid-emit can re-append a window on the
next run. Consumers must dedupe on `fingerprint`. `changes.json` and `report.md`
are last-write snapshots of the most recent run, written atomically (temp+rename).
