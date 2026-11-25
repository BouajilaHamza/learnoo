# Catalogue pilote — mathématiques

Statut : fondation technique, sans validation pédagogique du programme.

Le périmètre confirmé par le mainteneur est la section Mathématiques et la matière Mathématiques. La migration `0003_curriculum_catalog.sql` crée leurs identifiants stables (`math`, `bac-math-math`) et les tables nécessaires au catalogue. Elle ne renseigne aucun coefficient et ne crée aucun chapitre.

La liste de douze chapitres et les descriptions précédemment proposées n'étaient accompagnées ni de références officielles vérifiables ni d'une validation enseignante. Elles sont retirées de la migration et ne constituent pas un programme approuvé. Aucune autre section n'est ajoutée au pilote.

Avant d'ajouter des chapitres dans une nouvelle migration :

- Identifier la source officielle, sa date ou édition et les pages concernées.
- Faire confirmer le périmètre et l'ordre des chapitres par un enseignant identifié.
- Consigner la validation indépendante et les droits nécessaires pour les descriptions pédagogiques.
- Conserver les identifiants historiques ; les futures révisions publiées devront figer leurs métadonnées de programme.

Les tests utilisent uniquement des chapitres techniques explicitement non validés, insérés dans une base éphémère. Leur réussite ne prouve aucune conformité au programme officiel. Le choix du premier chapitre et les personnes responsables de sa validation restent à confirmer.
