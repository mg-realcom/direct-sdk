package statistics

import "github.com/mg-realcom/direct-sdk/common"

type ReportDefinition struct {
	Selection         *SelectionCriteria  `json:"SelectionCriteria,omitempty"`
	Goals             *[]string           `json:"Goals,omitempty"`
	AttributionModels *[]AttributionModel `json:"AttributionModels,omitempty"`
	FieldNames        []string            `json:"FieldNames"`
	Page              *common.Page        `json:"Page,omitempty"`
	OrderBy           *[]OrderBy          `json:"OrderBy,omitempty"`
	ReportName        string              `json:"ReportName,omitempty"`
	ReportType        ReportType          `json:"ReportType"`
	DateRangeType     DateRangeType       `json:"DateRangeType"`
	Format            common.Format       `json:"Format,omitempty"`
	IncludeVAT        common.YesNo        `json:"IncludeVAT,omitempty"`
	IncludeDiscount   *common.YesNo       `json:"IncludeDiscount,omitempty"`
}

type AttributionModel string

const (
	FC   AttributionModel = "FC"   // Первый переход.
	LC   AttributionModel = "LC"   // Последний переход.
	LSC  AttributionModel = "LSC"  // Последний значимый переход.
	LYDC AttributionModel = "LYDC" // Последний переход из Яндекс Директа.
)

type Field string

type SelectionCriteria struct {
	DateFrom string   `json:"DateFrom,omitempty"`
	DateTo   string   `json:"DateTo,omitempty"`
	Filter   []Filter `json:"Filter,omitempty"`
}

type Filter struct {
	Fields   string         `json:"Field"`
	Operator FilterOperator `json:"Operator"`
	Values   []string       `json:"Values"`
}

type FilterOperator string

const (
	Equals                        FilterOperator = "EQUALS"                              // Значение поля равно значению из Values
	NotEquals                     FilterOperator = "NOT_EQUALS"                          // Значение поля не равно значению из Values
	In                            FilterOperator = "IN"                                  // Значение поля равно любому значению из Values
	NotIn                         FilterOperator = "NOT_IN"                              // Значение поля не равно ни одному значению из Values
	LessThan                      FilterOperator = "LESS_THAN"                           // Значение поля меньше значения из Values
	GreaterThan                   FilterOperator = "GREATER_THAN"                        // Значение поля больше значения из Values
	StartsWithIgnoreCase          FilterOperator = "STARTS_WITH_IGNORE_CASE"             // Значение поля начинается со значения из Values
	DoesNotStartWithIgnoreCase    FilterOperator = "DOES_NOT_START_WITH_IGNORE_CASE"     // Значение поля не начинается со значения из Values
	StartsWithAnyIgnoreCase       FilterOperator = "STARTS_WITH_ANY_IGNORE_CASE"         // Значение поля начинается с любого из значений, указанных в Values
	DoesNotStartWithAllIgnoreCase FilterOperator = "DOES_NOT_START_WITH_ALL_IGNORE_CASE" // Значение поля не начинается ни с одного из значений, указанных в Values
)

type OrderBy struct {
	Field     string           `json:"field"`
	SortOrder OrderBySortOrder `json:"SortOrder"`
}

type OrderBySortOrder string

const (
	ASCENDING  OrderBySortOrder = "ASCENDING"  // По возрастанию
	DESCENDING OrderBySortOrder = "DESCENDING" // По убыванию
)

type ReportType string

const (
	AccountPerformanceReport           ReportType = "ACCOUNT_PERFORMANCE_REPORT"             // Статистика по аккаунту рекламодателя
	CampaignPerformanceReport          ReportType = "CAMPAIGN_PERFORMANCE_REPORT"            // Статистика по кампаниям
	AdgroupPerformanceReport           ReportType = "ADGROUP_PERFORMANCE_REPORT"             // Статистика по группам объявлений
	AdPerformanceReport                ReportType = "AD_PERFORMANCE_REPORT"                  // Статистика по объявлениям
	CriteriaPerformanceReport          ReportType = "CRITERIA_PERFORMANCE_REPORT"            // Статистика по условиям показа
	CustomReport                       ReportType = "CUSTOM_REPORT"                          // Статистика с произвольными группировками
	ReachAndFrequencyPerformanceReport ReportType = "REACH_AND_FREQUENCY_PERFORMANCE_REPORT" // Статистика по медийным кампаниям
	SearchQueryPerformanceReport       ReportType = "SEARCH_QUERY_PERFORMANCE_REPORT"        // Статистика по поисковым запросам
)

const (
	name = ReportType("name")
)

type DateRange struct {
	From string
	To   string
}

type DateRangeType string

const (
	DateRangeToday            DateRangeType = "TODAY"               // текущий день
	DateRangeThisWeekMonToday DateRangeType = "THIS_WEEK_MON_TODAY" // Текущая неделя начиная с понедельника, включая текущий день
	DateRangeThisWeekSunToday DateRangeType = "THIS_WEEK_SUN_TODAY" // Текущая неделя начиная с воскресенья
	DateRangeLastWeek         DateRangeType = "LAST_WEEK"           // Прошлая неделя с понедельника по воскресенье
	DateRangeLastBusinessWeek DateRangeType = "LAST_BUSINESS_WEEK"  // Прошлая рабочая неделя с понедельника по пятницу
	DateRangeLastWeekSunSat   DateRangeType = "LAST_WEEK_SUN_SAT"   // Прошлая неделя с воскресенья по субботу
	DateRangeThisMonth        DateRangeType = "THIS_MONTH"          // Текущий календарный месяц
	DateRangeLastMonth        DateRangeType = "LAST_MONTH"          // Полный предыдущий календарный месяц
	DateRangeAllTime          DateRangeType = "ALL_TIME"            // Вся доступная статистика, включая текущий день
	DateRangeCustomDate       DateRangeType = "CUSTOM_DATE"         // Произвольный период. При выборе этого значения необходимо указать даты начала и окончания периода
	DateRangeAuto             DateRangeType = "AUTO"                // Период, за который статистика показов и кликов могла измениться. Период выбирается автоматически в зависимости от того, произошла ли в предыдущий день корректировка статистики
)
