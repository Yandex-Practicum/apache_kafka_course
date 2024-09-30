from pydantic_avro.base import AvroBase

class Message(AvroBase):
    text: str
