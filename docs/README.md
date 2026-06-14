# 🎛️ Go Automate Docs

The [Go Automate](https://github.com/timmo001/go-automate) documentation site, built with Astro and Starlight.

## Site

The site is available at <https://go-automate.timmo.dev>.

## Project Structure

Content lives in `src/content/docs/` and is exposed as routes based on file names.

## Commands

All commands run from this `docs/` directory:

- `pnpm install`
- `pnpm dev`
- `pnpm build`
- `pnpm preview`

## Deployment

The site deploys to Vercel with the project **Root Directory** set to `docs/`.
Astro is auto-detected (build `astro build`, output `dist`).
