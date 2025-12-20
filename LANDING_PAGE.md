# Trigra Landing Page

Trigra features a modern, dark-themed, glassmorphic landing page designed to impress and convert.

## ✨ Features

- 🎨 **Modern Aesthetics** - Deep space dark mode with vibrant accents.
- ⚡ **Interactive Installer** - Highlighted installation command with one-click copy.
- 💎 **Glassmorphism** - Semi-transparent cards with frosted glass effects.
- 📱 **Mobile First** - Perfectly responsive for developers on the go.
- 🌊 **Smooth Animations** - Subtle micro-animations and floating background particles.

## 🚀 Deployment

### GitHub Pages (Recommended)

1. Ensure `index.html` is in your repository root.
2. Go to **Settings** → **Pages**.
3. Set Source to `Deploy from a branch`.
4. Branch: `main`, Folder: `/ (root)`.
5. Your page will be live at `https://taiwrash.github.io/trigra/`.

### Automated Docs Deployment
We also include a `.github/workflows/deploy-docs.yml` which can automate the deployment of the full Starlight-based documentation site found in the `/website` directory.

## 🛠 Local Preview

To view the landing page locally:

```bash
# Simple Python server
python3 -m http.server 8000
# Visit http://localhost:8000
```

## 🎨 Customization

You can easily adapt the landing page to your style:
- **Typography**: Uses modern sans-serif fonts from Google Fonts.
- **Color Palette**: Controlled via CSS variables in the `<style>` section.
- **Copy Target**: The `copyCommand` function handles the "Click to Copy" logic.

---
View the live [Trigra Documentation](https://taiwrash.github.io/trigra/docs).
