-- LRN-014: Official Tunisian Baccalauréat Curriculum Catalogue Seed
-- Extends the minimal subjects table with sections and chapters for the Tunisian Baccalauréat.
-- Pilot focus: Section Mathématiques & Sciences Expérimentales -> Mathématiques.

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
    ADD COLUMN coefficient integer DEFAULT 1,
    ADD COLUMN sort_order  integer DEFAULT 0;

CREATE INDEX subjects_section_id_idx ON subjects (section_id);

CREATE TABLE chapters (
    id             bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    subject_id     bigint      NOT NULL REFERENCES subjects (id) ON DELETE RESTRICT,
    code           text        NOT NULL,
    domain         text        NOT NULL, -- 'analyse', 'algebre-geometrie', 'probabilites-statistiques'
    title_fr       text        NOT NULL,
    title_ar       text        NOT NULL,
    description_tn text        NOT NULL, -- Tunisian Derja pedagogical guidance and revision tips
    sort_order     integer     NOT NULL DEFAULT 0,
    created_at     timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT chapters_code_format CHECK (code ~ '^[a-z0-9][a-z0-9-]{0,62}$'),
    CONSTRAINT chapters_domain_format CHECK (domain ~ '^[a-z0-9][a-z0-9-]{0,62}$'),
    CONSTRAINT chapters_title_fr_length CHECK (char_length(title_fr) BETWEEN 2 AND 200),
    CONSTRAINT chapters_title_ar_length CHECK (char_length(title_ar) BETWEEN 2 AND 200),
    CONSTRAINT chapters_description_tn_length CHECK (char_length(description_tn) BETWEEN 10 AND 2000)
);

CREATE UNIQUE INDEX chapters_subject_code_key ON chapters (subject_id, code);
CREATE INDEX chapters_subject_sort_idx ON chapters (subject_id, sort_order);

-- ----------------------------------------------------------------------------
-- SEED DATA: Official Tunisian Baccalauréat Sections
-- ----------------------------------------------------------------------------
INSERT INTO sections (code, title_fr, title_ar) VALUES
    ('math', 'Baccalauréat Mathématiques', 'باكالوريا رياضيات'),
    ('sciences', 'Baccalauréat Sciences Expérimentales', 'باكالوريا علوم تجريبية'),
    ('technique', 'Baccalauréat Sciences Techniques', 'باكالوريا علوم تقنية'),
    ('informatique', 'Baccalauréat Sciences de l''Informatique', 'باكالوريا علوم الإعلامية'),
    ('economie', 'Baccalauréat Économie et Gestion', 'باكالوريا اقتصاد وتصرف'),
    ('lettres', 'Baccalauréat Lettres', 'باكالوريا آداب');

-- ----------------------------------------------------------------------------
-- SEED DATA: Subjects for Math and Sciences Sections
-- ----------------------------------------------------------------------------
INSERT INTO subjects (code, section_id, title_fr, title_ar, coefficient, sort_order) VALUES
    ('bac-math-math', (SELECT id FROM sections WHERE code = 'math'), 'Mathématiques', 'الرياضيات', 4, 1),
    ('bac-sciences-math', (SELECT id FROM sections WHERE code = 'sciences'), 'Mathématiques', 'الرياضيات', 3, 1);

-- ----------------------------------------------------------------------------
-- SEED DATA: Complete 12 Chapters of Tunisian BAC Mathematics
-- Includes official French terms, Arabic titles, and authentic Tunisian Derja pedagogical notes.
-- ----------------------------------------------------------------------------

-- 1. Continuité et Limites
INSERT INTO chapters (subject_id, code, domain, title_fr, title_ar, description_tn, sort_order)
VALUES (
    (SELECT id FROM subjects WHERE code = 'bac-math-math'),
    'continuite-limites',
    'analyse',
    'Continuité et limites',
    'الاستمرارية والنهايات',
    'مراجعة شاملة للنهايات وقواعد الاستمرارية، مبرهنة القيم المتوسطة (TVI)، وطريقة إيجاد حلول المعادلات ونقاط الانعطاف في الباك التونسي.',
    1
);

-- 2. Dérivabilité et Théorèmes Fondamentaux
INSERT INTO chapters (subject_id, code, domain, title_fr, title_ar, description_tn, sort_order)
VALUES (
    (SELECT id FROM subjects WHERE code = 'bac-math-math'),
    'derivabilite-theoremes',
    'analyse',
    'Dérivabilité et théorèmes fondamentaux',
    'الاشتقاقية ونظريات الدوال',
    'قواعد الاشتقاق، مبرهنة رول (Théorème de Rolle)، مبرهنة التزايدات المنتهية (TAF) والتطبيقات متاعها في دراسة تغيرات الدوال وحصر المقادير.',
    2
);

-- 3. Fonctions Exponentielles et Logarithmes
INSERT INTO chapters (subject_id, code, domain, title_fr, title_ar, description_tn, sort_order)
VALUES (
    (SELECT id FROM subjects WHERE code = 'bac-math-math'),
    'fonctions-exponentielles-logarithmes',
    'analyse',
    'Fonctions exponentielles et logarithmes',
    'الدوال الأسية واللوغاريتمية',
    'دراسة دالة اللوغاريتم النيبري (ln) والدالة الأسية (exp)، النهايات الشهيرة، وحل المعادلات والمتراجحات اللي تتعاود ديما في امتحانات الباك التونسي.',
    3
);

-- 4. Suites Réelles
INSERT INTO chapters (subject_id, code, domain, title_fr, title_ar, description_tn, sort_order)
VALUES (
    (SELECT id FROM subjects WHERE code = 'bac-math-math'),
    'suites-reelles',
    'analyse',
    'Suites réelles',
    'المتتاليات العددية',
    'المتتاليات الحسابية والهندسية، المتتاليات التراجعية (Suites récurrentes)، البرهان بالتراجع (Récurrence)، وحساب النهايات والتقارب.',
    4
);

-- 5. Primitives et Calcul d'Intégrales
INSERT INTO chapters (subject_id, code, domain, title_fr, title_ar, description_tn, sort_order)
VALUES (
    (SELECT id FROM subjects WHERE code = 'bac-math-math'),
    'primitives-integrales',
    'analyse',
    'Primitives et calcul d''intégrales',
    'الدوال الأصلية وحساب التكامل',
    'جدول الدوال الأصلية، تقنيات التكامل بالتجزئة (Intégration par parties)، وحساب المساحات والحجوم الهندسية والتطبيقات العددية.',
    5
);

-- 6. Équations Différentielles
INSERT INTO chapters (subject_id, code, domain, title_fr, title_ar, description_tn, sort_order)
VALUES (
    (SELECT id FROM subjects WHERE code = 'bac-math-math'),
    'equations-differentielles',
    'analyse',
    'Équations différentielles',
    'المعادلات التفاضلية',
    'حل المعادلات التفاضلية من الرتبة الأولى والرتبة الثانية الخطية ذات المعاملات الثابتة مع الشروط الابتدائية وتطبيقاتها الفيزيائية.',
    6
);

-- 7. Nombres Complexes
INSERT INTO chapters (subject_id, code, domain, title_fr, title_ar, description_tn, sort_order)
VALUES (
    (SELECT id FROM subjects WHERE code = 'bac-math-math'),
    'nombres-complexes',
    'algebre-geometrie',
    'Nombres complexes',
    'الأعداد المركبة',
    'الشكل الجبري والمثلثي والآسي (Forme algébrique, trigonométrique, exponentielle)، حل معادلات الدرجة الثانية في C، والتوظيف الهندسي للأعداد المركبة.',
    7
);

-- 8. Arithmétique dans Z (Section Math)
INSERT INTO chapters (subject_id, code, domain, title_fr, title_ar, description_tn, sort_order)
VALUES (
    (SELECT id FROM subjects WHERE code = 'bac-math-math'),
    'arithmetique-z',
    'algebre-geometrie',
    'Arithmétique dans Z',
    'الحسابيات في Z',
    'خاص بشعبة الرياضيات: القسمة الإقليدية، التوافقات في Z (Congruences)، القاسم المشترك الأكبر PGCD، خوارزمية أقليدس، مبرهنة بوزو ومبرهنة غوص (Bézout et Gauss).',
    8
);

-- 9. Géométrie dans l'Espace
INSERT INTO chapters (subject_id, code, domain, title_fr, title_ar, description_tn, sort_order)
VALUES (
    (SELECT id FROM subjects WHERE code = 'bac-math-math'),
    'geometrie-espace',
    'algebre-geometrie',
    'Géométrie dans l''espace et produit vectoriel',
    'الهندسة في الفضاء والجداء المتجهي',
    'الجداء السلمي والمتجهي (Produit scalaire et vectoriel)، معادلات المستقيمات والمستويات، حساب المسافات، وتقاطع الفلكات والمستويات في الفضاء.',
    9
);

-- 10. Isométries et Similitudes Planes
INSERT INTO chapters (subject_id, code, domain, title_fr, title_ar, description_tn, sort_order)
VALUES (
    (SELECT id FROM subjects WHERE code = 'bac-math-math'),
    'isometries-similitudes',
    'algebre-geometrie',
    'Isométries et similitudes planes',
    'التقايسات والتشابه المستوي',
    'التحويلات النقطية: الانسحاب، الدوران، التناظر المحوري والعمودي، والتشابه المباشر وتحديد عناصره المميزة في المسائل الهندسية الوطنية.',
    10
);

-- 11. Dénombrement et Probabilités Finies
INSERT INTO chapters (subject_id, code, domain, title_fr, title_ar, description_tn, sort_order)
VALUES (
    (SELECT id FROM subjects WHERE code = 'bac-math-math'),
    'denombrement-probabilites',
    'probabilites-statistiques',
    'Dénombrement et probabilités finies',
    'العد والاحتمالات',
    'تقنيات السحب (مع الإرجاع، في آن واحد، بالتتابع دون إرجاع)، الترتيبات والتوفيقات، الاحتمال الشرطي، ومبرهنة الاحتمالات الكلية.',
    11
);

-- 12. Variables Aléatoires et Lois de Probabilité
INSERT INTO chapters (subject_id, code, domain, title_fr, title_ar, description_tn, sort_order)
VALUES (
    (SELECT id FROM subjects WHERE code = 'bac-math-math'),
    'variables-aleatoires-lois',
    'probabilites-statistiques',
    'Variables aléatoires et lois de probabilité',
    'المتغيرات العشوائية وقوانين الاحتمال',
    'قانون الاحتمال، الأمل الرياضي (Espérance)، التباين (Variance)، والقانون الحداني (Loi binomiale) والتطبيقات الإحصائية للباك.',
    12
);
