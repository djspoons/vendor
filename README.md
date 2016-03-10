# Vendor

Vendor copies go dependencies to ./vendor.

## Install

	$ go get github.com/bmizerany/vendor

## Use

Vendor all dependencies for a main package:

	$ vendor ./cmd/myapp

Update an already vendored dependency:

	$ vendor github.com/lib/pq

Got an admin or deploy script with dependencies not in the "main" packages?

	$ vendor admin.go

It just works. That's all.

## Motivation

We wanted a simple vendoring tool that stays out of our way and that makes less
assumptions about what we want and don't want in our vendor directory. Instead
vendor makes assumptions about how you use it and what you ask of it. These
assumptions are:

* Dependencies should be copied to ./vendor.

Vendor assumes you're running it where you vendor directory to exists under or
should be created.

* Dependencies in their working directories are just they way you want them to be under ./vendor.

Vendor does not check if there are uncommitted changes in the dependencies working
directory. The check is costly in terms of time and usually gets in the way
when you're iterating on one package and vendoring it in another to test.

* Dependencies aren't always `import`ed.

We have a script that `go install`s `vendor/github.com/backplaneio/tools/cmd/bpagent`
and runs it for development. Everyone on the team is using the same binary,
guaranteed.

Because this isn't imported, other dependency tools won't allow us to vendor it
and keep it up to date so we resulted to manually `cp -R` it and all of it's
dependencies. This was fine the first few times until @voutasaurus threaten to
burn the office down, so I made sure this tool could do it.

* Dependencies that exist in subdirectories of the working directory vendor is run in are to be ignored.

Vendor will not vendor anything that is already in the working tree.

* You're using source control

Vendor isn't careful because it assumes you are and can rollback any mistakes
you may make while using vendor. Vendor doesn't make mistakes. :)

## Goals

* Be fast
* Get out of the way
* Do what is asked of it

