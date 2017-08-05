package userValidator

import (
	"testing"
)

func TestValideUsername(t *testing.T) {
	tests := []struct {
		Username string
		Expected bool
	}{
		{
			"lol",
			true,
		},
		{
			"àçésdsd",
			true,
		},
		{
			"[}\\_",
			false,
		},
		{
			"^!",
			false,
		},
		{
			"test ",
			false,
		},
		{
			"test done",
			false,
		},
		{
			"れんちょん",
			true,
		},
		{
			"",
			false,
		},
	}

	for _, val := range tests {
		testVal := ValidateUsername(val.Username)
		if testVal != val.Expected {
			t.Errorf("The test returned a result %t instead of %t for '%s'", testVal, val.Expected, val.Username)
		}
	}
}

func TestEmailValidation(t *testing.T) {
	tests := [][]interface{}{
		{
			"lol@xd.uu",
			true,
		},
		{
			"someone.might.wann@this.mail",
			true,
		},
		{
			"disposable.style.email.with+symbol@example.com",
			true,
		},
		{
			"other.email-with-dash@example.com",
			true,
		},
		{
			"x@example.com",
			true,
		},
		{
			`"very.unusual.@.unusual.com"@example.com`,
			true,
		},
		{
			`very.(),:;<>[]\".VERY.\"very@\\ \"very\".unusual"@strange.example.com`,
			false,
		},
		{
			"example-indeed@strange-example.com",
			true,
		},
		{
			"admin@mailserver1",
			false,
		},
		{
			"#!$%&'*+-/=?^_`{}|~@example.org",
			true,
		},
		{
			`"()<>[]:,;@\\\"!#$%&'-/=?^_` + "`" + `{}| ~.a"@example.org`,
			true,
		},
		{
			`" "@example.org`,
			true,
		},
		{
			`example@s.solutions`,
			true,
		},
		{
			`Abc.example.com`,
			false,
		},
		{
			`A@b@c@example.com`,
			false,
		},
		{
			`a"b(c)d,e:f;g<h>i[j\k]l@example.com`,
			false,
		},
		{
			`just"not"right@example.com`,
			false,
		},
		{
			`this is"not\allowed@example.com`,
			false,
		},
		{
			`to\ to@xn-fsqu00a.xn-0zwm56d`,
			true,
		},
		{
			`1234567890123456789012345678901234567890123456789012345678901234+2323232@example.com`,
			true,
		},
		{
			`john..doe@example.com`,
			false,
		},
		{
			`john.doe@example..com`,
			false,
		},
		{
			`toto@[192.168.1.1]`,
			true,
		},
	}
	for _, val := range tests {
		testVal := EmailValidation(val[0].(string))
		if testVal != val[1] {
			t.Errorf("The test returned a result %t instead of %t for %s", testVal, val[1].(bool), val[0].(string))
		}
	}
}
