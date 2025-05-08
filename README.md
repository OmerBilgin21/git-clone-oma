# OMA

Git clone with limited commands, without remote and worktrees.  
I will not make use of any library for the core logic as long as I can help it (except the stdlib of course)  
I will however, use DB clients, migration tools...etc I have no interest in implementing those myself.  
So far, there's only one library I used for something that could be considered as core logic is go-runewidth.  
Because well, I'm really not interested in doing width calculations (at least more than I'm already doing) for non-UTF-8 chars and stuff.  
You can check out the dependencies from [here](https://github.com/OmerBilgin21/git-clone-oma/network/dependencies) beware that some of those are my installed libraries' dependencies  

Current accomplished tasks:
 * Initialize a repository, save the current snapshot of each file in the current directory
 * Ability to ignore files/directories by name
 * Commit changes as you go, and hold only actionable diffs to be able to rebuild an old versions if need be
 * Commit only the changes that're on top of the latest commits. Don't commit a whole diff of current file and cached text
 * See the diffs for files with appropriate colors and accounting for deletions and additions on top of the latest commit
 * Revert command with the ability to go back X versions
 * Log command to see the commit history
 * Somewhat decent argument parsing and making use of flags for commands when needed

Diff view with fixed widths and coloring:
![Screenshot from 2025-05-08 18-51-21](https://github.com/user-attachments/assets/c8e0f3c6-ebb0-4b59-873d-dfbc473480c3)
  
There is not much else to show visually.
  
  
Remaining tasks that I envisioned for this project are:
 * ~Change the file snapshop on each 5 commits as that'd increase the rebuild, commit operations
 and the possibility of doing history +-5 are are vastly less then within 5 operations.~ this is not a good idea with *possible* file deletions on reverts
 * Integration tests for some git commands where it makes sense (commit, revert...etc)

Improvements that I'm aware of but won't do:
 * Concurrently build/render the diffs for the diff command
 * Improve performance of the diff algorithm
