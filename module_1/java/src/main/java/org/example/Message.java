package org.example;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;

public class Message {
    private String userId;
    private String content;


    public Message() {
    }

    public Message(String userId, String content) {
        this.userId = userId;
        this.content = content;
    }

    public String getUserId() {
        return userId;
    }

    public void setUserId(String userId) {
        this.userId = userId;
    }

    public String getContent() {
        return content;
    }

    public void setContent(String content) {
        this.content = content;
    }

    // Метод для сериализации в JSON
    public static String serialize(Message message) throws JsonProcessingException {
        ObjectMapper objectMapper = new ObjectMapper();
        return objectMapper.writeValueAsString(message);
    }

    // Метод для десериализации из JSON в объект
    public static Message deserialize(String data) throws Exception {
        ObjectMapper objectMapper = new ObjectMapper();
        return objectMapper.readValue(data, Message.class);
    }

}
