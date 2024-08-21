CREATE TABLE working_hours (
                               kitchen_id UUID REFERENCES kitchens(id),
                               day_of_week INTEGER NOT NULL,
                               open_time TIME NOT NULL,
                               close_time TIME NOT NULL,
                               PRIMARY KEY (kitchen_id, day_of_week)
);
