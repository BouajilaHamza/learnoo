# Expérience locale de narration vidéo — non validée

Statut : prototype technique, non publié et sans validation enseignante. Ce dossier ne constitue ni un cours approuvé, ni une référence officielle du programme tunisien, ni une décision sur les futurs outils pédagogiques.

L'expérience initiale comparait deux narrations (Derja et arabe standard avec termes français) sur des images fixes. Les appréciations du contributeur sur le débit et la prononciation sont des observations exploratoires, pas une étude de compréhension ni une validation du niveau scolaire. Aucune conclusion générale sur la valeur pédagogique d'une langue ou d'une technique n'est établie.

## Corrections et limites

Les deux narrations utilisent désormais la même liste de diapositives dans `scripts/slides.py`. L'exemple numérique de la quatrième diapositive et la conclusion de la cinquième correspondent à la narration dans les deux variantes. Chaque image porte un avertissement de non-validation. Les chemins de sortie sont fournis en argument ; aucun chemin personnel n'est requis.

Les anciens fichiers MP3/MP4 sont retirés de cette proposition : ils précédaient ces corrections et ne prouvent pas le résultat actuel. Ils restent dans l'historique Git. Les scripts actuels doivent être exécutés puis les sorties écoutées et inspectées avant toute nouvelle comparaison. Les identifiants de contenu empêchent de combiner automatiquement des images ou des pistes issues d'une autre version des diapositives.

Le rendu emploie du HTML avec des caractères mathématiques simples et les polices locales. Ce n'est pas une mesure de KaTeX. Aucun composant du service Learnoo ne dépend de ce pipeline ; YouTube reste le mode de diffusion prévu.

## Exécution facultative

Prérequis : Python 3, Chrome/Chromium, `ffmpeg` et l'exécutable `edge-tts` installés séparément. La synthèse vocale transmet les textes à un service externe ; utiliser exclusivement ces textes techniques sans données personnelles. Les versions et les conditions de réutilisation du service et de ses sorties doivent être relevées avant toute diffusion. Les sorties ne sont jamais publiées automatiquement.

Depuis la racine du dépôt, remplacer les chemins de sortie et de navigateur selon la machine :

```sh
# HTML uniquement : hors ligne, sans dépendance Python supplémentaire.
python3 docs/fixtures/video-eval/scripts/generate_lesson_video.py --output-dir /tmp/learnoo-video --html-only

# Images, narration Derja et vidéo locale non validée.
python3 docs/fixtures/video-eval/scripts/generate_lesson_video.py --output-dir /tmp/learnoo-video --chrome google-chrome

# Narration MSA sur les mêmes images, puis assemblage.
python3 docs/fixtures/video-eval/scripts/generate_msa_audio.py --output-dir /tmp/learnoo-video
python3 docs/fixtures/video-eval/scripts/render_msa_video.py --output-dir /tmp/learnoo-video

# Contrôles locaux sans réseau ni outils multimédias.
python3 -m unittest discover -s docs/fixtures/video-eval/scripts -p 'test_*.py'
```

Les vidéos générées se trouvent dans `derja/unreviewed-lesson.mp4` et `msa/unreviewed-lesson.mp4` sous le répertoire choisi. Ne pas les enregistrer dans les assets de l'application. Les tests de scripts vérifient les chemins, l'association des segments et les entrées périmées ; ils ne prouvent ni la synchronisation audiovisuelle réelle ni la qualité pédagogique.

Prochaine étape : exécuter le pipeline dans un environnement documenté, écouter et inspecter chaque segment, puis obtenir une revue enseignante indépendante et clarifier les droits avant toute utilisation pédagogique. L'issue 48 reste une discussion exploratoire, pas une nouvelle fonctionnalité du produit.
