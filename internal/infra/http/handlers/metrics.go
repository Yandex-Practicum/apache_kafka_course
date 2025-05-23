package handlers

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"kafka-connect/internal/application/repo"
	"kafka-connect/internal/domain"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const tmpl = `{{range .}}<br/>
# HELP {{.Name}} {{.Description}}<br/>
# TYPE {{.Name}} {{.Type}}<br/>
{{.Name}} {{.Value}}<br/>
{{end}}<br/>`

type MetricsHandler struct {
	repo repo.Repo
	l    *zap.SugaredLogger
}

func (h *MetricsHandler) getMetric(ctx context.Context, name domain.MName) (domain.Metric, error) {
	v, err := h.repo.Get(ctx, name)
	if err != nil {
		return domain.Metric{}, err
	}
	return domain.NewMetric(v)
}

func (h *MetricsHandler) getAllMetrics(ctx context.Context) ([]domain.Metric, error) {
	var (
		data []domain.Metric
		v    domain.Metric
		err  error
	)

	for _, name := range domain.GetMetricNames() {
		v, err = h.getMetric(ctx, name)
		if err != nil {
			if !errors.Is(err, repo.NotExistsKey) {
				return nil, err
			}
		} else {
			data = append(data, v)
		}
	}

	return data, nil
}

func (h *MetricsHandler) Handler(res http.ResponseWriter, _ *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t, err := template.New("ml").Parse(tmpl)
	if err != nil {
		h.l.Fatal("Template parse error")
	}

	buf := bytes.Buffer{}

	data, err := h.getAllMetrics(ctx)
	if err != nil {
		h.l.Fatal("Template parse error")
	}

	if err := t.Execute(&buf, data); err != nil {
		h.l.Fatal("Template exec error")
	}

	res.Header().Set("Content-Type", "text/html")
	_, err = res.Write(buf.Bytes())
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
}

func NewMetricsHandler(repo repo.Repo, log *zap.SugaredLogger) *MetricsHandler {
	return &MetricsHandler{repo, log}
}
