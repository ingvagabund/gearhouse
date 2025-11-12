# Copyright 2025 The Gearhouse Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

IMAGE_VERSION=0.6

.PHONY: build
build:
	CGO_ENABLED=0 go build -o _output/bin/prlabeler ./cmd/prlabeler

image:
	podman build -t quay.io/jchaloup/prlabeler:$(IMAGE_VERSION) -f images/prlabeler.Dockerfile .

push:
	podman push quay.io/jchaloup/prlabeler:$(IMAGE_VERSION)

.PHONY: clean
clean:
	rm -rf _output
