# GitHub Actions Guide for VPS Deployment

This guide explains how to set up GitHub Actions for automatically building and deploying your Safecast New Map Go application to a VPS (Virtual Private Server). It includes a sample workflow that triggers on pushes to the `main` branch, builds the app, and deploys it via SSH.

## Prerequisites

1. **GitHub Repository**: Your project (`safecast-new-map`) must be in a GitHub repository. Push your code there if it's not already.

2. **VPS Access**:
   - SSH access to your VPS.
   - Generate an SSH key pair if needed: `ssh-keygen -t rsa -b 4096`.
   - Add the public key to your VPS's `~/.ssh/authorized_keys`.

3. **GitHub Secrets**:
   - In your GitHub repo: Settings → Secrets and variables → Actions → Add secrets:
     - `SSH_PRIVATE_KEY`: Paste your entire private SSH key (including `-----BEGIN OPENSSH PRIVATE KEY-----`).
     - `VPS_HOST`: VPS IP address or hostname (e.g., `192.168.1.100`).
     - `VPS_USER`: SSH username (e.g., `ubuntu`).
     - `VPS_PORT`: SSH port (default: `22`).
     - Optional: `VPS_PASSWORD` if using password auth (less secure than keys).

4. **VPS Setup**:
   - Ensure Go is installed on the VPS (for any build steps if needed).
   - Your app directory on VPS (e.g., `/home/ubuntu/safecast-new-map/`).
   - Ability to stop/start the app (e.g., via systemd: `sudo systemctl stop/start safecast-new-map`).

5. **Permissions**: The SSH user must have write access to the app directory and permissions to restart services.

## Step-by-Step Setup

### 1. Create the Workflow File

Create `.github/workflows/deploy.yml` in your repository root with the following content:

```yaml
name: Deploy to VPS

on:
  push:
    branches:
      - main  # Adjust to your default branch
  workflow_dispatch:  # Manual trigger option

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'  # Match your go.mod; check with `go version`

      - name: Build application
        run: |
          go mod tidy
          GOOS=linux GOARCH=amd64 go build -o safecast-new-map .  # Linux x64 binary; adjust GOARCH for your VPS

      - name: Deploy to VPS
        uses: appleboy/ssh-action@v1.0.3
        with:
          host: ${{ secrets.VPS_HOST }}
          username: ${{ secrets.VPS_USER }}
          key: ${{ secrets.SSH_PRIVATE_KEY }}
          port: ${{ secrets.VPS_PORT }}
          script: |
            # Stop current app (adjust for your setup)
            sudo systemctl stop safecast-new-map || pkill -f safecast-new-map || true
            
            # Backup old binary
            mv /path/to/your/app/safecast-new-map /path/to/your/app/safecast-new-map.old || true
            
            # Copy new binary
            scp -o StrictHostKeyChecking=no safecast-new-map ${{ secrets.VPS_USER }}@${{ secrets.VPS_HOST }}:/path/to/your/app/
            
            # Make executable and restart
            chmod +x /path/to/your/app/safecast-new-map
            sudo systemctl start safecast-new-map || /path/to/your/app/safecast-new-map &  # Adjust startup command
```

- **Key Notes**:
  - Replace `/path/to/your/app/` with your actual VPS app directory.
  - The build step creates a Linux binary. Test locally: `GOOS=linux GOARCH=amd64 go build -o safecast-new-map .`
  - Deployment uses SSH/SCP. If your VPS uses password auth, replace `key` with `password: ${{ secrets.VPS_PASSWORD }}`.

### 2. Commit and Push

- Add the file: `git add .github/workflows/deploy.yml`
- Commit: `git commit -m "Add GitHub Actions workflow for VPS deployment"`
- Push: `git push origin main`

The workflow will run on the next push.

### 3. Monitor and Test

- View runs: GitHub repo → Actions tab.
- Check logs for issues (e.g., SSH failures).
- Test manually: Trigger via "workflow_dispatch" or push a small change.
- Verify on VPS: Check if the app restarted and logs.

### Troubleshooting

- **SSH Issues**: Test connection manually: `ssh -i private-key user@vps-host`.
- **Build Failures**: Ensure Go version matches and dependencies are correct.
- **Permissions**: Confirm SSH user can execute commands on VPS.
- **Logs**: Use `journalctl -u safecast-new-map` on VPS for app logs.

### Advanced Customizations

- **Add Tests**: Insert `go test` before deployment.
- **Environments**: Add staging/prod environments for safer deploys.
- **Notifications**: Use actions like `8398a7/action-slack` for alerts.
- **Docker Deployment**: If containerized, use `appleboy/docker-action` instead of SSH.

Study this file, test in a staging environment, and adjust as needed. If you encounter issues, share error logs for help!