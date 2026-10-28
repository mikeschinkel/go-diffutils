# Plan: adopt gomion's generated workflows in go-diffutils

1. Use a gomion build that includes E-2454's `gomion repo` commands.
2. In go-diffutils, run `gomion repo setup`. It raises each go.mod below Go 1.27 and tidies it, then writes .github/workflows/test.yml and release.yml.
3. Run `gomion repo setup --check`; it must report the repo up to date.
4. Run each module's tests locally with GOWORK=off.
5. Commit the generated workflows and the raised go.mod/go.sum files together.

After the land: push, then read the annotations on the first Test run (`gh api repos/<owner>/<repo>/check-runs/<job id>/annotations` for each job), and on the next Release run once there is one. Each must show zero warnings and zero notices; any that remain are fixed in gomion's templates, not here.
