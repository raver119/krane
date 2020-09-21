**What is Krane?**

Krane is very simple parallel executor for Docker&trade; builder. The tool to build lots of docker images in parallel. 

**What for?**

If you have to build multiple images frequently, and you have powerful dev machine - you can get a significant speedup if you'll build images in parallel.

**How it works?**

Quite trivial: you provide configuration in YAML format, Krane applies basic dependency analysis, and executes build with respect to graph topology. 

**Got questions?**

File an issue right here, or drop me a line: [raver119@gmail.com](mailto:raver119@gmail.com)