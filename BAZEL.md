# Bazel Build Experiment

Notes on building Highlander with Bazel.

## What is Bazel

Bazel is a build system developed by Google. It is language-agnostic, meaning
it can run many kinds of builds. It is extensible via a Starlark build
langauage, which looks like a subset of Python3. But many use cases do not need
to write this language at all, and can simply use a CMake-ish config.

## Why try this?

There is a chance we can reduce the maintenance burden of adding languages,
container builds, and new build targets.

## Install Bazel CLI

Get an _installer script_ from the [GitHub releases page](https://github.com/bazelbuild/bazel/releases).

Run the installer script with `--user`. Adds bazel to **$HOME/bin** and drops a
**$HOME/.bazelrc**.

Bazel is written in Java, mostly, and bundles its own JDK.

## Concepts

WORKSPACE sits in the root of the repo. If a nested WORKSPACE is encountered, it
is ignored as another workspace that must be explicitly dependended on.

BUILD or BUILD.bazel files sit in the the root of packages. Packages define
related files and dependencies. The name of a package is the name of the
directory containing its BUILD file, relative to the top-level directory of the
source tree.

Packages contain **targets**. Targets are either:

* files
  * source files - handwritten, checked in
  * generated files - generated from source according to rules, not checked in
* rules
  * outputs of rules are always generated files
  * inputs may be source files, or generated files
  * outputs of one rule may be inputs to another (chaining)
  * inputs to a rule may be _other rules_. See [targets docs](https://docs.bazel.build/versions/1.1.0/build-ref.html#targets)

The name of a target is called a label. `@` refers to the root workspace, so the
following two labels are equivalent

```
@myrepo//my/app/main:app_binary
```

```
//my/app/main:app_binary
```

In this case, **my/app/main** is the package name and **app_binary** is the target.
Targets only ever belong to a single package.

Many portions of this fully-qualified label path can be omitted when using labels
inside BUILD files, not just the workspace. See the [label docs](https://docs.bazel.build/versions/1.1.0/build-ref.html#labels)

But from other packages or the command line, the full label path must be used.

So labels are kind of like URLs, with their own pathing spec.


## BUILD (or BUILD.bazel) files

BUILD files are a top-to-bottom sequence of statements _interpreted_ by the
Starlark language interpreter. Order matters. Vars must be defined before use.

BUILD files can import Starlark files like this to import rules, functions, or constants.

```
load("//foo/bar:file.bzl", "some_library")
```

Common dependencies are: src, deps, data

* src - files consumed directly by the rule, or a rule that outputs source files
* deps - rules pointing to separately compiled modules that provide symbols,
  headers, etc
* data - bazel runs in an isolated dir where only these data files are available
  so we can use these for golden test files. We don't need these for build, but
  we need them to run the test.

## Bazel command line

[build](https://docs.bazel.build/versions/1.1.0/guide.html#the-build-command)

Individual targets can be specified with labels

```
bazel build //my/app/main
```

Target patterns can specify more than one target at once. [See docs](https://docs.bazel.build/versions/1.1.0/guide.html#target-patterns).
An example:

```
bazel build //foo/...   # all rules in all pkgs in directory foo
```

A negation prefix `-` can be provided with a list of target patterns

```
bazel build -- foo/... -foo/bar/...
```


## Caching

Bazel caches in ~/.cache/bazel/_bazel_$USER/cache/repos/v1/ and this cache is
NEVER cleaned up automatically. This ensures that, if an upstream file goes away,
you could still recover it from the cache. The cache is shared across workspaces
and versions of bazel.


## Sandboxed execution

Sandboxing works on modern linux and macos. Some host configurations disable
sandboxing by default, though. See [this if you get weird errors](https://docs.bazel.build/versions/1.1.0/guide.html#sandboxing).

## Build Phases

[Main docs on build phases](https://docs.bazel.build/versions/1.1.0/guide.html#phases)


## Client/Server architecture

Running `bazel` in a workspace starts a server process, if one isn't already
started. The server shuts down after 3 hours of inactivity.`


## Gazelle

Gazelle is a tool for automatically generating BUILD files for Go packages. Bazel
wants a BUILD file in the root of every Go package, and this means a lot of files
and a lot of upfront cost.

We add Gazelle to our WORKSPACE and a Gazelle-specific BUILD to our repo root
that runs Gazelle magic. Then we call it with bazel itself

```
bazel run //:gazelle

```

We MUST specify manual resolution for our go packages in the highlander repo due
to the "paperless" import root that we declare.

## Run Go tests

```
bazel run --verbose_failures //go/src/paperless/cmd/argo-repo:go_default_test
```

## CI TODOs


## Patches

Q: do we need empty import prefix?

The `importpath` in `go_library` needs to start with "paperless"

Add annotations for new libs we add:

```
cat ERRORS.txt | go run gazelle_helper.go >> BUILD.bazel 
```

## Protobuf?

https://github.com/bazelbuild/rules_go/blob/master/proto/core.rst#option-2-use-pre-generated-pbgo-files

Add `# gazelle:proto disable_global`



## Building a Go Binary

TODO

## Building JS projects

TODO

