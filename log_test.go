/*
 * Copyright GoIIoT (https://github.com/goiiot)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package libmqtt

import "testing"

func Test_SilentLogger(t *testing.T) {
	if l := newLogger(Silent); l != nil {
		t.Fail()
	} else {
		l.v("test")
		l.d("test")
		l.i("test")
		l.w("test")
		l.e("test")
	}
}

func Test_ErrorLogger(t *testing.T) {
	if l := newLogger(Error); l == nil ||
		l.error == nil || l.warning != nil ||
		l.info != nil || l.debug != nil || l.verbose != nil {
		t.Fail()
	} else {
		l.v("test")
		l.d("test")
		l.i("test")
		l.w("test")
		l.e("test")
	}
}

func Test_WarningLogger(t *testing.T) {
	if l := newLogger(Warning); l == nil ||
		l.error == nil || l.warning == nil ||
		l.info != nil || l.debug != nil || l.verbose != nil {
		t.Log("failed at warning logger")
		t.Fail()
	} else {
		l.v("test")
		l.d("test")
		l.i("test")
		l.w("test")
		l.e("test")
	}
}

func Test_InfoLogger(t *testing.T) {
	if l := newLogger(Info); l == nil ||
		l.error == nil || l.warning == nil ||
		l.info == nil || l.debug != nil || l.verbose != nil {
		t.Log("failed at info logger")
		t.Fail()
	} else {
		l.v("test")
		l.d("test")
		l.i("test")
		l.w("test")
		l.e("test")
	}
}

func Test_DebugLogger(t *testing.T) {
	if l := newLogger(Debug); l == nil ||
		l.error == nil || l.warning == nil ||
		l.info == nil || l.debug == nil || l.verbose != nil {
		t.Log("failed at debug logger")
		t.Fail()
	} else {
		l.v("test")
		l.d("test")
		l.i("test")
		l.w("test")
		l.e("test")
	}

}

func Test_VerboseLogger(t *testing.T) {
	if l := newLogger(Verbose); l == nil ||
		l.error == nil || l.warning == nil ||
		l.info == nil || l.debug == nil || l.verbose == nil {
		t.Log("failed at verbose logger")
		t.Fail()
	} else {
		l.v("test")
		l.d("test")
		l.i("test")
		l.w("test")
		l.e("test")
	}
}
