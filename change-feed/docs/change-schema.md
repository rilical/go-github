# Change feed schema

The emitter writes three artifacts in the output directory:

- `changes.json`: pretty JSON snapshot of one `changes.ChangeBatch`.
- `feed.jsonl`: append-only JSON Lines stream; one change record per line.
- `report.md`: deterministic markdown summary for humans.

Each `feed.jsonl` record includes these fields:

- `id`
- `kind`
- `severity`
- `method`
- `path`
- `text`

`changes.json` also includes a summary object with counts:

- `Total`, `Breaking`, `Warn`, `Info`
- `Added`, `Removed`, `Modified`
