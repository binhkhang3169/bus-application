package repository

import (
	"context"
	"fmt"
	"go-bigquery-dashboard/internal/model" // <-- Nhớ kiểm tra lại đường dẫn module của bạn

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const (
	tripSearchTable = "trip_searches_raw"
)

type TripSearchRepository struct {
	client *bigquery.Client
	ctx    context.Context
}

func NewTripSearchRepository(credentialsFile string) (*TripSearchRepository, error) {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, fmt.Errorf("bigquery.NewClient: %v", err)
	}
	return &TripSearchRepository{
		client: client,
		ctx:    ctx,
	}, nil
}

func (r *TripSearchRepository) Close() {
	r.client.Close()
}

// (Các hàm GetTopRoutes, GetTopProvinces, GetSearchesByHourOfDay không thay đổi, giữ nguyên)

func (r *TripSearchRepository) GetTopRoutes(startDate, endDate string, limit int) ([]model.TopRouteStat, error) {
	query := fmt.Sprintf(`
		SELECT fromProvinceId, toProvinceId, COUNT(*) as search_count
		FROM `+"`%s.%s.%s`"+`
		WHERE searchTimestamp >= TIMESTAMP(@startDate) AND searchTimestamp < TIMESTAMP(@endDate)
		GROUP BY fromProvinceId, toProvinceId
		ORDER BY search_count DESC
		LIMIT @limit
	`, projectID, datasetID, tripSearchTable)

	q := r.client.Query(query)
	q.Parameters = []bigquery.QueryParameter{
		{Name: "startDate", Value: startDate},
		{Name: "endDate", Value: endDate},
		{Name: "limit", Value: limit},
	}
	it, err := q.Read(r.ctx)
	if err != nil {
		return nil, err
	}

	var results []model.TopRouteStat
	for {
		var row model.TopRouteStat
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, nil
}

func (r *TripSearchRepository) GetTopProvinces(startDate, endDate string, limit int) (*model.TopProvincesResponse, error) {
	params := []bigquery.QueryParameter{
		{Name: "startDate", Value: startDate},
		{Name: "endDate", Value: endDate},
		{Name: "limit", Value: limit},
	}

	whereClause := "WHERE searchTimestamp >= TIMESTAMP(@startDate) AND searchTimestamp < TIMESTAMP(@endDate)"

	originsQuery := fmt.Sprintf(`
		SELECT fromProvinceId as provinceId, COUNT(*) as search_count
		FROM `+"`%s.%s.%s`"+` %s
		GROUP BY provinceId ORDER BY search_count DESC LIMIT @limit
	`, projectID, datasetID, tripSearchTable, whereClause)
	qOrigins := r.client.Query(originsQuery)
	qOrigins.Parameters = params
	itOrigins, err := qOrigins.Read(r.ctx)
	if err != nil {
		return nil, err
	}
	var topOrigins []model.TopProvinceStat
	for {
		var row model.TopProvinceStat
		err := itOrigins.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		topOrigins = append(topOrigins, row)
	}

	destsQuery := fmt.Sprintf(`
		SELECT toProvinceId as provinceId, COUNT(*) as search_count
		FROM `+"`%s.%s.%s`"+` %s
		GROUP BY provinceId ORDER BY search_count DESC LIMIT @limit
	`, projectID, datasetID, tripSearchTable, whereClause)
	qDests := r.client.Query(destsQuery)
	qDests.Parameters = params
	itDests, err := qDests.Read(r.ctx)
	if err != nil {
		return nil, err
	}
	var topDests []model.TopProvinceStat
	for {
		var row model.TopProvinceStat
		err := itDests.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		topDests = append(topDests, row)
	}

	return &model.TopProvincesResponse{
		TopOrigins:      topOrigins,
		TopDestinations: topDests,
	}, nil
}

func (r *TripSearchRepository) GetSearchesByHourOfDay(startDate, endDate string) ([]model.HourlySearchStat, error) {
	query := fmt.Sprintf(`
		SELECT EXTRACT(HOUR FROM searchTimestamp AT TIME ZONE 'Asia/Ho_Chi_Minh') as hour, COUNT(*) as search_count
		FROM `+"`%s.%s.%s`"+`
		WHERE searchTimestamp >= TIMESTAMP(@startDate) AND searchTimestamp < TIMESTAMP(@endDate)
		GROUP BY hour
		ORDER BY hour ASC
	`, projectID, datasetID, tripSearchTable)
	q := r.client.Query(query)
	q.Parameters = []bigquery.QueryParameter{
		{Name: "startDate", Value: startDate},
		{Name: "endDate", Value: endDate},
	}
	it, err := q.Read(r.ctx)
	if err != nil {
		return nil, err
	}

	var results []model.HourlySearchStat
	for {
		var row model.HourlySearchStat
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, nil
}

// GetSearchesOverTime thống kê lượt tìm kiếm theo ngày.
func (r *TripSearchRepository) GetSearchesOverTime(startDate, endDate string) ([]model.TimeSeriesDataPoint, error) {
	// SỬA LỖI LẦN CUỐI: Dùng DATE_TRUNC và CAST sang DATE, để Go driver tự chuyển sang string.
	// Đây là cách làm ổn định nhất.
	query := fmt.Sprintf(`
		SELECT CAST(DATE_TRUNC(searchTimestamp, DAY, 'Asia/Ho_Chi_Minh') AS DATE) as date, COUNT(*) as value
		FROM `+"`%s.%s.%s`"+`
		WHERE searchTimestamp >= TIMESTAMP(@startDate) AND searchTimestamp < TIMESTAMP(@endDate)
		GROUP BY date
		ORDER BY date ASC
	`, projectID, datasetID, tripSearchTable)

	q := r.client.Query(query)
	q.Parameters = []bigquery.QueryParameter{
		{Name: "startDate", Value: startDate},
		{Name: "endDate", Value: endDate},
	}
	it, err := q.Read(r.ctx)
	if err != nil {
		return nil, fmt.Errorf("BigQuery Read failed for GetSearchesOverTime: %w", err)
	}

	var results []model.TimeSeriesDataPoint
	for {
		// Ở đây ta cần một struct trung gian vì kiểu trả về của BigQuery là DATE, không phải STRING
		var row struct {
			Date  bigquery.NullDate `bigquery:"date"`
			Value int64             `bigquery:"value"`
		}
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Row iteration failed for GetSearchesOverTime: %w", err)
		}
		// Chuyển đổi từ bigquery.NullDate sang string
		results = append(results, model.TimeSeriesDataPoint{
			Date:  row.Date.Date.String(),
			Value: float64(row.Value),
		})
	}
	return results, nil
}
