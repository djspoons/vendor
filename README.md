# Vendor

The Vendor tool copies go dependencies to the `./vendor` directory.

Vendor is the simplest way to ensure all necessary dependencies are located in
one place without complicated file management or hidden assumptions. It's all
about keeping it simple and straightforward.

NOTE: Before using Vendor, be sure the relevant dependencies are already in `GOPATH`.  

## Install

	$ go get github.com/bmizerany/vendor

## Use

After you've installed Vendor, it's easy to use.

To vendor all dependencies for a main package:

	$ vendor ./cmd/myapp

To update an already vendored dependency:

	$ vendor github.com/lib/pq

If you have an admin or deploy script with dependencies not in the "main" packages, use:

	$ vendor admin.go

It just works. That's all.

## Motivation

We wanted a simple vendoring tool that stays out of our way and that makes less
assumptions about what we want and don't want in our vendor directory. Instead
vendor makes assumptions about how you use it and what you ask of it. 

## Assumptions

Vendor assumes:

   * Dependencies will be copied to `./vendor`.

Vendor assumes you're running it in the parent directory of your vendor folder. If your vendor folder does not exist it will create it for you in the current directory.

   * Dependencies in their working directories are structured in they way you want them to be structured be under `./vendor`.

Vendor does not check if there are uncommitted changes in the dependencies'
working directory. The check usually gets in the way when you're iterating on
one package and vendoring it in another to test or experiment.

   * Dependencies that exist in subdirectories of the working directory vendor is run in are to be ignored.

Vendor will not vendor anything that is already in the working tree.

   * Dependencies are either not versioned or poorly versioned.

Vendor assumes that you would rather manage versions of dependencies yourself.

   * Dependencies that vendor dependencies are merely giving you a suggestion.

Vendor will not attempt to be clever and make a decision for you about what to
do when a dependency has vendored a dependency. Instead it reports and skips
them. You can run `vendor` again to fill the gaps in by copying your GOPATH
versions.

   * You're using source control.

The rollback function for Vendor is `git checkout -f -- vendor`. Vendor isn't
careful because it assumes you are and can rollback any mistakes you may make
while using vendor. Vendor doesn't make mistakes. :)

## Goals

* Be fast.
* Get out of the way.
* Do only what is asked of it.
