package rpc

import (
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// RuleItemsToProto converts model rule items to proto VenueDescriptionItem messages.
func RuleItemsToProto(items []models.RuleItem) []*apiv1.VenueDescriptionItem {
	result := make([]*apiv1.VenueDescriptionItem, 0, len(items))

	for _, item := range items {
		itemType := item.ItemType
		pbType, _ := itemType.ProtoEnum()

		pbItem := &apiv1.VenueDescriptionItem{
			Content:  item.Content,
			ItemType: pbType,
		}

		for _, aud := range item.Audiences {
			pe, _ := aud.ProtoEnum()
			pbItem.Audiences = append(pbItem.Audiences, pe)
		}

		result = append(result, pbItem)
	}

	return result
}

// ProtoToRuleItems converts proto VenueDescriptionItem messages to model rule items.
func ProtoToRuleItems(items []*apiv1.VenueDescriptionItem) []models.RuleItem {
	result := make([]models.RuleItem, 0, len(items))

	for i, item := range items {
		var itemType models.VenueDescriptionItemType

		_ = itemType.FromProtoEnum(item.GetItemType())

		var audiences []models.ScoringCategory

		for _, aud := range item.GetAudiences() {
			var sc models.ScoringCategory

			_ = sc.FromProtoEnum(aud)

			audiences = append(audiences, sc)
		}

		result = append(result, models.RuleItem{
			Content:   item.GetContent(),
			ItemType:  itemType,
			Audiences: audiences,
			Rank:      uint32(i),
		})
	}

	return result
}
