// Public domain.

// Package observation represents ground-based and space-based observations
// of moving objects against the sky.
package observation

import (
	"math"

	"github.com/soniakeys/astro"
	"github.com/soniakeys/coord"
)

// VMeas represents a "visual" measurment in units convenient for computations.
type VMeas struct {
	MJD        float64 // time of observation
	coord.Sphr         // components in radians
	// The apparent magnitude is represented as somehow normalized to "V."
	// The actual observed magnitude band is not represented here.
	VMag float64
	// Quality identifier.  Typically simply the 3 character MPC obscode,
	// but can be any string to associate the measurement with a quality
	// level.  This identifier is associated with the obserr keyword
	// in the digest2 config file, for example.
	Qual string
}

// VObs is a common interface for ground-based and spaced-based observations
type VObs interface {
	// underlying measurement--the actual observation
	Meas() *VMeas
	// location of observer relative to the center of the earth.
	// units are AU.
	EarthObserverVect() coord.Cart
}

// SiteObs represents an observation from a fixed ground-based observatory.
// It satisfies the VObs interface.
type SiteObs struct {
	VMeas // the observation
	// Parallax constants determine the observer location relative to the
	// center of the Earth.
	Par *ParallaxConst
}

// Meas satisfies a method of the VObs interface.
func (o *SiteObs) Meas() *VMeas {
	return &o.VMeas
}

// EarthObserverVect satisfies a method of the VObs interface.
func (o *SiteObs) EarthObserverVect() coord.Cart {
	return EarthObserverVect(o.MJD, o.Par)
}

// EarthObserverVect computes the geocenter-observer vector.
//
// Result units are AU.
func EarthObserverVect(mjd float64, p *ParallaxConst) coord.Cart {
	sth, cth := math.Sincos(astro.Lst(mjd, p.Longitude))
	return coord.Cart{
		X: p.RhoCosPhi * cth,
		Y: p.RhoCosPhi * sth,
		Z: p.RhoSinPhi,
	}
}

// SatObs represents an observation from an observatory in Earth orbit.
// It satisfies the VObs interface.
type SatObs struct {
	// Sat is typically the 3 character MPC obscode, but is not restricted
	// to these.
	Sat    string
	VMeas             // the observation
	Offset coord.Cart // offset from center of Earth, in AU
}

// Meas satisfies a method of the VObs interface.
func (o *SatObs) Meas() *VMeas { // Implement VObs
	return &o.VMeas
}

// EarthObserverVect satisfies a method of the VObs interface.
func (o *SatObs) EarthObserverVect() (c coord.Cart) {
	// to do:  use Offset.  need to verify orientation and such.
	// until then, geocentric is safer.
	return
}

// ParallaxConst represents a vector from the center of the Earth.
type ParallaxConst struct {
	Longitude float64 // unit is circles
	RhoCosPhi float64 // unit is AU
	RhoSinPhi float64 // unit is AU
}

// ParallaxMap is a mapping from strings to parallax constants.
// The map key is typically the 3 character MPC obscode, but
// is not restricted to these and can be anything convenient.
type ParallaxMap map[string]*ParallaxConst

// Arc, a sequence of observations of the same object.
type Arc struct {
	Desig string
	Obs   []VObs
}
