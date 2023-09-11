package otelsql_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"

	"github.com/AmeerIbrahimm/goagent/instrumentation/otelsql"
	"github.com/AmeerIbrahimm/goagent/instrumentation/otelsql/internal/test/oteltest"
	"github.com/AmeerIbrahimm/goagent/instrumentation/otelsql/internal/test/sqlmock"
)

func TestRecordStats(t *testing.T) {
	t.Parallel()

	expectedMetrics := expectedStatsMetric()

	oteltest.New(
		oteltest.MetricsEqualJSON(expectedMetrics),
		oteltest.MockDatabase(func(m sqlmock.Sqlmock) {
			m.ExpectPing()
		}),
	).
		Run(t, func(sc oteltest.SuiteContext) {
			db, err := newDB(sc.DatabaseDSN())
			require.NoError(t, err)

			err = otelsql.RecordStats(db,
				otelsql.WithMeterProvider(sc.MeterProvider()),
				otelsql.WithMinimumReadDBStatsInterval(100*time.Millisecond),
				otelsql.WithInstanceName("default"),
				otelsql.WithSystem(semconv.DBSystemPostgreSQL),
			)
			require.NoError(t, err)

			err = db.Ping()
			require.NoError(t, err)
		})
}

func expectedStatsMetric() string {
	return expectedMetricsFromFile("stats.json")
}
