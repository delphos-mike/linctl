# Plan: Documents, Project Writes, and Initiatives

- **Date:** 2026-07-14
- **Plan ID:** `2026-07-14-docs-projects-initiatives`
- **Spec:** `docs/specs/2026-07-14-docs-projects-initiatives.md`
- **Branch:** `mike-linctl-docs-projects-initiatives` (off `main`)
- **Notes:** `.local/notes/2026-07-14-docs-projects-initiatives/` (untracked)

Each numbered step is one commit. Verify before committing:
`make fmt && make build && go vet ./...`, then `./smoke_test.sh` (read-only).
Commit and push after each green capability.

## 1. Conventions

- Create `docs/specs/` and `docs/plans/` (with `.gitkeep`).
- Add `.local/` to `.gitignore`.
- Add a "Specs & Plans" section to `AGENTS.md`.
- Commit the spec and this plan.

Note: local pre-existing state converts tracked `CLAUDE.md` into a symlink to
untracked `AGENTS.md`. Stage `AGENTS.md` + the `CLAUDE.md` symlink together in
this commit (AGENTS.md becomes canonical). Leave `.sisyphus/` untracked.

## 2. `linctl graphql` raw passthrough

- `pkg/api/client.go`: add `ExecuteRaw(ctx, query, vars) (json.RawMessage, error)`
  if the existing `Execute` cannot return an untyped payload.
- `cmd/graphql.go`: query from `-q` / `--query-file` / stdin; `--var k=v`
  (repeatable) and `--vars-file` for variables; reuse `auth.GetAuthHeader()` +
  `api.NewClient`; print raw JSON via `output.JSON`.
- Register in `init()`.

## 3. Introspect Linear API (notes only, no commit)

Using the `graphql` command, confirm and record in
`.local/notes/2026-07-14-docs-projects-initiatives/`:
- `DocumentFilter` shape (project/initiative filter, full-text search field).
- `DocumentCreateInput` / `DocumentUpdateInput` fields.
- `ProjectCreateInput` / `ProjectUpdateInput` / `ProjectUpdateCreateInput`
  (health enum values).
- `InitiativeCreateInput` / `InitiativeUpdateInput`, initiative status enum,
  and `initiativeToProjectCreate` / `...Delete` shapes.

## 4. Documents CRUD

- `pkg/api/queries.go`: extend `Document` with `Project`/`Initiative`; add
  `Documents` pagination if not reused; add client methods `GetDocuments`,
  `GetDocument`, `CreateDocument`, `UpdateDocument`, `DeleteDocument`.
- `cmd/document.go`: `list/get/create/update/delete` with the flags in the spec.
- `cmd/root.go`-level helper `resolveBody(cmd)` (or in a shared `cmd` file)
  supporting `--body`, `--body-file`, and stdin.
- Register command + alias `doc`.

## 5. Project writes

- `pkg/api/queries.go`: `CreateProject`, `UpdateProject`, `CreateProjectUpdate`.
- `cmd/project.go`: add `create`, `update`, and a `status-update` subgroup with
  `create`. Resolve team/lead by key/email to IDs (reuse existing helpers).

## 6. Initiatives CRUD

- `pkg/api/queries.go`: expand `Initiative`; add `Initiatives`, `GetInitiatives`,
  `GetInitiative`, `CreateInitiative`, `UpdateInitiative`,
  `AddProjectToInitiative`, `RemoveProjectFromInitiative`.
- `cmd/initiative.go`: `list/get/create/update/add-project/remove-project`.
- Register command.

## 7. Tests + docs + follow-up

- Extend `smoke_test.sh`: `document list`, `initiative list` (read-only).
- Add `write_smoke_test.sh` gated by `LINCTL_WRITE_TESTS=1`: create→get→update→
  delete throwaway doc/project/initiative, asserting cleanup.
- Update `README.md`, `AGENTS.md` feature list + structure, and
  `master_api_ref.md` (Documents/Initiatives sections).
- Confirm the spec's Follow-ups note (Linear project/issues to migrate the PM
  skills) is present.

## Verification summary

`make fmt` → `make build` → `go vet ./...` → `./smoke_test.sh` per commit.
Write paths verified once via `LINCTL_WRITE_TESTS=1 ./write_smoke_test.sh`
against the real workspace using self-cleaning throwaway records.
