package repository

import (
	"context"
	"fmt"
	"go-bigquery-dashboard/internal/model"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const (
	projectID    = "dacntt-dfabb"
	datasetID    = "duancntt"
	ticketTable  = "tickets"
	invoiceTable = "invoices"
)

type DashboardRepository struct {
	client *bigquery.Client
	ctx    context.Context
}

func NewDashboardRepository(credentialsFile string) (*DashboardRepository, error) {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, fmt.Errorf("bigquery.NewClient: %v", err)
	}
	return &DashboardRepository{
		client: client,
		ctx:    ctx,
	}, nil
}

func (r *DashboardRepository) Close() {
	r.client.Close()
}

func (r *DashboardRepository) GetKPIs(startDate, endDate string) (*model.KpiData, error) {
	data := &model.KpiData{}
	var err error

	whereClause := "WHERE created_at >= @startDate AND created_at < @endDate"
	params := []bigquery.QueryParameter{
		{Name: "startDate", Value: startDate},
		{Name: "endDate", Value: endDate},
	}

	revQuery := fmt.Sprintf("SELECT CAST(COALESCE(SUM(final_amount), 0) AS FLOAT64) FROM `%s.%s.invoices` %s AND payment_status = '1'", projectID, datasetID, whereClause)
	q := r.client.Query(revQuery)
	q.Parameters = params
	it, err := q.Read(r.ctx)
	if err != nil {
		return nil, err
	}
	var revRow struct{ Value float64 }
	if err = it.Next(&revRow); err != nil && err != iterator.Done {
		return nil, err
	}
	data.TotalRevenue = revRow.Value

	ticketQuery := fmt.Sprintf("SELECT COUNT(*) FROM `%s.%s.tickets` %s", projectID, datasetID, whereClause)
	q = r.client.Query(ticketQuery)
	q.Parameters = params
	if data.TotalTickets, err = r.executeCountQuery(q); err != nil {
		return nil, err
	}

	invoiceQuery := fmt.Sprintf("SELECT COUNT(*) FROM `%s.%s.invoices` %s", projectID, datasetID, whereClause)
	q = r.client.Query(invoiceQuery)
	q.Parameters = params
	if data.TotalInvoices, err = r.executeCountQuery(q); err != nil {
		return nil, err
	}

	unpaidTicketsQuery := fmt.Sprintf(`
        SELECT COUNT(t.Ticket_Id) 
        FROM `+"`%s.%s.tickets`"+` t
        LEFT JOIN `+"`%s.%s.invoices`"+` i ON t.Ticket_Id = i.ticket_id
        WHERE t.created_at >= @startDate AND t.created_at < @endDate AND i.invoice_id IS NULL
    `, projectID, datasetID, projectID, datasetID)
	q = r.client.Query(unpaidTicketsQuery)
	q.Parameters = params
	if data.UnpaidTickets, err = r.executeCountQuery(q); err != nil {
		return nil, err
	}

	return data, nil
}

func (r *DashboardRepository) GetRevenueOverTime(startDate, endDate, groupBy string) ([]model.TimeSeriesDataPoint, error) {
	query := fmt.Sprintf(`
		SELECT
			format_date("%%Y-%%m-%%d", DATE_TRUNC(created_at, %s)) as date,
			CAST(SUM(final_amount) AS FLOAT64) as value
		FROM `+"`%s.%s.invoices`"+`
		WHERE created_at >= @startDate AND created_at < @endDate AND payment_status = '1'
		GROUP BY date
		ORDER BY date ASC
	`, groupBy, projectID, datasetID)

	q := r.client.Query(query)
	q.Parameters = []bigquery.QueryParameter{
		{Name: "startDate", Value: startDate},
		{Name: "endDate", Value: endDate},
	}
	it, err := q.Read(r.ctx)
	if err != nil {
		return nil, err
	}

	var results []model.TimeSeriesDataPoint
	for {
		var row model.TimeSeriesDataPoint
		if err := it.Next(&row); err == iterator.Done {
			break
		} else if err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, nil
}

func (r *DashboardRepository) GetTicketDistribution(startDate, endDate string) (*model.TicketDistributionData, error) {
	data := &model.TicketDistributionData{}
	var err error

	whereClause := "WHERE created_at >= @startDate AND created_at < @endDate"
	params := []bigquery.QueryParameter{
		{Name: "startDate", Value: startDate},
		{Name: "endDate", Value: endDate},
	}

	byChannelQuery := fmt.Sprintf(`
		SELECT 
			CASE Booking_Channel WHEN 0 THEN 'Web' WHEN 1 THEN 'App' WHEN 2 THEN 'Offline' ELSE 'Unknown' END as Category, 
			CAST(COUNT(*) AS FLOAT64) as Value 
		FROM `+"`%s.%s.tickets`"+` %s GROUP BY Category
	`, projectID, datasetID, whereClause)
	q := r.client.Query(byChannelQuery)
	q.Parameters = params
	data.ByChannel, err = r.executeStatQuery(q)
	if err != nil {
		return nil, err
	}

	byTypeQuery := fmt.Sprintf(`
		SELECT 
			CASE Type WHEN 0 THEN 'One-way' WHEN 1 THEN 'Round-trip' ELSE 'Unknown' END as Category, 
			CAST(COUNT(*) AS FLOAT64) as Value 
		FROM `+"`%s.%s.tickets`"+` %s GROUP BY Category
	`, projectID, datasetID, whereClause)
	q = r.client.Query(byTypeQuery)
	q.Parameters = params
	data.ByType, err = r.executeStatQuery(q)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *DashboardRepository) executeCountQuery(q *bigquery.Query) (int64, error) {
	it, err := q.Read(r.ctx)
	if err != nil {
		return 0, err
	}
	var row []bigquery.Value
	if err = it.Next(&row); err != nil {
		if err == iterator.Done {
			return 0, nil
		}
		return 0, err
	}
	if val, ok := row[0].(int64); ok {
		return val, nil
	}
	return 0, fmt.Errorf("could not convert count to int64")
}

func (r *DashboardRepository) executeStatQuery(q *bigquery.Query) ([]model.StatItem, error) {
	it, err := q.Read(r.ctx)
	if err != nil {
		return nil, err
	}

	var results []model.StatItem
	for {
		var row struct {
			Category string
			Value    float64
		}
		if err := it.Next(&row); err == iterator.Done {
			break
		} else if err != nil {
			return nil, err
		}
		if row.Category == "" {
			row.Category = "UNKNOWN"
		}
		results = append(results, model.StatItem{Category: row.Category, Value: row.Value})
	}
	return results, nil
}
