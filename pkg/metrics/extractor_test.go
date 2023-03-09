package metrics

import (
	"testing"
)

func TestToSnakeCase(t *testing.T) {
	v := toSnakeCase("StatusSNS")
	if v != "status_sns" {
		t.Fatalf("wrong value: %s", v)
	}
}

func TestExtract(t *testing.T) {
	json := `
{
  "Status": {
    "DeviceName": "nas_outlet",
    "DeviceName-alias": "outlet",
    "Topic": "nas-outlet",
    "SaveData": 1
  },
  "StatusSTS": {
    "UptimeSec": 1923,
    "Heap": 22,
    "Wifi": {
      "AP": 1,
      "SSId": "homenet_2G",
      "Channel": 8
    }
  }
}
`

	extracted := Extract([]byte(json))

	if extracted["status_sts_wifi_ap"] != 1. {
		t.Fatalf("wrong value: %s", extracted["status_sts_wifi_ap"])
	}

	if extracted["status_device_name"] != "nas_outlet" {
		t.Fatalf("wrong value: %s", extracted["status_device_name"])
	}

	if extracted["status_device_name_alias"] != "outlet" {
		t.Fatalf("wrong value: %s", extracted["status_device_name_alias"])
	}

}
