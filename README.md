ros7api
===

Experimental (don't use!) Mikrotik RouterOS 7 REST client for Go.

Attempting to autogenerate as much as possible - see `gen/kinds.proto` and `gen/types.text.pb`.

Regenerating
---

    $ nix-shell
    $ rm ros/zz*
    $ go generate

Bazel Integration
---

Soon, probably in [hscloud](https://hackdoc.hackerspace.pl/).
