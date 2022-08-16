-- Create a keyspace
CREATE KEYSPACE IF NOT EXISTS hasuniversity WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : '1' };

-- Create a UDT
CREATE TYPE IF NOT EXISTS hasuniversity.card (
front text,
back text
);

-- Create a table
CREATE TABLE IF NOT EXISTS hasuniversity.flashcard (
uuid text PRIMARY KEY,
author text,
name text,
tags list<text>,
cards list<FROZEN <card>>,
last_update_timestamp timestamp
);

-- Add dummy data
INSERT INTO hasuniversity.flashcard
(uuid, author, name, tags, cards, last_update_timestamp)
VALUES (
    'abc123', 
    'Hank Pecker', 
    'hasChud',
    ['US', 'Politics'],
    [ { front : '9/11', back : 'PepeLa' } ],
    toTimeStamp(now())
);
