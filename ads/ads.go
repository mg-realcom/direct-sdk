package ads

import (
	"github.com/mg-realcom/direct-sdk/common"
	"github.com/mg-realcom/direct-sdk/extension"
)

type Request struct {
	Method string `json:"method"`
	Params Params `json:"params"`
}

type Params struct {
	SelectionCriteria                      SelectionCriteria `json:"SelectionCriteria"`                      // Критерий отбора объявлений.
	FieldNames                             []string          `json:"FieldNames"`                             // Имена параметров верхнего уровня, которые требуется получить.
	TextAdFieldNames                       *[]string         `json:"TextAdFieldNames"`                       // Имена параметров текстово-графического объявления, которые требуется получить.
	TextAdPriceExtensionFieldNames         *[]string         `json:"TextAdPriceExtensionFieldNames"`         // Имена параметров цены товара или услуги в текстово-графическом объявлении, которые требуется получить.
	MobileAppAdFieldNames                  *[]string         `json:"MobileAppAdFieldNames"`                  // Имена параметров объявления для рекламы мобильных приложений, которые требуется получить.
	DynamicTextAdFieldNames                *[]string         `json:"DynamicTextAdFieldNames"`                // Имена параметров динамического объявления, которые требуется получить.
	TextImageAdFieldNames                  *[]string         `json:"TextImageAdFieldNames"`                  // Имена параметров графического объявления, созданного на основе изображения (в группе текстово-графических объявлений), которые требуется получить.
	MobileAppImageAdFieldNames             *[]string         `json:"MobileAppImageAdFieldNames"`             // Имена параметров графического объявления, созданного на основе изображения (в группе для рекламы мобильных приложений), которые требуется получить.
	TextAdBuilderAdFieldNames              *[]string         `json:"TextAdBuilderAdFieldNames"`              // Имена параметров графического объявления, созданного на основе креатива (в группе текстово-графических объявлений), которые требуется получить.
	MobileAppAdBuilderAdFieldNames         *[]string         `json:"MobileAppAdBuilderAdFieldNames"`         // Имена параметров графического объявления, созданного на основе креатива (в группе для рекламы мобильных приложений), которые требуется получить.
	MobileAppCpcVideoAdBuilderAdFieldNames *[]string         `json:"MobileAppCpcVideoAdBuilderAdFieldNames"` // Имена параметров видеообъявления, созданного на основе креатива (в группе для рекламы мобильных приложений), которые требуется получить.
	CpcVideoAdBuilderAdFieldNames          *[]string         `json:"CpcVideoAdBuilderAdFieldNames"`          // Имена параметров медийного баннера, которые требуется получить.
	CpmBannerAdBuilderAdFieldNames         *[]string         `json:"CpmBannerAdBuilderAdFieldNames"`         // Имена параметров медийного видеообъявления (в кампаниях с типом «Медийная кампания»), которые требуется получить.
	SmartAdBuilderAdFieldNames             *[]string         `json:"SmartAdBuilderAdFieldNames"`             // Имена параметров смарт-баннера, которые требуется получить.
	Page                                   *common.Page      `json:"Page"`                                   // Структура, задающая страницу при постраничной выборке данных.
}

type SelectionCriteria struct {
	Ids                         []int                        `json:"Ids"`                         //  Отбирать объявления с указанными идентификаторами. От 1 до 10 000 элементов в массиве.
	CampaignIDS                 []int                        `json:"CampaignIds"`                 // Отбирать объявления указанных групп. От 1 до 1000 элементов в массиве.
	AdGroupIDS                  []int                        `json:"AdGroupIds"`                  // Отбирать объявления указанных кампаний. От 1 до 10 элементов в массиве.
	States                      *[]AdState                   `json:"States"`                      // Отбирать объявления с указанными состояниями.
	Statuses                    *[]AdStatus                  `json:"Statuses"`                    // Отбирать объявления с указанными статусами.
	Types                       *[]AdType                    `json:"Types"`                       // Отбирать объявления с указанными типами.
	Mobile                      *bool                        `json:"Mobile"`                      // Отбирать объявления по признаку того, что объявление является мобильным:
	VCardIds                    *[]int                       `json:"VCardIds"`                    // Отбирать объявления с указанными визитками. От 1 до 50 элементов в массиве.
	SitelinkSetIds              *[]int                       `json:"SitelinkSetIds"`              // Отбирать объявления с указанными наборами быстрых ссылок. От 1 до 50 элементов в массиве.
	AdImageHashes               *[]string                    `json:"AdImageHashes"`               // Отбирать объявления с указанными изображениями. От 1 до 50 элементов в массиве.
	VCardModerationStatuses     *[]extension.StatusSelection `json:"VCardModerationStatuses"`     // Отбирать объявления по результату модерации визитки. Описание статусов
	SitelinksModerationStatuses *[]extension.StatusSelection `json:"SitelinksModerationStatuses"` // Отбирать объявления по результату модерации набора быстрых ссылок.
	AdImageModerationStatuses   *[]extension.StatusSelection `json:"AdImageModerationStatuses"`   // Отбирать объявления по результату модерации изображения.
	AdExtensionIds              *[]int                       `json:"AdExtensionIds"`              // Отбирать объявления с указанными расширениями. От 1 до 50 элементов в массиве.
}
