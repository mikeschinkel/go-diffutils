## Plan
1. Run `git status` and `git log <latest-tag>..HEAD` in the go-diffutils repo; list every module in it that has its own go.mod and tag prefix.
2. Commit any uncommitted work that belongs in the release. Ask the owner about anything unclear instead of guessing.
3. Run the repo's build, vet and tests. Fix or report any failure.
4. Push main to origin. Do not create or push tags.
5. For each module, compare the changes since its latest tag and propose a semver version, explaining whether the change is breaking.
6. Draft a release description for each module.
7. Report the proposed versions and descriptions to the owner. The owner tags and releases when ready.
