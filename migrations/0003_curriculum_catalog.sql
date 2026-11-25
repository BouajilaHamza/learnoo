-- LRN-014: catalogue foundation for the confirmed mathematics pilot
-- Extends the minimal subjects table with sections and chapters for the Tunisian Baccalauréat.
-- Only the maintainer-confirmed section and subject are seeded. Chapters await teacher validation.

CREATE TABLE sections (
    id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    code       text        NOT NULL,
    title_fr   text        NOT NULL,
    title_ar   text        NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT sections_code_format CHECK (code ~ '^[a-z0-9][a-z0-9-]{0,62}$'),
    CONSTRAINT sections_title_fr_length CHECK (char_length(title_fr) BETWEEN 2 AND 100),
    CONSTRAINT sections_title_ar_length CHECK (char_length(title_ar) BETWEEN 2 AND 100)
);

CREATE UNIQUE INDEX sections_code_key ON sections (code);

-- Extend the minimal subjects table created in 0002_permissions
ALTER TABLE subjects
    ADD COLUMN section_id  bigint REFERENCES sections (id) ON DELETE RESTRICT,
    ADD COLUMN title_fr    text,
    ADD COLUMN title_ar    text,
    ADD COLUMN coefficient integer CHECK (coefficient > 0),
    ADD COLUMN sort_order  integer DEFAULT 0;

CREATE INDEX subjects_section_id_idx ON subjects (section_id);

CREATE TABLE chapters (
    id             bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    subject_id     bigint      NOT NULL REFERENCES subjects (id) ON DELETE RESTRICT,
    code           text        NOT NULL,
    domain         text        NOT NULL, -- 'analyse', 'algebre-geometrie', 'probabilites-statistiques'
    title_fr       text        NOT NULL,
    title_ar       text        NOT NULL,
    description_tn text        NOT NULL DEFAULT '', -- Optional teacher-provided explanation
    sort_order     integer     NOT NULL DEFAULT 0,
    created_at     timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT chapters_code_format CHECK (code ~ '^[a-z0-9][a-z0-9-]{0,62}$'),
    CONSTRAINT chapters_domain_format CHECK (domain ~ '^[a-z0-9][a-z0-9-]{0,62}$'),
    CONSTRAINT chapters_title_fr_length CHECK (char_length(title_fr) BETWEEN 2 AND 200),
    CONSTRAINT chapters_title_ar_length CHECK (char_length(title_ar) BETWEEN 2 AND 200),
    CONSTRAINT chapters_description_tn_length CHECK (char_length(description_tn) BETWEEN 0 AND 2000)
);

CREATE UNIQUE INDEX chapters_subject_code_key ON chapters (subject_id, code);
CREATE INDEX chapters_subject_sort_idx ON chapters (subject_id, sort_order);

-- Stable identifiers for the confirmed pilot; no coefficient or chapter is inferred.
INSERT INTO sections (code, title_fr, title_ar)
VALUES ('math', 'Baccalauréat Mathématiques', 'باكالوريا رياضيات');

INSERT INTO subjects (code, section_id, title_fr, title_ar, sort_order)
VALUES ('bac-math-math', (SELECT id FROM sections WHERE code = 'math'),
        'Mathématiques', 'الرياضيات', 1);
