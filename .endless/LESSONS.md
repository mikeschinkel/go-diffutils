# Lessons Learned

A capture log of corrections, newest last. Entries are appended by
`endless lesson write`, which commits this file on the main checkout as it
writes — so recording a lesson never waits on a task landing, and never
dirties a worktree.

## Format

```
### [YYYY-MM-DD] One-line summary
<the lesson — what went wrong, why, and the rule that replaces it>
- **Project**: <project>
```

---

## Entries

### [2026-10-27] Release/CI workflow improvements for go-pkgs repos belong to the gomion project
Context: after E-2449 (publish go-diffutils) I filed E-2458 'Keep the root release marked latest after nested releases' against go-diffutils, covering go-diffutils and go-dt, and recommended folding the GitHub Actions warning fixes into it. Mike: 'Isn't that more of a Gomion project task? ... Again, wrong project, no?' Reason: the change is to the release workflow model shared across go-pkgs repos, which gomion owns (gomion already has E-2454 'Bring go-pkgs Release workflows to go-dt's model'); filing it in the repo where it was noticed scatters it. Applies to this kind of cross-repo release/CI workflow work; before filing, check whether gomion already has a task it fits.
- **Project**: go-diffutils
