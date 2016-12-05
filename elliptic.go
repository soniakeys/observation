// Public domain by author, Sonia Keys

package observation

import (
	"math"

	"github.com/soniakeys/astro"
	"github.com/soniakeys/unit"
)

// Formulas from Meeus Chapter 33, Elliptic Motion.

// AstrometricJ2000 computes astrometric coordinates of an observation.
//
// Argument jde is the time of observation, arguments sun and target are
// functions that return geocentric J2000 equatorial rectangular coodinates
// of the sun and the target object respectively.
//
// Results are J2000 right ascension, declination of target as observed from
// site at the given jde.  Also returns phase angle β, sun-object radius r and
// observer-object distance Δ.  Returned angles are in radians, distances in AU.
func AstrometricJ2000(jde float64, sun, target func(jde float64) (x, y, z, r float64)) (α unit.RA, δ, ψ, β unit.Angle, r, Δ float64) {
	// AstrometricJ2000 derived from similar function in meeus/elliptic,
	// modified here to take two functions and return additional
	// values β, r and Δ.
	// values phase angle β, sun-object radius r and observer-object distance Δ.
	//
	X, Y, Z, R := sun(jde)
	x, y, z, r := target(jde)
	// (33.10) p. 229
	ξ := X + x
	η := Y + y
	ζ := Z + z
	Δ = math.Sqrt(ξ*ξ + η*η + ζ*ζ)
	{
		// unit math to convert a distance Δ in AU to light time in τ days:
		// Δ AU   astro.AU  m           s         day
		// ---- . ----------- . --------- . ---------
		//                 AU   astro.C m   86400  s
		// units cancel to give the assignment statement below.
		τ := Δ * astro.AU / astro.C / 86400
		x, y, z, r = target(jde - τ)
		ξ = X + x
		η = Y + y
		ζ = Z + z
		Δ = math.Sqrt(ξ*ξ + η*η + ζ*ζ)
	}
	α = unit.RAFromRad(math.Atan2(η, ξ))
	δ = unit.Angle(math.Asin(ζ / Δ))
	ψ = unit.Angle(math.Acos((ξ*X + η*Y + ζ*Z) / R / Δ))
	β = unit.Angle(math.Acos((ξ*x + η*y + ζ*z) / r / Δ))
	return
}

func Vmag(H, G float64, β unit.Angle, r, Δ float64) float64 {
	if math.IsNaN(H) {
		return H
	}
	if math.IsNaN(G) {
		G = .15
	}
	tanβ2 := (β / 2).Tan()
	Φ1 := math.Exp(-3.33 * math.Pow(tanβ2, .63))
	Φ2 := math.Exp(-1.87 * math.Pow(tanβ2, 1.22))
	return H + 5*math.Log10(r*Δ) - 2.5*math.Log10((1-G)*Φ1+G*Φ2)
}
