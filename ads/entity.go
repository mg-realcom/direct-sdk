package ads

type AdState string

const (
	AdStateSuspended       AdState = "SUSPENDED"         // 	Показы объявления остановлены владельцем с помощью метода suspend или в веб-интерфейсе.
	AdStateOffByMonitoring AdState = "OFF_BY_MONITORING" // 	Показы объявления автоматически остановлены мониторингом доступности сайта.
	AdStateOn              AdState = "ON"                // 	Объявление активно, принадлежит к активной кампании и может быть показано (при наличии средств на кампании, в соответствии с настройками временного таргетинга и т. п.).
	AdStateOff             AdState = "OFF"               // 	Объявление неактивно (черновик, ожидает модерации, отклонено) или принадлежит к неактивной либо остановленной кампании.
	AdStateArchived        AdState = "ARCHIVED"          // 	Объявление помещено в архив (с помощью метода archive или пользователем в веб-интерфейсе) или принадлежит к архивной кампании.
)

type AdStatus string

const (
	AdStatusDraft       AdStatus = "DRAFT"       // Объявление создано и еще не отправлено на модерацию.
	AdStatusModeration  AdStatus = "MODERATION"  // 	Объявление находится на модерации.
	AdStatusPreaccepted AdStatus = "PREACCEPTED" // 	Объявление допущено к показам автоматически, но будет дополнительно проверено модератором.
	AdStatusAccepted    AdStatus = "ACCEPTED"    // 	Объявление принято модерацией.
	AdStatusRejected    AdStatus = "REJECTED"    // 	Объявление отклонено модерацией.
)

type AdGroupType string

const (
	AdGroupTypeTextAdGroup        AdGroupType = "TEXT_AD_GROUP"
	AdGroupTypeSmartAdGroup       AdGroupType = "SMART_AD_GROUP"
	AdGroupTypeMobileAppAdGroup   AdGroupType = "MOBILE_APP_AD_GROUP"
	AdGroupTypeDynamicTextAdGroup AdGroupType = "DYNAMIC_TEXT_AD_GROUP"
	AdGroupTypeCpmBannerAdGroup   AdGroupType = "CPM_BANNER_AD_GROUP"
	AdGroupTypeCpmVideoAdGroup    AdGroupType = "CPM_VIDEO_AD_GROUP"
)

type AdType string

const (
	AdTypeTextAd                       AdType = "TEXT_AD"
	AdTypeImageAd                      AdType = "IMAGE_AD"
	AdTypeTextImageAd                  AdType = "TEXT_IMAGE_AD"
	AdTypeTextAdBuilderAd              AdType = "TEXT_AD_BUILDER_AD"
	AdTypeCpcVideoAd                   AdType = "CPC_VIDEO_AD"
	AdTypeSmartAd                      AdType = "SMART_AD"
	AdTypeMobileAppAd                  AdType = "MOBILE_APP_AD"
	AdTypeMobileAppImageAd             AdType = "MOBILE_APP_IMAGE_AD"
	AdTypeMobileAppAdBuilderAd         AdType = "MOBILE_APP_AD_BUILDER_AD"
	AdTypeMobileAppCpcVideoAdBuilderAd AdType = "MOBILE_APP_CPC_VIDEO_AD_BUILDER_AD"
	AdTypeDynamicTextAd                AdType = "DYNAMIC_TEXT_AD"
	AdTypeCpmBannerAd                  AdType = "CPM_BANNER_AD"
	AdTypeCpmVideoAd                   AdType = "CPM_VIDEO_AD"
)
