workflow "CI" {
  on = "push"
  resolves = ["test"]
}

workflow "Release" {
  on = "push"
  resolves = ["goreleaser"]
}

action "is-tag" {
  uses = "actions/bin/filter@master"
  args = "tag"
}

action "not-tag" {
  uses = "actions/bin/filter@master"
  args = "not tag"
}

action "test" {
  uses = "./"
  args = "go test ./..."
  needs = ["not-tag"]
}

action "goreleaser" {
  uses = "./"
  secrets = [
    "GORELEASER_GITHUB_TOKEN"
  ]
  args = ["sh", "-c", "GITHUB_TOKEN=$GORELEASER_GITHUB_TOKEN goreleaser"]
  needs = ["is-tag"]
}
