
# Learn_Jenkins

This repository contains a simple Go-based backend used to demonstrate a Jenkins CI/CD pipeline, Docker deployment, and basic REST endpoints for user CRUD. The project is intentionally small and focuses on CI/CD flows and deployment.

## Table of contents
- **Project overview**
- **Quick start (development)**
- **Run tests**
- **Build & Docker**
- **CI/CD (Jenkins)**
- **Repository structure**
- **API / Postman**
- **Files referenced**


## Project overview

The application is a minimal Go HTTP server built with `gin`. It exposes endpoints to create and fetch users and is packaged into Docker images for three deployment targets: development, staging and production. The repository includes a `Jenkinsfile` that builds, tests, pushes Docker images, updates `docker-compose.*.yaml` files and deploys to a remote host via SSH.

## Quick start (development)

Requirements:
- Go 1.24+
- Docker & docker-compose (for local container runs)
- Make

1. Copy environment variables from `.env.example` to `.env` and set values.
2. Run tests (recommended): `go test ./...`.
3. Run locally: `go run main.go` (the app listens on `:$PORT`).
4. Or build binary: `make build` then `./Learn_Jenkins`.

## Run tests

Unit tests are included under `controllers`, `repositories`, and `services` directories. To run all tests:

```bash
go test ./...
```

## Build & Docker

Build Image locally:

```bash
make build
```


Use docker-compose files to run the app in different environments:

- `docker-compose.development.yaml` (maps host port 8001 -> container 8001)
- `docker-compose.staging.yaml` (maps host port 8002 -> container 8001)
- `docker-compose.production.yaml` (maps host port 8003 -> container 8001)

Start development stack:

```bash
docker-compose -f docker-compose.development.yaml up -d --build
```

## CI/CD (Jenkins)

This repository contains a `Jenkinsfile` implementing the following pipeline:

1. Checkout code and check commit/PR messages for skip tokens.
2. Determine target environment and image name based on branch (`develop` -> development, `staging` -> staging, `master/main` -> production).
3. Run unit tests inside a Docker container (Go image).
4. Build and push Docker image (for target branches).
5. Update `docker-compose.${TARGET}.yaml` with the new image tag.
6. Commit the updated compose file and push to GitHub (skipping CI for commit).
7. SSH to remote host and run `docker-compose -f docker-compose.${TARGET}.yaml up -d --remove-orphans --force-recreate`.

### Jenkins pipeline notes

- Credentials used in the pipeline (configured in Jenkins): `docker-credential`, `github-credential`, `ssh-key`, plus secrets for DB and host values.
- The pipeline builds images with a tag combining git short SHA and build number.

### CI/CD screenshot

Below is the current CI/CD screenshot captured from your environment showing three running containers for each target (production, development, staging):

```startLine:50:endLine:211:img/Screenshot 2025-10-03 011611.png
// Screenshot: CI/CD containers (production, development, staging)
```

> Note: The screenshot is attached to the repository under `img/`.

## Repository structure

High-level folders and important files:

- `config/` - database initialization and utilities (`database.go`, `database_util.go`).
- `controllers/` - HTTP controllers and tests.
- `domain/` - DTOs and models.
- `middlewares/` - HTTP middlewares used by Gin.
- `repositories/` - database access layer and tests.
- `routes/` - route wiring.
- `services/` - business logic and tests.
- `Dockerfile` - multi-stage Dockerfile to build the Go binary and produce a small runtime image.
- `Jenkinsfile` - pipeline for CI/CD (build, test, push, deploy).
- `docker-compose.*.yaml` - environment-specific compose files used during deployment.
- `postman.json` - sample requests for testing the API.

## Files referenced (quick descriptions)

- `main.go`: application bootstrap, database migration (`AutoMigrate`), route registration and server start.
- `Makefile`: convenience targets for building and running.
- `go.mod` / `go.sum`: module definitions and dependency checksums.

## API / Postman

Load the included `postman.json` collection and set `base_url` to `http://localhost:8001` (or the port mapped by docker-compose). It includes endpoints to create a user and fetch users.

## Notes

- Ensure secrets and credentials are configured securely in Jenkins and not checked into the repo.
- The pipeline uses `sed -i` to update compose files — on macOS use `sed -i ''` or adapt accordingly.

## VM setup

This project expects two VMs (or machines): an **app-server** (where the application runs) and a **jenkins-server** (CI/CD). The instructions below summarize the steps performed during provisioning and the minimum configuration required.

- **Firewall (recommended)**: allow `tcp` ports `22` (SSH), `8000`, `8001`, `8002`, `8003` (app targets), and `8080` (Jenkins UI).
- **Cloud SQL**: this repo assumes PostgreSQL (Cloud SQL or managed Postgres) — conceptually Postgres 17 was used during development.

### Jenkins server (basic install)

1. Update packages and install Java:

```bash
sudo apt update
sudo apt install fontconfig openjdk-17-jre -y
java -version
```

2. Add Jenkins apt repository and install:

```bash
sudo wget -O /usr/share/keyrings/jenkins-keyring.asc \
  https://pkg.jenkins.io/debian-stable/jenkins.io-2023.key
echo "deb [signed-by=/usr/share/keyrings/jenkins-keyring.asc] https://pkg.jenkins.io/debian-stable binary/" | sudo tee /etc/apt/sources.list.d/jenkins.list > /dev/null
sudo apt-get update
sudo apt-get install jenkins -y
```

3. Open Jenkins web UI (port `8080`) and get the initial admin password:

```bash
sudo cat /var/lib/jenkins/secrets/initialAdminPassword
```

4. Recommended Jenkins plugins: **Blue Ocean**, **Docker**, **Go**, **SSH Agent**.

5. Configure Jenkins credentials (IDs used in `Jenkinsfile`):

- `docker-credential` (username/password)
- `github-credential` (username/password or PAT)
- `ssh-key` (SSH Username with private key)
- Additional secret text or username/password credentials for DB and host values as used in the `Jenkinsfile`.

### App host VM (deploy target)

1. Prepare the OS and install Docker & Docker Compose:

```bash
sudo apt update
sudo apt install apt-transport-https ca-certificates curl software-properties-common -y
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt update
sudo apt install docker-ce -y
sudo systemctl enable --now docker
sudo usermod -aG docker ${USER}
newgrp docker || true
sudo apt install docker-compose -y
```

2. Add the Jenkins public SSH key to the app VM `~/.ssh/authorized_keys` so the Jenkins server can SSH and deploy. In the `Jenkinsfile`, the SSH credential id `ssh-key` is used.

### GitHub webhook and Jenkins

- Configure a GitHub webhook on the repository to point to `http://<jenkins-host>:8080/github-webhook/` (or use Jenkins GitHub plugin configuration) so pushes and PRs trigger builds.
- In multibranch or pipeline jobs, add the repository and credentials. Use the `Jenkinsfile` from the repo as the pipeline script.

### Notes and troubleshooting

- Make sure the Jenkins user has the SSH private key (credential `ssh-key`) and that the corresponding public key is in the app VM `authorized_keys`.
- Ensure the database (Cloud SQL) allows connections from the app VM, or use a private connection/VPC as required.
- When testing locally, `docker-compose.development.yaml` maps host port `8001` to container port `8001` — adjust `postman.json` base_url accordingly.

---

