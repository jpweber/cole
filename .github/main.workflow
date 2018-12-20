workflow "Build and Push" {
  on = "push"
  resolves = [
    "Docker Login",
    "docker push",
  ]
}

action "docker build" {
  uses = "actions/docker/cli@76ff57a"
  args = "build -t cole ."
}

action "Docker Login" {
  uses = "actions/docker/login@76ff57a"
  secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
}

action "Docker Tag" {
  uses = "actions/docker/tag@76ff57a"
  args = ["$IMAGE_NAME", "$CONTAINER_REGISTRY_PATH/$IMAGE_NAME"]
  needs = ["docker build", "Docker Login"]
  env = {
    IMAGE_NAME = "cole"
    CONTAINER_REGISTRY_PATH = "jpweber"
  }
}

action "docker push" {
  uses = "actions/docker/cli@76ff57a"
  needs = ["Docker Tag"]
  secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
  args = "[\"push\",\"$CONTAINER_REGISTRY_PATH/$IMAGE_NAME\"]"
    env = {
    IMAGE_NAME = "cole"
    CONTAINER_REGISTRY_PATH = "jpweber"
  }
}
