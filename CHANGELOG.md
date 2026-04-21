# Changelog

## Unreleased

- Add manifest-driven sync for shared maintainer files in downstream lesson repositories.
- Add reusable GitHub Actions workflows, `release-please`, and Commitizen-based conventional-commit enforcement.
- Document the `sync-template-files.sh` update flow and the preferred local `prek` hook setup.

## [0.3.0](https://github.com/oer-particle-physics/hugo-styles/compare/v0.2.1...v0.3.0) (2026-04-21)


### Features

* add lesson metadata shortcode ([dac1650](https://github.com/oer-particle-physics/hugo-styles/commit/dac16507de2d343f7dfde3f2a37f2dfcbe8cb905)), closes [#15](https://github.com/oer-particle-physics/hugo-styles/issues/15)
* make challenge subtitles optional ([790b0a4](https://github.com/oer-particle-physics/hugo-styles/commit/790b0a4eb82ffbdf1e0cd1b968135e82ba9dbb77)), closes [#16](https://github.com/oer-particle-physics/hugo-styles/issues/16)


### Bug Fixes

* render keypoints as inline markdown ([2d72e67](https://github.com/oer-particle-physics/hugo-styles/commit/2d72e676c70455a8a6544415a88c68fcadce8d90)), closes [#17](https://github.com/oer-particle-physics/hugo-styles/issues/17)

## [0.2.1](https://github.com/oer-particle-physics/hugo-styles/compare/v0.2.0...v0.2.1) (2026-04-19)


### Bug Fixes

* quote refresh workflow PR title ([b8d9d81](https://github.com/oer-particle-physics/hugo-styles/commit/b8d9d8189b545d355e307a1c2e154168198ac166))
* run upstream sync wrapper via bash ([4ecabf2](https://github.com/oer-particle-physics/hugo-styles/commit/4ecabf295d6d6ef308a1f3b99dab6676d51cd8a4))

## [0.2.0](https://github.com/oer-particle-physics/hugo-styles/compare/v0.1.0...v0.2.0) (2026-04-19)


### Features

* add lesson schedule and overview as well as authors ([4e6a2e6](https://github.com/oer-particle-physics/hugo-styles/commit/4e6a2e6250de1a7ed9e5f0750eab4bad2463b145))
* add shared template sync and release automation ([5bdae95](https://github.com/oer-particle-physics/hugo-styles/commit/5bdae95894ed4590796f57cfe524b3a75dbcbff8))
* add versioning documentation and streamline docs ([070265b](https://github.com/oer-particle-physics/hugo-styles/commit/070265b1d60f7d4ed8c69ce810235c78aed7a316))
* enable config-driven versioned site builds ([9320a98](https://github.com/oer-particle-physics/hugo-styles/commit/9320a9870a3fb59ac7ed3a835b4d4cfb572708f1))
* Improve footer with copyright and citation ([98b037c](https://github.com/oer-particle-physics/hugo-styles/commit/98b037c79f20b147c4b7baa0fe4b6a1273bba694))
* **layout:** Further improvements to single pages ([a82c845](https://github.com/oer-particle-physics/hugo-styles/commit/a82c8457c2c7845b02cdeb8502facf4829388ce7))
* **layout:** Reduce layout customisation ([a97999d](https://github.com/oer-particle-physics/hugo-styles/commit/a97999d65832eb25128e9c3d7af8e641ffb97355))


### Bug Fixes

* another attempt to fix image handling for migration ([fb26f68](https://github.com/oer-particle-physics/hugo-styles/commit/fb26f680421dde8ab8ca5d45fda7c5dd90ae5418))
* **cmd:** improve migration script ([33a7323](https://github.com/oer-particle-physics/hugo-styles/commit/33a7323522163bcf026cc718cbedfb5897684663))
* **cmd:** migration images regression ([5d06579](https://github.com/oer-particle-physics/hugo-styles/commit/5d0657936c2d71cf201f99b2f1e61bb24be345af))
* correct version menu links ([105b5e2](https://github.com/oer-particle-physics/hugo-styles/commit/105b5e2796c75cba0a0318dd38c7c5ddc05c0ec1))
* **docs:** Improve documentation ([afbe79a](https://github.com/oer-particle-physics/hugo-styles/commit/afbe79a9a724a27e764e328ce344cd90cfc6191a))
* Further layout fixes ([bf388ce](https://github.com/oer-particle-physics/hugo-styles/commit/bf388ced380af81ea8879653424dbcec2ac4e948))
* **layout:** Further layout fixes ([1f7e4e6](https://github.com/oer-particle-physics/hugo-styles/commit/1f7e4e6d88094e5b869eaa080cfa6e4c93579994))
* **layout:** navbar to show episodes ([5766e41](https://github.com/oer-particle-physics/hugo-styles/commit/5766e411918c1449facb17ef8098f118b87a3b65))
* **migration:** better support for images ([44e9e47](https://github.com/oer-particle-physics/hugo-styles/commit/44e9e47969cfce80ca3e11b530d6a5775ec75ab0))
* Page layout ([25c5f5c](https://github.com/oer-particle-physics/hugo-styles/commit/25c5f5caf3c1ae250d1fe0dd3a31f4fbb964da63))
* use plain semver tags for release-please ([6355f4c](https://github.com/oer-particle-physics/hugo-styles/commit/6355f4c4e22eff93d7cebb784a242eb55e7dc3d7))

## 0.1.0 - 2026-04-11

- Initial shared-module release.
- Added lesson layouts, pedagogy shortcodes, aggregated resource pages, and audience toggle.
- Added Hextra-aligned tabs, search, and aggregated resource navigation.
- Added Hugo-native lesson validation for metadata, duplicate weights, glossary/profile refs, and image alt text.
- Added regression fixtures and tests for the migration/check command.
- Expanded the built-in documentation set for authoring, deployment, troubleshooting, migration, and updates.
