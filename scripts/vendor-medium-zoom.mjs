import { copyFile, mkdir, readFile } from "node:fs/promises";
import { createRequire } from "node:module";
import path from "node:path";
import process from "node:process";
import { fileURLToPath } from "node:url";

const require = createRequire(import.meta.url);
const scriptDir = path.dirname(fileURLToPath(import.meta.url));
const rootDir = path.resolve(scriptDir, "..");
const sourcePath = require.resolve("medium-zoom");
const targetPath = path.join(rootDir, "assets", "js", "vendor", "medium-zoom.min.js");
const checkOnly = process.argv.includes("--check");

await mkdir(path.dirname(targetPath), { recursive: true });

if (checkOnly) {
  const [sourceContent, targetContent] = await Promise.all([
    readFile(sourcePath, "utf8"),
    readFile(targetPath, "utf8"),
  ]);

  if (sourceContent !== targetContent) {
    console.error("Vendored Medium Zoom bundle is out of date. Run `npm run vendor:medium-zoom`.");
    process.exit(1);
  }

  console.log("Vendored Medium Zoom bundle matches the pinned npm package.");
} else {
  await copyFile(sourcePath, targetPath);
  console.log(`Updated ${path.relative(rootDir, targetPath)} from ${sourcePath}.`);
}
