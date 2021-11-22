#!/bin/bash

# Copyright The Shipwright Contributors
#
# SPDX-License-Identifier: Apache-2.0

# Verifies if a developer has forgot to run the
# `make generate` so that all the changes in the
# clientset should also be pushed

if [[ -n "$(git status --porcelain -- pkg/client pkg/apis)" ]]; then
  echo "The pkg/client or pkg/apis package contains changes:"
  git --no-pager diff --name-only -- pkg/client pkg/apis
  echo
  echo "Run make generate to those commit changes!"
  exit 1
fi
