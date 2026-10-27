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

### [2026-10-27] Mike wants CI warnings and notices eliminated, not dismissed as informational
Context: reviewing GitHub Actions annotations on go-diffutils' first Test and Release runs (Node 20 deprecation warnings, setup-go cache warnings for a missing root go.sum, ubuntu-latest moving to Ubuntu 26 notices). I called the Ubuntu notice 'informational' and offered 'accept it' as an option. Mike: 'I am a STRONG believer in making the effort to eliminate warnings and notices whenever possible because then when I see a warning or notice I pay attention rather than assuming it is nothing to worry about. So just ignoring them as just informational is not a preferred option.' Reason: a clean output keeps every new warning meaningful. When presenting options for a warning or notice, lead with how to make it go away; offer leaving it only when nothing can remove it, and say so.
- **Project**: go-diffutils
