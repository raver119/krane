**What is Krane?**

Krane is very simple parallel executor for Docker&trade; builder. The tool to build lots of docker images in parallel. 

**What for?**

If you have to build multiple images frequently, and you have powerful dev machine - you can get a significant speedup if you'll build images in parallel.

**How it works?**

Quite trivial: you provide configuration in YAML format, Krane applies basic dependency analysis, and executes build with respect to graph topology. 

**How to use?**
```
git cone https://github.com/raver119/krane
cd krane
go build -o krane .
```

Once you have the binary - write the build configuration in YAML format like this:

```yaml
build:
  - containerName: organiation/image:latest
    dockerpath: /path/to/Folder
    noCache: false
  - containerName: organiation/other_image:stable
    dockerpath: /path/to/OtherFolder
    noCache: false
threads: 12
```

Then just run it:

```
krane -f Path/To/File.yaml
```

If everything is ok, you'll see something like this:

```
Successfully built 6 images
```

**Is Minikube supported?**

Minikube has no need in any kind of special treatment. Just run `eval $(minikube docker-env)` before running Krane, and all new images in this session will use Minukube's internal registry. 

**Got questions?**

File an issue right here, or drop me a line: [raver119@gmail.com](mailto:raver119@gmail.com)