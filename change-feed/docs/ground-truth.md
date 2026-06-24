# Ground Truth - captured from the real system (anti-slop reference)

> Captured 2026-06-24 by running the validated spike (`files/spike/`) against the real POC spec
> pair (`files/poc/base.json` -> `head.json`), oasdiff v1.20.0, Go 1.26.4. **Every test assertion
> in the implementation plan must match THESE observed values, not imagined ones.** If a future
> oasdiff bump changes these, re-capture and update the plan - do not "fix" the test to a guess.

## Why this file exists

The first draft of the plan invented oasdiff rule IDs (`response-property-removed`,
`api-removed-without-deprecation`) that the real system NEVER emits. The real IDs are different.
Asserting against invented IDs is exactly the vibe-coded slop we are preventing. This file is the
observed reality. Tests cite it.

## Headline counts (the anchor regression target)

```
total changes : 348
breaking (ERR): 36
warn (WARN)   : 70
info (INFO)   : 242
endpoints added (EndpointsDiff.Added)  : 6
endpoints deleted (EndpointsDiff.Deleted): 2
serialized json bytes (sanity): 125078
```

Window: base `1d6267568cc6...` -> head `3e08e45a052d...`.

## First serialized record (exact shape - field names are ground truth)

```json
{
  "id": "response-property-became-nullable",
  "text": "the response property `items/description` became nullable for the status `200`",
  "level": 3,
  "operation": "GET",
  "operationId": "code-security/get-configurations-for-enterprise",
  "path": "/enterprises/{enterprise}/code-security/configurations",
  "section": "paths",
  "fingerprint": "221396095a4c"
}
```

Go struct field names (from oasdiff source `formatters/changes.go`): `Id, Text, Level
(checker.Level), Operation, OperationId, Path, Section, Fingerprint`. `Comment` is omitempty;
`IsBreaking` has `json:"-"`. `diff.Endpoint{Method, Path}` backs EndpointsDiff.

## The COMPLETE rule-ID x severity distribution (all 16 real IDs, sums to 348)

| Level | Count | Rule ID | Maps to ChangeKind (heuristic) |
|------:|------:|---------|--------------------------------|
| L1 INFO | 216 | `response-optional-property-added` | response.changed |
| L2 WARN | 49 | `response-property-enum-value-added` | response.changed |
| L2 WARN | 21 | `response-optional-property-removed` | response.changed |
| L3 ERR | 20 | `response-required-property-removed` | response.changed |
| L3 ERR | 13 | `response-property-became-nullable` | response.changed |
| L1 INFO | 6 | `endpoint-added` | operation.added |
| L1 INFO | 5 | `api-schema-removed` | schema.changed |
| L1 INFO | 4 | `new-optional-request-property` | request.body.changed |
| L1 INFO | 3 | `request-property-list-of-types-widened` | request.body.changed |
| L1 INFO | 2 | `request-property-one-of-added` | request.body.changed |
| L1 INFO | 2 | `response-property-any-of-added` | response.changed |
| L3 ERR | 2 | `api-path-removed-without-deprecation` | operation.removed |
| L1 INFO | 2 | `request-property-became-optional` | request.body.changed |
| L1 INFO | 1 | `response-body-all-of-added` | response.changed |
| L3 ERR | 1 | `response-body-type-changed` | response.changed |
| L1 INFO | 1 | `request-property-enum-value-added` | request.body.changed |

Severity check: INFO 1+1+2+2+2+3+4+5+6+216 = **242**; WARN 21+49 = **70**;
ERR 1+2+13+20 = **36**. Total **348**. Matches headline. ✓

## Design facts this forces (not optional)

1. **Removal rule for ops is `api-path-removed-without-deprecation`** (L3). The id
   `api-removed-without-deprecation` (no `-path-`) is the oasdiff id for a different shape; both
   exist in the ruleset, so the deprecation-severity index must handle BOTH spellings. Never assert
   only the guessed one.
2. **`endpoint-added` count (6) == EndpointsDiff.Added (6)** and
   **`api-path-removed-without-deprecation` count (2) == EndpointsDiff.Deleted (2)**. This is why
   the normalizer dedups checker add/remove rows against the EndpointsDiff sets - otherwise every
   added op is double-counted (once as endpoint-added INFO, once as the EndpointsDiff add).
3. **oasdiff ships ~992 rule keys; this real window exercises only 16.** A hand-written exhaustive
   `kindMap` of literal IDs will be wrong and rot. Use a **prefix/substring heuristic**
   (`endpoint-added`->added; `api-*removed*`/`api-path-*removed*`->operation.removed;
   `api-schema-removed`->schema.changed; `request-*`->request.body/param; `response-*`->response.changed;
   `*deprecat*`->operation.deprecated; security->security.changed) with an explicit `other` fallback.
   Keep a SMALL exact-match override table only where the heuristic is wrong. This is the DRY +
   explicit-over-clever choice and survives oasdiff growth.
4. **`api-schema-removed` is INFO (L1)**, not breaking - a removed reusable schema component is not
   by itself a breaking endpoint change in oasdiff's model. Don't assume "removed == breaking".

## How to re-capture (when oasdiff is bumped)

```bash
cd files/spike && go run .         # headline counts + first record
# rule-id histogram: re-add the throwaway TestDumpRealRuleIDs from session history,
# run `go test -v -run TestDumpRealRuleIDs`, then delete it (it is throwaway exploration).
```
