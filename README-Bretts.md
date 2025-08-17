# Brett's Notes to Self

**Update**: Windows compilation issues have been resolved! The project now supports native Windows compilation using stub dependency resolvers.

When I originally tried to work on this, I had trouble getting it to compile and run from a Windows machine. So instead I went to a Xubuntu machine (Ellie's old Dell Precision) and had no trouble. I think some of the libraries needed don't work right off the bat in Windows.

**Resolution**: The tree-sitter dependency issues on Windows have been solved by implementing Windows-specific stub resolvers that bypass the problematic CGO dependencies while maintaining core functionality.

# Compilation Instructions

See [How to Compile](NotesForDeveloper/How%20to%20Compile.md) for step-by-step instructions on building this project for Linux and Windows. This guide explains how to compile the Go code and how to cross-compile for Windows from Linux.
