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
