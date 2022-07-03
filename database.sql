/* 1. create db_links database */
DROP DATABASE IF EXISTS db_links;
CREATE DATABASE db_links
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;


/* 2. create t_links table */
DROP TABLE IF EXISTS t_links;
CREATE TABLE IF NOT EXISTS public.t_links
(
    id SERIAL NOT NULL,
    url character varying(100) NOT NULL,
    link character varying(20) NOT NULL,
    requests bigint DEFAULT 0,
    deleted bit(1) DEFAULT '0'::"bit",
    CONSTRAINT t_links_pkey PRIMARY KEY (id)
    )


/* 3. create function fn_create_link for links creation */
/*
fn_create_link returns int value and string. There could be 3 options:
- 0, "" - URL was created
- 1, "valid code" - URL is already in DB
- 2, "" - linkID is already in DB
*/

DROP FUNCTION IF EXISTS fn_create_link;
CREATE FUNCTION fn_create_link
(p_url character,
 p_link_id character)
    RETURNS RECORD AS $$
DECLARE
ret RECORD;
BEGIN
    IF (SELECT COUNT(*) FROM t_links WHERE url = p_url) <> 0
    THEN
SELECT 1, (SELECT link FROM t_links WHERE url = p_url)::text INTO ret;
RETURN ret;
END IF;

    IF (SELECT COUNT(*) FROM t_links WHERE link =  p_link_id) <> 0
    THEN
SELECT 2, '' INTO ret;
RETURN ret;
END IF;

INSERT INTO public.t_links(
    url, link)
VALUES (p_url, p_link_id);

SELECT 0, '' INTO ret;
RETURN ret;
END;$$ LANGUAGE plpgsql;