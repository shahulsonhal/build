#!/bin/bash

# Copyright The Shipwright Contributors
#
# SPDX-License-Identifier: Apache-2.0

# Verifies if a developer has forgot to run the
# `make generate-crds` so that all the changes in the
# CRD (deploy/crds) should also be pushed

if [[ -n "$(git status --porcelain -- deploy/crds)" ]]; then
 echo "The deploy/crds contains changes:"
  git --no-pager diff --name-only -- deploy/crds
  echo
  echo "Run make generate-crds to those commit changes!"
  exit 1
fi
