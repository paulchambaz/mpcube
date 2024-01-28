default:
  @just --list

watch:
  cargo watch -x run

watch-test:
  cargo watch -x test

run:
  cargo run

build:
  cargo build --release

test:
  cargo test

coverage:
  cargo tarpaulin

fmt:
  cargo fmt

vhs:
  nix shell -c vhs demo.tape

publish:
  # first we should run the full list of tests to be sure
  # it would be good to have the full list of operations to do to update here
  # that includes making all the modifications
  # we need to increase the version (probably with vipe)
  # we need to write a changelog (grab the version, then vipe)
  # then merging those modifications in master
  # then pushing the repo to github

publish-cargo:
  # the cargo version and changelog have been updated
  cargo publish

publish-nix:
  # TODO

publish-aur:
  # increase the build version - get it from Cargo.toml or even better from nix flake
  makepkg --printsrcinfo > .SRCINFO
  # git add PKGBUILD .SRCINFO
  # git commit -m "aur: $(LATEST_COMMIT_MESSAGE)"
  # git push aur master
