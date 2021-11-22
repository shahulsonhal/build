#!/bin/bash

# Copyright The Shipwright Contributors
#
# SPDX-License-Identifier: Apache-2.0

# Verifies if a developer has forgot to run the
# `make generate` or `make generate-crd` so that all the changes in the
# clientset and CRD (deploy/crds) should also be pushed

if [[ -n "$(git status --porcelain -- pkg/client pkg/apis)" ]]; then
  echo "The pkg/client or pkg/apis package contains changes:"
  git --no-pager diff --name-only -- pkg/client pkg/apis
  echo
  echo "Run make generate to those commit changes!"
  exit 1
fi

if [[ -n "$(git status --porcelain -- deploy/crds)" ]]; then
 echo "The deploy/crds contains changes:"
  git --no-pager diff --name-only -- deploy/crds
  echo
  echo "Run make generate-crds to those commit changes!"
  exit 1
fi
