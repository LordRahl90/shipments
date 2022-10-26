package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringification(t *testing.T) {
	table := []struct {
		name string
		item Weight
		exp  string
	}{
		{
			name: "small",
			item: WeightSmall,
			exp:  "small",
		},
		{
			name: "medium",
			item: WeightMedium,
			exp:  "medium",
		},
		{
			name: "large",
			item: WeightLarge,
			exp:  "large",
		},
		{
			name: "huge",
			item: WeightHuge,
			exp:  "huge",
		},
		{
			name: "Unknown",
			item: WeightUnknown,
			exp:  "unknown",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.item.String()
			require.Equal(t, tt.exp, got)
		})
	}
}

func TestGetItemFromString(t *testing.T) {
	table := []struct {
		name string
		item string
		exp  Weight
	}{
		{
			name: "Small",
			item: "small",
			exp:  WeightSmall,
		},
		{
			name: "Medium",
			item: "medium",
			exp:  WeightMedium,
		},
		{
			name: "large",
			item: "large",
			exp:  WeightLarge,
		},
		{
			name: "Unknown",
			item: "Unknown",
			exp:  WeightUnknown,
		},
		{
			name: "huge",
			item: "huge",
			exp:  WeightHuge,
		},
		{
			name: "Hummer Jeep",
			item: "Hummer Jeep",
			exp:  WeightUnknown,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := WeightFromString(tt.item)
			require.Equal(t, tt.exp, got)
		})
	}
}

func TestWeightFromSize(t *testing.T) {
	table := []struct {
		name string
		size float64
		exp  Weight
	}{
		{
			name: "0kg",
			size: 0,
			exp:  WeightSmall,
		},
		{
			name: "5kg",
			size: 5,
			exp:  WeightSmall,
		},
		{
			name: "10kg",
			size: 10,
			exp:  WeightSmall,
		},
		{
			name: "11kg",
			size: 11,
			exp:  WeightMedium,
		},
		{
			name: "20kg",
			size: 20,
			exp:  WeightMedium,
		},
		{
			name: "25kg",
			size: 25,
			exp:  WeightMedium,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := WeightFromSize(tt.size)
			require.Equal(t, tt.exp, got)
		})
	}
}

func TestPriceFromSize(t *testing.T) {
	table := []struct {
		name                 string
		size                 float64
		origin, dest, errMsg string
		expErr               bool
		exp                  float64
	}{
		{
			name:   "0-local",
			size:   0,
			origin: "se",
			dest:   "se",
			expErr: false,
			exp:    100.0,
		},
		{
			name:   "0-eu",
			size:   0,
			origin: "se",
			dest:   "dk",
			exp:    150.0,
		},
		{
			name:   "0-one-outside-eu",
			size:   0,
			origin: "dk",
			dest:   "ng",
			exp:    250.0,
		},
		{
			name:   "0-both-outside-eu",
			size:   0,
			origin: "gh",
			dest:   "ng",
			exp:    250.0,
		},
		{
			name:   "0-same-outside-eu",
			size:   0,
			origin: "ng",
			dest:   "ng",
			exp:    250.0,
		},
		{
			name:   "undefined-size",
			size:   5000,
			origin: "ng",
			dest:   "ng",
			exp:    0.0,
			expErr: true,
			errMsg: "weight price not defined unknown",
		},
		{
			name:   "sample-test",
			size:   45,
			origin: "us",
			dest:   "se",
			exp:    1250.0,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PriceFromSize(tt.size, tt.origin, tt.dest)
			if tt.expErr {
				require.EqualError(t, err, tt.errMsg)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.exp, got)

		})
	}
}

func TestPriceFromWeight(t *testing.T) {
	table := []struct {
		name                 string
		weight               Weight
		origin, dest, errMsg string
		expErr               bool
		exp                  float64
	}{
		{
			name:   "0-local",
			weight: WeightSmall,
			origin: "se",
			dest:   "se",
			exp:    100.0,
		},
		{
			name:   "0-eu",
			weight: WeightSmall,
			origin: "se",
			dest:   "dk",
			exp:    150.0,
		},
		{
			name:   "0-one-outside-eu",
			weight: WeightSmall,
			origin: "dk",
			dest:   "ng",
			exp:    250.0,
		},
		{
			name:   "0-both-outside-eu",
			weight: WeightSmall,
			origin: "gh",
			dest:   "ng",
			exp:    250.0,
		},
		{
			name:   "0-same-outside-eu",
			weight: WeightSmall,
			origin: "ng",
			dest:   "ng",
			exp:    250.0,
		},
		{
			name:   "25-same-outside-eu",
			weight: WeightMedium,
			origin: "ng",
			dest:   "ng",
			exp:    750.0,
		},
		{
			name:   "0-unknown-weight",
			weight: WeightUnknown,
			origin: "ng",
			dest:   "ng",
			expErr: true,
			errMsg: "weight price not defined unknown",
			exp:    0.0,
		},
		{
			name:   "0-invalid-country",
			weight: WeightSmall,
			origin: "10",
			dest:   "ng",
			expErr: true,
			exp:    0.0,
			errMsg: "gountries error. Could not find country with code %s: 10",
		},
		{
			name:   "sample-test",
			weight: WeightLarge,
			origin: "us",
			dest:   "se",
			exp:    1250.0,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PriceFromWeight(tt.weight, tt.origin, tt.dest)
			if tt.expErr {
				require.EqualError(t, err, tt.errMsg)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.exp, got)

		})
	}
}
