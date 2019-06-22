package postgres

const dumpTemplate = `
--
-- PostgreSQL database dump
--

-- Dumped from database version {{ .Version }}

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '{{ .DBScheme }}', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';
SET default_with_oids = false;

{{range .Tables}}

DROP TABLE IF EXISTS {{ .Name }};
--
-- Name: {{ .Name }}; Type: TABLE; Schema: {{ .DBScheme }}; Owner: -
--
{{ .SequenceScheme }}

--
-- Name: {{ .Name }}; Type: SEQUENCE; Schema: {{ .DBScheme }}; Owner: -
--

{{ .Schema }}

--
-- Data for Name: {{ .Name }}; Type: TABLE DATA; Schema: {{ .DBScheme }}; Owner: -
--

{{ if .Data }}
INSERT INTO {{ .Name }} VALUES {{ .Data }};
{{ end }}

{{ end }}

-- Dump completed on {{ .EndTime }}
`
