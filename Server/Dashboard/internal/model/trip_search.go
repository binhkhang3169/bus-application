package model

// TopRouteStat đại diện cho thống kê của một tuyến đường.
type TopRouteStat struct {
	FromProvinceID int64 `json:"from_province_id" bigquery:"fromProvinceId"`
	ToProvinceID   int64 `json:"to_province_id" bigquery:"toProvinceId"`
	SearchCount    int64 `json:"search_count" bigquery:"search_count"`
}

// TopProvinceStat đại diện cho thống kê của một tỉnh (đi hoặc đến).
type TopProvinceStat struct {
	ProvinceID  int64 `json:"province_id" bigquery:"provinceId"`
	SearchCount int64 `json:"search_count" bigquery:"search_count"`
}

// TopProvincesResponse là cấu trúc trả về cho API top tỉnh.
type TopProvincesResponse struct {
	TopOrigins      []TopProvinceStat `json:"top_origins"`
	TopDestinations []TopProvinceStat `json:"top_destinations"`
}

// HourlySearchStat đại diện cho thống kê tìm kiếm trong một giờ.
type HourlySearchStat struct {
	Hour        int64 `json:"hour" bigquery:"hour"`
	SearchCount int64 `json:"search_count" bigquery:"search_count"`
}
