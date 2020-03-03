package envflag

import (
	"os"
	"testing"
	"time"
)

// Test parsing a environment variables
func TestParseEnv(t *testing.T) {
	os.Setenv("BOOL", "true")
	os.Setenv("INT", "22")
	os.Setenv("INT64", "0x23")
	os.Setenv("UINT", "24")
	os.Setenv("UINT64", "25")
	os.Setenv("STRING", "hello")
	os.Setenv("FLOAT64", "2718e28")
	os.Setenv("DURATION", "2m")

	f := NewFlagSet(os.Args[0], ContinueOnError)

	boolFlag := f.Bool("bool", false, "bool value")
	intFlag := f.Int("int", 0, "int value")
	int64Flag := f.Int64("int64", 0, "int64 value")
	uintFlag := f.Uint("uint", 0, "uint value")
	uint64Flag := f.Uint64("uint64", 0, "uint64 value")
	stringFlag := f.String("string", "0", "string value")
	float64Flag := f.Float64("float64", 0, "float64 value")
	durationFlag := f.Duration("duration", 5*time.Second, "time.Duration value")

	err := f.parseEnv(os.Environ())
	if err != nil {
		t.Fatal("expected no error; got ", err)
	}
	if *boolFlag != true {
		t.Error("bool flag should be true, is ", *boolFlag)
	}
	if *intFlag != 22 {
		t.Error("int flag should be 22, is ", *intFlag)
	}
	if *int64Flag != 0x23 {
		t.Error("int64 flag should be 0x23, is ", *int64Flag)
	}
	if *uintFlag != 24 {
		t.Error("uint flag should be 24, is ", *uintFlag)
	}
	if *uint64Flag != 25 {
		t.Error("uint64 flag should be 25, is ", *uint64Flag)
	}
	if *stringFlag != "hello" {
		t.Error("string flag should be `hello`, is ", *stringFlag)
	}
	if *float64Flag != 2718e28 {
		t.Error("float64 flag should be 2718e28, is ", *float64Flag)
	}
	if *durationFlag != 2*time.Minute {
		t.Error("duration flag should be 2m, is ", *durationFlag)
	}
}

func TestFlagSetParseErrors(t *testing.T) {
	fs := NewFlagSet("test", ContinueOnError)
	fs.Int("int", 0, "int value")

	args := []string{"-int", "bad"}
	expected := `invalid value "bad" for flag -int: parse error`
	if err := fs.Parse(args); err == nil || err.Error() != expected {
		t.Errorf("expected error %q parsing from args, got: %v", expected, err)
	}

	if err := os.Setenv("INT", "bad"); err != nil {
		t.Fatalf("error setting env: %s", err.Error())
	}
	expected = `invalid value "bad" for environment variable int: parse error`
	if err := fs.Parse([]string{}); err == nil || err.Error() != expected {
		t.Errorf("expected error %q parsing from env, got: %v", expected, err)
	}
}
