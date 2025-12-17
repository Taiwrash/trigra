# TRIGRA Landing Page

A modern, dark-themed landing page for the TRIGRA project.

## Features

- 🎨 **Dark Mode Design** - Beautiful gradient background with animations
- 📋 **One-Click Copy** - Installation command with copy button
- ⚡ **Animated** - Smooth transitions and floating background elements
- 📱 **Responsive** - Works perfectly on all devices
- ⌨️ **Keyboard Shortcut** - Press Ctrl/Cmd + K to copy command

## Preview

Open `index.html` in your browser to see the landing page.

## Deployment

### GitHub Pages

1. Push `index.html` to your repository
2. Go to Settings → Pages
3. Select branch: `main`
4. Select folder: `/ (root)`
5. Save

Your page will be available at: `https://taiwrash.github.io/trigra/`

### Custom Domain

Add a `CNAME` file with your domain:

```bash
echo "gitops.yourdomain.com" > CNAME
git add CNAME
git commit -m "Add custom domain"
git push
```

## Local Development

```bash
# Open in browser
open index.html

# Or serve with Python
python3 -m http.server 8000
# Visit http://localhost:8000
```

## Customization

Edit `index.html` to customize:

- Colors: Change CSS variables in `:root`
- Features: Update the `.features` grid
- Links: Modify GitHub URLs
- Command: Update the install command

## Installation Command

The page displays:
```bash
curl -fsSL https://raw.githubusercontent.com/Taiwrash/trigra/main/install.sh | bash
```

Make sure `install.sh` is in your repository root!
