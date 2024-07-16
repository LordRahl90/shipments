package core

import (
	"context"
	"fmt"

	"shipments/domains/tracing"

	"github.com/pariz/gountries"
	"go.opentelemetry.io/otel/attribute"
)

// Weight special type created for items
type Weight int

const (
	// WeightUnknown unknown item
	WeightUnknown Weight = iota
	// WeightSmall small weight
	WeightSmall
	// WeightMedium medium weight
	WeightMedium
	// WeightLarge large weight
	WeightLarge
	// WeightHuge huge weight
	WeightHuge

	local = "se"
)

var (
	// WeightPrices keeps track of the prices for each weight classs.
	// this is not a hot data and could be kept in the db/cache as the company grows
	WeightPrices = map[Weight]float64{
		WeightSmall:  100,
		WeightMedium: 300,
		WeightLarge:  500,
		WeightHuge:   2000,
	}

	query = gountries.New()
)

// String returns the stringified version of the weight
func (i Weight) String() string {
	switch i {
	case WeightSmall:
		return "small"
	case WeightMedium:
		return "medium"
	case WeightLarge:
		return "large"
	case WeightHuge:
		return "huge"
	default:
		return "unknown"
	}
}

// WeightFromString returns a Weight from the equivalent string representation
func WeightFromString(item string) Weight {
	switch item {
	case WeightSmall.String():
		return WeightSmall
	case WeightMedium.String():
		return WeightMedium
	case WeightLarge.String():
		return WeightLarge
	case WeightHuge.String():
		return WeightHuge
	default:
		return WeightUnknown
	}
}

// WeightFromSize return a weight from a given size
func WeightFromSize(size float64) Weight {
	switch {
	case size >= 0 && size <= 10:
		return WeightSmall
	case size > 10 && size <= 25:
		return WeightMedium
	case size > 25 && size <= 50:
		return WeightLarge
	case size > 50 && size <= 1000:
		return WeightHuge
	default:
		return WeightUnknown
	}
}

// Multiplier defines the price multiplier based on the destination country
func Multiplier(ctx context.Context, origin, destination string) (float64, error) {
	_, span := tracing.Tracer().Start(ctx, "core:multiplier")
	span.SetAttributes(attribute.KeyValue{Key: "origin", Value: attribute.StringValue(origin)})
	span.SetAttributes(attribute.KeyValue{Key: "destination", Value: attribute.StringValue(destination)})
	defer span.End()

	if origin == local && destination == local {
		return 1.0, nil
	}
	originInEU, err := isInEU(origin)
	if err != nil {
		return 0, err
	}
	destInEU, err := isInEU(destination)
	if err != nil {
		return 0, err
	}
	// both countries are in the EU
	if originInEU && destInEU {
		return 1.5, nil
	}
	// if it's only one that's in the EU and the other is not
	// then 2.5 applies
	return 2.5, nil
}

func isInEU(countryCode string) (bool, error) {
	country, err := query.FindCountryByAlpha(countryCode)
	if err != nil {
		return false, err
	}
	return country.EuMember, nil
}

// PriceFromSize returns the total price based on the size and destination
func PriceFromSize(ctx context.Context, size float64, origin, destination string) (float64, error) {
	ctx, span := tracing.Tracer().Start(ctx, "core:price-from-size")
	span.SetAttributes(attribute.KeyValue{Key: "weight", Value: attribute.StringValue(fmt.Sprintf("%f", size))})
	span.SetAttributes(attribute.KeyValue{Key: "origin", Value: attribute.StringValue(origin)})
	span.SetAttributes(attribute.KeyValue{Key: "destination", Value: attribute.StringValue(destination)})
	defer span.End()
	w := WeightFromSize(size)
	p, ok := WeightPrices[w]
	if !ok {
		return 0, fmt.Errorf("weight price not defined %s", w.String())
	}
	m, err := Multiplier(ctx, origin, destination)
	if err != nil {
		return 0, err
	}
	total := p * m
	return total, nil
}

// PriceFromWeight return the price from a given weight and destination country.
func PriceFromWeight(ctx context.Context, w Weight, origin, destination string) (float64, error) {
	ctx, span := tracing.Tracer().Start(ctx, "core:price-from-weight")
	span.SetAttributes(attribute.KeyValue{Key: "weight", Value: attribute.StringValue(w.String())})
	span.SetAttributes(attribute.KeyValue{Key: "origin", Value: attribute.StringValue(origin)})
	span.SetAttributes(attribute.KeyValue{Key: "destination", Value: attribute.StringValue(destination)})
	defer span.End()

	p, ok := WeightPrices[w]
	if !ok {
		return 0, fmt.Errorf("weight price not defined %s", w.String())
	}
	m, err := Multiplier(ctx, origin, destination)
	if err != nil {
		return 0, err
	}
	total := p * m
	return total, nil
}
