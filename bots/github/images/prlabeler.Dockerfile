# Copyright 2025 The Kubernetes Authors.
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
FROM golang:1.24.2 AS builder

WORKDIR /app
COPY ../ .
RUN CGO_ENABLED=0 go build -o _output/bin/prlabeler ./cmd/prlabeler

FROM scratch

MAINTAINER Johhny Cottage <???@???.??>

LABEL org.opencontainers.image.source https://github.com/ingvagabund/gearhouse

USER 1000

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /app/_output/bin/prlabeler /bin/prlabeler

CMD ["/bin/prlabeler", "--help"]
