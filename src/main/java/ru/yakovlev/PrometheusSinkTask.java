package ru.yakovlev;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.MapperFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.json.JsonMapper;
import lombok.extern.slf4j.Slf4j;
import org.apache.kafka.common.utils.AppInfoParser;
import org.apache.kafka.connect.sink.SinkRecord;
import org.apache.kafka.connect.sink.SinkTask;
import ru.yakovlev.model.MetricEvent;

import java.util.Collection;
import java.util.HashMap;
import java.util.Map;

@Slf4j
public class PrometheusSinkTask extends SinkTask {
    private static final String PORT = "port";
    private static final ObjectMapper MAPPER = buildMapper();
    private final static String METRIC_FORMAT_LINE_1 = "# HELP %s %s";
    private final static String METRIC_FORMAT_LINE_2 = "# TYPE %s %s";
    private final static String METRIC_FORMAT_LINE_3 = "%s %s";

    private PrometheusHttpServer httpServer;

    public PrometheusSinkTask() {
    }

    @Override
    public String version() {
        return AppInfoParser.getVersion();
    }

    @Override
    public void start(Map<String, String> map) {
        StringBuilder sb = new StringBuilder();
        map.forEach((key, value) ->
                sb.append("\t").append(key).append(" = ").append(value).append(System.lineSeparator()));
        log.info("Starting PrometheusSinkTask with the properties: {}{}", System.lineSeparator(), sb);
        int port = Integer.parseInt(map.get(PORT));
        try {
            httpServer = PrometheusHttpServer.getInstance(port);
        } catch (Exception e) {
            log.error("Failed to start PrometheusSinkConnector. Reason: {}", e.getMessage());
        }
    }

    @Override
    public void put(Collection<SinkRecord> records) {
        log.info("Received records count {}", records.size());
        for (SinkRecord next : records) {
            String value = (String) next.value();
            log.info("Received metric: {}", value);
            try {
                Map<String, MetricEvent> metricEvents = MAPPER.readValue(value, new TypeReference<HashMap<String, MetricEvent>>() {
                });
                for (Map.Entry<String, MetricEvent> eventEntry : metricEvents.entrySet()) {
                    String metric = buildSingleMetric(eventEntry.getValue());
                    log.info("Built metric: {}", metric);
                    httpServer.addMetric(eventEntry.getKey(), metric);
                }
            } catch (Exception e) {
                log.error("Failed to read value. Reason: {}", e.getMessage());
            }
        }
    }

    @Override
    public void stop() {
        httpServer.stop();
    }

    private String buildSingleMetric(MetricEvent metricEvent) {
        StringBuilder sb = new StringBuilder();
        sb.append(String.format(METRIC_FORMAT_LINE_1, metricEvent.getName(), metricEvent.getDescription())).append(System.lineSeparator());
        sb.append(String.format(METRIC_FORMAT_LINE_2, metricEvent.getName(), metricEvent.getType())).append(System.lineSeparator());
        sb.append(String.format(METRIC_FORMAT_LINE_3, metricEvent.getName(), metricEvent.getValue())).append(System.lineSeparator());
        return sb.toString();
    }

    private static ObjectMapper buildMapper() {
        return JsonMapper.builder()
                .disable(MapperFeature.ACCEPT_CASE_INSENSITIVE_PROPERTIES)
                .build();
    }
}
