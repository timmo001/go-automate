// Generates the Open Graph share image at public/og.png.
// Run with: pnpm og
import { readFile } from 'node:fs/promises';
import { fileURLToPath } from 'node:url';
import path from 'node:path';
import sharp from 'sharp';

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');

const W = 1200;
const H = 630;

// Embed the logo as a nested SVG at a fixed position and size.
const logoRaw = (await readFile(path.join(root, 'src/assets/logo.svg'), 'utf8')).trim();
const logo = logoRaw.replace(
  '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 36 36">',
  '<svg x="110" y="165" width="300" height="300" viewBox="0 0 36 36">',
);

const fontStack = "Inter, 'Liberation Sans', 'DejaVu Sans', Arial, sans-serif";

const svg = `<svg xmlns="http://www.w3.org/2000/svg" width="${W}" height="${H}" viewBox="0 0 ${W} ${H}">
  <rect width="${W}" height="${H}" fill="#0e7490"/>
  ${logo}
  <g font-family="${fontStack}" fill="#ffffff">
    <text x="470" y="300" font-size="96" font-weight="700">Go Automate</text>
    <text x="474" y="372" font-size="38" fill-opacity="0.88">Run common tasks and trigger</text>
    <text x="474" y="420" font-size="38" fill-opacity="0.88">Home Assistant from anywhere.</text>
    <text x="474" y="500" font-size="30" fill-opacity="0.7">go-automate.timmo.dev</text>
  </g>
</svg>`;

await sharp(Buffer.from(svg)).png().toFile(path.join(root, 'public/og.png'));
console.log('Wrote public/og.png');
