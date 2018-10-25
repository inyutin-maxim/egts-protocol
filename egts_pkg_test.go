package main

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

/*
Packet data:
 0100030B0023008A0001491800610099B00902000202101500D53F01106F1C059E7AB53C3501D0872C0100000000CC27

EGTS Transport Layer:
---------------------
  Validating result   - 0 (OK)

  Protocol Version    - 1
  Security Key ID     - 0
  Flags               - 00000011b (0x03)
       Prefix         - 00
       Route          -   0
       Encryption Alg -    00
       Compression    -      0
       Priority       -       11 (low)
  Header Length       - 11
  Header Encoding     - 0
  Frame Data Length   - 35
  Packet ID           - 138
  No route info       -
  Header Check Sum    - 0x49

EGTS Service Layer:
---------------------
  Validating result   - 0 (OK)

  Packet Type         - EGTS_PT_APPDATA
  Service Layer CS    - 0x27CC

    Service Layer Record:
    ---------------------
    Validating Result    - 0 (OK)

    Record Length               - 24
    Record Number               - 97
    Record flags                -     10011001b (0x99)
        Sourse Service On Device    - 1
        Recipient Service On Device -  0
        Group Flag                  -   0
        Record Processing Priority  -    11 (low)
        Time Field Exists           -      0
        Event ID Field Exists       -       0
        Object ID Field Exists      -        1
    Object Identifier           - 133552
    Source Service Type         - 2 (EGTS_TELEDATA_SERVICE) from ST
    Recipient Service Type      - 2 (EGTS_TELEDATA_SERVICE)

       Subrecord Data:
       ------------------
       Validating Result   - 150 (Unknown service)

       Subrecord Type      - 16 (unspecified)
       Subrecord Length    - 21
*/

var (
	egtsPkgValid = EgtsPackage{
		ProtocolVersion:  1,
		SecurityKeyID:    0,
		Prefix:           "00",
		Route:            "0",
		EncryptionAlg:    "00",
		Compression:      "0",
		Priority:         "11",
		HeaderLength:     11,
		HeaderEncoding:   0,
		FrameDataLength:  35,
		PacketIdentifier: 138,
		PacketType:       1,
		HeaderCheckSum:   73,
		ServicesFrameData: &ServiceDataSet{
			ServiceDataRecord{
				RecordLength:             24,
				RecordNumber:             97,
				SourceServiceOnDevice:    "1",
				RecipientServiceOnDevice: "0",
				Group: "0",
				RecordProcessingPriority: "11",
				TimeFieldExists:          "0",
				EventIDFieldExists:       "0",
				ObjectIDFieldExists:      "1",
				ObjectIdentifier:         133552,
				SourceServiceType:        2,
				RecipientServiceType:     2,
				RecordDataSet: RecordDataSet{
					RecordData{
						SubrecordType:   16,
						SubrecordLength: 21,
						SubrecordData: &EgtsSrPosData{
							NavigationTime:      time.Date(2018, time.July, 4, 20, 8, 53, 0, time.UTC),
							Latitude:            55,
							Longitude:           37,
							ALTE:                "0",
							LOHS:                "0",
							LAHS:                "0",
							MV:                  "0",
							BB:                  "0",
							CS:                  "0",
							FIX:                 "0",
							VLD:                 "1",
							DirectionHighestBit: 1,
							AltitudeSign:        0,
							Speed:               200,
							Direction:           44,
							Odometer:            []byte{0x01, 0x00, 0x00},
							DigitalInputs:       0,
							Source:              0,
						},
					},
				},
			},
		},
		ServicesFrameDataCheckSum: 10188, //52263
	}
)

func TestEgtsPackage_Encode(t *testing.T) {
	testEgtsPkgBytes := []byte{0x01, 0x00, 0x03, 0x0B, 0x00, 0x23, 0x00, 0x8A, 0x00, 0x01, 0x49, 0x18, 0x00, 0x61,
		0x00, 0x99, 0xB0, 0x09, 0x02, 0x00, 0x02, 0x02, 0x10, 0x15, 0x00, 0xD5, 0x3F, 0x01, 0x10, 0x1b, 0xc7, 0x71, 0x9c,
		0xf4, 0x49, 0x9f, 0x34, 0x01, 0xD0, 0x87, 0x2C, 0x01, 0x00, 0x00, 0x00, 0x00, 0xAC, 0xC9}

	posDataBytes, err := egtsPkgValid.Encode()
	if err != nil {
		t.Errorf("Ошибка кодирования: %v\n", err)
	}

	if !bytes.Equal(posDataBytes, testEgtsPkgBytes) {
		t.Errorf("Байтовые строки не совпадают: %v != %v ", posDataBytes, testEgtsPkgBytes)
	}
}


func TestEgtsPackage_Decode(t *testing.T) {
	egtsPkgBytes := []byte{0x01, 0x00, 0x03, 0x0B, 0x00, 0x23, 0x00, 0x8A, 0x00, 0x01, 0x49, 0x18, 0x00, 0x61,
		0x00, 0x99, 0xB0, 0x09, 0x02, 0x00, 0x02, 0x02, 0x10, 0x15, 0x00, 0xD5, 0x3F, 0x01, 0x10, 0x6F, 0x1C, 0x05, 0x9E,
		0x7A, 0xB5, 0x3C, 0x35, 0x01, 0xD0, 0x87, 0x2C, 0x01, 0x00, 0x00, 0x00, 0x00, 0xCC, 0x27}

	egtsPkg := EgtsPackage{}

	if _, err := egtsPkg.Decode(egtsPkgBytes); err != nil {
		t.Errorf("Ошибка декадирования: %v\n", err)
	}

	if !reflect.DeepEqual(egtsPkg, egtsPkgValid) {
		t.Errorf("Запись ServicesFrameData не совпадают: %v != %v ", egtsPkg, egtsPkgValid)
	}
}
