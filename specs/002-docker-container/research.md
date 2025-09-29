# Research Findings: 002-docker-container

## Dockerizing Go CLI Applications

**Decision**: Use multi-stage Docker build with Go base image for building and Alpine for runtime to minimize image size.

**Rationale**: Go CLIs are statically linked, allowing small runtime images. Multi-stage builds separate build dependencies from runtime.

**Alternatives Considered**:
- Single-stage with Ubuntu: Larger image (~500MB vs ~20MB).
- Scratch base: No shell for debugging, harder troubleshooting.

**Best Practices**:
- Use official Go images for build stage.
- Copy only binary to runtime stage.
- Expose no ports (CLI tool).
- Set USER for security.

## CI Automation with GitHub Actions

**Decision**: Use GitHub Actions workflow triggered on push to main, building and pushing to GitHub Container Registry.

**Rationale**: Native integration with GitHub, supports container registry, follows project example workflows.

**Alternatives Considered**:
- GitLab CI: Not applicable.
- Jenkins: Overkill for simple push.
- Manual pushes: Error-prone, not automated.

**Best Practices**:
- Use actions/checkout and actions/setup-go.
- Build with Docker Buildx for multi-platform.
- Login to registry with secrets.
- Tag with 'latest' on main pushes.