#!/usr/bin/env bash

# push-docker-tags.sh
#
# Run from ci to tag images based on the current branch or tag name. 
# A bit like dockerhub autobuild config, but somewhere we can version control it.
# 
# The `docker-build` job in .circleci/config.yml builds the current commit
# in docker and tags it as ipfs/go-ipfs:wip
#
# Then the `docker-publish` job runs this script to decide what tag, if any,
# to publish to dockerhub.
#
# Usage:
#   ./push-docker-tags.sh <build number> <git commit sha1> <git branch name> [git tag name] [dry run]
#
# Example:
#   # dry run. pass a 5th arg to have it print what it would do rather than do it.
<<<<<<< HEAD
#   ./push-docker-tags.sh $(date -u +%F) testingsha master "" dryrun
#    
#   # push tag for the master branch
#   ./push-docker-tags.sh $(date -u +%F) testingsha master
#
#   # push tag for a release tag
#   ./push-docker-tags.sh $(date -u +%F) testingsha release v0.5.0
#
#   # Serving suggestion in circle ci - https://circleci.com/docs/2.0/env-vars/#built-in-environment-variables
#   ./push-docker-tags.sh $(date -u +%F) "$CIRCLE_SHA1" "$CIRCLE_BRANCH" "$CIRCLE_TAG"
=======
#   ./push-docker-tags.sh 1 testingsha mybranch v1.0 dryrun
#    
#   # push tag for the master branch
#   ./push-docker-tags.sh 1 testingsha master
#
#   # push tag for a release tag
#   ./push-docker-tags.sh 1 testingsha release v0.5.0
#
#   # Surving suggestion in circle ci - https://circleci.com/docs/2.0/env-vars/#built-in-environment-variables
#   ./push-docker-tags.sh "$CIRCLE_BUILD_NUM" "$CIRCLE_SHA1" "$CIRCLE_BRANCH" "$CIRCLE_TAG"
>>>>>>> feat: docker build and tag from ci
#
set -euo pipefail

if [[ $# -lt 3 ]] ; then
  echo 'At least 3 args required. Pass 5 args for a dry run.'
  echo 'Usage:'
  echo './push-docker-tags.sh <build number> <git commit sha1> <git branch name> [git tag name] [dry run]'
  exit 1
fi

BUILD_NUM=$1
GIT_SHA1=$2
GIT_SHA1_SHORT=$(echo "$GIT_SHA1" | cut -c 1-7)
GIT_BRANCH=$3
GIT_TAG=${4:-""}
DRY_RUN=${5:-false}

WIP_IMAGE_TAG=${WIP_IMAGE_TAG:-wip}
IMAGE_NAME=${IMAGE_NAME:-ipfs/go-ipfs}

pushTag () {
  local IMAGE_TAG=$1
  if [ "$DRY_RUN" != false ]; then
    echo "DRY RUN! I would have tagged and pushed the following..."
    echo docker tag "$IMAGE_NAME:$WIP_IMAGE_TAG" "$IMAGE_NAME:$IMAGE_TAG"
<<<<<<< HEAD
    echo docker push "$IMAGE_NAME:$IMAGE_TAG"
  else
    echo "Tagging $IMAGE_NAME:$IMAGE_TAG and pushing to dockerhub"
=======
    echo docker push "$IMAGE_NAME:$IMAGE_TAG"  
  else 
>>>>>>> feat: docker build and tag from ci
    docker tag "$IMAGE_NAME:$WIP_IMAGE_TAG" "$IMAGE_NAME:$IMAGE_TAG"
    docker push "$IMAGE_NAME:$IMAGE_TAG"
  fi
}

<<<<<<< HEAD
if [[ $GIT_TAG =~ ^v[0-9]+ ]]; then
  pushTag "$GIT_TAG"
  pushTag "latest"

elif [ "$GIT_BRANCH" = "feat/stabilize-dht" ]; then
  pushTag "bifrost-${BUILD_NUM}-${GIT_SHA1_SHORT}"
  pushTag "bifrost-latest"

elif [ "$GIT_BRANCH" = "master" ]; then
  pushTag "master-${BUILD_NUM}-${GIT_SHA1_SHORT}"
  pushTag "master-latest"
=======
if [[ $GIT_TAG =~ ^v[0-9]+ ]]; then 
  pushTag "$GIT_TAG"

elif [[ $GIT_BRANCH =~ ^cluster ]]; then 
  pushTag "$GIT_BRANCH"

elif [ "$GIT_BRANCH" = "feat/stabilize-dht" ]; then 
  pushTag "bifrost-${BUILD_NUM}-${GIT_SHA1_SHORT}"

elif [ "$GIT_BRANCH" = "release" ]; then 
  pushTag "release"
  pushTag "latest"

elif [ "$GIT_BRANCH" = "master" ]; then 
  pushTag "master"
>>>>>>> feat: docker build and tag from ci

else
  echo "Nothing to do. No docker tag defined for branch: $GIT_BRANCH, tag: $GIT_TAG"

fi
