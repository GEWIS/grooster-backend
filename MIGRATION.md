# Monorepo migration

In October 2026, [GEWIS/grooster-frontend](https://github.com/GEWIS/grooster-frontend) was merged into this repository (formerly backend-only).

## Layout

| Path                        | Origin                                                                        |
|-----------------------------|-------------------------------------------------------------------------------|
| `backend/`                  | Backend root of this repository (moved in a regular commit, history intact)   |
| `packages/grooster-client/` | `client/` of this repository (moved in a regular commit, history intact)      |
| `frontend/`                 | GEWIS/grooster-frontend `develop` @ `002501c`, merged with full history       |

The backend history was not rewritten: every commit up to and including `33a4c96` keeps its original hash.
Use `git log --follow <path>` to see the history of files from before the moves.

The frontend history was rewritten with `git filter-repo --to-subdirectory-filter frontend` before merging, so
frontend commits have new hashes. Their content, authors, dates and messages are unchanged.

## Tags

The backend tags (`v1` … `v1.21.1`) are unchanged.
The frontend tags were prefixed to avoid collisions: `v1.23.0` in the frontend repo is now `frontend/v1.23.0` here.
The GitHub releases of the frontend remain in the archived GEWIS/grooster-frontend repository.

From the merge on, backend, frontend and API client share one version, starting at `v2.0.0`.

## Branches

Only frontend `develop` (identical to `main`) was merged. All other frontend branches had already been merged and were not
carried over.

## Issues

Frontend issues were transferred to this repository and labelled `For: Frontend`. Issue numbers changed. References
such as `#62` in frontend commit messages and in old issue or PR texts refer to the **old** frontend numbers below.

| Old (grooster-frontend) | New (grooster-backend) | Title |
|-------------------------|------------------------|-------|
| _to be filled in after the transfer_ | | |
