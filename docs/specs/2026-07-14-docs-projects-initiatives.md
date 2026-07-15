# Spec: Documents, Project Writes, and Initiatives

- **Date:** 2026-07-14
- **Status:** Accepted
- **Author:** Mike Beaumier
- **Plan:** `docs/plans/2026-07-14-docs-projects-initiatives.md`

## Problem

`linctl` can read projects but cannot manage Linear **documents**, cannot
**write** projects (create/update/status-updates), and has only a read-only stub
for **initiatives**. The project-management skill suite (`create-project`,
`project-checkpoint`, `project-status`, `project-sitrep`, `project-conventions`,
`create-initiative`) uses `linctl`/GraphQL as its fallback path when the Linear
MCP server is unavailable, and that path currently cannot fulfill the suite's
core operations:

- read/write a project's `Project State` document,
- read/write the workspace `Project Conventions` document,
- create and update projects and post project status updates,
- create, read, and update initiatives and attach projects to them.

## Goals

- Add first-class document management: `list`, `get`, `create`, `update`,
  `delete`, scoped to a project or initiative or the workspace.
- Add project write operations: `create`, `update`, and `status-update create`.
- Add initiative management: `list`, `get`, `create`, `update`, and
  `add-project` / `remove-project`.
- Add a raw `graphql` passthrough that reuses `linctl` auth — the documented
  escape hatch and an introspection tool.
- Keep the command surface consistent with the existing agent-first CLI
  (`list/get/create/update/delete`, `--json`/`--plaintext`/table output, filter
  flags, `pkg/api` client methods with GraphQL in `pkg/api/queries.go`).

## Non-goals

- Rewiring the project-management skills to call these commands. That is a
  **follow-up** (see Follow-ups) tracked as a separate Linear project/issues.
- OAuth flows, webhook changes, or issue-side changes.
- Rich TUI / interactive editing. Body content is supplied non-interactively.

## Command surface

Consistent with existing verbs and the three output modes (`--json`,
`--plaintext`, table). Body/markdown content for docs, project descriptions,
and status updates is resolved from (in precedence order) `--body`,
`--body-file <path>`, then piped **stdin**.

### `linctl graphql`
```
linctl graphql -q '<query>' [--query-file f] [--var k=v ...] [--vars-file j.json]
```
Reads a query/mutation from `-q`, `--query-file`, or stdin; reuses
`auth.GetAuthHeader()`; prints the raw JSON `data`/`errors` payload.

### `linctl document` (alias `doc`)
- `document list [--project ID] [--initiative ID] [--query TEXT] [--limit N] [--sort updatedAt|createdAt]`
  — output exposes the document's **project association** so callers can
  distinguish workspace-resident from project-resident docs.
- `document get ID`
- `document create --title T (--body|--body-file|stdin) [--project ID] [--initiative ID] [--icon I] [--color C]`
- `document update ID [--title T] [--body|--body-file|stdin] [--icon I] [--color C]`
- `document delete ID`

Exact-title matching and `updatedAt` "newest wins" selection are left to the
**caller** (generic surface), matching how `project list` exposes raw filters
rather than baking in one consumer's contract.

### `linctl project` (additions)
- `project create --name N [--description D] [--content|--body-file|stdin] [--team KEY] [--lead EMAIL] [--start-date D] [--target-date D] [--state S]`
- `project update PROJECT-ID [same flags as create]`
- `project status-update create PROJECT-ID --health <onTrack|atRisk|offTrack> (--body|--body-file|stdin)`

### `linctl initiative` (alias omitted to avoid `git init`-style confusion)
- `initiative list [--query TEXT] [--limit N] [--sort updatedAt|createdAt]`
- `initiative get ID`
- `initiative create --name N [--description D] [--content|--body-file|stdin] [--owner EMAIL] [--target-date D] [--status S]`
- `initiative update ID [same flags]`
- `initiative add-project ID --project PROJECT-ID`
- `initiative remove-project ID --project PROJECT-ID`

## API surface (Linear GraphQL)

Confirmed via introspection using the new `graphql` command before typed
methods are written. Expected operations:

- Documents: `documents(filter, first, after, orderBy)`, `document(id)`,
  `documentCreate(input: DocumentCreateInput!)`,
  `documentUpdate(id, input: DocumentUpdateInput!)`,
  `documentDelete(id)`.
- Projects: `projectCreate(input: ProjectCreateInput!)`,
  `projectUpdate(id, input: ProjectUpdateInput!)`,
  `projectUpdateCreate(input: ProjectUpdateCreateInput!)` (status updates).
- Initiatives: `initiatives(...)`, `initiative(id)`,
  `initiativeCreate(input)`, `initiativeUpdate(id, input)`,
  `initiativeToProjectCreate(input)` / `initiativeToProjectDelete(id)`.

Struct changes in `pkg/api/queries.go`:
- `Document`: add `Project` (id/title) and `Initiative` (id) associations.
- `Initiative`: expand beyond id/name/description to owner, status, dates, url,
  timestamps; add `Initiatives` paginated collection.

## Testing

- `smoke_test.sh` gains read-only checks: `document list`, `initiative list`.
- New `write_smoke_test.sh`, opt-in via `LINCTL_WRITE_TESTS=1`, creates → gets →
  updates → deletes a throwaway document, project, and initiative, asserting
  self-cleanup so the real workspace is not polluted.

## Follow-ups

- **Create a Linear project + issues** to migrate the project-management skill
  suite (`project-checkpoint`, `project-conventions`, `project-status`,
  `project-sitrep`, `project-context`, `create-project`, `create-initiative`)
  onto these new `linctl` commands as the MCP-fallback path, replacing the
  hand-rolled GraphQL references in their skill docs.
