package utils

import (
	"math"
)

func Round(val float64, prec int) float64 { 

    var rounder float64 
    intermed := val * math.Pow(10, float64(prec)) 

    if val >= 0.5 { 
        rounder = math.Ceil(intermed) 
    } else { 
        rounder = math.Floor(intermed) 
    } 

    return rounder / math.Pow(10, float64(prec)) 

} 	