package iteration

import "testing"

func TestRepeat(t *testing.T)  {
    t.Run("no number iterations passed", func(t *testing.T) {
        repeated := Repeat("a", 0)
        expected := "aaaaa"

        if repeated != expected {
            t.Errorf("expected %q but got %q", expected, repeated)
        }
    })
    t.Run("iterations number passed: 2", func(t *testing.T) {
        repeated := Repeat("a", 2)
        expected := "aa"

        if repeated != expected {
            t.Errorf("expected %q but got %q", expected, repeated)
        }
    })
   }

func BenchmarkRepeat(b *testing.B)  {
    for i := 0; i < b.N; i++ {
        Repeat("a", 0)
    } 
    
}
