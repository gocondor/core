package env

import (
	"os"
	"testing"
)

func TestGetVar(t *testing.T) {
	os.Setenv("testKey11", "testVal")
	v := GetVar("testKey11")
	if v != "testVal" {
		t.Error("failed testing get var")
	}
}

func TestGetVarOtherwiseDefault(t *testing.T) {
	v := GetVarOtherwiseDefault("testKey12", "defaultVal")
	if v != "defaultVal" {
		t.Error("failed testing get default")
	}
	os.Setenv("testKey12", "testVal")
	v = GetVarOtherwiseDefault("testKey12", "defaultVal")
	if v != "testVal" {
		t.Error("failed testing get default val")
	}
}

func TestIsSet(t *testing.T) {
	i := IsSet("testKey13")
	if i == true {
		t.Error("failed testing is set")
	}
	os.Setenv("testKey13", "testVal")
	i = IsSet("testKey13")
	if i == false {
		t.Error("filed testing is set")
	}
}

func TestSetEnvVars(t *testing.T) {
	envVars := map[string]string{
		"key14": "testVal14",
		"key15": "testVal15",
	}

	SetEnvVars(envVars)
	if GetVar("key14") != "testVal14" {
		t.Error("failed testing set vars")
	}
	if GetVar("key15") != "testVal15" {
		t.Error("failed testing set vars")
	}
}
