// Copyright The Shipwright Contributors
//
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"fmt"

	buildv1alpha1 "github.com/shipwright-io/build/pkg/apis/build/v1alpha1"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

type sourceResult struct {
	defaultSource,
	bundleSource buildv1alpha1.SourceResult
}

const (
	defaultSourceName       = "default"
	commitSHAResult         = "commit-sha"
	commitAuthorResult      = "commit-author"
	bundleImageDigestResult = "bundle-image-digest"
	imageDigestResult       = "image-digest"
	imageSizeResult         = "image-size"
)

// UpdateBuildRunUsingTaskResults surface the task results
// to the buildrun
func UpdateBuildRunUsingTaskResults(
	buildRun *buildv1alpha1.BuildRun,
	lastTaskRun *v1beta1.TaskRun,
) {
	var sources sourceResult

	// Initializing source result
	sources.defaultSource.Git = &buildv1alpha1.GitSourceResult{}
	sources.bundleSource.Bundle = &buildv1alpha1.BundleSourceResult{}

	// Initializing output result
	buildRun.Status.Output = &buildv1alpha1.Output{}

	for _, result := range lastTaskRun.Status.TaskRunResults {
		updateBuildRunStatus(buildRun, result, &sources)
	}

	// Appending the source result only if the defined source
	// from build spec emitting the results
	if sources.defaultSource.Name != "" {
		buildRun.Status.Sources = append(buildRun.Status.Sources, sources.defaultSource)
	}

	if sources.bundleSource.Name != "" {
		buildRun.Status.Sources = append(buildRun.Status.Sources, sources.bundleSource)
	}
}

func updateBuildRunStatus(
	buildRun *buildv1alpha1.BuildRun,
	result v1beta1.TaskRunResult,
	sources *sourceResult,
) {
	switch result.Name {
	case generateSourceResultName(defaultSourceName, commitSHAResult):
		// Source name is default as `spec.source` has no name field
		sources.defaultSource.Name = defaultSourceName
		sources.defaultSource.Git.CommitSha = result.Value
	case generateSourceResultName(defaultSourceName, commitAuthorResult):
		// Source name is default as `spec.source` has no name field
		sources.defaultSource.Name = defaultSourceName
		sources.defaultSource.Git.CommitAuthor = result.Value
	case generateSourceResultName(defaultSourceName, bundleImageDigestResult):
		// Source name is default as `spec.source` has no name field
		sources.bundleSource.Name = defaultSourceName
		sources.bundleSource.Bundle.Digest = result.Value
	case generateOutputResultName(imageDigestResult):
		buildRun.Status.Output.Digest = result.Value
	case generateOutputResultName(imageSizeResult):
		buildRun.Status.Output.Size = result.Value
	}
}

func generateSourceResultName(source string, resultName string) string {
	return fmt.Sprintf("%s-source-%s-%s", prefixParamsResultsVolumes, defaultSourceName, resultName)
}

func generateOutputResultName(resultName string) string {
	return fmt.Sprintf("%s-%s", prefixParamsResultsVolumes, resultName)
}
