# Grax Homebrew Tap

This tap publishes the `grax` Homebrew formula and package assets.

- Formula: [`Formula/grax.rb`](Formula/grax.rb)
- Package files:
  - [`Modelfile`](Modelfile)
  - [`docs/grax_package.md`](docs/grax_package.md)
  - [`model/parsers/grax.go`](model/parsers/grax.go)
  - [`model/renderers/grax.go`](model/renderers/grax.go)
- Install with explicit Homebrew 6.0.0 tap trust as needed

## Install

```shell
brew tap thebreathwright/grax
brew trust --formula thebreathwright/grax/grax
brew install grax
```
