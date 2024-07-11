CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE public.embed_domain (
	embed_domain_id uuid DEFAULT uuid_generate_v4() NOT NULL,
	"domain" text NOT NULL,
	email text NOT NULL,
	created_at timestamptz DEFAULT now() NOT NULL,
	ads bool DEFAULT true NOT NULL,
	updated_at timestamptz DEFAULT now() NOT NULL,
	CONSTRAINT embed_domain_pk PRIMARY KEY (embed_domain_id),
	CONSTRAINT embed_domain_unique UNIQUE (domain)
);

CREATE OR REPLACE FUNCTION public.update_updated_at()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
   NEW.updated_at = now(); 
   RETURN NEW;
END;
$function$
;

create trigger update_updated_at before
update
    on
    public.embed_domain for each row execute function update_updated_at();