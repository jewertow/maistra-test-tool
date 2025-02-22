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

package authorization

import (
	"testing"

	"github.com/maistra/maistra-test-tool/pkg/app"
	"github.com/maistra/maistra-test-tool/pkg/util/hack"
	"github.com/maistra/maistra-test-tool/pkg/util/oc"
	"github.com/maistra/maistra-test-tool/pkg/util/test"
)

func TestAuthorizationDenyAllow(t *testing.T) {
	test.NewTest(t).Id("T23").Groups(test.Full, test.InterOp).Run(func(t test.TestHelper) {
		hack.DisableLogrusForThisTest(t)

		ns := "foo"
		t.Cleanup(func() {
			oc.RecreateNamespace(t, ns)
		})

		t.Log("This test validates authorization policies with a deny action")

		t.LogStep("Install httpbin and sleep")
		app.InstallAndWaitReady(t, app.Httpbin(ns), app.Sleep(ns))

		t.LogStep("Check if httpbin returns 200 OK when no authorization policies are in place")
		assertHttpbinRequestSucceeds(t, ns, httpbinRequest("GET", "/ip"))

		t.NewSubTest("explicitly deny request").Run(func(t test.TestHelper) {
			t.Cleanup(func() {
				oc.DeleteFromString(t, ns, DenyGETPolicy)
			})
			t.LogStep("Apply policy that denies all GET requests to httpbin")
			oc.ApplyString(t, ns, DenyGETPolicy)

			t.LogStep("Verify that GET request is denied")
			assertRequestDenied(t, ns, httpbinRequest("GET", "/get"), "403")

			t.LogStep("Verify that POST request is allowed")
			assertRequestAccepted(t, ns, httpbinRequest("POST", "/post"))
		})

		t.NewSubTest("deny request header").Run(func(t test.TestHelper) {
			t.Cleanup(func() {
				oc.DeleteFromString(t, ns, DenyHeaderNotAdminPolicy)
			})
			t.LogStep("Apply policy that denies GET requests unless the HTTP header 'x-token: admin' is present")
			oc.ApplyString(t, ns, DenyHeaderNotAdminPolicy)

			t.LogStep("Verify that GET request with HTTP header 'x-token: admin' is allowed")
			assertRequestAccepted(t, ns, httpbinRequest("GET", "/get", "x-token: admin"))

			t.LogStep("Verify that GET request with HTTP header 'x-token: guest' is denied")
			assertRequestDenied(t, ns, httpbinRequest("GET", "/get", "x-token: guest"), "403")
		})

		t.NewSubTest("allow request path").Run(func(t test.TestHelper) {
			t.Cleanup(func() {
				oc.DeleteFromString(t, ns, DenyHeaderNotAdminPolicy)
				oc.DeleteFromString(t, ns, AllowPathIPPolicy)
			})
			t.LogStep("Apply policy that denies GET requests unless the HTTP header 'x-token: admin' is present")
			oc.ApplyString(t, ns, DenyHeaderNotAdminPolicy)

			t.LogStep("Apply policy that allows requests with the path '/ip'")
			oc.ApplyString(t, ns, AllowPathIPPolicy)

			t.LogStep("Verify that GET request with the HTTP header 'x-token: guest' at path '/ip' is denied")
			assertRequestDenied(t, ns, httpbinRequest("GET", "/ip", "x-token: guest"), "403")

			t.LogStep("Verify that GET request with HTTP header 'x-token: admin' at path '/ip' is allowed")
			assertRequestAccepted(t, ns, httpbinRequest("GET", "/ip", "x-token: admin"))

			t.LogStep("Verify that GET request with HTTP header 'x-token: admin' at path '/get' is denied")
			assertRequestDenied(t, ns, httpbinRequest("GET", "/get", "x-token: admin"), "403")
		})
	})
}

const (
	DenyGETPolicy = `
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: deny-method-get
  namespace: foo
spec:
  selector:
    matchLabels:
      app: httpbin
  action: DENY
  rules:
  - to:
    - operation:
        methods: ["GET"]`

	DenyHeaderNotAdminPolicy = `
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: deny-method-get
  namespace: foo
spec:
  selector:
    matchLabels:
      app: httpbin
  action: DENY
  rules:
  - to:
    - operation:
        methods: ["GET"]
    when:
    - key: request.headers[x-token]
      notValues: ["admin"]`

	AllowPathIPPolicy = `
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: allow-path-ip
  namespace: foo
spec:
  selector:
    matchLabels:
      app: httpbin
  action: ALLOW
  rules:
  - to:
    - operation:
        paths: ["/ip"]`
)
