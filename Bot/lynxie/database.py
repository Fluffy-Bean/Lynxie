from sqlalchemy import create_engine, Column, Integer, String, DateTime
from sqlalchemy.orm import declarative_base, sessionmaker
from lynxie.config import DATABASE_URI


Base = declarative_base()


class CommandHistory(Base):
    __tablename__ = "command_history"

    id = Column(Integer, primary_key=True)
    command = Column(String)
    user = Column(Integer)
    channel = Column(Integer)
    guild = Column(Integer)
    timestamp = Column(DateTime)


class Database:
    def __init__(self):
        self.engine = create_engine(DATABASE_URI)
        self.session = sessionmaker(bind=self.engine)
        self.session = self.session()

    def make_database(self):
        Base.metadata.create_all(self.engine)


if __name__ == "__main__":
    db = Database()
    db.make_database()
