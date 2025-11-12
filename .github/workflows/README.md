# GitHub Actions Workflows Documentation

This directory contains GitHub Actions workflow configurations for the ME-Chain project.

## Workflow List

### 1. `build-push-private-net.yml`
Automatically builds private network Docker images and pushes them to Harbor registry.

**Trigger Conditions:**
- Push Git tags with `v*` format (e.g., v1.0.0, v2.1.3)

**Features:**
- Build private network images with pre-configured genesis accounts
- Push images to Harbor registry
- Automatically tag with version and latest labels
- Generate build summary

**Runner:**
- Uses self-hosted Linux runner

## Required GitHub Secrets

Before using these workflows, you need to configure the following Secrets in your GitHub repository:

### Harbor Registry Configuration

Navigate to repository settings: `Settings` > `Secrets and variables` > `Actions` > `New repository secret`

#### 1. `HARBOR_REGISTRY`
- **Description:** Harbor registry address
- **Example value:** `harbor.example.com` or `192.168.0.79`
- **Purpose:** Hostname or IP address of the Docker registry

#### 2. `HARBOR_USERNAME`
- **Description:** Harbor login username
- **Example value:** `admin` or `robot$project+deployer`
- **Purpose:** Authentication for docker login

#### 3. `HARBOR_PASSWORD`
- **Description:** Harbor login password or Robot Token
- **Example value:** `your_password` or Robot Token
- **Purpose:** Authentication for docker login
- **Note:** Recommended to use Harbor Robot Account Token instead of personal password

## Configuration Steps

### 1. Create Harbor Robot Account (Recommended)

Create a Robot Account in Harbor:

1. Login to Harbor Web UI
2. Navigate to project `st-chain`
3. Click `Robot Accounts` tab
4. Click `+ NEW ROBOT ACCOUNT`
5. Configure:
   - **Name:** `github-actions-deployer`
   - **Expiration time:** Set as needed
   - **Permissions:** Check `Push Artifact`
6. Save the displayed Token after creation

### 2. Configure Secrets in GitHub

```bash
# Example values (replace with actual values)
HARBOR_REGISTRY: harbor.example.com
HARBOR_USERNAME: robot$st-chain+github-actions-deployer
HARBOR_PASSWORD: eyJhbGc... (Robot Token)
```

Configuration location:
```
https://github.com/YOUR_ORG/YOUR_REPO/settings/secrets/actions
```

### 3. Verify Configuration

Push a test tag:

```bash
# Create and push tag
git tag v0.0.1-test
git push origin v0.0.1-test

# View Actions run status
# https://github.com/YOUR_ORG/YOUR_REPO/actions
```

## Usage Instructions

### Publishing a New Version

1. **Prepare Code**
   ```bash
   git checkout main
   git pull origin main
   ```

2. **Create Version Tag**
   ```bash
   # Format: v<major>.<minor>.<patch>
   git tag v1.0.0
   ```

3. **Push Tag**
   ```bash
   git push origin v1.0.0
   ```

4. **Monitor Build**
   - Navigate to `Actions` tab
   - View "Build and Push Private Network Docker Image" workflow
   - Wait for build completion (approximately 2-3 minutes)

5. **Verify Image**
   ```bash
   # Pull image
   docker pull harbor.example.com/st-chain/me_hub:v1.0.0_private_net

   # Run test
   docker run -d \
     -p 36657:36657 -p 1318:1318 -p 9545:9545 -p 8090:8090 \
     --name mechain-test \
     harbor.example.com/st-chain/me_hub:v1.0.0_private_net
   ```

## Image Tag Rules

The workflow creates two tags for each image:

1. **Version Tag**
   - Format: `{HARBOR_REGISTRY}/st-chain/me_hub:{git_tag}_private_net`
   - Example: `harbor.example.com/st-chain/me_hub:v1.0.0_private_net`
   - Purpose: Specific version of the private network image

2. **Latest Tag**
   - Format: `{HARBOR_REGISTRY}/st-chain/me_hub:latest_private_net`
   - Example: `harbor.example.com/st-chain/me_hub:latest_private_net`
   - Purpose: Always points to the latest built private network image

## Genesis Accounts Configuration

The private network image comes pre-configured with the following genesis accounts:

| Account Name | Initial Balance | Purpose |
|--------------|----------------|---------|
| alice | 10,000 MEC | Test Account 1 |
| bob | 10,000 MEC | Test Account 2 |
| charlie | 10,000 MEC | Test Account 3 |
| david | 10,000 MEC | Test Account 4 |

Additionally, the following default accounts are included:
- `global_dao` - DAO account
- `pools` - AMM pool account
- `user` - Regular user account
- `sequencer` - Sequencer account

## Troubleshooting

### Issue 1: Build Failed - "make: command not found"

**Cause:** Self-hosted runner missing required tools

**Solution:**
```bash
# Install on runner machine
sudo apt-get update
sudo apt-get install -y make build-essential
```

### Issue 2: Push Failed - "unauthorized: authentication required"

**Cause:** Harbor authentication failure

**Check:**
1. Verify Secrets are configured correctly
2. Confirm Robot Account has push permissions
3. Check if Token has expired

**Solution:**
```bash
# Manually test login
echo "YOUR_PASSWORD" | docker login YOUR_REGISTRY -u YOUR_USERNAME --password-stdin
```

### Issue 3: Build Timeout

**Cause:** Self-hosted runner resource insufficiency

**Solution:**
1. Check runner machine resources (CPU, memory, disk)
2. Clean up Docker cache:
   ```bash
   docker system prune -a -f
   ```

### Issue 4: Go Module Download Failed

**Cause:** Network issues or proxy configuration

**Solution:**
```bash
# Configure Go proxy on runner machine
export GOPROXY=https://goproxy.cn,direct
# or
export GOPROXY=https://proxy.golang.org,direct
```

## Advanced Configuration

### Modify Genesis Accounts

Edit `.github/workflows/build-push-private-net.yml`:

```yaml
env:
  GENESIS_ACCOUNTS: "alice:10000000000000000000000umec,bob:10000000000000000000000umec,custom:50000000000000000000000umec"
```

### Add Additional Image Tags

Add to the "Tag and push image" step in the workflow:

```yaml
- name: Tag and push image
  run: |
    docker tag me-hub/private-net:build-temp ${{ steps.image.outputs.tag }}
    docker push ${{ steps.image.outputs.tag }}

    # Add custom tag
    docker tag me-hub/private-net:build-temp ${{ secrets.HARBOR_REGISTRY }}/st-chain/me_hub:stable_private_net
    docker push ${{ secrets.HARBOR_REGISTRY }}/st-chain/me_hub:stable_private_net
```

### Add Notifications

Install notification step (e.g., send to Slack):

```yaml
- name: Send notification
  if: success()
  run: |
    curl -X POST ${{ secrets.SLACK_WEBHOOK_URL }} \
      -H 'Content-Type: application/json' \
      -d '{
        "text": "✅ Private network image built successfully: ${{ steps.image.outputs.tag }}"
      }'
```

## Security Recommendations

1. **Use Robot Accounts**
   - Do not use personal Harbor accounts
   - Create dedicated Robot Accounts for CI/CD
   - Set reasonable expiration times

2. **Principle of Least Privilege**
   - Grant Robot Accounts only necessary push permissions
   - Do not grant administrator privileges

3. **Rotate Credentials Regularly**
   - Regularly update Robot Account Tokens
   - Record Token update times

4. **Monitor Build Logs**
   - Regularly review Actions logs
   - Watch for unusual build activity

## Related Documentation

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Docker Build Documentation](../docker/README.md)
- [Genesis Accounts Configuration](../docker/GENESIS_ACCOUNTS.md)
- [Harbor Documentation](https://goharbor.io/docs/)

## Support

If you have issues:
1. Check the troubleshooting section in this document
2. View GitHub Actions run logs
3. Submit an Issue in the project repository

---

**Last Updated:** 2025-11-12
**Maintainer:** Me-Chain Team
