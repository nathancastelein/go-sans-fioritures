package main

import (
	"context"
	"log/slog"
)

type ReportRepository interface {
	ReportSuspiciousActivity(ctx context.Context, report Report)
}

type Report struct {
	Stone  string
	Report string
}

func (r Report) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("stone", r.Stone),
		slog.String("report", r.Report))
}

type InMemoryReportRepository struct{}

func NewInMemoryReportRepository() ReportRepository {
	return &InMemoryReportRepository{}
}

func (i *InMemoryReportRepository) ReportSuspiciousActivity(ctx context.Context, report Report) {
	slog.Info("suspicious activity reported", slog.Any("report", report))
}
