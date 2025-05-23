package ru.yakovlev;

import jakarta.servlet.http.HttpServlet;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;

import java.io.IOException;
import java.util.concurrent.ConcurrentHashMap;

public class MetricsServlet extends HttpServlet {
    private final ConcurrentHashMap<String, String> metrics;

    public MetricsServlet(ConcurrentHashMap<String, String> metrics) {
        this.metrics = metrics;
    }

    @Override
    protected void doGet(HttpServletRequest req, HttpServletResponse resp) throws IOException {
        resp.setContentType("text/plain");
        StringBuilder response = new StringBuilder();
        metrics.values().forEach(response::append);
        resp.getWriter().write(response.toString());
    }
}
