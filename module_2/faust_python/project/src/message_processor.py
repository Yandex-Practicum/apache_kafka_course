import logging.config

import faust

from settings.logging import LOGGING_CONFIG

logging.config.dictConfig(LOGGING_CONFIG)
logger = logging.getLogger("message_processor")

app = faust.App(
    "message-processor-app",
    broker="localhost:29092",
    value_serializer="raw",
    store="rocksdb://",
)


def load_banned_words(file_path: str) -> list[str]:
    """
    Чтение списка запрещенных слов.

    :param file_path: Путь к файлу.
    :return: Список запрещенных слов.
    """
    result = list()

    try:
        with open(file_path, "r", encoding="utf-8") as file:
            for line in file.readlines():
                result.append(line.strip())

        logger.info(f"Loaded banned words: {result}")
    except Exception as exp:
        logger.error(f"Error loading banned words from file: {file_path}. {exp}")

    return result


# Таблица заблокированных пользователей
blocked_users = app.Table(
    "blocked_users",
    partitions=1,   # Количество партиций
    default=None,   # Функция или тип для пропущенных ключей
)
# Топик входящих сообщений
messages_topic = app.topic("messages", key_type=str, value_type=str)
# Топик отфильтрованных сообщений
filtered_messages_topic = app.topic("filtered_messages", key_type=str, value_type=str)


@app.agent(messages_topic)
async def process_messages(stream):
    """
    Агент для фильтрации сообщений.
    """
    banned_words = load_banned_words("resources/banned_words.txt")

    async for direction, message in stream.items():
        logger.info(f"New message: {direction=}, {message=}")
        if not message:
            # Исключаем пустые сообщения
            continue

        sender, receiver = direction.split("->")
        logger.info(f"{sender=}, {receiver=}, {message=}")

        # Проверка списка заблокированных пользователей
        blocked_list = blocked_users.get(receiver)
        if blocked_list:
            logger.info(f"Check that {sender} is not in {blocked_list}")
            if sender in blocked_list:
                # Если получатель заблокировал отправителя, то игнорируем сообщение
                logger.info("This message will be ignored")
                continue

        # Замена запрещенных слов на звездочки
        for banned_word in banned_words:
            message = message.replace(banned_word, "***")

        # Отправка отфильтрованных сообщений в `filtered_messages`
        logger.info(f"Filtered key={direction}, {message=}")
        await filtered_messages_topic.send(key=direction, value=message)


if __name__ == '__main__':
    app.main()
