#!/usr/bin/env bash

os=$(go env GOOS)
arch=$(go env GOARCH)

if [[ -d /usr/local/kubebuilder ]]; then
  /usr/local/kubebuilder/bin/kubebuilder version
  exit 0
fi

DEFAULT_RELEASE_VERSION=2.3.0
RELEASE_VERSION=${1:-$DEFAULT_RELEASE_VERSION}

# download kubebuilder and extract it to tmp
curl -L https://go.kubebuilder.io/dl/${RELEASE_VERSION}/${os}/${arch} | tar -xz -C /tmp/

# move to a long-term location and put it on your path
sudo mv /tmp/kubebuilder_${RELEASE_VERSION}_${os}_${arch} /usr/local/kubebuilder
