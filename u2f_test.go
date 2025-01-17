// Go FIDO U2F Library
// Copyright 2015 The Go FIDO U2F Library Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package u2f

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestFull(t *testing.T) {
	testFunc := func(t *testing.T, useConfig bool) {
		// These are actual responses from a Yubikey with Chrome.

		const appID = "http://localhost:3483"

		timeNow := time.Now()

		var config *Config
		if useConfig {
			config = &Config{
				Time: func() time.Time { return timeNow },
			}
		}

		challengeConfig := func(cenc string) *Config {
			c, err := decodeBase64(cenc)
			if err != nil {
				panic(err)
			}
			return &Config{
				Rand: bytes.NewBuffer(c),
				Time: func() time.Time { return timeNow.Add(-time.Minute) },
			}
		}

		registerChallenge, err := NewChallengeConfig(appID, []string{appID}, challengeConfig("s4UJ3wkN80p4wLjyI2Guv-_a-s7LV54Ic9PAZvHo_lM"))
		if err != nil {
			t.Error(err)
		}

		const regRespJSON = "{\"registrationData\":\"BQTD17IP7bZ3Gcd7l5Ao4qqohsUcm0bcXgHLpn0pv2VWNl7SBtNFo0wEoAdMrHlFXGzJgQz_bRZaKXZfHyd3fAo0QJmZkSv9ZbTKz7TVO6jnOcKGrSHb15JDatMMFxHxN5BR56CE3sj10jtGOY7szQIi4RGU6kONIuriAarxuEFJ5IswggIcMIIBBqADAgECAgQk26tAMAsGCSqGSIb3DQEBCzAuMSwwKgYDVQQDEyNZdWJpY28gVTJGIFJvb3QgQ0EgU2VyaWFsIDQ1NzIwMDYzMTAgFw0xNDA4MDEwMDAwMDBaGA8yMDUwMDkwNDAwMDAwMFowKzEpMCcGA1UEAwwgWXViaWNvIFUyRiBFRSBTZXJpYWwgMTM1MDMyNzc4ODgwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAQCsJS-NH1HeUHEd46-xcpN7SpHn6oeb-w5r-veDCBwy1vUvWnJanjjv4dR_rV5G436ysKUAXUcsVe5fAnkORo2oxIwEDAOBgorBgEEAYLECgEBBAAwCwYJKoZIhvcNAQELA4IBAQCjY64OmDrzC7rxLIst81pZvxy7ShsPy2jEhFWEkPaHNFhluNsCacNG5VOITCxWB68OonuQrIzx70MfcqwYnbIcgkkUvxeIpVEaM9B7TI40ZHzp9h4VFqmps26QCkAgYfaapG4SxTK5k_lCPvqqTPmjtlS03d7ykkpUj9WZlVEN1Pf02aTVIZOHPHHJuH6GhT6eLadejwxtKDBTdNTv3V4UlvjDOQYQe9aL1jUNqtLDeBHso8pDvJMLc0CX3vadaI2UVQxM-xip4kuGouXYj0mYmaCbzluBDFNsrzkNyL3elg3zMMrKvAUhoYMjlX_-vKWcqQsgsQ0JtSMcWMJ-umeDMEQCIApTYovLr8citOpIKkyNidCQz7UeSOWNMlPBB-s3r4G9AiAskXkh7iale4QDe6a-675L3xzohYb8Fcvz3gH6dkDLvw\",\"version\":\"U2F_V2\",\"challenge\":\"s4UJ3wkN80p4wLjyI2Guv-_a-s7LV54Ic9PAZvHo_lM\",\"appId\":\"http://localhost:3483\",\"clientData\":\"eyJ0eXAiOiJuYXZpZ2F0b3IuaWQuZmluaXNoRW5yb2xsbWVudCIsImNoYWxsZW5nZSI6InM0VUozd2tOODBwNHdManlJMkd1di1fYS1zN0xWNTRJYzlQQVp2SG9fbE0iLCJvcmlnaW4iOiJodHRwOi8vbG9jYWxob3N0OjM0ODMiLCJjaWRfcHVia2V5IjoiIn0\"}"
		var regResp RegisterResponse
		if err := json.Unmarshal([]byte(regRespJSON), &regResp); err != nil {
			t.Error(err)
		}

		reg, err := Register(regResp, *registerChallenge, config)
		if err != nil {
			t.Error(err)
		}

		authChallenge, err := NewChallengeConfig(appID, []string{appID}, challengeConfig("PzN6SGiUaeypErE3SCHeRlkRxVwfWlGVi35gfq6LsdY"))
		if err != nil {
			t.Error(err)
		}

		const signRespJSON = "{\"keyHandle\":\"mZmRK_1ltMrPtNU7qOc5woatIdvXkkNq0wwXEfE3kFHnoITeyPXSO0Y5juzNAiLhEZTqQ40i6uIBqvG4QUnkiw\",\"clientData\":\"eyJ0eXAiOiJuYXZpZ2F0b3IuaWQuZ2V0QXNzZXJ0aW9uIiwiY2hhbGxlbmdlIjoiUHpONlNHaVVhZXlwRXJFM1NDSGVSbGtSeFZ3ZldsR1ZpMzVnZnE2THNkWSIsIm9yaWdpbiI6Imh0dHA6Ly9sb2NhbGhvc3Q6MzQ4MyIsImNpZF9wdWJrZXkiOiIifQ\",\"signatureData\":\"AQAAAAYwRAIgBuyafOXoc9Q7fARcs2JbCZdtnMzVCyeJC-J-2Im1IBsCIDxkzmvPX9RCY8uts4wM1y4wEX9LmNH2Mz_VFd-JdyGE\"}"
		var signResp SignResponse
		if err := json.Unmarshal([]byte(signRespJSON), &signResp); err != nil {
			t.Error(err)
		}

		regAuthenticate := reg.Authenticate
		if useConfig {
			regAuthenticate = func(resp SignResponse, c Challenge, counter uint32) (newCounter uint32, err error) {
				return reg.AuthenticateConfig(resp, c, counter, config)
			}
		}

		newCounter, err := regAuthenticate(signResp, *authChallenge, 0)
		if err != nil {
			t.Error(err)
		}
		if newCounter != 6 {
			t.Errorf("Wrong new counter: %d", newCounter)
		}

		newCounter, err = regAuthenticate(signResp, *authChallenge, 7)
		if err == nil {
			t.Errorf("Expected error due to decreasing counter")
		}
	}

	for _, useConfig := range []bool{false, true} {
		t.Run(fmt.Sprintf("useConfig=%v", useConfig), func(t *testing.T) {
			testFunc(t, useConfig)
		})
	}
}
