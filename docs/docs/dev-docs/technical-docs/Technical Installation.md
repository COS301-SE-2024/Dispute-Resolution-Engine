# Install Guide


**Requirements:**

- [Docker](https://www.docker.com/) installed
- Loki plugin for Docker (install with `docker plugin install grafana/loki-docker-driver:2.9.2 --alias loki --grant-all-permissions`)

## Building and Running
1. Clone the repository from GitHub with the following command:

```sh
git clone https://github.com/COS301-SE-2024/Dispute-Resolution-Engine
cd Dispute-Resolution-Engine
```

2. Build and run the project using Docker

```sh
docker compose up -d
```

Which would automatically download, install, and run the relevant containers to be able
to use the system.

**WARNING:** The application requires some environment variable files to function properly. These
files are confidential and can be obtained by contacting the development team.
