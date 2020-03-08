package burgers

const BurgersDDL = `CREATE TABLE IF NOT EXISTS burgers
(
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    price       INT  NOT NULL CHECK ( price > 0 ),
    removed     BOOLEAN DEFAULT FALSE
);`

const GetBurgers = "SELECT id, name, price FROM burgers WHERE removed = FALSE;"

const SaveBurger = "INSERT INTO burgers (name, price) VALUES ($1, $2);"

const RemoveBurger = "UPDATE burgers SET removed = true WHERE id = $1;"
