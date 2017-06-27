package ipfix_test

import (
	"bytes"
	"encoding/hex"
	"io"
	"sync"
	"testing"

	"github.com/calmh/ipfix"
)

func TestCanCreateSession(t *testing.T) {
	p := ipfix.NewSession()

	if p == nil {
		t.Error("New session can't be nil")
	}
}

func TestParseTemplateSet(t *testing.T) {
	packet, _ := hex.DecodeString("000a008c51ec4264000000000b20bdbe0002007c283b0008001c0010800c000400003c258003000800003c258004000800003c258012ffff00003c258001ffff00003c25801cffff00003c25001b0010c2ac0008000c0004800c000400003c258003000800003c258004000800003c258012ffff00003c258001ffff00003c25801cffff00003c2500080004")
	p := ipfix.NewSession()

	r := bytes.NewBuffer(packet)
	msg, err := p.ParseReader(r)
	if err != nil {
		t.Fatal("ParseReader failed", err)
	}

	if msg.Header.Version != uint16(0xa) {
		t.Errorf("version mismatch %d != 10", msg.Header.Version)
	}
	if len(msg.DataRecords) != 0 {
		t.Error("Incorrect number of data records", len(msg.DataRecords))
	}
	if len(msg.TemplateRecords) != 2 {
		t.Error("Incorrect number of template records", len(msg.TemplateRecords))
	}
	if id := msg.TemplateRecords[0].TemplateID; id != uint16(10299) {
		t.Error("Incorrect template ID", id)
	}
	if id := msg.TemplateRecords[1].TemplateID; id != uint16(49836) {
		t.Error("Incorrect template ID", id)
	}
}

func TestParseTemplateIDAliasing(t *testing.T) {
	packet, _ := hex.DecodeString("000a017c51ec4264000000000b20bdbe0002016c283b0008001c0010800c000400003c258003000800003c258004000800003c258012ffff00003c258001ffff00003c25801cffff00003c25001b0010c2ac0008000c0004800c000400003c258003000800003c258004000800003c258012ffff00003c258001ffff00003c25801cffff00003c250008000412340008001c0010800c000400003c258003000800003c258004000800003c258012ffff00003c258001ffff00003c25801cffff00003c25001b0010abcd0008000c0004800c000400003c258003000800003c258004000800003c258012ffff00003c258001ffff00003c25801cffff00003c250008000412340008001c0010800c000400003c258003000800003c258004000800003c258012ffff00003c258001ffff00003c25801cffff00003c25001b0010abcd0008000c0004800c000400003c258003000800003c258004000800003c258012ffff00003c258001ffff00003c25801cffff00003c2500080004")
	p := ipfix.NewSession(ipfix.WithIDAliasing(true))

	r := bytes.NewBuffer(packet)
	msg, err := p.ParseBuffer(r.Bytes())
	if err != nil {
		t.Fatal("ParseBuffer failed", err)
	}

	if msg.Header.Version != uint16(0xa) {
		t.Errorf("version mismatch %d != 10", msg.Header.Version)
	}
	if len(msg.DataRecords) != 0 {
		t.Error("Incorrect number of data records", len(msg.DataRecords))
	}
	if len(msg.TemplateRecords) != 6 {
		t.Error("Incorrect number of template records", len(msg.TemplateRecords))
	}

	for i := 0; i < len(msg.TemplateRecords); i += 2 {
		if id := msg.TemplateRecords[i].TemplateID; id != uint16(256) {
			t.Error("Incorrect template ID", id)
		}
		if id := msg.TemplateRecords[i+1].TemplateID; id != uint16(257) {
			t.Error("Incorrect template ID", id)
		}
	}
}

func TestParseDataSet(t *testing.T) {
	testParseDataSet(false, t)
}

func TestParseDataSetWithAliasing(t *testing.T) {
	testParseDataSet(true, t)
}

func testParseDataSet(withAliasing bool, t *testing.T) {
	p0, _ := hex.DecodeString("000a008c51ec4264000000000b20bdbe0002007c283b0008001c0010800c000400003c258003000800003c258004000800003c258012ffff00003c258001ffff00003c25801cffff00003c25001b0010c2ac0008000c0004800c000400003c258003000800003c258004000800003c258012ffff00003c258001ffff00003c25801cffff00003c2500080004")
	p1, _ := hex.DecodeString("000a05b051ec4270000000000b20bdbec2ac05a0ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043000116fcb8ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b525043005e489f46ac10200300000026000000000000019f0000000000000160000e4265696e6720616e616c797a656400c27ef905ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043007aa7519c0808080800000000000000000000008d00000000000000550003444e5300ac102082ac10200f0000000000000000000000940000000000000147000f426974546f7272656e74204b52504300b228265c1859c1570000000000000000000000000000000000000064000f426974546f7272656e74204b52504300ac10200fac10200f0000000000000000000000920000000000000145000f426974546f7272656e74204b525043007b75a68ad92bb37f00000000000000000000006e0000000000000064000f426974546f7272656e74204b52504300ac10200fac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043004f972c247449d8f200000000000000000000006e0000000000000064000f426974546f7272656e74204b52504300ac10200fac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b5250430048b682a4ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b52504300595cc40dac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b5250430057451cc1ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b525043005465e5a8ac1020ff00000000000000000000000000000000000000af001a44726f70626f78204c414e2073796e6320646973636f766572790764726f70626f78ac102013ac10200f00000000000000000000008f000000000000014b000f426974546f7272656e74204b5250430001ab3c06ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b52504300befcacc8ffffffff00000000000000000000000000000000000000af001a44726f70626f78204c414e2073796e6320646973636f766572790764726f70626f78ac102013ac10200300000025000000000000019e0000000000000167000e4265696e6720616e616c797a656400c27ef905ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043006ca28bcdac10200f000000000000000000000091000000000000011c000f426974546f7272656e74204b52504300b13531caac10200f000000000000000000000068000000000000005f000f426974546f7272656e74204b5250430053df9212ac10200f0000000000000000000000940000000000000159000f426974546f7272656e74204b525043005f43f0b2ac10200f0000000000000000000001220000000000000252000f426974546f7272656e74204b52504300567ce6fbac10200100000000000000000000005a000000000000005a00034e545000ac102080ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b5250430055550ef7ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b52504300ba9322a2ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043004579e7114b01bf5300000000000000000000006e0000000000000064000f426974546f7272656e74204b52504300ac10200fac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043005cf46adf")
	b := new(bytes.Buffer)
	p := ipfix.NewSession(ipfix.WithIDAliasing(withAliasing))

	// Handle a data set with an unknown template first
	b.Write(p1)
	b.Write(p0)
	b.Write(p1)

	msg, err := p.ParseReader(b)
	if err != nil {
		t.Fatal("ParseReader failed", err)
	}

	if len(msg.DataRecords) != 0 {
		t.Error("Incorrect number of data records", len(msg.DataRecords))
	}
	if len(msg.TemplateRecords) != 0 {
		t.Error("Incorrect number of template records", len(msg.TemplateRecords))
	}

	msg, err = p.ParseReader(b)
	if err != nil {
		t.Fatal("ParseReader failed", err)
	}

	if len(msg.DataRecords) != 0 {
		t.Error("Incorrect number of data records", len(msg.DataRecords))
	}
	if len(msg.TemplateRecords) != 2 {
		t.Error("Incorrect number of template records", len(msg.TemplateRecords))
	}

	msg, err = p.ParseReader(b)
	if err != nil {
		t.Fatal("ParseReader failed", err)
	}

	if len(msg.DataRecords) != 31 {
		t.Error("Incorrect number of data records", len(msg.DataRecords))
	}
	if len(msg.TemplateRecords) != 0 {
		t.Error("Incorrect number of template records", len(msg.TemplateRecords))
	}
}

func TestReadParseBuffer(t *testing.T) {
	p0, _ := hex.DecodeString("000a008c51ec4264000000000b20bdbe0002007c283b0008001c0010800c000400003c258003000800003c258004000800003c258012ffff00003c258001ffff00003c25801cffff00003c25001b0010c2ac0008000c0004800c000400003c258003000800003c258004000800003c258012ffff00003c258001ffff00003c25801cffff00003c2500080004")
	p1, _ := hex.DecodeString("000a05b051ec4270000000000b20bdbec2ac05a0ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043000116fcb8ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b525043005e489f46ac10200300000026000000000000019f0000000000000160000e4265696e6720616e616c797a656400c27ef905ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043007aa7519c0808080800000000000000000000008d00000000000000550003444e5300ac102082ac10200f0000000000000000000000940000000000000147000f426974546f7272656e74204b52504300b228265c1859c1570000000000000000000000000000000000000064000f426974546f7272656e74204b52504300ac10200fac10200f0000000000000000000000920000000000000145000f426974546f7272656e74204b525043007b75a68ad92bb37f00000000000000000000006e0000000000000064000f426974546f7272656e74204b52504300ac10200fac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043004f972c247449d8f200000000000000000000006e0000000000000064000f426974546f7272656e74204b52504300ac10200fac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b5250430048b682a4ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b52504300595cc40dac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b5250430057451cc1ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b525043005465e5a8ac1020ff00000000000000000000000000000000000000af001a44726f70626f78204c414e2073796e6320646973636f766572790764726f70626f78ac102013ac10200f00000000000000000000008f000000000000014b000f426974546f7272656e74204b5250430001ab3c06ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b52504300befcacc8ffffffff00000000000000000000000000000000000000af001a44726f70626f78204c414e2073796e6320646973636f766572790764726f70626f78ac102013ac10200300000025000000000000019e0000000000000167000e4265696e6720616e616c797a656400c27ef905ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043006ca28bcdac10200f000000000000000000000091000000000000011c000f426974546f7272656e74204b52504300b13531caac10200f000000000000000000000068000000000000005f000f426974546f7272656e74204b5250430053df9212ac10200f0000000000000000000000940000000000000159000f426974546f7272656e74204b525043005f43f0b2ac10200f0000000000000000000001220000000000000252000f426974546f7272656e74204b52504300567ce6fbac10200100000000000000000000005a000000000000005a00034e545000ac102080ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b5250430055550ef7ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b52504300ba9322a2ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043004579e7114b01bf5300000000000000000000006e0000000000000064000f426974546f7272656e74204b52504300ac10200fac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043005cf46adf")
	b := new(bytes.Buffer)
	p := ipfix.NewSession()

	// Handle a data set with an unknown template first
	b.Write(p1)
	b.Write(p0)
	b.Write(p1)

	bs, _, err := ipfix.Read(b, nil)
	if err != nil {
		t.Fatal(err)
	}

	msg, err := p.ParseBuffer(bs)
	if err != nil {
		t.Fatal("ParseBuffer failed", err)
	}

	if len(msg.DataRecords) != 0 {
		t.Error("Incorrect number of data records", len(msg.DataRecords))
	}
	if len(msg.TemplateRecords) != 0 {
		t.Error("Incorrect number of template records", len(msg.TemplateRecords))
	}

	bs, _, err = ipfix.Read(b, bs)
	if err != nil {
		t.Fatal(err)
	}

	msg, err = p.ParseBuffer(bs)
	if err != nil {
		t.Fatal("ParseBuffer failed", err)
	}

	if len(msg.DataRecords) != 0 {
		t.Error("Incorrect number of data records", len(msg.DataRecords))
	}
	if len(msg.TemplateRecords) != 2 {
		t.Error("Incorrect number of template records", len(msg.TemplateRecords))
	}

	bs, _, err = ipfix.Read(b, bs)
	if err != nil {
		t.Fatal(err)
	}

	msg, err = p.ParseBuffer(bs)
	if err != nil {
		t.Fatal("ParseBuffer failed", err)
	}

	if len(msg.DataRecords) != 31 {
		t.Error("Incorrect number of data records", len(msg.DataRecords))
	}
	if len(msg.TemplateRecords) != 0 {
		t.Error("Incorrect number of template records", len(msg.TemplateRecords))
	}
}

func TestParallellParseBuffer(t *testing.T) {
	p0, _ := hex.DecodeString("000a008c51ec4264000000000b20bdbe0002007c283b0008001c0010800c000400003c258003000800003c258004000800003c258012ffff00003c258001ffff00003c25801cffff00003c25001b0010c2ac0008000c0004800c000400003c258003000800003c258004000800003c258012ffff00003c258001ffff00003c25801cffff00003c2500080004")
	p1, _ := hex.DecodeString("000a05b051ec4270000000000b20bdbec2ac05a0ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043000116fcb8ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b525043005e489f46ac10200300000026000000000000019f0000000000000160000e4265696e6720616e616c797a656400c27ef905ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043007aa7519c0808080800000000000000000000008d00000000000000550003444e5300ac102082ac10200f0000000000000000000000940000000000000147000f426974546f7272656e74204b52504300b228265c1859c1570000000000000000000000000000000000000064000f426974546f7272656e74204b52504300ac10200fac10200f0000000000000000000000920000000000000145000f426974546f7272656e74204b525043007b75a68ad92bb37f00000000000000000000006e0000000000000064000f426974546f7272656e74204b52504300ac10200fac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043004f972c247449d8f200000000000000000000006e0000000000000064000f426974546f7272656e74204b52504300ac10200fac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b5250430048b682a4ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b52504300595cc40dac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b5250430057451cc1ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b525043005465e5a8ac1020ff00000000000000000000000000000000000000af001a44726f70626f78204c414e2073796e6320646973636f766572790764726f70626f78ac102013ac10200f00000000000000000000008f000000000000014b000f426974546f7272656e74204b5250430001ab3c06ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b52504300befcacc8ffffffff00000000000000000000000000000000000000af001a44726f70626f78204c414e2073796e6320646973636f766572790764726f70626f78ac102013ac10200300000025000000000000019e0000000000000167000e4265696e6720616e616c797a656400c27ef905ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043006ca28bcdac10200f000000000000000000000091000000000000011c000f426974546f7272656e74204b52504300b13531caac10200f000000000000000000000068000000000000005f000f426974546f7272656e74204b5250430053df9212ac10200f0000000000000000000000940000000000000159000f426974546f7272656e74204b525043005f43f0b2ac10200f0000000000000000000001220000000000000252000f426974546f7272656e74204b52504300567ce6fbac10200100000000000000000000005a000000000000005a00034e545000ac102080ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b5250430055550ef7ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b52504300ba9322a2ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043004579e7114b01bf5300000000000000000000006e0000000000000064000f426974546f7272656e74204b52504300ac10200fac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043005cf46adf")

	p := ipfix.NewSession()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for i := 0; i < 1000; i++ {
			p.ParseBuffer(p0)
		}
	}()

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			p.ParseBuffer(p0)
			for i := 0; i < 1000; i++ {
				msg, err := p.ParseBuffer(p1)
				if err != nil {
					t.Fatal("ParseReader failed", err)
				}

				if len(msg.DataRecords) != 31 {
					t.Error("Incorrect number of data records", len(msg.DataRecords))
				}
				if len(msg.TemplateRecords) != 0 {
					t.Error("Incorrect number of template records", len(msg.TemplateRecords))
				}
			}
		}()
	}
}

func TestEOFError(t *testing.T) {
	truncated, _ := hex.DecodeString("000a008c51ec4264000000000b20bdbe0002007c283b0008001c0010800c000400003c258003000800003c258004000800003c258012ffff00003c258001ffff00003c25801cffff00003c25001b0010c2ac0008000c0004800c000400003c258003000800003c258004000800003c258")
	b := new(bytes.Buffer)
	p := ipfix.NewSession()

	b.Write(truncated)

	_, err := p.ParseReader(b)
	if err != io.EOF {
		t.Fatalf("Received %v instead of io.EOF error", err)
	}
}

func TestVersionError(t *testing.T) {
	truncated, _ := hex.DecodeString("000a05b051ec4270000000000b20bdbec2ac05a0ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043000116fcb8ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b525043005e489f46ac10200300000026000000000000019f0000000000000160000e4265696e6720616e616c797a656400c27ef905ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043007aa7519c0808080800000000000000000000008d00000000000000550003444e5300ac102082ac10200f0000000000000000000000940000000000000147000f426974546f7272656e74204b52504300b228265c1859c1570000000000000000000000000000000000000064000f426974546f7272656e74204b52504300ac10200fac10200f0000000000000000000000920000000000000145000f426974546f7272656e74204b525043007b75a68ad92bb37f00000000000000000000006e0000000000000064000f426974546f7272656e74204b52504300ac10200fac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043004f972c247449d8f200000000000000000000006e0000000000000064000f426974546f7272656e74204b52504300ac10200fac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b5250430048b682a4ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b52504300595cc40dac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b5250430057451cc1ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b525043005465e5a8ac1020ff00000000000000000000000000000000000000af001a44726f70626f78204c414e2073796e6320646973636f766572790764726f70626f78ac102013ac10200f00000000000000000000008f000000000000014b000f426974546f7272656e74204b5250430001ab3c06ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b52504300befcacc8ffffffff00000000000000000000000000000000000000af001a44726f70626f78204c414e2073796e6320646973636f766572790764726f70626f78ac102013ac10200300000025000000000000019e0000000000000167000e4265696e6720616e616c797a656400c27ef905ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043006ca28bcdac10200f000000000000000000000091000000000000011c000f426974546f7272656e74204b52504300b13531caac10200f000000000000000000000068000000000000005f000f426974546f7272656e74204b5250430053df9212ac10200f0000000000000000000000940000000000000159000f426974546f7272656e74204b525043005f43f0b2ac10200f0000000000000000000001220000000000000252000f426974546f7272656e74204b52504300567ce6fbac10200100000000000000000000005a000000000000005a00034e545000ac102080ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b5250430055550ef7ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b52504300ba9322a2ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043004579e7114b01bf5300000000000000000000006e0000000000000064000f426974546f7272656e74204b52504300ac10200fac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043005cf46a")
	p1, _ := hex.DecodeString("000a05b051ec4270000000000b20bdbec2ac05a0ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043000116fcb8ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b525043005e489f46ac10200300000026000000000000019f0000000000000160000e4265696e6720616e616c797a656400c27ef905ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043007aa7519c0808080800000000000000000000008d00000000000000550003444e5300ac102082ac10200f0000000000000000000000940000000000000147000f426974546f7272656e74204b52504300b228265c1859c1570000000000000000000000000000000000000064000f426974546f7272656e74204b52504300ac10200fac10200f0000000000000000000000920000000000000145000f426974546f7272656e74204b525043007b75a68ad92bb37f00000000000000000000006e0000000000000064000f426974546f7272656e74204b52504300ac10200fac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043004f972c247449d8f200000000000000000000006e0000000000000064000f426974546f7272656e74204b52504300ac10200fac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b5250430048b682a4ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b52504300595cc40dac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b5250430057451cc1ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b525043005465e5a8ac1020ff00000000000000000000000000000000000000af001a44726f70626f78204c414e2073796e6320646973636f766572790764726f70626f78ac102013ac10200f00000000000000000000008f000000000000014b000f426974546f7272656e74204b5250430001ab3c06ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b52504300befcacc8ffffffff00000000000000000000000000000000000000af001a44726f70626f78204c414e2073796e6320646973636f766572790764726f70626f78ac102013ac10200300000025000000000000019e0000000000000167000e4265696e6720616e616c797a656400c27ef905ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043006ca28bcdac10200f000000000000000000000091000000000000011c000f426974546f7272656e74204b52504300b13531caac10200f000000000000000000000068000000000000005f000f426974546f7272656e74204b5250430053df9212ac10200f0000000000000000000000940000000000000159000f426974546f7272656e74204b525043005f43f0b2ac10200f0000000000000000000001220000000000000252000f426974546f7272656e74204b52504300567ce6fbac10200100000000000000000000005a000000000000005a00034e545000ac102080ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b5250430055550ef7ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b52504300ba9322a2ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043004579e7114b01bf5300000000000000000000006e0000000000000064000f426974546f7272656e74204b52504300ac10200fac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043005cf46adf")
	b := new(bytes.Buffer)
	p := ipfix.NewSession()

	b.Write(truncated)
	b.Write(p1)

	p.ParseReader(b)
	_, err := p.ParseReader(b)
	if err != ipfix.ErrVersion {
		t.Fatalf("Received %v instead of ipfix.ErrVersion error", err)
	}
}

func TestBadEncodingError(t *testing.T) {
	p0, _ := hex.DecodeString("000a009c520239cc002488cc0b20bdbe0002008c283b0008001c0010800c000400003c258003000800003c258004000800003c258012ffff00003c258001ffff00003c25801cffff00003c25001b00104f4d000b000c00040097000400960004800c000400003c258016ffff00003c258003000800003c258004000800003c258012ffff00003c258001ffff00003c25801cffff00003c2500080004")
	p1, _ := hex.DecodeString("000a05ae520239cc0024889e0b20bdbe4f4d029cac10200f520239cc520239ac000000000000000000000000910000000000000136000f426974546f7272656e74204b525043001b2065fbac102003520239cc520239c0000000220000000000000001ab0000000000000168000e4265696e6720616e616c797a656400c27ef905ac1020ff520239cc520239ac00000000000000000000000000000000000000005c00144e657442696f73204e616d65205365727669636500ac102082ac10200f520239cc520239ac000000000000000000000000910000000000000136000f426974546f7272656e74204b52504300dcee3ce9ac10200f520239cc520239ac0000000000000000000000008c000000000000013a000f426974546f7272656e74204b525043006fe9d570ac10200f520239cc520239ad000000000000000000000000910000000000000136000f426974546f7272656e74204b525043005853e3e4ac10200f520239cc520239ae000000000000000000000000910000000000000136000f426974546f7272656e74204b52504300b1621b8b18345cdc520239cc520239ae000000000000000000000000000000000000000064000f426974546f7272656e74204b52504300ac10200fac10200f520239cc520239ae0000000000000000000000008f000000000000014b000f426974546f7272656e74204b5250430025faf170ac10200f520239cc520239ae0000000000000000000000006b0000000000000059000f426974546f7272656e74204b52504300539525cdac10200f520239cc520239ae000000000000000000000000940000000000000147000f426974546f7272656e74204b52504300bc069bc5ac10200f520239cc520239ae0000000000000000000000008f000000000000014b000f426974546f7272656e74204b5250430056931e40283b004920010470002804d6000000000000000400000000000000000000004a0000000000000062000e4265696e6720616e616c797a65640020010470deeb003280db765c7c72c6934f4d02b9c27ef904520239cc520239c20000001dff0001687474703a2f2f6e796d2e73652f737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737373737300000000000004ff000000000000048f066e796d2e7365044854545000ac102082ac10200f520239cc520239a50000000000000000000000011e00000000000002ba000f426974546f7272656e74204b5250430001ace2e4ac10200f520239cc520239b10000000000000000000000006b0000000000000059000f426974546f7272656e74204b525043007894a711ac10200f520239cc520239b10000000000000000000000009f0000000000000134000f426974546f7272656e74204b525043005f8bc343ac10200f520239cc520239b1000000000000000000000000920000000000000145000f426974546f7272656e74204b5250430079e53f6eac10200f520239cc520239af000000000000000000000000910000000000000136000f426974546f7272656e74204b525043005be27a31ac10200f520239cc520239af000000000000000000000000940000000000000147000f426974546f7272656e74204b5250430070d14d8eac10200f520239cc520239b1000000000000000000000000920000000000000145000f426974546f7272656e74204b52504300de5e3a01")
	p2, _ := hex.DecodeString("000a05a6520239f9002489e30b20bdbe4f4d0596ac10200f520239f9520239e0000000000000000000000000f1000000000000005b000f426974546f7272656e74204b525043006dab2a88ac10200f520239f9520239e00000000000000000000000008f000000000000015d000f426974546f7272656e74204b525043007b778163ac10200f520239f9520239e0000000000000000000000000910000000000000136000f426974546f7272656e74204b525043005fb2e498ac10200f520239f9520239da000000000000000000000001ad000000000000042f000f426974546f7272656e74204b525043006ee75810ac10200f520239f9520239df000000000000000000000000910000000000000136000f426974546f7272656e74204b5250430005526896ac102003520239f9520239f3000000220000000000000001ad0000000000000167000e4265696e6720616e616c797a656400c27ef905ac102003520239f9520239f4000000240000000000000001a80000000000000167000e4265696e6720616e616c797a656400c27ef905ac10200f520239f9520239e1000000000000000000000000910000000000000136000f426974546f7272656e74204b525043000e2b6855ac10200f520239f9520239d5000000000000000000000001d00000000000000347000f426974546f7272656e74204b5250430077739721ac10200f520239f9520239dc000000000000000000000000fd0000000000000198000f426974546f7272656e74204b525043005bc8cb08ac10200f520239f9520239e10000000000000000000000008c000000000000013a000f426974546f7272656e74204b5250430057fcb5ceac10200f520239f9520239e1000000000000000000000000910000000000000136000f426974546f7272656e74204b52504300dfccf34aac10200f520239f9520239d90000000000000000000000011e0000000000000296000f426974546f7272656e74204b5250430075414462ac10200f520239f9520239e20000000000000000000000008f000000000000014b000f426974546f7272656e74204b525043003d5b581bac10200f520239f9520239e2000000000000000000000000910000000000000136000f426974546f7272656e74204b52504300050c9e43ac10200f520239f9520239e2000000000000000000000000910000000000000136000f426974546f7272656e74204b52504300ae0335b2ac102003520239f9520239f5000000250000000000000001ab0000000000000167000e4265696e6720616e616c797a656400c27ef905ac10200f520239f9520239e2000000000000000000000000910000000000000136000f426974546f7272656e74204b525043004e61163cac10200f520239f9520239d800000000000000000000000091000000000000019a000f426974546f7272656e74204b525043007ab7e026ac10200f520239f9520239e2000000000000000000000000910000000000000136000f426974546f7272656e74204b5250430029d08ffcac10200f520239f9520239dd000000000000000000000002820000000000000378000f426974546f7272656e74204b52504300b71ec470ac10200f520239f9520239e20000000000000000000000008f0000000000000134000f426974546f7272656e74204b525043005ccb6c11ac10200f520239f9520239e20000000000000000000000008f0000000000000195000f426974546f7272656e74204b5250430071007099ac10200f520239f9520239e20000000000000000000000008f0000000000000134000f426974546f7272656e74204b52504300c4000455ac102003520239f9520239f60000001e0000000000000001a90000000000000167000e4265696e6720616e616c797a656400c27ef905ac10200f520239f9520239e3000000000000000000000000910000000000000136000f426974546f7272656e74204b525043005fe287cd")
	b := new(bytes.Buffer)
	p := ipfix.NewSession()

	b.Write(p0)
	b.Write(p1)
	b.Write(p2)

	p.ParseReader(b)             // The template
	p.ParseReader(b)             // The broken message
	msg, err := p.ParseReader(b) // The good message
	if err != nil {
		t.Fatalf("Received error %v even though we should be back in sync", err)
	}
	if len(msg.DataRecords) != 26 {
		t.Error("Incorrect number of data records", len(msg.DataRecords))
	}
	if len(msg.TemplateRecords) != 0 {
		t.Error("Incorrect number of template records", len(msg.TemplateRecords))
	}
}

var msg ipfix.Message

func BenchmarkParseReader(b *testing.B) {
	p0, _ := hex.DecodeString("000a008c51ec4264000000000b20bdbe0002007c283b0008001c0010800c000400003c258003000800003c258004000800003c258012ffff00003c258001ffff00003c25801cffff00003c25001b0010c2ac0008000c0004800c000400003c258003000800003c258004000800003c258012ffff00003c258001ffff00003c25801cffff00003c2500080004")
	p1, _ := hex.DecodeString("000a05b051ec4270000000000b20bdbec2ac05a0ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043000116fcb8ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b525043005e489f46ac10200300000026000000000000019f0000000000000160000e4265696e6720616e616c797a656400c27ef905ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043007aa7519c0808080800000000000000000000008d00000000000000550003444e5300ac102082ac10200f0000000000000000000000940000000000000147000f426974546f7272656e74204b52504300b228265c1859c1570000000000000000000000000000000000000064000f426974546f7272656e74204b52504300ac10200fac10200f0000000000000000000000920000000000000145000f426974546f7272656e74204b525043007b75a68ad92bb37f00000000000000000000006e0000000000000064000f426974546f7272656e74204b52504300ac10200fac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043004f972c247449d8f200000000000000000000006e0000000000000064000f426974546f7272656e74204b52504300ac10200fac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b5250430048b682a4ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b52504300595cc40dac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b5250430057451cc1ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b525043005465e5a8ac1020ff00000000000000000000000000000000000000af001a44726f70626f78204c414e2073796e6320646973636f766572790764726f70626f78ac102013ac10200f00000000000000000000008f000000000000014b000f426974546f7272656e74204b5250430001ab3c06ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b52504300befcacc8ffffffff00000000000000000000000000000000000000af001a44726f70626f78204c414e2073796e6320646973636f766572790764726f70626f78ac102013ac10200300000025000000000000019e0000000000000167000e4265696e6720616e616c797a656400c27ef905ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043006ca28bcdac10200f000000000000000000000091000000000000011c000f426974546f7272656e74204b52504300b13531caac10200f000000000000000000000068000000000000005f000f426974546f7272656e74204b5250430053df9212ac10200f0000000000000000000000940000000000000159000f426974546f7272656e74204b525043005f43f0b2ac10200f0000000000000000000001220000000000000252000f426974546f7272656e74204b52504300567ce6fbac10200100000000000000000000005a000000000000005a00034e545000ac102080ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b5250430055550ef7ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b52504300ba9322a2ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043004579e7114b01bf5300000000000000000000006e0000000000000064000f426974546f7272656e74204b52504300ac10200fac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043005cf46adf")
	pb := new(bytes.Buffer)
	pb.Write(p0)

	p := ipfix.NewSession()
	_, err := p.ParseReader(pb)
	if err != nil {
		b.Fatal("ParseReader failed", err)
	}

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(1)

	for i := 0; i < b.N; {
		pb.Write(p1)
		msg, err = p.ParseReader(pb)
		if err != nil {
			b.Error("ParseReader failed", err)
		}
		i += len(msg.DataRecords) + len(msg.TemplateRecords)
	}
}

func BenchmarkParseBuffer(b *testing.B) {
	p0, _ := hex.DecodeString("000a008c51ec4264000000000b20bdbe0002007c283b0008001c0010800c000400003c258003000800003c258004000800003c258012ffff00003c258001ffff00003c25801cffff00003c25001b0010c2ac0008000c0004800c000400003c258003000800003c258004000800003c258012ffff00003c258001ffff00003c25801cffff00003c2500080004")
	p1, _ := hex.DecodeString("000a05b051ec4270000000000b20bdbec2ac05a0ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043000116fcb8ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b525043005e489f46ac10200300000026000000000000019f0000000000000160000e4265696e6720616e616c797a656400c27ef905ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043007aa7519c0808080800000000000000000000008d00000000000000550003444e5300ac102082ac10200f0000000000000000000000940000000000000147000f426974546f7272656e74204b52504300b228265c1859c1570000000000000000000000000000000000000064000f426974546f7272656e74204b52504300ac10200fac10200f0000000000000000000000920000000000000145000f426974546f7272656e74204b525043007b75a68ad92bb37f00000000000000000000006e0000000000000064000f426974546f7272656e74204b52504300ac10200fac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043004f972c247449d8f200000000000000000000006e0000000000000064000f426974546f7272656e74204b52504300ac10200fac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b5250430048b682a4ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b52504300595cc40dac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b5250430057451cc1ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b525043005465e5a8ac1020ff00000000000000000000000000000000000000af001a44726f70626f78204c414e2073796e6320646973636f766572790764726f70626f78ac102013ac10200f00000000000000000000008f000000000000014b000f426974546f7272656e74204b5250430001ab3c06ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b52504300befcacc8ffffffff00000000000000000000000000000000000000af001a44726f70626f78204c414e2073796e6320646973636f766572790764726f70626f78ac102013ac10200300000025000000000000019e0000000000000167000e4265696e6720616e616c797a656400c27ef905ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043006ca28bcdac10200f000000000000000000000091000000000000011c000f426974546f7272656e74204b52504300b13531caac10200f000000000000000000000068000000000000005f000f426974546f7272656e74204b5250430053df9212ac10200f0000000000000000000000940000000000000159000f426974546f7272656e74204b525043005f43f0b2ac10200f0000000000000000000001220000000000000252000f426974546f7272656e74204b52504300567ce6fbac10200100000000000000000000005a000000000000005a00034e545000ac102080ac10200f00000000000000000000008c000000000000013a000f426974546f7272656e74204b5250430055550ef7ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b52504300ba9322a2ac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043004579e7114b01bf5300000000000000000000006e0000000000000064000f426974546f7272656e74204b52504300ac10200fac10200f0000000000000000000000910000000000000136000f426974546f7272656e74204b525043005cf46adf")

	p := ipfix.NewSession()
	_, err := p.ParseBuffer(p0)
	if err != nil {
		b.Fatal("ParseReader failed", err)
	}

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(1)

	for i := 0; i < b.N; {
		msg, err = p.ParseBuffer(p1)
		if err != nil {
			b.Error("ParseReader failed", err)
		}
		i += len(msg.DataRecords) + len(msg.TemplateRecords)
	}
}

func TestParsingTemplateAndDataRecords(t *testing.T) {
	packet, _ := hex.DecodeString("000a00405685b3700000000000bc614e000200140100000300080004000c0004000200040100001cc0a800c9c0a80001000000ebc0a800cac0a800010000002a")
	p := ipfix.NewSession()

	msg, err := p.ParseBuffer(packet)
	if err != nil {
		t.Fatal("ParseBuffer failed", err)
	}

	if len(msg.TemplateRecords) != 1 {
		t.Error("Incorrect number of template records", len(msg.TemplateRecords))
	}
	if len(msg.DataRecords) != 2 {
		t.Error("Incorrect number of data records", len(msg.DataRecords))
	}
}
