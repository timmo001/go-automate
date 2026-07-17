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
- `pnpm deploy` (deploy the built site to Cloudflare Workers)
- `pnpm deploy:preview` (upload a preview version without promoting it)
- `pnpm preview`

## Deployment

The site deploys to Cloudflare Workers as static assets. Workers Builds uses:

- Root directory: `docs`
- Production branch: `main`
- Build command: `pnpm build`
- Deploy command: `pnpm deploy`
- Non-production deploy command: `pnpm deploy:preview`

`wrangler.jsonc` owns the Worker name, compatibility date, custom domain, asset directory, and 404 behaviour. The site is fully static, so it does not use an Astro adapter or invoke Worker code for page requests.
