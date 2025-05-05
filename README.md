# OMA

Git clone with limited commands, without remote and worktrees.  
I will not make use of any library for the core logic as long as I can help it (except the stdlib of course)  
I will however, use DB clients, migration tools...etc I have no interest in implementing those myself.  
So far, there's only one library I used for something that could be considered as core logic is go-runewidth.  
Because well, I'm really not interested in doing width calculations (at least more than I'm already doing) for non-UTF-8 chars and stuff.  
You can check out the dependencies from [here](https://github.com/OmerBilgin21/git-clone-oma/network/dependencies) beware that some of those are my installed libraries' dependencies  

Current accomplished tasks:
 * Initialize a repository, save the current snapshot of the current folder
 * Commit changes as you go, and hold only actionable diffs to be able to rebuild the old versions if need be
 * See the diffs for files with appropriate colors and accounting for deletions, additions and moves.
 * Build old versions of files using the version actions I have
 * Add a revert command with the ability to go back X versions
 * Maybe a log command to see the commit history?
 * Somewhat decent argument parsing and making use of flags for commands when needed
  
  
Remaining tasks that I envisioned for this project are:
 * ~Change the file snapshop on each 5 commits as that'd increase the rebuild, commit operations
 and the possibility of doing history +-5 are are vastly less then within 5 operations.~ this is not a good idea with *possible* file deletions on reverts
 * fix the FIXMEs and do the TODOs, then done.

Improvements that I'm aware of but won't do:
 * Concurrently build/render the diffs for the diff command
