package ru.yakovlev;

import lombok.extern.slf4j.Slf4j;
import org.eclipse.jetty.server.Server;
import org.eclipse.jetty.servlet.ServletContextHandler;
import org.eclipse.jetty.servlet.ServletHolder;

import java.util.concurrent.ConcurrentHashMap;

@Slf4j
public class PrometheusHttpServer {

    private final ConcurrentHashMap<String, String> metrics = new ConcurrentHashMap<>();
    private static PrometheusHttpServer instance;
    private Server server;

    private PrometheusHttpServer() {
    }

    public static PrometheusHttpServer getInstance(int port) throws Exception {
        if (instance == null) {
            instance = new PrometheusHttpServer();
            log.info("Starting server at port {}", port);
            instance.start(port);
        }
        return instance;
    }


    public void start(int port) throws Exception {
        log.info("Starting a new server at port {}", port);
        this.server = new Server(port);
        ServletContextHandler handler = new ServletContextHandler();
        handler.addServlet(new ServletHolder(new MetricsServlet(metrics)), "/metrics");
        server.setHandler(handler);
        server.start();
        server.join();
    }

    public void addMetric(String name, String data) {
        metrics.put(name, data);
    }

    public void stop() {
        if (instance == null) return;
        try {
            this.server.stop();
        } catch (Exception e) {
            log.error("Failed to stop server. Reason: {}", e.getMessage());
        }
    }
}
