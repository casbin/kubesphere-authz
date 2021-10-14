package audit

import (
	//"sync"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

type AuditItem struct {
	Time        time.Time
	Namespace   string
	Resource    string
	RequestBody string
	Verb        string
	Result      bool
	Info        string
}

func (a *AuditItem) String() string {
	resString := "APPROVED"
	if !a.Result {
		resString = "REJECTED"
	}
	res := fmt.Sprintf("%s %s %s namespace: %s, resource:%s ,", a.Time, resString, a.Verb, a.Namespace, a.Resource)
	if a.Info != "" {
		res += "info: " + a.Info + ","
	}
	if a.RequestBody != "" {
		res += "requestBody: %s" + a.RequestBody
	}
	return res
}

type Auditor struct {
	//sync.Mutex
	items           []AuditItem
	bufferSize      int
	recv            chan AuditItem
	stop            chan struct{}
	showRequestBody bool
	logPath         string
}

func NewAuditor(bufferSize int, showRequestBody bool, logPath string) *Auditor {
	if bufferSize <= 0 {
		bufferSize = 1000
	}
	var res = Auditor{
		items:           make([]AuditItem, 0, bufferSize),
		bufferSize:      bufferSize,
		recv:            make(chan AuditItem),
		showRequestBody: showRequestBody,
		logPath:         logPath,
	}
	return &res
}

func (a *Auditor) Stop() {
	a.stop <- struct{}{}
}

func (a *Auditor) Insert(data []byte, approved bool, err error) {
	var requestBody v1.AdmissionReview
	var decoder runtime.Decoder = serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer()
	decoder.Decode(data, nil, &requestBody)

	var item = AuditItem{
		Time:      time.Now(),
		Namespace: requestBody.Request.Namespace,
		Resource:  requestBody.Request.Resource.String(),
		Verb:      string(requestBody.Request.Operation),
		Result:    approved,
	}
	if err != nil {
		item.Info = err.Error()
	}
	if a.showRequestBody {
		item.RequestBody = string(data)
	}
	a.recv <- item
}

func (a *Auditor) Run() {
	defer a.cleanUp()
	for {
		select {
		case <-a.stop:
			return
		case item := <-a.recv:
			a.items = append(a.items, item)
			if len(a.items) >= a.bufferSize {
				a.WriteIntoLog()
			}
		}
	}

}

func (a *Auditor) cleanUp() {
	if len(a.items) != 0 {
		a.WriteIntoLog()
	}
}

func (a *Auditor) WriteIntoLog() {
	currentTime := time.Now().Unix()
	fileName := "audit_" + strconv.FormatInt(currentTime, 10) + ".log"
	fullFileName := filepath.Join(a.logPath, fileName)
	file, err := os.OpenFile(fullFileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("Auditor: failed to open %s: %s", fullFileName, err.Error())
		return
	}
	for _, item := range a.items {
		io.WriteString(file, item.String()+"\n")
	}
	a.items = a.items[0:0]
	file.Close()
}
