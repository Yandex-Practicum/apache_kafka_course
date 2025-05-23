import json
import logging

logger = logging.getLogger(__name__)


class Message:
    def __init__(self, user_id: str = None, content: str = None):
        self.user_id = user_id
        self.content = content

    def serialize(self) -> str:
        """
        Метод для сериализации в JSON.

        :return: JSON в виде строки.
        """
        try:
            return json.dumps({"user_id": self.user_id, "content": self.content})
        except Exception as e:
            logger.error(f"Serialization error: {e}")
            raise

    @staticmethod
    def deserialize(data: str) -> 'Message':
        """
        Метод для десериализации из JSON в объект.

        :param data: JSON в виде строки.
        :return: JSON, как объект класса Message.
        """
        try:
            json_data = json.loads(data)
            return Message(json_data["user_id"], json_data["content"])
        except Exception as e:
            logger.error(f"Deserialization error: {e}")
            raise
