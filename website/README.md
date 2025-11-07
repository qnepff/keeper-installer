# QuickNEasy Website

This directory contains the source files for **https://quickneasy.info**

## Files

- `index.html` - Main landing page
- `style.css` - Styles

## Automatic Deployment

These files are automatically deployed to BunnyCDN on every push to `main`:

1. GitHub Actions builds installers (Linux, Windows, macOS)
2. Uploads installers to BunnyCDN storage zone `qne-installers-cdn`
3. Uploads website files to BunnyCDN storage zone `quickneasy-info`
4. Purges CDN cache

**CDN URL:** https://quickneasy-info-cdn.b-cdn.net  
**Domain:** https://quickneasy.info (after DNS)

## Download URLs

The website detects user's OS and provides the correct installer:

- **Linux:** `https://qne-installers-cdn-cdn.b-cdn.net/keeper-installer-linux-amd64`
- **Windows:** `https://qne-installers-cdn-cdn.b-cdn.net/keeper-installer-windows-amd64.exe`
- **macOS:** `https://qne-installers-cdn-cdn.b-cdn.net/keeper-installer-darwin-amd64`

## Local Testing

```bash
# Serve locally
python3 -m http.server 8000

# Visit: http://localhost:8000
```

## Manual Upload

If needed, you can manually upload:

```bash
export BUNNY_API_KEY=your-key

curl -X PUT \
  -H "AccessKey: $BUNNY_API_KEY" \
  -T index.html \
  "https://storage.bunnycdn.com/quickneasy-info/index.html"

curl -X PUT \
  -H "AccessKey: $BUNNY_API_KEY" \
  -T style.css \
  "https://storage.bunnycdn.com/quickneasy-info/style.css"
```
