package ru.yakovlev;

import lombok.extern.slf4j.Slf4j;
import org.apache.kafka.common.config.ConfigDef;
import org.apache.kafka.common.utils.AppInfoParser;
import org.apache.kafka.connect.connector.Task;
import org.apache.kafka.connect.sink.SinkConnector;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;

@Slf4j
public class PrometheusSinkConnector extends SinkConnector {

    private static final int DEFAULT_PORT = 8095;
    private Map<String, String> props;

    @Override
    public String version() {
        return AppInfoParser.getVersion();
    }

    @Override
    public void start(Map<String, String> props) {
        StringBuilder sb = new StringBuilder();
        props.forEach((key, value) ->
                sb.append("\t").append(key).append(" = ").append(value).append(System.lineSeparator()));
        this.props = props;
        log.info("Starting PrometheusSinkConnector with the properties: {}{}", System.lineSeparator(), sb);
    }

    @Override
    public Class<? extends Task> taskClass() {
        return PrometheusSinkTask.class;
    }

    @Override
    public List<Map<String, String>> taskConfigs(int maxTasks) {
        List<Map<String, String>> configs = new ArrayList<>();
        for (int i = 0; i < maxTasks; i++) {
            configs.add(props);
        }
        return configs;
    }

    @Override
    public void stop() {
    }

    @Override
    public ConfigDef config() {
        return new ConfigDef()
                .define("port", ConfigDef.Type.INT, DEFAULT_PORT, ConfigDef.Importance.HIGH, "Prometheus port");
    }
}
