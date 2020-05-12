CREATE TABLE IF NOT EXISTS  user_info (
    id INT NOT NULL AUTO_INCREMENT,
    username VARCHAR(20),
    password VARCHAR(80),
    email VARCHAR(50),
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS centers (
  id INT NOT NULL AUTO_INCREMENT,
  center_name TEXT NOT NULL,
  twitter_access_token TEXT,
  twitter_secret TEXT,
  facebook_page_name  TEXT,
  facebook_page_id TEXT,
  facebook_page_authtoken TEXT,
  discord_channel_id TEXT,
  discord_guild_id TEXT,
  gg_leap_link TEXT,
  PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS schedules (
    id INT NOT NULL AUTO_INCREMENT,
    time_to_post TEXT,
    day_of_week TEXT,
    game TEXT,
    center_id int,
    FOREIGN KEY (center_id) REFERENCES centers(id),
    PRIMARY KEY(id)
);