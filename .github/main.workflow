workflow "Benchmark" {
  on = "push"
  resolves = ["gobenchdata to gh-pages"]
}

action "filter" {
  uses = "actions/bin/filter@master"
  args = "branch master"
}

action "gobenchdata to gh-pages" {
  uses = "bobheadxi/gobenchdata@v0.2.0"
  needs = ["filter"]
  secrets = ["GITHUB_TOKEN"]
  env = {
    PRUNE = "20"
    GO_BENCHMARK_PKGS = "./benchmarks/..."
  }
}
