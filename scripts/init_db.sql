CREATE TABLE Users (
                       id SERIAL PRIMARY KEY,
                       email VARCHAR(255) NOT NULL UNIQUE,
                       password_hash VARCHAR(255) NOT NULL,
                       name VARCHAR(255),
                       phone_number VARCHAR(15),
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Items (
                       id SERIAL PRIMARY KEY,
                       user_id INTEGER REFERENCES Users(id) ON DELETE CASCADE,
                       title VARCHAR(255) NOT NULL,
                       description TEXT,
                       price_per_day DECIMAL(10, 2) NOT NULL,
                       category VARCHAR(255),
                       available BOOLEAN DEFAULT TRUE,
                       location VARCHAR(255),
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
m
CREATE TABLE Bookings (
                          id SERIAL PRIMARY KEY,
                          item_id INTEGER REFERENCES Items(id) ON DELETE CASCADE,
                          user_id INTEGER REFERENCES Users(id) ON DELETE CASCADE,
                          start_date DATE NOT NULL,
                          end_date DATE NOT NULL,
                          total_price DECIMAL(10, 2) NOT NULL,
                          status VARCHAR(50) DEFAULT 'pending',
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);