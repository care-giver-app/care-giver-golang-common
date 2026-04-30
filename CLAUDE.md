# care-giver-golang-common

Shared Go library imported by all three Go services in the CareToSher monorepo. Contains domain models, DynamoDB repository layer, and event type configs. There is no `main` package — this is not deployed directly.

## Module

```
github.com/care-giver-app/care-giver-golang-common
```

## Consumers

All three services import this module. When you change and tag a new version, update `go.mod` in each:

| Service | Current version ref |
|---|---|
| `care-giver-api` | v0.4.7 |
| `care-giver-notification-executor` | v0.3.0 |
| `care-giver-notification-orchestrator` | v0.0.1 |

**Never edit consumer `go.mod` files from within this repo.** Do it from each consumer's directory after tagging.

## Commands

```bash
make test          # run all tests with coverage
make test-report   # open coverage HTML in browser
make lint          # run golangci-lint
```

Dependencies are vendored — run `go mod vendor` after adding any new dependency, and commit the `vendor/` changes.

## Package Structure

```
pkg/
  event/          - Event domain model, factory, config reader
    types/*.json  - One JSON file per event type (embedded at compile time)
  repository/     - DynamoDB CRUD for each entity; each has a *Provider interface
  user/           - User domain model + factory
  receiver/       - Receiver domain model + factory
  relationship/   - Relationship model + authorization helpers
  dynamo/         - DynamoDB client factory + Mock for tests
  awsconfig/      - AWS SDK config loader
  log/            - Zap logger factory + log key constants
```

## Adding a New Event Type

1. Create `pkg/event/types/<snake_case_name>.json` — the filename **must** be the event `type` string lowercased with spaces replaced by underscores (e.g. `"Doctor Appointment"` → `doctor_appointment.json`).
2. Add a test case to `TestGetAllConfigs` in `pkg/event/event_test.go` asserting the new type appears.
3. Run `make test` — the embed happens at compile time so a missing or misnamed file will cause a test failure.

### Minimal event type JSON

```json
{
    "type": "My Event",
    "icon": "assets/my-event-icon.svg",
    "color": { "primary": "#RRGGBB", "secondary": "#RRGGBB" },
    "isTrackable": false
}
```

### Optional fields in the config schema

| Field | Purpose |
|---|---|
| `"data": { "name": "...", "unit": "..." }` | Single numeric data point (used by Walk, Weight) |
| `"graph": { "type": "line\|scatter", "title": "..." }` | Enables chart rendering on the stats page |
| `"fields": [ ... ]` | Multiple typed form fields (used by Doctor Appointment) |

`fields` entries support `inputType`: `"text"`, `"textarea"`, `"number"`, `"date"`.

## Adding a New Repository

1. Add a `*Provider` interface in the same file as the implementation (see `EventRepositoryProvider` as a model).
2. Use `pkg/dynamo/Mock` for tests — do not hit a real DynamoDB table in unit tests.
3. Register the new log key constants in `pkg/log/log.go` and use them instead of raw strings.

## Testing Conventions

- Table-driven tests using `testify/assert`.
- All repository tests use `pkg/dynamo.Mock` — never a real AWS client.
- Test files live alongside the source file they test (`event_test.go` next to `event.go`).
- Run `make test` before tagging any release.

## Releasing

Merging a PR to `main` triggers the release workflow automatically. No manual tagging needed.

The new version is determined by the **PR title prefix**:

| PR title prefix | Version bump | Example |
|---|---|---|
| `Major: ` | `vX+1.0.0` | Breaking change to an interface |
| `Minor: ` | `vX.Y+1.0` | New package or exported function |
| `Patch: ` | `vX.Y.Z+1` | Bug fix, config change, doc update |

The PR validation workflow enforces this format — a PR title that doesn't start with `Major: `, `Minor: `, or `Patch: ` will fail CI and cannot be merged.

After the merge, the release workflow tags the commit, creates a GitHub Release, and the new version is immediately available to consumers.

**To pick up the new version in a consumer repo:**

```bash
go get github.com/care-giver-app/care-giver-golang-common@v0.x.y
go mod vendor
```

Do this in each of the three consumer repos (`care-giver-api`, `care-giver-notification-executor`, `care-giver-notification-orchestrator`).
