# grooster

A GEWIS rooster (roster/schedule) maker.

## Repository layout

| Path                        | Contents                                                        |
|-----------------------------|-----------------------------------------------------------------|
| `backend/`                  | Go backend, see [backend/README.md](backend/README.md)          |
| `frontend/`                 | Vue frontend, see [frontend/README.md](frontend/README.md)      |
| `packages/grooster-client/` | Generated TypeScript API client (`@gewis/grooster-backend-ts`)  |

The API client is generated from the backend's OpenAPI spec (see
[backend/README.md](backend/README.md#generating-docs-and-client)) and published to npm.
