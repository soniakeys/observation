// Public domain by author, Sonia Keys

// Formulas from Meeus Chapter 33, Elliptic Motion.
package observation

// position cut and paste from meeus/elliptic, modified to return additional
// values.
//
// Position returns observed equatorial coordinates of a body with Keplerian elements.
//
// Argument e must be a valid V87Planet object for Earth.
//
// Results are right ascension and declination α and δ, and elongation ψ,
// all in radians.  In addition, returns phase angle β, the sun-object distance
// r in AU and the earth-object distance Δ in AU.
func (o *Orbit) PositionFromEarth(jde float64, e *pp.V87Planet) (α, δ, ψ, β, r, Δ float64) {
	return astrometricJ2000(o.Position, jde, e)
}

// An ObjPosFunc returns J2000 equatorial coordinates for an object
// at the given jde.
type objPosFunc func(jde float64) (x, y, z, r float64)

// astrometricJ2000 cut and paste from meeus/elliptic, modified to return
// additional values phase angle β and observer-object range Δ.
//
// AstrometricJ2000 is a utility function for computing astrometric coordinates.
//
// Argument f is a function that returns J2000 equatorial rectangular
// coodinates of a body.
//
// Results are J2000 right ascention, declination, and elongation.
func astrometricJ2000(object, earth) (α, δ, ψ, β, r, Δ float64) {
	X, Y, Z := solarxyz.PositionJ2000(e, jde)
	x, y, z, r := f(jde)
	// (33.10) p. 229
	ξ := X + x
	η := Y + y
	ζ := Z + z
	Δ = math.Sqrt(ξ*ξ + η*η + ζ*ζ)
	{
		τ := base.LightTime(Δ)
		x, y, z, r = f(jde - τ)
		ξ = X + x
		η = Y + y
		ζ = Z + z
		Δ = math.Sqrt(ξ*ξ + η*η + ζ*ζ)
	}
	α = math.Atan2(η, ξ)
	if α < 0 {
		α += 2 * math.Pi
	}
	δ = math.Asin(ζ / Δ)
	R := math.Sqrt(X*X + Y*Y + Z*Z)
	ψ = math.Acos((ξ*X + η*Y + ζ*Z) / R / Δ)
	β = math.Acos((ξ*x + η*y + ζ*z) / r / Δ)
	return
}

func Vmag(H, G, β, r, Δ float64) float64 {
	if math.IsNaN(H) {
		return H
	}
	if math.IsNaN(G) {
		G = .15
	}
	tanβ2 := math.Tan(β / 2)
	Φ1 := math.Exp(-3.33 * math.Pow(tanβ2, .63))
	Φ2 := math.Exp(-1.87 * math.Pow(tanβ2, 1.22))
	return H + 5*math.Log10(r*Δ) - 2.5*math.Log10((1-G)*Φ1+G*Φ2)
}
