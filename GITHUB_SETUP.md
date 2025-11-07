# GitHub Setup for Automated Deployment

## Required: GitHub Secret

The CI/CD workflow needs your BunnyCDN API key to upload files.

### Step 1: Get BunnyCDN API Key

Your API key should be stored in: `~/Dropbox/chameleon-info/bunny/`

If not there, get it from:
1. Login to https://panel.bunny.net
2. Go to Account → API
3. Copy the API Key

### Step 2: Add Secret to GitHub

1. Go to: https://github.com/qnepff/keeper-installer/settings/secrets/actions

2. Click **"New repository secret"**

3. Add secret:
   - **Name:** `BUNNY_API_KEY`
   - **Value:** Your BunnyCDN API key (paste it)

4. Click **"Add secret"**

### Step 3: Verify Setup

After adding the secret, the next push to `main` will:

1. ✅ Build installers (Linux, Windows, macOS)
2. ✅ Upload installers to BunnyCDN automatically
3. ✅ Upload website to BunnyCDN automatically
4. ✅ Purge CDN cache

**No manual steps needed!**

## How It Works

### On Every Push to `main`:

```
GitHub Actions
├── Build Linux installer → Upload to CDN
├── Build Windows installer → Upload to CDN
├── Build macOS installer → Upload to CDN
└── Deploy website → Upload to CDN → Purge cache
```

### CDN Endpoints

**Installers:** https://qne-installers-cdn-cdn.b-cdn.net/
- `keeper-installer-linux-amd64`
- `keeper-installer-windows-amd64.exe`
- `keeper-installer-darwin-amd64`

**Website:** https://quickneasy-info-cdn.b-cdn.net/
- `index.html`
- `style.css`

### Workflow File

Location: `.github/workflows/build.yml`

Key features:
- Multi-platform builds (3 OS)
- Direct BunnyCDN API uploads
- Cache purging
- GitHub Releases on tags

## Testing

1. Make a small change (e.g., edit README)
2. Commit and push to `main`
3. Watch Actions: https://github.com/qnepff/keeper-installer/actions
4. Verify files appear on CDN within 5-10 minutes

## Troubleshooting

**"Error: secrets.BUNNY_API_KEY not found"**
- Go back to Step 2 and add the secret

**"401 Unauthorized"**
- API key is wrong, update the secret

**"404 Not Found"**
- Storage zone name is wrong in workflow
- Check: `qne-installers-cdn` and `quickneasy-info`

**CDN not updating**
- Cache purge might take 1-5 minutes
- Force refresh: Ctrl+Shift+R

## Manual Override

If automation fails, you can still upload manually:

```bash
cd /home/paulf/qne-repos/keeper-installer
./upload-to-cdn.sh
```

## Next Steps

After secret is added:
1. ✅ Push this commit
2. ✅ Watch Actions tab for build
3. ✅ Verify CDN URLs work
4. ✅ Test website: https://quickneasy-info-cdn.b-cdn.net
5. ✅ Switch DNS when ready
