package main

// версия сервера метрик без использования храинлища прома
// метрики сохраняются в слайс и выдаются прому по запросу
// добавление делается в начало, чтобы в случае дубликата пром получил последнюю версию
// проверено - пром допускает дубликаты метрик но берет первую
// если интервалы постов и скрейпов будут одинаковы (15с) - дубликаты не появятся
// после каждого скрейпа хранилище сбрасывается
// этот сервер не хранит данные
// для разведение метрик разных клиентов каждый клиент должен помечать свои метрики своей уникальной меткой (serviceID)
// поддерживаются только вектора
// обычные метрики будут пробрасываться но результат непредсказуем (метрики клиентов пересекутся, и в зачет пойдет только последний клиент в цикле)

import (
	"encoding/json"
	"fmt"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"io"
	"log"
	"net/http"
	"sync"
)

type MetricsCache struct {
	mu             sync.Mutex
	metricFamilies []*dto.MetricFamily
}

func (m *MetricsCache) HandlerPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err := m.readMetrics(r.Body); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	r.Body.Close()
}

func (m *MetricsCache) HandlerGet(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()
	contentType := expfmt.Negotiate(r.Header)
	w.Header().Set("Content-Type", string(contentType))
	enc := expfmt.NewEncoder(w, contentType)
	m.export(enc, w)
	m.clear()
}

func (m *MetricsCache) readMetrics(r io.Reader) error {
	var inFamilies []*dto.MetricFamily
	if err := json.NewDecoder(r).Decode(&inFamilies); err != nil {
		fmt.Errorf(err.Error())
	}
	m.append(inFamilies)
	return nil
}

func (m *MetricsCache) export(enc expfmt.Encoder, w http.ResponseWriter) {
	for _, mf := range m.metricFamilies {
		if err := enc.Encode(mf); err != nil {
			http.Error(w, "An error has occurred during metrics encoding:\n\n"+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (m *MetricsCache) append(newMetricFamilies []*dto.MetricFamily) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.metricFamilies = append(newMetricFamilies, m.metricFamilies...)
}

func (m *MetricsCache) clear() {
	m.metricFamilies = nil
}

func NewMetricCache() *MetricsCache {
	return &MetricsCache{}
}

func main() {
	mCache := NewMetricCache()

	http.HandleFunc("/metrics", mCache.HandlerGet)
	http.HandleFunc("/metrics/", mCache.HandlerPost)
	fmt.Println("listen on http://127.0.0.1:8030/metrics")
	log.Fatal(http.ListenAndServe(":8030", nil))
}
