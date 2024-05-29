CREATE TABLE public.events(
    id uuid PRIMARY KEY,
    creator_id int NOT NULL,
    title varchar(255) NOT NULL,
    description text NOT NULL,
    start_dt timestamp NOT NULL,
    end_dt timestamp NOT NULL,
    notify_before bigint NOT NULL
)