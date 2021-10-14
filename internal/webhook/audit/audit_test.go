package audit

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	. "github.com/agiledragon/gomonkey"
	. "github.com/smartystreets/goconvey/convey"
)

const testData1 string = `
{"kind":"AdmissionReview","apiVersion":"admission.k8s.io/v1","request":{"uid":"86ce1fec-24c6-4c0f-b5b4-705e62513797","kind":{"group":"","version":"v1","kind":"Endpoints"},"resource":{"group":"","version":"v1","resource":"endpoints"},"requestKind":{"group":"","version":"v1","kind":"Endpoints"},"requestResource":{"group":"","version":"v1","resource":"endpoints"},"name":"k8s.io-minikube-hostpath","namespace":"kube-system","operation":"UPDATE","userInfo":{"username":"system:serviceaccount:kube-system:storage-provisioner","uid":"f76978a7-db88-44fe-bf8e-1d7c903d04e7","groups":["system:serviceaccounts","system:serviceaccounts:kube-system","system:authenticated"],"extra":{"authentication.kubernetes.io/pod-name":["storage-provisioner"],"authentication.kubernetes.io/pod-uid":["cea5cc69-3d0c-4d97-98bb-ea58ffc1d912"]}},"object":{"kind":"Endpoints","apiVersion":"v1","metadata":{"name":"k8s.io-minikube-hostpath","namespace":"kube-system","uid":"493c5775-dba4-4a32-ae6e-51dee1d26dec","resourceVersion":"796","creationTimestamp":"2021-10-12T10:16:47Z","annotations":{"control-plane.alpha.kubernetes.io/leader":"{\"holderIdentity\":\"minikube_60a2b6d4-7f58-4c67-aadb-c2fda56b01bb\",\"leaseDurationSeconds\":15,\"acquireTime\":\"2021-10-12T10:16:47Z\",\"renewTime\":\"2021-10-12T10:24:22Z\",\"leaderTransitions\":0}"},"managedFields":[{"manager":"storage-provisioner","operation":"Update","apiVersion":"v1","time":"2021-10-12T10:16:47Z","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:control-plane.alpha.kubernetes.io/leader":{}}}}}]}},"oldObject":{"kind":"Endpoints","apiVersion":"v1","metadata":{"name":"k8s.io-minikube-hostpath","namespace":"kube-system","uid":"493c5775-dba4-4a32-ae6e-51dee1d26dec","resourceVersion":"796","creationTimestamp":"2021-10-12T10:16:47Z","annotations":{"control-plane.alpha.kubernetes.io/leader":"{\"holderIdentity\":\"minikube_60a2b6d4-7f58-4c67-aadb-c2fda56b01bb\",\"leaseDurationSeconds\":15,\"acquireTime\":\"2021-10-12T10:16:47Z\",\"renewTime\":\"2021-10-12T10:24:20Z\",\"leaderTransitions\":0}"}}},"dryRun":false,"options":{"kind":"UpdateOptions","apiVersion":"meta.k8s.io/v1"}}}

`
const testData2 string = `
{"kind":"AdmissionReview","apiVersion":"admission.k8s.io/v1","request":{"uid":"ad6a5637-8d86-49f2-b8c8-dd24a1c638c8","kind":{"group":"coordination.k8s.io","version":"v1","kind":"Lease"},"resource":{"group":"coordination.k8s.io","version":"v1","resource":"leases"},"requestKind":{"group":"coordination.k8s.io","version":"v1","kind":"Lease"},"requestResource":{"group":"coordination.k8s.io","version":"v1","resource":"leases"},"name":"minikube","namespace":"kube-node-lease","operation":"UPDATE","userInfo":{"username":"system:node:minikube","groups":["system:nodes","system:authenticated"]},"object":{"kind":"Lease","apiVersion":"coordination.k8s.io/v1","metadata":{"name":"minikube","namespace":"kube-node-lease","uid":"e333edd6-3349-42ea-9bcd-f1e2ab945f72","resourceVersion":"789","creationTimestamp":"2021-10-12T10:16:31Z","ownerReferences":[{"apiVersion":"v1","kind":"Node","name":"minikube","uid":"f2b3bd6f-463c-4513-a967-c290fbd11014"}],"managedFields":[{"manager":"kubelet","operation":"Update","apiVersion":"coordination.k8s.io/v1","time":"2021-10-12T10:16:31Z","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:ownerReferences":{".":{},"k:{\"uid\":\"f2b3bd6f-463c-4513-a967-c290fbd11014\"}":{".":{},"f:apiVersion":{},"f:kind":{},"f:name":{},"f:uid":{}}}},"f:spec":{"f:holderIdentity":{},"f:leaseDurationSeconds":{},"f:renewTime":{}}}}]},"spec":{"holderIdentity":"minikube","leaseDurationSeconds":40,"renewTime":"2021-10-12T10:24:21.171624Z"}},"oldObject":{"kind":"Lease","apiVersion":"coordination.k8s.io/v1","metadata":{"name":"minikube","namespace":"kube-node-lease","uid":"e333edd6-3349-42ea-9bcd-f1e2ab945f72","resourceVersion":"789","creationTimestamp":"2021-10-12T10:16:31Z","ownerReferences":[{"apiVersion":"v1","kind":"Node","name":"minikube","uid":"f2b3bd6f-463c-4513-a967-c290fbd11014"}]},"spec":{"holderIdentity":"minikube","leaseDurationSeconds":40,"renewTime":"2021-10-12T10:24:10.849561Z"}},"dryRun":false,"options":{"kind":"UpdateOptions","apiVersion":"meta.k8s.io/v1"}}}
`

func TestInsert(t *testing.T) {
	Convey("testInsert", t, func() {
		auditor := NewAuditor(2, false, ".")
		go auditor.Run()
		auditor.Insert([]byte(testData1), true, nil)
		//sleep to avoid data race

		time.Sleep(1 * time.Second)
		So(len(auditor.items), ShouldEqual, 1)
		tmp := auditor.items[0]
		So(tmp.Namespace, ShouldEqual, "kube-system")
		So(tmp.Resource, ShouldEqual, "/v1, Resource=endpoints")
		So(tmp.Verb, ShouldEqual, "UPDATE")
		var buffer bytes.Buffer
		patches := ApplyMethod(reflect.TypeOf(auditor), "WriteIntoLog", func(a *Auditor) {

			for _, item := range a.items {
				io.WriteString(&buffer, item.String()+"\n")
			}
			a.items = a.items[0:0]
		})
		defer patches.Reset()

		auditor.Insert([]byte(testData2), false, nil)
		//sleep to avoid data race
		time.Sleep(3 * time.Second)
		So(len(auditor.items), ShouldEqual, 0)
		out := buffer.String()
		lines := strings.Split(out, "\n")
		So(len(lines), ShouldEqual, 3)

	})

}
func TestMultipleInsert(t *testing.T) {
	Convey("testMultipleInsert", t, func() {
		auditor := NewAuditor(300, false, ".")
		go auditor.Run()

		var buffer bytes.Buffer
		patches := ApplyMethod(reflect.TypeOf(auditor), "WriteIntoLog", func(a *Auditor) {

			for _, item := range a.items {
				io.WriteString(&buffer, item.String()+"\n")
			}
			a.items = a.items[0:0]
		})
		defer patches.Reset()
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				auditor.Insert([]byte(testData1), true, nil)
				auditor.Insert([]byte(testData2), false, nil)
				wg.Done()
			}()

		}

		wg.Wait()
		//sleep to avoid data race
		time.Sleep(3 * time.Second)
		So(len(auditor.items), ShouldEqual, 200)
	})
}
