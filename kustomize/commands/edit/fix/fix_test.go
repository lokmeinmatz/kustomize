// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package fix

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	testutils_test "sigs.k8s.io/kustomize/kustomize/v4/commands/internal/testutils"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

func TestFix(t *testing.T) {
	fSys := filesys.MakeFsInMemory()
	testutils_test.WriteTestKustomizationWith(fSys, []byte(`nameprefix: some-prefix-`))

	cmd := NewCmdFix(fSys, os.Stdout)
	assert.NoError(t, cmd.RunE(cmd, nil))

	content, err := testutils_test.ReadTestKustomization(fSys)
	assert.NoError(t, err)

	assert.Contains(t, string(content), "apiVersion: ")
	assert.Contains(t, string(content), "kind: Kustomization")
}

func TestFixOutdatedPatchesFieldTitle(t *testing.T) {
	kustomizationContentWithOutdatedPatchesFieldTitle := []byte(`
patchesJson6902:
- path: patch1.yaml
  target:
    kind: Service
- path: patch2.yaml
  target:
    group: apps
    kind: Deployment
    version: v1
`)

	expected := []byte(`
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
patches:
- path: patch1.yaml
  target:
    kind: Service
- path: patch2.yaml
  target:
    group: apps
    kind: Deployment
    version: v1
`)
	fSys := filesys.MakeFsInMemory()
	testutils_test.WriteTestKustomizationWith(fSys, kustomizationContentWithOutdatedPatchesFieldTitle)
	cmd := NewCmdFix(fSys, os.Stdout)
	assert.NoError(t, cmd.RunE(cmd, nil))

	content, err := testutils_test.ReadTestKustomization(fSys)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "apiVersion: ")
	assert.Contains(t, string(content), "kind: Kustomization")

	if diff := cmp.Diff(expected, content); diff != "" {
		t.Errorf("Mismatch (-expected, +actual):\n%s", diff)
	}
}

func TestRenameAndKeepOutdatedPatchesField(t *testing.T) {
	kustomizationContentWithOutdatedPatchesFieldTitle := []byte(`
patchesJson6902:
- path: patch1.yaml
  target:
    kind: Deployment
patches:
- path: patch2.yaml
  target:
    kind: Deployment
- path: patch3.yaml
  target:
    kind: Service
`)

	expected := []byte(`
patches:
- path: patch2.yaml
  target:
    kind: Deployment
- path: patch3.yaml
  target:
    kind: Service
- path: patch1.yaml
  target:
    kind: Deployment
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
`)
	fSys := filesys.MakeFsInMemory()
	testutils_test.WriteTestKustomizationWith(fSys, kustomizationContentWithOutdatedPatchesFieldTitle)
	cmd := NewCmdFix(fSys, os.Stdout)
	assert.NoError(t, cmd.RunE(cmd, nil))

	content, err := testutils_test.ReadTestKustomization(fSys)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "apiVersion: ")
	assert.Contains(t, string(content), "kind: Kustomization")

	if diff := cmp.Diff(expected, content); diff != "" {
		t.Errorf("Mismatch (-expected, +actual):\n%s", diff)
	}
}

func TestFixOutdatedPatchesStrategicMergeFieldTitle(t *testing.T) {
	kustomizationContentWithOutdatedPatchesFieldTitle := []byte(`
patchesStrategicMerge:
- |-
  apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: nginx
  spec:
    template:
      spec:
        containers:
          - name: nginx
            image: nignx:latest
`)

	expected := []byte(`
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
patches:
- patch: |-
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: nginx
    spec:
      template:
        spec:
          containers:
            - name: nginx
              image: nignx:latest
`)
	fSys := filesys.MakeFsInMemory()
	testutils_test.WriteTestKustomizationWith(fSys, kustomizationContentWithOutdatedPatchesFieldTitle)
	cmd := NewCmdFix(fSys, os.Stdout)
	assert.NoError(t, cmd.RunE(cmd, nil))

	content, err := testutils_test.ReadTestKustomization(fSys)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "apiVersion: ")
	assert.Contains(t, string(content), "kind: Kustomization")

	if diff := cmp.Diff(expected, content); diff != "" {
		t.Errorf("Mismatch (-expected, +actual):\n%s", diff)
	}
}

func TestFixAndMergeOutdatedPatchesStrategicMergeFieldTitle(t *testing.T) {
	kustomizationContentWithOutdatedPatchesFieldTitle := []byte(`
patchesStrategicMerge:
- |-
  apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: nginx
  spec:
    template:
      spec:
        containers:
          - name: nginx
            image: nignx:latest
patches:
- path: patch2.yaml
  target:
    kind: Deployment
`)

	expected := []byte(`
patches:
- path: patch2.yaml
  target:
    kind: Deployment
- patch: |-
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: nginx
    spec:
      template:
        spec:
          containers:
            - name: nginx
              image: nignx:latest
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
`)
	fSys := filesys.MakeFsInMemory()
	testutils_test.WriteTestKustomizationWith(fSys, kustomizationContentWithOutdatedPatchesFieldTitle)
	cmd := NewCmdFix(fSys, os.Stdout)
	assert.NoError(t, cmd.RunE(cmd, nil))

	content, err := testutils_test.ReadTestKustomization(fSys)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "apiVersion: ")
	assert.Contains(t, string(content), "kind: Kustomization")

	if diff := cmp.Diff(expected, content); diff != "" {
		t.Errorf("Mismatch (-expected, +actual):\n%s", diff)
	}
}

func TestFixOutdatedPatchesStrategicMergeToPathFieldTitle(t *testing.T) {
	kustomizationContentWithOutdatedPatchesFieldTitle := []byte(`
patchesStrategicMerge:
- deploy.yaml
`)

	expected := []byte(`
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
patches:
- path: deploy.yaml
`)
	fSys := filesys.MakeFsInMemory()
	testutils_test.WriteTestKustomizationWith(fSys, kustomizationContentWithOutdatedPatchesFieldTitle)
	cmd := NewCmdFix(fSys, os.Stdout)
	assert.NoError(t, cmd.RunE(cmd, nil))

	content, err := testutils_test.ReadTestKustomization(fSys)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "apiVersion: ")
	assert.Contains(t, string(content), "kind: Kustomization")

	if diff := cmp.Diff(expected, content); diff != "" {
		t.Errorf("Mismatch (-expected, +actual):\n%s", diff)
	}
}

func TestFixOutdatedPatchesStrategicMergeToPathFieldYMLTitle(t *testing.T) {
	kustomizationContentWithOutdatedPatchesFieldTitle := []byte(`
patchesStrategicMerge:
- deploy.yml
`)

	expected := []byte(`
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
patches:
- path: deploy.yml
`)
	fSys := filesys.MakeFsInMemory()
	testutils_test.WriteTestKustomizationWith(fSys, kustomizationContentWithOutdatedPatchesFieldTitle)
	cmd := NewCmdFix(fSys, os.Stdout)
	assert.NoError(t, cmd.RunE(cmd, nil))

	content, err := testutils_test.ReadTestKustomization(fSys)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "apiVersion: ")
	assert.Contains(t, string(content), "kind: Kustomization")

	if diff := cmp.Diff(expected, content); diff != "" {
		t.Errorf("Mismatch (-expected, +actual):\n%s", diff)
	}
}

func TestFixOutdatedPatchesStrategicMergeFieldPatchEndOfYamlTitle(t *testing.T) {
	kustomizationContentWithOutdatedPatchesFieldTitle := []byte(`
patchesStrategicMerge:
- |-
  apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: nginx
  spec:
    template:
      spec:
        containers:
          - name: nginx
            env:
              - name: CONFIG_FILE_PATH
                value: home.yaml
`)

	expected := []byte(`
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
patches:
- patch: |-
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: nginx
    spec:
      template:
        spec:
          containers:
            - name: nginx
              env:
                - name: CONFIG_FILE_PATH
                  value: home.yaml
`)
	fSys := filesys.MakeFsInMemory()
	testutils_test.WriteTestKustomizationWith(fSys, kustomizationContentWithOutdatedPatchesFieldTitle)
	cmd := NewCmdFix(fSys, os.Stdout)
	assert.NoError(t, cmd.RunE(cmd, nil))

	content, err := testutils_test.ReadTestKustomization(fSys)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "apiVersion: ")
	assert.Contains(t, string(content), "kind: Kustomization")

	if diff := cmp.Diff(expected, content); diff != "" {
		t.Errorf("Mismatch (-expected, +actual):\n%s", diff)
	}
}

func TestFixOutdatedCommonLabels(t *testing.T) {
	kustomizationContentWithOutdatedCommonLabels := []byte(`
commonLabels:
  foo: bar
labels:
- pairs:
    a: b
`)

	expected := []byte(`
labels:
- pairs:
    a: b
- includeSelectors: true
  pairs:
    foo: bar
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
`)
	fSys := filesys.MakeFsInMemory()
	testutils_test.WriteTestKustomizationWith(fSys, kustomizationContentWithOutdatedCommonLabels)
	cmd := NewCmdFix(fSys, os.Stdout)
	assert.NoError(t, cmd.RunE(cmd, nil))

	content, err := testutils_test.ReadTestKustomization(fSys)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "apiVersion: ")
	assert.Contains(t, string(content), "kind: Kustomization")

	if diff := cmp.Diff(expected, content); diff != "" {
		t.Errorf("Mismatch (-expected, +actual):\n%s", diff)
	}
}

func TestFixOutdatedCommonLabelsDuplicate(t *testing.T) {
	kustomizationContentWithOutdatedCommonLabels := []byte(`
commonLabels:
  foo: bar
labels:
- pairs:
    foo: baz
    a: b
`)

	fSys := filesys.MakeFsInMemory()
	testutils_test.WriteTestKustomizationWith(fSys, kustomizationContentWithOutdatedCommonLabels)
	cmd := NewCmdFix(fSys, os.Stdout)
	err := cmd.RunE(cmd, nil)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "label name 'foo' exists in both commonLabels and labels")
}
