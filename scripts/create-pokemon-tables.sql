DROP TABLE IF EXISTS pokemon_cards;
DROP TABLE IF EXISTS card_sets;

CREATE TABLE card_sets (
  name VARCHAR(255) NOT NULL,
  series VARCHAR(255) NOT NULL,
  total INT NOT NULL,
  PRIMARY KEY (name, series)
);

CREATE TABLE pokemon_cards (
  id INT AUTO_INCREMENT NOT NULL,
  name VARCHAR(255) NOT NULL,
  number VARCHAR(255) NOT NULL,
  rarity VARCHAR(255) NOT NULL,
  value DECIMAL(10,2) NOT NULL,
  image_url VARCHAR(255) NOT NULL,
  card_set VARCHAR(255) NOT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (card_set) REFERENCES card_sets(name)
);
