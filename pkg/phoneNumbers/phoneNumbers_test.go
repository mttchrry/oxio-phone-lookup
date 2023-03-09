package phoneNumbers

import (
	"context"
	"testing"
)

func TestPhoneNumbers_SimplifyString_NotValid(t *testing.T) {
	p := New()
	ctx := context.Background()

	_, err := p.simplifyString(ctx, "")
	if err == nil {
		t.Error("Expected empty string to be invalid")
	}

	_, err = p.simplifyString(ctx, "1  4145551122")
	if err == nil {
		t.Error("Expected double space to be invalid")
	}

	_, err = p.simplifyString(ctx, "141-45551122")
	if err == nil {
		t.Error("Expected misplaced dash to be invalid")
	}

	_, err = p.simplifyString(ctx, "43121 4145551122")
	if err == nil {
		t.Error("Expect long country code to be invalid")
	}

	_, err = p.simplifyString(ctx, "+121 4145 5511 22")
	if err == nil {
		t.Error("Expect non-standard spacing to be invalid")
	}

	_, err = p.simplifyString(ctx, "+121 4145 5511 22")
	if err == nil {
		t.Error("Expect space and parenthesis to be invalid")
	}
}

func TestPhoneNumbers_SimplifyString_Valid(t *testing.T) {
	p := New()
	ctx := context.Background()

	expectedSimpleNumber := "14145551122"
	str, err := p.simplifyString(ctx, "+14145551122")
	if err != nil {
		t.Errorf("Unexpected err: %v", err)
	}
	if str != expectedSimpleNumber {
		t.Errorf("expecting %v, got %v", expectedSimpleNumber, str)
	}

	str, err = p.simplifyString(ctx, "1 414 555 1122")
	if err != nil {
		t.Errorf("Unexpected err: %v", err)
	}
	if str != expectedSimpleNumber {
		t.Errorf("expecting %v, got %v", expectedSimpleNumber, str)
	}

	str, err = p.simplifyString(ctx, "+1(414)555-1122")
	if err != nil {
		t.Errorf("Unexpected err: %v", err)
	}
	if str != expectedSimpleNumber {
		t.Errorf("expecting %v, got %v", expectedSimpleNumber, str)
	}

	str, err = p.simplifyString(ctx, "+1(414)555-1122")
	if err != nil {
		t.Errorf("Unexpected err: %v", err)
	}
	if str != expectedSimpleNumber {
		t.Errorf("expecting %v, got %v", expectedSimpleNumber, str)
	}

	longNumber := "12344145551122"
	str, err = p.simplifyString(ctx, "1234(414)555-1122")
	if err != nil {
		t.Errorf("Unexpected err: %v", err)
	}
	if str != longNumber {
		t.Errorf("expecting %v, got %v", longNumber, str)
	}

	tenDigitNumber := "4145551122"
	str, err = p.simplifyString(ctx, "(414)555 1122")
	if err != nil {
		t.Errorf("Unexpected err: %v", err)
	}
	if str != tenDigitNumber {
		t.Errorf("expecting %v, got %v", tenDigitNumber, str)
	}

}