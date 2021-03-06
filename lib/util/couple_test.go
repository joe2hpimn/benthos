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

package util

import (
	"testing"

	"github.com/jeffail/benthos/lib/input"
	"github.com/jeffail/benthos/lib/output"
	"github.com/jeffail/benthos/lib/types"
)

//--------------------------------------------------------------------------------------------------

func TestCouple(t *testing.T) {
	inMock := input.MockType{MsgChan: make(chan types.Message)}
	outMock := output.MockType{ResChan: make(chan types.Response)}

	if inMock.MsgChan == outMock.MsgChan {
		t.Errorf("Message channels should not yet match: %v == %v", inMock.MsgChan, outMock.MsgChan)
	}
	if inMock.ResChan == outMock.ResChan {
		t.Errorf("Response channels should not yet match: %v == %v", inMock.ResChan, outMock.ResChan)
	}

	if err := Couple(&inMock, &outMock); err != nil {
		t.Error(err)
	}
	if inMock.MsgChan != outMock.MsgChan {
		t.Errorf("Non-matching message channels: %v != %v", inMock.MsgChan, outMock.MsgChan)
	}
	if inMock.ResChan != outMock.ResChan {
		t.Errorf("Non-matching response channels: %v != %v", inMock.ResChan, outMock.ResChan)
	}
}

//--------------------------------------------------------------------------------------------------
