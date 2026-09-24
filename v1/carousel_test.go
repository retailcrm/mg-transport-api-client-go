package v1

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProductCarouselTemplateContracts(t *testing.T) {
	const carousel = `{"cards":[
		{"header":{"type":"product"},"buttons":{"items":[{"type":"spm"}]}},
		{"header":{"type":"product"},"buttons":{"items":[{"type":"spm"}]}}
	]}`
	input := []byte(`{"channel_id":123,"name":"product_carousel","lang":"en_US","category":"marketing","body":"See our products","type":"carousel","carousel":` + carousel + `}`)

	for _, target := range []any{
		&TemplateCreateWebhookData{}, &TemplateUpdateWebhookData{},
		&Template{}, &ActivateTemplateRequest{}, &UpdateTemplateRequest{},
	} {
		require.NoError(t, json.Unmarshal(input, target))
		data, err := json.Marshal(target)
		require.NoError(t, err)
		var fields map[string]json.RawMessage
		require.NoError(t, json.Unmarshal(data, &fields))
		require.JSONEq(t, carousel, string(fields["carousel"]))
	}

	var created TemplateCreateWebhookData
	require.NoError(t, json.Unmarshal(input, &created))
	require.Len(t, created.Carousel.Cards, 2)
	require.Equal(t, TemplateCarouselCardProduct, created.Carousel.Cards[0].Header.Type)
	require.IsType(t, &SPMButton{}, created.Carousel.Cards[0].Buttons.Items[0])
}

func TestProductCarouselMessageContracts(t *testing.T) {
	const carousel = `{"cards":[
		{"header":{"type":"product","product_retailer_id":"sku-1"},"buttons":[{"type":"url","title":"Open","args":["first"]}]},
		{"header":{"type":"product","product_retailer_id":"sku-2"}}
	]}`
	input := []byte(`{"code":"product_carousel","category":"marketing","carousel":` + carousel + `}`)
	var info TemplateInfo
	require.NoError(t, json.Unmarshal(input, &info))
	require.Len(t, info.Carousel.Cards, 2)
	require.Equal(t, TemplateCarouselCardProduct, info.Carousel.Cards[0].Header.Type)
	require.Equal(t, "sku-1", info.Carousel.Cards[0].Header.ProductRetailerID)
	require.Nil(t, info.Carousel.Cards[0].Body)
	require.Equal(t, "first", info.Carousel.Cards[0].Buttons[0].Args[0])

	data, err := json.Marshal(info)
	require.NoError(t, err)
	var fields map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(data, &fields))
	require.JSONEq(t, carousel, string(fields["carousel"]))
}
