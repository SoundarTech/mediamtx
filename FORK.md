# Soundar MediaMTX Fork Policy

This repository is SoundarTech's private, minimally modified fork of
[MediaMTX](https://github.com/bluenviron/mediamtx).  It remains under the
upstream MIT license; see [LICENSE](LICENSE).

## Provenance

| Item | Value |
| --- | --- |
| Upstream repository | `https://github.com/bluenviron/mediamtx.git` |
| Imported upstream tag | `v1.20.1` |
| Imported upstream commit | `883194a19b7244355c9bc975c0574c9842733637` |
| Soundar maintenance branch | `soundar/v1.20.1` |
| Soundar baseline tag | `v1.20.1-soundar.0` |
| Governed fork release tag | `v1.20.1-soundar.1` |
| Functional delta at baseline | None |

The baseline and governed release tags are annotated but not cryptographically
signed.  A release manager must replace them with organization-signed release
tags before the fork is used in a production build or release artifact.

## Maintenance Rules

- Preserve upstream history and retain upstream license and notices.
- Integrate upstream releases explicitly into a new `soundar/vX.Y.Z` branch;
  do not merge directly into an existing released baseline.
- Keep Soundar changes minimal, isolated, reviewed and covered by tests.
- Record the upstream tag, commit, Soundar delta and security review in every
  release tag or release note.
- Do not enable a MediaMTX listener, control API or debug endpoint without an
  owning Platform adapter, authenticated boundary and updated threat model.
- Platform consumers must lock an exact Soundar tag and vendor the approved
  source; they must not consume the moving maintenance branch.

## Planned Platform Boundary

The fork will expose a narrow, testable embedding facade for Platform's media
adapter.  The facade must not call `os.Exit`, install global signal handlers,
perform self-updates or expose upstream internals as Platform API.
