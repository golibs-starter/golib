package example

// ==================================================
// ========= Example for register a informer ========
// ==================================================

import "gitlab.id.vin/vincart/golib/actuator"

// NewSampleInformer
// In this example, the /actuator/info endpoint with return:
//{
//  "meta": {
//    "code": 200,
//    "message": "Successful"
//  },
//  "data": {
//    "service_name": "Sample Service",
//    "info": {
//      "sample": {
//        "key1": "val1"
//      }
//    }
//  }
//}
func NewSampleInformer() actuator.Informer {
	return &SampleInformer{}
}

type SampleInformer struct {
}

func (s SampleInformer) Key() string {
	return "sample"
}

func (s SampleInformer) Value() interface{} {
	return map[string]interface{}{
		"key1": "val1",
	}
}
