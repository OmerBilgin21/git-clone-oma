# OMA

Current accomplished tasks:
 * Initialize a repository, save the current snapshot of the current folder
 * Commit changes as you go, and hold only actionable diffs to be able to rebuild the old versions if need be
 * See the diffs for files with appropriate colors and accounting for deletions, additions and moves.
 * Build old versions of files using the version actions I have
 * Add a revert command with the ability to go back X versions
 * Maybe a log command to see the commit history? Not entirely sure if I want to do this one.
  
  
Remaining tasks that I envisioned for this project are:
 * ~Change the file snapshop on each 5 commits as that'd increase the rebuild, commit operations
 and the possibility of doing history +-5 are are vastly less then within 5 operations.~ this is not a good idea with file deletions on reverts
 * Concurrently build/render the diffs for the diff command

