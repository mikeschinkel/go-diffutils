gomion E-2454 added `gomion repo setup` and `gomion repo workflows`. They generate a repo's .github/workflows/test.yml and release.yml from the go.mod files git tracks, and raise every go directive below Go 1.27 to 1.27. The generated workflows follow go-dt's Release model:
- a module choice, with vX.Y.Z tags for root and <dir>/vX.Y.Z for nested modules;
- refusal of a wrong module path, a replace directive or an existing tag;
- GOWORK=off tests and per-module release notes, with .endless excluded from root's notes.

They also fix the model's gaps:
- root kept marked latest;
- Go pinned to 1.27, so GOEXPERIMENT=jsonv2 is no longer set;
- ubuntu-26.04 instead of ubuntu-latest;
- Node 24 action releases;
- no setup-go cache, which warned when there was no root go.sum.

Mike's rule is that Actions warnings and notices are removed wherever possible, never accepted as informational.

On 2026-10-28 Mike decided that E-2454 changes only gomion and files one task per go-pkgs project, so that each change to go-diffutils is made by a session for go-diffutils itself. This is that task. The generated files carry a do-not-edit header; regenerate them with gomion rather than editing them by hand.
