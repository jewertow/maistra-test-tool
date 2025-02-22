// Copyright 2021 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ingress

import (
	"github.com/maistra/maistra-test-tool/pkg/util"
	"github.com/maistra/maistra-test-tool/pkg/util/env"
)

var (
	httpbinSampleServerCertKey = env.GetRootDir() + "/sampleCerts/httpbin.example.com/httpbin.example.com.key"
	httpbinSampleServerCert    = env.GetRootDir() + "/sampleCerts/httpbin.example.com/httpbin.example.com.crt"
	httpbinSampleCACert        = env.GetRootDir() + "/sampleCerts/httpbin.example.com/example.com.crt"
	httpbinSampleClientCert    = env.GetRootDir() + "/sampleCerts/httpbin.example.com/httpbin-client.example.com.crt"
	httpbinSampleClientCertKey = env.GetRootDir() + "/sampleCerts/httpbin.example.com/httpbin-client.example.com.key"

	helloworldServerCertKey = env.GetRootDir() + "/sampleCerts/helloworldv1/helloworld-v1.example.com.key"
	helloworldServerCert    = env.GetRootDir() + "/sampleCerts/helloworldv1/helloworld-v1.example.com.crt"

	nginxServerCACert = env.GetRootDir() + "/sampleCerts/nginx.example.com/example.com.crt"
)

var (
	// OCP4.x
	meshNamespace        string = env.Getenv("MESHNAMESPACE", "istio-system")
	gatewayHTTP, _              = util.ShellSilent(`kubectl get routes -n %s istio-ingressgateway -o jsonpath='{.spec.host}'`, meshNamespace)
	secureIngressPort, _        = util.ShellSilent(`kubectl -n %s get service %s -o jsonpath='{.spec.ports[?(@.name=="https")].port}'`, meshNamespace, "istio-ingressgateway")
)
