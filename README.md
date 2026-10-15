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

## Running the project

The project is run with [Task](https://taskfile.dev) (see [installation](https://taskfile.dev/installation/)),
using the tasks defined in [`Taskfile.yml`](Taskfile.yml).

Before the first run, set up the backend's `.env` file (see
[backend/README.md](backend/README.md)):

```bash
cp backend/.env-example backend/.env
```

Then install dependencies, seed the database and start the dev servers:

```bash
task install   # install Go and npm dependencies
task seed      # seed the database
task dev       # run the backend and frontend dev servers
```

| Task                 | Description                                         |
|----------------------|-----------------------------------------------------|
| `task install`       | Install all Go and npm dependencies                 |
| `task seed`          | Seed the database                                   |
| `task dev`           | Run the backend and frontend dev servers together   |
| `task dev:backend`   | Run only the backend dev server                     |
| `task dev:frontend`  | Run only the frontend dev server                    |

Run `task --list` to see all available tasks.
