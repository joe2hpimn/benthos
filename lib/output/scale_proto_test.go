/*
Copyright (c) 2014 Ashley Jeffs

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package output

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/pull"
	"github.com/go-mangos/mangos/transport/tcp"

	"github.com/jeffail/benthos/lib/types"
	"github.com/jeffail/util/log"
	"github.com/jeffail/util/metrics"
)

//--------------------------------------------------------------------------------------------------

func TestScaleProtoBasic(t *testing.T) {
	nTestLoops := 1000

	sendChan := make(chan types.Message)

	conf := NewConfig()
	conf.ScaleProto.Address = "tcp://localhost:1324"
	conf.ScaleProto.Bind = true
	conf.ScaleProto.SocketType = "PUSH"

	s, err := NewScaleProto(conf, log.NewLogger(os.Stdout, logConfig), metrics.DudType{})
	if err != nil {
		t.Error(err)
		return
	}

	if err = s.StartReceiving(sendChan); err != nil {
		t.Error(err)
		return
	}

	socket, err := pull.NewSocket()
	if err != nil {
		t.Error(err)
		return
	}

	socket.AddTransport(tcp.NewTransport())
	socket.SetOption(mangos.OptionRecvDeadline, time.Second)

	if err = socket.Dial("tcp://localhost:1324"); err != nil {
		t.Error(err)
		return
	}

	for i := 0; i < nTestLoops; i++ {
		testStr := fmt.Sprintf("test%v", i)
		testMsg := types.Message{Parts: [][]byte{[]byte(testStr)}}

		select {
		case sendChan <- testMsg:
		case <-time.After(time.Second):
			t.Errorf("Action timed out")
			return
		}

		data, err := socket.Recv()
		if err != nil {
			t.Error(err)
			return
		}
		if res := string(data); res != testStr {
			t.Errorf("Wrong value on output: %v != %v", res, testStr)
		}

		select {
		case res := <-s.ResponseChan():
			if res.Error() != nil {
				t.Error(res.Error())
				return
			}
		case <-time.After(time.Second):
			t.Errorf("Action timed out")
			return
		}
	}

	for i := 0; i < nTestLoops; i++ {
		testStr := fmt.Sprintf("test%v", i)
		testMsg := types.Message{Parts: [][]byte{
			[]byte(testStr + "PART-A"),
			[]byte(testStr + "PART-B"),
		}}

		select {
		case sendChan <- testMsg:
		case <-time.After(time.Second):
			t.Errorf("Action timed out")
			return
		}

		data, err := socket.Recv()
		if err != nil {
			t.Error(err)
			return
		}
		msg, err := types.FromBytes(data)
		if err != nil {
			t.Error(err)
			return
		}
		if exp, actual := 2, len(msg.Parts); exp != actual {
			t.Errorf("Unexpected message parts received: %v != %v", exp, actual)
			return
		}
		if exp, actual := testStr+"PART-A", string(msg.Parts[0]); exp != actual {
			t.Errorf("Unexpected message received: %v != %v", exp, actual)
			return
		}
		if exp, actual := testStr+"PART-B", string(msg.Parts[1]); exp != actual {
			t.Errorf("Unexpected message received: %v != %v", exp, actual)
			return
		}

		select {
		case res := <-s.ResponseChan():
			if res.Error() != nil {
				t.Error(res.Error())
				return
			}
		case <-time.After(time.Second):
			t.Errorf("Action timed out")
			return
		}
	}
}

//--------------------------------------------------------------------------------------------------
