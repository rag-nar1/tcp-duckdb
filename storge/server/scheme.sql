-- Create the 'user' table
CREATE TABLE user (
    userid INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    usertype TEXT NOT NULL -- [super , norm]
);

-- Create the 'DB' table
CREATE TABLE DB (
    dbid INTEGER PRIMARY KEY AUTOINCREMENT,
    dbname TEXT NOT NULL UNIQUE
);

-- Create the 'tables' table
CREATE TABLE tables (
    tableid INTEGER PRIMARY KEY AUTOINCREMENT,
    tablename TEXT NOT NULL,
    dbid INTEGER NOT NULL,
    FOREIGN KEY (dbid) REFERENCES DB(dbid) ON DELETE CASCADE
);

-- Create the 'dbprivilege' table
CREATE TABLE dbprivilege (
    dbid INTEGER NOT NULL,
    userid INTEGER NOT NULL,
    privilegetype TEXT NOT NULL, -- [read , write]
    PRIMARY KEY (dbid, userid, privilegetype),
    FOREIGN KEY (dbid) REFERENCES DB(dbid) ON DELETE CASCADE,
    FOREIGN KEY (userid) REFERENCES user(userid) ON DELETE CASCADE
);

-- Create the 'tableprivilege' table
CREATE TABLE tableprivilege (
    tableid INTEGER NOT NULL,
    userid INTEGER NOT NULL,
    tableprivilege TEXT NOT NULL, --[select, insert, update, delete]
    PRIMARY KEY (tableid, userid, tableprivilege),
    FOREIGN KEY (tableid) REFERENCES tables(tableid) ON DELETE CASCADE,
    FOREIGN KEY (userid) REFERENCES user(userid) ON DELETE CASCADE
);

CREATE TABLE postgres (
    dbid INTEGER NOT NULL UNIQUE,
    connstr TEXT NOT NULL,
    PRIMARY KEY (dbid),
    FOREIGN KEY (dbid) REFERENCES DB(dbid) ON DELETE CASCADE
);

CREATE TABLE keys (
    dbid INTEGER NOT NULL UNIQUE,
    key CHAR(32) NOT NULL,
    PRIMARY KEY (dbid),
    FOREIGN KEY (dbid) REFERENCES DB(dbid) ON DELETE CASCADE
);