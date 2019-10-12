package flaggio_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/victorkohl/flaggio/internal/flaggio"
)

func TestList_DefaultWhenOn(t *testing.T) {
	tt := []struct {
		desc            string
		list            flaggio.VariantList
		expectedVariant *flaggio.Variant
	}{
		{
			desc: "returns DefaultWhenOn variant",
			list: flaggio.VariantList{
				{ID: "1", DefaultWhenOn: false},
				{ID: "2", DefaultWhenOn: true},
				{ID: "3", DefaultWhenOn: false},
			},
			expectedVariant: &flaggio.Variant{ID: "2", DefaultWhenOn: true},
		},
		{
			desc: "returns first variant when multiple DefaultWhenOn",
			list: flaggio.VariantList{
				{ID: "1", DefaultWhenOn: false},
				{ID: "2", DefaultWhenOn: true},
				{ID: "3", DefaultWhenOn: true},
			},
			expectedVariant: &flaggio.Variant{ID: "2", DefaultWhenOn: true},
		},
		{
			desc: "returns nil when no DefaultWhenOn",
			list: flaggio.VariantList{
				{ID: "1", DefaultWhenOn: false},
				{ID: "2", DefaultWhenOn: false},
				{ID: "3", DefaultWhenOn: false},
			},
			expectedVariant: nil,
		},
	}

	for _, test := range tt {
		t.Run(test.desc, func(t *testing.T) {
			vrnt := test.list.DefaultWhenOn()
			assert.Equal(t, test.expectedVariant, vrnt)
		})
	}
}

func TestList_DefaultWhenOff(t *testing.T) {
	tt := []struct {
		desc            string
		list            flaggio.VariantList
		expectedVariant *flaggio.Variant
	}{
		{
			desc: "returns DefaultWhenOff variant",
			list: flaggio.VariantList{
				{ID: "1", DefaultWhenOff: false},
				{ID: "2", DefaultWhenOff: true},
				{ID: "3", DefaultWhenOff: false},
			},
			expectedVariant: &flaggio.Variant{ID: "2", DefaultWhenOff: true},
		},
		{
			desc: "returns first variant when multiple DefaultWhenOff",
			list: flaggio.VariantList{
				{ID: "1", DefaultWhenOff: false},
				{ID: "2", DefaultWhenOff: true},
				{ID: "3", DefaultWhenOff: true},
			},
			expectedVariant: &flaggio.Variant{ID: "2", DefaultWhenOff: true},
		},
		{
			desc: "returns nil when no DefaultWhenOff",
			list: flaggio.VariantList{
				{ID: "1", DefaultWhenOff: false},
				{ID: "2", DefaultWhenOff: false},
				{ID: "3", DefaultWhenOff: false},
			},
			expectedVariant: nil,
		},
	}

	for _, test := range tt {
		t.Run(test.desc, func(t *testing.T) {
			vrnt := test.list.DefaultWhenOff()
			assert.Equal(t, test.expectedVariant, vrnt)
		})
	}
}
