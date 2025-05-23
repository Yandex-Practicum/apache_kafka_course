package ru.yakovlev.model;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;

@Getter
@AllArgsConstructor
@NoArgsConstructor
public class MetricEvent {
    private String name;
    private String type;
    private String description;
    private Double value;
}
