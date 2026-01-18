# Tailwind CSS v4.1 Setup Guide

Cacto CMS uses **Tailwind CSS v4.1** with PostCSS for optimal performance and minimal dependencies.

## 🏗️ Architecture

**Why PostCSS?**
- ✅ **Minimal dependency**: Only PostCSS + Tailwind v4 plugin
- ✅ **Go-friendly**: No build tool required (Vite not needed)
- ✅ **Tree-shaking**: Only used classes in `output.css`
- ✅ **Enterprise-grade**: Production-ready, optimized builds

```
web/static/css/
├── input.css          # Tailwind v4 source (@import "tailwindcss")
└── output.css         # Generated CSS (only used classes, gitignored)

postcss.config.js      # PostCSS configuration
package.json           # Node.js dependencies (minimal)
```

## 🚀 Quick Start

### 1. Install Dependencies

```bash
# Install Go + Node.js dependencies
make install

# Or install Node.js dependencies manually
npm install
```

### 2. Development

**Option A: Separate terminals (recommended)**
```bash
# Terminal 1: Watch CSS changes (only used classes)
make css-watch

# Terminal 2: Run Go server with hot reload
make dev
```

**Option B: Single terminal**
```bash
# Run CSS watch in background
make css-watch &

# Run Go server
make dev
```

### 3. Production Build

```bash
# Build CSS + Go binaries
make build
```

This will:
- Generate optimized, minified CSS (only used classes)
- Build Go binaries
- Ready for deployment

## 📝 Makefile Commands

| Command | Description |
|---------|-------------|
| `make css` | Build Tailwind CSS v4 (production, only used classes) |
| `make css-watch` | Watch CSS changes (development, JIT) |
| `make dev` | Run Go server (reminds to run css-watch) |
| `make build` | Build everything (CSS + Go) |

## 🎨 Tailwind v4.1 Features

### What's New in v4.1

- **`@import "tailwindcss"`**: Single import replaces `@tailwind` directives
- **Automatic content detection**: Scans your templates automatically
- **CSS-first approach**: More CSS-native, less JS config
- **Better tree-shaking**: Only used classes in output

### Customization

Edit `web/static/css/input.css`:

```css
@import "tailwindcss";

@layer components {
  .my-custom-button {
    @apply px-4 py-2 bg-blue-500 text-white rounded;
  }
}
```

### Content Scanning

Tailwind v4 automatically scans:
- `./app/interfaces/templates/**/*.{templ,go}`
- `./app/interfaces/http/**/*.go`
- `./web/static/**/*.html`

No `tailwind.config.js` needed for basic usage! (Optional for theme customization)

## 🔧 Advanced Usage

### Manual CSS Build

```bash
# Development (JIT, unminified)
npx postcss ./web/static/css/input.css -o ./web/static/css/output.css --watch

# Production (minified, only used classes)
NODE_ENV=production npx postcss ./web/static/css/input.css -o ./web/static/css/output.css
```

### Custom Theme (Optional)

Create `tailwind.config.js` if needed:

```js
/** @type {import('tailwindcss').Config} */
module.exports = {
  theme: {
    extend: {
      colors: {
        primary: {
          600: '#2563eb',
        },
      },
    },
  },
}
```

## 🚢 Deployment

### Production Checklist

- [ ] Run `make build` (builds CSS + Go)
- [ ] Verify `web/static/css/output.css` exists
- [ ] Check CSS is minified and contains only used classes
- [ ] Test in production environment

### CI/CD Integration

```yaml
# Example GitHub Actions
- name: Install dependencies
  run: |
    make install
    
- name: Build
  run: |
    make build
```

## 📦 Dependencies

**Minimal Setup:**
- `@tailwindcss/postcss` - Tailwind v4 PostCSS plugin
- `postcss` - CSS processing
- `postcss-cli` - CLI tool
- `autoprefixer` - Vendor prefixes
- `cssnano` - Production minification (optional)

**Total size:** ~30MB (node_modules)

## 🎯 Best Practices

1. **Use @layer**: Organize custom styles in layers (base, components, utilities)
2. **Component classes**: Extract repeated patterns with `@apply`
3. **Content scanning**: Tailwind v4 auto-detects, but ensure templates are in standard locations
4. **Production**: Always minify CSS in production
5. **Watch mode**: Use `css-watch` during development for JIT compilation

## 🐛 Troubleshooting

### CSS not updating?

1. Check `css-watch` is running
2. Verify `output.css` exists
3. Hard refresh browser (Cmd+Shift+R / Ctrl+Shift+R)
4. Check browser console for 404 errors

### Classes not working?

1. Verify class is in template files (auto-detected)
2. Rebuild CSS: `make css`
3. Check PostCSS config is correct

### Build fails?

1. Check Node.js version: `node --version` (needs >=18)
2. Reinstall: `rm -rf node_modules && npm install`
3. Check `package.json` scripts

## 📚 Resources

- [Tailwind CSS v4 Docs](https://tailwindcss.com/docs)
- [PostCSS](https://postcss.org/)
