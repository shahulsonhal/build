// Copyright The Shipwright Contributors
//
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"

	core "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/pointer"

	buildv1alpha1 "github.com/shipwright-io/build/pkg/apis/build/v1alpha1"
)

const (
	pollCreateInterval = 1 * time.Second
	pollCreateTimeout  = 10 * time.Second
)

type buildPrototype struct{ build buildv1alpha1.Build }
type buildRunPrototype struct{ buildRun buildv1alpha1.BuildRun }

func NewBuildPrototype() *buildPrototype {
	return &buildPrototype{
		build: buildv1alpha1.Build{},
	}
}

func (b *buildPrototype) Name(name string) *buildPrototype {
	b.build.ObjectMeta.Name = name
	return b
}

func (b *buildPrototype) Namespace(namespace string) *buildPrototype {
	b.build.ObjectMeta.Namespace = namespace
	return b
}

func (b *buildPrototype) ClusterBuildStrategy(name string) *buildPrototype {
	var cbs = buildv1alpha1.ClusterBuildStrategyKind
	b.build.Spec.Strategy = buildv1alpha1.Strategy{
		Kind: &cbs,
		Name: name,
	}
	return b
}

func (b *buildPrototype) SourceBundle(image string) *buildPrototype {
	if b.build.Spec.Source.BundleContainer == nil {
		b.build.Spec.Source.BundleContainer = &buildv1alpha1.BundleContainer{}
	}
	b.build.Spec.Source.BundleContainer.Image = image
	return b
}

func (b *buildPrototype) SourceContextDir(contextDir string) *buildPrototype {
	b.build.Spec.Source.ContextDir = pointer.StringPtr(contextDir)
	return b
}

func (b *buildPrototype) Dockerfile(dockerfile string) *buildPrototype {
	b.build.Spec.Dockerfile = &dockerfile
	return b
}

func (b *buildPrototype) OutputImage(image string) *buildPrototype {
	b.build.Spec.Output.Image = image
	return b
}

func (b *buildPrototype) OutputImageCredentials(name string) *buildPrototype {
	if name != "" {
		b.build.Spec.Output.Credentials = &core.LocalObjectReference{Name: name}
	}

	return b
}

func (b buildPrototype) Create() (build *buildv1alpha1.Build, err error) {
	ctx := context.Background()

	_, err = testBuild.
		BuildClientSet.
		ShipwrightV1alpha1().
		Builds(b.build.Namespace).
		Create(ctx, &b.build, meta.CreateOptions{})

	if err != nil {
		return nil, err
	}

	err = wait.PollImmediate(pollCreateInterval, pollCreateTimeout, func() (done bool, err error) {
		build, err = testBuild.BuildClientSet.ShipwrightV1alpha1().Builds(b.build.Namespace).Get(ctx, b.build.Name, meta.GetOptions{})
		if err != nil {
			return false, err
		}

		return build.Status.Registered != nil && *build.Status.Registered == v1.ConditionTrue, nil
	})

	return
}

func NewBuildRunPrototype() *buildRunPrototype {
	return &buildRunPrototype{buildRun: buildv1alpha1.BuildRun{}}
}

func (b *buildRunPrototype) Name(name string) *buildRunPrototype {
	b.buildRun.ObjectMeta.Name = name
	return b
}

func (b *buildRunPrototype) ForBuild(build *buildv1alpha1.Build) *buildRunPrototype {
	b.buildRun.Spec.BuildRef = buildv1alpha1.BuildRef{Name: build.Name}
	b.buildRun.ObjectMeta.Namespace = build.Namespace
	return b
}

func (b *buildRunPrototype) GenerateServiceAccount() *buildRunPrototype {
	if b.buildRun.Spec.ServiceAccount == nil {
		b.buildRun.Spec.ServiceAccount = &buildv1alpha1.ServiceAccount{}
	}
	b.buildRun.Spec.ServiceAccount.Generate = pointer.BoolPtr(true)
	return b
}

func (b *buildRunPrototype) Create() (*buildv1alpha1.BuildRun, error) {
	return testBuild.
		BuildClientSet.
		ShipwrightV1alpha1().
		BuildRuns(b.buildRun.Namespace).
		Create(context.Background(), &b.buildRun, meta.CreateOptions{})
}

// Logf logs data
func Logf(format string, args ...interface{}) {
	currentTime := time.Now().UTC().Format(time.RFC3339)

	fmt.Fprintf(
		GinkgoWriter,
		fmt.Sprintf("%s %d %s\n", currentTime, getGinkgoNode(), format),
		args...,
	)
}

func getArg(argName string) (bool, string) {
	for i, arg := range os.Args {
		if arg == argName {
			return true, os.Args[i+1]
		} else if strings.HasPrefix(arg, argName+"=") {
			argAndValue := strings.SplitN(arg, "=", 2)
			return true, argAndValue[1]
		}
	}
	return false, ""
}

func getGinkgoNode() int {
	defined, ginkgoNodeString := getArg("--ginkgo.parallel.node")
	if !defined {
		return 1
	}
	ginkgoNode, err := strconv.Atoi(ginkgoNodeString)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return 0
	}
	return ginkgoNode
}
