import { copyFile, mkdir, readFile } from "node:fs/promises";
import { createRequire } from "node:module";
import path from "node:path";
import process from "node:process";
import { fileURLToPath } from "node:url";

const require = createRequire(import.meta.url);
const scriptDir = path.dirname(fileURLToPath(import.meta.url));
const rootDir = path.resolve(scriptDir, "..");
const packageEntry = require.resolve("flexsearch");
const sourcePath = path.join(path.dirname(packageEntry), "flexsearch.bundle.min.js");
const targetPath = path.join(rootDir, "assets", "js", "vendor", "flexsearch.bundle.min.js");
const checkOnly = process.argv.includes("--check");

await mkdir(path.dirname(targetPath), { recursive: true });

if (checkOnly) {
  const [sourceContent, targetContent] = await Promise.all([
    readFile(sourcePath, "utf8"),
    readFile(targetPath, "utf8"),
  ]);

  if (sourceContent !== targetContent) {
    console.error("Vendored FlexSearch bundle is out of date. Run `npm run vendor:flexsearch`.");
    process.exit(1);
  }

  console.log("Vendored FlexSearch bundle matches the pinned npm package.");
} else {
  await copyFile(sourcePath, targetPath);
  console.log(`Updated ${path.relative(rootDir, targetPath)} from ${sourcePath}.`);
}
