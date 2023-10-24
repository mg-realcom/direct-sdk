package extension

type StatusSelection string

const (
	DRAFT      StatusSelection = "DRAFT"      // Дополнение не отправлялось на модерацию.
	MODERATION StatusSelection = "MODERATION" // Дополнение находится на модерации.
	ACCEPTED   StatusSelection = "ACCEPTED"   // Дополнение принято модерацией. Объявление при показе будет содержать дополнение.
	REJECTED   StatusSelection = "REJECTED"   // Дополнение отклонено модерацией.
	UNKNOWN    StatusSelection = "UNKNOWN"    // Неизвестный статус. Используется для обеспечения обратной совместимости и отображения статусов, не поддерживаемых в данной версии API.
)
