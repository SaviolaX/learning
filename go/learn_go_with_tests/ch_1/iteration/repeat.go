package iteration

import "strings"

const repeatCount = 5

func Repeat(character string, repeats int) string {
    var repeated strings.Builder
    var ni int
    if repeats == 0 {
        ni = repeatCount
    } else {
        ni = repeats
    }
    for i := 0; i < ni; i++ {
        repeated.WriteString(character) 
    }
   return repeated.String() 
}
