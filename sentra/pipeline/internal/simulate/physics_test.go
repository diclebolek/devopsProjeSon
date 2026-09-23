package simulate

import (
	"testing"
	"time"
)

func TestRegimesAreSeparable(t *testing.T) {
	capturedAt := time.UnixMilli(1_710_000_000_000).UTC()
	normal, normalLabel, err := Synthesize(RegimeNormal, 1, 0, capturedAt, "motor-lab-01", "0.1.0")
	if err != nil {
		t.Fatal(err)
	}
	overload, _, err := Synthesize(RegimeOverload, 1, 0, capturedAt, "motor-lab-01", "0.1.0")
	if err != nil {
		t.Fatal(err)
	}
	unbalanced, _, err := Synthesize(RegimeUnbalanced, 1, 0, capturedAt, "motor-lab-01", "0.1.0")
	if err != nil {
		t.Fatal(err)
	}
	rough, _, err := Synthesize(RegimeHighVibration, 1, 0, capturedAt, "motor-lab-01", "0.1.0")
	if err != nil {
		t.Fatal(err)
	}

	if normalLabel.Condition != "NORMAL" {
		t.Fatalf("label %s", normalLabel.Condition)
	}
	if overload.CurrentA < normal.CurrentA*2 {
		t.Fatalf("overload current %v normal %v", overload.CurrentA, normal.CurrentA)
	}
	if overload.TemperatureC < normal.TemperatureC+20 {
		t.Fatalf("overload temperature %v", overload.TemperatureC)
	}
	if unbalanced.AccelRMS < normal.AccelRMS*3 {
		t.Fatalf("unbalanced rms %v normal %v", unbalanced.AccelRMS, normal.AccelRMS)
	}
	if rough.Kurtosis < normal.Kurtosis+1 {
		t.Fatalf("high vibration kurtosis %v normal %v", rough.Kurtosis, normal.Kurtosis)
	}
	if normal.DominantFreqHz < 20 || normal.DominantFreqHz > 30 {
		t.Fatalf("normal dominant frequency %v", normal.DominantFreqHz)
	}
}

func TestCycleRotatesEveryPeriod(t *testing.T) {
	if got := Resolve(RegimeCycle, 0, 2); got != RegimeNormal {
		t.Fatal(got)
	}
	if got := Resolve(RegimeCycle, 2, 2); got != RegimeOverload {
		t.Fatal(got)
	}
	if got := Resolve(RegimeOverload, 9, 2); got != RegimeOverload {
		t.Fatal(got)
	}
}
