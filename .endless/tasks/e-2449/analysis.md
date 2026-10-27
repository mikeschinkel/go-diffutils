Decisions from Mike on 2026-10-27 (asked from go-dt's ES-1369):

- License: Mike's goal is to follow go-difflib's license in letter and spirit; he doesn't monetize the package. go-difflib is BSD 3-Clause (its LICENSE has the non-endorsement clause), not BSD 2-Clause as the old LICENSE said, so switching go-diffutils to BSD 2-Clause is not possible. MIT plus Mezard's retained notice meets BSD 3-Clause's terms (keep the copyright line, conditions and disclaimer); switching to BSD 3-Clause would add only the non-endorsement clause for Mike's own code. A clean-room rewrite was rejected as not worth the effort. Mike delegated writing the license text.
- Versions: v0.1.0 for all three, like his other go-pkgs.
- dmp and sgdiff require root by pseudo-version (option A) rather than a second land after the root tag (option B), a follow-up task (option C) or a hand-made root tag (option D). A pseudo-version is lower than v0.1.0, so a consumer requiring root v0.1.0 alongside dmp gets v0.1.0 by minimum version selection.
- doterr.go in root and sgdiff is the per-package doterr pattern, not duplication.
- go-tealeaves imports only the root package, so E-2440 needs only root's tag.
